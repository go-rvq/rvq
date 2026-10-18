package vuetifyx

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// The namespace vx of the vue-i18n: the module's words, and the months of
// the i18n's default module (one place for them), January first.
func TestVueMessages(t *testing.T) {
	b, err := json.Marshal(VueMessages(context.Background()))
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{`"filterClearAll":"Clear Filters"`, `"monthNames":["January",`, `"December"]`,
		`"abbrMonthNames":["Jan",`, `"filterBy":"Filter by {filter}"`} {
		if !strings.Contains(s, want) {
			t.Errorf("no %s in %.300s", want, s)
		}
	}
	if strings.Contains(s, "ModuleDescription") || strings.Contains(s, `"monthNames":[""`) {
		t.Errorf("the description, or the month 0: %.300s", s)
	}
}
