package admin

import (
	"context"
	"encoding/hex"
	"fmt"

	h "github.com/go-rvq/htmlgo"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// installChild mounts the revisions as a nested model under the versioned model
// (parentMb.AddChild), so its routes live at /<parent>/{id}/revisions/… — its
// listing is the history, its detail renders the record as of that revision.
func (mh *ModelHistory) installChild() {
	b := mh.mb.Builder()

	child := b.Model(&histmodels.Revision{}, presets.ModelWithID(mh.table))
	child.URIName("revisions")

	// One Revision Go type is mapped to every model's own <table>_revisions: force
	// the table on each query, and scope reads to the parent record from the path.
	child.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).
			WithReadCallbacks(func(cb *gorm2op.Callbacks[*gorm2op.DataOperatorBuilder]) {
				cb.Pre(func(state *gorm2op.CallbackState) error {
					state.DB = state.DB.Table(mh.table)
					if key := parentRecordKey(state.Ctx); key != "" {
						state.DB = state.DB.Where("record_key = ?", key)
					}
					return nil
				})
			})
	})

	mh.mb.AddChild(child)
	mh.configChildListing(child)
	mh.configChildDetailing(child)
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
// current" compares against.
func (mh *ModelHistory) latestHash(recordKey string) (histmodels.Hash, error) {
	var hash histmodels.Hash
	err := mh.db.Table(mh.table).
		Select("hash").
		Where("record_key = ?", recordKey).
		Order("created_at DESC").
		Limit(1).
		Scan(&hash).Error
	return hash, err
}

func (mh *ModelHistory) configChildListing(child *presets.ModelBuilder) {
	l := child.Listing("Hash", "Creator", "CreatedAt", "Published", "Tag", "AccessCount").
		OrderBy("created_at DESC")

	l.Field("Hash").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		rev := field.Obj.(*histmodels.Revision)
		return h.Td(h.Code(shortHash(rev.Hash)))
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
		comp, _ := mh.detailComponent(obj, rev.RecordKey, ctx)
		return comp
	})
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
