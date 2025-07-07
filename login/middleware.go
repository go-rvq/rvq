package login

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ContextUserKey int

const (
	UserKey ContextUserKey = iota
	loginWIPKey
)

type MiddlewareConfig interface {
	middlewareConfig()
}

// LoginNotRequired executes the next handler regardless of whether the user is logged in or not
type LoginNotRequired struct{}

func (*LoginNotRequired) middlewareConfig() {}

// DisableAutoRedirectToHomePage makes it possible to visit login page when user is logged in
type DisableAutoRedirectToHomePage struct{}

func (*DisableAutoRedirectToHomePage) middlewareConfig() {}

func MockCurrentUser(user any) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserKey, user)))
		})
	}
}

func (b *Builder) ParseClaims(r *http.Request) (*UserClaims, error) {
	return parseUserClaimsFromCookie(r, b.authCookieName, b.secret)
}

func (b *Builder) ParseRequestUser(r *http.Request) (user any, err, userStateErr error) {
	var claims *UserClaims
	if claims, err = b.ParseClaims(r); err != nil {
		return
	}

	defer func() {
		if err != nil || userStateErr != nil {
			user = nil
		}
	}()

	var secureSalt string
	if b.userModel != nil {
		if user, err = b.findUserByID(claims.UserID); err != nil {
			if err == ErrUserNotFound {
				err = nil
				return
			}
			return
		}

		if claims.Provider == "" {
			if user.(UserPasser).GetPasswordUpdatedAt() != claims.PassUpdatedAt {
				userStateErr = ErrPasswordChanged
				return
			}
			if user.(UserPasser).GetLocked() {
				userStateErr = ErrUserLocked
				return
			}
		} else {
			user.(OAuthUser).SetAvatar(claims.AvatarURL)
		}

		if b.sessionSecureEnabled {
			secureSalt = user.(SessionSecurer).GetSecure()
			if _, err = parseBaseClaimsFromCookie(r, b.authSecureCookieName, b.secret+secureSalt); err != nil {
				return
			}
		}
	}

	return
}

func (b *Builder) BasichAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			claims, err = b.ParseClaims(r)

			user       any
			secureSalt string

			setMessage = func(msg string) {
				w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm=%q`, msg))
				http.Error(w, msg, http.StatusUnauthorized)
			}
			errorMessage = func(err error) {
				setMessage("ERROR: " + err.Error())
			}
		)

		if _, ok := b.whiteList[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}

		if err != nil {
			if err.Error() == "no token string" {
				if b.userPassEnabled {
					if user, pass, ok := r.BasicAuth(); ok {
						if user, err := b.authUserPass(user, pass); err != nil {
							errorMessage(err)
							return
						} else {
							u := user.(UserPasser)
							userID := objectID(user)
							claims = &UserClaims{
								UserID:           userID,
								PassUpdatedAt:    u.GetPasswordUpdatedAt(),
								RegisteredClaims: b.genBaseSessionClaim(userID, u.GetAccountName() == b.initialUserAccount),
							}

							if user, err = b.findUserByID(claims.UserID); err != nil {
								errorMessage(err)
								return
							} else if err = b.setSecureCookiesByClaims(w, user, *claims); err != nil {
								errorMessage(err)
								return
							}

							goto ok
						}
					}
				}
				setMessage("Authentication")
			} else {
				errorMessage(err)
			}
			return
		}

		if b.userModel != nil {
			var err error
			user, err = b.findUserByID(claims.UserID)
			if err == nil {
				if claims.Provider == "" {
					if user.(UserPasser).GetPasswordUpdatedAt() != claims.PassUpdatedAt {
						err = ErrPasswordChanged
					}
					if user.(UserPasser).GetLocked() {
						err = ErrUserLocked
					}
				}
			}

			if err != nil {
				errorMessage(err)
				return
			}

			if b.sessionSecureEnabled {
				secureSalt = user.(SessionSecurer).GetSecure()
				_, err := parseBaseClaimsFromCookie(r, b.authSecureCookieName, b.secret+secureSalt)
				if err != nil {
					errorMessage(err)
					return
				}
			}
		} else {
			user = claims
		}

	ok:
		if b.autoExtendSession && time.Now().Sub(claims.IssuedAt.Time).Seconds() > float64(b.sessionMaxAge)/10 {
			oldSessionToken := b.mustGetSessionToken(*claims)

			claims.RegisteredClaims = b.genBaseSessionClaim(claims.UserID, user.(UserPasser).GetAccountName() != b.initialUserAccount)
			b.setAuthCookiesFromUserClaims(w, claims, secureSalt)

			if b.afterExtendSessionHook != nil {
				setCookieForRequest(r, &http.Cookie{Name: b.authCookieName, Value: b.mustGetSessionToken(*claims)})
				if err := b.afterExtendSessionHook(r, user, oldSessionToken); err != nil {
					errorMessage(err)
					return
				}
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), UserKey, user))

		next.ServeHTTP(w, r)
	})
}

func (b *Builder) Middleware(cfgs ...MiddlewareConfig) func(next http.Handler) http.Handler {
	mustLogin := true
	autoRedirectToHomePage := true
	for _, cfg := range cfgs {
		switch cfg.(type) {
		case *LoginNotRequired:
			mustLogin = false
		case *DisableAutoRedirectToHomePage:
			autoRedirectToHomePage = false
		}
	}

	b.WhiteList(
		b.oauthBeginURL,
		b.oauthCallbackURL,
		b.oauthCallbackCompleteURL,
		b.passwordLoginURL,
		b.forgetPasswordPageURL,
		b.sendResetPasswordLinkURL,
		b.resetPasswordLinkSentPageURL,
		b.resetPasswordURL,
		b.resetPasswordPageURL,
		b.validateTOTPURL,
	)

	staticFileRe := regexp.MustCompile(`\.(css|js|gif|jpg|jpeg|png|ico|svg|ttf|eot|woff|woff2|map;?)$`)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if staticFileRe.MatchString(strings.ToLower(r.URL.Path)) {
				next.ServeHTTP(w, r)
				return
			}

			if _, ok := b.whiteList[r.URL.Path]; ok {
				next.ServeHTTP(w, r)
				return
			}

			var (
				path        = strings.TrimRight(r.URL.Path, "/")
				claims, err = b.ParseClaims(r)

				user       any
				secureSalt string
			)

			if err != nil {
				if !mustLogin {
					next.ServeHTTP(w, r)
					return
				}
				if r.Method == http.MethodGet {
					b.setContinueURL(w, r)
				}
				if path == b.loginPageURL || !b.requireForRequest(r) {
					next.ServeHTTP(w, r)
				} else {
					http.Redirect(w, r, b.loginPageURL, http.StatusFound)
				}
				return
			}

			if b.userModel != nil {
				var err error
				user, err = b.findUserByID(claims.UserID)
				if err == nil {
					if claims.Provider == "" {
						if user.(UserPasser).GetPasswordUpdatedAt() != claims.PassUpdatedAt {
							err = ErrPasswordChanged
						}
						if user.(UserPasser).GetLocked() {
							err = ErrUserLocked
						}
					} else {
						user.(OAuthUser).SetAvatar(claims.AvatarURL)
					}
				}
				if err != nil {
					if !mustLogin {
						next.ServeHTTP(w, r)
						return
					}
					switch err {
					case ErrUserNotFound:
						setFailCodeFlash(w, FailCodeUserNotFound)
					case ErrUserLocked:
						setFailCodeFlash(w, FailCodeUserLocked)
					case ErrPasswordChanged:
						isSelfChange := false
						if c, err := r.Cookie(infoCodeFlashCookieName); err == nil {
							v, _ := strconv.Atoi(c.Value)
							if InfoCode(v) == InfoCodePasswordSuccessfullyChanged {
								isSelfChange = true
							}
						}
						if !isSelfChange {
							setWarnCodeFlash(w, WarnCodePasswordHasBeenChanged)
						}
					default:
						panic(err)
					}
					if path == b.logoutURL {
						next.ServeHTTP(w, r)
					} else {
						http.Redirect(w, r, b.logoutURL, http.StatusFound)
					}
					return
				}

				if b.sessionSecureEnabled {
					secureSalt = user.(SessionSecurer).GetSecure()
					_, err := parseBaseClaimsFromCookie(r, b.authSecureCookieName, b.secret+secureSalt)
					if err != nil {
						if !mustLogin {
							next.ServeHTTP(w, r)
							return
						}
						if path == b.logoutURL {
							next.ServeHTTP(w, r)
						} else {
							http.Redirect(w, r, b.logoutURL, http.StatusFound)
						}
						return
					}
				}
			} else {
				user = claims
			}

			if b.autoExtendSession && time.Now().Sub(claims.IssuedAt.Time).Seconds() > float64(b.sessionMaxAge)/10 {
				oldSessionToken := b.mustGetSessionToken(*claims)

				claims.RegisteredClaims = b.genBaseSessionClaim(claims.UserID, user.(UserPasser).GetAccountName() != b.initialUserAccount)
				b.setAuthCookiesFromUserClaims(w, claims, secureSalt)

				if b.afterExtendSessionHook != nil {
					setCookieForRequest(r, &http.Cookie{Name: b.authCookieName, Value: b.mustGetSessionToken(*claims)})
					if herr := b.afterExtendSessionHook(r, user, oldSessionToken); herr != nil {
						if !mustLogin {
							next.ServeHTTP(w, r)
							return
						}
						setNoticeOrPanic(w, herr)
						http.Redirect(w, r, b.logoutURL, http.StatusFound)
						return
					}
				}
			}

			r = r.WithContext(context.WithValue(r.Context(), UserKey, user))

			if path == b.logoutURL {
				next.ServeHTTP(w, r)
				return
			}

			if claims.Provider == "" && b.totpEnabled {
				if !user.(UserPasser).GetIsTOTPSetup() {
					if path == b.loginPageURL {
						next.ServeHTTP(w, r)
						return
					}
					r = r.WithContext(context.WithValue(r.Context(), loginWIPKey, true))
					if path == b.totpSetupPageURL {
						next.ServeHTTP(w, r)
						return
					}
					http.Redirect(w, r, b.totpSetupPageURL, http.StatusFound)
					return
				}

				if !claims.TOTPValidated {
					if path == b.loginPageURL {
						next.ServeHTTP(w, r)
						return
					}
					r = r.WithContext(context.WithValue(r.Context(), loginWIPKey, true))
					if path == b.totpValidatePageURL {
						next.ServeHTTP(w, r)
						return
					}
					http.Redirect(w, r, b.totpValidatePageURL, http.StatusFound)
					return
				}
			}

			if autoRedirectToHomePage {
				if path == b.loginPageURL || path == b.totpSetupPageURL || path == b.totpValidatePageURL {
					http.Redirect(w, r, b.homePageURLFunc(r, user), http.StatusFound)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetCurrentUser(r *http.Request) (u interface{}) {
	return r.Context().Value(UserKey)
}

// IsLoginWIP indicates whether the user is in an intermediate step of login process,
// such as on the TOTP validation page
func IsLoginWIP(r *http.Request) bool {
	v, ok := r.Context().Value(loginWIPKey).(bool)
	if !ok {
		return false
	}
	return v
}
