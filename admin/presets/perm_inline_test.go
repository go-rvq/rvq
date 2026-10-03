package presets

import (
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
)

type inlineConfig struct {
	ID    uint
	Title string
}

type inlineGmail struct {
	User string
}

type inlinePost struct {
	ID     uint
	Title  string
	Config inlineConfig
	Gmail  inlineGmail
}

// A field holding a model edited in place is "&Name", its fields under it; a
// nested struct's fields are under it as fields ("#Gmail:#User"). The
// resource the runtime asks and the one of the dump are the same.
func TestInlinePerm(t *testing.T) {
	b := New(i18n.New()).URIPrefix("/admin")
	b.Permission(perm.New())
	config := NewModelBuilder(b, &inlineConfig{}, ModelConfig().SetId("config"))
	posts := b.Model(&inlinePost{}, ModelWithID("posts"))
	ed := posts.Editing("Title", "Config", "Gmail")
	ed.Field("Config").AutoNested(config, &config.Editing("Title").FieldsBuilder)
	// a nested struct with no model of its own
	gmail := NewFieldsBuilder(b).Model(&inlineGmail{})
	gmail.Field("User")
	ed.Field("Gmail").Nested(NestedStruct(nil, gmail))

	p := posts.Permissioner()
	for field, want := range map[string]string{
		"Title":        "#Title",
		"Config":       "&Config",
		"Config.Title": "&Config:#Title",
		"Gmail.User":   "#Gmail:#User",
	} {
		if got := joinParts(p.FieldPermParts(field)); got != want {
			t.Errorf("%s: %q, want %q", field, got, want)
		}
	}
	dump := map[string]bool{}
	for _, e := range b.BuildPermissions().Tree().Zip() {
		dump[e.Unique] = true
	}
	for _, res := range []string{":posts:<*>:&Config:", ":posts:<*>:&Config:#Title:", ":posts:<*>:#Gmail:#User:"} {
		if !dump[PermModule+res] {
			t.Errorf("%s not in the dump", res)
		}
	}
}

func joinParts(parts []string) (s string) {
	for i, p := range parts {
		if i > 0 {
			s += ":"
		}
		s += p
	}
	return
}
