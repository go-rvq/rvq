package helper

import (
	"context"
	"reflect"
	"strings"

	"github.com/gad-lang/gad"
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/helper/tagparser"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/tag"
	"github.com/go-rvq/rvq/web/zeroer"
	"github.com/sunfmin/reflectutils"
)

type (
	KeyValueArray gad.KeyValueArray
)

func (kva KeyValueArray) HasFlag(key gad.Object) (ok bool) {
	for _, value := range kva {
		if value.K.Equal(key) {
			switch vt := value.V.(type) {
			case gad.Bool:
				ok = bool(vt)
			case gad.Flag:
				ok = bool(vt)
			}
		}
	}
	return
}

var AdminTagKey = reflect.TypeOf((*AdminTag)(nil)).Elem()

func TagOf(fb *presets.FieldBuilder) *AdminTag {
	if data := fb.GetData(AdminTagKey); data != nil {
		return data.(*AdminTag)
	}
	return nil
}

type AdminTag struct {
	Mode                  presets.FieldMode
	Required              presets.FieldMode
	RequiredBy            string
	ReadOnly              bool
	ListHint              string
	WriteHint             string
	DetailHint            string
	EditorJS              bool
	TipTap                bool
	EditComponentHandlers []func(ctx *presets.FieldContext, comp h.HTMLComponent) h.HTMLComponent
}

func (t *AdminTag) Parse(sf *reflect.StructField, s string) (valid bool) {
	if s == "" {
		*t = AdminTag{}
		return true
	}
	typ := sf.Type

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.Slice, reflect.Map:
		t.Mode = presets.WRITE
	case reflect.Func, reflect.Interface, reflect.Chan:
		return
	default:
		t.Mode = presets.ALL
	}

	valid = true

	if s == "" {
		return
	}

	var (
		kva gad.KeyValueArray
		vm  *gad.VM
		err error
	)

	if vm, kva, err = tagparser.Parse(s); err != nil {
		panic(err)
	}

	na := gad.NewNamedArgs(kva)
	if v := na.GetValueOrNil("mode"); v != nil {
		switch vt := v.(type) {
		case gad.Uint:
			t.Mode = presets.FieldMode(vt)
		}
	}
	if v := na.GetValueOrNil("required"); v != nil {
		switch vt := v.(type) {
		case gad.Uint:
			t.Required = presets.FieldMode(vt)
		default:
			if !v.IsFalsy() {
				t.Required = presets.WRITE
			}
		}
	}
	if v := na.GetValueOrNil("required_by"); v != nil {
		t.RequiredBy = v.ToString()
	}
	if v := na.GetValueOrNil("edit_attrs"); v != nil {
		tag := func(comp h.HTMLComponent) *h.HTMLTagBuilder {
			return comp.(tag.TagGetter).GetHTMLTagBuilder()
		}
		for k, v := range vm.ToInterface(v).(map[string]any) {
			var handler func(ctx *presets.FieldContext, comp h.HTMLComponent) h.HTMLComponent
			if k[0] == '_' {
				k = ":" + k[1:]
			}
			k = strings.ReplaceAll(k, "_", "-")

			switch vt := v.(type) {
			case gad.Flag:
				handler = func(ctx *presets.FieldContext, comp h.HTMLComponent) h.HTMLComponent {
					tag(comp).SetAttr(k, !vt.IsFalsy())
					return comp
				}
			default:
				handler = func(ctx *presets.FieldContext, comp h.HTMLComponent) h.HTMLComponent {
					tag(comp).SetAttr(k, vt)
					return comp
				}
			}

			if handler != nil {
				t.EditComponentHandlers = append(t.EditComponentHandlers, handler)
			}
		}
	}
	if v, ok := na.GetValueOrNil("hint").(gad.Str); ok {
		t.ListHint = string(v)
		t.DetailHint = string(v)
		t.WriteHint = string(v)
	}
	if v, ok := na.GetValueOrNil("detail_hint").(gad.Str); ok {
		t.DetailHint = string(v)
	}
	if v, ok := na.GetValueOrNil("write_hint").(gad.Str); ok {
		t.WriteHint = string(v)
	}
	if v, ok := na.GetValueOrNil("list_hint").(gad.Str); ok {
		t.ListHint = string(v)
	}

	t.ReadOnly = !na.GetValue("ro").IsFalsy()
	t.EditorJS = !na.GetValue("editorjs").IsFalsy()
	t.TipTap = !na.GetValue("tiptap").IsFalsy()

	return
}

type Hint struct {
	value string
	fixed bool
}

func (h *Hint) Fixed() *Hint {
	h.fixed = true
	return h
}

var tstr = reflect.TypeOf("")

func FieldReadyHandle(field *presets.FieldBuilder) {
	var (
		tag      AdminTag
		tagValue = strings.TrimSpace(field.StructField().Tag.Get("admin"))
	)

	if tagValue == "-" {
		field.Enable(false)
		return
	}

	if !tag.Parse(&field.StructField().StructField, tagValue) {
		return
	}

	if field.GetData(AdminTagKey) != nil {
		return
	}

	field.SetData(AdminTagKey, &tag)

	if field.StructField().Type == tstr {
		field.ValueFormatters.Append(presets.FieldValueFormatterFunc(func(field *presets.FieldContext) (err error) {
			field.ValueOverride = strings.TrimSpace(field.ValueOverride.(string))
			return
		}))
	}

	if tag.Required != 0 {
		field.ContextSetup(func(ctx *presets.FieldContext) {
			if !ctx.Mode.Dot().HasAny(presets.DETAIL, presets.LIST) {
				if ctx.Label != "" {
					ctx.Label += " *"
				}
			}
		})
		field.Validators.Append(presets.FieldValidatorFunc(func(field *presets.FieldContext) (err web.ValidationErrors) {
			if field.ReadOnly || (field.EventContext != nil && ContextIsSkipFieldRequirementCheck(presets.GetSaveContext(field.EventContext), field.FormKey)) {
				return
			}
			if zeroer.IsZero(field.Value()) || len(strings.TrimSpace(field.StringValue())) == 0 {
				if st := field.Field.StructField(); st != nil {
					typ := st.Type
					if typ.Kind() == reflect.Pointer {
						typ = typ.Elem()
					}
					if typ.Kind() == reflect.Struct {
						if v, err2 := reflectutils.Get(field.Obj, field.Name+"ID"); err2 == nil && v != nil {
							if !gad.IsZero(v) {
								return
							}
						}
					}
				}
				err.FieldError(field.FormKey, GetMessages(field.EventContext.Context()).ErrFieldRequired.Error())
			}
			return
		}))
	}

	if tag.Mode != 0 {
		field.SetMode(tag.Mode)
	}

	if tag.ReadOnly {
		field.ContextSetup(func(ctx *presets.FieldContext) {
			ctx.ReadOnly = true
		})
	}

	if len(tag.DetailHint) > 0 {
		field.ContextSetup(func(ctx *presets.FieldContext) {
			if ctx.Mode.IsDetail() {
				ctx.Hint = tag.DetailHint
			}
		})
	}

	if len(tag.WriteHint) > 0 {
		field.ContextSetup(func(ctx *presets.FieldContext) {
			if ctx.Mode.IsWrite() {
				ctx.Hint = tag.WriteHint
			}
		})
	}

	if len(tag.ListHint) > 0 {
		field.ContextSetup(func(ctx *presets.FieldContext) {
			if ctx.Mode.IsList() {
				ctx.Hint = tag.ListHint
			}
		})
	}

	if tag.EditorJS {
		field.ComponentFunc(presets.EditorJSComponentFunc)
	}

	if len(tag.EditComponentHandlers) > 0 {
		field.ContextSetup(func(ctx *presets.FieldContext) {
			if !ctx.ReadOnly {
				ctx.ComponentHandler(tag.EditComponentHandlers...)
			}
		})
	}

	return
}

type ctxKey string

func ContextIsSkipFieldRequirementCheck(ctx context.Context, fieldKey string) (v bool) {
	if ctx == nil {
		return
	}
	v, _ = ctx.Value(ctxKey("skipFieldRequirementCheck:" + fieldKey)).(bool)
	return
}

func SkipFieldRequirementCheck(ctx web.ContextValuer, fieldKey ...string) {
	for _, k := range fieldKey {
		ctx.WithContextValue(ctxKey("skipFieldRequirementCheck:"+k), true)
	}
}
