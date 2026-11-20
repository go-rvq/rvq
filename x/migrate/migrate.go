// Package migrate is the migrations of the application: one graph
// (x/lifecycle) for all the packages compiled, each registering its own in
// its init().
//
//	func init() {
//		migrate.Graph.Add(lifecycle.Node{
//			Key:   worker.Key,
//			After: []lifecycle.Key{i18n.Key},
//			Start: migrate.Func(func(ctx context.Context, db *gorm.DB) error {
//				return db.AutoMigrate(&QorJob{}, …)
//			}),
//		})
//	}
//
// A migration is a node whose Start migrates and returns no Closer: the
// database it migrates comes in the ctx (WithDB). Registered in init(), only
// what was compiled migrates — a feature is a build tag: one out of the
// binary registers nothing, and who depends on it says so in Optional.
//
// Graph is a registry: each run is a graph of its own, of its nodes and the
// application's (presets.Builder.Migrations) — two applications or tests run
// at once, each on the database of its ctx. `app migrate` (Command) runs
// them; the serve, as its node Key (Node), before what needs the tables.
package migrate

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/go-rvq/rvq/cli"
	"github.com/go-rvq/rvq/x/lifecycle"
	"gorm.io/gorm"
)

// Key is the node of a command that runs the migrations (Node): what needs
// the tables starts after it.
const Key lifecycle.Key = "rvq/migrate"

// Graph is the migrations of every package compiled, registered in their
// init().
var Graph lifecycle.Graph

type dbKey struct{}

// WithDB is ctx with the database the migrations migrate.
func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey{}, db)
}

// DB is the database of ctx (WithDB); nil, none.
func DB(ctx context.Context) *gorm.DB {
	db, _ := ctx.Value(dbKey{}).(*gorm.DB)
	return db
}

// ErrNoDB is a migration run without a database in its ctx.
var ErrNoDB = errors.New("migrate: no database in the ctx (WithDB)")

// Func is the Start of a migration: f with the database of the ctx, a new
// session of it; no Closer.
func Func(f func(ctx context.Context, db *gorm.DB) error) func(ctx context.Context) (io.Closer, error) {
	return func(ctx context.Context) (io.Closer, error) {
		db := DB(ctx)
		if db == nil {
			return nil, ErrNoDB
		}
		return nil, f(ctx, db.Session(&gorm.Session{NewDB: true, Context: ctx}))
	}
}

// Run runs the migrations of every package (Graph) and the extra ones (the
// application's: presets.Builder.Migrations) on db, in the order of one
// graph.
func Run(ctx context.Context, db *gorm.DB, extra ...lifecycle.Node) error {
	return run(ctx, &Graph, db, extra)
}

// graph is the graph of a run: the nodes of base and the extra ones.
func graph(base *lifecycle.Graph, extra []lifecycle.Node) *lifecycle.Graph {
	g := &lifecycle.Graph{}
	g.Add(base.Nodes()...)
	g.Add(extra...)
	return g
}

func run(ctx context.Context, base *lifecycle.Graph, db *gorm.DB, extra []lifecycle.Node) error {
	if db == nil {
		return ErrNoDB
	}
	return graph(base, extra).Run(WithDB(ctx, db), nil)
}

func nodes(extra func() []lifecycle.Node) []lifecycle.Node {
	if extra == nil {
		return nil
	}
	return extra()
}

// Extra is extra (nil: none) as the extra nodes of Node and Command: what
// it holds when they run.
func Extra(extra *lifecycle.Nodes) func() []lifecycle.Node {
	return func() []lifecycle.Node {
		if extra == nil {
			return nil
		}
		return *extra
	}
}

// Node is the node of a command that runs the migrations (Key) — of every
// package and the ones extra gives (nil: none) when it starts —, on the
// database db gives then; after the keys given (the database's).
func Node(db func() *gorm.DB, extra func() []lifecycle.Node, after ...lifecycle.Key) lifecycle.Node {
	return node(&Graph, db, extra, after...)
}

func node(base *lifecycle.Graph, db func() *gorm.DB, extra func() []lifecycle.Node, after ...lifecycle.Key) lifecycle.Node {
	return lifecycle.Node{Key: Key, After: after, Start: func(ctx context.Context) (io.Closer, error) {
		return nil, run(ctx, base, db(), nodes(extra))
	}}
}

// Command is `migrate [-plan]`: the migrations of every package and the
// ones extra gives (nil: none) run on the database db gives; -plan, their
// order, nothing run. extra is asked when the command runs: an application
// built then gives its migrations then.
func Command(db func() *gorm.DB, extra func() []lifecycle.Node) *cli.Command {
	return command(&Graph, db, extra)
}

func command(base *lifecycle.Graph, db func() *gorm.DB, extra func() []lifecycle.Node) *cli.Command {
	var plan bool
	return &cli.Command{
		Name:        "migrate",
		Description: "migrates the database (-plan: the order, nothing run)",
		New: func(ctx *cli.CommandContext) error {
			ctx.Flags().BoolVar(&plan, "plan", false, "print the order of the migrations, run nothing")
			return nil
		},
		ParseArgs: func(ctx *cli.CommandContext) error {
			return ctx.Args.Eq(0)
		},
		Run: func(ctx *cli.CommandContext) error {
			g := graph(base, nodes(extra))
			if plan {
				p, err := g.Plan()
				if err != nil {
					return err
				}
				_, err = io.WriteString(ctx.Out, p)
				return err
			}
			sorted, err := g.Sorted()
			if err != nil {
				return err
			}
			d := db()
			if d == nil {
				return ErrNoDB
			}
			if err = g.Run(WithDB(ctx.Context, d), nil); err != nil {
				return err
			}
			_, err = fmt.Fprintf(ctx.Out, "migrated: %d\n", len(sorted))
			return err
		},
	}
}
