package presets

type FieldFilter func(f *FieldBuilder) bool

type FieldFilters []FieldFilter

func (ff FieldFilters) Accept(f *FieldBuilder) bool {
	for _, filter := range ff {
		if !filter(f) {
			return false
		}
	}
	return true
}

func FieldRenderable() FieldFilter {
	return func(f *FieldBuilder) bool {
		return f.compFunc != nil
	}
}

func (b FieldBuilders) Filter(filter ...FieldFilter) (out FieldBuilders) {
	filters := FieldFilters(filter)
	for _, f := range b {
		if !filters.Accept(f) {
			continue
		}
		out = append(out, f)
	}
	return
}
