package presets

import (
	"context"
	"net/http"
	"time"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/timezone"
)

// TimeLocationFunc sets the time zone of the admin: the one its instants
// (time.Time fields, filters) are shown and typed in, by request — f returns
// nil for UTC. They are stored in UTC.
func (b *Builder) TimeLocationFunc(f func(r *http.Request) *time.Location) *Builder {
	b.timeLocationFunc = f
	return b
}

// TimeLocation is the time zone of the admin for r, UTC without one.
func (b *Builder) TimeLocation(r *http.Request) *time.Location {
	if b.timeLocationFunc != nil {
		if loc := b.timeLocationFunc(r); loc != nil {
			return loc
		}
	}
	return time.UTC
}

// withTimeLocation is r carrying the time zone of the admin (a request that
// carries one already keeps it).
func (b *Builder) withTimeLocation(r *http.Request) *http.Request {
	if b.timeLocationFunc == nil {
		return r
	}
	if _, ok := r.Context().Value(timeLocationSetKey{}).(bool); ok {
		return r
	}
	ctx := timezone.WithLocation(r.Context(), b.TimeLocation(r))
	return r.WithContext(context.WithValue(ctx, timeLocationSetKey{}, true))
}

type timeLocationSetKey struct{}

// TimeLocationOf is the time zone of the admin in the event ctx, UTC without
// one.
func TimeLocationOf(ctx *web.EventContext) *time.Location {
	if ctx == nil {
		return time.UTC
	}
	return timezone.Location(ctx.Context())
}

// InTimeLocation is t in the time zone of the admin in ctx.
func InTimeLocation(ctx *web.EventContext, t time.Time) time.Time {
	return t.In(TimeLocationOf(ctx))
}
