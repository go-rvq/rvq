package presets

import (
	"strings"

	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type MultiSelectConfig struct {
	// AvailableKeys list of pairs of key and label [[key, label], ...]
	AvailableKeys   func(ctx *FieldContext) []string
	KeyLabels       func(ctx *FieldContext, keys []string) []string
	SelectedKeys    func(ctx *FieldContext) []string
	SetSelectedKeys func(ctx *FieldContext, keys []string) (err error)
	ToStringWrap    func(ctx *FieldContext, keys, labels []string) (newLabels []string)
	ToStringSep     string
}

type FieldMultiSelectConfigor func() *MultiSelectConfig

func MultiSelectWriteComponentFunc(field *FieldContext, _ *web.EventContext) (comp h.HTMLComponent) {
	var (
		cfg      = field.Field.GetData(multiSelectFieldDataConfigor).(*MultiSelectConfig)
		keys     = cfg.AvailableKeys(field)
		selected = cfg.SelectedKeys(field)
		data     = make([]any, len(keys))
	)

	for i, label := range cfg.KeyLabels(field, keys) {
		data[i] = map[string]any{"key": keys[i], "label": label}
	}

	comp = h.Div(
		VDataTable(
			web.Slot().Name("bottom"),
		).
			FormField(field.FormKey, selected).
			Attr(
				":headers", `[{key:"label"}]`,
				":item-value", "item => item.key",
				":items", data,
				"show-select", true,
				":items-per-page", len(keys),
			),
	).Class("v-input__details")

	if field.Label != "" {
		comp = h.HTMLComponents{
			h.Div(
				h.Label(field.Label).Class("v-label theme--light"),
			),
			comp,
		}
	}

	comp = h.Div(comp)

	return comp
}

func MultiSelectReadComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	var (
		comp h.HTMLComponents
		cfg  = field.Field.GetData(multiSelectFieldDataConfigor).(*MultiSelectConfig)
		keys = cfg.SelectedKeys(field)
	)
	if len(keys) == 0 {
		return nil
	}

	if field.Label != "" {
		comp = append(comp, VLabel(h.Text(field.Label)))
	}

	comp = append(comp, h.Div(h.Text(strings.Join(cfg.KeyLabels(field, keys), cfg.ToStringSep))).Class("pt-1 mb-3"))

	return comp
}

func MultiSelectToStringComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	var (
		cfg  = field.Field.GetData(multiSelectFieldDataConfigor).(*MultiSelectConfig)
		keys = cfg.SelectedKeys(field)
	)

	if len(keys) == 0 {
		return nil
	}

	labels := cfg.KeyLabels(field, keys)

	if cfg.ToStringWrap != nil {
		labels = cfg.ToStringWrap(field, keys, labels)
	}

	return h.Text(strings.Join(labels, cfg.ToStringSep))
}

type multiSelectFieldDataType uint8

const multiSelectFieldDataConfigor multiSelectFieldDataType = iota

func RegisterMultiSelectType(p *Builder, typ any, cfg *MultiSelectConfig) {
	if cfg.ToStringSep == "" {
		cfg.ToStringSep = ", "
	}

	p.FieldDefaults(LIST).
		FieldType(typ).
		SetData(multiSelectFieldDataConfigor, cfg).
		ComponentFunc(func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
			comp := MultiSelectToStringComponentFunc(field, ctx)
			if comp != nil {
				comp = h.Td(comp)
			}
			return comp
		})

	p.FieldDefaults(WRITE).
		FieldType(typ).
		SetData(multiSelectFieldDataConfigor, cfg).
		ComponentFunc(MultiSelectWriteComponentFunc)

	p.FieldDefaults(DETAIL).
		FieldType(typ).
		SetData(multiSelectFieldDataConfigor, cfg).
		ComponentFunc(MultiSelectReadComponentFunc)
}
