package presets

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/timezone"
)

type timeLocationRecord struct {
	At  time.Time
	Day *time.Time
}

// TestTimeLocation: an instant typed with no zone is of the time zone of the
// admin, stored in UTC, and shown back in it; a day has no time zone.
func TestTimeLocation(t *testing.T) {
	sp, _ := time.LoadLocation("America/Sao_Paulo")
	b := New(i18n.New()).TimeLocationFunc(func(*http.Request) *time.Location { return sp })
	r := httptest.NewRequest("POST", "/", nil)
	r.Form = url.Values{"At": {"2025-05-01 21:30"}, "Day": {"2025-05-01"}}
	r = b.withTimeLocation(r)
	if timezone.Location(r.Context()) != sp {
		t.Fatal("the request carries the time zone of the admin")
	}
	ctx := &web.EventContext{R: r}
	var rec timeLocationRecord
	if err := DateTimeComponentFuncSetter(&rec, &FieldContext{Name: "At", FormKey: "At"}, ctx); err != nil {
		t.Fatal(err)
	}
	if got := rec.At.Format(time.RFC3339); got != "2025-05-02T00:30:00Z" {
		t.Errorf("21:30 in São Paulo, in UTC: %s", got)
	}
	if got := ToStringContext(ctx, rec.At); got != "2025-05-01 21:30:00" {
		t.Errorf("shown in São Paulo: %s", got)
	}
	if got := ToStringContext(nil, rec.At); got != "2025-05-02 00:30:00" {
		t.Errorf("without a request, UTC: %s", got)
	}
	if err := DateComponentFuncSetter(&rec, &FieldContext{Name: "Day", FormKey: "Day"}, ctx); err != nil {
		t.Fatal(err)
	}
	if got := rec.Day.Format(time.RFC3339); got != "2025-05-01T00:00:00Z" {
		t.Errorf("a day, at midnight UTC: %s", got)
	}
	if New(i18n.New()).TimeLocation(r) != time.UTC {
		t.Error("an admin without one: UTC")
	}
}
