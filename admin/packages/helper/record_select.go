package helper

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// RecordSelectBuilder is the choice of a record of a model out of a form —
// the select of a ModelSelectorBuilder with no field, in its dialog (the
// search, the list, the pages: as the menu's choice of an organization):
// its records searched by the model's listing (the ListData event at its
// address, under the parents given: the records the user may list, as the
// listing scopes them — the ones they do not reach are not offered), their
// titles by the model's encoder (SetRecordEncoderFactory); what to do when
// one is chosen, OnChange. Its messages: EnsureMessages, before the admin
// serves.
type RecordSelectBuilder struct {
	mb       *presets.ModelBuilder
	label    string
	parents  []model.ID
	value    string
	onChange string
	attrs    []any
}

// RecordSelect is the choice of a record of mb.
func RecordSelect(mb *presets.ModelBuilder) *RecordSelectBuilder {
	return &RecordSelectBuilder{mb: mb}
}

// Label is its label (the model's title, singular, by default).
func (b *RecordSelectBuilder) Label(v string) *RecordSelectBuilder {
	b.label = v
	return b
}

// Parents are the records the model's are under (its listing's address).
func (b *RecordSelectBuilder) Parents(ids ...model.ID) *RecordSelectBuilder {
	b.parents = ids
	return b
}

// Value is the id of the record chosen ("" none).
func (b *RecordSelectBuilder) Value(id string) *RecordSelectBuilder {
	b.value = id
	return b
}

// OnChange is the script run when a record is chosen: the id chosen its
// $event ("" none).
func (b *RecordSelectBuilder) OnChange(js string) *RecordSelectBuilder {
	b.onChange = js
	return b
}

// Attr adds attributes to the select (name, value…).
func (b *RecordSelectBuilder) Attr(attrs ...any) *RecordSelectBuilder {
	b.attrs = append(b.attrs, attrs...)
	return b
}

// Listable says the user of the request may list the records of the model
// under the parents: none, no select (nothing to choose).
func (b *RecordSelectBuilder) Listable(ctx *web.EventContext) bool {
	return !b.mb.Permissioner().Lister(ctx.R, b.parents...).Denied()
}

// Component is the select for the request of ctx, in its dialog: the record
// of Value fetched for its title (one the user may not read: none chosen);
// nil when the user may not list the records (Listable).
func (b *RecordSelectBuilder) Component(ctx *web.EventContext) *vx.VXAdvancedSelectBuilder {
	if !b.Listable(ctx) {
		return nil
	}
	label := b.label
	if label == "" {
		label = b.mb.TTitle(ctx.Context())
	}
	searcher := web.GET().EventFunc(actions.ListData).URL(b.mb.Info().ListingHref(b.parents...))
	configuror, _ := b.mb.GetData(ModelSelectorConfigurorKey).(ModelSelectorConfiguror)
	enc, _ := b.mb.GetData(ModelSelectorEncoderKey).(*RecordEncodeFactory)
	if enc != nil {
		configuror = enc.Configure(b.mb.Listing(), searcher, configuror)
	}
	msgs := GetMessages(ctx.Context())
	sel := vx.VXSelectOne().Label(label).ItemValue("ID").ItemText("Text").ItemsSearcher(searcher).
		Dialog(true).DialogTexts(label, msgs.SelectorSelect, msgs.SelectorChange, msgs.SelectorSearch,
		msgs.SelectorNoItems, msgs.SelectorClose)
	if configuror != nil {
		configuror(sel)
	}
	if itemText := sel.GetAttr("item-text"); itemText != nil {
		searcher.Query(presets.ParamsItemTextKey, itemText.Value.(string))
	}
	value := b.chosen(ctx, sel, enc)
	sel.Attr(":model-value", h.JSONString(value))
	if b.onChange != "" {
		sel.Attr("@update:model-value", b.onChange)
	}
	if len(b.attrs) > 0 {
		sel.Attr(b.attrs...)
	}
	return sel
}

// Chosen is the id of the record of Value when the user may choose it — list
// the records (Listable), read it under the parents and fetch it (the
// model's scope of reading) —, and the record; "" none.
func (b *RecordSelectBuilder) Chosen(ctx *web.EventContext) (string, any) {
	if b.value == "" || !b.Listable(ctx) {
		return "", nil
	}
	id, err := b.mb.ParseRecordID(b.value)
	if err != nil || b.mb.Permissioner().Reader(ctx.R, id, b.parents...).Denied() {
		return "", nil
	}
	rec := b.mb.NewModel()
	if fetch := b.mb.Detailing().GetFetchFunc(); fetch == nil || fetch(rec, id, ctx) != nil {
		return "", nil
	}
	return id.String(), rec
}

// chosen is Chosen, its item in sel.
func (b *RecordSelectBuilder) chosen(ctx *web.EventContext, sel *vx.VXAdvancedSelectBuilder, enc *RecordEncodeFactory) string {
	id, rec := b.Chosen(ctx)
	if id == "" {
		return ""
	}
	if enc != nil {
		sel.Items(enc.REF.EncodeSlice(ctx, rec))
	} else {
		sel.Items([]map[string]string{{"ID": id, "Text": fmt.Sprint(rec)}})
	}
	return id
}
