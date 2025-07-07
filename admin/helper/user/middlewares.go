package user

import (
	"net/http"

	"github.com/qor5/admin/v3/role"
	"github.com/qor5/x/v3/login"
	"gorm.io/gorm"
)

type Middlewares struct {
	b  *Builder
	db *gorm.DB

	logoutURL                    string
	checkIsTokenValidFromRequest func(db *gorm.DB, r *http.Request, userID uint) (valid bool, err error)
	dev                          bool
}

func (b *Builder) Middlewares(db *gorm.DB, logoutURL string, checkIsTokenValidFromRequest func(db *gorm.DB, r *http.Request, userID uint) (valid bool, err error)) *Middlewares {
	return &Middlewares{
		b:                            b,
		db:                           db,
		logoutURL:                    logoutURL,
		checkIsTokenValidFromRequest: checkIsTokenValidFromRequest,
	}
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
		if u == nil {
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
		if user == nil {
			next.ServeHTTP(w, r)
			return
		}
		if login.IsLoginWIP(r) {
			next.ServeHTTP(w, r)
			return
		}

		valid, err := b.checkIsTokenValidFromRequest(b.db, r, user.GetID())
		if err != nil || !valid {
			if r.URL.Path == b.logoutURL {
				next.ServeHTTP(w, r)
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
