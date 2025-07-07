package presets

import (
	"fmt"

	"github.com/go-rvq/rvq/web"
)

type (
	FieldWalkState uint8

	FieldWalkHandleOptions struct {
		SkipPerssionCheck bool
		SkipNestedNil     bool
		SkipMode          bool
		InitializeSlices  bool
		InitializeObjects bool
		Handler           FieldWalkHandle
	}

	FieldWalkHandle func(field *FieldContext) (s FieldWalkState)

	FieldWalker interface {
		Walk(fctx *FieldContext, opts *FieldWalkHandleOptions) (s FieldWalkState)
	}
)

const (
	FieldWalkNext FieldWalkState = iota
	FieldWalkSkipSiblings
	FieldWalkSkipChildren
	FieldWalkStop
)

func (b *FieldsBuilder) Walk(info *ModelInfo, obj interface{}, mode FieldModeStack, ctx *web.EventContext, handle FieldWalkHandle) {
	b.walk(info, obj, mode, nil, "", ctx, &FieldWalkHandleOptions{Handler: handle})
}

func (b *FieldsBuilder) WalkOptions(info *ModelInfo, obj interface{}, mode FieldModeStack, ctx *web.EventContext, opts *FieldWalkHandleOptions) {
	b.walk(info, obj, mode, nil, "", ctx, opts)
}

func (b *FieldsBuilder) walk(info *ModelInfo, obj interface{}, mode FieldModeStack, path FieldPath, parentFormValueKey string, ctx *web.EventContext, opts *FieldWalkHandleOptions) (s FieldWalkState) {
	var (
		layout     = b.CurrentLayout()
		fieldsChan = make(chan string)
	)

	if obj == nil && opts.InitializeObjects {
		obj = b.model
	}

	// if not is embedded
	if !opts.SkipMode && info != nil && len(info.Schema().PrimaryFields()) > 0 {
		if !info.mb.singleton && !mode.Dot().Is(LIST, DETAIL) {
			id, _, _ := info.LookupID(obj)
			if id.IsZero() {
				mode = append(mode, NEW)
			} else {
				mode = append(mode, EDIT)
			}
		}
	}

	go func() {
		defer close(fieldsChan)

		for _, iv := range layout {
			switch t := iv.(type) {
			case string:
				fieldsChan <- t

			case []string:
				for _, s2 := range t {
					fieldsChan <- s2
				}
			case *FieldsSection:
				for _, row := range t.Rows {
					for _, n := range row {
						fieldsChan <- n
					}
				}
			default:
				panic("unknown fields layout, must be string/[]string/*FieldsSection")
			}
		}
	}()

	for fieldName := range fieldsChan {
		s = b.walkField(info, obj, mode, path, parentFormValueKey, ctx, fieldName, opts)
		if s == FieldWalkStop {
			return s
		} else if s == FieldWalkSkipSiblings {
			return FieldWalkNext
		}
	}
	return
}

func (b *FieldsBuilder) walkField(info *ModelInfo, obj interface{}, mode FieldModeStack, path FieldPath, parentFormValueKey string, ctx *web.EventContext, name string, opts *FieldWalkHandleOptions) (s FieldWalkState) {
	var (
		f              = b.GetFieldOrDefault(name)
		contextKeyPath = f.name
	)

	if parentFormValueKey != "" {
		contextKeyPath = fmt.Sprintf("%s.%s", parentFormValueKey, f.name)
	}

	path.Append(f.name)

	fctx := &FieldContext{
		ToComponentOptions: &ToComponentOptions{},
		Field:              f,
		Mode:               mode,
		Obj:                obj,
		EventContext:       ctx,
		ModelInfo:          info,
		Name:               f.name,
		Path:               path,
		FormKey:            contextKeyPath,
		Nested:             f.nested,
		Context:            f.context,
	}

	s = opts.Handler(fctx)
	if s == FieldWalkSkipChildren {
		s = FieldWalkNext
	} else if s == FieldWalkNext {
		if f.nested != nil {
			s = f.nested.Walk(fctx, opts)
		}
	}
	return
}
