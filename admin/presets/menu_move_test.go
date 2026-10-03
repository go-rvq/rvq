package presets

import (
	"strings"
	"testing"
)

type menuGamma struct{ ID int }
type menuDelta struct{ ID int }

// keysIn is the group's items by key, in order — what the menu will render.
func keysIn(g *MenuGroupBuilder) []string {
	r := make([]string, len(g.Items()))
	for i, it := range g.Items() {
		r[i] = it.Key()
	}
	return r
}

func joinKeys(g *MenuGroupBuilder) string { return strings.Join(keysIn(g), " ") }

// An item is in exactly one place. Every move takes it out of where it was,
// whatever that was — the root, another group, or the same group again.
func TestMoveBetweenEveryPairOfPlaces(t *testing.T) {
	b := menuBuilder()
	b.Model(&menuGamma{})
	root := b.MenuTree()
	a := b.MenuGroup("a")
	c := b.MenuGroup("c")
	inner := b.MenuGroup("inner")
	a.Add(inner)

	key := ModelItem("menu_gammas").Key()
	// Registration leaves it at the root.
	if got := b.MenuItems()[key].Parent(); got != root {
		t.Fatalf("depois do registro o pai é %v, queria a raiz", got)
	}

	for _, step := range []struct {
		name string
		dst  *MenuGroupBuilder
	}{
		{"raiz → grupo", a},
		{"grupo → outro grupo", c},
		{"grupo → subgrupo", inner},
		{"subgrupo → raiz", root},
		{"raiz → subgrupo", inner},
	} {
		if err := b.MoveMenuItem(ModelItem("menu_gammas"), step.dst); err != nil {
			t.Fatalf("%s: %v", step.name, err)
		}
		if got := b.MenuItems()[key].Parent(); got != step.dst {
			t.Errorf("%s: pai = %v", step.name, got)
		}
		// Exactly one copy, across the whole tree.
		if n := countKey(root, key); n != 1 {
			t.Errorf("%s: o item aparece %d vezes na árvore", step.name, n)
		}
	}
}

// countKey counts the item across the tree, to catch a move that copies.
func countKey(g *MenuGroupBuilder, key string) (n int) {
	for _, it := range g.Items() {
		if it.Key() == key {
			n++
		}
		if child := it.Group(); child != nil {
			n += countKey(child, key)
		}
	}
	return
}

// A move with no destination puts the item back at the root.
func TestMoveToNilIsTheRoot(t *testing.T) {
	b := menuBuilder()
	b.Model(&menuGamma{})
	a := b.MenuGroup("a").Add(ModelItem("menu_gammas"))

	if err := b.MoveMenuItem(ModelItem("menu_gammas"), nil); err != nil {
		t.Fatal(err)
	}
	if len(a.Items()) != 0 {
		t.Error("o item continua no grupo")
	}
	if b.MenuItems()[ModelItem("menu_gammas").Key()].Parent() != b.MenuTree() {
		t.Error("o item não voltou para a raiz")
	}
}

// Order is what the caller asked for, and an index puts the item at a place
// instead of at the end.
func TestMoveKeepsAndChoosesOrder(t *testing.T) {
	b := menuBuilder()
	g := b.MenuGroup("g").Add(ModelItem("one"), ModelItem("two"), ModelItem("three"))

	if got, want := joinKeys(g), "m:one m:two m:three"; got != want {
		t.Errorf("ordem = %q, want %q", got, want)
	}

	// To the front.
	if err := b.MoveMenuItem(ModelItem("three"), g, 0); err != nil {
		t.Fatal(err)
	}
	if got, want := joinKeys(g), "m:three m:one m:two"; got != want {
		t.Errorf("depois do índice 0 = %q, want %q", got, want)
	}

	// To the middle.
	if err := b.MoveMenuItem(ModelItem("three"), g, 1); err != nil {
		t.Fatal(err)
	}
	if got, want := joinKeys(g), "m:one m:three m:two"; got != want {
		t.Errorf("depois do índice 1 = %q, want %q", got, want)
	}

	// An index past the end, or negative, means the end.
	for _, idx := range []int{99, -1} {
		if err := b.MoveMenuItem(ModelItem("one"), g, idx); err != nil {
			t.Fatal(err)
		}
		if got, want := joinKeys(g), "m:three m:two m:one"; got != want {
			t.Errorf("índice %d = %q, want %q", idx, got, want)
		}
	}
}

// Moving a group takes its contents with it, and every descendant's path
// follows — which is what a model's URI is built from.
func TestMovingAGroupCarriesItsSubtree(t *testing.T) {
	b := menuBuilder()
	mb := b.Model(&menuDelta{})

	a := b.MenuGroup("a")
	mid := b.MenuGroup("mid")
	leaf := b.MenuGroup("leaf")

	a.Add(mid)
	mid.Add(leaf)
	leaf.Add(ModelItem(mb.id))

	if got, want := mb.URI(), "a/mid/leaf/"+mb.uriName; got != want {
		t.Fatalf("URI() = %q, want %q", got, want)
	}

	// Move the middle group to the root: the leaf and the model come along.
	if err := b.MoveMenuItem(GroupItem("mid"), nil); err != nil {
		t.Fatal(err)
	}
	if got, want := leaf.Path(), "mid/leaf"; got != want {
		t.Errorf("Path() do leaf = %q, want %q", got, want)
	}
	if got, want := mb.URI(), "mid/leaf/"+mb.uriName; got != want {
		t.Errorf("URI() = %q, want %q", got, want)
	}
	if got := mb.MenuGroup(); got != leaf {
		t.Errorf("o modelo saiu do leaf: %v", got)
	}

	// And back under a deeper place: the whole chain updates again.
	other := b.MenuGroup("other")
	if err := b.MoveMenuItem(GroupItem("mid"), other); err != nil {
		t.Fatal(err)
	}
	if got, want := mb.URI(), "other/mid/leaf/"+mb.uriName; got != want {
		t.Errorf("URI() = %q, want %q", got, want)
	}
}

// Moving a key that nobody registered is allowed: it reserves the place. The
// value arrives later, wherever the last move left the reservation.
func TestMovePendingThenRegister(t *testing.T) {
	b := menuBuilder()
	first := b.MenuGroup("first")
	second := b.MenuGroup("second")

	if err := b.MoveMenuItem(ModelItem("menu_gammas"), first); err != nil {
		t.Fatal(err)
	}
	if err := b.MoveMenuItem(ModelItem("menu_gammas"), second, 0); err != nil {
		t.Fatal(err)
	}

	mb := b.Model(&menuGamma{})
	if got := mb.MenuGroup(); got != second {
		t.Errorf("grupo = %v, want second", got)
	}
	if got, want := joinKeys(second), "m:menu_gammas"; got != want {
		t.Errorf("second = %q, want %q", got, want)
	}
	if len(first.Items()) != 0 {
		t.Error("a reserva ficou também no primeiro grupo")
	}
}

// A page moves like anything else, and its URL follows the group it lands in.
func TestMovePageUpdatesItsURL(t *testing.T) {
	b := menuBuilder()
	page := HttpPage("/report")
	b.PagesRegistrator().AddHttpPage(page)

	a := b.MenuGroup("a")
	deep := b.MenuGroup("deep")
	a.Add(deep)

	for _, c := range []struct {
		dst  *MenuGroupBuilder
		want string
	}{
		{nil, "/admin/report"},
		{a, "/admin/a/report"},
		{deep, "/admin/a/deep/report"},
	} {
		if err := b.MoveMenuItem(PageItem("/report"), c.dst); err != nil {
			t.Fatal(err)
		}
		page.Build("admin")
		if got := page.FullPath(); got != c.want {
			t.Errorf("FullPath() = %q, want %q", got, c.want)
		}
	}
}

// A refused move changes nothing: the item stays where it was.
func TestRejectedMoveKeepsThePlace(t *testing.T) {
	b := menuBuilder()
	a := b.MenuGroup("a")
	inner := b.MenuGroup("inner")
	a.Add(inner)

	before := joinKeys(a)
	if err := b.MoveMenuItem(GroupItem("a"), inner); err == nil {
		t.Fatal("o ciclo passou")
	}
	if got := joinKeys(a); got != before {
		t.Errorf("a árvore mudou: %q -> %q", before, got)
	}
	if inner.Parent() != a {
		t.Error("o subgrupo saiu do lugar")
	}
	if a.Parent() != nil {
		t.Error("o grupo de cima ganhou um pai")
	}
}
