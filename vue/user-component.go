package vue

import (
	"context"
	"strings"

	"github.com/qor5/web/v3/tag"
	h "github.com/theplant/htmlgo"
)

type UserComponentAssigner struct {
	Dst    Var
	Values map[string]any
}

func (a *UserComponentAssigner) Set(key string, value any) *UserComponentAssigner {
	if a.Values == nil {
		a.Values = make(map[string]any)
	}
	a.Values[key] = value
	return a
}

type UserComponentBuilder struct {
	scopeNames  []string
	scopeValues [][]any
	setupFunc   string
	onUmount    string
	onMounted   string
	assign      map[Var]*UserComponentAssigner
	*h.HTMLTagBuilder
}

func UserComponent(children ...h.HTMLComponent) *UserComponentBuilder {
	return &UserComponentBuilder{HTMLTagBuilder: h.Tag("user-component").Children(h.Tag("template").Children(children...))}
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
	b.setupFunc = s
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

func (b *UserComponentBuilder) Component() *h.HTMLTagBuilder {
	return b.HTMLTagBuilder
}

func (b *UserComponentBuilder) Template() *h.HTMLTagBuilder {
	return (*tag.Children(b.HTMLTagBuilder))[0].(*h.HTMLTagBuilder)
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
		children := tag.Children(comp)
		*children = *tag.Children(template)
	}

	if len(b.assign) > 0 {
		var (
			assign = make([]string, len(b.assign))
			i      int
		)

		for _, a := range b.assign {
			assign[i] = "[" + string(a.Dst) + "," + h.JSONString(a.Values) + "]"
			i++
		}

		comp.Attr(":assign", "["+strings.Join(assign, ", ")+"]")
	}

	if b.setupFunc != "" {
		comp.Attr(":setup", b.setupFunc)
	}

	if b.onMounted != "" {
		comp.Attr("@mounted", b.onMounted)
	}

	if b.onUmount != "" {
		comp.Attr("@unmount", b.onUmount)
	}

	return comp.MarshalHTML(ctx)
}
