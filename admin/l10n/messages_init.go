package l10n

import "fmt"

// MessageInit is one localized-message definition produced by an initializer,
// shaped like an entry of the site's messages config file. The application layer
// (e.g. hermon-cms) maps it onto its LocaleMessage model and upserts it.
//
// Origin says where a message comes from, when it is not an ordinary text of
// the site (from the template scan or the config file, editable): the
// application names its origins — hermon-cms, for instance, has "system" for a
// label the Go code owns and translates (read-only in the admin) and
// "locale_config" for the words of a form a config key declares. "" is an
// ordinary text.
type MessageInit struct {
	Key    string
	Type   string // text | html | yaml ("" means text)
	Value  string // the value for the locale being initialized
	Help   string
	Hint   string
	Origin string
	// Lang is the language Value is written in, when it is not the locale's
	// own — a fallback ("en") the code gives a language it has no text for,
	// or "auto" when it is not known: the value is then a source to
	// translate into the locale, not the locale's text. Empty: Value is the
	// locale's own.
	Lang string
}

// MessageInitFunc produces the message definitions for one locale (BCP-47 code).
type MessageInitFunc func(locale string) ([]MessageInit, error)

type messageInitializer struct {
	name string
	fn   MessageInitFunc
}

// RegisterMessageInitializer adds a message initializer. Initializers run in the
// order registered; a later one overrides an earlier one for the same key (so the
// config file can override the template scan, and application/system definitions
// can override both). The template scan and the config file are themselves
// registered as initializers by the application.
func (b *Builder) RegisterMessageInitializer(name string, fn MessageInitFunc) *Builder {
	b.messageInits = append(b.messageInits, messageInitializer{name: name, fn: fn})
	return b
}

// HasMessageInitializers reports whether any initializer is registered.
func (b *Builder) HasMessageInitializers() bool {
	return len(b.messageInits) > 0
}

// CollectMessages runs every registered initializer for a locale and returns the
// merged definitions in first-seen key order, with later initializers overriding
// earlier ones field by field (a later empty field — Origin included — does not
// erase an earlier value).
func (b *Builder) CollectMessages(locale string) ([]MessageInit, error) {
	byKey := map[string]*MessageInit{}
	var order []string
	for _, it := range b.messageInits {
		defs, err := it.fn(locale)
		if err != nil {
			return nil, fmt.Errorf("l10n message initializer %q: %w", it.name, err)
		}
		for _, d := range defs {
			if d.Key == "" {
				continue
			}
			cur := byKey[d.Key]
			if cur == nil {
				cp := d
				byKey[d.Key] = &cp
				order = append(order, d.Key)
				continue
			}
			if d.Type != "" {
				cur.Type = d.Type
			}
			if d.Value != "" {
				// a value and the language it is written in go together
				cur.Value, cur.Lang = d.Value, d.Lang
			}
			if d.Help != "" {
				cur.Help = d.Help
			}
			if d.Hint != "" {
				cur.Hint = d.Hint
			}
			if d.Origin != "" {
				cur.Origin = d.Origin
			}
		}
	}
	out := make([]MessageInit, 0, len(order))
	for _, k := range order {
		out = append(out, *byKey[k])
	}
	return out, nil
}
