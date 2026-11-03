# Keep me signed in on this device

A session lasts `SessionMaxAge` (1 hour by default). It is renewed while it is
used (`AutoExtendSession`), so a session unused for that long ends: an
automatic logout.

`Builder.KeepSignedIn(maxAge)` turns on a box in the login form, **"Keep me
signed in on this device"** (the field `keep_signed_in`, `KeepSignedInParam`).
A session that asks for it lasts `maxAge` seconds instead, renewed in the same
way while it is used. `0` turns the option off, and the box is not shown.

```go
loginBuilder.KeepSignedIn(30 * 24 * 60 * 60) // 30 days
```

## Where it goes

- **The claims:** `UserClaims.KeepSignedIn`. The renewal of the middleware
  keeps it, and so does the second factor (TOTP).
- **The cookies of the session:** their max age is the time left until the
  claims expire (`cookieMaxAge`), across restarts of the browser.
- **The max age of a request:** `Builder.SessionMaxAgeOf(r)` is the one of the
  session of `r`. The log of the sessions
  (`admin/packages/login_session`) uses it for `expired_at`, so that check does
  not end a long session either.
- **The initial user:** when it is logged (`initialUserLogged`), its session
  ends at the end of the day, whatever the box says.

## In HermonCMS

`LOGIN_KEEP_SIGNED_IN_DAYS` (default 30; `0` turns it off) in the `env:` of
the `app.yml`.

## Tests

`keep_signed_in_test.go` covers how long a session lasts with and without the
option, the max age of the cookie, and `SessionMaxAgeOf` read from the cookie
of a request.
