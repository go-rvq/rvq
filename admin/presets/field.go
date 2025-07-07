package presets

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/reflect_utils"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/datafield"
	"github.com/go-rvq/rvq/web/zeroer"
	"github.com/go-rvq/rvq/x/i18n"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/iancoleman/strcase"
	"github.com/sunfmin/reflectutils"
)

func FieldPathIndex(i int) string {
	return fmt.Sprintf("[%d]", i)
}

type FieldPath []string

func (p *FieldPath) AppendIndex(i int) {
	*p = append((*p), FieldPathIndex(i))
}

func (p *FieldPath) Append(v ...string) {
	*p = append(*p, v...)
}

func (p FieldPath) NoIndex() (r FieldPath) {
	for _, s := range p {
		if s[0] != '[' {
			r = append(r, s)
		}
	}
	return r
}

func (p FieldPath) Fqn() string {
	return strings.Join(p, ".")
}

type FieldContext struct {
	ToComponentOptions *ToComponentOptions
	Parent             *FieldContext
	Mode               FieldModeStack
	Field              *FieldBuilder
	EventContext       *web.EventContext
	Obj                interface{}
	Name               string
	FormKey            string
	Path               FieldPath
	Label              string
	Hint               func() string
	Errors             []string
	SliceErrors        map[int][]string
	ModelInfo          *ModelInfo
	Nested             Nested
	Context            context.Context
	ReadOnly           bool
	Required           bool
	Disabled           bool
	ValueOverride      interface{}
	ComponentHandlers  []func(ctx *FieldContext, comp h.HTMLComponent) h.HTMLComponent
}

func (fc *FieldContext) Root() *FieldContext {
	for fc.Parent != nil {
		fc = fc.Parent
	}
	return fc
}

func (fc *FieldContext) ChildFieldFormKey(child string) (s string) {
	if len(fc.FormKey) == 0 {
		return child
	}
	return fc.FormKey + "." + child
}

func (fc *FieldContext) ComponentHandler(f ...func(ctx *FieldContext, comp h.HTMLComponent) h.HTMLComponent) *FieldContext {
	fc.ComponentHandlers = append(fc.ComponentHandlers, f...)
	return fc
}

func (fc *FieldContext) StringValue() (r string) {
	v := fc.Value()

	if v == nil {
		return ""
	}

	if s, _ := v.(FieldStringer); s != nil {
		return s.FieldString(fc)
	}

	return ToStringContext(fc.EventContext, v)
}

func (fc *FieldContext) RawValue() (r interface{}) {
	fieldName := fc.Name
	return reflectutils.MustGet(fc.Obj, fieldName)
}

func (fc *FieldContext) Value() (r interface{}) {
	if fc.ValueOverride != nil {
		return fc.ValueOverride
	}
	return fc.RawValue()
}

func (fc *FieldContext) SetContextValue(key, value interface{}) {
	if fc.Context == nil {
		fc.Context = context.WithValue(context.Background(), key, value)
	} else {
		fc.Context = context.WithValue(fc.Context, key, value)
	}
}

func (fc *FieldContext) ContextValue(key interface{}) (r interface{}) {
	if fc.Context == nil {
		return
	}
	return fc.Context.Value(key)
}

func (fc *FieldContext) Error(err error) *FieldContext {
	if err != nil {
		fc.Errors = append(fc.Errors, err.Error())
	}
	return fc
}

func (fc *FieldContext) AddError(err string) *FieldContext {
	fc.Errors = append(fc.Errors, err)
	return fc
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
	mode             FieldMode
	structField      *reflect_utils.IndexableStructField
	compFunc         FieldComponentFunc
	setterFunc       FieldSetterFunc
	context          context.Context
	rt               reflect.Type
	nested           Nested
	disabled         bool
	enabled          func(ctx *FieldContext) bool
	Setup            FieldContextSetups
	ToComponentSetup FieldContextSetups
	Validators       FieldValidators
	ValueFormatters  FieldValueFormatters
	defaultValuer    func()
	audited          bool
	hint             bool

	datafield.DataField[*FieldBuilder]
}

func (n *FieldBuilder) SetLabelKey(labelKey string) *FieldBuilder {
	n.labelKey = labelKey
	return n
}

func (n *FieldBuilder) SetI18nLabel(i18nLabel func(ctx context.Context) string) *FieldBuilder {
	n.i18nLabel = i18nLabel
	return n
}

func (n *FieldBuilder) SetHiddenLabel(hiddenLabel bool) *FieldBuilder {
	n.hiddenLabel = hiddenLabel
	return n
}

func (n *FieldBuilder) SetHint(hint bool) *FieldBuilder {
	n.hint = hint
	return n
}

func (n *FieldBuilder) Hint() bool {
	return n.hint
}

func (b *FieldBuilder) String() string {
	return b.name
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

func (b *FieldBuilder) Enabled() func(ctx *FieldContext) bool {
	return b.enabled
}

func (b *FieldBuilder) SetEnabled(enabled func(ctx *FieldContext) bool) *FieldBuilder {
	b.enabled = enabled
	return b
}
func (b *FieldBuilder) Enable(v bool) *FieldBuilder {
	b.enabled = func(*FieldContext) bool {
		return v
	}
	return b
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

func (b *FieldBuilder) SetDisabled(v bool) *FieldBuilder {
	b.disabled = v
	return b
}

func (b *FieldBuilder) Disabled() bool {
	return b.disabled
}

func (b *FieldBuilder) GetCompFunc() FieldComponentFunc {
	return b.compFunc
}

func (b *FieldBuilder) Label(v string) (r *FieldBuilder) {
	b.label = v
	return b
}

func (b *FieldBuilder) Audited() bool {
	return b.audited
}

func (b *FieldBuilder) SetAudited(audited bool) *FieldBuilder {
	b.audited = audited
	return b
}

func (b *FieldBuilder) SetData(key, value any) *FieldBuilder {
	b.DataField.SetData(key, value)
	return b
}

func (b FieldBuilder) Clone() *FieldBuilder {
	b.Setup = append(FieldContextSetups{}, b.Setup...)
	b.ToComponentSetup = append(FieldContextSetups{}, b.ToComponentSetup...)
	b.ValueFormatters = append(FieldValueFormatters{}, b.ValueFormatters...)
	b.Validators = append(FieldValidators{}, b.Validators...)
	b.DataField = b.DataField.Clone()
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

func (b *FieldBuilder) DisableZeroRender() *FieldBuilder {
	return b.WrapComponentFunc(func(old FieldComponentFunc) FieldComponentFunc {
		return func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
			if zeroer.IsZero(field.Value()) {
				return nil
			}
			return old(field, ctx)
		}
	})
}

func (b *FieldBuilder) SetterFunc(v FieldSetterFunc) (r *FieldBuilder) {
	b.setterFunc = v
	return b
}

func (b *FieldBuilder) GetSetterFunc() FieldSetterFunc {
	return b.setterFunc
}

func (b *FieldBuilder) WithContextValue(key interface{}, val interface{}) (r *FieldBuilder) {
	if b.context == nil {
		b.context = context.Background()
	}
	b.context = context.WithValue(b.context, key, val)
	return b
}

func (b *FieldBuilder) ContextLabel(info *ModelInfo, ctx context.Context, fallback ...func(ctx context.Context, nameLabel NameLabel) string) (label string) {
	if b.hiddenLabel {
		return ""
	}

	label = b.labelKey

	if label == "" {
		if b.i18nLabel != nil {
			return b.i18nLabel(ctx)
		}

		if b.label != "" {
			for _, f := range fallback {
				if label = f(ctx, NameLabel{name: label}); label != "" {
					return label
				}
			}
			label = b.label
		} else {
			for _, f := range fallback {
				if label = f(ctx, b.NameLabel); label != "" {
					return label
				}
			}

			label = b.name
		}
	}

	if info != nil {
		label = i18n.Translate(info.mb.FieldTranslator(), ctx, label)
	} else {
		msgr := MustGetMessages(ctx)
		label = msgr.Common.Get(label)
	}

	if label == "" {
		label = HumanizeString(b.name)
	}

	return label
}

func (b *FieldBuilder) DefaultContextLabel(ctx context.Context) string {
	return b.ContextLabel(nil, ctx)
}

func (b *FieldBuilder) ContextHint(info *ModelInfo, ctx context.Context) string {
	if info != nil {
		return i18n.TranslateD(info.mb.FieldHintTranslator(), nil, ctx, b.name+"_Hint")
	}

	msgr := MustGetMessages(ctx)
	return msgr.Common.Get(b.name)
}

func (b *FieldBuilder) DefaultContextHint(ctx context.Context) string {
	return b.ContextHint(nil, ctx)
}

func (b *FieldBuilder) SetupFunc(fc FieldContextSetup) *FieldBuilder {
	b.Setup.Add(fc)
	return b
}

func (b *FieldBuilder) ToComponent(ctx *FieldContext) (comp h.HTMLComponent) {
	panics := true
	defer func() {
		if panics {
			if r := recover(); r != nil {
				comp = FieldComponentWrapper(func(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
					return v.VAlert().
						Type(v.TypeError).
						Density(v.DensityCompact).
						Variant(v.VariantTonal).
						Text(fmt.Sprintf("%v", r))
				})(ctx, ctx.EventContext)
			}
		}
	}()
	comp = b.compFunc(ctx, ctx.EventContext)
	for _, f := range ctx.ComponentHandlers {
		comp = f(ctx, comp)
	}
	panics = false
	return
}

func (b *FieldBuilder) ContextSetup(f func(ctx *FieldContext)) *FieldBuilder {
	b.ToComponentSetup = append(b.ToComponentSetup, f)
	return b
}

type FieldBuilders []*FieldBuilder

func (b FieldBuilders) Get(name string) *FieldBuilder {
	for _, fb := range b {
		if fb.name == name {
			return fb
		}
	}
	return nil
}

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

func (b FieldBuilders) EachHavesComponent(cb func(fb *FieldBuilder) bool) {
	for _, fb := range b {
		if fb.compFunc != nil {
			if !cb(fb) {
				break
			}
		}
	}
}

func (b FieldBuilders) FieldsFromLayout(layout []interface{}, filter ...FieldFilter) (res FieldBuilders) {
	for _, f := range b.FieldsGenFromLayout(layout, filter...) {
		res = append(res, f)
	}
	return
}

func (b FieldBuilders) FieldsGenFromLayout(layout []interface{}, filter ...FieldFilter) func(func(int, *FieldBuilder) bool) {
	accept := FieldFilters(filter).Accept

	return func(yield_ func(int, *FieldBuilder) bool) {
		var (
			i     int
			yield = func(f *FieldBuilder) bool {
				if !accept(f) {
					return true
				}
				if !yield_(i, f) {
					return false
				}
				i++
				return true
			}
		)
		for _, iv := range layout {
			switch t := iv.(type) {
			case string:
				if f := b.Get(t); f != nil {
					if !yield(f) {
						return
					}
				}
			case []string:
				for _, n := range t {
					if f := b.Get(n); f != nil {
						if !yield(f) {
							return
						}
					}
				}
			case *FieldsSection:
				for _, row := range t.Rows {
					for _, n := range row {
						if f := b.Get(n); f != nil {
							if !yield(f) {
								return
							}
						}
					}
				}
			default:
				panic("unknown fields layout, must be string/[]string/*FieldsSection")
			}
		}
	}
}
