package gitedit

import (
	"testing"

	"github.com/go-rvq/rvq/x/i18n/msgpath"
)

// The texts of the IDE: their tags parse (an apostrophe in a hint breaks the
// admin's start), every one translated.
func TestIdeMessagesTags(t *testing.T) {
	for name, m := range map[string]*IdeMessages{"en": IdeMessages_en_US, "pt-BR": IdeMessages_pt_BR} {
		entries, _, err := msgpath.Walk(m)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		for _, e := range entries {
			if v, _ := msgpath.Get(m, e.Path); v == "" && e.Path != "ModuleDescription" {
				t.Errorf("%s: %s has no text", name, e.Path)
			}
		}
	}
}
