package admin

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/origin"
	histmodels "github.com/go-rvq/rvq/admin/packages/history/models"
)

// OriginFunc is where the request r — a save, a Record — came from: the
// revision's Origin (origin.Func).
type OriginFunc = origin.Func

// SetOriginFunc sets how the revisions — and every other user of the origin
// package, as the trash — learn where their author was: origin.SetFunc.
func SetOriginFunc(f OriginFunc) { origin.SetFunc(f) }

// DefaultOrigin is the address r came from and its browser: origin.Default.
func DefaultOrigin(r *http.Request) histmodels.Origin { return origin.Default(r) }

// ClientIP is the address r came from: origin.ClientIP.
func ClientIP(r *http.Request) string { return origin.ClientIP(r) }

// originOf is where r came from; zero with no request.
func originOf(r *http.Request) histmodels.Origin { return origin.Of(r) }
