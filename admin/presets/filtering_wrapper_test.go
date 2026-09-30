package presets

import (
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-rvq/rvq/web"
)

// A wrapper of the tabs adds its tab after the model's, whether the model's
// were set before it or after.
func TestFilterTabsWrapper(t *testing.T) {
	trash := func(_ *web.EventContext, tabs []*FilterTab) []*FilterTab {
		if len(tabs) == 0 {
			tabs = append(tabs, &FilterTab{Label: "All", Default: true})
		}
		return append(tabs, &FilterTab{ID: "trash", Label: "Trash"})
	}
	own := func(*web.EventContext) []*FilterTab {
		return []*FilterTab{{ID: "pending", Label: "Pending", Query: url.Values{"status": {"pending"}}, Default: true}}
	}
	ctx := &web.EventContext{R: httptest.NewRequest("GET", "/", nil)}
	labels := func(tabs []*FilterTab) (r []string) {
		for _, tab := range tabs {
			r = append(r, tab.Label)
		}
		return
	}

	b := &ListingBuilder{}
	b.AppendFilterTabsWrapper(trash)
	if got := labels(b.allFilterTabs(ctx)); len(got) != 2 || got[0] != "All" || got[1] != "Trash" {
		t.Errorf("no tabs of its own: %v", got)
	}
	b.FilterTabsFunc(own)
	tabs := b.allFilterTabs(ctx)
	if got := labels(tabs); len(got) != 2 || got[0] != "Pending" || got[1] != "Trash" {
		t.Errorf("its tabs, set after the wrapper: %v", got)
	}
	if tabs[0].Query.Get("f_status") != "pending" {
		t.Errorf("its tab's filter: %v", tabs[0].Query)
	}
	if def := defaultFilterTab(ctx, tabs); def == nil || def.ID != "pending" {
		t.Errorf("the default tab: %+v", def)
	}
}
