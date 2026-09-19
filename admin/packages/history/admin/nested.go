package admin

import (
	"bytes"
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
	// Revisions are recorded automatically on save; they are never created or
	// edited by hand. Read-only also makes a row click open the detail (the
	// compare-with-current view) instead of an edit form.
	child.SetReadonly(true)

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

// recordKey resolves the parent record key for a request: the id in the path
// (the revisions route captures it), or — for a singleton parent, whose route
// carries no {id} segment (see ModelInfo.DetailingHrefCtx: a singleton's href is
// just the base) — the singleton's own record id, loaded through the parent's
// fetcher and read with MustRecordID (the same key createRevision stores).
func (mh *ModelHistory) recordKey(ctx *web.EventContext) string {
	if k := parentRecordKey(ctx); k != "" {
		return k
	}
	if mh.mb.GetSingleton() {
		obj := mh.mb.NewModel()
		if err := mh.mb.Editing().Fetcher(obj, mh.mb.MustRecordID(obj), ctx); err == nil {
			if id := mh.mb.MustRecordID(obj); !id.IsZero() {
				return id.String()
			}
		}
	}
	return ""
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

// latestHashCacheKey scopes the per-request latest-hash cache to this model's
// table (several histories may render on one request).
type latestHashCacheKey string

// cachedLatestHash resolves the record's newest revision hash once per request,
// caching it on the EventContext (the listing calls it for every row). A query
// error yields nil (no highlight) rather than failing the whole listing.
func (mh *ModelHistory) cachedLatestHash(recordKey string, ctx *web.EventContext) histmodels.Hash {
	key := latestHashCacheKey(mh.table + "\x00" + recordKey)
	if v := ctx.ContextValue(key); v != nil {
		hash, _ := v.(histmodels.Hash)
		return hash
	}
	hash, err := mh.latestHash(recordKey)
	if err != nil {
		hash = nil
	}
	ctx.WithContextValue(key, hash)
	return hash
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
	l := child.Listing("Hash", "ChangedFields", "Creator", "CreatedAt", "Published", "Tag", "AccessCount").
		OrderBy("created_at DESC")

	// Highlight the CURRENT record's row (its newest revision) with the success
	// color. latestHash is resolved once per request (cached on the EventContext),
	// then each <tr> whose hash matches gets the "bg-success" class.
	l.RowWrapperFunc(func(row h.MutableAttrHTMLComponent, id string, obj interface{}, dataTableID string, ctx *web.EventContext) h.HTMLComponent {
		rev, ok := obj.(*histmodels.Revision)
		if !ok {
			return row
		}
		cur := mh.cachedLatestHash(rev.RecordKey, ctx)
		if len(cur) > 0 && bytes.Equal(cur, rev.Hash) {
			row.SetAttr("class", "bg-success")
		}
		return row
	})

	// Changed fields: the nested label of what changed in this revision. The
	// column takes the remaining width and truncates with an ellipsis on one
	// line; clicking it toggles to wrap and show the rest.
	l.Field("ChangedFields").
		SetI18nLabel(func(c context.Context) string { return getMessages(c).Changes }).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			rev := field.Obj.(*histmodels.Revision)
			label := mh.ChangedFieldsLabel(rev, ctx)
			// A truncating inline-block (not td width:100%, which would collapse the
			// other columns and let Hash overlap the selection checkbox): wide cap,
			// ellipsis on one line; clicking toggles to wrap and show the rest.
			return h.Td(
				web.Scope(
					h.Div(h.Text(label)).
						Attr("title", label).
						Attr("@click", "locals.expanded = !locals.expanded").
						Attr(":style", "locals.expanded ? 'white-space:normal;word-break:break-word;cursor:pointer' : 'display:inline-block;max-width:40vw;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;vertical-align:bottom;cursor:pointer'"),
				).LocalsInit("{ expanded: false }"),
			)
		})

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
			return mh.compare(mh.recordKey(ctx), a, bb, ctx)
		})

}

// configChildDetailing makes /<parent>/{id}/revisions/{hash} render the revision
// **compared with the current record** — the same view the row's "compare with
// current" used to open. So the revision detail is the comparison (Current |
// Revision, with whole-record and per-field/partial revert), not a static
// snapshot; the separate compare action is therefore not needed.
func (mh *ModelHistory) configChildDetailing(child *presets.ModelBuilder) {
	child.Detailing("Snapshot")

	// Fetch the revision through our own table-qualified query. gorm2op's default
	// Fetch does state.DB.First(&Revision{}), and First appends ORDER BY on the
	// Revision schema's default table "revisions" — but the FROM is <table>_revisions
	// (set via .Table), so Postgres rejects it ("missing FROM-clause entry for table
	// revisions"). mh.Revision uses .Table(mh.table).Take (no schema-qualified order).
	child.Detailing().FetchFunc(func(obj interface{}, id presets.ID, ctx *web.EventContext) error {
		hash, err := decodeHash(id.String())
		if err != nil {
			return err
		}
		rev, err := mh.Revision(mh.recordKey(ctx), hash)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return presets.ErrRecordNotFound
			}
			return err
		}
		*(obj.(*histmodels.Revision)) = *rev
		return nil
	})

	child.Detailing().Field("Snapshot").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		rev := field.Obj.(*histmodels.Revision)
		cur, err := mh.latestHash(rev.RecordKey)
		if err != nil {
			return v.VAlert(h.Text(err.Error())).Type(v.TypeError).Variant(v.VariantTonal)
		}
		comp, cerr := mh.compare(rev.RecordKey, rev.Hash, cur, ctx)
		if cerr != nil {
			return v.VAlert(h.Text(cerr.Error())).Type(v.TypeError).Variant(v.VariantTonal)
		}
		return comp
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
