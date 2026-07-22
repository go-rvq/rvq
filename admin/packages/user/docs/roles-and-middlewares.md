# Roles & middlewares

## Roles

Two built-in role names are defined by the package:

| Constant | Value |
| --- | --- |
| `user.RoleAdministrador` | `"Administrador"` |
| `user.RoleLogged` | `"Logged"` |

### Resolving a user's roles

```go
names := ub.UserRoles(u)   // []string
```

`UserRoles` returns the effective role names for permission checks:

- the **initial user account** always resolves to `[Administrador]`;
- everyone else gets their assigned role names plus `Logged`.

### Managing roles

| Call | Purpose |
| --- | --- |
| `user.FindRoles(db, u)` | Load a user's `*role.Role` associations. |
| `(*Builder) AppendRoles(names...)` | Add role names to the set the builder ensures exist. |
| `(*Builder) GrantUserRole(db, userID, roleName)` | Attach a role to a user (writes `user_role_join`). |
| `(*Builder) InitDefaultRoles(db)` | Create any of `Builder.Roles` that don't exist yet. |

`Builder.Roles` (default `[Administrador]`) is the set seeded by
`InitDefaultRoles`; `Builder.UserManagerRoles` marks roles whose users are hidden
from non-admins in the listing.

## HTTP middlewares

Build them from the user builder:

```go
mw := ub.Middlewares(db, logoutURL, sessionMgr.CheckIsTokenValidFromRequest)
router.Use(mw.MiddlewareMD())   // or compose the pieces individually
```

`Middlewares` bundles four concerns:

| Middleware | Effect |
| --- | --- |
| `WithRoles` | Loads the current user's roles from `user_role_join` and sets them on the user, so downstream permission checks see them. |
| `Security` | Adds HSTS + no-cache headers. |
| `ValidateSessionToken` | Validates the request's session token (via the injected `checkIsTokenValidFromRequest`, typically `login_session.Manager.CheckIsTokenValidFromRequest`); redirects to the logout URL when invalid. Skipped during login-in-progress. |
| `Middleware` | Composes `ValidateSessionToken(WithRoles(Security(next)))`. |

`SetDevMode(true)` makes `MiddlewareMD()` a no-op wrapper (bypass auth in
development). Each middleware also has an `*MD()` variant returning the
`func(http.Handler) http.Handler` form for router `Use`.
