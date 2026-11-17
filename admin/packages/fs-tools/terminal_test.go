package fs_tools

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/packages/fs-tools/sandbox"
	"github.com/go-rvq/rvq/x/appconfig"
	"github.com/gorilla/websocket"
)

func TestParseSize(t *testing.T) {
	for v, want := range map[string]int64{"": 0, "10": 10, "2K": 2048, "512M": 512 << 20, "1g": 1 << 30, "3GB": 3 << 30} {
		if n, err := ParseSize(v); err != nil || n != want {
			t.Errorf("%q: %v %v", v, n, err)
		}
	}
	for _, v := range []string{"M", "1X", "1MM", "-1"} {
		if _, err := ParseSize(v); err == nil {
			t.Errorf("%q accepted", v)
		}
	}
}

func TestTerminalConfig(t *testing.T) {
	appconfig.Reset()
	defer appconfig.Reset()
	appconfig.Load([]byte("fstools:\n  terminal:\n    sandbox: off\n    limits: {cpu: 1m}\n    idle: 5m\n"))
	tm := loadTerminal()
	if tm.off != nil || tm.mode != sandbox.ModeOff || tm.limits.CPU != time.Minute || tm.idle != 5*time.Minute || tm.limits.Memory != 2<<30 {
		t.Errorf("%+v", tm)
	}
	appconfig.Load([]byte("fstools:\n  terminal: {sandbox: nope}\n"))
	if loadTerminal().off == nil {
		t.Error("unknown sandbox accepted")
	}
	appconfig.Load([]byte("fstools:\n  terminal: {sandbox: off, limits: {memory: lots}}\n"))
	if loadTerminal().off == nil {
		t.Error("bad memory accepted")
	}
}

// TestTerminal: a shell in the sandbox, by the websocket — the folder asked,
// FILES_ROOT, what it reaches.
func TestTerminal(t *testing.T) {
	mode, err := sandbox.Detect(sandbox.ModeAuto)
	if err != nil {
		t.Skip(err)
	}
	data, outside := t.TempDir(), t.TempDir()
	os.MkdirAll(filepath.Join(data, "sub"), 0o755)
	os.WriteFile(filepath.Join(outside, "secret"), []byte("S"), 0o644)
	old := EnvDirs
	EnvDirs = "w,/data," + data
	defer func() { EnvDirs = old }()

	b := &Builder{log: slog.Default(), term: &terminal{mode: mode, shell: "/bin/sh", idle: time.Minute}}
	// behind a writer that only unwraps (the log's)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { b.runTerminal(unwrapOnly{w}, r) }))
	defer srv.Close()
	c, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http")+"/?dir=/data/sub&cols=80&rows=24", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	script := `echo "at $(pwd)"; echo "root $FILES_ROOT"; echo hi > $FILES_ROOT/data/w && echo written; ` +
		`cat ` + outside + `/secret 2>/dev/null || echo refused; exit 3` + "\n"
	if err := c.WriteMessage(websocket.BinaryMessage, []byte(script)); err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	var exit string
	var info *termInfo
	_ = c.SetReadDeadline(time.Now().Add(20 * time.Second))
	for exit == "" {
		mt, p, err := c.ReadMessage()
		if err != nil {
			t.Fatal(err, out.String())
		}
		if mt == websocket.TextMessage {
			var m struct {
				Exit string
				Info *termInfo
			}
			_ = json.Unmarshal(p, &m)
			if m.Info != nil {
				// the first message, before the shell's output
				if out.Len() > 0 || info != nil {
					t.Error("info not first")
				}
				info = m.Info
			}
			exit = m.Exit
			continue
		}
		out.Write(p)
	}
	if info == nil || info.Label == "" || len(info.Details) < 2 {
		t.Fatalf("no info: %+v", info)
	}
	if got := info.Details[1].Values; len(got) != 1 || got[0].Text != "$FILES_ROOT/data — "+Messages_en_US.ReadWrite {
		t.Errorf("folders %+v", got)
	}
	got := out.String()
	for _, want := range []string{"/root/data/sub", "root /", "written", "refused"} {
		if !strings.Contains(got, want) {
			t.Errorf("%q not in %q", want, got)
		}
	}
	if !strings.Contains(exit, "3") {
		t.Errorf("exit %q", exit)
	}
	if b, err := os.ReadFile(filepath.Join(data, "w")); err != nil || string(b) != "hi\n" {
		t.Error("not written", err)
	}
	// the session's folder gone (before the end was told)
	if m, _ := filepath.Glob(filepath.Join(os.TempDir(), "rvq-terminal-*", "root", "data")); len(m) > 0 {
		for _, p := range m {
			if l, _ := os.Readlink(p); l == data {
				t.Error("left", p)
			}
		}
	}
}

type unwrapOnly struct{ w http.ResponseWriter }

func (u unwrapOnly) Header() http.Header         { return u.w.Header() }
func (u unwrapOnly) Write(p []byte) (int, error) { return u.w.Write(p) }
func (u unwrapOnly) WriteHeader(code int)        { u.w.WriteHeader(code) }
func (u unwrapOnly) Unwrap() http.ResponseWriter { return u.w }

func TestTermInfo(t *testing.T) {
	m := Messages_en_US
	s := &session{folders: []sessionFolder{{At: "files", Write: true}, {At: "docs"}}}
	off := (&terminal{mode: sandbox.ModeOff, idle: 30 * time.Minute}).info(m, s)
	if off.Label != m.TermNotIsolated || off.Color != "warning" || !off.Details[0].Values[0].Warn {
		t.Errorf("off: %+v", off)
	}
	if v := off.Details[1].Values; len(v) != 2 || v[1].Text != "$FILES_ROOT/docs — "+m.ReadOnly {
		t.Errorf("folders: %+v", v)
	}
	tm := &terminal{mode: sandbox.ModeLandlock, idle: 30 * time.Minute,
		limits: sandbox.Limits{CPU: 10 * time.Minute, Memory: 2 << 30, FileSize: 1536 << 20}}
	in := tm.info(m, s)
	if in.Label != "Isolated · Landlock" || in.Color != "info" {
		t.Errorf("landlock: %+v", in)
	}
	var limits, idle string
	for _, d := range in.Details {
		switch d.Label {
		case m.TermLimits:
			limits = d.Values[0].Text
		case m.TermIdle:
			idle = d.Values[0].Text
		}
	}
	if limits != "CPU 10 min · memory 2 GiB · file 1.5 GiB" || idle != "closes after 30 min" {
		t.Errorf("%q %q", limits, idle)
	}
	// the network: on, or off and who turns it on
	network := func(in termInfo) string {
		for _, d := range in.Details {
			if d.Label == m.TermNetwork {
				return d.Values[0].Text
			}
		}
		return ""
	}
	if got := network(in); got != m.TermNetOff {
		t.Errorf("network off: %q", got)
	}
	tm.net = true
	if got := network(tm.info(m, s)); got != m.TermNetOn {
		t.Errorf("network on: %q", got)
	}
}

func TestFmt(t *testing.T) {
	for d, want := range map[time.Duration]string{45 * time.Second: "45 s", 10 * time.Minute: "10 min", time.Hour: "1 h", 90 * time.Minute: "1 h 30 min"} {
		if got := fmtDuration(d); got != want {
			t.Errorf("%v: %q", d, got)
		}
	}
	for n, want := range map[int64]string{512: "512 B", 512 << 20: "512 MiB", 2 << 30: "2 GiB", 1536 << 20: "1.5 GiB"} {
		if got := fmtSize(n); got != want {
			t.Errorf("%d: %q", n, got)
		}
	}
}
