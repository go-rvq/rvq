//go:build linux

package sandbox

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Detect is the mode m is here (ModeAuto: the first that works), or why
// none is (ErrUnavailable).
func Detect(m Mode) (Mode, error) {
	switch m {
	case ModeOff:
		return m, nil
	case ModeBwrap:
		if bwrapWorks() {
			return m, nil
		}
		return "", fmt.Errorf("%w: bwrap does not run (namespaces of users refused?)", ErrUnavailable)
	case ModeLandlock:
		if LandlockABI() > 0 {
			return m, nil
		}
		return "", fmt.Errorf("%w: Landlock is not in the kernel", ErrUnavailable)
	case ModeAuto, "":
		if bwrapWorks() {
			return ModeBwrap, nil
		}
		if LandlockABI() > 0 {
			return ModeLandlock, nil
		}
		return "", fmt.Errorf("%w: neither bwrap nor Landlock works here", ErrUnavailable)
	}
	return "", fmt.Errorf("%w: unknown mode %q", ErrUnavailable, m)
}

var (
	bwrapOnce sync.Once
	bwrapPath string
)

// bwrapWorks: bwrap is found and starts its namespaces.
func bwrapWorks() bool {
	bwrapOnce.Do(func() {
		p, err := exec.LookPath("bwrap")
		if err != nil {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if exec.CommandContext(ctx, p, "--ro-bind", "/", "/", "--unshare-all", "--die-with-parent", "--", "/bin/true").Run() == nil {
			bwrapPath = p
		}
	})
	return bwrapPath != ""
}

var (
	abiOnce sync.Once
	abi     int
)

// LandlockABI is the version of Landlock of the kernel, 0 without it.
func LandlockABI() int {
	abiOnce.Do(func() {
		v, _, e := unix.Syscall(unix.SYS_LANDLOCK_CREATE_RULESET, 0, 0, createRulesetVersion)
		if e == 0 {
			abi = int(v)
		}
	})
	return abi
}

// Proc is a program started.
type Proc struct {
	Pid  int
	done chan struct{}
	err  error
}

// Wait waits for its end.
func (p *Proc) Wait() error {
	<-p.done
	return p.err
}

// Done is closed at its end.
func (p *Proc) Done() <-chan struct{} { return p.done }

// Signal signals its group (its session's processes).
func (p *Proc) Signal(sig syscall.Signal) {
	_ = syscall.Kill(-p.Pid, sig)
}

// gate holds the program until its limits are set: sh reads a line of fd
// 3, then becomes it.
const gate = `read _ <&3; exec 3<&-; exec "$@"`

// Start starts s in mode m (one Detect gave).
func Start(m Mode, s Spec) (*Proc, error) {
	if len(s.Argv) == 0 {
		return nil, errors.New("sandbox: no program")
	}
	argv := s.Argv
	var restrict func() error
	switch m {
	case ModeBwrap:
		if !bwrapWorks() {
			return nil, fmt.Errorf("%w: bwrap", ErrUnavailable)
		}
		argv = append(bwrapArgs(s), argv...)
	case ModeLandlock:
		restrict = func() error { return landlock(s) }
	case ModeOff:
	default:
		return nil, fmt.Errorf("%w: mode %q", ErrUnavailable, m)
	}
	gr, gw, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	cmd := exec.Command("/bin/sh", append([]string{"-c", gate, "sh"}, argv...)...)
	cmd.Env, cmd.Dir = s.Env, s.Dir
	cmd.Stdin, cmd.Stdout, cmd.Stderr = s.TTY, s.TTY, s.TTY
	cmd.ExtraFiles = []*os.File{gr}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true, Setctty: s.Setctty, Pdeathsig: syscall.SIGKILL}

	p := &Proc{done: make(chan struct{})}
	started := make(chan error, 1)
	go func() {
		// the thread of the restriction (Landlock, no_new_privs: of a
		// thread, inherited by its children) is never given back: it ends
		// with the goroutine; Pdeathsig is of it too, so it waits here
		runtime.LockOSThread()
		if restrict != nil {
			if err := restrict(); err != nil {
				started <- err
				return
			}
		}
		if err := cmd.Start(); err != nil {
			started <- err
			return
		}
		started <- nil
		p.err = cmd.Wait()
		close(p.done)
	}()
	if err := <-started; err != nil {
		gw.Close()
		return nil, err
	}
	p.Pid = cmd.Process.Pid
	err = setLimits(p.Pid, s.Limits)
	if err == nil {
		_, err = gw.Write([]byte("go\n"))
	}
	gw.Close()
	if err != nil {
		p.Signal(syscall.SIGKILL)
		<-p.done
		return nil, err
	}
	return p, nil
}

func setLimits(pid int, l Limits) error {
	set := func(res int, v uint64) error {
		if v == 0 {
			return nil
		}
		return unix.Prlimit(pid, res, &unix.Rlimit{Cur: v, Max: v}, nil)
	}
	return errors.Join(
		set(unix.RLIMIT_CPU, uint64(l.CPU/time.Second)),
		set(unix.RLIMIT_AS, uint64(l.Memory)),
		set(unix.RLIMIT_FSIZE, uint64(l.FileSize)),
		set(unix.RLIMIT_NPROC, uint64(l.Processes)),
		set(unix.RLIMIT_CORE, 0),
	)
}

// bwrapArgs: the system read, the paths of s, a /dev, /proc and /tmp of its
// own, no namespace shared but the network (Network).
func bwrapArgs(s Spec) []string {
	a := []string{bwrapPath, "--die-with-parent", "--unshare-all"}
	if s.Network {
		a = append(a, "--share-net")
	}
	for _, p := range SystemPaths {
		a = append(a, "--ro-bind-try", p, p)
	}
	if s.Network {
		a = append(a, "--ro-bind-try", "/run/systemd/resolve", "/run/systemd/resolve")
	}
	a = append(a, "--dev", "/dev", "--proc", "/proc", "--tmpfs", "/tmp")
	for _, p := range s.Paths {
		if p.Write {
			a = append(a, "--bind", p.Path, p.Path)
		} else {
			a = append(a, "--ro-bind", p.Path, p.Path)
		}
	}
	if s.Dir != "" {
		a = append(a, "--chdir", s.Dir)
	}
	return append(a, "--")
}

// The constants of Landlock (linux/landlock.h): of the kernel, the same in
// every version of x/sys.
const (
	fsExecute               = 0x1
	fsWriteFile             = 0x2
	fsReadFile              = 0x4
	fsReadDir               = 0x8
	fsRefer                 = 0x2000
	fsTruncate              = 0x4000
	fsIoctlDev              = 0x8000
	netBindTCP              = 0x1
	netConnectTCP           = 0x2
	createRulesetVersion    = 0x1
	rulePathBeneath         = 1
	scopeAbstractUnixSocket = 0x1
	scopeSignal             = 0x2
)

const (
	fsRead = fsExecute | fsReadFile | fsReadDir
	// fsFile: the rights of a file (the others are of folders)
	fsFile = fsExecute | fsWriteFile | fsReadFile |
		fsTruncate | fsIoctlDev
	fsDevice = fsReadFile | fsWriteFile |
		fsTruncate | fsIoctlDev
)

// devices are the files of /dev a shell uses.
var devices = []string{"/dev/null", "/dev/zero", "/dev/full", "/dev/random", "/dev/urandom", "/dev/tty", "/dev/pts", "/dev/ptmx"}

// fsHandled is what Landlock of the version v restricts of the files.
func fsHandled(v int) uint64 {
	h := uint64(1<<13 - 1)
	if v >= 2 {
		h |= fsRefer
	}
	if v >= 3 {
		h |= fsTruncate
	}
	if v >= 5 {
		h |= fsIoctlDev
	}
	return h
}

// rulesetAttr is struct landlock_ruleset_attr (its fields by version: fs,
// net from 4, scoped from 6); pathBeneathAttr, landlock_path_beneath_attr
// (packed: 12 bytes). Of the kernel, not of x/sys, whose names changed.
type rulesetAttr struct{ fs, net, scoped uint64 }

type pathBeneathAttr struct {
	allowed uint64
	parent  int32
}

// landlock restricts the calling thread (and what it starts) to s.
func landlock(s Spec) error {
	v := LandlockABI()
	if v == 0 {
		return fmt.Errorf("%w: Landlock", ErrUnavailable)
	}
	handled := fsHandled(v)
	attr := rulesetAttr{fs: handled}
	size := unsafe.Sizeof(attr.fs)
	if v >= 4 {
		size += unsafe.Sizeof(attr.net)
		if !s.Network {
			attr.net = netBindTCP | netConnectTCP
		}
	}
	if v >= 6 {
		size += unsafe.Sizeof(attr.scoped)
		attr.scoped = scopeAbstractUnixSocket | scopeSignal
	}
	fd, _, e := unix.Syscall(unix.SYS_LANDLOCK_CREATE_RULESET, uintptr(unsafe.Pointer(&attr)), size, 0)
	if e != 0 {
		return fmt.Errorf("landlock_create_ruleset: %w", e)
	}
	defer unix.Close(int(fd))
	add := func(p string, access uint64) error {
		f, err := unix.Open(p, unix.O_PATH|unix.O_CLOEXEC, 0)
		if err != nil {
			if errors.Is(err, unix.ENOENT) {
				return nil
			}
			return fmt.Errorf("%s: %w", p, err)
		}
		defer unix.Close(f)
		var st unix.Stat_t
		if err := unix.Fstat(f, &st); err != nil {
			return err
		}
		if st.Mode&unix.S_IFMT != unix.S_IFDIR {
			access &= fsFile
		}
		rule := pathBeneathAttr{allowed: access & handled, parent: int32(f)}
		if _, _, e := unix.Syscall6(unix.SYS_LANDLOCK_ADD_RULE, fd, rulePathBeneath,
			uintptr(unsafe.Pointer(&rule)), 0, 0, 0); e != 0 {
			return fmt.Errorf("landlock_add_rule %s: %w", p, e)
		}
		return nil
	}
	for _, p := range SystemPaths {
		if err := add(p, fsRead); err != nil {
			return err
		}
	}
	if s.Network {
		if err := add("/run/systemd/resolve", fsRead); err != nil {
			return err
		}
	}
	for _, p := range devices {
		if err := add(p, fsDevice); err != nil {
			return err
		}
	}
	for _, p := range s.Paths {
		access := uint64(fsRead)
		if p.Write {
			access = handled
		}
		if err := add(p.Path, access); err != nil {
			return err
		}
	}
	if err := unix.Prctl(unix.PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0); err != nil {
		return fmt.Errorf("no_new_privs: %w", err)
	}
	if _, _, e := unix.Syscall(unix.SYS_LANDLOCK_RESTRICT_SELF, fd, 0, 0); e != 0 {
		return fmt.Errorf("landlock_restrict_self: %w", e)
	}
	return nil
}

// Describe says what m keeps apart here (for the log).
func Describe(m Mode) string {
	switch m {
	case ModeLandlock:
		v := LandlockABI()
		s := "landlock ABI " + strconv.Itoa(v) + ": files"
		if v >= 4 {
			s += ", TCP"
		}
		if v >= 6 {
			s += ", signals"
		} else {
			s += " (signals not isolated: Linux < 6.12)"
		}
		return s
	case ModeBwrap:
		return "bwrap: namespaces of mounts, processes, network"
	}
	return string(m)
}
