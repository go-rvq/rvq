package admin

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	h "github.com/go-rvq/htmlgo"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"gorm.io/gorm"
)

// installChild mounts the revisions as a nested model under the versioned model
// (parentMb.AddChild), so its routes live at /<parent>/{id}/revisions/… — its
// listing is the history, its detail renders the record as of that revision.
func (mh *ModelHistory) installChild() {
	b := mh.mb.Builder()

	// NewModelBuilder (not b.Model): a child must NOT join the builder's
	// top-level models, or its nested routes would be set up twice — once by the
	// top-level loop and once by the parent (which registers its children),
	// conflicting. AddChild registers it under the parent only.
	child := presets.NewModelBuilder(b, &histmodels.Revision{}, presets.ModelWithID(mh.table))
	child.URIName("revisions")

	// One Revision Go type is mapped to every model's own <table>_revisions: force
	// the table on each query, and scope reads to the parent record from the path.
	child.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).
			WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
				cb.Pre(func(state *gorm2op.CallbackState) error {
					state.DB = state.DB.Table(mh.table)
					key := parentRecordKey(state.Ctx)
					if key != "" {
						state.DB = state.DB.Where("record_key = ?", key)
					}
					// Listing scoped to one field (?field=): keep only the
					// revisions in which that field actually changed.
					if state.SearchParams != nil && key != "" {
						if field := fieldParam(state.Ctx); field != "" {
							hist, err := mh.FieldHistory(key, field)
							if err != nil {
								return err
							}
							hashes := make([]histmodels.Hash, len(hist))
							for i, fr := range hist {
								hashes[i] = fr.Revision.Hash
							}
							state.DB = state.DB.Where("hash IN ?", hashes)
						}
					}
					return nil
				})
			})
	})

	child.RegisterEventFunc(mh.revertEventName(), mh.revertEvent)
	child.RegisterEventFunc(mh.revertHunksEventName(), mh.revertHunksEvent)

	mh.mb.AddChild(child)
	mh.configChildListing(child)
	mh.configChildDetailing(child)
}

func (mh *ModelHistory) revertEventName() string { return "history_revert_" + mh.table }

// revertEvent restores the parent record to the chosen revision (git-revert:
// records a new revision), then navigates to the parent's detail showing it.
func (mh *ModelHistory) revertEvent(ctx *web.EventContext) (r web.EventResponse, err error) {
	recordKey := parentRecordKey(ctx)
	hash, err := decodeHash(ctx.R.FormValue("hash"))
	if err != nil {
		return
	}
	id, err := mh.mb.ParseRecordID(recordKey)
	if err != nil {
		return
	}
	obj := mh.mb.NewModel()
	id.SetTo(obj)
	if err = mh.db.First(obj).Error; err != nil {
		return
	}
	if err = mh.RevertRecord(obj, hash, ctx); err != nil {
		return
	}
	presets.ShowMessage(&r, getMessages(ctx.Context()).Reverted, "success")
	r.PushState = web.Location(nil).URL(mh.mb.Info().DetailingHref(recordKey))
	return
}

// parentRecordKey is the id of the parent record in the request path — the
// revisions' RecordKey (the same string capture stores).
func parentRecordKey(ctx *web.EventContext) string {
	ids := presets.ParentsModelID(ctx.R)
	if len(ids) == 0 {
		return ""
	}
	return ids.Last().String()
}

func decodeHash(s string) (histmodels.Hash, error) {
	b, err := hex.DecodeString(s)
	return histmodels.Hash(b), err
}

// latestHash is the record's current (newest) revision — what "compare with
// current" compares against; nil when the record has no revision yet.
func (mh *ModelHistory) latestHash(recordKey string) (histmodels.Hash, error) {
	return scanLatestHash(mh.db, mh.table, recordKey)
}

// scanLatestHash reads the newest revision's hash of a record from a given db
// (so it works inside a publish transaction). It scans through database/sql's
// Row, so the bytea runs through Hash's sql.Scanner. A gorm .Scan(&hash) would
// instead see Hash's []byte kind as a slice of rows and try to scan the 32-byte
// value into a single byte ("converting []uint8 to a uint8: invalid syntax").
func scanLatestHash(db *gorm.DB, table, recordKey string) (histmodels.Hash, error) {
	var hash histmodels.Hash
	err := db.Table(table).
		Select("hash").
		Where("record_key = ?", recordKey).
		Order("created_at DESC").
		Limit(1).
		Row().
		Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return hash, err
}

func (mh *ModelHistory) configChildListing(child *presets.ModelBuilder) {
	l := child.Listing("Hash", "Creator", "CreatedAt", "Published", "Tag", "AccessCount").
		OrderBy("created_at DESC")

	l.Field("Hash").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		rev := field.Obj.(*histmodels.Revision)
		return h.Td(h.Code(shortHash(rev.Hash)))
	})

	// A "Field" select scopes the listing to one field's history (the actual
	// filtering — keeping only revisions where that field changed — is done in
	// the read callback via FieldHistory; this only provides the picker).
	l.FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		opts := make([]*vx.SelectItem, len(mh.resolved))
		for i, f := range mh.resolved {
			opts[i] = &vx.SelectItem{Text: f, Value: f}
		}
		return vx.FilterData{
			{
				Key:      "field",
				Label:    getMessages(ctx.Context()).Field,
				ItemType: vx.ItemTypeSelect,
				Options:  opts,
				SQLConditionFunc: func(val, mod string) (bool, string, []any) {
					return false, "", nil // filtered in the read callback
				},
			},
		}
	})

	msgr := func(ctx *web.EventContext) *Messages { return getMessages(ctx.Context()) }

	// Two selected → field-by-field diff between them.
	l.BulkAction("Diff").
		SetI18nLabel(func(c context.Context) string { return getMessages(c).Compare }).
		ComponentFunc(func(selectedIds []string, ctx *web.EventContext) (h.HTMLComponent, error) {
			if len(selectedIds) != 2 {
				return v.VAlert(h.Text(msgr(ctx).SelectTwoHint)).Type("warning").Variant(v.VariantTonal), nil
			}
			a, err := decodeHash(selectedIds[0])
			if err != nil {
				return nil, err
			}
			bb, err := decodeHash(selectedIds[1])
			if err != nil {
				return nil, err
			}
			return mh.diffDialog(parentRecordKey(ctx), a, bb, ctx)
		})

	// One selected → compare it with the current revision.
	l.ItemAction("CompareCurrent").
		Icon("mdi-compare").
		SetI18nLabel(func(c context.Context) string { return getMessages(c).CompareCurrent }).
		ComponentFunc(func(id string, ctx *web.EventContext) (h.HTMLComponent, error) {
			a, err := decodeHash(id)
			if err != nil {
				return nil, err
			}
			key := parentRecordKey(ctx)
			cur, err := mh.latestHash(key)
			if err != nil {
				return nil, err
			}
			return mh.diffDialog(key, a, cur, ctx)
		})
}

func (mh *ModelHistory) diffDialog(recordKey string, a, b histmodels.Hash, ctx *web.EventContext) (h.HTMLComponent, error) {
	sec, err := mh.compare(recordKey, a, b, ctx)
	if err != nil {
		return nil, err
	}
	return vx.VXDialog().Title(getMessages(ctx.Context()).History).Width("1200").SlotBody(sec), nil
}

// configChildDetailing makes /<parent>/{id}/revisions/{hash} render the record
// as it was at that revision — the parent's own detail fields, fed an object
// with the revision's snapshot applied.
func (mh *ModelHistory) configChildDetailing(child *presets.ModelBuilder) {
	child.Detailing("Snapshot")
	child.Detailing().Field("Snapshot").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		rev := field.Obj.(*histmodels.Revision)
		obj, err := mh.applied(rev)
		if err != nil {
			return v.VAlert(h.Text(err.Error())).Type(v.TypeError).Variant(v.VariantTonal)
		}
		fieldName := fieldParam(ctx)
		// Text/HTML fields render through the detail; structured fields render a
		// readable JSON summary (their interactive detail widgets don't belong in
		// this reconstructed view).
		simple, structured := mh.splitFields(fieldName)
		comp, _ := mh.detailComponent(obj, rev.RecordKey, simple, ctx)
		m, _ := fieldMap(rev)
		msgr := getMessages(ctx.Context())
		revertBtn := v.VBtn(msgr.Revert).
			Color("warning").Variant(v.VariantTonal).PrependIcon("mdi-history").
			Attr("@click", web.Plaid().
				EventFunc(mh.revertEventName()).
				Query("hash", rev.Hash.String()).
				Go())

		out := h.HTMLComponents{
			h.Div(revertBtn).Class("d-flex justify-end mb-3"),
			comp,
			structuredRows(m, structured),
		}
		// Scoped to a single partial-capable field: offer hunk-level revert
		// against the live current value.
		if fieldName != "" && mh.AcceptsPartial(fieldName) {
			if cur, cerr := mh.currentRecord(rev.RecordKey); cerr == nil {
				out = append(out, mh.hunkSelectPanel(
					rev.RecordKey, fieldName, fieldStringValue(cur, fieldName), fieldValue(m, fieldName), rev.Hash.String(), ctx))
			}
		}
		return h.Div(out...)
	})
}

// currentRecord loads the live parent record for recordKey.
func (mh *ModelHistory) currentRecord(recordKey string) (any, error) {
	id, err := mh.mb.ParseRecordID(recordKey)
	if err != nil {
		return nil, err
	}
	obj := mh.mb.NewModel()
	id.SetTo(obj)
	if err := mh.db.First(obj).Error; err != nil {
		return nil, err
	}
	return obj, nil
}

// applied loads the parent record and overlays the revision's versioned-field
// snapshot onto it, so it reads as it did at that revision.
func (mh *ModelHistory) applied(rev *histmodels.Revision) (any, error) {
	id, err := mh.mb.ParseRecordID(rev.RecordKey)
	if err != nil {
		return nil, err
	}
	obj := mh.mb.NewModel()
	id.SetTo(obj)
	if err := mh.db.First(obj).Error; err != nil {
		return nil, err
	}
	m, err := fieldMap(rev)
	if err != nil {
		return nil, err
	}
	for _, f := range mh.resolved {
		raw, ok := m[f]
		if !ok {
			continue
		}
		if err := setFieldFromJSON(obj, f, raw); err != nil {
			return nil, fmt.Errorf("history: apply field %q: %w", f, err)
		}
	}
	return obj, nil
}
