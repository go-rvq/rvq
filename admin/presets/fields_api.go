package presets

type (
	SetObjectFieldsHandler interface {
		Handler(fromObj interface{}, toObj interface{}, parent *FieldContext) (err error)
	}

	SetObjectFieldsHandlerFunc func(fromObj interface{}, toObj interface{}, parent *FieldContext) (err error)
	SetObjectFieldsHandlers    []SetObjectFieldsHandler
)

func (f SetObjectFieldsHandlerFunc) Handler(fromObj interface{}, toObj interface{}, parent *FieldContext) (err error) {
	return f(fromObj, toObj, parent)
}

func (h SetObjectFieldsHandlers) Handler(fromObj interface{}, toObj interface{}, parent *FieldContext) (err error) {
	for _, handler := range h {
		if err = handler.Handler(fromObj, toObj, parent); err != nil {
			return
		}
	}
	return
}

func (h *SetObjectFieldsHandlers) Append(handler ...SetObjectFieldsHandler) {
	*h = append(*h, handler...)
}

func (h *SetObjectFieldsHandlers) AppendFunc(handler ...SetObjectFieldsHandlerFunc) {
	for _, f := range handler {
		h.Append(f)
	}
}
