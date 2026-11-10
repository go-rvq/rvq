package userdocs

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

type listingActionsModel struct{ ID uint }

// The actions of a listing — of its top and of each row — are documented as
// the detailing's: admin.action finds them, by their title, and each is a
// node under the model's actions.
func TestListingActions(t *testing.T) {
	p := presets.New(i18n.New())
	mb := p.Model(&listingActionsModel{})
	mb.Listing().Action("top").SetI18nLabel(func(context.Context) string { return "Top" })
	mb.Listing().ItemAction("row").SetI18nLabel(func(context.Context) string { return "Row" })
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}

	if got := listingActions(mb); len(got) != 2 || got[0].Name() != "top" || got[1].Name() != "row" {
		t.Fatalf("the listing's actions: %v", got)
	}
	for name, want := range map[string]string{"top": "Top", "row": "Row"} {
		if title, ok := actionTitle(mb, name, ctx); !ok || title != want {
			t.Errorf("%s: %q %v", name, title, ok)
		}
	}
	if _, ok := actionTitle(mb, "none", ctx); ok {
		t.Error("an action there is not")
	}
}
