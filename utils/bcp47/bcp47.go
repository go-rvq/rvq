// Package bcp47 offers helpers for BCP-47 language tags: validation, display
// names in the language's own tongue (autonyms), and admin field components (a
// searchable select and read-only renderers). The offered tag list is the full
// CLDR set (golang.org/x/text) plus the common language-region variants.
package bcp47

import (
	"sort"
	"sync"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// commonRegionVariants are language-region tags people actually store as a
// locale (en-US, pt-BR, …). display.Supported.Tags() lists base languages and a
// few regional forms but omits these, so they are added to the offered list —
// otherwise a stored ID like "en-US" would have no matching item and the select
// would show it blank.
var commonRegionVariants = []string{
	"en-US", "en-GB", "en-CA", "en-AU",
	"pt-BR", "pt-PT",
	"es-ES", "es-MX", "es-AR",
	"fr-FR", "fr-CA",
	"de-DE", "de-AT",
	"it-IT", "nl-NL",
	"zh-CN", "zh-TW", "zh-HK",
	"ja-JP", "ko-KR", "ru-RU",
}

var (
	once      sync.Once
	selfNamer display.Namer
	enNamer   display.Namer
	sorted    []entry
	byCode    map[string]string
)

type entry struct {
	Code  string
	Label string
}

// nameOf returns the display label for any BCP-47 code — the language name in its
// OWN language (the autonym) plus the code, e.g. "português (pt-BR)", "中文
// (zh-CN)". The code disambiguates variants that share an autonym (pt / pt-BR are
// both "português"). Falls back to the English name, then to the bare code. Works
// for tags outside the offered list too, so a stored value still shows a name.
func nameOf(code string) string {
	t := language.Make(code)
	name := selfNamer.Name(t)
	if name == "" {
		name = enNamer.Name(t)
	}
	if name == "" {
		return code
	}
	return name + " (" + code + ")"
}

func build() {
	selfNamer = display.Self
	enNamer = display.English.Tags()

	// Base languages from CLDR, plus the common language-region variants people
	// store as a locale (en-US, pt-BR, …), which the base list omits.
	seen := map[string]bool{}
	var codes []string
	for _, t := range display.Supported.Tags() {
		c := t.String()
		if !seen[c] {
			seen[c] = true
			codes = append(codes, c)
		}
	}
	for _, c := range commonRegionVariants {
		if !seen[c] {
			seen[c] = true
			codes = append(codes, c)
		}
	}

	sorted = make([]entry, 0, len(codes))
	byCode = make(map[string]string, len(codes))
	for _, code := range codes {
		label := nameOf(code)
		sorted = append(sorted, entry{Code: code, Label: label})
		byCode[code] = label
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Label < sorted[j].Label })
}

// NameOnly returns just the language name (the autonym, then the English name),
// without the code — for a "CODE - name" rendering where the code is shown
// separately. Empty when the namer has none.
func NameOnly(code string) string {
	once.Do(build)
	t := language.Make(code)
	if n := selfNamer.Name(t); n != "" {
		return n
	}
	return enNamer.Name(t)
}

// CodeLabel returns a tag as "CODE - name" (e.g. "pt-BR - português"), the code
// first and then the language's own name — the label the admin's locale switcher
// shows. Just the code when no name can be derived.
func CodeLabel(code string) string {
	if code == "" {
		return ""
	}
	if name := NameOnly(code); name != "" {
		return code + " - " + name
	}
	return code
}

// Codes returns every offered tag code, ordered by display name.
func Codes() []string {
	once.Do(build)
	codes := make([]string, len(sorted))
	for i, e := range sorted {
		codes[i] = e.Code
	}
	return codes
}

// Label returns the "name (code)" display label for a tag. A tag outside the
// offered list (a stored en-US, a hand-typed one) is still named from the CLDR
// namer, so the select never shows a bare code where a name exists.
func Label(code string) string {
	if code == "" {
		return ""
	}
	once.Do(build)
	if l, ok := byCode[code]; ok {
		return l
	}
	return nameOf(code)
}

// Items returns the select items ([{value,title}]) for the autocomplete. A
// non-empty include code not already in the list is prepended (named from the
// CLDR namer), so the current value always has a matching item and the select
// shows its label instead of a blank.
func Items(include string) []map[string]string {
	once.Do(build)
	items := make([]map[string]string, 0, len(sorted)+1)
	if include != "" {
		if _, ok := byCode[include]; !ok {
			items = append(items, map[string]string{"value": include, "title": nameOf(include)})
		}
	}
	for _, e := range sorted {
		items = append(items, map[string]string{"value": e.Code, "title": e.Label})
	}
	return items
}

// Valid reports whether code is a non-empty, parseable BCP-47 language tag (the
// shape an HTML lang attribute needs).
func Valid(code string) bool {
	if code == "" {
		return false
	}
	tag, err := language.Parse(code)
	return err == nil && tag != language.Und
}

// AutocompleteComponentFunc renders a searchable BCP-47 select (a VAutocomplete
// over the full tag list) bound to the field, so a language is picked by typing.
func AutocompleteComponentFunc(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
	current := field.StringValue()

	ac := v.VAutocomplete().
		Label(field.InputLabel()).
		Items(Items(current)).
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
		body = h.Text(Label(code))
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
		if name := NameOnly(code); name != "" {
			text = code + " - " + name
		}
		body = h.Text(text)
	}
	if field.Mode.Dot().IsList() {
		return h.Td(body)
	}
	return h.Div(body)
}
