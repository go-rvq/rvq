// Package bcp47field holds the admin field components for BCP-47 language tags —
// a searchable select and read-only renderers — built on the pure helpers in the
// bcp47 package. It is separate so bcp47 stays dependency-light (importable by
// presets), while these components depend on presets/web/vuetify.
package bcp47field

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/utils/bcp47"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// AutocompleteComponentFunc renders a searchable BCP-47 select (a VAutocomplete
// over the full tag list) bound to the field, so a language is picked by typing.
func AutocompleteComponentFunc(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	current := field.StringValue()

	ac := v.VAutocomplete().
		Label(field.InputLabel()).
		Items(bcp47.Items(current)).
		ItemValue("value").
		ItemTitle("title").
		Variant("underlined").
		Clearable(true).
		Attr(web.VField(field.FormKey, current)...)

	if len(field.Errors) > 0 {
		ac.Attr(":error-messages", h.JSONString(field.Errors))
	}
	if hint := field.CheckHint().Hint; hint != "" {
		ac.Hint(hint).PersistentHint(true)
	}
	if field.Disabled {
		ac.Attr("disabled", true)
	}
	return ac
}

// ReadonlyComponentFunc renders a BCP-47 field read-only: the tag's display label
// ("name (code)"), or empty when unset. Mode-aware — a table cell in a listing,
// plain text in a detail view (which is not a table).
func ReadonlyComponentFunc(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	var (
		body h.HTMLComponent
		code = field.StringValue()
	)
	if code != "" {
		body = h.Text(bcp47.Label(code))
	}
	if field.Mode.Dot().IsList() {
		return h.Td(body)
	}
	return h.Div(body)
}

// FlagLabelReadonlyComponentFunc renders a BCP-47 field read-only as "🏳 name
// (code)" (the region flag, when the tag has one). Mode-aware — a table cell in a
// listing, plain text in a detail view.
func FlagLabelReadonlyComponentFunc(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	var (
		body h.HTMLComponent
		code = field.StringValue()
	)
	if code != "" {
		body = h.Text(bcp47.FlagLabel(code))
	}
	if field.Mode.Dot().IsList() {
		return h.Td(body)
	}
	return h.Div(body)
}

// CodeLabelReadonlyComponentFunc renders a BCP-47 field read-only as "CODE - name"
// (e.g. "pt-BR - português"). Mode-aware (a table cell in a listing, plain text in
// a detail view).
func CodeLabelReadonlyComponentFunc(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	var (
		body h.HTMLComponent
		code = field.StringValue()
	)
	if code != "" {
		text := code
		if name := bcp47.NameOnly(code); name != "" {
			text = code + " - " + name
		}
		body = h.Text(text)
	}
	if field.Mode.Dot().IsList() {
		return h.Td(body)
	}
	return h.Div(body)
}
