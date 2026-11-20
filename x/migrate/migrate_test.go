package migrate

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/cli"
	"github.com/go-rvq/rvq/x/lifecycle"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

func dummyDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(tests.DummyDialector{}, &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

// recorder registers migrations that record their key, and the database
// they were given.
type recorder struct {
	keys []string
	dbs  []*gorm.DB
}

func (r *recorder) node(key lifecycle.Key, after ...lifecycle.Key) lifecycle.Node {
	return lifecycle.Node{Key: key, After: after, Start: Func(func(ctx context.Context, db *gorm.DB) error {
		r.keys = append(r.keys, string(key))
		r.dbs = append(r.dbs, db)
		return nil
	})}
}

// The migrations in the order of the graph, each with a session of the
// database of the ctx.
func TestRun(t *testing.T) {
	var (
		g  lifecycle.Graph
		r  recorder
		db = dummyDB(t)
	)
	g.Add(r.node("worker", "i18n"), r.node("i18n"), r.node("media", "i18n", "worker"))
	if err := run(context.Background(), &g, db, nil); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(r.keys, ","); got != "i18n,worker,media" {
		t.Errorf("order: %s", got)
	}
	for _, d := range r.dbs {
		if d == nil || d == db || d.Dialector != db.Dialector {
			t.Errorf("not a session of the database: %p", d)
		}
	}
}

// No database: ErrNoDB, from Run or from a migration run without one.
func TestNoDB(t *testing.T) {
	var g lifecycle.Graph
	if err := run(context.Background(), &g, nil, nil); !errors.Is(err, ErrNoDB) {
		t.Errorf("Run: %v", err)
	}
	start := Func(func(context.Context, *gorm.DB) error { return nil })
	if _, err := start(context.Background()); !errors.Is(err, ErrNoDB) {
		t.Errorf("Func: %v", err)
	}
	if DB(context.Background()) != nil {
		t.Error("DB without one")
	}
}

// A migration that fails: its error with its key; the ones after it not run.
func TestRunFails(t *testing.T) {
	var (
		g    lifecycle.Graph
		r    recorder
		boom = errors.New("boom")
	)
	g.Add(r.node("a"), r.node("c", "b"), lifecycle.Node{Key: "b", After: []lifecycle.Key{"a"},
		Start: Func(func(context.Context, *gorm.DB) error { return boom })})
	err := run(context.Background(), &g, dummyDB(t), nil)
	if !errors.Is(err, boom) || !strings.Contains(err.Error(), "start b") {
		t.Errorf("err: %v", err)
	}
	if got := strings.Join(r.keys, ","); got != "a" {
		t.Errorf("run: %s", got)
	}
}

// The node of a command (the serve's): the migrations when it starts, on the
// database given then; nothing to close.
func TestNode(t *testing.T) {
	var (
		g     lifecycle.Graph
		r     recorder
		serve lifecycle.Graph
		db    = dummyDB(t)
		order []string
	)
	g.Add(r.node("i18n"))
	serve.Add(
		lifecycle.Node{Key: "db", Start: func(context.Context) (io.Closer, error) { order = append(order, "db"); return nil, nil }},
		node(&g, func() *gorm.DB { order = append(order, "migrate"); return db }, nil, "db"),
	)
	if err := serve.Run(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(append(order, r.keys...), ","); got != "db,migrate,i18n" {
		t.Errorf("run: %s", got)
	}
}

// runCmd runs `app migrate args…`: its output.
func runCmd(t *testing.T, cmd *cli.Command, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	root := (&cli.Command{Name: "app"}).Sub(cmd)
	ctx, err := root.Parse(&cli.CommandContext{InputArgs: append([]string{"migrate"}, args...), Out: &out, Err: &out})
	if err != nil {
		return "", err
	}
	err = ctx.Run()
	return out.String(), err
}

// `migrate`: the migrations run; `migrate -plan`: their order, none run.
func TestCommand(t *testing.T) {
	var (
		g lifecycle.Graph
		r recorder
	)
	g.Add(r.node("i18n"), r.node("worker", "i18n"))
	db := dummyDB(t)
	cmd := command(&g, func() *gorm.DB { return db }, nil)

	out, err := runCmd(t, cmd, "-plan")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "1. i18n\n  2. worker (after i18n)") || len(r.keys) > 0 {
		t.Errorf("plan: %q; run: %v", out, r.keys)
	}

	out, err = runCmd(t, command(&g, func() *gorm.DB { return db }, nil))
	if err != nil {
		t.Fatal(err)
	}
	if out != "migrated: 2\n" || strings.Join(r.keys, ",") != "i18n,worker" {
		t.Errorf("migrate: %q; run: %v", out, r.keys)
	}

	if _, err = runCmd(t, command(&g, func() *gorm.DB { return db }, nil), "extra"); err == nil {
		t.Error("an argument taken")
	}
}

// The application's migrations (presets.Builder.Migrations) with the
// packages', in one order; the registry of the packages not sealed — run
// again, another graph.
func TestExtra(t *testing.T) {
	var (
		g     lifecycle.Graph
		r     recorder
		extra lifecycle.Nodes
		db    = dummyDB(t)
	)
	g.Add(r.node("activity-model", "rvq/activity"), r.node("rvq/activity"))
	extra.Add(r.node("app/trash", "activity-model"))
	cmd := func() *cli.Command { return command(&g, func() *gorm.DB { return db }, Extra(&extra)) }
	out, err := runCmd(t, cmd(), "-plan")
	if err != nil || !strings.Contains(out, "3. app/trash (after activity-model)") {
		t.Fatalf("plan: %q %v", out, err)
	}
	for range 2 {
		if _, err := runCmd(t, cmd()); err != nil {
			t.Fatal(err)
		}
	}
	if got := strings.Join(r.keys, ","); got != "rvq/activity,activity-model,app/trash,rvq/activity,activity-model,app/trash" {
		t.Errorf("run: %s", got)
	}
}

// A database managed: building on it migrates nothing (Auto); another, its
// sessions too, still does.
func TestManaged(t *testing.T) {
	db, other := dummyDB(t), dummyDB(t)
	Manage(db)
	ran := 0
	f := func(*gorm.DB) error { ran++; return nil }
	if err := Auto(db.Session(&gorm.Session{}), f); err != nil || ran != 0 {
		t.Errorf("managed: ran %d, %v", ran, err)
	}
	if err := Auto(other, f); err != nil || ran != 1 || Managed(other) {
		t.Errorf("not managed: ran %d, %v", ran, err)
	}
}
