package helper

import (
	"strings"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// ModelFilterItem is the filter item of the listing by records of mb chosen
// in a dialog (vx.ItemTypeRecordPicker): the search is the records event of
// the listing of mb, at the parents of the request; the condition, column IN
// the keys chosen (one or, many, more). The records chosen in the request
// come as its options, their titles (the labels of the filter). key is the
// key of the item as the FilterDataFunc declares it (no FilterKeyPrefix).
func ModelFilterItem(mb *presets.ModelBuilder, ctx *web.EventContext, key, label, column string, many bool) *vx.FilterItem {
	ensureMessages(mb.Builder().I18n())
	msgs := GetMessages(ctx.Context())

	cfg := &vx.RecordPickerConfig{
		URL:       mb.Info().ListingHrefCtx(ctx),
		ItemValue: "ID",
		ItemText:  "Text",
		Many:      many,
		Query:     map[string]string{presets.ParamsItemTextKey: "Text"},
		Texts: vx.RecordPickerTexts{
			Select:  msgs.SelectorSelect,
			Change:  msgs.SelectorChange,
			Search:  msgs.SelectorSearch,
			NoItems: msgs.SelectorNoItems,
			Close:   msgs.SelectorClose,
		},
	}
	// the titles by the encoder of the model, when it has one (ModelSelect's)
	if enc, _ := mb.GetData(ModelSelectorEncoderKey).(*RecordEncodeFactory); enc != nil {
		mb.Listing().SetRecordEncoderFactory(enc.Name, enc.REF)
		cfg.Query = map[string]string{presets.ParamListingEncoder: enc.Name}
		cfg.ItemText = enc.TextKey
	}

	it := &vx.FilterItem{
		Key:          key,
		Label:        label,
		ItemType:     vx.ItemTypeRecordPicker,
		Modifier:     vx.ModifierIn,
		SQLCondition: column + " " + vx.SQLOperatorPlaceholder + " ?",
		Config:       cfg,
	}

	// the records chosen: their titles
	if v := presets.FilterQueryValue(ctx, key, "in"); v != "" {
		for _, id := range strings.Split(v, ",") {
			if r := fetchByID(mb, id, ctx); r != nil {
				it.Options = append(it.Options, &vx.SelectItem{Value: id, Text: mb.RecordTitle(r, ctx)})
			}
		}
	}
	return it
}
