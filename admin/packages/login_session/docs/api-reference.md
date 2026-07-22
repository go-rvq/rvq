# API reference

## Types

| Symbol | Purpose |
| --- | --- |
| `Manager` | Records / validates / expires login sessions; created by `NewManager`. |
| `LoginSession` | The session row (`UserID uuid.UUID`, `Device`, `IP`, `TokenHash`, `ExpiredAt`, computed `Time`/`Status`). |
| `Messages` | i18n messages; `MessagesKey` is the module key. |

## Manager

| Method | Purpose |
| --- | --- |
| `NewManager(lb *login.Builder) *Manager` | Construct; registers i18n messages. |
| `AddSessionLogByUserID(db, r, userID) error` | Record a new session for the request's token. |
| `UpdateCurrentSessionLog(db, r, userID, oldToken) error` | Move the log from an old token to the current one (session extension). |
| `ExpireCurrentSessionLog(db, r, userID) error` | Expire the session matching the request's token. |
| `ExpireAllSessionLogs(db, userID) error` | Expire every session of the user. |
| `ExpireOtherSessionLogs(db, r, userID) error` | Expire all sessions except the request's current one. |
| `CheckIsTokenValidFromRequest(db, r, userID) (bool, error)` | Whether the request's token matches a live session (for auth middleware). |
| `Sessions(db, ctx, userID) (h.HTMLComponent, error)` | Render the user's sessions list. |

## Helpers

| Function | Purpose |
| --- | --- |
| `GetIP(r) string` | Client IP (X-Forwarded-For aware). |
| `GetProxy(r) []string` | The X-Forwarded-For chain. |
| `GetStringHash(v string, len int) string` | First `len` hex chars of `sha256(v)`. |
| `IsTokenValid(s LoginSession) bool` | True when `now` is past the session's `ExpiredAt`. |

## Constants

| Symbol | Value |
| --- | --- |
| `LoginTokenHashLen` | `8` — stored token-hash length. |
