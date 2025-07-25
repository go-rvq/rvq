package vuetifyx

import (
	"fmt"

	v "github.com/go-rvq/rvq/x/ui/vuetify"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

type VXTextFieldBuilder struct {
	label     string
	readOnly  bool
	dense     string
	vField    vField
	text      string
	class     string
	valueType string
	suffix    string
}

type vField struct {
	formKey string
	value   interface{}
}

func VXTextField() *VXTextFieldBuilder {
	return &VXTextFieldBuilder{}
}

func (b *VXTextFieldBuilder) Label(label string) *VXTextFieldBuilder {
	b.label = label
	return b
}

func (b *VXTextFieldBuilder) ReadOnly(readOnly bool) *VXTextFieldBuilder {
	b.readOnly = readOnly
	return b
}

func (b *VXTextFieldBuilder) Dense(dense string) *VXTextFieldBuilder {
	b.dense = dense
	return b
}

func (b *VXTextFieldBuilder) Text(value string) *VXTextFieldBuilder {
	b.text = value
	b.readOnly = true
	return b
}

func (b *VXTextFieldBuilder) VField(formKey string, value interface{}) *VXTextFieldBuilder {
	b.vField.formKey = formKey
	b.vField.value = value
	return b
}

func (b *VXTextFieldBuilder) Class(class string) *VXTextFieldBuilder {
	b.class = class
	return b
}

func (b *VXTextFieldBuilder) Type(valueType string) *VXTextFieldBuilder {
	b.valueType = valueType
	return b
}

func (b *VXTextFieldBuilder) Suffix(suffix string) *VXTextFieldBuilder {
	b.suffix = suffix
	return b
}

func (b *VXTextFieldBuilder) Write(ctx *h.Context) (err error) {
	var labelStyle string = "font-size:16px; font-weight:500;"
	var label h.HTMLComponent
	if b.label != "" {
		label = h.Div(h.Span(b.label).Style(labelStyle)).Class("mb-2")
	}
	if b.readOnly {
		div := h.Div().Class(b.class)
		if b.label != "" {
			div.AppendChildren(label)
		}
		if b.suffix != "" {
			b.text = fmt.Sprintf("%s %s", b.text, b.suffix)
		}
		div.AppendChildren(
			h.Div(h.Span(b.text)),
		)
		return div.Write(ctx)
	}

	var valueType string = "text"
	if b.valueType != "" {
		valueType = b.valueType
	}
	content := v.VTextField().HideDetails(true).Type(valueType).
		Variant(v.VariantOutlined).Density(v.DensityCompact).
		Suffix(b.suffix).
		Attr(web.VField(b.vField.formKey, b.vField.value)...)
	return h.Div(
		label,
		content,
	).Class(b.class).Write(ctx)
}
