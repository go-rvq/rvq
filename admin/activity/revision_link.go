package activity

import (
	"context"
	"net/http"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
)

// RevisionRef points at the history revision the current save produced.
type RevisionRef struct {
	Table string
	Hash  []byte
}

type revisionRefKey struct{}

// WithRevisionRef records, on r, the revision the current save produced. The
// history plugin calls it right after capturing; the edit log then references
// that revision instead of storing a duplicate diff. Returns the updated
// request (its context carries the ref).
func WithRevisionRef(r *http.Request, table string, hash []byte) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), revisionRefKey{}, &RevisionRef{Table: table, Hash: hash}))
}

// revisionRefFromContext returns the revision ref set for the current save, or
// nil when history did not (or has not yet) captured one.
func revisionRefFromContext(ctx context.Context) *RevisionRef {
	ref, _ := ctx.Value(revisionRefKey{}).(*RevisionRef)
	return ref
}

// RevisionDiffFunc renders the change carried by a history revision
// (table + hash). The history plugin sets it in its Configure, so the activity
// log view can show the diff from the revision without activity importing
// history (which would be an import cycle). Nil when history is not mounted.
var RevisionDiffFunc func(table string, hash []byte, ctx *web.EventContext) h.HTMLComponent
