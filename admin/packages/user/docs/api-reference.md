# API reference

## Types

| Symbol | Purpose |
| --- | --- |
| `User` | Interface a user model must implement (`login.UserPasser` + uuid id, name/email/status/account/roles/registration date). |
| `Builder` | Configures the user resource; created by `New`. Fields: `ExpireAllSessionLogs`, `LoginInitialUserEmail`, `Roles`, `UserManagerRoles`. |
| `Middlewares` | Auth/role/security HTTP middlewares; created by `(*Builder).Middlewares`. |
| `ChangePassword` | Form for the change-password action (`OldPassword`, `NewPassword`, `ConfirmPassword`). |
| `Messages` | i18n messages; `MessagesKey` is the module key. |

## Construction & model registration

| Symbol | Purpose |
| --- | --- |
| `New(db, lb, mb, loginInitialUserEmail) *Builder` | Wire the user resource (listing/editing/detailing + actions). |
| `Model(p, v, opts...) *presets.ModelBuilder` | Register a model bound to this package's i18n module. |
| `NewModel(p, v, opts...) *presets.ModelBuilder` | Same, via `presets.NewModelBuilder`. |
| `DefaultModelOptions(opts...)` | Model options that set the module key. |

## Users & roles

| Symbol | Purpose |
| --- | --- |
| `GetCurrentUser(r) User` | The authenticated user of a request (nil if none). |
| `(*Builder) GenInitialUser(db) User` | Create the initial admin from the login builder's initial account/password. |
| `(*Builder) UserRoles(u) []string` | Effective role names (initial user ⇒ `[Administrador]`, else roles + `Logged`). |
| `FindRoles(db, u) []*role.Role` | Load a user's role associations. |
| `(*Builder) GrantUserRole(db, userID, roleName)` | Attach a role to a user. |
| `(*Builder) InitDefaultRoles(db)` | Create missing default roles. |
| `(*Builder) AppendRoles(names...) *Builder` | Extend the default role set. |
| `RoleAdministrador` / `RoleLogged` | Built-in role names. |

## Middlewares

| Symbol | Purpose |
| --- | --- |
| `(*Builder) Middlewares(db, logoutURL, checkIsTokenValidFromRequest) *Middlewares` | Build the middleware set. |
| `WithRoles` / `WithRolesMD` | Load the current user's roles onto the user. |
| `Security` / `SecurityMD` | HSTS + no-cache headers. |
| `ValidateSessionToken` / `ValidateSessionTokenMD` | Validate the session token, redirect to logout when invalid. |
| `Middleware` / `MiddlewareMD` | Composition of the above (`MiddlewareMD` is a no-op in dev mode). |
| `SetDevMode(bool)` / `DevMode()` | Toggle/read dev mode. |

## Actions & i18n

| Symbol | Purpose |
| --- | --- |
| `ConfigureChangePasswordAction(lb, currentPasswordCheck, d, getUser)` | Add the change-password detail action. |
| `GetMessages(ctx) *Messages` / `ConfigureMessages(*i18n.Builder)` | i18n access / registration. |

The detailing configured by `New` also registers the event funcs
`eventUnlockUser`, `eventSendResetPasswordEmail`, and (when a session expirer is
set) `eventRevokeTOTP`, surfaced as the unlock / send-reset-password /
revoke-TOTP actions.
