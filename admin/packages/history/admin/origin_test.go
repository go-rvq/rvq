package admin

import (
	"context"
	"strings"

	"net/http"
	"net/http/httptest"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"

	"github.com/go-rvq/rvq/admin/model"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
	"github.com/google/uuid"
)

// A revision keeps where its author was: the address — the first of
// X-Forwarded-For behind a proxy — and the browser, of the save's request
// (an admin save) or of the one given to Record (a change outside the
// admin); the place, by SetOriginFunc. Nothing known without a request.
func TestRevisionOrigin(t *testing.T) {
	db, mb, mh := setupHistory(t)

	// an admin save: the request's address and browser
	obj := &histDoc{Title: "A", Body: "<p>one</p>"}
	r := httptest.NewRequest("POST", "/admin/docs", nil)
	r.Header.Set("X-Forwarded-For", "200.1.2.3, 10.0.0.1")
	r.Header.Set("User-Agent", "Firefox/140")
	if err := mb.Editing().Saver(obj, model.ID{}, webCtx(r)); err != nil {
		t.Fatal(err)
	}
	latest := func() histmodels.Revision {
		t.Helper()
		revs, err := mh.Chain(mb.MustRecordID(obj).String())
		if err != nil || len(revs) == 0 {
			t.Fatalf("chain: %v, %d", err, len(revs))
		}
		return revs[len(revs)-1]
	}
	if o := latest().Origin; o.IP != "200.1.2.3" || o.UserAgent != "Firefox/140" || o.City != "" {
		t.Errorf("admin save: %+v", o)
	}

	// Record, the place by the application's origin function
	SetOriginFunc(func(r *http.Request) histmodels.Origin {
		o := DefaultOrigin(r)
		lat, lng := -20.7546, -42.8825
		o.City, o.Region, o.CountryName, o.Country, o.Latitude, o.Longitude = "Viçosa", "Minas Gerais", "Brazil", "BR", &lat, &lng
		return o
	})
	t.Cleanup(func() { SetOriginFunc(nil) })
	obj.Body = "<p>two</p>"
	db.Save(obj)
	r = httptest.NewRequest("POST", "/api/docs", nil)
	r.RemoteAddr = "177.4.5.6:4431"
	if _, created, err := mh.Record(obj, uuid.New(), "Ana", r); err != nil || !created {
		t.Fatalf("record: %v, created %v", err, created)
	}
	o := latest().Origin
	if o.IP != "177.4.5.6" || o.City != "Viçosa" || o.Latitude == nil {
		t.Errorf("record: %+v", o)
	}
	if got, want := o.String(), "Viçosa, Minas Gerais, Brazil (177.4.5.6)"; got != want {
		t.Errorf("origin %q, want %q", got, want)
	}

	// no request: nothing known
	obj.Body = "<p>three</p>"
	db.Save(obj)
	if _, _, err := mh.Record(obj, uuid.Nil, "", nil); err != nil {
		t.Fatal(err)
	}
	if o := latest().Origin; o != (histmodels.Origin{}) {
		t.Errorf("no request: %+v", o)
	}
}

// The Origin action is the revisions' detail's — in the permission tree, and
// in the dialog of its menu —: its body shows the place and, with its
// coordinates, the map.
func TestOriginActionBody(t *testing.T) {
	lat, lng := 51.5142, -0.0931
	body := originBody(Messages_en_US, histmodels.Origin{IP: "81.2.69.160", City: "London", CountryName: "United Kingdom", Latitude: &lat, Longitude: &lng})
	html := renderString(t, body)
	for _, want := range []string{"81.2.69.160", "London, United Kingdom", "51.5142, -0.0931", "rvq-geo-map"} {
		if !contains(html, want) {
			t.Errorf("no %q in the Origin body", want)
		}
	}
	if html := renderString(t, originBody(Messages_en_US, histmodels.Origin{IP: "10.0.0.1"})); contains(html, "rvq-geo-map") || !contains(html, Messages_en_US.OriginNoMap) {
		t.Error("an origin without coordinates has a map")
	}
}

func webCtx(r *http.Request) *web.EventContext { return &web.EventContext{R: r} }

func renderString(t *testing.T, c h.HTMLComponent) string {
	t.Helper()
	return h.MustString(c, context.Background())
}

func contains(s, sub string) bool { return strings.Contains(s, sub) }
