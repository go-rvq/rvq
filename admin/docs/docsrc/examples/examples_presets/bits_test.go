package examples_presets

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web/multipartestutils"
	"github.com/go-rvq/rvq/x/i18n"
)

// The example of a field of bits: the chips of its channels, each with the
// description of the messages; the set ones marked in the detail.
func TestPresetsBitsField(t *testing.T) {
	pb := presets.New(i18n.New())
	PresetsBitsField(pb, TestDB)
	s := &Subscriber{Name: "Ana", Channels: ChannelEmail | ChannelPush}
	if err := TestDB.Create(s).Error; err != nil {
		t.Fatal(err)
	}
	do := func(event string) string {
		t.Helper()
		w := httptest.NewRecorder()
		pb.ServeHTTP(w, multipartestutils.NewMultipartBuilder().
			PageURL(fmt.Sprintf("/subscribers?__execute_event__=%s&id=%d", event, s.ID)).BuildEventFuncRequest())
		if w.Code != 200 {
			t.Fatalf("%s: %d", event, w.Code)
		}
		return w.Body.String()
	}
	edit := do("presets_Edit")
	for _, want := range []string{"v-chip-group", "Email", "SMS", "Push", "Notified by a text message", "mdi-help-circle-outline"} {
		if !strings.Contains(edit, want) {
			t.Errorf("edit: no %q", want)
		}
	}
	if n := strings.Count(do("presets_Detailing"), "data-bit-set"); n != 2 {
		t.Errorf("detail: %d marked, want 2", n)
	}
}
