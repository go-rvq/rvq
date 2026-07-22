# Getting started

## The `User` interface

Your user model must implement `user.User`:

```go
type User interface {
    login.UserPasser                 // account + password (x/login)
    GetID() uuid.UUID                // uuid primary key
    GetName() string
    SetName(string)
    SetEmail(string)
    SetRegistrationDate(time.Time)
    GetStatus() string               // "active" / "inactive"
    GetAccountName() string          // login account (email)
    GetRoles() role.Roles
    SetRoles(role.Roles)
}
```

A typical model embeds `login.UserPass` (and optionally `login.SessionSecure`
for TOTP) and a `uuid.UUID` primary key, with a many2many `Roles` association.

## Registering the model

Use the package helpers so the model is bound to this package's i18n module:

```go
mb := user.Model(presetsBuilder, &models.User{})     // or user.NewModel(...)
```

`Model` / `NewModel` apply `DefaultModelOptions`, which set the module key
(`MessagesKey`) used for translations.

## Wiring the builder

```go
ub := user.New(db, lb, mb, loginInitialUserEmail)
```

`New`:

- registers the i18n messages;
- configures the listing / editing / detailing (see the README);
- returns a `*Builder` you can further tune with `AppendRoles`, `Roles`,
  `UserManagerRoles`, and `ExpireAllSessionLogs` (needed for the "revoke TOTP"
  action, which also expires the user's session logs).

## The initial user

```go
ub.GenInitialUser(db)
```

Reads the initial account/password from the `login.Builder`
(`GetInitialUserAccount` / `GetInitialPassword`); if that account does not exist
yet it is created, encrypted, and granted the `Administrador` role.
`GenInitialUser` also ensures the default roles exist (`InitDefaultRoles`).

The initial user account is treated as protected throughout the admin: it cannot
be edited, have its roles changed, or be deleted.

## Current user

```go
u := user.GetCurrentUser(r)   // *http.Request -> User (nil if not logged in)
```

Resolves the authenticated user from the login session.
