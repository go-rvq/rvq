# user

Admin user management for the `go-rvq/rvq` stack: it wires a user model into
`presets` (listing / editing / detailing with roles, status and account), seeds
the initial administrator, manages roles, and provides the auth middlewares and
the current-user helper.

Users are identified by `uuid.UUID` (`User.GetID()`), matching the rest of the
stack (activity, notes, sessions, sharing).

## Quick start

```go
import "github.com/go-rvq/rvq/admin/packages/user"

// mb is the presets.ModelBuilder for your *User model; lb is the login.Builder.
ub := user.New(db, lb, mb, loginInitialUserEmail)
ub.GenInitialUser(db)   // create the initial admin (from lb's InitialUserAccount/Password)
```

Your model must implement the `user.User` interface:

```go
type User interface {
    login.UserPasser
    GetID() uuid.UUID
    GetName() string
    SetName(string)
    SetEmail(string)
    SetRegistrationDate(time.Time)
    GetStatus() string
    GetAccountName() string
    GetRoles() role.Roles
    SetRoles(role.Roles)
}
```

## Documentation

Detailed docs live under [`docs/`](docs/):

- [Getting started](docs/getting-started.md) — the `User` interface, `New`, model registration, the initial user.
- [Roles & middlewares](docs/roles-and-middlewares.md) — role resolution, granting roles, and the HTTP middlewares.
- [API reference](docs/api-reference.md).

## What `New` configures

`New(db, lb, mb, loginInitialUserEmail)` sets up the user resource:

- **Listing** — columns (ID, Name, Email, Status, Notes) + filters (created,
  name, status, registration date); non-admins cannot see admin/manager users.
- **Editing** — Name, Account (email), Roles (multi-select), Status; the initial
  user account is read-only and cannot be edited, re-roled or deleted.
- **Detailing** — roles view + actions: change password, send reset-password
  email, unlock a locked user, revoke TOTP.
- Registration date is stamped on save.
