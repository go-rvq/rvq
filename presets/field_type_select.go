package presets

import (
	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/vue"
	v "github.com/qor5/x/v3/ui/vuetify"
	h "github.com/theplant/htmlgo"
)

type (
	SelectConfigor interface {
		AvailableKeys(ctx *FieldContext) []string
		KeyLabels(ctx *FieldContext, key []string) []string
		SelectedKey(ctx *FieldContext) string
		SetSelectedKey(ctx *FieldContext, key string) (err error)
		ToStringWrap(ctx *FieldContext, key, label string) (newLabel string)
	}

	GetSelectConfigor interface {
		SelectConfigor() SelectConfigor
	}
)

type SelectConfig struct {
	// AvailableKeys list of pairs of key and label [[key, label], ...]
	AvailableKeysFunc  func(ctx *FieldContext) []string
	KeyLabelsFunc      func(ctx *FieldContext, key []string) []string
	SelectedKeyFunc    func(ctx *FieldContext) string
	SetSelectedKeyFunc func(ctx *FieldContext, key string) (err error)
	ToStringWrapFunc   func(ctx *FieldContext, key, label string) (newLabel string)
}

func (s *SelectConfig) AvailableKeys(ctx *FieldContext) []string {
	return s.AvailableKeysFunc(ctx)
}

func (s *SelectConfig) KeyLabels(ctx *FieldContext, key []string) []string {
	return s.KeyLabelsFunc(ctx, key)
}

func (s *SelectConfig) SelectedKey(ctx *FieldContext) string {
	return s.SelectedKeyFunc(ctx)
}

func (s *SelectConfig) SetSelectedKey(ctx *FieldContext, key string) (err error) {
	return s.SetSelectedKeyFunc(ctx, key)
}

func (s *SelectConfig) ToStringWrap(ctx *FieldContext, key, label string) (newLabel string) {
	if s.ToStringWrapFunc != nil {
		return s.ToStringWrapFunc(ctx, key, label)
	}
	return label
}

type FieldSelectConfigor func() *SelectConfig

func selectorConfigorFromFieldContext(field *FieldContext) SelectConfigor {
	switch t := field.RawValue().(type) {
	case GetSelectConfigor:
		return t.SelectConfigor()
	case SelectConfigor:
		return t
	default:
		if cfg, _ := field.Field.GetData(selectFieldDataConfigor).(SelectConfigor); cfg != nil {
			return cfg
		}
		panic("not a SelectConfigor")
	}
}

func SelectWriteComponentFunc(field *FieldContext, _ *web.EventContext) (comp h.HTMLComponent) {
	var (
		cfg      = selectorConfigorFromFieldContext(field)
		keys     = cfg.AvailableKeys(field)
		selected = cfg.SelectedKey(field)
		options  = make([]any, len(keys))
	)

	for i, label := range cfg.KeyLabels(field, keys) {
		options[i] = map[string]any{"value": keys[i], "title": label}
	}

	return vue.FormField(
		v.VSelect().
			Label(field.Label).
			Items(options).
			Density(v.DensityCompact).
			ItemValue(`value`).
			ItemTitle(`title`).
			Attr("v-model", "fieldValue.value").
			HideDetails(true),
	).
		Value(field.FormKey, selected)
}

func SelectReadComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	var (
		comp        h.HTMLComponents
		cfg         = selectorConfigorFromFieldContext(field)
		selectedKey = cfg.SelectedKey(field)
	)

	if cfg == nil {
		cfg = field.Value().(SelectConfigor)
	}

	if len(selectedKey) == 0 {
		return nil
	}

	if field.Label != "" {
		comp = append(comp, v.VLabel(h.Text(field.Label)))
	}

	comp = append(comp, h.Div(h.Text(cfg.KeyLabels(field, []string{selectedKey})[0])).Class("pt-1 mb-3"))

	return comp
}

func SelectToStringComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	var (
		cfg = selectorConfigorFromFieldContext(field)
		key = cfg.SelectedKey(field)
	)

	if len(key) == 0 {
		return h.Text("")
	}

	labels := cfg.KeyLabels(field, []string{key})

	if cfg.ToStringWrap != nil {
		labels[0] = cfg.ToStringWrap(field, key, labels[0])
	}

	return h.Text(labels[0])
}

type selectFieldDataType uint8

const selectFieldDataConfigor selectFieldDataType = iota

func RegisterSelectType(p *Builder, typ any, cfg *SelectConfig) {
	p.FieldDefaults(LIST).
		FieldType(typ).
		SetData(selectFieldDataConfigor, cfg).
		ComponentFunc(func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
			comp := SelectToStringComponentFunc(field, ctx)
			if comp != nil {
				comp = h.Td(comp)
			}
			return comp
		})

	p.FieldDefaults(WRITE).
		FieldType(typ).
		SetData(selectFieldDataConfigor, cfg).
		ComponentFunc(SelectWriteComponentFunc)

	p.FieldDefaults(DETAIL).
		FieldType(typ).
		SetData(selectFieldDataConfigor, cfg).
		ComponentFunc(SelectReadComponentFunc)
}
