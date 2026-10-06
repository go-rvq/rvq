// Package place is where an address is: the country, region and city, and
// the coordinates — approximate, a city's, never an address —, as an
// application's geolocation knows it (hermon's geo, a GeoLite2 database).
package place

import (
	"fmt"
	"strings"
)

// Place is where an address is; zero when not known. Embedded
// (`gorm:"embedded"`), its columns are country, country_name, region, city,
// latitude and longitude.
type Place struct {
	Country     string // ISO code: "BR"
	CountryName string
	Region      string
	City        string
	Latitude    *float64
	Longitude   *float64
}

// Located reports whether anything is known of the place.
func (p Place) Located() bool { return p.Country != "" || p.City != "" }

// String is the place in words: "Viçosa, Minas Gerais, Brazil"; "" when it is
// not known.
func (p Place) String() string {
	country := p.CountryName
	if country == "" {
		country = p.Country
	}
	var parts []string
	for _, s := range []string{p.City, p.Region, country} {
		if s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, ", ")
}

// Coordinates are the place's, "-20.7546, -42.8825"; "" when not known.
func (p Place) Coordinates() string {
	if p.Latitude == nil || p.Longitude == nil {
		return ""
	}
	return fmt.Sprintf("%.4f, %.4f", *p.Latitude, *p.Longitude)
}

// MapURL is the place on OpenStreetMap; "" without coordinates.
func (p Place) MapURL() string {
	if p.Latitude == nil || p.Longitude == nil {
		return ""
	}
	return fmt.Sprintf("https://www.openstreetmap.org/?mlat=%.4f&mlon=%.4f#map=11/%.4f/%.4f",
		*p.Latitude, *p.Longitude, *p.Latitude, *p.Longitude)
}
