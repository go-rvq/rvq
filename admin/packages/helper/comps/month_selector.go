package comps

import (
	"net/http"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

type MonthSelector struct {
	getLabels func(r *http.Request) [time.December + 1]string
}

func NewMonthSelector(getLabels func(r *http.Request) [time.December + 1]string) *MonthSelector {
	return &MonthSelector{getLabels: getLabels}
}

type mes struct {
	ID   uint
	Name string
}

func (m *MonthSelector) Of(mb *presets.ModelBuilder, fieldName string) *MonthSelectorField {
	return &MonthSelectorField{s: m, mb: mb, fieldName: fieldName, getLabels: m.getLabels}
}

func (m *MonthSelector) SetDefaults(b *presets.Builder) {
	b.FieldDefaults(presets.LIST).
		FieldType(time.January).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return m.CellValue(field, ctx)
		})
	b.FieldDefaults(presets.WRITE).
		FieldType(time.January).
		ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
			return m.Input(field, ctx)
		})
	b.FieldDefaults(presets.DETAIL).
		FieldType(time.January).
		ComponentFunc(m.Detail)
}

func (m *MonthSelector) CellValue(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	val := field.Value().(time.Month)
	if val == 0 {
		return h.Td()
	}
	return h.Td(h.Text(m.getLabels(ctx.R)[val]))
}

func (m *MonthSelector) Input(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	var (
		labels = m.getLabels(ctx.R)
		meses  = make([]mes, 12)
	)

	for i := 0; i < 12; i++ {
		meses[i] = mes{
			ID:   uint(i + 1),
			Name: labels[i+1],
		}
	}

	var (
		val      any
		selected = field.Value().(time.Month)
	)

	if selected > 0 {
		val = int(selected)
	}

	return VAutocomplete().
		Label(field.Label).
		Name(field.Name).
		ItemTitle("Name").
		ItemValue("ID").
		Clearable(true).
		Items(meses).
		Attr(web.VField(field.FormKey, val)...).
		ErrorMessages(field.Errors...).
		Disabled(field.ReadOnly)
}

func (m *MonthSelector) Detail(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	val := field.Value().(time.Month)
	if val < 1 {
		val = 0
	}
	field.ValueOverride = m.getLabels(ctx.R)[val]
	return presets.ReadonlyComponentFunc(field, ctx)
}

type MonthSelectorField struct {
	s         *MonthSelector
	mb        *presets.ModelBuilder
	fieldName string
	getLabels func(r *http.Request) [time.December + 1]string
}

func (m *MonthSelectorField) InputField(field *presets.FieldBuilder) *MonthSelectorField {
	field.ComponentFunc(m.s.Input)
	return m
}

func (m *MonthSelectorField) Input(fb *presets.FieldsBuilder) *MonthSelectorField {
	return m.InputField(fb.Field(m.fieldName))
}

func (m *MonthSelectorField) CellValueField(field *presets.FieldBuilder) *MonthSelectorField {
	field.ComponentFunc(m.s.Input)
	return m
}

func (m *MonthSelectorField) Text(fb *presets.FieldsBuilder) *MonthSelectorField {
	return m.TextField(fb.Field(m.fieldName))
}

func (m *MonthSelectorField) TextField(field *presets.FieldBuilder) *MonthSelectorField {
	field.ContextSetup(func(ctx *presets.FieldContext) {
		val := ctx.Value().(time.Month)
		if val < 1 {
			val = 0
		}
		ctx.ValueOverride = m.getLabels(ctx.EventContext.R)[val]
	})
	field.ComponentFunc(presets.ReadonlyComponentFunc)
	return m
}

func (m *MonthSelectorField) CellValue(fb *presets.FieldsBuilder) *MonthSelectorField {
	return m.CellValueField(fb.Field(m.fieldName))
}

func (m *MonthSelectorField) Editing() *MonthSelectorField {
	return m.Input(&m.mb.Editing().FieldsBuilder)
}

func (m *MonthSelectorField) Listing() *MonthSelectorField {
	return m.CellValue(&m.mb.Listing().FieldsBuilder)
}

func (m *MonthSelectorField) Detailing() *MonthSelectorField {
	return m.Text(&m.mb.Detailing().FieldsBuilder)
}

func (m *MonthSelectorField) All() *MonthSelectorField {
	m.Editing().Listing()
	if m.mb.HasDetailing() {
		m.Detailing()
	}
	return m
}
