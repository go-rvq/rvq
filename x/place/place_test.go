package place

import "testing"

// A place in words, its coordinates and its map; nothing of a place not known.
func TestPlace(t *testing.T) {
	lat, lng := -20.75461, -42.88253
	p := Place{Country: "BR", CountryName: "Brazil", Region: "Minas Gerais", City: "Viçosa", Latitude: &lat, Longitude: &lng}
	if !p.Located() || p.String() != "Viçosa, Minas Gerais, Brazil" || p.Coordinates() != "-20.7546, -42.8825" ||
		p.MapURL() != "https://www.openstreetmap.org/?mlat=-20.7546&mlon=-42.8825#map=11/-20.7546/-42.8825" {
		t.Errorf("%v · %q · %q · %q", p.Located(), p.String(), p.Coordinates(), p.MapURL())
	}
	if s := (Place{Country: "US"}).String(); s != "US" {
		t.Errorf("only the country's code: %q", s)
	}
	var none Place
	if none.Located() || none.String() != "" || none.Coordinates() != "" || none.MapURL() != "" {
		t.Error("a place not known says something")
	}
}
