# Getting started

## The manager

```go
mgr := login_session.NewManager(lb)   // lb is the *login.Builder
```

`NewManager` registers the package's i18n messages and returns a `*Manager` that
reads session config (token, max age) from the login builder.

Auto-migrate the table once:

```go
db.AutoMigrate(&login_session.LoginSession{})
```

## Lifecycle

Wire the manager into the login builder's hooks (see the app `auth` setup):

| Event | Call |
| --- | --- |
| after login | `AddSessionLogByUserID(db, r, userID)` |
| extend session | `UpdateCurrentSessionLog(db, r, userID, oldToken)` |
| logout | `ExpireCurrentSessionLog(db, r, userID)` |
| reset / change password | `ExpireAllSessionLogs(db, userID)` |
| "sign out other devices" | `ExpireOtherSessionLogs(db, r, userID)` |

Each records or expires a `LoginSession` row for the given user uuid. The device
string is parsed from the `User-Agent`; the client IP honours `X-Forwarded-For`.

## Validating a request

Use the manager as the token check in the `user` auth middleware:

```go
mw := userBuilder.Middlewares(db, logoutURL, mgr.CheckIsTokenValidFromRequest)
```

`CheckIsTokenValidFromRequest(db, r, userID)` returns whether the request's
session token matches a live (non-expired) session for that user. The
middleware redirects to the logout URL when it doesn't.

## Listing sessions (profile)

```go
comp, err := mgr.Sessions(db, ctx, userID)
```

Returns an `h.HTMLComponent` listing the user's sessions (device, IP, time,
status) — used by the profile package's login-sessions dialog.

## Helpers

| Function | Purpose |
| --- | --- |
| `GetIP(r)` | Client IP (honours `X-Forwarded-For`, else `RemoteAddr`). |
| `GetProxy(r)` | The `X-Forwarded-For` chain. |
| `GetStringHash(v, len)` | First `len` hex chars of `sha256(v)` (used for `TokenHash`). |
| `IsTokenValid(s)` | Whether a `LoginSession` is past its `ExpiredAt`. |
