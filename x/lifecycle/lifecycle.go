// Package lifecycle is what a command starts and closes: a graph of nodes,
// each registered by its package with what it depends on, started in the
// order of the graph and closed in the inverse one.
//
// Building an application executes nothing: its packages add nodes (Add) and
// the command that declares the graph runs them (Run) — `app migrate` the
// migrations, `app serve` the services and the HTTP —; a command without one
// (`app password`) builds the application and runs alone.
//
//	g.Add(lifecycle.Node{
//		Key:   worker.Key,
//		After: []lifecycle.Key{migrate.Key},
//		Start: func(ctx context.Context) (io.Closer, error) {
//			w.Listen()
//			return w, nil // a Shutdowner: closed within CloseTimeout
//		},
//	})
//
// # The order
//
// A package knows only the keys of what it depends on (After), exported
// constants of their packages. Before puts a node before others that do not
// know it: an edge the other way. The order is topological, ties by the key:
// always the same. A cycle, a dependency not registered (unless Optional) or
// a key added twice are errors; the first Sorted seals the graph — an Add
// after it panics —, so its registration is over before it runs.
//
// # Run
//
// Run starts the nodes in order, each Closer returned deferred: they close
// in the inverse order (the defers are LIFO) — what depends closes before
// its dependency. A Start that fails returns its error, the defers closing
// the nodes started. A Closer that is a Shutdowner is closed by Shutdown,
// within CloseTimeout; a Waiter (the HTTP server: its Listen) keeps the
// command alive — the first one to end, as the ctx done, ends it.
package lifecycle

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"sync"
	"time"
)

// Key names a node: its package's, an exported constant ("rvq/worker").
type Key string

// Node is what a package registers: its key, its order and its Start.
type Node struct {
	Key Key
	// After are the keys this node starts after: its dependencies.
	After []Key
	// Before are the keys this node starts before, they not knowing it.
	Before []Key
	// Optional are the keys of After and Before that may not be registered
	// (of a feature not compiled); any other one missing is an error.
	Optional []Key

	// Start starts the node; the Closer it returns (nil: nothing to close)
	// is closed when the command ends, in the inverse order. Nil: a node
	// only of order (a point others depend on).
	Start func(ctx context.Context) (io.Closer, error)
}

// Nodes are nodes offered while an application is built — before the command
// that runs them exists —: the command adds them to the graph of its run
// (cli.CommandContext.Lifecycle).
type Nodes []Node

// Add offers the nodes.
func (ns *Nodes) Add(nodes ...Node) {
	*ns = append(*ns, nodes...)
}

// Chain offers n after the node offered last (its After and that one's
// key): the nodes chained run in the order they were offered — the order an
// application is built in, when what it offers depends on what it built
// before.
func (ns *Nodes) Chain(n Node) {
	if l := len(*ns); l > 0 {
		n.After = append(slices.Clone(n.After), (*ns)[l-1].Key)
	}
	ns.Add(n)
}

// Has tells whether a node of key was offered.
func (ns Nodes) Has(key Key) bool {
	for _, n := range ns {
		if n.Key == key {
			return true
		}
	}
	return false
}

// Shutdowner is a Closer closed within a deadline: Run prefers it to Close.
// *http.Server is one.
type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

// Waiter is a Closer that keeps the command alive while it runs (an HTTP
// server): Run waits for the waiters, the first one to end ending the
// command.
type Waiter interface {
	Wait() error
}

// DefaultCloseTimeout is the CloseTimeout of a Graph that has none.
var DefaultCloseTimeout = 30 * time.Second

// Graph is the nodes of a command. The zero value is ready to use.
type Graph struct {
	// CloseTimeout is the deadline of each Shutdowner and of the waiters to
	// end after the closers (zero: DefaultCloseTimeout).
	CloseTimeout time.Duration

	mu     sync.Mutex
	nodes  []Node
	keys   map[Key]bool
	sealed bool
	sorted []Node
	err    error
}

// Add registers the nodes. A key empty or added twice, or an Add after the
// graph is sealed (Sorted), panics: a mistake of the code.
func (g *Graph) Add(nodes ...Node) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, n := range nodes {
		if g.sealed {
			panic(fmt.Sprintf("lifecycle: add %s: the graph is sealed (it ran)", n.Key))
		}
		if n.Key == "" {
			panic("lifecycle: add: a node without a key")
		}
		if g.keys[n.Key] {
			panic(fmt.Sprintf("lifecycle: add %s: added twice", n.Key))
		}
		if g.keys == nil {
			g.keys = map[Key]bool{}
		}
		g.keys[n.Key] = true
		g.nodes = append(g.nodes, n)
	}
}

// Has tells whether key was added.
func (g *Graph) Has(key Key) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.keys[key]
}

// Nodes are the nodes added, in the order of the Adds (a copy): to compose
// the graph of a run from it and others. It does not seal the graph.
func (g *Graph) Nodes() []Node {
	g.mu.Lock()
	defer g.mu.Unlock()
	return slices.Clone(g.nodes)
}

// MissingError is a dependency (After or Before) of a node not registered
// and not Optional.
type MissingError struct {
	Node, Dep Key
}

func (e *MissingError) Error() string {
	return fmt.Sprintf("lifecycle: %s: depends on %s, not registered", e.Node, e.Dep)
}

// CycleError is a cycle of the graph: Path ends where it begins.
type CycleError struct {
	Path []Key
}

func (e *CycleError) Error() string {
	s := make([]string, len(e.Path))
	for i, k := range e.Path {
		s[i] = string(k)
	}
	return "lifecycle: a cycle: " + strings.Join(s, " → ")
}

// Sorted is the nodes in the order they start, ties by the key; it seals the
// graph. Its errors: a MissingError for each dependency missing, a
// CycleError.
func (g *Graph) Sorted() ([]Node, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.sealed {
		g.sealed = true
		g.sorted, g.err = sortNodes(g.nodes, g.keys)
	}
	return g.sorted, g.err
}

func sortNodes(nodes []Node, keys map[Key]bool) ([]Node, error) {
	byKey := make(map[Key]Node, len(nodes))
	for _, n := range nodes {
		byKey[n.Key] = n
	}
	// next[a] are the nodes that start after a
	next := map[Key][]Key{}
	in := map[Key]int{}
	var errs []error
	edge := func(from, to Key) {
		next[from] = append(next[from], to)
		in[to]++
	}
	for _, n := range nodes {
		for _, deps := range [2][]Key{n.After, n.Before} {
			for _, d := range deps {
				if !keys[d] && !slices.Contains(n.Optional, d) {
					errs = append(errs, &MissingError{Node: n.Key, Dep: d})
				}
			}
		}
		for _, a := range n.After {
			if keys[a] {
				edge(a, n.Key)
			}
		}
		for _, b := range n.Before {
			if keys[b] {
				edge(n.Key, b)
			}
		}
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	// Kahn: the ready ones by the key
	var ready []Key
	for _, n := range nodes {
		if in[n.Key] == 0 {
			ready = append(ready, n.Key)
		}
	}
	sorted := make([]Node, 0, len(nodes))
	for len(ready) > 0 {
		slices.Sort(ready)
		k := ready[0]
		ready = ready[1:]
		sorted = append(sorted, byKey[k])
		for _, m := range next[k] {
			if in[m]--; in[m] == 0 {
				ready = append(ready, m)
			}
		}
	}
	if len(sorted) < len(nodes) {
		return nil, &CycleError{Path: findCycle(next, in)}
	}
	return sorted, nil
}

// findCycle is a cycle among the nodes left (in > 0): from the least key, a
// depth-first walk until a node on its path comes again.
func findCycle(next map[Key][]Key, in map[Key]int) []Key {
	var left []Key
	for k, n := range in {
		if n > 0 {
			left = append(left, k)
		}
	}
	slices.Sort(left)
	var (
		path  []Key
		onIdx = map[Key]int{}
		seen  = map[Key]bool{}
		walk  func(k Key) []Key
	)
	walk = func(k Key) []Key {
		if i, ok := onIdx[k]; ok {
			return append(slices.Clone(path[i:]), k)
		}
		if seen[k] {
			return nil
		}
		seen[k] = true
		onIdx[k] = len(path)
		path = append(path, k)
		nx := slices.Clone(next[k])
		slices.Sort(nx)
		for _, m := range nx {
			if in[m] > 0 {
				if c := walk(m); c != nil {
					return c
				}
			}
		}
		path = path[:len(path)-1]
		delete(onIdx, k)
		return nil
	}
	for _, k := range left {
		if c := walk(k); c != nil {
			return c
		}
	}
	return left
}

// Plan is the order of the graph, readable: the nodes as they start (with
// what they come after), closed the other way round.
func (g *Graph) Plan() (string, error) {
	nodes, err := g.Sorted()
	if err != nil {
		return "", err
	}
	var b strings.Builder
	b.WriteString("start:\n")
	for i, n := range nodes {
		fmt.Fprintf(&b, "  %d. %s", i+1, n.Key)
		if len(n.After) > 0 {
			fmt.Fprintf(&b, " (after %s)", joinKeys(n.After))
		}
		if len(n.Before) > 0 {
			fmt.Fprintf(&b, " (before %s)", joinKeys(n.Before))
		}
		b.WriteString("\n")
	}
	b.WriteString("close: the inverse order\n")
	return b.String(), nil
}

func joinKeys(keys []Key) string {
	s := make([]string, len(keys))
	for i, k := range keys {
		s[i] = string(k)
	}
	return strings.Join(s, ", ")
}

func (g *Graph) closeTimeout() time.Duration {
	if g.CloseTimeout > 0 {
		return g.CloseTimeout
	}
	return DefaultCloseTimeout
}

// Run starts the nodes in order, then run (nil: none), then waits for the
// waiters; it ends — the closers deferred, in the inverse order — when run
// returned and there are no waiters, when the first waiter ends, or when ctx
// is done. A Start or run that fails returns its error, the nodes started
// closed. Its error joins the closers' and the waiters', each with its key.
func (g *Graph) Run(ctx context.Context, run func(ctx context.Context) error) (err error) {
	nodes, err := g.Sorted()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ws := &waiters{timeout: g.closeTimeout()}
	// the last of the defers: after every closer, the waiters ended
	defer func() { err = errors.Join(err, ws.join()) }()

	for _, n := range nodes {
		if n.Start == nil {
			continue
		}
		c, e := n.Start(ctx)
		if e != nil {
			return fmt.Errorf("lifecycle: start %s: %w", n.Key, e)
		}
		if c == nil {
			continue
		}
		defer func() { err = errors.Join(err, g.close(ctx, n.Key, c)) }()
		if w, ok := c.(Waiter); ok {
			ws.add(n.Key, w, cancel)
		}
	}
	if run != nil {
		if err = run(ctx); err != nil {
			return err
		}
	}
	if ws.n > 0 {
		<-ctx.Done()
	}
	return nil
}

// close closes c: by Shutdown within the timeout when it is a Shutdowner.
func (g *Graph) close(ctx context.Context, key Key, c io.Closer) (err error) {
	if s, ok := c.(Shutdowner); ok {
		sctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), g.closeTimeout())
		defer cancel()
		err = s.Shutdown(sctx)
	} else {
		err = c.Close()
	}
	if err != nil {
		return fmt.Errorf("lifecycle: close %s: %w", key, err)
	}
	return nil
}

// waiters are the Waiters of a Run: each in a goroutine, the first to end
// canceling the Run's ctx.
type waiters struct {
	timeout time.Duration
	wg      sync.WaitGroup
	n       int
	mu      sync.Mutex
	errs    []error
	running map[Key]bool
}

func (ws *waiters) add(key Key, w Waiter, cancel context.CancelFunc) {
	ws.n++
	ws.mu.Lock()
	if ws.running == nil {
		ws.running = map[Key]bool{}
	}
	ws.running[key] = true
	ws.mu.Unlock()
	ws.wg.Go(func() {
		err := w.Wait()
		ws.mu.Lock()
		delete(ws.running, key)
		if err != nil {
			ws.errs = append(ws.errs, fmt.Errorf("lifecycle: %s: %w", key, err))
		}
		ws.mu.Unlock()
		cancel()
	})
}

// join waits for the waiters to end, within the timeout: their errors, and
// one for each still running.
func (ws *waiters) join() error {
	if ws.n == 0 {
		return nil
	}
	done := make(chan struct{})
	go func() {
		ws.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(ws.timeout):
	}
	ws.mu.Lock()
	defer ws.mu.Unlock()
	errs := slices.Clone(ws.errs)
	var still []Key
	for k := range ws.running {
		still = append(still, k)
	}
	slices.Sort(still)
	for _, k := range still {
		errs = append(errs, fmt.Errorf("lifecycle: %s: still running %s after its close", k, ws.timeout))
	}
	return errors.Join(errs...)
}

// CloserFunc is a function as an io.Closer.
type CloserFunc func() error

func (f CloserFunc) Close() error { return f() }

// Go runs f in a goroutine with a ctx of its own; the Closer cancels that
// ctx and waits for f to return — within the deadline, as a Shutdowner. A
// loop in the background (a polling, a ticker): not a Waiter, it does not
// keep the command alive.
func Go(ctx context.Context, f func(ctx context.Context)) io.Closer {
	ctx, cancel := context.WithCancel(ctx)
	g := &goCloser{cancel: cancel, done: make(chan struct{})}
	go func() {
		defer close(g.done)
		f(ctx)
	}()
	return g
}

type goCloser struct {
	cancel context.CancelFunc
	done   chan struct{}
}

func (g *goCloser) Close() error {
	g.cancel()
	<-g.done
	return nil
}

func (g *goCloser) Shutdown(ctx context.Context) error {
	g.cancel()
	select {
	case <-g.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
