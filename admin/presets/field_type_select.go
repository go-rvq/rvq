package presets

import (
	"net/url"

	"github.com/go-playground/form"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

type (
	SelectConfigor interface {
		AvailableKeys(ctx *FieldContext) []string
		KeyLabels(ctx *FieldContext, key []string) []string
		KeyLabelsAndHints(ctx *FieldContext, key []string) [][2]string
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
	AvailableKeysFunc     func(ctx *FieldContext) []string
	KeyLabelsFunc         func(ctx *FieldContext, key []string) []string
	KeyLabelsAndHintsFunc func(ctx *FieldContext, key []string) [][2]string
	SelectedKeyFunc       func(ctx *FieldContext) string
	SetSelectedKeyFunc    func(ctx *FieldContext, key string) (err error)
	ToStringWrapFunc      func(ctx *FieldContext, key, label string) (newLabel string)
}

func (s *SelectConfig) AvailableKeys(ctx *FieldContext) []string {
	return s.AvailableKeysFunc(ctx)
}

func (s *SelectConfig) KeyLabels(ctx *FieldContext, key []string) (r []string) {
	if s.KeyLabelsFunc == nil {
		labelsAndHints := s.KeyLabelsAndHintsFunc(ctx, key)
		r = make([]string, len(labelsAndHints))
		for i, label := range labelsAndHints {
			r[i] = label[0]
		}
		return
	}
	if s.KeyLabelsFunc == nil {
		return key
	}
	return s.KeyLabelsFunc(ctx, key)
}

func (s *SelectConfig) KeyLabelsAndHints(ctx *FieldContext, key []string) (r [][2]string) {
	if s.KeyLabelsAndHintsFunc == nil {
		var labels []string
		if s.KeyLabelsFunc == nil {
			labels = key
		} else {
			labels = s.KeyLabelsFunc(ctx, key)
		}
		r = make([][2]string, len(labels))
		for i, label := range labels {
			r[i] = [2]string{label, ""}
		}
		return
	}
	return s.KeyLabelsAndHintsFunc(ctx, key)
}

func (s *SelectConfig) SelectedKey(ctx *FieldContext) string {
	if s.SelectedKeyFunc != nil {
		return s.SelectedKeyFunc(ctx)
	}
	return ctx.StringValue()
}

func (s *SelectConfig) SetSelectedKey(ctx *FieldContext, key string) (err error) {
	if s.SetSelectedKeyFunc != nil {
		return s.SetSelectedKeyFunc(ctx, key)
	}
	d := form.NewDecoder()
	return d.Decode(ctx.Obj, url.Values{
		ctx.Name: []string{key},
	})
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

	for i, labelAndHint := range cfg.KeyLabelsAndHints(field, keys) {
		options[i] = map[string]any{
			"value": keys[i],
			"title": labelAndHint[0],
			"hint":  labelAndHint[1],
		}
	}

	sel := v.VSelect().
		Label(field.Label).
		Items(options).
		Density(v.DensityComfortable).
		ItemValue(`value`).
		ItemTitle(`title`).
		ErrorMessages(field.Errors...).
		Attr("v-model", "fieldValue.value")

	if field.Field.hint {
		sel.Hint(field.HintLoader()).PersistentHint(true)
	}

	return vue.FormField(
		sel,
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

	labels[0] = cfg.ToStringWrap(field, key, labels[0])

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

func ConfigureSelectField(mb *ModelBuilder, name string, mode FieldMode, cfg SelectConfigor) {
	if mode == 0 {
		mode = ALL
	}
	if mode.Has(LIST) {
		mb.Listing().Field(name).
			SetData(selectFieldDataConfigor, cfg).
			ComponentFunc(func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
				comp := SelectToStringComponentFunc(field, ctx)
				if comp != nil {
					comp = h.Td(comp)
				}
				return comp
			})
	}

	if mode.Has(DETAIL) && mb.hasDetailing {
		mb.detailing.Field(name).
			SetData(selectFieldDataConfigor, cfg).
			ComponentFunc(SelectReadComponentFunc)
	}

	setter := func(obj interface{}, field *FieldContext, ctx *web.EventContext) (err error) {
		return cfg.SetSelectedKey(field, field.EventContext.R.Form.Get(field.FormKey))
	}

	if mb.Editing().HasCreatingBuilder() {
		if mode.Has(NEW) {
			mb.Editing().CreatingBuilder().Field(name).
				SetData(selectFieldDataConfigor, cfg).
				ComponentFunc(func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
					if mode.Has(field.Mode.Dot()) {
						return SelectWriteComponentFunc(field, ctx)
					}
					return nil
				}).
				SetterFunc(setter)
		}

		if mode.Has(EDIT) {
			mb.Editing().Field(name).
				SetData(selectFieldDataConfigor, cfg).
				ComponentFunc(func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
					if mode.Has(field.Mode.Dot()) {
						return SelectWriteComponentFunc(field, ctx)
					}
					return nil
				}).
				SetterFunc(setter)
		}
	} else if mode.Has(EDIT) || mode.Has(NEW) {
		mb.Editing().Field(name).
			SetData(selectFieldDataConfigor, cfg).
			ComponentFunc(func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
				if mode.Has(field.Mode.Dot()) {
					return SelectWriteComponentFunc(field, ctx)
				}
				return nil
			}).
			SetterFunc(setter)
	}
}
