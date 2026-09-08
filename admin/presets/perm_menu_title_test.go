package presets

import (
	"context"
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

type permTitleThing struct {
	ID int
}

// The permission tree carries each group's title. It used to assign the field
// to itself, so every group arrived titleless.
func TestBuildPermissionsCarriesTheGroupTitle(t *testing.T) {
	b := New(i18n.New().SupportLanguages(language.English))

	b.Model(&permTitleThing{}).SetMenuGroupName("settings")
	b.MenuGroup("settings").Title("Configurações")

	root := b.BuildPermissions()

	var found *PermMenu
	for _, child := range root.Children {
		if child.Name != "" && len(child.Resources) > 0 {
			found = child
		}
	}
	if found == nil {
		t.Fatal("the group did not reach the permission tree")
	}
	if found.Title == nil {
		t.Fatal("the group arrived with no title")
	}
	if got := found.Title(context.Background()); got != "Configurações" {
		t.Errorf("title = %q, want %q", got, "Configurações")
	}
}

// A group that was never given a title still answers, humanizing its name.
func TestBuildPermissionsTitlesAnUnnamedGroup(t *testing.T) {
	b := New(i18n.New().SupportLanguages(language.English))

	b.Model(&permTitleThing{}).SetMenuGroupName("site_settings")
	b.MenuGroup("site_settings")

	root := b.BuildPermissions()
	for _, child := range root.Children {
		if len(child.Resources) == 0 {
			continue
		}
		if child.Title == nil {
			t.Fatal("no title func for a group without a title")
		}
		if got := child.Title(context.Background()); got == "" {
			t.Error("the humanized fallback produced nothing")
		}
	}
}
