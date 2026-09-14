package login

import (
	"context"
	"net/http"
)

func CurrentUserFromContext(ctx context.Context) any {
	return ctx.Value(UserKey)
}

func CurrentUserFromRequest(r *http.Request) (u interface{}) {
	return CurrentUserFromContext(r.Context())
}

// trustedUserKey marks a request whose user was set by trusted in-process code
// (e.g. a CLI command that dispatches into the mux), not by a login cookie. It
// lives only in the Go context — never a header — so it cannot be forged over
// HTTP: an incoming network request can never carry it.
type trustedUserContextKey int

const trustedUserKey trustedUserContextKey = 0

// WithTrustedUser returns a context carrying user as the pre-authenticated
// current user. The login Middleware honours it verbatim and skips the cookie
// check, so an in-process caller can act as that user without a session.
func WithTrustedUser(ctx context.Context, user any) context.Context {
	ctx = context.WithValue(ctx, trustedUserKey, user)
	return context.WithValue(ctx, UserKey, user)
}

// trustedUserFromContext returns the trusted user set by WithTrustedUser, and
// whether one was set.
func trustedUserFromContext(ctx context.Context) (any, bool) {
	u := ctx.Value(trustedUserKey)
	return u, u != nil
}

// HasTrustedUser reports whether ctx carries a pre-authenticated user set by
// trusted in-process code (WithTrustedUser). Such a user has no login session,
// so session-token validation must be skipped for it.
func HasTrustedUser(ctx context.Context) bool {
	_, ok := trustedUserFromContext(ctx)
	return ok
}
