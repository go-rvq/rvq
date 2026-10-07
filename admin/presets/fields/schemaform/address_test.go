package schemaform

import (
	"encoding/json"
	"net/url"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/web"
)

// An address is drawn by <vx-address-field> with the builder's key (else the
// package's), bound to the field; posted as its fields, it is read back as a
// record of AddressFields — the coordinates numbers —, nil when its words are
// empty.
func TestAddress(t *testing.T) {
	b := New().MapsKey(func(*web.EventContext) string { return "k-123" })
	got := render(t, b, `interface { [label="Where"] where? address }`)
	for _, want := range []string{"<vx-address-field", "label='Where'", "api-key='k-123'", `v-model='form["Value"].where'`, "data-schemaform-address='where'"} {
		if !strings.Contains(got, want) {
			t.Errorf("no %s:\n%s", want, got)
		}
	}

	was := MapsKey
	MapsKey = func(*web.EventContext) string { return "pkg-key" }
	t.Cleanup(func() { MapsKey = was })
	if got := render(t, New(), `interface { where? address }`); !strings.Contains(got, "api-key='pkg-key'") {
		t.Errorf("the package's key:\n%s", got)
	}

	s, err := New().Parse(`interface { where? address }`)
	if err != nil {
		t.Fatal(err)
	}
	v := url.Values{"Value.where.formatted": {"Rua A, 10 - Viçosa, MG"}, "Value.where.city": {"Viçosa"},
		"Value.where.lat": {"-20.75"}, "Value.where.lng": {"-42.88"}, "Value.where.placeId": {"p1"}}
	out, _ := json.Marshal(New().decoders[AddressType](v, "Value.where"))
	if want := `{"formatted":"Rua A, 10 - Viçosa, MG","street":"","number":"","complement":"","neighborhood":"","city":"Viçosa","region":"","postalCode":"","country":"","countryCode":"","lat":-20.75,"lng":-42.88,"placeId":"p1"}`; string(out) != want {
		t.Errorf("decoded:\n%s\nwant\n%s", out, want)
	}
	if rec := s.Decode(url.Values{"Value.where.city": {"x"}}, "Value"); !strings.Contains(string(mustJSON(t, rec)), `"where":null`) {
		t.Errorf("no words, no address: %s", mustJSON(t, rec))
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
