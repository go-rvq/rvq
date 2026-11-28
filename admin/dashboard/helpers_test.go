package dashboard

import (
	"context"
	"time"
)

// newTest is a builder of tests (no admin, no database) at the clock now,
// its home panel of every group.
func newTest() (*Builder, *Panel) {
	b := New(nil, nil)
	b.Now = func() time.Time { return now }
	return b, b.Panel(PanelHome).AllGroups()
}

func titled(s string) func(context.Context) string { return func(context.Context) string { return s } }
