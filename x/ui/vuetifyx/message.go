package vuetifyx

import "encoding/json"

type MessageDetails struct {
	Raw      string
	Template *ComponentDefinition
}

func (m *MessageDetails) MarshalJSON() ([]byte, error) {
	if m.Template != nil {
		return json.Marshal(m.Template)
	}
	return json.Marshal(m.Raw)
}

type Message struct {
	Value  string          `json:"value,omitempty"`
	Type   string          `json:"type,omitempty"`
	Detail *MessageDetails `json:"detail,omitempty"`
}
