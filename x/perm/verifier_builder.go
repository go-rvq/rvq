package perm

import (
	"context"
	"strings"
)

type PermVerifierFunc func(v *Verifier) *Verifier

type PermVerifierBuilder struct {
	v     *Verifier
	vFunc PermVerifierFunc
	title func(context context.Context) string
}

func PermVerifier(v ...*Verifier) *PermVerifierBuilder {
	b := &PermVerifierBuilder{}
	for _, v := range v {
		b.v = v
	}
	return b
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
