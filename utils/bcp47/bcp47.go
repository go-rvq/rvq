// Package bcp47 offers pure helpers for BCP-47 language tags: validation, display
// names in the language's own tongue (autonyms), region flags, and the select
// item list. The offered tag list is the full CLDR set (golang.org/x/text) plus
// the common language-region variants. The admin field components that render
// these live in the bcp47field subpackage (which imports presets), keeping this
// package dependency-light so anything — presets included — can import it.
package bcp47

import (
	"sort"
	"sync"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

// capFirst upper-cases the first rune. CLDR autonyms follow each language's own
// casing — Portuguese and others lower-case their language name ("português"),
// English capitalizes it ("English") — so a label made of them reads
// inconsistently. Capitalizing just the first rune (not title-casing every word,
// which would wrongly upper-case "de"/"do") gives "Português", "Español de
// México", while leaving non-cased scripts ("中文") untouched.
func capFirst(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	if !unicode.IsLetter(r) || unicode.IsUpper(r) {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}

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
// nameOnly resolves the language name (autonym, then English), first rune
// capitalized. No lazy build — the namers must already be set (it is called from
// build itself and from NameOnly after once.Do).
func nameOnly(code string) string {
	t := language.Make(code)
	name := selfNamer.Name(t)
	if name == "" {
		name = enNamer.Name(t)
	}
	return capFirst(name)
}

func nameOf(code string) string {
	name := nameOnly(code)
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
	return nameOnly(code)
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

// FlagCodeLabel is CodeLabel prefixed with the region flag when the tag has one
// (e.g. "🇧🇷 pt-BR - português"); without a region it is just CodeLabel.
func FlagCodeLabel(code string) string {
	label := CodeLabel(code)
	if flag := Flag(code); flag != "" {
		return flag + " " + label
	}
	return label
}

// FlagLabel is Label prefixed with the region flag when the tag has one (e.g.
// "🇧🇷 Português (pt-BR)"); without a region it is just Label.
func FlagLabel(code string) string {
	label := Label(code)
	if flag := Flag(code); flag != "" {
		return flag + " " + label
	}
	return label
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

// Items returns the select items ([{value,title}]) for the autocomplete, each
// title as "🏳 name (code)" (the region flag, when the tag has one). A non-empty
// include code not already in the list is prepended (named from the CLDR namer),
// so the current value always has a matching item and the select shows its label
// instead of a blank.
func Items(include string) []map[string]string {
	once.Do(build)
	flagged := func(code, label string) string {
		if f := Flag(code); f != "" {
			return f + " " + label
		}
		return label
	}
	items := make([]map[string]string, 0, len(sorted)+1)
	if include != "" {
		if _, ok := byCode[include]; !ok {
			items = append(items, map[string]string{"value": include, "title": flagged(include, nameOf(include))})
		}
	}
	for _, e := range sorted {
		items = append(items, map[string]string{"value": e.Code, "title": flagged(e.Code, e.Label)})
	}
	return items
}

// Flag returns the flag emoji for a tag's region subtag (e.g. "pt-BR" → "🇧🇷",
// "en-US" → "🇺🇸"), built from the two regional-indicator symbols. It is empty
// when the tag names no two-letter region (a bare "pt", or a numeric UN M49
// region): a language is not a country, so there is nothing to show. The flag
// stands for the tag's region, not for where the language is spoken.
func Flag(code string) string {
	t := language.Make(code)
	region, conf := t.Region()
	if conf == language.No {
		return ""
	}
	s := region.String()
	if len(s) != 2 {
		return "" // numeric M49 region (e.g. "419"): no flag
	}
	r := []rune(s)
	for _, c := range r {
		if c < 'A' || c > 'Z' {
			return ""
		}
	}
	return string([]rune{0x1F1E6 + (r[0] - 'A'), 0x1F1E6 + (r[1] - 'A')})
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
