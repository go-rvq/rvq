package presets

import (
	"strings"
	"testing"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

type menuAlpha struct{ ID int }
type menuBeta struct{ ID int }

func menuBuilder() *Builder {
	return New(i18n.New().SupportLanguages(language.English))
}

// keyOf is what the tree files an item under: its type and its name. The type
// is half the key, so the three kinds never collide.
func TestMenuKeysAreIndependentPerType(t *testing.T) {
	b := menuBuilder()

	model := b.MenuItemOf(ModelItem("x"))
	page := b.MenuItemOf(PageItem("/x"))
	group := b.MenuItemOf(GroupItem("x"))

	keys := map[string]bool{model.Key(): true, page.Key(): true, group.Key(): true}
	if len(keys) != 3 {
		t.Errorf("as três chaves colidiram: %v", keys)
	}
	if model == page || model == group || page == group {
		t.Error("chaves diferentes devolveram o mesmo item")
	}
}

// A key is unique: registering a value over a value is a mistake, not a silent
// replacement.
func TestRegisterMenuItemRejectsASecondValue(t *testing.T) {
	b := menuBuilder()

	if _, err := b.RegisterMenuItem(MenuItemPage, "/report", &HttpPageBuilder{path: "/report"}); err != nil {
		t.Fatalf("primeiro registro: %v", err)
	}

	_, err := b.RegisterMenuItem(MenuItemPage, "/report", &HttpPageBuilder{path: "/report"})
	if err == nil {
		t.Fatal("o segundo registro passou")
	}
	if !strings.Contains(err.Error(), "already registered") {
		t.Errorf("erro = %v, queria que dissesse que já estava registrado", err)
	}

	// Registering the same value again is the same registration, not a second.
	page := &HttpPageBuilder{path: "/twice"}
	if _, err := b.RegisterMenuItem(MenuItemPage, "/twice", page); err != nil {
		t.Fatal(err)
	}
	if _, err := b.RegisterMenuItem(MenuItemPage, "/twice", page); err != nil {
		t.Errorf("registrar o mesmo valor deu erro: %v", err)
	}
}

// A key referenced before it exists holds its place and waits; registering
// fills it where it stands, in either order.
func TestReferenceBeforeRegistrationInEitherOrder(t *testing.T) {
	t.Run("referência antes", func(t *testing.T) {
		b := menuBuilder()
		site := b.MenuGroup("site").Add(ModelItem("menu_alphas"))

		it := b.MenuItems()[ModelItem("menu_alphas").Key()]
		if it == nil || it.Registered() {
			t.Fatal("a referência não criou a espera")
		}
		if it.Parent() != site {
			t.Error("a espera não ficou no grupo que a nomeou")
		}

		mb := b.Model(&menuAlpha{})
		if !it.Registered() {
			t.Error("o registro não preencheu a espera")
		}
		if it.Value != mb {
			t.Error("a espera foi preenchida com outro valor")
		}
		if it.Parent() != site {
			t.Error("o registro tirou o item do grupo")
		}
	})

	t.Run("registro antes", func(t *testing.T) {
		b := menuBuilder()
		mb := b.Model(&menuAlpha{})
		site := b.MenuGroup("site").Add(ModelItem("menu_alphas"))

		if got := mb.MenuGroup(); got != site {
			t.Errorf("MenuGroup() = %v, want site", got)
		}
	})
}

// Two early references to the same key: the last one wins, and the item is in
// exactly one place.
func TestLastReferenceWins(t *testing.T) {
	b := menuBuilder()
	first := b.MenuGroup("first").Add(ModelItem("menu_alphas"))
	second := b.MenuGroup("second").Add(ModelItem("menu_alphas"))

	if len(first.Items()) != 0 {
		t.Errorf("o primeiro grupo ficou com %d itens, queria 0", len(first.Items()))
	}
	if len(second.Items()) != 1 {
		t.Fatalf("o segundo grupo ficou com %d itens, queria 1", len(second.Items()))
	}

	b.Model(&menuAlpha{})
	if got := b.GetModelByID("menu_alphas").MenuGroup(); got != second {
		t.Errorf("grupo = %v, want second", got)
	}
}

// Moving is by key, and never duplicates the item.
func TestMoveDoesNotDuplicate(t *testing.T) {
	b := menuBuilder()
	b.Model(&menuAlpha{})
	a := b.MenuGroup("a")
	c := b.MenuGroup("c")

	for _, dst := range []*MenuGroupBuilder{a, c, a} {
		if err := b.MoveMenuItem(ModelItem("menu_alphas"), dst); err != nil {
			t.Fatal(err)
		}
	}

	if len(a.Items()) != 1 || len(c.Items()) != 0 {
		t.Errorf("a=%d c=%d, queria a=1 c=0", len(a.Items()), len(c.Items()))
	}
	// And gone from the root, where registration first put it.
	for _, it := range b.MenuTree().Items() {
		if it.Key() == ModelItem("menu_alphas").Key() {
			t.Error("o item continua na raiz")
		}
	}
}

// A group cannot go inside itself or inside its own descendant.
func TestMoveRejectsCycles(t *testing.T) {
	b := menuBuilder()
	a := b.MenuGroup("a")
	bb := b.MenuGroup("b")
	a.Add(bb)

	if err := b.MoveMenuItem(GroupItem("a"), a); err == nil {
		t.Error("um grupo entrou em si mesmo")
	}
	if err := b.MoveMenuItem(GroupItem("a"), bb); err == nil {
		t.Error("um grupo entrou no próprio descendente")
	}
}

// Three groups deep, with a model at the fourth level: the model's URI, its
// permission path and its breadcrumb all carry the whole chain, not only the
// innermost group.
func TestThreeLevelsDeepModelKeepsTheWholePath(t *testing.T) {
	b := menuBuilder()
	mb := b.Model(&menuBeta{})

	a := b.MenuGroup("a")
	c := b.MenuGroup("c")
	d := b.MenuGroup("d")
	a.Add(c)
	c.Add(d)
	d.Add(ModelItem(mb.id))

	if got, want := d.Path(), "a/c/d"; got != want {
		t.Errorf("Path() = %q, want %q", got, want)
	}
	if got, want := mb.MenuGroupName(), "a/c/d"; got != want {
		t.Errorf("MenuGroupName() = %q, want %q", got, want)
	}
	if got, want := mb.URI(), "a/c/d/"+mb.uriName; got != want {
		t.Errorf("URI() = %q, want %q", got, want)
	}
}

// A page nested in a group is served under the group's whole path.
func TestNestedPageFullPath(t *testing.T) {
	b := menuBuilder()
	page := HttpPage("/report")
	b.PagesRegistrator().AddHttpPage(page)

	a := b.MenuGroup("a")
	inner := b.MenuGroup("inner")
	a.Add(inner)
	inner.Add(PageItem("/report"))

	page.Build("admin")
	if got, want := page.FullPath(), "/admin/a/inner/report"; got != want {
		t.Errorf("FullPath() = %q, want %q", got, want)
	}
}
