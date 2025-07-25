package presets

import (
	"fmt"

	"github.com/go-rvq/rvq/web"
)

type CtxField string

type CtxFieldOptions struct {
	Label       string
	Hint        string
	Hidden      bool
	HiddenLabel bool
	ReadOnly    bool
}

func ContextWithFieldOptions(ctx web.ContextValuer, fqn CtxField, opts *CtxFieldOptions) {
	ctx.WithContextValue(fqn, opts)
}

func FieldOptionsFromContext(ctx web.ContextValuer, fqn CtxField, init ...bool) (opts *CtxFieldOptions) {
	opts, _ = ctx.ContextValue(fqn).(*CtxFieldOptions)
	if opts == nil {
		for _, b := range init {
			if b {
				opts = &CtxFieldOptions{}
			}
		}
	}
	return
}

func ContextWrapFieldOptions(ctx web.ContextValuer, fqn CtxField, f func(opts *CtxFieldOptions)) {
	opts := FieldOptionsFromContext(ctx, fqn)
	if opts == nil {
		opts = &CtxFieldOptions{}
	} else {
		cp := *opts
		opts = &cp
	}
	f(opts)
	ctx.WithContextValue(fqn, opts)
}

func (b *FieldBuilder) NewContext(info *ModelInfo, ctx *web.EventContext, parent *FieldContext, obj any) (fctx *FieldContext) {
	return NewFieldContextBuilder(b, info, ctx, parent, obj).Build()
}

func (b *FieldBuilder) NewContextBuilder(info *ModelInfo, ctx *web.EventContext, parent *FieldContext, obj any) *FieldContextBuilder {
	return NewFieldContextBuilder(b, info, ctx, parent, obj)
}

type FieldContextBuilder struct {
	field  *FieldBuilder
	info   *ModelInfo
	ctx    *web.EventContext
	parent *FieldContext
	obj    any
	slice  bool
	index  int
}

func NewFieldContextBuilder(field *FieldBuilder, info *ModelInfo, ctx *web.EventContext, parent *FieldContext, obj any) *FieldContextBuilder {
	return &FieldContextBuilder{field: field, info: info, ctx: ctx, parent: parent, obj: obj}
}

func (b *FieldContextBuilder) Index(index int) *FieldContextBuilder {
	b.slice = true
	b.index = index
	return b
}

func (b *FieldContextBuilder) Build() (fctx *FieldContext) {
	var (
		keyPath = b.field.name
		path    FieldPath
	)
	if b.parent != nil {
		if b.parent.FormKey != "" {
			keyPath = fmt.Sprintf("%s.%s", b.parent.FormKey, keyPath)
		}
		path = b.parent.Path
	}

	path.Append(b.field.name)

	if b.slice {
		keyPath = fmt.Sprintf("%s[%d]", keyPath, b.index)
		path.AppendIndex(b.index)
	}

	var (
		finfo = FieldOptionsFromContext(b.ctx, CtxField(keyPath), true)
		label string
	)

	if !finfo.Hidden && !finfo.HiddenLabel && !b.field.hiddenLabel {
		label = finfo.Label
		if label == "" {
			label = b.field.ContextLabel(b.info, b.ctx.Context())
		}
	}

	fctx = &FieldContext{
		ToComponentOptions: &ToComponentOptions{},
		Parent:             b.parent,
		ModelInfo:          b.info,
		Obj:                b.obj,
		Field:              b.field,
		EventContext:       b.ctx,
		FormKey:            keyPath,
		Path:               path,
		Name:               b.field.name,
		Label:              label,
		Hint: func() string {
			if finfo.Hint != "" {
				return finfo.Hint
			}
			return b.field.ContextHint(b.info, b.ctx.Context())
		},
		Nested:  b.field.nested,
		Context: b.field.context,
	}

	return
}
