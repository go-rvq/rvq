# profile

The logged-in user's own profile for the `go-rvq/rvq` admin: a singleton
"My Profile" resource, the sidebar brand/account widget, and the login-sessions
dialog (sign out of other devices).

It builds on [`user`](../user) (current user) and
[`login_session`](../login_session) (session management).

## Quick start

```go
import "github.com/go-rvq/rvq/admin/packages/profile"

pf := profile.New(db, sessionMgr, lb, userMb)
pf.Install(presetsBuilder)                 // registers the "my-profile" singleton
presetsBuilder.BrandFunc(pf.Brand(presetsBuilder))   // account widget in the sidebar
```

## What it provides

- **`Install`** — registers a singleton `Profile` model (id `my-profile`, menu
  hidden) that shows the current user's name, account, status and roles;
  read-only except the editable name; grants the `Logged` role access to it.
- **`Brand`** — the sidebar component: avatar + name + first role, a logout
  button, and an optional notification slot; shows an "Entrar" (sign in) button
  when nobody is logged in.
- **`ProfilePage`** — a standalone profile page with "Login sessions" and
  "Sign out" actions.
- **Login-sessions dialog** — lists the user's active sessions and lets them
  sign out of all of them.

## Documentation

- [Getting started](docs/getting-started.md)
- [API reference](docs/api-reference.md)
