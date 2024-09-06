package presets

import (
	"fmt"

	"github.com/qor5/web/v3"
)

type (
	FieldWalkState uint8

	FieldWalkHandle func(field *FieldContext) (s FieldWalkState)

	FieldWalker interface {
		Walk(fctx *FieldContext, handle FieldWalkHandle) (s FieldWalkState)
	}
)

const (
	FieldWalkNext FieldWalkState = iota
	FieldWalkSkipSiblings
	FieldWalkSkipChildren
	FieldWalkStop
)

func (b *FieldsBuilder) Walk(info *ModelInfo, obj interface{}, mode FieldModeStack, ctx *web.EventContext, handle FieldWalkHandle) {
	b.walk(info, obj, mode, "", ctx, handle)
}

func (b *FieldsBuilder) walk(info *ModelInfo, obj interface{}, mode FieldModeStack, parentFormValueKey string, ctx *web.EventContext, handle FieldWalkHandle) (s FieldWalkState) {
	var (
		layout     = b.CurrentLayout()
		fieldsChan = make(chan string)
	)

	if info != nil && !info.mb.singleton && !mode.Dot().Is(LIST, DETAIL) {
		id, _, _ := info.GetID(obj)
		if id.IsZero() {
			mode = append(mode, NEW)
		} else {
			mode = append(mode, EDIT)
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
		s = b.walkField(info, obj, mode, parentFormValueKey, ctx, fieldName, handle)
		if s == FieldWalkStop {
			return s
		} else if s == FieldWalkSkipSiblings {
			return FieldWalkNext
		}
	}
	return
}

func (b *FieldsBuilder) walkField(info *ModelInfo, obj interface{}, mode FieldModeStack, parentFormValueKey string, ctx *web.EventContext, name string, handle FieldWalkHandle) (s FieldWalkState) {
	var (
		f              = b.GetFieldOrDefault(name)
		contextKeyPath = f.name
	)

	if parentFormValueKey != "" {
		contextKeyPath = fmt.Sprintf("%s.%s", parentFormValueKey, f.name)
	}

	fctx := &FieldContext{
		Field:        f,
		Mode:         mode,
		Obj:          obj,
		EventContext: ctx,
		ModelInfo:    info,
		Name:         f.name,
		FormKey:      contextKeyPath,
		Nested:       f.nested,
		Context:      f.context,
	}

	s = handle(fctx)
	if s == FieldWalkSkipChildren {
		s = FieldWalkNext
	} else if s == FieldWalkNext {
		if f.nested != nil {
			s = f.nested.Walk(fctx, handle)
		}
	}
	return
}
