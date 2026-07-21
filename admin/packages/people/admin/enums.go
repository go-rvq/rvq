package admin

import (
	"github.com/go-rvq/rvq/admin/packages/people/messages"
	"github.com/go-rvq/rvq/admin/packages/people/models"
	"github.com/go-rvq/rvq/admin/presets"
)

// enumSelect builds a SelectConfig for a string-based enum whose labels come
// from the people messages of the request language.
func enumSelect[T ~string](values []T, label func(m *messages.Messages, key string) string) *presets.SelectConfig {
	keys := make([]string, len(values))
	for i, v := range values {
		keys[i] = string(v)
	}
	return &presets.SelectConfig{
		AvailableKeysFunc: func(*presets.FieldContext) []string { return keys },
		KeyLabelsFunc: func(f *presets.FieldContext, ks []string) []string {
			m := messages.Get(f.EventContext.Context())
			out := make([]string, len(ks))
			for i, k := range ks {
				out[i] = label(m, k)
			}
			return out
		},
	}
}

// registerEnumSelects registers translated select components for the people
// enum types (document type).
func registerEnumSelects(b *presets.Builder) {
	presets.RegisterSelectType(b, models.DocumentType(""),
		enumSelect(models.DocumentTypes, func(m *messages.Messages, k string) string {
			switch models.DocumentType(k) {
			case models.DocumentCPF:
				return m.DocumentTypeCPF
			case models.DocumentCNPJ:
				return m.DocumentTypeCNPJ
			case models.DocumentOther:
				return m.DocumentTypeOther
			}
			return k
		}))
}
