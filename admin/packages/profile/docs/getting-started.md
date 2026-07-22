# Getting started

## Dependencies

`profile` composes two sibling packages:

- [`user`](../../user) — resolves the current user (`user.GetCurrentUser`) and
  defines the `user.User` interface + roles.
- [`login_session`](../../login_session) — the `*login_session.Manager` used to
  list and revoke the user's sessions.

## Construction

```go
pf := profile.New(db, sessionMgr, lb, userMb)
```

- `db` — the database.
- `sessionMgr` — a `*login_session.Manager`.
- `lb` — the `*login.Builder` (login/logout URLs).
- `userMb` — the user `*presets.ModelBuilder` (used for its i18n builder).

`New` registers the profile i18n messages and returns a `*Builder`.

## Install

```go
if err := pf.Install(presetsBuilder); err != nil { ... }
```

`Install` registers a **singleton** `Profile` model with id `my-profile`
(`ModelID`), hidden from the menu. It:

- allows the `Logged` role to access `*:my_profile:*`;
- fetches the current user on open and maps it to the read-only `Profile` view
  (name is editable, account/status/roles are read-only);
- forbids deletion.

## Sidebar brand widget

```go
presetsBuilder.BrandFunc(pf.Brand(presetsBuilder))
```

Renders the account card (avatar, name, first role, logout). Attach a
notification component with:

```go
pf.SetNotificationComponentFunc(func(ctx *web.EventContext, u user.User) h.HTMLComponent {
    return myBell(ctx, u)
})
```

## Standalone profile page

```go
mux.Handle("/profile", pf.ProfilePage(presetsBuilder))
```

Serves a page with "Login sessions" (opens the sessions dialog) and "Sign out".

## The `Profile` view model

```go
type Profile struct {
    ID          uuid.UUID `admin:"ro" gorm:"type:uuid"`
    Name        string
    AccountName string `admin:"ro"`
    Status      string `admin:"ro"`
    Roles       string `admin:"ro"`
}
```

A read projection of the current user (not a separate table row you manage).
