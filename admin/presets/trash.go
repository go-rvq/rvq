package presets

import (
	"github.com/go-rvq/rvq/web"
)

// The trash: the deleted records, listed (LIST | TRASH) and shown
// (DETAIL | TRASH). A deleted record is read only: it is not edited, and
// none of its actions runs but those available in the trash (TrashPolicy) —
// its menu keeps the models nested in it and its pages.
//
// Whoever keeps deleted records says where they are: the listing, by
// ListingBuilder.SetInTrashFunc (a tab of it, say); a record, by
// ModelBuilder.SetDeletedFunc. The mode of a request — LIST, DETAIL, with
// TRASH or not — is in its context (GetFieldMode), and what is offered there
// follows it.

// TrashPolicy is where an action is available: out of the trash — the
// default —, in it too, or in it only.
type TrashPolicy uint8

const (
	// TrashOut: out of the trash only (the default)
	TrashOut TrashPolicy = iota
	// TrashToo: in the trash too
	TrashToo
	// TrashOnly: in the trash only (bringing a record back)
	TrashOnly
)

// Available says the policy lets the action be in the trash (inTrash) or
// out of it.
func (p TrashPolicy) Available(inTrash bool) bool {
	if inTrash {
		return p != TrashOut
	}
	return p != TrashOnly
}

type fieldModeKeyType int

const fieldModeKey fieldModeKeyType = iota

// WithFieldMode sets the mode of the request of ctx: the listing's (LIST),
// a detail's (DETAIL), each with TRASH when of the deleted records.
func WithFieldMode(ctx *web.EventContext, mode FieldMode) {
	ctx.WithContextValue(fieldModeKey, mode)
}

// GetFieldMode is the mode of the request of ctx (WithFieldMode), 0 when
// none was set.
func GetFieldMode(ctx *web.EventContext) FieldMode {
	m, _ := ctx.ContextValue(fieldModeKey).(FieldMode)
	return m
}

// InTrash says the request of ctx is of the trash.
func InTrash(ctx *web.EventContext) bool {
	return GetFieldMode(ctx).IsTrash()
}

// SetDeletedFunc sets what says a record of the model is deleted: in the
// trash, read only.
func (mb *ModelBuilder) SetDeletedFunc(f func(obj any) bool) *ModelBuilder {
	mb.deletedFunc = f
	return mb
}

// IsDeleted says obj is a deleted record (SetDeletedFunc).
func (mb *ModelBuilder) IsDeleted(obj any) bool {
	return mb.deletedFunc != nil && obj != nil && mb.deletedFunc(obj)
}

// SetInTrashFunc sets what says the listing of a request is of the trash —
// it lists the deleted records.
func (b *ListingBuilder) SetInTrashFunc(f func(ctx *web.EventContext) bool) *ListingBuilder {
	b.inTrashFunc = f
	return b
}

// InTrash says the listing of the request of ctx is of the trash.
func (b *ListingBuilder) InTrash(ctx *web.EventContext) bool {
	return b.inTrashFunc != nil && b.inTrashFunc(ctx)
}

// requestMode is the mode of the listing of the request of ctx: the one of
// its context when it is a listing's, else its own (Mode) — set in the
// context when it has none (a listing in a detail keeps the detail's).
func (b *ListingBuilder) requestMode(ctx *web.EventContext) FieldMode {
	m := GetFieldMode(ctx)
	if m.IsList() {
		return m
	}
	own := b.Mode(ctx)
	if m == 0 {
		WithFieldMode(ctx, own)
	}
	return own
}

// Mode is the mode of the listing of the request: LIST, with TRASH when it
// is of the trash.
func (b *ListingBuilder) Mode(ctx *web.EventContext) FieldMode {
	if b.InTrash(ctx) {
		return LIST | TRASH
	}
	return LIST
}

// SetTrash sets where the action is available (TrashPolicy).
func (b *ActionBuilder) SetTrash(p TrashPolicy) *ActionBuilder {
	b.trash = p
	return b
}

// Trash is where the action is available (TrashPolicy).
func (b *ActionBuilder) Trash() TrashPolicy { return b.trash }

// SetTrash sets where the bulk action is available (TrashPolicy).
func (b *BulkActionBuilder) SetTrash(p TrashPolicy) *BulkActionBuilder {
	b.trash = p
	return b
}

// Trash is where the bulk action is available (TrashPolicy).
func (b *BulkActionBuilder) Trash() TrashPolicy { return b.trash }

// SetTrash sets where the item is available (TrashPolicy). An item that
// opens a nested model is in the trash whatever it says.
func (b *RowMenuItemBuilder) SetTrash(p TrashPolicy) *RowMenuItemBuilder {
	b.trash = p
	return b
}

// availableIn says the item is in the menu of a record in the trash
// (inTrash) or out of it.
func (b *RowMenuItemBuilder) availableIn(inTrash bool) bool {
	return b.child != nil || b.trash.Available(inTrash)
}

// availableBulkActions are the bulk actions the listing of the request
// offers: allowed, and available where it is — in the trash or out of it.
func (b *ListingBuilder) availableBulkActions(ctx *web.EventContext) (as []*BulkActionBuilder) {
	inTrash := b.requestMode(ctx).IsTrash()
	for _, ba := range b.bulkActions {
		if ba.Verifier(b.mb.permissioner.ReqList(ctx.R)).Denied() || !ba.trash.Available(inTrash) {
			continue
		}
		as = append(as, ba)
	}
	return
}

// trashAllows says the action may run where it is asked for: on a deleted
// record (fetched by id) or in a listing of the trash, only if available in
// the trash; elsewhere only if available out of it.
func (b *ActionBuilder) trashAllows(baseModel *ModelBuilder, id string, ctx *web.EventContext) bool {
	if baseModel == nil {
		return true
	}
	var obj any
	if id != "" && baseModel.deletedFunc != nil && b.db != nil {
		obj, _ = b.db.Fetch(id, ctx)
	}
	return b.trashAllowsObj(baseModel, obj, ctx)
}

// trashAllowsObj says the action may run on obj — a deleted record, only if
// available in the trash —, or, with no record, in the listing of the
// request.
func (b *ActionBuilder) trashAllowsObj(baseModel *ModelBuilder, obj any, ctx *web.EventContext) bool {
	inTrash := false
	if obj != nil {
		inTrash = baseModel.IsDeleted(obj)
	} else if b.typ == ActionTypeList && baseModel.listing != nil {
		inTrash = baseModel.listing.InTrash(ctx)
	}
	return b.trash.Available(inTrash)
}

// trashAllows says the bulk action may run in the listing of the request:
// of the trash, only if available there; elsewhere only if available out
// of it.
func (b *BulkActionBuilder) trashAllows(ctx *web.EventContext) bool {
	return b.l == nil || b.trash.Available(b.l.InTrash(ctx))
}

// detailMode is the mode of the detail of obj: DETAIL, with TRASH when it is
// deleted.
func (mb *ModelBuilder) detailMode(obj any) FieldMode {
	if mb.IsDeleted(obj) {
		return DETAIL | TRASH
	}
	return DETAIL
}
