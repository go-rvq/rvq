package cli

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-rvq/rvq/x/lifecycle"
)

// recorder records what the nodes and the Run do, in order.
type recorder struct {
	mu sync.Mutex
	s  []string
}

func (r *recorder) add(s string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.s = append(r.s, s)
}

func (r *recorder) String() string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return strings.Join(r.s, ",")
}

func (r *recorder) node(key lifecycle.Key, after ...lifecycle.Key) lifecycle.Node {
	return lifecycle.Node{Key: key, After: after, Start: func(context.Context) (io.Closer, error) {
		r.add("start " + string(key))
		return lifecycle.CloserFunc(func() error { r.add("close " + string(key)); return nil }), nil
	}}
}

func parse(t *testing.T, root *Command, out io.Writer, args ...string) *CommandContext {
	t.Helper()
	ctx, err := root.Parse(&CommandContext{name: "app", InputArgs: args, Out: out, Err: out})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

// The Out and Err given are the command's — not os.Stdout and os.Stderr —,
// a sub-command's too; -rvq:cli:out names another.
func TestParseKeepsWriters(t *testing.T) {
	var out bytes.Buffer
	root := (&Command{Name: "app"}).Sub(&Command{Name: "say", Run: func(ctx *CommandContext) error {
		_, err := io.WriteString(ctx.Out, "hi")
		return err
	}})
	ctx := parse(t, root, &out, "say")
	if ctx.Out != &out || ctx.Err != &out {
		t.Fatal("the writers given replaced")
	}
	if err := ctx.Run(); err != nil || out.String() != "hi" {
		t.Errorf("out: %q %v", out.String(), err)
	}

	file := filepath.Join(t.TempDir(), "out")
	if err := parse(t, root, &out, "say", "-rvq:cli:out", file).Run(); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(file); string(b) != "hi" {
		t.Errorf("the file: %q", b)
	}
}

// A command with a Lifecycle: its nodes started in order around its Run,
// closed in the inverse order.
func TestLifecycleRun(t *testing.T) {
	var (
		r recorder
		g lifecycle.Nodes
	)
	g.Add(r.node("worker", "db"), r.node("db"))
	root := &Command{
		Name:      "app",
		Lifecycle: func(ctx *CommandContext) error { ctx.Lifecycle().Add(g...); return nil },
		Run: func(ctx *CommandContext) error {
			if ctx.Context.Err() != nil {
				t.Error("the ctx of the Run done")
			}
			r.add("run")
			return nil
		},
	}
	if err := parse(t, root, io.Discard).Run(); err != nil {
		t.Fatal(err)
	}
	if got := r.String(); got != "start db,start worker,run,close worker,close db" {
		t.Errorf("run: %s", got)
	}
}

// -plan: the order, nothing run.
func TestLifecyclePlan(t *testing.T) {
	var (
		r   recorder
		g   lifecycle.Nodes
		out bytes.Buffer
	)
	g.Add(r.node("worker", "db"), r.node("db"))
	root := &Command{Name: "app", Lifecycle: func(ctx *CommandContext) error { ctx.Lifecycle().Add(g...); return nil },
		Run: func(*CommandContext) error { r.add("run"); return nil }}
	if err := parse(t, root, &out, "-plan").Run(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "1. db\n  2. worker (after db)") || r.String() != "" {
		t.Errorf("plan: %q; run: %s", out.String(), r.String())
	}
}

// A command with its own -plan keeps it; the Lifecycle's is not added.
func TestLifecycleOwnPlan(t *testing.T) {
	var (
		g    lifecycle.Nodes
		plan bool
	)
	root := &Command{
		Name:      "app",
		Lifecycle: func(ctx *CommandContext) error { ctx.Lifecycle().Add(g...); return nil },
		New: func(ctx *CommandContext) error {
			ctx.Flags().BoolVar(&plan, "plan", false, "its own")
			return nil
		},
		Run: func(*CommandContext) error { return nil },
	}
	if err := parse(t, root, io.Discard, "-plan").Run(); err != nil {
		t.Fatal(err)
	}
	if !plan {
		t.Error("its own -plan not set")
	}
}

// A sub-command runs none of its parent's nodes: without a Lifecycle, its
// Run alone; with one, its own graph (empty here).
func TestLifecycleNone(t *testing.T) {
	var (
		r recorder
		g lifecycle.Nodes
	)
	g.Add(r.node("db"))
	root := (&Command{Name: "app", Lifecycle: func(ctx *CommandContext) error { ctx.Lifecycle().Add(g...); return nil }}).
		Sub(&Command{Name: "passwd", Run: func(*CommandContext) error { r.add("passwd"); return nil }}).
		Sub(&Command{Name: "own", Lifecycle: func(*CommandContext) error { return nil },
			Run: func(*CommandContext) error { r.add("own"); return nil }})
	for _, name := range []string{"passwd", "own"} {
		if err := parse(t, root, io.Discard, name).Run(); err != nil {
			t.Fatal(err)
		}
	}
	if got := r.String(); got != "passwd,own" {
		t.Errorf("run: %s", got)
	}
}

// waiter keeps the command alive until closed.
type waiter chan struct{}

func (w waiter) Close() error { close(w); return nil }
func (w waiter) Wait() error  { <-w; return nil }

// A waiter keeps the command alive; the ctx of the command done ends it.
func TestLifecycleWaits(t *testing.T) {
	var g lifecycle.Nodes
	g.Add(lifecycle.Node{Key: "http", Start: func(context.Context) (io.Closer, error) { return make(waiter), nil }})
	root := &Command{Name: "app", Lifecycle: func(ctx *CommandContext) error { ctx.Lifecycle().Add(g...); return nil }}
	c, cancel := context.WithCancel(context.Background())
	ctx, err := root.Parse(&CommandContext{name: "app", Context: c, Out: io.Discard, Err: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error)
	go func() { done <- ctx.Run() }()
	select {
	case err := <-done:
		t.Fatalf("returned with its waiter running: %v", err)
	case <-time.After(50 * time.Millisecond):
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("did not end")
	}
}
