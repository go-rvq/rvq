package integration_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	. "github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

type tabsRecord struct {
	Sender   string
	Prefix   string
	Server   string
	Mailbox  string
	Unplaced string
}

// A group drawn as tabs (NewFieldsTabs): a tab of each of its groups, by
// its title, every field of every tab in the form (eager windows); the tab
// of a field in error is marked, and is the one open.
func TestFieldsTabs(t *testing.T) {
	pb := New(i18n.New())
	pb.Model(&tabsRecord{})
	fb := NewFieldsBuilder(pb).Model(&tabsRecord{}).Defaults(NewFieldDefaults(WRITE))
	fb = fb.Only(NewFieldsTabs("Settings",
		&FieldsGroup{Name: "Sending", Title: func(context.Context) string { return "Envio" },
			Items: FieldsLayout{"Sender", "Prefix"}},
		&FieldsGroup{Name: "Answers", Title: func(context.Context) string { return "Respostas" },
			Items: FieldsLayout{"Server", "Mailbox"}},
	))

	render := func(vErr *web.ValidationErrors) string {
		ctx := &web.EventContext{R: httptest.NewRequest("POST", "/", nil)}
		if vErr != nil {
			ctx.Flash = vErr
		}
		return h.MustString(fb.ToComponent(nil, nil, &tabsRecord{Server: "imap.example.com"}, FieldModeStack{WRITE}, ctx), context.TODO())
	}

	got := render(nil)
	for _, want := range []string{
		"<v-tabs", "Envio", "Respostas", "<v-window", "data-fields-tabs='Settings'",
		`form["Sender"]`, `form["Prefix"]`, `form["Server"]`, `form["Mailbox"]`, "imap.example.com",
		"{tab: 0}",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the tabs have no %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, `form["Unplaced"]`) {
		t.Error("a field out of the layout is drawn")
	}
	if strings.Count(got, ":eager='true'") != 2 {
		t.Errorf("each tab's fields eager: %d", strings.Count(got, ":eager='true'"))
	}
	if strings.Contains(got, "mdi-alert-circle") {
		t.Error("a tab marked with no error")
	}

	// a field of the second tab in error: it is marked, and open
	vErr := &web.ValidationErrors{}
	vErr.FieldError("Server", "required")
	got = render(vErr)
	if !strings.Contains(got, "mdi-alert-circle") || !strings.Contains(got, "{tab: 1}") {
		t.Errorf("the tab in error is not marked and open:\n%s", got)
	}
}
