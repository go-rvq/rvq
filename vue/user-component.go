package vue

import (
	"context"
	"strings"

	"github.com/qor5/web/v3/js"
	h "github.com/theplant/htmlgo"
)

type UserComponentAssigner struct {
	Dst    Var
	Values js.Object
	Merges []string
}

func (a *UserComponentAssigner) String() string {
	var s []string
	for _, merge := range a.Merges {
		s = append(s, "..."+merge)
	}
	if len(a.Values) > 0 {
		s = append(s, "..."+a.Values.String())
	}
	return "{" + strings.Join(s, ", ") + "}"
}

func (a *UserComponentAssigner) Set(key string, value any) *UserComponentAssigner {
	if a.Values == nil {
		a.Values = make(js.Object)
	}
	a.Values[key] = value
	return a
}

func (a *UserComponentAssigner) Merge(value string) *UserComponentAssigner {
	a.Merges = append(a.Merges, value)
	return a
}

type UserComponentBuilder struct {
	scopeNames  []string
	scopeValues [][]any
	setupFuncs  js.RawSlice
	onUmount    string
	onMounted   string
	assign      map[Var]*UserComponentAssigner
	*h.HTMLTagBuilder
}

func UserComponent(children ...h.HTMLComponent) *UserComponentBuilder {
	return &UserComponentBuilder{HTMLTagBuilder: h.Tag("user-component").Children(h.Tag("template").Children(children...))}
}

func (b *UserComponentBuilder) GetChildren() []h.HTMLComponent {
	return b.HTMLTagBuilder.Childs
}

func (b *UserComponentBuilder) Scope(name string, value ...any) *UserComponentBuilder {
	b.scopeNames = append(b.scopeNames, name)
	b.scopeValues = append(b.scopeValues, value)
	return b
}

func (b *UserComponentBuilder) ScopeVar(name string, value string) *UserComponentBuilder {
	return b.Scope(name, Var(value))
}

func (b *UserComponentBuilder) Setup(s string) *UserComponentBuilder {
	b.setupFuncs = append(b.setupFuncs, s)
	return b
}

func (b *UserComponentBuilder) OnMounted(s string) *UserComponentBuilder {
	b.onMounted = s
	return b
}

func (b *UserComponentBuilder) OnUmount(s string) *UserComponentBuilder {
	b.onUmount = s
	return b
}

func (b *UserComponentBuilder) Assigner(dst Var) *UserComponentAssigner {
	if b.assign == nil {
		b.assign = make(map[Var]*UserComponentAssigner)
	}
	if assigner, ok := b.assign[dst]; ok {
		return assigner
	}
	assigner := &UserComponentAssigner{Dst: dst}
	b.assign[dst] = assigner
	return assigner
}

func (b *UserComponentBuilder) Assign(dst Var, key string, val any) *UserComponentBuilder {
	b.Assigner(dst).Set(key, val)
	return b
}

func (b *UserComponentBuilder) AssignMany(dst Var, val string) *UserComponentBuilder {
	b.Assigner(dst).Merge(val)
	return b
}

func (b *UserComponentBuilder) Component() *h.HTMLTagBuilder {
	return b.HTMLTagBuilder
}

func (b *UserComponentBuilder) Template() *h.HTMLTagBuilder {
	return b.Childs[0].(*h.HTMLTagBuilder)
}

func (b *UserComponentBuilder) AppendChild(h ...h.HTMLComponent) *UserComponentBuilder {
	b.Template().AppendChildren(h...)
	return b
}

func (b *UserComponentBuilder) MarshalHTML(ctx context.Context) ([]byte, error) {
	scopeValues := make([]string, len(b.scopeValues))

	for i, v := range b.scopeValues {
		if len(v) == 1 {
			switch v := v[0].(type) {
			case Var:
				scopeValues[i] = string(v)
			default:
				scopeValues[i] = h.JSONString(v)
			}
		} else {

		}
	}

	comp := b.HTMLTagBuilder
	template := b.Template()

	if len(scopeValues) > 0 {
		var scope []string
		for i, name := range b.scopeNames {
			v := scopeValues[i]
			if len(v) > 0 {
				scope = append(scope, name+": "+v)
			}
		}
		comp.Attr(":scope", "{"+strings.Join(scope, ", ")+"}")
		template.Attr("v-slot", "{"+strings.Join(b.scopeNames, ", ")+"}")
	} else {
		comp.Childs = template.Childs
	}

	if len(b.assign) > 0 {
		var (
			assign = make([]string, len(b.assign))
			i      int
		)

		for _, a := range b.assign {
			v := a.String()
			assign[i] = "[" + string(a.Dst) + "," + v + "]"
			i++
		}

		comp.Attr(":assign", "["+strings.Join(assign, ", ")+"]")
	}

	if len(b.setupFuncs) > 0 {
		comp.Attr(":setup", b.setupFuncs.String())
	}

	if b.onMounted != "" {
		comp.Attr("@mounted", b.onMounted)
	}

	if b.onUmount != "" {
		comp.Attr("@unmount", b.onUmount)
	}

	return comp.MarshalHTML(ctx)
}
