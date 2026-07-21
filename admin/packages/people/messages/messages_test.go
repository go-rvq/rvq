package messages

import (
	"reflect"
	"testing"
)

// TestMessagesCompleteness ensures the pt-BR and en-US variants fill every field
// of the Messages struct, so no label renders blank.
func TestMessagesCompleteness(t *testing.T) {
	for name, m := range map[string]*Messages{"pt_BR": Messages_pt_BR, "en_US": Messages_en_US} {
		v := reflect.ValueOf(*m)
		tp := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).String() == "" {
				t.Errorf("%s: field %s is empty", name, tp.Field(i).Name)
			}
		}
	}
}
