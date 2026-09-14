package web

import (
	"context"
	"net/http"
)

// RequestSourceCLI marks a request that a trusted in-process CLI command
// dispatched into the mux, rather than one that arrived over the network.
const RequestSourceCLI = "cli"

// requestSourceKey carries the request source. It lives only in the Go context
// — never a header — so a network request can never forge it: only in-process
// code holding this package's unexported key can set it.
type requestSourceContextKey int

const requestSourceKey requestSourceContextKey = 0

// WithRequestSource returns a context tagged with the given request source
// (e.g. RequestSourceCLI). Trusted in-process callers use it to record where a
// dispatched request originated.
func WithRequestSource(ctx context.Context, source string) context.Context {
	return context.WithValue(ctx, requestSourceKey, source)
}

// RequestSource returns the source recorded by WithRequestSource, or "" when the
// request arrived over the network without one.
func RequestSource(r *http.Request) string {
	s, _ := r.Context().Value(requestSourceKey).(string)
	return s
}

// IsCLIRequest reports whether r was dispatched by a trusted in-process CLI
// command (see RequestSourceCLI).
func IsCLIRequest(r *http.Request) bool {
	return RequestSource(r) == RequestSourceCLI
}
