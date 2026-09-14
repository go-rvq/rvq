package user

import (
	"net/http"

	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/x/login"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Middlewares struct {
	b  *Builder
	db *gorm.DB

	logoutURL                    string
	checkIsTokenValidFromRequest func(db *gorm.DB, r *http.Request, userID uuid.UUID) (valid bool, err error)
	dev                          bool
}

func (b *Builder) Middlewares(db *gorm.DB, logoutURL string, checkIsTokenValidFromRequest func(db *gorm.DB, r *http.Request, userID uuid.UUID) (valid bool, err error)) *Middlewares {
	return &Middlewares{
		b:                            b,
		db:                           db,
		logoutURL:                    logoutURL,
		checkIsTokenValidFromRequest: checkIsTokenValidFromRequest,
	}
}

// loginBuilder is the login these middlewares belong to, when there is one.
func (b *Middlewares) loginBuilder() *login.Builder {
	if b.b == nil {
		return nil
	}
	return b.b.lb
}

func (b *Middlewares) SetDevMode(v bool) *Middlewares {
	b.dev = v
	return b
}

func (b *Middlewares) DevMode() bool {
	return b.dev
}

func (b *Middlewares) WithRoles(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := GetCurrentUser(r)
		// The anonymous user carries no roles (it is "not logged in"); skip the
		// role lookup as for a nil user.
		if u == nil || u.Anonymous() {
			next.ServeHTTP(w, r)
			return
		}

		var roleIDs []uint
		if err := b.db.Table("user_role_join").Select("role_id").Where("user_id=?", u.GetID()).Scan(&roleIDs).Error; err != nil {
			panic(err)
		}
		if len(roleIDs) > 0 {
			var roles []*role.Role
			if err := b.db.Where("id in (?)", roleIDs).Find(&roles).Error; err != nil {
				panic(err)
			}
			u.SetRoles(roles)
		}
		next.ServeHTTP(w, r)
	})
}

func (b *Middlewares) WithRolesMD() func(next http.Handler) http.Handler {
	return b.WithRoles
}

func (b *Middlewares) Security(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Add("Cache-control", "no-cache, no-store, max-age=0, must-revalidate")
		w.Header().Add("Pragma", "no-cache")

		next.ServeHTTP(w, req)
	})
}

func (b *Middlewares) SecurityMD() func(next http.Handler) http.Handler {
	return b.Security
}

func (b *Middlewares) ValidateSessionToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetCurrentUser(r)
		// The anonymous user means "nobody is logged in" — it has no session to
		// validate, so treat it like nil here (otherwise validation fails and the
		// request is bounced to logout, looping the login page).
		if user == nil || user.Anonymous() {
			next.ServeHTTP(w, r)
			return
		}
		// A trusted in-process user (WithTrustedUser — e.g. a CLI command) is
		// pre-authenticated and has no login session, so there is no session token
		// to validate; skipping it here avoids bouncing the request to logout.
		if login.HasTrustedUser(r.Context()) {
			next.ServeHTTP(w, r)
			return
		}
		if login.IsLoginWIP(r) {
			next.ServeHTTP(w, r)
			return
		}

		// the login page and its neighbours have no session to validate, and
		// checking one there would leave nowhere to go
		if b.loginBuilder() != nil && b.loginBuilder().WhiteListed(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		valid, err := b.checkIsTokenValidFromRequest(b.db, r, user.GetID())
		if err != nil || !valid {
			if r.URL.Path == b.logoutURL {
				next.ServeHTTP(w, r)
				return
			}
			// The token is still good but the session behind it ended (an expiry,
			// a logout from somewhere else). For a request fired by a page the
			// user is working on, going to the login page would take the page —
			// and everything typed into it — along, so that one is answered in
			// place. See login.UnauthorizedResponder.
			if b.loginBuilder() != nil && b.loginBuilder().RespondUnauthorized(w, r) {
				return
			}
			http.Redirect(w, r, b.logoutURL, http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (b *Middlewares) ValidateSessionTokenMD() func(next http.Handler) http.Handler {
	return b.ValidateSessionToken
}

func (b *Middlewares) Middleware(next http.Handler) http.Handler {
	return b.ValidateSessionToken(b.WithRoles(b.Security(next)))
}

func (b *Middlewares) MiddlewareMD() func(next http.Handler) http.Handler {
	if b.dev {
		return func(next http.Handler) http.Handler {
			return next
		}
	}
	return b.Middleware
}
