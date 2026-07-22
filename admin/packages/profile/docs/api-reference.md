# API reference

## Types

| Symbol | Purpose |
| --- | --- |
| `Builder` | The profile feature; created by `New`. |
| `Options` | `CurrentUser func(*http.Request) user.User` and `NewUser func() user.User` hooks. |
| `Profile` | Read-only view model of the current user (id, name, account, status, roles). |
| `Messages` | i18n messages; `MessagesKey` is the module key. |

## Functions & methods

| Symbol | Purpose |
| --- | --- |
| `New(db, mgr *login_session.Manager, lb *login.Builder, userMb *presets.ModelBuilder) *Builder` | Construct the profile builder. |
| `(*Builder) Install(p *presets.Builder) error` | Register the `my-profile` singleton resource. |
| `(*Builder) Brand(p *presets.Builder) func(ctx) h.HTMLComponent` | Sidebar account widget. |
| `(*Builder) ProfilePage(p *presets.Builder) http.Handler` | Standalone profile page. |
| `(*Builder) NotificationComponentFunc()` / `SetNotificationComponentFunc(fn)` | Get/set the notification slot rendered in the brand widget. |
| `Model(p, v, opts...)` / `NewModel(p, v, opts...)` / `DefaultModelOptions(opts...)` | Register a model bound to this package's i18n module. |
| `GetMessages(ctx)` / `ConfigureMessages(*i18n.Builder)` | i18n access / registration. |

## Constants

| Symbol | Value |
| --- | --- |
| `ModelID` | `"my-profile"` — the singleton resource id. |
