package userdocs

import (
	"net/url"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/x/i18n"
)

type helpPost struct {
	ID       uint
	Title    string
	Comments []*helpComment
}

type helpComment struct {
	ID   uint
	Text string
}

type helpSettings struct {
	ID   uint
	Name string
}

// The address of a part of the admin — a page, or the event that filled a
// drawer or a dialog — is the node of its document: the deepest there is.
func TestNodeOfURI(t *testing.T) {
	p := presets.New(i18n.New()).URIPrefix("/admin")
	post := p.Model(&helpPost{})
	post.Detailing()
	post.TakeFieldAsChild("Comments")
	settings := p.Model(&helpSettings{}).Singleton(true)
	b := &Builder{p: p}

	pn, sn := b.modelNode(post), b.modelNode(settings)
	cn := b.modelNode(post.Children()[0])
	tree := []*Node{
		{ID: pn, Children: []*Node{
			{ID: pn + "/forms/new"}, {ID: pn + "/forms/edit"}, {ID: pn + "/forms/detail"},
			{ID: pn + "/actions", Children: []*Node{{ID: pn + "/actions/Publish"}}},
			{ID: cn, Children: []*Node{{ID: cn + "/forms/edit"}}},
		}},
		{ID: sn, Children: []*Node{{ID: sn + "/forms/edit"}}},
		{ID: "pages/site-files"},
	}
	listing := post.Info().ListingHref()
	ev := func(path, event string, q ...string) string {
		v := url.Values{"__execute_event__": {event}}
		for i := 0; i+1 < len(q); i += 2 {
			v.Set(q[i], q[i+1])
		}
		return path + "?" + v.Encode()
	}
	for _, c := range []struct{ uri, node string }{
		{listing, pn},
		{listing + "?page=2", pn},
		{listing + "/new", pn + "/forms/new"},
		{listing + "/7", pn + "/forms/detail"},
		{listing + "/7/edit", pn + "/forms/edit"},
		{ev(listing, actions.New), pn + "/forms/new"},
		{ev(listing, actions.Edit), pn + "/forms/new"},
		{ev(listing, actions.Edit, "id", "7"), pn + "/forms/edit"},
		{ev(listing, actions.Detailing, "id", "7"), pn + "/forms/detail"},
		{ev(listing, actions.DoAction, "id", "7", "action", "Publish"), pn + "/actions/Publish"},
		// an action not documented: its model's
		{ev(listing, actions.DoAction, "action", "Other"), pn},
		// nested: under its parent's record
		{listing + "/7/comments", cn},
		{ev(listing+"/7/comments", actions.Edit, "id", "3"), cn + "/forms/edit"},
		// the singleton: its form
		{settings.Info().ListingHref(), sn + "/forms/edit"},
		// a page of the menu
		{"/admin/site-files", "pages/site-files"},
		// of no part
		{"/admin/nothing/here", ""},
		{"/elsewhere", ""},
	} {
		if got := b.NodeOfURI(tree, c.uri); got != c.node {
			t.Errorf("%s: %q, not %q", c.uri, got, c.node)
		}
	}
}
