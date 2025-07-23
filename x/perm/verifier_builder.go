package perm

import (
	"context"
	"iter"
	"sort"
	"strings"
)

type PermVerifierFunc func(v *Verifier) *Verifier

type ActionsWithTitleMap map[string]func(context context.Context) string

type PermVerifierBuilder struct {
	v       *Verifier
	vFunc   PermVerifierFunc
	title   func(context context.Context) string
	actions ActionsWithTitleMap
}

func PermVerifier(v ...*Verifier) *PermVerifierBuilder {
	b := &PermVerifierBuilder{}
	for _, v := range v {
		b.v = v
	}
	return b
}

func (b *PermVerifierBuilder) ClearActions() *PermVerifierBuilder {
	b.actions = nil
	return b
}

func (b *PermVerifierBuilder) Action(name string, title func(context context.Context) string) *PermVerifierBuilder {
	if b.actions == nil {
		b.actions = make(ActionsWithTitleMap)
	}
	b.actions[name] = title
	return b
}

func (b *PermVerifierBuilder) GetActions() ActionsWithTitleMap {
	return b.actions
}

func (b *PermVerifierBuilder) HasActions() bool {
	return len(b.actions) > 0
}

func (b *PermVerifierBuilder) ActionVerifier(action string) *PermVerifierBuilder {
	b2 := *b
	b2.actions = nil
	b2.title = b.actions[action]
	b2.vFunc = func(v *Verifier) *Verifier {
		return b.BuildDo(v, action)
	}
	return &b2
}

func (b *PermVerifierBuilder) SplitForActions() iter.Seq[*PermVerifierBuilder] {
	if len(b.actions) == 0 {
		return func(yield func(*PermVerifierBuilder) bool) {}
	}

	var actions []string
	for name := range b.actions {
		actions = append(actions, name)
	}
	sort.Strings(actions)

	return func(yield func(*PermVerifierBuilder) bool) {
		for _, action := range actions {
			if !yield(b.ActionVerifier(action)) {
				return
			}
		}
	}
}

func (b *PermVerifierBuilder) AcceptAction(action string) bool {
	if len(b.actions) == 0 {
		return true
	}
	_, ok := b.actions[action]
	return ok
}

func (b *PermVerifierBuilder) Func(verifier PermVerifierFunc) *PermVerifierBuilder {
	b.vFunc = verifier
	return b
}

func (b *PermVerifierBuilder) GetFunc() PermVerifierFunc {
	return b.vFunc
}

func (b *PermVerifierBuilder) Verifier(v *Verifier) *PermVerifierBuilder {
	b.v = v
	return b
}

func (b *PermVerifierBuilder) GetVerifier() *Verifier {
	return b.v
}

func (b *PermVerifierBuilder) Title(f func(ctx context.Context) string) *PermVerifierBuilder {
	b.title = f
	return b
}

func (b *PermVerifierBuilder) TitleString(s string) *PermVerifierBuilder {
	return b.Title(func(ctx context.Context) string {
		return s
	})
}

func (b *PermVerifierBuilder) GetTitle() func(context context.Context) string {
	return b.title
}

func (b *PermVerifierBuilder) TTitle(context context.Context) string {
	if b.title == nil {
		return ""
	}
	return b.title(context)
}

func (b *PermVerifierBuilder) Valid() bool {
	return b.v != nil || b.vFunc != nil
}

func (b *PermVerifierBuilder) Wrap(wrap func(old PermVerifierFunc) PermVerifierFunc) *PermVerifierBuilder {
	b.vFunc = wrap(b.vFunc)
	return b
}

func (b *PermVerifierBuilder) Path(pth string) *PermVerifierBuilder {
	b.vFunc = func(v *Verifier) *Verifier {
		return v.On("/" + strings.Trim(pth, "/"))
	}
	return b
}

func (b *PermVerifierBuilder) BuildDo(dot *Verifier, action string) *Verifier {
	v := b.Build(dot)
	if len(b.actions) == 0 || b.actions[action] != nil {
		return v.On(action)
	}
	return v.Deny()
}

func (b *PermVerifierBuilder) Build(dot *Verifier) *Verifier {
	if b.v != nil {
		return b.v
	}
	if b.vFunc != nil {
		return b.vFunc(dot)
	}
	return dot
}

type PermVerifiers []*PermVerifierBuilder

func (b *PermVerifiers) Add(v ...*PermVerifierBuilder) {
	*b = append(*b, v...)
}

type PermVerifierBuilderNode struct {
	Parent *PermVerifierBuilderNode
	Elem   *PermVerifierBuilder
}

func WalkPermVerififierBuilders(s ...iter.Seq[*PermVerifierBuilder]) iter.Seq[*PermVerifierBuilderNode] {
	return func(yield func(*PermVerifierBuilderNode) bool) {
		for _, s := range s {
			if !WalkPermVerififierBuildersTree(nil, s, yield) {
				break
			}
		}
	}
}

func WalkPermVerififierBuildersTree(parent *PermVerifierBuilderNode, s iter.Seq[*PermVerifierBuilder], yield func(*PermVerifierBuilderNode) bool) bool {
	for b := range s {
		node := &PermVerifierBuilderNode{
			Parent: parent,
			Elem:   b,
		}

		if !yield(node) {
			return false
		}

		if !WalkPermVerififierBuildersTree(node, b.SplitForActions(), yield) {
			return false
		}
	}
	return true
}
