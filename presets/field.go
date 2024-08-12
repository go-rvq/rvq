package presets

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/qor5/admin/v3/reflect_utils"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	"github.com/sunfmin/reflectutils"
)

type FieldContext struct {
	Mode          FieldModeStack
	Field         *FieldBuilder
	EventContext  *web.EventContext
	Obj           interface{}
	Name          string
	FormKey       string
	Label         string
	Errors        []string
	ModelInfo     *ModelInfo
	Nested        *NestedFieldBuilder
	Context       context.Context
	ReadOnly      bool
	Required      bool
	Disabled      bool
	ValueOverride interface{}
}

func (fc *FieldContext) StringValue(obj interface{}) (r string) {
	val := fc.Value(obj)
	switch vt := val.(type) {
	case []rune:
		return string(vt)
	case []byte:
		return string(vt)
	case time.Time:
		return vt.Format("2006-01-02 15:04:05")
	case *time.Time:
		if vt == nil {
			return ""
		}
		return vt.Format("2006-01-02 15:04:05")
	}
	return fmt.Sprint(val)
}

func (fc *FieldContext) RawValue(obj interface{}) (r interface{}) {
	fieldName := fc.Name
	return reflectutils.MustGet(obj, fieldName)
}

func (fc *FieldContext) Value(obj interface{}) (r interface{}) {
	if fc.ValueOverride != nil {
		return fc.ValueOverride
	}
	return fc.RawValue(obj)
}

func (fc *FieldContext) ContextValue(key interface{}) (r interface{}) {
	if fc.Context == nil {
		return
	}
	return fc.Context.Value(key)
}

type FieldContextSetup func(ctx *FieldContext)

type FieldContextSetups []FieldContextSetup

func (f *FieldContextSetups) Add(fc FieldContextSetup) {
	*f = append(*f, fc)
}

func (f FieldContextSetups) Setup(ctx *FieldContext) {
	for _, setup := range f {
		setup(ctx)
	}
}

type FieldBuilder struct {
	NameLabel
	mode                FieldMode
	structField         *reflect_utils.IndexableStructField
	compFunc            FieldComponentFunc
	setterFunc          FieldSetterFunc
	context             context.Context
	rt                  reflect.Type
	nestedFieldsBuilder *NestedFieldBuilder
	enabled             func(ctx *FieldContext) bool
	data                map[any]any
	Setup               FieldContextSetups
	ToComponentSetup    FieldContextSetups
	Validators          FieldValidators
	ValueFormatters     FieldValueFormatters
	defaultValuer       func()
	audited             bool
}

func (b *FieldBuilder) ColumnName() string {
	return strcase.ToSnake(b.NameLabel.name)
}

func (b *FieldBuilder) Mode() FieldMode {
	return b.mode
}

func (b *FieldBuilder) SetMode(mode FieldMode) *FieldBuilder {
	b.mode = mode
	return b
}

func (b *FieldBuilder) StructField() *reflect_utils.IndexableStructField {
	return b.structField
}

func (b *FieldBuilder) SetData(key, value any) *FieldBuilder {
	if b.data == nil {
		b.data = make(map[any]any)
	}
	b.data[key] = value
	return b
}

func (b *FieldBuilder) GetData(key any) any {
	if b.data == nil {
		return nil
	}
	return b.data[key]
}

func (b *FieldBuilder) Enabled() func(ctx *FieldContext) bool {
	return b.enabled
}

func (b *FieldBuilder) SetEnabled(enabled func(ctx *FieldContext) bool) {
	b.enabled = enabled
}

func (b *FieldBuilder) WrapEnabled(do func(old func(ctx *FieldContext) bool) func(ctx *FieldContext) bool) *FieldBuilder {
	old := b.enabled
	if old == nil {
		old = func(ctx *FieldContext) bool {
			return true
		}
	}
	b.enabled = do(old)
	return b
}

func (b *FieldBuilder) IsEnabled(ctx *FieldContext) bool {
	if b.enabled != nil {
		return b.enabled(ctx)
	}
	return true
}

func (b *FieldBuilder) GetCompFunc() FieldComponentFunc {
	return b.compFunc
}

func (b *FieldBuilder) Label(v string) (r *FieldBuilder) {
	b.label = v
	return b
}

func (b *FieldBuilder) SetHiddenLabel(hiddenLabel bool) *FieldBuilder {
	b.hiddenLabel = hiddenLabel
	return b
}

func (b *FieldBuilder) Audited() bool {
	return b.audited
}

func (b *FieldBuilder) SetAudited(audited bool) *FieldBuilder {
	b.audited = audited
	return b
}

func (b FieldBuilder) Clone() *FieldBuilder {
	b.Setup = append(FieldContextSetups{}, b.Setup...)
	b.ToComponentSetup = append(FieldContextSetups{}, b.ToComponentSetup...)
	b.ValueFormatters = append(FieldValueFormatters{}, b.ValueFormatters...)
	b.Validators = append(FieldValidators{}, b.Validators...)

	if b.data != nil {
		data := make(map[any]any, len(b.data))
		for k, v := range b.data {
			data[k] = v
		}
		b.data = data
	}
	return &b
}

func (b *FieldBuilder) ComponentFunc(v FieldComponentFunc) (r *FieldBuilder) {
	if v == nil {
		panic("value required")
	}
	b.compFunc = v
	return b
}

func (b *FieldBuilder) WrapComponentFunc(v func(old FieldComponentFunc) FieldComponentFunc) (r *FieldBuilder) {
	if v == nil {
		panic("value required")
	}
	b.compFunc = v(b.compFunc)
	return b
}

func (b *FieldBuilder) SetterFunc(v FieldSetterFunc) (r *FieldBuilder) {
	b.setterFunc = v
	return b
}

func (b *FieldBuilder) WithContextValue(key interface{}, val interface{}) (r *FieldBuilder) {
	if b.context == nil {
		b.context = context.Background()
	}
	b.context = context.WithValue(b.context, key, val)
	return b
}

func (b *FieldBuilder) RequestLabel(fb *FieldsBuilder, info *ModelInfo, r *http.Request) string {
	if b.hiddenLabel {
		return ""
	}

	var label = b.labelKey
	if label == "" {
		if b.i18nLabel != nil {
			return b.i18nLabel(r)
		}

		msgr := MustGetMessages(r)
		if label = msgr.CommonFieldLabels.Get(b.name); label != "" {
			return label
		}

		label = fb.getLabel(b.NameLabel)
	}

	if info != nil {
		return i18n.PT(r, ModelsI18nModuleKey, info.Label(), label)
	}
	return label
}

type FieldBuilders []*FieldBuilder

func (b FieldBuilders) Len() int {
	return len(b)
}

func (b FieldBuilders) HasMode(mode FieldMode, cb ...func(fb *FieldBuilder) bool) (r FieldBuilders) {
	b.Filter(func(fb *FieldBuilder) bool {
		if fb.mode.Has(mode) {
			for _, f := range cb {
				if !f(fb) {
					return true
				}
			}
			r = append(r, fb)
		}
		return true
	})
	return
}

func (b FieldBuilders) Filter(f func(fb *FieldBuilder) bool) (ret FieldBuilders) {
	for _, fb := range b {
		if f(fb) {
			ret = append(ret, fb)
		}
	}
	return ret
}

func (b FieldBuilders) FirstFilter(f func(fb *FieldBuilder) bool) *FieldBuilder {
	for _, fb := range b {
		if f(fb) {
			return fb
		}
	}
	return nil
}

func (b FieldBuilders) First() *FieldBuilder {
	if len(b) == 0 {
		return nil
	}
	return b[0]
}

func (b FieldBuilders) Last() *FieldBuilder {
	if len(b) == 0 {
		return nil
	}
	return b[len(b)-1]
}

func (b FieldBuilders) Renderable() (r FieldBuilders) {
	for _, fb := range b {
		if fb.compFunc != nil {
			r = append(r, fb)
		}
	}
	return r
}

func (b FieldBuilders) EachHavesComponent(cb func(fb *FieldBuilder) bool) {
	for _, fb := range b {
		if fb.compFunc != nil {
			if !cb(fb) {
				break
			}
		}
	}
}
