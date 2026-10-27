package integration_test

import (
	"fmt"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/admin/example/admin"
	"github.com/go-rvq/rvq/admin/example/models"
	. "github.com/go-rvq/rvq/web/multipartestutils"
)

// The example of a field of bits (InputDemo.Notifications1, BitsComponentFunc):
// to edit, a chip a channel with the `?` of its description, the value the sum
// of the bits; in the detail, the set ones marked; saved, the bits chosen.
func TestInputDemoBits(t *testing.T) {
	h := admin.TestHandler(TestDB)
	d := &models.InputDemo{Notifications1: models.NotifyEmail | models.NotifyPush}
	if err := TestDB.Create(d).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { TestDB.Delete(d) })
	do := func(event string, fields ...string) string {
		t.Helper()
		b := NewMultipartBuilder().PageURL(fmt.Sprintf("/featured-models/input-demos?__execute_event__=%s&id=%d", event, d.ID))
		for i := 0; i+1 < len(fields); i += 2 {
			b = b.AddField(fields[i], fields[i+1])
		}
		w := httptest.NewRecorder()
		h.ServeHTTP(w, b.BuildEventFuncRequest())
		if w.Code != 200 {
			t.Fatalf("%s: %d %.300s", event, w.Code, w.Body.String())
		}
		return w.Body.String()
	}

	edit := do("presets_Edit")
	for _, want := range []string{"v-chip-group", "Email", "SMS", "Push", "Notifies by a text message",
		"mdi-help-circle-outline", "Notifications1"} {
		if !strings.Contains(edit, want) {
			t.Errorf("edit: no %q", want)
		}
	}

	detail := do("presets_Detailing")
	if n := strings.Count(detail, "data-bit-set"); n != 2 {
		t.Errorf("detail: %d chips marked, want 2 (email, push)", n)
	}

	// the form's stamp (its __formSign), as the form sends it
	sign := regexp.MustCompile(`__formSign\\?":\s*\\?"([^"\\]+)`).FindStringSubmatch(edit)
	if sign == nil {
		t.Fatal("the form carries no __formSign")
	}
	do("presets_Update", "__formSign", sign[1], "Notifications1", fmt.Sprint(int(models.NotifySMS)))
	var got models.InputDemo
	TestDB.First(&got, d.ID)
	if got.Notifications1 != models.NotifySMS {
		t.Errorf("saved: %03b, want SMS", got.Notifications1)
	}
}
