// Package timezone carries the time zone an instant is shown and typed in:
// the one of a request (its context), UTC without one. Instants are stored in
// UTC; a time zone is how a reader sees them.
package timezone

import (
	"context"
	"strings"
	"time"
)

type contextKey struct{}

// WithLocation is ctx carrying loc, the time zone of what is shown in it.
func WithLocation(ctx context.Context, loc *time.Location) context.Context {
	if loc == nil {
		return ctx
	}
	return context.WithValue(ctx, contextKey{}, loc)
}

// Location is the time zone ctx carries, UTC without one.
func Location(ctx context.Context) *time.Location {
	if ctx != nil {
		if loc, _ := ctx.Value(contextKey{}).(*time.Location); loc != nil {
			return loc
		}
	}
	return time.UTC
}

// Load is the time zone of the IANA name (America/Sao_Paulo); an empty name
// is UTC.
func Load(name string) (*time.Location, error) {
	if name = strings.TrimSpace(name); name == "" {
		return time.UTC, nil
	}
	return time.LoadLocation(name)
}

// LoadOr is the time zone of name, def when it is empty or not known.
func LoadOr(name string, def *time.Location) *time.Location {
	if strings.TrimSpace(name) != "" {
		if loc, err := time.LoadLocation(strings.TrimSpace(name)); err == nil {
			return loc
		}
	}
	if def == nil {
		return time.UTC
	}
	return def
}

// Common are the usual time zones, offered by a field of one (another may be
// typed).
var Common = []string{
	"UTC", "America/Sao_Paulo", "America/Manaus", "America/Fortaleza", "America/Recife", "America/Belem",
	"America/Cuiaba", "America/Porto_Velho", "America/Rio_Branco", "America/Noronha",
	"America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles",
	"America/Mexico_City", "America/Bogota", "America/Lima", "America/Santiago", "America/Argentina/Buenos_Aires",
	"Europe/Lisbon", "Europe/London", "Europe/Madrid", "Europe/Paris", "Europe/Berlin",
	"Asia/Tokyo", "Asia/Shanghai", "Australia/Sydney",
}
