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
