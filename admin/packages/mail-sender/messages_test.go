package mail_sender

import (
	"reflect"
	"testing"

	"github.com/go-rvq/rvq/x/i18n/msgpath"
)

// Every message's tag parses — an apostrophe in a hint ('the SMTP's') breaks
// it, and the application panics at boot registering the module.
func TestMessagesTags(t *testing.T) {
	typ := reflect.TypeOf(Messages{})
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if _, err := msgpath.ParseTag(f.Tag); err != nil {
			t.Errorf("%s: %v", f.Name, err)
		}
	}
}
