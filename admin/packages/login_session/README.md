# login_session

Server-side login-session tracking for the `go-rvq/rvq` admin: it records one
row per active login (device, IP, hashed token, expiry) so the app can list a
user's sessions, validate the current one, and revoke sessions on logout,
password reset or from the profile page.

Sessions are keyed by the user's `uuid.UUID`.

## Quick start

```go
import "github.com/go-rvq/rvq/admin/packages/login_session"

mgr := login_session.NewManager(lb)

// on login:
mgr.AddSessionLogByUserID(db, r, user.GetID())

// as an auth middleware check:
valid, err := mgr.CheckIsTokenValidFromRequest(db, r, user.GetID())

// on logout:
mgr.ExpireCurrentSessionLog(db, r, user.GetID())
```

Auto-migrate the model:

```go
db.AutoMigrate(&login_session.LoginSession{})
```

## Documentation

- [Getting started](docs/getting-started.md)
- [API reference](docs/api-reference.md)

## The model

```go
type LoginSession struct {
    gorm.Model
    UserID    uuid.UUID
    Device    string      // "<browser> - <os>", parsed from the user agent
    IP        string
    TokenHash string      // first 8 hex chars of sha256(session token)
    ExpiredAt time.Time
    // Time / Status are computed (gorm:"-") for display
}
```

The full session token is never stored — only a short hash (`TokenHash`) used to
match the current request's token against the logged sessions.
