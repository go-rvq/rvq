package presets

import "github.com/go-rvq/rvq/x/lifecycle"

// Lifecycle are the nodes the packages of the admin offer while it is built
// (x/lifecycle): what runs while the application serves — the worker, the
// renewal of the secure key —, started by the command that serves, which
// adds them to the graph of its run (cli.CommandContext.Lifecycle), never by
// building the admin.
func (b *Builder) Lifecycle() *lifecycle.Nodes {
	return &b.lifecycle
}

// Migrations are the migrations of the admin that the application decides —
// a model or a table of its own (the activity's log, the revisions of a
// model, the trash of one) —: x/migrate runs them with the ones of every
// package (migrate.Graph), in one graph (migrate.Node, migrate.Command).
func (b *Builder) Migrations() *lifecycle.Nodes {
	return &b.migrations
}
