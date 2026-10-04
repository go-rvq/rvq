package perm

import "context"

type restrictionKey struct{}

// WithRestriction restricts the requests of ctx to what the policies of
// subjects allow too: a request is allowed what its own subjects allow (its
// user's roles) only when one of subjects is allowed it as well — an
// intersection. The policies of an access key ("key:<id>") restrict so: only
// permissive, they never give the user more than the user has.
func WithRestriction(ctx context.Context, subjects ...string) context.Context {
	return context.WithValue(ctx, restrictionKey{}, append([]string{}, subjects...))
}

// RestrictionOf are the subjects ctx is restricted to (WithRestriction); nil
// when it is not.
func RestrictionOf(ctx context.Context) []string {
	rs, _ := ctx.Value(restrictionKey{}).([]string)
	return rs
}
