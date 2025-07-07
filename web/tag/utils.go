package tag

import (
	h "github.com/go-rvq/htmlgo"
)

func FirstValidComponent(c h.HTMLComponent) h.HTMLComponent {
	switch t := c.(type) {
	case h.HTMLComponents:
		for _, comp := range t {
			if comp != nil {
				return FirstValidComponent(comp)
			}
		}
	case interface{ GetChildren() []h.HTMLComponent }:
		return FirstValidComponent(h.HTMLComponents(t.GetChildren()))
	case *h.HTMLTagBuilder:
		return FirstValidComponent(h.HTMLComponents((&TagBuilder[any]{tag: t}).GetChildren()))
	}
	return c
}

// Simplify Simplifies components walking over nested HTMLComponents and calls cb if
// component not is nil
func Simplify(c h.HTMLComponent, cb func(c h.HTMLComponent)) {
	if c == nil {
		return
	}
	switch t := c.(type) {
	case h.HTMLComponents:
		for _, c := range t {
			Simplify(c, cb)
		}
	default:
		cb(c)
	}
}

type WalkState uint8

const (
	SkipNext WalkState = iota + 1
	SkipAll
)

func WalkS(c h.HTMLComponent, cb func(c h.HTMLComponent) (state WalkState)) (state WalkState) {
	var items []h.HTMLComponent
	switch t := c.(type) {
	case h.HTMLComponents:
		items = t
	case interface{ GetChildren() []h.HTMLComponent }:
		switch cb(c) {
		case SkipNext:
			return 0
		case SkipAll:
			return SkipAll
		}
		items = t.GetChildren()
	case *h.HTMLTagBuilder:
		items = t.Childs
	}

	for _, comp := range items {
		if comp != nil {
			state = WalkS(comp, cb)
			switch state {
			case SkipNext:
				return 0
			case SkipAll:
				return
			}
		}
	}
	return
}

func Walk(c h.HTMLComponent, cb func(c h.HTMLComponent) (state WalkState)) {
	WalkS(c, cb)
}
