package presets

import (
	"context"
	"fmt"
	"strings"
)

// MenuItemType says what a menu key names. It is half of the key, so a model, a
// page and a group may share a name without colliding.
type MenuItemType string

const (
	MenuItemGroup MenuItemType = "g"
	MenuItemModel MenuItemType = "m"
	MenuItemPage  MenuItemType = "p"
)

// menuKey is the identity of a menu item: its type and its name.
func menuKey(typ MenuItemType, name string) string { return string(typ) + ":" + name }

// MenuRef names an item without requiring it to exist yet. It is what the
// group-building API takes, so a group can name a model, a page or another
// group before any of them is registered.
type MenuRef struct {
	Type MenuItemType
	Name string
}

// ModelItem, PageItem and GroupItem are the typed references. Use them instead
// of a bare string, which cannot say which of the three it means.
func ModelItem(id string) MenuRef   { return MenuRef{Type: MenuItemModel, Name: id} }
func PageItem(path string) MenuRef  { return MenuRef{Type: MenuItemPage, Name: path} }
func GroupItem(name string) MenuRef { return MenuRef{Type: MenuItemGroup, Name: name} }

func (r MenuRef) Key() string { return menuKey(r.Type, r.Name) }

// MenuItem is one entry of the menu tree.
//
// Value is what the entry renders — a *ModelBuilder, a *HttpPageBuilder or a
// *MenuGroupBuilder — and is nil while the key has only been referenced: a
// group that names an item before it exists puts the item in place anyway,
// valueless, and it waits there. Registering the key later fills the value
// without moving anything, so the order in which an application configures its
// models, pages and groups does not matter.
type MenuItem struct {
	Type   MenuItemType
	Name   string
	Value  any
	parent *MenuGroupBuilder
}

// Key is the item's identity.
func (i *MenuItem) Key() string { return menuKey(i.Type, i.Name) }

// Ref is the item as a reference.
func (i *MenuItem) Ref() MenuRef { return MenuRef{Type: i.Type, Name: i.Name} }

// Registered reports whether the item's value has arrived. An unregistered item
// holds a place in the tree and renders nothing.
func (i *MenuItem) Registered() bool { return i.Value != nil }

// Parent is the group the item sits in, or nil at the root.
func (i *MenuItem) Parent() *MenuGroupBuilder { return i.parent }

// Group returns the item's value as a group, or nil when it is not one.
func (i *MenuItem) Group() *MenuGroupBuilder {
	g, _ := i.Value.(*MenuGroupBuilder)
	return g
}

type MenuGroupBuilder struct {
	item  *MenuItem
	title func(ctx context.Context) string
	// description is what is in the group, in the language of the request
	description func(ctx context.Context) string
	name        string
	icon        string
	items       []*MenuItem
	// uriName is the group's segment in its models' URIs; nil is its name.
	uriName *string

	// b is the builder the group belongs to; it owns the key registry, so the
	// group needs it to resolve the names its sub-items refer to.
	b *Builder
}

func (b *MenuGroupBuilder) Name() string { return b.name }

// Item is the group's own entry in the tree — nil for the root sentinel.
func (b *MenuGroupBuilder) Item() *MenuItem { return b.item }

// Items are the group's entries, in order.
func (b *MenuGroupBuilder) Items() []*MenuItem { return b.items }

// Parent is the group this one sits in, or nil for a top-level group. The root
// sentinel is not a group anyone is in, so a group directly under it has no
// parent — the same answer Ancestors and a model's MenuGroup give.
func (b *MenuGroupBuilder) Parent() *MenuGroupBuilder {
	if b.item == nil || b.item.parent == nil || b.item.parent.item == nil {
		return nil
	}
	return b.item.parent
}

func (b *MenuGroupBuilder) TitleFunc(f func(ctx context.Context) string) *MenuGroupBuilder {
	b.title = f
	return b
}

func (b *MenuGroupBuilder) Title(s string) *MenuGroupBuilder {
	b.title = func(context.Context) string {
		return s
	}
	return b
}

func (b *MenuGroupBuilder) TTitle(ctx context.Context) string {
	if b == nil {
		return ""
	}
	if b.title != nil {
		return b.title(ctx)
	}
	return HumanizeString(b.name)
}

func (b *MenuGroupBuilder) DescriptionFunc(f func(ctx context.Context) string) *MenuGroupBuilder {
	b.description = f
	return b
}

func (b *MenuGroupBuilder) Description(s string) *MenuGroupBuilder {
	b.description = func(context.Context) string {
		return s
	}
	return b
}

// TDescription is what is in the group, in the language of ctx; "" when it
// has none.
func (b *MenuGroupBuilder) TDescription(ctx context.Context) string {
	if b == nil || b.description == nil {
		return ""
	}
	return b.description(ctx)
}

func (b *MenuGroupBuilder) Icon(v string) (r *MenuGroupBuilder) {
	b.icon = v
	return b
}

// Ancestors are the groups above this one, outermost first. The root sentinel
// is not among them.
func (b *MenuGroupBuilder) Ancestors() (r []*MenuGroupBuilder) {
	for p := b.Parent(); p != nil; p = p.Parent() {
		r = append([]*MenuGroupBuilder{p}, r...)
	}
	return
}

// PathNames are the names from the outermost ancestor down to this group.
func (b *MenuGroupBuilder) PathNames() (r []string) {
	if b == nil || b.item == nil {
		return nil
	}
	for _, a := range b.Ancestors() {
		r = append(r, a.name)
	}
	return append(r, b.name)
}

// URIName sets the group's segment in the URIs of the models under it — its
// name by default. Empty, the group adds no segment: its models' URIs are
// those of its parent's. Only the URIs change: the group's name, its key in
// the menu and the permissions' resources stay its name.
func (b *MenuGroupBuilder) URIName(v string) *MenuGroupBuilder {
	b.uriName = &v
	return b
}

// GetURIName is the group's segment in its models' URIs (URIName).
func (b *MenuGroupBuilder) GetURIName() string {
	if b.uriName != nil {
		return *b.uriName
	}
	return b.name
}

// URIPathNames are the URI segments from the outermost ancestor down to this
// group, the empty ones left out.
func (b *MenuGroupBuilder) URIPathNames() (r []string) {
	if b == nil || b.item == nil {
		return nil
	}
	for _, g := range append(b.Ancestors(), b) {
		if n := g.GetURIName(); n != "" {
			r = append(r, n)
		}
	}
	return
}

// URIPath is what a model under the group prefixes its URI with: the
// URIName of every group of the chain — "a/b/c" for c in b in a.
func (b *MenuGroupBuilder) URIPath() string {
	return strings.Join(b.URIPathNames(), "/")
}

// Path is the group's place in the tree — "a/b/c" for a group c nested in b
// nested in a —, by the groups' names: what the permissions' resources are
// named after. The models' URIs take URIPath, which is the same unless a group
// says otherwise (URIName).
func (b *MenuGroupBuilder) Path() string {
	return strings.Join(b.PathNames(), "/")
}

// contains reports whether g is this group or sits below it. It is what keeps a
// move from making a cycle.
func (b *MenuGroupBuilder) contains(g *MenuGroupBuilder) bool {
	for p := g; p != nil; p = p.rawParent() {
		if p == b {
			return true
		}
	}
	return false
}

// rawParent is the group above, the root sentinel included. contains walks with
// it so the chain does not stop one short of the top.
func (b *MenuGroupBuilder) rawParent() *MenuGroupBuilder {
	if b.item == nil {
		return nil
	}
	return b.item.parent
}

// indexOf is the position of key among the group's items, or -1.
func (b *MenuGroupBuilder) indexOf(key string) int {
	for i, it := range b.items {
		if it.Key() == key {
			return i
		}
	}
	return -1
}

// Add puts items in the group, in the order given. Each is a MenuRef — from
// ModelItem/PageItem/GroupItem — a *MenuGroupBuilder, or, for the older API, a
// string naming a page (leading "/") or a model.
//
// An item that does not exist yet is still placed: it waits here for its
// registration. An item already elsewhere is moved, never duplicated.
func (b *MenuGroupBuilder) Add(items ...any) *MenuGroupBuilder {
	for _, item := range items {
		if err := b.b.MoveMenuItem(menuRefOf(item), b); err != nil {
			panic(err)
		}
	}
	return b
}

// SubItems sets the group's sub-items by name, a page when the name starts with
// "/" and a model otherwise.
//
// It cannot say "the group named x", which is what the ambiguity of a bare
// string costs; use Add with GroupItem for that.
func (b *MenuGroupBuilder) SubItems(ss ...string) (r *MenuGroupBuilder) {
	items := make([]any, len(ss))
	for i, s := range ss {
		items[i] = s
	}
	return b.Add(items...)
}

// SubItemsAny sets the group's sub-items, each a name (string), a MenuRef or a
// nested *MenuGroupBuilder.
func (b *MenuGroupBuilder) SubItemsAny(items ...interface{}) (r *MenuGroupBuilder) {
	return b.Add(items...)
}

// menuRefOf reads one item of the group-building API as a reference. A bare
// string is the older form: a page when it starts with "/", a model otherwise.
func menuRefOf(item any) MenuRef {
	switch it := item.(type) {
	case MenuRef:
		return it
	case *MenuGroupBuilder:
		return GroupItem(it.name)
	case *MenuItem:
		return it.Ref()
	case string:
		if strings.HasPrefix(it, "/") {
			return PageItem(it)
		}
		return ModelItem(it)
	default:
		panic(fmt.Sprintf("unknown menu item type: %T", item))
	}
}
