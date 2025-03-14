package presets

import (
	"fmt"
	"time"

	"github.com/qor5/web/v3"
	"github.com/qor5/web/v3/zeroer"
	"github.com/qor5/x/v3/i18n"
	"github.com/qor5/x/v3/ui/editorjs"
	. "github.com/qor5/x/v3/ui/vuetify"
	vx "github.com/qor5/x/v3/ui/vuetifyx"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
)

func TDReadonlyBoolComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	var c h.HTMLComponent
	if zeroer.IsZero(field.Value()) {
		c = h.Text("-")
	} else {
		c = VIcon("mdi-check")
	}
	return h.Td(c)
}

func TDStringComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return h.Td(h.Text(field.StringValue()))
}

func CheckboxComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return VCheckbox().
		Attr(web.VField(field.FormKey, field.Value().(bool))...).
		Label(field.Label).
		ErrorMessages(field.Errors...).
		Disabled(field.ReadOnly)
}

func CheckboxReadonlyComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	if zeroer.IsZero(field.Value()) {
		return nil
	}
	return FormFieldComponentWrapper(VChip(h.Text(field.Label)).PrependIcon("mdi-check"))
}

func SwitchComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return VSwitch().
		Attr(web.VField(field.FormKey, field.Value().(bool))...).
		Label(field.Label).
		ErrorMessages(field.Errors...).
		Disabled(field.ReadOnly)
}

func NumberComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return VTextField().
		Type("number").
		Variant(FieldVariantUnderlined).
		Attr(web.VField(field.FormKey, field.StringValue())...).
		Label(field.Label).
		Hint(field.Hint()).
		ErrorMessages(field.Errors...).
		Disabled(field.ReadOnly)
}

func TimeComponentFunc(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
	msgr := i18n.MustGetModuleMessages(ctx.Context(), CoreI18nModuleKey, Messages_en_US).(*Messages)
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
	return vx.VXDateTimePicker().
		Label(field.Label).
		Attr(web.VField(field.FormKey, val)...).
		Value(val).
		TimePickerProps(vx.TimePickerProps{
			Format:     "24hr",
			Scrollable: true,
		}).
		DialogWidth(640).
		ClearText(msgr.Clear).
		OkText(msgr.OK)
}

func TimeReadonlyComponentFunc(field *FieldContext, ctx *web.EventContext) h.HTMLComponent {
	var t *time.Time
	if v := field.Value(); v != nil {
		switch vt := v.(type) {
		case time.Time:
			if !vt.IsZero() {
				t = &vt
			}
		case *time.Time:
			if !vt.IsZero() {
				t = vt
			}
		default:
			panic(fmt.Sprintf("unknown time type: %T\n", v))
		}
	}

	if t == nil {
		return nil
	}

	msgr := i18n.MustGetModuleMessages(ctx.Context(), CoreI18nModuleKey, Messages_en_US).(*Messages)
	val := t.Format(msgr.TimeFormats.DateTime)
	return vx.VXReadonlyField().
		Label(field.Label).
		Value(val)
}

func TimeComponentFuncSetter(obj interface{}, field *FieldContext, ctx *web.EventContext) (err error) {
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

func TextFieldComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return VTextField().
		Type("text").
		Variant(FieldVariantUnderlined).
		Attr(web.VField(field.FormKey, field.StringValue())...).
		Label(field.Label).
		Hint(field.Hint()).
		ErrorMessages(field.Errors...).
		Disabled(field.ReadOnly)
}

func LongTextFieldComponentFunc(field *FieldContext, _ *web.EventContext) *VTextareaBuilder {
	return VTextarea().
		Type("text").
		Variant(FieldVariantUnderlined).
		Attr(web.VField(field.FormKey, field.StringValue())...).
		Label(field.Label).
		Hint(field.Hint()).
		ErrorMessages(field.Errors...).
		Disabled(field.ReadOnly)
}

func RuneFieldComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	if r, ok := field.Value().(rune); ok {
		if r == 0 {
			field.ValueOverride = ""
		} else {
			field.ValueOverride = string(r)
		}
	}

	if field.ReadOnly {
		return ReadonlyTextComponentFunc(field, nil)
	}
	return TextFieldComponentFunc(field, nil).(*VTextFieldBuilder).MaxLength(1)
}

func ReadonlyTextComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return vx.VXReadonlyField().
		Label(field.Label).
		Value(field.StringValue())
}

func FileFieldComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	return VFileInput().
		Variant(FieldVariantUnderlined).
		Attr(web.VField(field.FormKey, "")...).
		Label(field.Label).
		Hint(field.Hint()).
		ErrorMessages(field.Errors...).
		Disabled(field.ReadOnly)
}

func EditorJSComponentFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	if mode := field.Mode.Dot(); mode.HasAny(EDIT, NEW) {
		return EditorJSComponentWriteFunc(field, nil)
	} else if mode.Has(LIST) {
		s, _ := field.Value().(string)
		var (
			comp h.HTMLComponent
			err  error
		)
		if s, err = editorjs.Htmlify(s); err != nil {
			comp = h.Div(h.RawHTML(err.Error())).Class("text-error")
		} else {
			comp = h.RawHTML(s)
		}
		return h.Td(comp)
	}
	return EditorJSComponentReadFunc(field, nil)
}

func EditorJSComponentReadFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	v := field.Value()
	if v == nil {
		return nil
	}

	s, _ := v.(string)

	var (
		comp h.HTMLComponent
		err  error
	)

	if len(s) > 0 && s[0] == '<' {
		comp = h.RawHTML(s)
	} else if s, err = editorjs.Htmlify(s); err != nil {
		comp = h.Div(h.RawHTML(err.Error())).Class("text-error")
	} else {
		comp = h.RawHTML(s)
	}

	return h.HTMLComponents{
		h.Div(
			h.Label(field.Label),
		),
		VCard(
			VCardText(comp).Class("editorjs-body"),
		).MaxHeight(700).Class("overflow-auto"),
	}
}

func EditorJSComponentWriteFunc(field *FieldContext, _ *web.EventContext) h.HTMLComponent {
	s, _ := field.Value().(string)
	_, err := editorjs.Parse([]byte(s))
	return vx.EditorJS().Label(field.Label).FormField(field.FormKey, field.Value().(string)).Errors(err)
}
