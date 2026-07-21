package comps

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/l10n"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/thirdpart/gorm/datatypes"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/i18n"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/sunfmin/reflectutils"
)

// jsonMapCellString renders a JSON map value for an editor cell: plain strings
// are shown as-is; nested objects/arrays (and other non-string values) are shown
// as indented JSON so they can be edited and round-tripped. This lets the flat
// data-table editor hold nested structures — e.g. the `{LANG: {KEY: VALUE}}`
// shape — by editing the nested part as JSON text.
func jsonMapCellString(val interface{}) string {
	switch t := val.(type) {
	case nil:
		return ""
	case string:
		return t
	default:
		if b, err := json.MarshalIndent(t, "", "  "); err == nil {
			return string(b)
		}
		return fmt.Sprint(val)
	}
}

// jsonMapCellParse parses an editor cell back into a JSON map value: a cell whose
// (trimmed) content is a JSON object or array is stored as the nested structure;
// anything else is stored as a plain string.
func jsonMapCellParse(s string) interface{} {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		var out interface{}
		if err := json.Unmarshal([]byte(trimmed), &out); err == nil {
			return out
		}
	}
	return s
}

type Pair struct {
	Name  string
	Value string
}

func ConfigureFieldJSONMap(mb *presets.ModelBuilder, field *presets.FieldBuilder, keyComp func(fctx *presets.FieldContext) h.HTMLComponent) *presets.FieldBuilder {
	t := mb.Translator()

	return field.
		SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
			fieldValue := datatypes.NullJSONMap{}

			for k, v := range ctx.R.PostForm {
				if v[0] != "" && strings.HasPrefix(k, field.FormKey+"[") && strings.HasSuffix(k, "].Name") {
					kValue := strings.TrimSuffix(k, "Name") + "Value"
					if fv := ctx.R.PostForm.Get(kValue); fv != "" {
						fieldValue[v[0]] = jsonMapCellParse(fv)
					}
				}
			}

			reflectutils.Set(obj, field.Name, fieldValue)

			return nil
		}).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			value, _ := field.Value().(datatypes.NullJSONMap)
			return configureFieldJSONMapInput(t, field.EventContext, field.Label, field.FormKey, value, keyComp(field))
		})
}

func configureFieldJSONMapInput(tr i18n.Translator, ctx *web.EventContext, label, fieldName string, value datatypes.NullJSONMap, keyComp h.HTMLComponent) *vue.UserComponentBuilder {
	if value == nil {
		value = datatypes.NullJSONMap{}
	}

	var keys []string

	for k := range value {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var items = make([]*Pair, len(keys))
	for i, key := range keys {
		items[i] = &Pair{Name: key, Value: jsonMapCellString(value[key])}
	}

	t := func(s string) string {
		return i18n.Translate(tr, ctx.Context(), s)
	}

	return vue.UserComponent(
		v.VDataTableVirtual(
			web.Slot(
				v.VToolbar(
					v.VToolbarTitle(label).Style("font-size:1rem"),
					v.VSpacer(),
					v.VBtn("").Icon("mdi-plus").Flat(true).Density("compact").Attr("@click", `self.add()`),
				).Flat(true).Density("compact"),
			).Name("top"),
			web.Slot(
				keyComp,
			).Name("item.Name").Scope("{item, value}"),

			web.Slot(
				v.VTextarea().
					Density("compact").
					Attr("v-model", `item.Value`).
					HideDetails(true).
					Rows(1),
			).Name("item.Value").Scope("{item, value}"),
			web.Slot(
				v.VBtn("").Icon("mdi-delete").Density("compact").Attr("@click", `self.delete(item)`),
			).Name("item.actions").Scope("{item}"),
		).Density("compact").
			Attr(":items", fmt.Sprintf("form[%q]", fieldName)).
			Attr(":headers", fmt.Sprintf(`[{key:"Name",title:%q},{key:"Value",title:%q},{key:"actions", width:0}]`, t("Name"), t("Value"))),
	).Scope("self", vue.Var(`{
editedIndex: -1,
dialog: false,
dialogDelete: false,
editedItem: {
	Name: '',
	Value: ''
},
defaultItem: {
	Name: '',
	Value: ''
},
items: []
}`)).
		Setup(`({scope, computed, ref, watch, $scope, window, toRaw}) => {
form["`+fieldName+`"] = form["`+fieldName+`"] || []
const items = ref([])

scope.self.delete = (item) => {
	let removed = false
	form["`+fieldName+`"] = form["`+fieldName+`"].filter((e) => {
		if (removed) return true
		if (e.Name == item.Name && e.Value == item.Value) {
			removed = true
			return false
		}
		return true
	})
}

scope.self.add = () => {
	form["`+fieldName+`"].push({Name:'', Value:''})
}
}`).
		Assign("form", fieldName, items)

	// TODO: Parent Page, Galery, Sliders, Config.Address
}

type JSONMapFieldBuilder struct {
	nameKey    string
	nameLabel  func(ctx context.Context) string
	valueKey   string
	valueLabel func(ctx context.Context) string
	mb         *presets.ModelBuilder
	field      *presets.FieldBuilder
	keyComp    func(fctx *presets.FieldContext, nameKey string) h.HTMLComponent
}

func JSONMapField(mb *presets.ModelBuilder, field *presets.FieldBuilder, keyComp func(fctx *presets.FieldContext, nameKey string) h.HTMLComponent) *JSONMapFieldBuilder {
	return &JSONMapFieldBuilder{
		nameKey:  "Name",
		valueKey: "Value",
		mb:       mb,
		field:    field,
		keyComp:  keyComp,
	}
}

func L10nJSONMapField(mb *presets.ModelBuilder, field *presets.FieldBuilder, keyComp func(fctx *presets.FieldContext, nameKey string) h.HTMLComponent) *JSONMapFieldBuilder {
	return JSONMapField(mb, field, keyComp).
		NameKey("Location").
		NameLabelFunc(func(ctx context.Context) string {
			return l10n.MustGetMessages(ctx).Location
		})
}

func (b *JSONMapFieldBuilder) NameKey(name string) *JSONMapFieldBuilder {
	b.nameKey = name
	return b
}

func (b *JSONMapFieldBuilder) NameLabel(s string) *JSONMapFieldBuilder {
	b.nameLabel = func(ctx context.Context) string {
		return i18n.Translate(b.mb.Translator(), ctx, s)
	}
	return b
}

func (b *JSONMapFieldBuilder) NameLabelFunc(f func(ctx context.Context) string) *JSONMapFieldBuilder {
	b.nameLabel = f
	return b
}

func (b *JSONMapFieldBuilder) ValueKey(v string) *JSONMapFieldBuilder {
	b.valueKey = v
	return b
}

func (b *JSONMapFieldBuilder) ValueLabel(s string) *JSONMapFieldBuilder {
	b.valueLabel = func(ctx context.Context) string {
		return i18n.Translate(b.mb.Translator(), ctx, s)
	}
	return b
}

func (b *JSONMapFieldBuilder) ValueLabelFunc(f func(ctx context.Context) string) *JSONMapFieldBuilder {
	b.valueLabel = f
	return b
}

func (b *JSONMapFieldBuilder) Build() *presets.FieldBuilder {
	t := b.mb.Translator()

	return b.field.
		SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) (err error) {
			fieldValue := datatypes.NullJSONMap{}

			for k, v := range ctx.R.PostForm {
				if v[0] != "" && strings.HasPrefix(k, field.FormKey+"[") && strings.HasSuffix(k, "]."+b.nameKey) {
					kValue := strings.TrimSuffix(k, b.nameKey) + b.valueKey
					if fv := ctx.R.PostForm.Get(kValue); fv != "" {
						fieldValue[v[0]] = jsonMapCellParse(fv)
					}
				}
			}

			reflectutils.Set(obj, field.Name, fieldValue)

			return nil
		}).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			value, _ := field.Value().(datatypes.NullJSONMap)
			return b.input(t, field.EventContext, field.Label, field.FormKey, value, b.keyComp(field, b.nameKey))
		})
}

func (b *JSONMapFieldBuilder) input(tr i18n.Translator, ctx *web.EventContext, label, fieldName string, value datatypes.NullJSONMap, keyComp h.HTMLComponent) *vue.UserComponentBuilder {
	if value == nil {
		value = datatypes.NullJSONMap{}
	}

	var keys []string

	for k := range value {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	var items = make([]map[string]string, len(keys))
	for i, key := range keys {
		items[i] = map[string]string{
			b.nameKey:  key,
			b.valueKey: jsonMapCellString(value[key]),
		}
	}

	var (
		t = func(s string) string {
			return i18n.Translate(tr, ctx.Context(), s)
		}
		nameLabel  = b.nameLabel
		valueLabel = b.valueLabel
	)

	if nameLabel == nil {
		nameLabel = func(ctx context.Context) string {
			return t(b.nameKey)
		}
	}

	if valueLabel == nil {
		valueLabel = func(ctx context.Context) string {
			return t(b.valueKey)
		}
	}

	return vue.UserComponent(
		v.VDataTableVirtual(
			web.Slot(
				v.VToolbar(
					v.VToolbarTitle(label).Style("font-size:1rem"),
					v.VSpacer(),
					v.VBtn("").Icon("mdi-plus").Flat(true).Density("compact").Attr("@click", `self.add()`),
				).Flat(true).Density("compact"),
			).Name("top"),
			web.Slot(
				keyComp,
			).Name("item."+b.nameKey).Scope("{item, value}"),

			web.Slot(
				v.VTextarea().
					Density("compact").
					Attr("v-model", `item.`+b.valueKey).
					HideDetails(true).
					Rows(1),
			).Name("item."+b.valueKey).Scope("{item, value}"),
			web.Slot(
				v.VBtn("").Icon("mdi-delete").Density("compact").Attr("@click", `self.delete(item)`),
			).Name("item.actions").Scope("{item}"),
		).Density("compact").
			Attr(":items", fmt.Sprintf("form[%q]", fieldName)).
			Attr(":headers", fmt.Sprintf(`[{key:%q,title:%q},{key:%q,title:%q},{key:"actions", width:0}]`,
				b.nameKey,
				nameLabel(ctx.Context()),
				b.valueKey,
				valueLabel(ctx.Context()))),
	).Scope("self", vue.Var(`{
editedIndex: -1,
dialog: false,
dialogDelete: false,
editedItem: {
	`+b.nameKey+`: "",
	`+b.valueKey+`: ""
},
defaultItem: {
	`+b.nameKey+`: "",
	`+b.valueKey+`: ""
},
items: []
}`)).
		Setup(`({scope, computed, ref, watch, $scope, window, toRaw}) => {
form["`+fieldName+`"] = form["`+fieldName+`"] || []
const items = ref([])

scope.self.delete = (item) => {
	let removed = false
	form["`+fieldName+`"] = form["`+fieldName+`"].filter((e) => {
		if (removed) return true
		if (e.`+b.nameKey+` == item.`+b.nameKey+` && e.`+b.valueKey+` == item.`+b.valueKey+`) {
			removed = true
			return false
		}
		return true
	})
}

scope.self.add = () => {
	form["`+fieldName+`"].push({`+b.nameKey+`:"", `+b.valueKey+`:""})
}
}`).
		Assign("form", fieldName, items)

	// TODO: Parent Page, Galery, Sliders, Config.Address
}
