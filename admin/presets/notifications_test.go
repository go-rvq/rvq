package presets

import (
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

// The notification center sums its sources and shows their contents in
// order, NotificationFunc's first; a source with nothing is left out.
func TestNotifications(t *testing.T) {
	p := New(i18n.New())
	ctx := &web.EventContext{}
	if total, content := p.Notifications(ctx); total != 0 || len(content) != 0 {
		t.Fatalf("no source: %d %d", total, len(content))
	}
	p.AddNotificationSource(NotificationSource{
		Name:    "invites",
		Count:   func(*web.EventContext) int { return 2 },
		Content: func(*web.EventContext) h.HTMLComponent { return h.Text("invites") },
	})
	p.AddNotificationSource(NotificationSource{
		Name:    "none",
		Count:   func(*web.EventContext) int { return 0 },
		Content: func(*web.EventContext) h.HTMLComponent { return nil },
	})
	p.NotificationFunc(func(*web.EventContext) h.HTMLComponent { return h.Text("notes") },
		func(*web.EventContext) int { return 3 })
	total, content := p.Notifications(ctx)
	if total != 5 || len(content) != 2 {
		t.Fatalf("total %d, contents %d", total, len(content))
	}
	if content[0].(h.RawHTML) != "notes" || content[1].(h.RawHTML) != "invites" {
		t.Errorf("the order: %v", content)
	}
}
