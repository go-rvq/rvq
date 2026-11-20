package lifecycle

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

func keys(nodes []Node) string {
	s := make([]string, len(nodes))
	for i, n := range nodes {
		s[i] = string(n.Key)
	}
	return strings.Join(s, ",")
}

// The order: the dependencies first, ties by the key — whatever the order of
// the Adds.
func TestSortedOrder(t *testing.T) {
	var g Graph
	g.Add(
		Node{Key: "http", After: []Key{"worker", "contact"}},
		Node{Key: "worker", After: []Key{"migrate"}},
		Node{Key: "contact", After: []Key{"migrate"}},
		Node{Key: "migrate", After: []Key{"db"}},
		Node{Key: "db"},
	)
	nodes, err := g.Sorted()
	if err != nil {
		t.Fatal(err)
	}
	if got := keys(nodes); got != "db,migrate,contact,worker,http" {
		t.Errorf("order: %s", got)
	}
}

// Before: an edge the other way — the other node does not know it.
func TestSortedBefore(t *testing.T) {
	var g Graph
	g.Add(
		Node{Key: "a"},
		Node{Key: "b", After: []Key{"a"}},
		Node{Key: "x", Before: []Key{"a"}},
	)
	nodes, err := g.Sorted()
	if err != nil {
		t.Fatal(err)
	}
	if got := keys(nodes); got != "x,a,b" {
		t.Errorf("order: %s", got)
	}
}

// A dependency not registered: an error, unless Optional.
func TestSortedMissing(t *testing.T) {
	var g Graph
	g.Add(
		Node{Key: "a", After: []Key{"gone"}},
		Node{Key: "b", After: []Key{"off"}, Before: []Key{"off2"}, Optional: []Key{"off", "off2"}},
	)
	_, err := g.Sorted()
	var me *MissingError
	if !errors.As(err, &me) || me.Node != "a" || me.Dep != "gone" {
		t.Fatalf("missing: %v", err)
	}
	if strings.Contains(err.Error(), "off") {
		t.Errorf("an optional one reported: %v", err)
	}

	var g2 Graph
	g2.Add(Node{Key: "b", After: []Key{"off"}, Optional: []Key{"off"}})
	if nodes, err := g2.Sorted(); err != nil || keys(nodes) != "b" {
		t.Errorf("optional: %v %v", keys(nodes), err)
	}
}

// A cycle: its path, ending where it begins.
func TestSortedCycle(t *testing.T) {
	var g Graph
	g.Add(
		Node{Key: "root"},
		Node{Key: "a", After: []Key{"c", "root"}},
		Node{Key: "b", After: []Key{"a"}},
		Node{Key: "c", After: []Key{"b"}},
	)
	_, err := g.Sorted()
	var ce *CycleError
	if !errors.As(err, &ce) {
		t.Fatalf("cycle: %v", err)
	}
	if got := ce.Error(); got != "lifecycle: a cycle: a → b → c → a" {
		t.Errorf("the path: %s", got)
	}
}

// A key added twice, an Add after the graph ran: a mistake of the code.
func TestAddPanics(t *testing.T) {
	mustPanic := func(name string, f func()) {
		t.Helper()
		defer func() {
			if recover() == nil {
				t.Errorf("%s: no panic", name)
			}
		}()
		f()
	}
	var g Graph
	g.Add(Node{Key: "a"})
	mustPanic("twice", func() { g.Add(Node{Key: "a"}) })
	mustPanic("no key", func() { g.Add(Node{}) })
	if _, err := g.Sorted(); err != nil {
		t.Fatal(err)
	}
	mustPanic("sealed", func() { g.Add(Node{Key: "b"}) })
	if !g.Has("a") || g.Has("b") {
		t.Error("Has")
	}
}

func TestPlan(t *testing.T) {
	var g Graph
	g.Add(Node{Key: "db"}, Node{Key: "worker", After: []Key{"db"}})
	p, err := g.Plan()
	if err != nil {
		t.Fatal(err)
	}
	if want := "start:\n  1. db\n  2. worker (after db)\nclose: the inverse order\n"; p != want {
		t.Errorf("plan:\n%s", p)
	}
}

// log records what the nodes do, in order.
type log struct {
	mu sync.Mutex
	s  []string
}

func (l *log) add(s string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.s = append(l.s, s)
}

func (l *log) String() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return strings.Join(l.s, ",")
}

func node(l *log, key Key, after ...Key) Node {
	return Node{Key: key, After: after, Start: func(context.Context) (io.Closer, error) {
		l.add("start " + string(key))
		return CloserFunc(func() error { l.add("close " + string(key)); return nil }), nil
	}}
}

// Started in order, closed in the inverse one; run between them.
func TestRun(t *testing.T) {
	var (
		l log
		g Graph
	)
	g.Add(node(&l, "b", "a"), node(&l, "a"), Node{Key: "point", After: []Key{"a"}},
		Node{Key: "noclose", After: []Key{"b"}, Start: func(context.Context) (io.Closer, error) {
			l.add("start noclose")
			return nil, nil
		}})
	err := g.Run(context.Background(), func(context.Context) error { l.add("run"); return nil })
	if err != nil {
		t.Fatal(err)
	}
	if got := l.String(); got != "start a,start b,start noclose,run,close b,close a" {
		t.Errorf("run: %s", got)
	}
}

// A Start that fails: its error with its key; the started closed, inverse;
// nothing after it started, run not run.
func TestRunStartFails(t *testing.T) {
	var (
		l    log
		g    Graph
		boom = errors.New("boom")
	)
	g.Add(node(&l, "a"), node(&l, "b", "a"), node(&l, "d", "c"),
		Node{Key: "c", After: []Key{"b"}, Start: func(context.Context) (io.Closer, error) { return nil, boom }})
	err := g.Run(context.Background(), func(context.Context) error { l.add("run"); return nil })
	if !errors.Is(err, boom) || !strings.Contains(err.Error(), "start c") {
		t.Fatalf("err: %v", err)
	}
	if got := l.String(); got != "start a,start b,close b,close a" {
		t.Errorf("run: %s", got)
	}
}

// run's error returned; the nodes closed anyway; the closers' errors joined.
func TestRunErrors(t *testing.T) {
	var (
		g      Graph
		runErr = errors.New("run")
		clErr  = errors.New("close")
	)
	g.Add(Node{Key: "a", Start: func(context.Context) (io.Closer, error) {
		return CloserFunc(func() error { return clErr }), nil
	}})
	err := g.Run(context.Background(), func(context.Context) error { return runErr })
	if !errors.Is(err, runErr) || !errors.Is(err, clErr) || !strings.Contains(err.Error(), "close a") {
		t.Errorf("err: %v", err)
	}
}

type shutdowner struct {
	deadline bool
	closed   bool
}

func (s *shutdowner) Close() error { s.closed = true; return nil }
func (s *shutdowner) Shutdown(ctx context.Context) error {
	_, s.deadline = ctx.Deadline()
	return nil
}

// A Shutdowner: Shutdown, with a deadline, not Close.
func TestRunShutdowner(t *testing.T) {
	var g Graph
	s := &shutdowner{}
	g.Add(Node{Key: "a", Start: func(context.Context) (io.Closer, error) { return s, nil }})
	if err := g.Run(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if !s.deadline || s.closed {
		t.Errorf("shutdown with a deadline: %v; closed: %v", s.deadline, s.closed)
	}
}

// server is a Waiter: it runs until closed (or it fails by itself).
type server struct {
	stop chan struct{}
	once sync.Once
	err  error
}

func newServer() *server { return &server{stop: make(chan struct{})} }

func (s *server) end(err error) {
	s.once.Do(func() {
		s.err = err
		close(s.stop)
	})
}
func (s *server) Close() error { s.end(nil); return nil }
func (s *server) Wait() error  { <-s.stop; return s.err }

// The command lives while its waiter does: ctx done, the nodes closed.
func TestRunWaitsUntilCtxDone(t *testing.T) {
	var (
		l   log
		g   Graph
		srv = newServer()
	)
	g.Add(node(&l, "worker"), Node{Key: "http", After: []Key{"worker"}, Start: func(context.Context) (io.Closer, error) {
		return srv, nil
	}})
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error)
	go func() { done <- g.Run(ctx, nil) }()
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
	if got := l.String(); got != "start worker,close worker" {
		t.Errorf("run: %s", got)
	}
}

// A waiter that ends by itself (the port taken): it ends the command, the
// others closed; its error returned.
func TestRunWaiterEnds(t *testing.T) {
	var (
		l    log
		g    Graph
		srv  = newServer()
		port = errors.New("address already in use")
	)
	g.Add(node(&l, "worker"), Node{Key: "http", After: []Key{"worker"}, Start: func(context.Context) (io.Closer, error) {
		go srv.end(port)
		return srv, nil
	}})
	done := make(chan error)
	go func() { done <- g.Run(context.Background(), nil) }()
	select {
	case err := <-done:
		if !errors.Is(err, port) || !strings.Contains(err.Error(), "http") {
			t.Errorf("err: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("did not end")
	}
	if got := l.String(); got != "start worker,close worker" {
		t.Errorf("run: %s", got)
	}
}

// stubborn is a Waiter its Close does not stop.
type stubborn struct{ stop chan struct{} }

func (s *stubborn) Close() error { return nil }
func (s *stubborn) Wait() error  { <-s.stop; return nil }

// A waiter still running after its close: an error, once the timeout passed.
func TestRunWaiterStillRunning(t *testing.T) {
	g := Graph{CloseTimeout: 20 * time.Millisecond}
	s := &stubborn{stop: make(chan struct{})}
	defer close(s.stop)
	g.Add(Node{Key: "http", Start: func(context.Context) (io.Closer, error) { return s, nil }})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := g.Run(ctx, nil)
	if err == nil || !strings.Contains(err.Error(), "http: still running") {
		t.Errorf("err: %v", err)
	}
}

// Go: the loop stops by its ctx when closed; Close waits for it.
func TestGo(t *testing.T) {
	var (
		g       Graph
		stopped = make(chan struct{})
	)
	g.Add(Node{Key: "poll", Start: func(ctx context.Context) (io.Closer, error) {
		return Go(ctx, func(ctx context.Context) {
			<-ctx.Done()
			close(stopped)
		}), nil
	}})
	if err := g.Run(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	select {
	case <-stopped:
	default:
		t.Error("the loop still running after Run")
	}

	c := Go(context.Background(), func(ctx context.Context) { <-ctx.Done() })
	if err := c.Close(); err != nil {
		t.Error(err)
	}
	// Shutdown past its deadline: the loop that does not stop
	block := make(chan struct{})
	defer close(block)
	c = Go(context.Background(), func(context.Context) { <-block })
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	if err := c.(Shutdowner).Shutdown(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("shutdown: %v", err)
	}
}

// Graphs are their own: two run at once, nothing shared.
func TestRunConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := range 8 {
		wg.Go(func() {
			var (
				l log
				g Graph
			)
			g.Add(node(&l, "a"), node(&l, "b", "a"))
			if err := g.Run(context.Background(), nil); err != nil {
				t.Error(i, err)
			}
			if got := l.String(); got != "start a,start b,close b,close a" {
				t.Error(i, got)
			}
		})
	}
	wg.Wait()
}

// Chain: each after the one offered before it — the order offered.
func TestNodesChain(t *testing.T) {
	var ns Nodes
	ns.Chain(Node{Key: "c"})
	ns.Chain(Node{Key: "a", After: []Key{"x"}})
	ns.Chain(Node{Key: "b"})
	var g Graph
	g.Add(Node{Key: "x"})
	g.Add(ns...)
	nodes, err := g.Sorted()
	if err != nil {
		t.Fatal(err)
	}
	if got := keys(nodes); got != "c,x,a,b" {
		t.Errorf("order: %s", got)
	}
	if !ns.Has("a") || ns.Has("x") {
		t.Error("Has")
	}
}
