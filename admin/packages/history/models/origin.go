package models

import (
	"fmt"
	"strings"
)

// Origin is where a revision was made from: the author's address and browser
// and — when the application locates addresses (admin SetOriginFunc) — where
// the address is, approximate (a city's). Zero when nothing is known: a
// revision of the boot's seed, one recorded with no request.
type Origin struct {
	IP          string `gorm:"type:varchar(64)"`
	UserAgent   string
	Country     string `gorm:"type:varchar(2)"` // ISO code: "BR"
	CountryName string
	Region      string
	City        string
	Latitude    *float64
	Longitude   *float64
}

// Place is where the origin is in words — "Viçosa, Minas Gerais, Brazil" —;
// "" when it is not known.
func (o Origin) Place() string {
	country := o.CountryName
	if country == "" {
		country = o.Country
	}
	var parts []string
	for _, s := range []string{o.City, o.Region, country} {
		if s != "" {
			parts = append(parts, s)
		}
	}
	return strings.Join(parts, ", ")
}

// String is the origin on one line: its place and its address — "Viçosa,
// Minas Gerais, Brazil (200.1.2.3)" —, the address alone when the place is
// not known; "" when nothing is.
func (o Origin) String() string {
	switch place := o.Place(); {
	case place != "" && o.IP != "":
		return fmt.Sprintf("%s (%s)", place, o.IP)
	case place != "":
		return place
	}
	return o.IP
}
