package userdocs

import (
	"strings"
	"testing"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

// The sets in the field of a permission: offered by the + of a resource,
// with what they give there; a line of one, a row of the set under its
// resource, of its prefix, with its lines — a set of nothing there, an
// error —; what it gives, active; a line of one added, removed as any.
func TestPermValueSets(t *testing.T) {
	t.Cleanup(func() { perm.SetSetExpander(nil) })
	pb, ib := schemaApp(t)
	inRequest(t, ib, func(ctx *web.EventContext) {
		ch, err := PermChoicesOf(ctx.Context(), pb, "admin:posts:", false, nil)
		if err != nil {
			t.Fatal(err)
		}
		var readers *PermChoice
		for i, c := range ch.Choices {
			if c.Segment == "%readers" {
				readers = &ch.Choices[i]
			}
		}
		if readers == nil || readers.Kind != PermKindSet || readers.Resource != "admin:posts:%readers" ||
			strings.Join(readers.Lines, ",") != "@list,<*>:@get" || readers.Title == "readers" {
			t.Fatalf("readers offered: %+v", readers)
		}

		value, err := EditValue("", PermEditAdd, []string{"admin:posts:%readers"}, "")
		if err != nil {
			t.Fatal(err)
		}
		value, _ = EditValue(value, PermEditAdd, []string{"admin:settings:%removers"}, "deny")
		items, err := PermValueTree(ctx, pb, value, PermValueOptions{View: PermTreeByVerb})
		if err != nil {
			t.Fatal(err)
		}
		var sets []*PermTreeItem
		var list *PermTreeItem
		walkItems(items, func(it *PermTreeItem) {
			if it.Kind == PermKindSet {
				sets = append(sets, it)
			}
			if it.Edit == "admin:posts:@list" {
				list = it
			}
		})
		if len(sets) != 2 {
			t.Fatalf("the rows of the sets: %d (%s)", len(sets), value)
		}
		r, rm := sets[0], sets[1]
		if r.Subtitle != "%readers" || r.Prefix != "grant" || r.Error || strings.Join(r.Lines, ",") != "@list,<*>:@get" {
			t.Errorf("readers: %+v", r)
		}
		// removers of a singleton: nothing there, ignored
		if rm.Subtitle != "%removers" || rm.Prefix != "deny" || !rm.Error || len(rm.Lines) != 0 || rm.Description == "" {
			t.Errorf("removers: %+v", rm)
		}
		if list == nil || list.Active == PermActiveNone || list.Active == "" {
			t.Errorf("what the set gives not active: %+v", list)
		}
		if value, _ = EditValue(value, PermEditRemove, []string{r.Edit}, ""); strings.Contains(value, "readers") {
			t.Errorf("removed: %s", value)
		}
	})
}
