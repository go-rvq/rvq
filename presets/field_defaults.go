package presets

import (
	"fmt"
	"path/filepath"
	"reflect"
	"time"

	"github.com/qor5/admin/v3/reflect_utils"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	. "github.com/qor5/x/v3/ui/vuetify"
	"github.com/qor5/x/v3/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
)

type FieldDefaultBuilder struct {
	valType    reflect.Type
	compFunc   FieldComponentFunc
	setterFunc FieldSetterFunc
}

type FieldMode uint8

func (f FieldMode) Is(m ...FieldMode) bool {
	for _, other := range m {
		if f == other {
			return true
		}
	}
	return false
}

func (m FieldMode) Has(f FieldMode) bool {
	return m&f != 0
}
func (m FieldMode) HasAny(f ...FieldMode) bool {
	for _, f := range f {
		if m.Has(f) {
			return true
		}
	}
	return false
}

func (f FieldMode) Split() (modes []FieldMode) {
	if f.Has(NEW) {
		modes = append(modes, NEW)
	}
	if f.Has(EDIT) {
		modes = append(modes, EDIT)
	}
	if f.Has(DETAIL) {
		modes = append(modes, DETAIL)
	}
	if f.Has(LIST) {
		modes = append(modes, LIST)
	}
	return
}

const (
	NONE FieldMode = iota << 1
	LIST
	DETAIL
	NEW
	EDIT

	WRITE = NEW | EDIT
	ALL   = LIST | DETAIL | WRITE
)

type FieldModeStack []FieldMode

func (s FieldModeStack) Global() FieldMode {
	return s[0]
}

func (s FieldModeStack) Dot() FieldMode {
	return s[len(s)-1]
}

func (s FieldModeStack) DotStack() FieldModeStack {
	return FieldModeStack{s.Dot()}
}

func NewFieldDefault(t reflect.Type) (r *FieldDefaultBuilder) {
	r = &FieldDefaultBuilder{valType: t}
	return
}

func (b *FieldDefaultBuilder) ComponentFunc(v FieldComponentFunc) (r *FieldDefaultBuilder) {
	b.compFunc = v
	return b
}

func (b *FieldDefaultBuilder) SetterFunc(v FieldSetterFunc) (r *FieldDefaultBuilder) {
	b.setterFunc = v
	return b
}

var numberVals = []interface{}{
	int(0), int8(0), int16(0), int32(0), int64(0),
	uint(0), uint8(8), uint16(0), uint32(0), uint64(0),
	float32(0.0), float64(0.0),
}

var stringVals = []interface{}{
	"",
	[]rune(""),
	[]byte(""),
}

var timeVals = []interface{}{
	time.Now(),
	ptrTime(time.Now()),
}

type FieldDefaults struct {
	mode             FieldMode
	fieldTypes       []*FieldDefaultBuilder
	excludesPatterns []string
}

func NewFieldDefaults(t FieldMode) (r *FieldDefaults) {
	r = &FieldDefaults{
		mode: t,
	}

	r.builtInFieldTypes()
	return
}

func (b *FieldDefaults) FieldType(v interface{}) (r *FieldDefaultBuilder) {
	return b.fieldTypeByTypeOrCreate(reflect.TypeOf(v))
}

func (b *FieldDefaults) Exclude(patterns ...string) (r *FieldDefaults) {
	b.excludesPatterns = patterns
	return b
}

func (b *FieldDefaults) InspectFields(val interface{}, setup ...FieldSetuper) (r *FieldsBuilder) {
	_, fields := reflect_utils.UniqueFieldsOfReflectType(reflect.TypeOf(val))
	fieldBuilders := b.SetupFields(b.NewFieldBuilders(fields), setup...)

	return &FieldsBuilder{
		model:    val,
		fields:   fieldBuilders,
		defaults: b,
	}
}

func (b *FieldDefaults) NewFieldBuilders(fields reflect_utils.IndexableStructFields) (fbs FieldBuilders) {
	fbs = make(FieldBuilders, len(fields))
	var validCount int

	for _, f := range fields {
		if hasMatched(b.excludesPatterns, f.Name) {
			continue
		}
		fbs[validCount] = &FieldBuilder{
			rt: f.Type,
			NameLabel: NameLabel{
				name: f.Name,
			},
			mode:        ALL,
			structField: f,
		}
		validCount++
	}
	return fbs[:validCount]
}

func (b *FieldDefaults) SetupFields(fbs FieldBuilders, setupers ...FieldSetuper) (result FieldBuilders) {
	var (
		applyFt = func(f *FieldBuilder, ft *FieldDefaultBuilder) {
			if f.compFunc == nil {
				f.ComponentFunc(ft.compFunc)
			}

			if f.setterFunc == nil {
				f.SetterFunc(ft.setterFunc)
			}
		}
		withFt = func(f *FieldBuilder, rt reflect.Type) {
			if ft := b.fieldTypeByType(rt); ft != nil {
				applyFt(f, ft)
			}
		}
	)

	setuper := FieldSetupers(setupers)

	for _, f := range fbs {
		setuper.InitField(f)
		withFt(f, f.structField.Type)
		setuper.ConfigureField(f)
		result = append(result, f)
	}

	return result
}

func hasMatched(patterns []string, name string) bool {
	for _, p := range patterns {
		ok, err := filepath.Match(p, name)
		if err != nil {
			panic(err)
		}
		if ok {
			return true
		}
	}
	return false
}

func (b *FieldDefaults) String() string {
	var types []string
	for _, t := range b.fieldTypes {
		types = append(types, fmt.Sprint(t.valType))
	}
	return fmt.Sprintf("mode: %d, types %v", b.mode, types)
}

func (b *FieldDefaults) fieldTypeByType(tv reflect.Type) (r *FieldDefaultBuilder) {
	for _, ft := range b.fieldTypes {
		if ft.valType == tv {
			return ft
		}
	}
	return nil
}

func (b *FieldDefaults) fieldTypeByTypeOrCreate(tv reflect.Type) (r *FieldDefaultBuilder) {
	if r = b.fieldTypeByType(tv); r != nil {
		return
	}

	r = NewFieldDefault(tv)

	if b.mode == LIST {
		r.ComponentFunc(cfTextTd)
	} else {
		r.ComponentFunc(cfTextField)
	}
	b.fieldTypes = append(b.fieldTypes, r)
	return
}

func cfTextTd(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
	return h.Td(h.Text(field.StringValue()))
}

func cfCheckbox(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return VCheckbox().
		Attr(web.VField(field.FormKey, field.Value().(bool))...).
		Label(field.Label).
		ErrorMessages(field.Errors...).
		Disabled(field.ReadOnly)
}

func cfNumber(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return VTextField().
		Type("number").
		Variant(FieldVariantUnderlined).
		Attr(web.VField(field.FormKey, field.StringValue())...).
		Label(field.Label).
		ErrorMessages(field.Errors...).
		Disabled(field.ReadOnly)
}

func cfTime(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
	msgr := i18n.MustGetModuleMessages(ctx.R, CoreI18nModuleKey, Messages_en_US).(*Messages)
	val := ""
	if v := field.Value(); v != nil {
		switch vt := v.(type) {
		case time.Time:
			val = vt.Format("2006-01-02 15:04")
		case *time.Time:
			val = vt.Format("2006-01-02 15:04")
		default:
			panic(fmt.Sprintf("unknown time type: %T\n", v))
		}
	}
	return vuetifyx.VXDateTimePicker().
		Label(field.Label).
		Attr(web.VField(field.FormKey, val)...).
		Value(val).
		TimePickerProps(vuetifyx.TimePickerProps{
			Format:     "24hr",
			Scrollable: true,
		}).
		DialogWidth(640).
		ClearText(msgr.Clear).
		OkText(msgr.OK)
}

func cfTimeReadonly(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
	msgr := i18n.MustGetModuleMessages(ctx.R, CoreI18nModuleKey, Messages_en_US).(*Messages)
	val := ""
	if v := field.Value(); v != nil {
		switch vt := v.(type) {
		case time.Time:
			val = vt.Format(msgr.TimeFormats.DateTime)
		case *time.Time:
			val = vt.Format(msgr.TimeFormats.DateTime)
		default:
			panic(fmt.Sprintf("unknown time type: %T\n", v))
		}
	}
	return vuetifyx.VXReadonlyField().
		Label(field.Label).
		Value(val)
}

func cfTimeSetter(obj interface{}, field *FieldContext, ctx *web.EventContext) (err error) {
	v := ctx.R.Form.Get(field.FormKey)
	if v == "" {
		return reflectutils.Set(obj, field.Name, nil)
	}
	t, err := time.ParseInLocation("2006-01-02 15:04", v, time.Local)
	if err != nil {
		return err
	}
	return reflectutils.Set(obj, field.Name, t)
}

func cfTextField(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
	return VTextField().
		Type("text").
		Variant(FieldVariantUnderlined).
		Attr(web.VField(field.FormKey, field.StringValue())...).
		Label(field.Label).
		ErrorMessages(field.Errors...).
		Disabled(field.ReadOnly)
}

func CFReadonlyText(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
	return vuetifyx.VXReadonlyField().
		Label(field.Label).
		Value(field.StringValue())
}

func cfReadonlyCheckbox(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
	return vuetifyx.VXReadonlyField().
		Label(field.Label).
		Value(field.Value()).
		Checkbox(true)
}

func (b *FieldDefaults) builtInFieldTypes() {
	if b.mode == LIST {
		b.FieldType(true).
			ComponentFunc(cfTextTd)

		for _, v := range numberVals {
			b.FieldType(v).
				ComponentFunc(cfTextTd)
		}

		for _, v := range stringVals {
			b.FieldType(v).
				ComponentFunc(cfTextTd)
		}
		return
	}

	if b.mode == DETAIL {
		b.FieldType(true).
			ComponentFunc(cfReadonlyCheckbox)

		for _, v := range numberVals {
			b.FieldType(v).
				ComponentFunc(CFReadonlyText)
		}

		for _, v := range stringVals {
			b.FieldType(v).
				ComponentFunc(CFReadonlyText)
		}

		for _, v := range timeVals {
			b.FieldType(v).
				ComponentFunc(cfTimeReadonly)
		}
		return
	}

	b.FieldType(true).
		ComponentFunc(cfCheckbox)

	for _, v := range numberVals {
		b.FieldType(v).
			ComponentFunc(cfNumber)
	}

	for _, v := range stringVals {
		b.FieldType(v).
			ComponentFunc(cfTextField)
	}

	for _, v := range timeVals {
		b.FieldType(v).
			ComponentFunc(cfTime).
			SetterFunc(cfTimeSetter)
	}

	b.Exclude("ID")
	return
}
