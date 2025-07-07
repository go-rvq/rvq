package web

import (
	"context"
	"fmt"
	"sort"
	"strings"

	h "github.com/go-rvq/htmlgo"
	. "github.com/go-rvq/rvq/web/tag"
	"github.com/rs/xid"
)

type ScopeBuilder struct {
	TagBuilder[*ScopeBuilder]
	observers []Observer
	slots     map[string]any
	initVal   map[string][]string
	name      string
}

func Scope(children ...h.HTMLComponent) *ScopeBuilder {
	return NewTag(&ScopeBuilder{}, "go-plaid-scope", children...)
}

func (b *ScopeBuilder) Locals() *ScopeBuilder {
	return b.Slot("locals")
}

func (b *ScopeBuilder) Form() *ScopeBuilder {
	return b.Slot("form")
}

func (b *ScopeBuilder) Vars() *ScopeBuilder {
	return b.Slot("vars")
}

func (b *ScopeBuilder) Closer() *ScopeBuilder {
	return b.Slot("closer")
}

func (b *ScopeBuilder) Fullscreen() *ScopeBuilder {
	return b.Slot("fullscreen")
}

func (b *ScopeBuilder) init(attr string, vs ...interface{}) (r *ScopeBuilder) {
	if b.initVal == nil {
		b.initVal = make(map[string][]string)
	}
	if len(vs) == 0 {
		b.initVal[attr] = append(b.initVal[attr], "{}")
	} else {
		for _, v := range vs {
			switch vt := v.(type) {
			case string:
				b.initVal[attr] = append(b.initVal[attr], vt)
			default:
				b.initVal[attr] = append(b.initVal[attr], h.JSONString(v))
			}
		}
	}
	return b
}

func (b *ScopeBuilder) LocalsInit(vs ...interface{}) (r *ScopeBuilder) {
	return b.init(":locals", vs...).Locals()
}

func (b *ScopeBuilder) FormInit(vs ...interface{}) (r *ScopeBuilder) {
	return b.init(":form", vs...).Form()
}

func (b *ScopeBuilder) CloserInit(vs ...interface{}) (r *ScopeBuilder) {
	return b.init(":closer", vs...).Closer()
}

func (b *ScopeBuilder) FullscreenInit(vs ...interface{}) (r *ScopeBuilder) {
	return b.init(":fullscreen", vs...).Fullscreen()
}

func (b *ScopeBuilder) Setup(callback string) (r *ScopeBuilder) {
	return b.init(":setup", []any{callback})
}

func (b *ScopeBuilder) Name(name string) (r *ScopeBuilder) {
	b.name = name
	return b
}

func (b *ScopeBuilder) OnChange(v string) (r *ScopeBuilder) {
	return b.Attr("@change-debounced", fmt.Sprintf(`({locals, form, oldLocals, oldForm}) => { %s }`, v)).
		Attr(":use-debounce", 800)
}

func (b *ScopeBuilder) UseDebounce(v int) (r *ScopeBuilder) {
	return b.Attr(":use-debounce", v)
}

func (b *ScopeBuilder) AppendInit(v string) (r *ScopeBuilder) {
	t := b.GetAttr(":init")
	if t == nil {
		return b.LocalsInit(v)
	}
	t.Override(func(old any) any {
		s := old.(string)
		s = strings.TrimSpace(s[1 : len(s)-1])
		if s != "" {
			s += ", "
		}
		return "{" + s + v + "}"
	})
	return b
}

func (b *ScopeBuilder) Slot(v ...string) (r *ScopeBuilder) {
	if b.slots == nil {
		b.slots = make(map[string]any)
	}
	for _, v := range v {
		v = strings.TrimSpace(v)
		if v[0] == '{' {
			v = v[1 : len(v)-1]
		}
		for _, v := range strings.Split(v, ",") {
			v = strings.TrimSpace(v)
			b.slots[v] = nil
		}
	}
	return b
}

func (b *ScopeBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	if len(b.observers) > 0 {
		b.Attr(":observers", h.JSONString(b.observers))
	}
	if len(b.slots) > 0 {
		var names []string
		for k := range b.slots {
			names = append(names, k)
		}
		sort.Strings(names)
		b.Attr("v-slot", "{ "+strings.Join(names, ", ")+" }")
	}
	if len(b.initVal) > 0 {
		var names []string
		for k := range b.initVal {
			names = append(names, k)
		}
		sort.Strings(names)

		for _, name := range names {
			v := "[" + strings.Join(b.initVal[name], ", ") + "]"
			b.Attr(name, v)
		}
	}
	if b.name == "" {
		b.name = xid.New().String()
	}
	b.Attr("scope-name", b.name)
	return b.TagBuilder.MarshalHTML(ctx)
}

type Observer struct {
	Name   string `json:"name"`
	Script string `json:"script"` // available parameters: name payload vars locals form plaid
}

func (b *ScopeBuilder) Observer(name string, script string) (r *ScopeBuilder) {
	b.observers = append(b.observers, Observer{name, script})
	return b
}

func (b *ScopeBuilder) Observers(vs ...Observer) (r *ScopeBuilder) {
	b.observers = append(b.observers, vs...)
	return b
}

func CloserScope(comp h.HTMLComponent, show ...bool) (r *ScopeBuilder) {
	if r, _ = comp.(*ScopeBuilder); r == nil {
		r = Scope(
			comp,
		)
	}

	r.Closer().
		CloserInit()

	for _, s := range show {
		if s {
			r.Setup("setTimeout(() => {closer.show = true},100)")
		}
	}

	return
}

type notification struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Payload any    `json:"payload"`
}

func NotifyScript(name string, payload any) string {
	return fmt.Sprintf(`vars.__notification = %s`, h.JSONString(notification{
		ID:      xid.New().String(),
		Name:    name,
		Payload: payload,
	}))
}
