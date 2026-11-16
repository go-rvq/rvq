package fs_tools

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
	"github.com/go-rvq/rvq/admin/packages/fs-tools/fs"
	"github.com/go-rvq/rvq/admin/packages/fs-tools/sandbox"
	"github.com/go-rvq/rvq/x/appconfig"
	"github.com/go-rvq/rvq/x/features"
	"github.com/go-rvq/rvq/x/login"
	"github.com/gorilla/websocket"
)

// TerminalFeature is the terminal of the files' page: on unless app.yml
// turns it off (`features: {rvq.fstools.terminal: false}`), not compiled
// with the tag rvq.fstools.noterminal.
var TerminalFeature = features.Register("rvq.fstools.terminal", "Terminal of the files",
	"a shell in the folders of the file system (Administration > File system), in a sandbox (fstools.terminal of app.yml)").
	CompiledIf(terminalCompiled).DefaultOn()

// TerminalConfig is `fstools: {terminal: …}` of app.yml.
type TerminalConfig struct {
	// Sandbox: auto (bwrap, else Landlock; neither: the terminal is off),
	// bwrap, landlock, off (none: the shell reaches what the server does).
	Sandbox sandbox.Mode `yaml:"sandbox"`
	// Network: the shell reaches the network (TCP).
	Network bool `yaml:"network"`
	// Shell: its program (default bash, else sh).
	Shell  string `yaml:"shell"`
	Limits struct {
		CPU       string `yaml:"cpu"`
		Memory    string `yaml:"memory"`
		FileSize  string `yaml:"file_size"`
		Processes int    `yaml:"processes"`
	} `yaml:"limits"`
	// Idle: with no input for so long, it closes.
	Idle string `yaml:"idle"`
}

// DefaultTerminalConfig is what app.yml does not say.
func DefaultTerminalConfig() (c TerminalConfig) {
	c.Sandbox = sandbox.ModeAuto
	c.Limits.CPU, c.Limits.Memory, c.Limits.FileSize = "10m", "2G", "1G"
	c.Idle = "30m"
	return
}

// terminal is the terminal as configured.
type terminal struct {
	mode   sandbox.Mode
	off    error // why it is off
	shell  string
	net    bool
	limits sandbox.Limits
	idle   time.Duration
}

func loadTerminal() (t *terminal) {
	t = &terminal{}
	if !TerminalFeature.Compiled() || !TerminalFeature.Enabled() {
		t.off = errors.New("feature rvq.fstools.terminal off")
		return
	}
	var cfg struct {
		Terminal TerminalConfig `yaml:"terminal"`
	}
	cfg.Terminal = DefaultTerminalConfig()
	if _, err := appconfig.Section("fstools", &cfg); err != nil {
		t.off = fmt.Errorf("fstools of app.yml: %w", err)
		return
	}
	c := cfg.Terminal
	var errs []error
	parseDur := func(name, v string) time.Duration {
		if v == "" || v == "0" {
			return 0
		}
		d, err := time.ParseDuration(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
		return d
	}
	parseSize := func(name, v string) int64 {
		n, err := ParseSize(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
		return n
	}
	t.limits = sandbox.Limits{
		CPU:       parseDur("limits.cpu", c.Limits.CPU),
		Memory:    parseSize("limits.memory", c.Limits.Memory),
		FileSize:  parseSize("limits.file_size", c.Limits.FileSize),
		Processes: c.Limits.Processes,
	}
	if t.idle = parseDur("idle", c.Idle); t.idle <= 0 {
		t.idle = 30 * time.Minute
	}
	if len(errs) > 0 {
		t.off = fmt.Errorf("fstools.terminal of app.yml: %w", errors.Join(errs...))
		return
	}
	t.net = c.Network
	if t.shell = c.Shell; t.shell == "" {
		t.shell = "/bin/bash"
		if _, err := os.Stat(t.shell); err != nil {
			t.shell = "/bin/sh"
		}
	}
	t.mode, t.off = sandbox.Detect(c.Sandbox)
	return
}

// ParseSize reads "512M", "2G", "1048576" (bytes; K, M, G, T of 1024).
func ParseSize(v string) (int64, error) {
	v = strings.TrimSpace(strings.ToUpper(v))
	if v == "" || v == "0" {
		return 0, nil
	}
	mul := int64(1)
	if i := strings.IndexAny(v, "KMGT"); i > 0 {
		mul = int64(1) << (10 * (1 + strings.IndexByte("KMGT", v[i])))
		if rest := strings.TrimSuffix(v[i+1:], "B"); rest != "" {
			return 0, fmt.Errorf("invalid size %q", v)
		}
		v = v[:i]
	}
	n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("invalid size %q", v)
	}
	return n * mul, nil
}

// terminalPath is where the terminal is: <prefix>/fs-tools/terminal?dir=….
const terminalPath = explorerPath + "terminal"

// terminalOn says r may open the terminal: it is on (feature, sandbox) and
// r has the page's !terminal.
func (b *Builder) terminalOn(r *http.Request) bool {
	return b.term != nil && b.term.off == nil && b.terminalAllowed(r)
}

// terminalAllowed: r has the page's !terminal (a shell: no permissions
// configured, no terminal).
func (b *Builder) terminalAllowed(r *http.Request) bool {
	if b.page == nil {
		return false
	}
	ver := b.page.Page().ActionVerifier(r, "!terminal")
	return ver != nil && ver.Allowed()
}

// terminalOff is why the terminal is off, to one who may open it; nil when
// it is on or r may not.
func (b *Builder) terminalOff(r *http.Request) error {
	if b.term == nil || b.term.off == nil || !b.terminalAllowed(r) {
		return nil
	}
	return b.term.off
}

// mountTerminal serves the terminal (a websocket: binary frames, the bytes;
// text, {"resize":[cols,rows]}; at its end {"exit":"…"}).
func (b *Builder) mountTerminal(prefix string, mux *http.ServeMux) {
	b.term = loadTerminal()
	if b.term.off != nil {
		b.log.Info("terminal off", "reason", b.term.off)
	} else {
		b.log.Info("terminal on", "sandbox", sandbox.Describe(b.term.mode))
	}
	mux.Handle(prefix+terminalPath, b.recordAccess(http.HandlerFunc(b.serveTerminal), AccessTerminal))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  32 << 10,
	WriteBufferSize: 32 << 10,
	// the page is of this host (its port aside: a proxy may drop it)
	CheckOrigin: func(r *http.Request) bool {
		o := r.Header.Get("Origin")
		if o == "" {
			return true
		}
		u, err := url.Parse(o)
		if err != nil {
			return false
		}
		host := r.Host
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		}
		return strings.EqualFold(u.Hostname(), host)
	},
}

// session is the folders of a terminal: a temporary folder holding root
// (FILES_ROOT: a link to each folder of the OS the request reaches, at its
// place in the file system) and tmp (its HOME, its TMPDIR).
type session struct {
	dir, root, tmp string
	paths          []sandbox.Path
}

// newSession is the session of r; none when r reaches no folder of the OS.
func (b *Builder) newSession(r *http.Request) (*session, error) {
	type link struct {
		at, local string
		write     bool
	}
	var links []link
	dirs, err := fs.Dirs(EnvDirs)
	if err != nil {
		return nil, err
	}
	for _, d := range dirs {
		p := "/" + d.Mount
		if !b.allowed(r, p, "", false) {
			continue
		}
		links = append(links, link{d.Mount, d.Local, !d.ReadOnly && b.allowed(r, p, "", true)})
	}
	for _, m := range b.mounts {
		if m.Allowed != nil && !m.Allowed(r, MountAccess{Path: "/"}) {
			continue
		}
		dir, err := m.Dir(r)
		if err != nil {
			continue
		}
		dir, _ = filepath.Abs(dir)
		links = append(links, link{m.Name, dir, m.Allowed == nil || m.Allowed(r, MountAccess{Path: "/", Write: true})})
	}
	if len(links) == 0 {
		return nil, errors.New("no folder of the system here")
	}
	tmp, err := os.MkdirTemp("", "rvq-terminal-")
	if err != nil {
		return nil, err
	}
	s := &session{dir: tmp, root: filepath.Join(tmp, "root"), tmp: filepath.Join(tmp, "tmp")}
	s.paths = []sandbox.Path{{Path: tmp, Write: true}}
	for _, l := range links {
		at := filepath.Join(s.root, filepath.FromSlash(l.at))
		if err = os.MkdirAll(filepath.Dir(at), 0o700); err == nil {
			err = os.Symlink(l.local, at)
		}
		if err != nil {
			s.close()
			return nil, err
		}
		s.paths = append(s.paths, sandbox.Path{Path: l.local, Write: l.write})
	}
	if err = os.MkdirAll(s.tmp, 0o700); err != nil {
		s.close()
		return nil, err
	}
	return s, nil
}

// at is the folder of the OS of p (a path of the file system) in s: its
// root when p is not a folder reached.
func (s *session) at(p string) string {
	d := filepath.Join(s.root, filepath.FromSlash(path.Clean("/"+p)))
	if fi, err := os.Stat(d); err == nil && fi.IsDir() {
		return d
	}
	return s.root
}

func (s *session) close() { _ = os.RemoveAll(s.dir) }

func (b *Builder) serveTerminal(w http.ResponseWriter, r *http.Request) {
	if !b.terminalOn(r) {
		http.Error(w, "terminal not allowed", http.StatusForbidden)
		return
	}
	b.runTerminal(w, r)
}

// runTerminal is the terminal of r, allowed.
func (b *Builder) runTerminal(w http.ResponseWriter, r *http.Request) {
	sess, err := b.newSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	defer sess.close()
	wc, err := upgrader.Upgrade(hijacker{w}, r, nil)
	if err != nil {
		return
	}
	defer wc.Close()
	var wmu sync.Mutex
	write := func(mt int, p []byte) error {
		wmu.Lock()
		defer wmu.Unlock()
		_ = wc.SetWriteDeadline(time.Now().Add(15 * time.Second))
		return wc.WriteMessage(mt, p)
	}
	exit := func(s string) {
		p, _ := json.Marshal(map[string]string{"exit": s})
		_ = write(websocket.TextMessage, p)
	}
	q := r.URL.Query()
	master, slave, err := pty.Open()
	if err != nil {
		exit("could not start: " + err.Error())
		return
	}
	defer master.Close()
	_ = pty.Setsize(master, &pty.Winsize{Cols: uint16(atoiOr(q.Get("cols"), 120)), Rows: uint16(atoiOr(q.Get("rows"), 32))})
	t := b.term
	dir := sess.at(q.Get("dir"))
	env := []string{
		// its folder by the links of FILES_ROOT (the shell keeps it)
		"PWD=" + dir,
		"PATH=/usr/local/bin:/usr/bin:/bin", "TERM=xterm-256color", "COLORTERM=truecolor",
		"HOME=" + sess.tmp, "TMPDIR=" + sess.tmp, "FILES_ROOT=" + sess.root, "SHELL=" + t.shell,
	}
	for _, k := range []string{"LANG", "LC_ALL", "TZ", "USER", "LOGNAME"} {
		if v, ok := os.LookupEnv(k); ok {
			env = append(env, k+"="+v)
		}
	}
	proc, err := sandbox.Start(t.mode, sandbox.Spec{
		Argv: []string{t.shell, "-l"}, Env: env, Dir: dir,
		Paths: sess.paths, Network: t.net, Limits: t.limits, TTY: slave, Setctty: true,
	})
	slave.Close()
	if err != nil {
		exit("could not start: " + err.Error())
		return
	}
	who := whoOf(r)
	b.log.Info("terminal opened", "user", who, "dir", q.Get("dir"), "sandbox", string(t.mode))

	// pty → page
	outDone := make(chan struct{})
	go func() {
		defer close(outDone)
		buf := make([]byte, 32<<10)
		for {
			n, err := master.Read(buf)
			if n > 0 && write(websocket.BinaryMessage, buf[:n]) != nil {
				return
			}
			if err != nil {
				return
			}
		}
	}()
	// page → pty; idle: closed
	inDone := make(chan struct{})
	go func() {
		defer close(inDone)
		for {
			_ = wc.SetReadDeadline(time.Now().Add(t.idle))
			mt, p, err := wc.ReadMessage()
			if err != nil {
				return
			}
			switch mt {
			case websocket.BinaryMessage:
				_, _ = master.Write(p)
			case websocket.TextMessage:
				var c struct {
					Resize []int `json:"resize"`
				}
				if json.Unmarshal(p, &c) == nil && len(c.Resize) == 2 {
					_ = pty.Setsize(master, &pty.Winsize{Cols: uint16(atoiOr(strconv.Itoa(c.Resize[0]), 120)), Rows: uint16(atoiOr(strconv.Itoa(c.Resize[1]), 32))})
				}
			}
		}
	}()
	select {
	case <-proc.Done():
	case <-inDone:
		// the page went away: the shell hung up, killed if it stays
		proc.Signal(syscall.SIGHUP)
		select {
		case <-proc.Done():
		case <-time.After(5 * time.Second):
			proc.Signal(syscall.SIGKILL)
			<-proc.Done()
		}
	}
	// what the session started and left
	proc.Signal(syscall.SIGKILL)
	time.Sleep(100 * time.Millisecond)
	master.Close()
	<-outDone
	status := "0"
	if err := proc.Wait(); err != nil {
		status = err.Error()
	}
	b.log.Info("terminal closed", "user", who, "exit", status)
	exit(status)
}

// hijacker is w hijacked through the writers wrapping it (their Unwrap):
// the log's, the login's.
type hijacker struct{ http.ResponseWriter }

func (w hijacker) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(w.ResponseWriter).Hijack()
}

func atoiOr(s string, def int) int {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 || n > 1000 {
		return def
	}
	return n
}

// whoOf is the account of r's user, for the log.
func whoOf(r *http.Request) string {
	switch u := login.CurrentUserFromRequest(r).(type) {
	case nil:
		return ""
	case interface{ GetAccountName() string }:
		return u.GetAccountName()
	default:
		return fmt.Sprint(u)
	}
}
