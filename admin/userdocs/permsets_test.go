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

// The permissions a set gives, marked on their items (Via): of each set that
// gives them, of the prefix it gives them with; a line of its own beside;
// shown with "only what is defined"; not a line of their own.
func TestPermValueVia(t *testing.T) {
	t.Cleanup(func() { perm.SetSetExpander(nil) })
	pb, ib := schemaApp(t)
	inRequest(t, ib, func(ctx *web.EventContext) {
		value := "admin:{posts:{%readers, %editors, @list}, settings:{deny %editors}}"
		items, err := PermValueTree(ctx, pb, value, PermValueOptions{View: PermTreeByVerb, Only: true})
		if err != nil {
			t.Fatal(err)
		}
		via := map[string]string{}
		explicit := map[string]bool{}
		walkItems(items, func(it *PermTreeItem) {
			if it.Kind == PermKindSet || it.Edit == "" {
				return
			}
			var vs []string
			for _, v := range it.Via {
				if v.Title == "" {
					t.Errorf("%s: a via of no title", it.Edit)
				}
				vs = append(vs, v.Prefix+":"+v.Set)
			}
			if len(vs) > 0 {
				via[it.Edit] = strings.Join(vs, ",")
			}
			explicit[it.Edit] = it.Explicit
		})
		for edit, want := range map[string]string{
			"admin:posts:@list":    "grant:readers,grant:editors",
			"admin:posts:<*>:@get": "grant:readers,grant:editors",
			"admin:posts:@create":  "grant:editors",
			"admin:settings:@edit": "deny:editors",
		} {
			if got := via[edit]; got != want {
				t.Errorf("%s: %q, want %q", edit, got, want)
			}
		}
		if !explicit["admin:posts:@list"] || explicit["admin:posts:@create"] {
			t.Errorf("the lines of their own: %v", explicit)
		}
	})
}
