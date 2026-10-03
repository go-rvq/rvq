package presets

import (
	"errors"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

type trashRecord struct {
	ID      uint
	Title   string
	Deleted bool
}

type trashChild struct {
	ID    uint
	Title string
}

// TRASH qualifies LIST and DETAIL: a listing of the trash is a listing
// (IsList, Has), not LIST (Is compares the whole mode).
func TestFieldModeTrash(t *testing.T) {
	for _, c := range []struct {
		mode                    FieldMode
		list, detail, trash, is bool
	}{
		{LIST, true, false, false, true},
		{LIST | TRASH, true, false, true, false},
		{DETAIL | TRASH, false, true, true, false},
		{EDIT, false, false, false, false},
	} {
		if c.mode.IsList() != c.list || c.mode.IsDetail() != c.detail || c.mode.IsTrash() != c.trash || c.mode.Is(LIST) != c.is {
			t.Errorf("%v: list %v detail %v trash %v Is(LIST) %v", c.mode, c.mode.IsList(), c.mode.IsDetail(), c.mode.IsTrash(), c.mode.Is(LIST))
		}
	}
	stack := FieldModeStack{DETAIL | TRASH}
	if !stack.IsDetail() || !stack.IsTrash() || stack.IsList() || (FieldModeStack{}).IsTrash() {
		t.Error("the stack's mode is its dot's")
	}
}

func TestTrashPolicy(t *testing.T) {
	for _, c := range []struct {
		p       TrashPolicy
		in, out bool
	}{
		{TrashOut, false, true},
		{TrashToo, true, true},
		{TrashOnly, true, false},
	} {
		if c.p.Available(true) != c.in || c.p.Available(false) != c.out {
			t.Errorf("policy %d: in %v out %v", c.p, c.p.Available(true), c.p.Available(false))
		}
	}
}

// trashApp is a model whose listing is of the trash with ?tab=trash, a
// record deleted when it says so, with a nested model, an item of its row
// menu, bulk actions and actions of each policy.
func trashApp() (*Builder, *ModelBuilder, map[uint]*trashRecord) {
	records := map[uint]*trashRecord{1: {ID: 1, Title: "kept"}, 2: {ID: 2, Title: "gone", Deleted: true}}
	b := New(i18n.New()).URIPrefix("/admin")
	mb := b.Model(&trashRecord{})
	mb.SetDeletedFunc(func(obj any) bool { return obj.(*trashRecord).Deleted })
	l := mb.Listing("Title")
	l.SetInTrashFunc(func(ctx *web.EventContext) bool { return ctx.R.URL.Query().Get("tab") == "trash" })
	l.BulkAction("publish")
	l.BulkAction("restore").SetTrash(TrashOnly)
	l.BulkAction("export").SetTrash(TrashToo)
	l.RowMenu().RowMenuItem("Localize")
	mb.AddChild(b.Model(&trashChild{}))
	d := mb.Detailing("Title")
	d.FetchFunc(func(obj any, id ID, ctx *web.EventContext) error {
		r, ok := records[id.GetValue("ID").(uint)]
		if !ok {
			return ErrRecordNotFound
		}
		*obj.(*trashRecord) = *r
		return nil
	})
	d.Action("archive")
	d.Action("origin").SetTrash(TrashOnly)
	d.Action("log").SetTrash(TrashToo)
	return b, mb, records
}

func trashCtx(url string) *web.EventContext {
	return &web.EventContext{R: httptest.NewRequest("GET", url, nil)}
}

func bulkNames(as []*BulkActionBuilder) (names []string) {
	for _, a := range as {
		names = append(names, a.name)
	}
	return
}

// The listing of the trash (LIST | TRASH) offers the bulk actions available
// there, and sets its mode in the context; out of it, the others.
func TestListingTrash(t *testing.T) {
	_, mb, _ := trashApp()
	l := mb.Listing()

	out := trashCtx("/admin/trash-records")
	if got := bulkNames(l.availableBulkActions(out)); !slices.Equal(got, []string{"publish", "export"}) {
		t.Errorf("out of the trash: %v", got)
	}
	if GetFieldMode(out) != LIST || InTrash(out) {
		t.Errorf("out of the trash, the mode %v", GetFieldMode(out))
	}

	in := trashCtx("/admin/trash-records?tab=trash")
	if got := bulkNames(l.availableBulkActions(in)); !slices.Equal(got, []string{"restore", "export"}) {
		t.Errorf("in the trash: %v", got)
	}
	if GetFieldMode(in) != LIST|TRASH || !InTrash(in) {
		t.Errorf("in the trash, the mode %v", GetFieldMode(in))
	}

	// a listing in a detail keeps the detail's mode in the context
	nested := trashCtx("/admin/trash-records?tab=trash")
	WithFieldMode(nested, DETAIL)
	if m := l.requestMode(nested); m != LIST|TRASH || GetFieldMode(nested) != DETAIL {
		t.Errorf("a listing in a detail: %v, the context %v", m, GetFieldMode(nested))
	}
}

// In the trash the menu of a row keeps the nested models, not its items of
// actions.
func TestRowMenuTrash(t *testing.T) {
	_, mb, _ := trashApp()
	var in, out []string
	for _, it := range mb.Listing().RowMenu().Items() {
		if it.availableIn(true) {
			in = append(in, it.Name())
		}
		if it.availableIn(false) {
			out = append(out, it.Name())
		}
	}
	if slices.Contains(in, "Localize") || slices.Contains(in, "Delete") || len(in) != 1 {
		t.Errorf("in the trash: %v", in)
	}
	if !slices.Contains(out, "Localize") || len(out) != len(mb.Listing().RowMenu().Items()) {
		t.Errorf("out of the trash: %v", out)
	}
}

// A deleted record is read only: not edited, not deleted; its detail is
// DETAIL | TRASH.
func TestDeletedRecord(t *testing.T) {
	_, mb, records := trashApp()
	ctx := trashCtx("/admin/trash-records/2")
	kept, gone := records[1], records[2]
	if mb.detailMode(kept) != DETAIL || mb.detailMode(gone) != DETAIL|TRASH {
		t.Errorf("the modes %v, %v", mb.detailMode(kept), mb.detailMode(gone))
	}
	if !mb.EditingRestriction.CanObj(kept, ctx) || mb.EditingRestriction.CanObj(gone, ctx) {
		t.Error("the model's editing restriction")
	}
	if mb.Detailing().CanEditObj(gone, ctx) || mb.Editing().CanEditObj(gone, ctx) || mb.Listing().CanEditObj(gone, ctx) {
		t.Error("a builder edits a deleted record")
	}
	if mb.DeletingRestriction.CanObj(gone, ctx) || mb.Detailing().DeletingRestriction.CanObj(gone, ctx) ||
		mb.Listing().DeletingRestriction.CanObj(gone, ctx) {
		t.Error("a deleted record is deleted again")
	}
}

// The actions of a record: in the menu, those available where it is; run,
// refused where they are not.
func TestActionsTrash(t *testing.T) {
	_, mb, records := trashApp()
	d := mb.Detailing()
	titles := func(obj any, ctx *web.EventContext) (names []string) {
		items, _ := BuildMenuItemCompomentsOfActions("p", ctx, mb, "", obj, d.GetActions()...)
		for range items {
			names = append(names, "item")
		}
		return
	}
	if n := len(titles(records[1], trashCtx("/admin/trash-records/1"))); n != 2 { // archive, log
		t.Errorf("a record kept: %d items", n)
	}
	if n := len(titles(records[2], trashCtx("/admin/trash-records/2"))); n != 2 { // origin, log
		t.Errorf("a deleted record: %d items", n)
	}

	for _, c := range []struct {
		action string
		id     uint
		want   bool
	}{
		{"archive", 1, true}, {"archive", 2, false},
		{"origin", 1, false}, {"origin", 2, true},
		{"log", 1, true}, {"log", 2, true},
	} {
		a := getAction(d.GetActions(), c.action)
		if got := a.trashAllowsObj(mb, records[c.id], trashCtx("/admin/trash-records")); got != c.want {
			t.Errorf("%s on %d: %v", c.action, c.id, got)
		}
	}

	// a bulk action not available where it is asked for is refused
	publish := mb.Listing().BulkAction("publish")
	if err := publish.Do(nil, trashCtx("/admin/trash-records?tab=trash"), &web.EventResponse{}); !errors.Is(err, ErrActionNotAllowed) {
		t.Errorf("publish in the trash: %v", err)
	}
	restore := mb.Listing().BulkAction("restore")
	if err := restore.Do(nil, trashCtx("/admin/trash-records"), &web.EventResponse{}); !errors.Is(err, ErrActionNotAllowed) {
		t.Errorf("restore out of the trash: %v", err)
	}
}

// A trailing field (AppendTrailingFields) is shown after a layout set after
// it (Only), which leaves it out of the listing's fields.
func TestTrailingFieldsAfterOnly(t *testing.T) {
	b := New(i18n.New()).URIPrefix("/admin")
	mb := b.Model(&trashRecord{})
	l := mb.Listing()
	l.Field("Deleted").SetEnabled(func(ctx *FieldContext) bool { return ctx.Mode.IsTrash() })
	l.AppendTrailingFields("Deleted")
	l.SetInTrashFunc(func(ctx *web.EventContext) bool { return ctx.R.URL.Query().Get("tab") == "trash" })
	l.Only("Title")
	names := func(url string) (r []string) {
		for _, c := range l.Columns(trashCtx(url)) {
			r = append(r, c.Name)
		}
		return
	}
	if got := names("/admin/trash-records?tab=trash"); !slices.Equal(got, []string{"Title", "Deleted"}) {
		t.Errorf("in the trash: %v", got)
	}
	if got := names("/admin/trash-records"); !slices.Equal(got, []string{"Title"}) {
		t.Errorf("out of the trash: %v", got)
	}
}
