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
	child.MenuIcon("mdi-history")

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
					// Listing scoped to fields (?field=): keep only the revisions
					// in which ANY of those fields changed.
					if state.SearchParams != nil && key != "" {
						if fields := fieldsParam(state.Ctx); len(fields) > 0 {
							seen := map[string]histmodels.Hash{}
							for _, field := range fields {
								hist, err := mh.FieldHistory(key, field)
								if err != nil {
									return err
								}
								for _, fr := range hist {
									seen[string(fr.Revision.Hash)] = fr.Revision.Hash
								}
							}
							hashes := make([]histmodels.Hash, 0, len(seen))
							for _, hh := range seen {
								hashes = append(hashes, hh)
							}
							state.DB = state.DB.Where("hash IN ?", hashes)
						}
					}
					return nil
				})
			})
	})

	child.RegisterEventFunc(mh.compareEventName(), mh.compareEvent)

	mh.mb.AddChild(child)
	mh.configChildListing(child)
	mh.configChildDetailing(child)
	mh.installRevertConfirm(child)
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

	// Column labels: the Revision struct's fields have no app-level
	// (ModelsI18nModuleKey) translations, so set them from this plugin's own
	// i18n messages instead of leaving the framework to warn about missing keys.
	l.Field("Hash").
		SetI18nLabel(func(c context.Context) string { return getMessages(c).Hash }).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			rev := field.Obj.(*histmodels.Revision)
			return h.Td(h.Code(shortHash(rev.Hash)))
		})
	l.Field("Creator").SetI18nLabel(func(c context.Context) string { return getMessages(c).Author })
	l.Field("CreatedAt").SetI18nLabel(func(c context.Context) string { return getMessages(c).When })
	l.Field("Published").SetI18nLabel(func(c context.Context) string { return getMessages(c).Published })
	l.Field("Tag").SetI18nLabel(func(c context.Context) string { return getMessages(c).Tag })
	l.Field("AccessCount").SetI18nLabel(func(c context.Context) string { return getMessages(c).Accesses })

	// A "Field" select scopes the listing to one field's history (the actual
	// filtering — keeping only revisions where that field changed — is done in
	// the read callback via FieldHistory; this only provides the picker).
	l.FilterDataFunc(func(ctx *web.EventContext) vx.FilterData {
		paths := mh.fieldPaths()
		opts := make([]*vx.SelectItem, len(paths))
		for i, f := range paths {
			opts[i] = &vx.SelectItem{Text: f, Value: f}
		}
		return vx.FilterData{
			{
				Key:      "field",
				Label:    getMessages(ctx.Context()).Field,
				ItemType: vx.ItemTypeTreeSelect,
				Options:  opts,
				SQLConditionFunc: func(val, mod string) (bool, string, []any) {
					return false, "", nil // filtered in the read callback
				},
			},
		}
	})

	msgr := func(ctx *web.EventContext) *Messages { return getMessages(ctx.Context()) }

	// Two selected → field-by-field diff between them. Icon-only trigger, its
	// label as the tooltip.
	l.BulkAction("Diff").
		Icon("mdi-compare").
		SetI18nLabel(func(c context.Context) string { return getMessages(c).Compare }).
		ButtonCompFunc(func(ctx *web.EventContext, title func() string, onclick *web.VueEventTagBuilder) h.HTMLComponent {
			return v.VBtn("").Icon("mdi-compare").
				Color(v.ColorSecondary).Variant(v.VariantFlat).Density(v.DensityComfortable).Class("ml-2").
				Attr("title", title()).
				Attr("@click", onclick.Go())
		}).
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
			return mh.compare(parentRecordKey(ctx), a, bb, ctx)
		})

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
		fields := fieldsParam(ctx)
		// Text/HTML fields render through the detail; structured fields render a
		// readable JSON summary (their interactive detail widgets don't belong in
		// this reconstructed view).
		simple, structured := mh.splitFields(fields)
		comp, _ := mh.detailComponent(obj, rev.RecordKey, simple, ctx)
		m, _ := fieldMap(rev)
		msgr := getMessages(ctx.Context())
		revertBtn := v.VBtn(msgr.Revert).
			Color("warning").Variant(v.VariantTonal).PrependIcon("mdi-history").
			Attr("title", msgr.Revert).
			Attr("@click", mh.revertButtonClick(rev.RecordKey, rev.Hash.String(), "", "", ""))

		out := h.HTMLComponents{
			h.Div(revertBtn).Class("d-flex justify-end mb-3"),
			comp,
			structuredRows(m, structured),
		}
		// Scoped to a single partial-capable field: offer hunk-level revert
		// against the live current value. The field's whole-field Revert button
		// leads the panel's top bar (beside the select-all toggle and counter).
		if len(fields) == 1 && mh.AcceptsPartial(fields[0]) {
			fieldName := fields[0]
			if cur, cerr := mh.currentRecord(rev.RecordKey); cerr == nil {
				leading := v.VBtn("").Icon("mdi-history").
					Variant(v.VariantText).Size(v.SizeSmall).Color("warning").
					Attr("title", msgr.Revert).
					Attr("@click", mh.revertButtonClick(rev.RecordKey, rev.Hash.String(), fieldName, "", ""))
				out = append(out, mh.hunkSelectPanel(
					rev.RecordKey, fieldName, fieldStringValue(cur, fieldName), fieldValue(m, fieldName), rev.Hash.String(), leading, ctx))
			}
		}
		return h.Div(out...)
	})

	// "Compare with current" as a per-row action. It must be a Detailing action
	// with ShowInList, not a listing ItemAction: the row menu fires the shared
	// presets_Action event, which the detailing resolves against its OWN actions
	// (formAction → parseRequestAction). A listing ItemAction there faults with
	// "action required" because that handler never looks at the listing's item
	// actions. View() wraps the returned component in the dialog itself.
	child.Detailing().
		Action("CompareCurrent").
		ShowInList().
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
			return mh.compare(key, a, cur, ctx)
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
