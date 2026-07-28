# When the session dies under an open page

Somebody has a form half filled in, in a dialog. The session expires (or someone
logged everybody out). They hit save.

Without any of this, the middleware answers `302` to the login page and the
browser replaces everything — the form goes, and so does everything typed into
it. They log in again and land on the home page, without their work.

What happens now: **the page does not move**. The login opens in a dialog over
it, and when it succeeds the interrupted request runs again, as if nothing had
happened.

## Who is asking

Every request made by `plaid()` carries

```
X-Plaid-Request: 1
```

(`web.PlaidRequestHeader`, set in `js/corejs/src/builder.ts`). That is what tells
the two cases apart:

| who asks | with no session |
|---|---|
| the browser navigating (`GET /admin/products`) | `302` to the login page — there is no page to keep |
| the page, through `plaid()` | `401` + `loginURI` — the page stays standing |

## The dialog, in two requests

1. The middleware finds the session gone on a `plaid()` request and answers
   **`401`** — which is what the server log should show anyway — with **one
   field**:

   ```json
   { "loginURI": "/admin/login-dialog" }
   ```

   (`web.EventResponse.LoginURI`.) Nothing else: no portals, no script. The
   client acts on nothing else in such a response.

2. `plaid()` asks for that address in a **new** request, carrying `vars` and,
   in its scope, the way back:

   ```js
   plaid()
     .vars(vars)
     .parseUrl(loginURI)
     .scope({ onLoginSuccess: () => { /* replays the original request */ } })
     .go()
   ```

   What comes back is the dialog, in the login portal
   (`presets.LoginPortalName`, at the layout root, which is why it covers
   everything), plus a script that:

   - puts the dialogs already open out of sight with CSS
     (`rvq-hidden-by-login`, `visibility: hidden` only — they stay mounted, with
     everything the user typed);
   - opens this one (`vars.presetsLoginDialog = true`);
   - waits for the `rvq:login-done` message, then restores those dialogs, closes
     the login and calls `onLoginSuccess()`.

3. `onLoginSuccess()` runs the original `plaid()` request again — same URL, same
   form. The user's click finally goes through.

The dialog holds the **real login page** in an `iframe`, so everything that page
does keeps working (form protection, reCAPTCHA when configured, OAuth buttons,
error messages). On success it lands on `presets.LoginDoneURI`, a tiny page whose
whole job is `postMessage("rvq:login-done")`.

`LoginDialogURI` and `LoginDoneURI` are whitelisted in the login middleware: they
are asked for precisely when there is no session.

## The two middlewares

Two different things turn a request away for want of a session, and both go
through the same hook (`login.UnauthorizedResponder`):

- **`x/login`** — there is no JWT, or it is no longer valid;
- **`admin/packages/user`, `ValidateSessionToken`** — the JWT is fine, but the
  session behind it ended (`login_sessions`, which IPCD's `make app-logout`
  expires).

The answer itself is installed by `admin/login` (`installLoginDialog`), because
only whoever owns the page can build a dialog on it.

> In dev (`-tags=dev`) the session middleware in `admin/packages/user` is a
> no-op: none of this shows up until the server runs without the tag.

## Where it lives

| what | where |
|---|---|
| header and response field | `web/page.go` (`PlaidRequestHeader`, `IsPlaidRequest`), `web/api.go` (`EventResponse.LoginURI`) |
| client side | `js/corejs/src/builder.ts` (`fetch`, `goPre`, `login`) |
| dialog, script and the 401 answer | `admin/login/dialog.go` |
| middleware hook | `x/login/builder.go` (`UnauthorizedResponder`, `RespondUnauthorized`, `WhiteListed`) |
| tests | `admin/login/dialog_test.go`, `js/integration_tests/admin/presets/login.test.ts` |
