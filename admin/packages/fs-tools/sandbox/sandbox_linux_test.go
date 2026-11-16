//go:build linux

package sandbox

import (
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// run runs script by sh in m, its output.
func run(t *testing.T, m Mode, s Spec, script string) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	s.Argv = []string{"/bin/sh", "-c", script}
	s.TTY = w
	p, err := Start(m, s)
	w.Close()
	if err != nil {
		t.Fatal(err)
	}
	out, _ := io.ReadAll(r)
	select {
	case <-p.Done():
	case <-time.After(20 * time.Second):
		t.Fatal("not ended")
	}
	return string(out)
}

func modes(t *testing.T) (out []Mode) {
	for _, m := range []Mode{ModeLandlock, ModeBwrap} {
		if _, err := Detect(m); err == nil {
			out = append(out, m)
		} else {
			t.Log(m, err)
		}
	}
	if len(out) == 0 {
		t.Skip("no sandbox here")
	}
	return
}

func TestSandbox(t *testing.T) {
	for _, m := range modes(t) {
		t.Run(string(m), func(t *testing.T) {
			rw, ro, out := t.TempDir(), t.TempDir(), t.TempDir()
			os.WriteFile(filepath.Join(ro, "a"), []byte("A"), 0o644)
			os.WriteFile(filepath.Join(out, "secret"), []byte("S"), 0o644)
			s := Spec{
				Env:    []string{"PATH=/usr/bin:/bin"},
				Dir:    rw,
				Paths:  []Path{{Path: rw, Write: true}, {Path: ro}},
				Limits: Limits{CPU: 7 * time.Second, FileSize: 1 << 20},
			}
			got := run(t, m, s, `
cat `+ro+`/a && echo " read-ok"
echo x > `+rw+`/w && echo write-ok
echo x > `+ro+`/w 2>/dev/null || echo ro-refused
cat `+out+`/secret 2>/dev/null || echo out-refused
echo "cpu $(ulimit -t)"
pwd`)
			for _, want := range []string{"A read-ok", "write-ok", "ro-refused", "out-refused", "cpu 7", rw} {
				if !strings.Contains(got, want) {
					t.Errorf("%q not in %q", want, got)
				}
			}
			if _, err := os.Stat(filepath.Join(ro, "w")); err == nil {
				t.Error("written in the read-only")
			}
		})
	}
}

func TestSandboxNetwork(t *testing.T) {
	for _, m := range modes(t) {
		if m == ModeLandlock && LandlockABI() < 4 {
			continue
		}
		t.Run(string(m), func(t *testing.T) {
			l, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer l.Close()
			go func() {
				for {
					c, err := l.Accept()
					if err != nil {
						return
					}
					c.Close()
				}
			}()
			py := `python3 -c 'import socket,sys; socket.create_connection(("127.0.0.1", int(sys.argv[1])), 3); print("connected")' ` +
				strings.TrimPrefix(l.Addr().String(), "127.0.0.1:") + ` 2>/dev/null || echo refused`
			if got := run(t, m, Spec{Env: []string{"PATH=/usr/bin:/bin"}, Dir: "/"}, py); !strings.Contains(got, "refused") {
				t.Errorf("network reached: %q", got)
			}
			if m == ModeLandlock {
				if got := run(t, m, Spec{Env: []string{"PATH=/usr/bin:/bin"}, Dir: "/", Network: true}, py); !strings.Contains(got, "connected") {
					t.Errorf("network refused: %q", got)
				}
			}
		})
	}
}
