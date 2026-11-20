package worker

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/lifecycle"
	"github.com/go-rvq/rvq/x/migrate"
	"github.com/google/uuid"
)

// countQueue counts its listens and shutdowns.
type countQueue struct {
	Queue
	listens, shutdowns atomic.Int32
}

func (q *countQueue) Listen([]*JobDefinition, func(uuid.UUID) (QueJobInterface, error)) error {
	q.listens.Add(1)
	return nil
}

func (q *countQueue) Shutdown(context.Context) error {
	q.shutdowns.Add(1)
	return nil
}

// The worker as a node: listening when it starts — once, whoever asks —,
// shut down when it closes; after the migrations, if they run.
func TestNode(t *testing.T) {
	q := &countQueue{}
	b := &Builder{q: q}
	n := b.Node("db")
	if n.Key != Key || len(n.After) != 2 || n.After[0] != migrate.Key || n.Optional[0] != migrate.Key {
		t.Fatalf("the node: %+v", n)
	}

	var g lifecycle.Graph
	g.Add(lifecycle.Node{Key: "db"}, n)
	err := g.Run(context.Background(), func(context.Context) error {
		b.Listen() // the application too: once
		if q.listens.Load() != 1 || q.shutdowns.Load() != 0 {
			t.Errorf("running: listens %d, shutdowns %d", q.listens.Load(), q.shutdowns.Load())
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if q.listens.Load() != 1 || q.shutdowns.Load() != 1 {
		t.Errorf("after: listens %d, shutdowns %d", q.listens.Load(), q.shutdowns.Load())
	}
}

// Installed, the worker offers its node to the commands; building the admin
// starts nothing.
func TestInstallOffersTheNode(t *testing.T) {
	q := &countQueue{}
	b := &Builder{q: q}
	pb := presets.New(i18n.New()).URIPrefix("/admin")
	if err := b.Install(pb); err != nil {
		t.Fatal(err)
	}
	nodes := *pb.Lifecycle()
	if len(nodes) != 1 || nodes[0].Key != Key {
		t.Fatalf("offered: %v", nodes)
	}
	if q.listens.Load() != 0 {
		t.Error("listening by being installed")
	}
}
