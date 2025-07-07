package presets

import (
	"context"
	"fmt"
)

type Menu []interface{}

func (b *Menu) Order(items ...interface{}) {
	for _, item := range items {
		switch v := item.(type) {
		case string:
			*b = append(*b, v)
		case *MenuGroupBuilder:
			if b.isMenuGroupInOrder(v) {
				b.removeMenuGroupInOrder(v)
			}
			*b = append(*b, v)
		default:
			panic(fmt.Sprintf("unknown menu order item type: %T\n", item))
		}
	}
}

func (b *Menu) Group(name string) *MenuGroupBuilder {
	mgb := b.Group(name)
	if !b.isMenuGroupInOrder(mgb) {
		*b = append(*b, mgb)
	}
	return mgb
}

func (b *Menu) isMenuGroupInOrder(mgb *MenuGroupBuilder) bool {
	for _, v := range *b {
		if v == mgb {
			return true
		}
	}
	return false
}

func (b *Menu) removeMenuGroupInOrder(mgb *MenuGroupBuilder) {
	for i, om := range *b {
		if om == mgb {
			*b = append((*b)[:i], (*b)[i+1:]...)
			break
		}
	}
}

type MenuGroupBuilder struct {
	title func(ctx context.Context) string
	name  string
	icon  string
	// item can be Slug name, model name
	// the underlying logic is using Slug name,
	// so if the Slug name is customized, item must be the Slug name
	subMenuItems []string
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
	if b.title != nil {
		return b.title(ctx)
	}
	return HumanizeString(b.name)
}

func (b *MenuGroupBuilder) Icon(v string) (r *MenuGroupBuilder) {
	b.icon = v
	return b
}

func (b *MenuGroupBuilder) SubItems(ss ...string) (r *MenuGroupBuilder) {
	b.subMenuItems = ss
	return b
}

type MenuGroups struct {
	menuGroups []*MenuGroupBuilder
}

func (g *MenuGroups) MenuGroup(name string) (r *MenuGroupBuilder) {
	for _, mg := range g.menuGroups {
		if mg.name == name {
			return mg
		}
	}
	r = &MenuGroupBuilder{name: name}
	g.menuGroups = append(g.menuGroups, r)
	return
}
