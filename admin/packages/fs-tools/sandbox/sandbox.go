// Package sandbox starts a program kept away from the rest of the system:
// it reads only the system's folders and those it is given (Path), writes
// only those given to write, reaches no network (unless Network) and has
// limits (Limits).
//
// By bubblewrap (ModeBwrap: its namespaces — of the mounts, the processes,
// the network —; it needs the namespaces of users unprivileged, often
// refused: AppArmor, the seccomp of Docker), else by Landlock (ModeLandlock:
// the kernel's restriction of the files, from Linux 5.13, of TCP from 6.7
// (ABI 4) and of the signals from 6.12 (ABI 6) — below it, the program may
// signal the processes of its user, the server's). ModeAuto takes the first
// that works; ModeOff is no sandbox, the limits only.
package sandbox

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// Mode is how the program is kept apart.
type Mode string

const (
	ModeAuto     Mode = "auto"
	ModeBwrap    Mode = "bwrap"
	ModeLandlock Mode = "landlock"
	ModeOff      Mode = "off"
)

// Path is a folder (or a file) the program reaches: read, or Write.
type Path struct {
	Path  string
	Write bool
}

// Limits of the program and those it starts (0: none): the time of CPU,
// the memory (virtual), the size of a file written, the processes of the
// user (the server's included: RLIMIT_NPROC counts them all).
type Limits struct {
	CPU       time.Duration
	Memory    int64
	FileSize  int64
	Processes int
}

// Spec is a program to start.
type Spec struct {
	// Argv[0] is an absolute path.
	Argv []string
	Env  []string
	// Dir is where it starts.
	Dir string
	// Paths it reaches besides the system's (/usr, /etc, …).
	Paths   []Path
	Network bool
	Limits  Limits
	// TTY: its input, its output, its errors; with Setctty, a terminal (its
	// slave) it controls, in a session of its own.
	TTY     *os.File
	Setctty bool
}

// ErrUnavailable: the mode asked does not work here.
var ErrUnavailable = errors.New("sandbox unavailable")

// SystemPaths are what every program reads: its libraries, its
// configuration.
var SystemPaths = []string{"/usr", "/bin", "/sbin", "/lib", "/lib32", "/lib64", "/libx32", "/etc", "/opt"}

// Coverage is what a mode keeps apart here: the files (only the paths
// given), TCP (no connection, no port opened, unless Spec.Network), the
// signals (no process out of the sandbox signaled). ABI is Landlock's
// version (ModeLandlock).
type Coverage struct {
	Mode                Mode
	Files, TCP, Signals bool
	ABI                 int
}

// Isolated: something is kept apart (not ModeOff).
func (c Coverage) Isolated() bool { return c.Files || c.TCP || c.Signals }

// String says it, for the log: "landlock ABI 4: files, TCP (signals not
// isolated: Linux < 6.12)".
func (c Coverage) String() string {
	switch c.Mode {
	case ModeBwrap:
		return "bwrap: namespaces of mounts, processes, network"
	case ModeLandlock:
		s := "landlock ABI " + strconv.Itoa(c.ABI) + ": files"
		if c.TCP {
			s += ", TCP"
		}
		if c.Signals {
			s += ", signals"
		} else {
			s += " (signals not isolated: Linux < 6.12)"
		}
		return s
	}
	return string(c.Mode)
}

// Describe says what m keeps apart here (for the log).
func Describe(m Mode) string { return CoverageOf(m).String() }
