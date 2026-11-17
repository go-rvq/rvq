package fs_tools

import (
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/x/i18n/msgpath"
)

// Every message's tag parses — a wrong one panics the application at boot,
// registering the module.
func TestMessagesTags(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeOf(Messages{}), reflect.TypeOf(FilesMessages{})} {
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if _, err := msgpath.ParseTag(f.Tag); err != nil {
				t.Errorf("%s.%s: %v", typ.Name(), f.Name, err)
			}
		}
	}
}

// Every word of the components has its translation, and keeps the
// placeholders of vue-i18n ({n}) of the English.
func TestFilesMessages(t *testing.T) {
	en, pt := reflect.ValueOf(*FilesMessages_en_US), reflect.ValueOf(*FilesMessages_pt_BR)
	ph := regexp.MustCompile(`\{\w+\}`)
	for i := 0; i < en.NumField(); i++ {
		name := en.Type().Field(i).Name
		e, p := en.Field(i).String(), pt.Field(i).String()
		if e == "" || p == "" {
			t.Errorf("%s: empty", name)
			continue
		}
		if a, b := ph.FindAllString(e, -1), ph.FindAllString(p, -1); strings.Join(a, ",") != strings.Join(b, ",") {
			t.Errorf("%s: placeholders %v, %v", name, a, b)
		}
	}
}
