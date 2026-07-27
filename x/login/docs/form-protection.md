# Protecting the login form without reCAPTCHA

The login and forget-password forms accept a POST from anyone. With
`Builder.Recaptcha(true, …)` configured, Google answers for them. Without it
they used to accept anything — so an installation with no Google account had a
form a script could hammer all night.

Now, when there is no reCAPTCHA, a **built-in protection** takes over
([`challenge.go`](../challenge.go)). It is on by default and needs no
configuration: it signs with the same secret as the session (`Builder.Secret`).

## Why not a proof of work, or a modern captcha

Everything the browser side of those needs — `crypto.subtle`,
`navigator.credentials`, `PublicKeyCredential` — exists only in a **secure
context**: HTTPS, or `localhost`. An admin served over plain HTTP inside a
network (which these are, often) would be left with no protection at all, or
with a broken login page.

So this one is decided entirely by the SERVER, in plain HTML, with **no
JavaScript**: it behaves the same over HTTP and over HTTPS.

## What it checks

| | |
| --- | --- |
| **Signed challenge** | The form carries `__lc`: a random nonce, the instant it was issued, and the HMAC-SHA256 of both. A POST that did not come from a form this server rendered has none, and cannot make one up. |
| **Single use** | A challenge already spent is refused, so one captured from a rendered page cannot be replayed. |
| **Time to fill** | Under `ChallengeMinAge` (900 ms) nobody read and typed the form; over `ChallengeMaxAge` (30 min) it is a page left open. |
| **Honeypot** | `email_confirmation`, a field taken out of the reader's way with CSS (and out of a screen reader's, with `aria-hidden`/`tabindex`). Anything that fills it is not a reader. |

Each is cheap, and none asks the person logging in to do anything.

## Configuring it

```go
lb := login.New(i18nB).Secret(secret)   // the challenge signs with this secret

lb.ChallengeAges(2*time.Second, time.Hour) // stricter/looser than the defaults
lb.ChallengeDisabled(true)                 // off (a proxy already filters, say)
```

It turns itself off when reCAPTCHA is configured — the two do the same job.
`Builder.ChallengeEnabled()` reports which one is on.

Rendering it is one call, already wired into both forms of both view sets
(`x/login/views.go` and `admin/login/views.go`):

```go
vh.ChallengeFields()   // nil when the protection is off
```

A custom login view MUST include it, or its POSTs are refused.

## What the user sees when it refuses

`FailCodeIncorrectChallenge` — *"Could not confirm this form was filled in by a
person. Please try again."* — and `FailCodeChallengeExpired` — *"This page has
been open for too long. Please try again."* Both redirect back to the form,
which comes back with a fresh challenge, so a person who hit either one just
submits again.

## The secure key: getting past it on purpose

Automation needs a way in that the form protection does not stand in front of: a
health check, a deploy script, an operator with a terminal. That way is a FILE
on the server, [`secure_key.go`](../secure_key.go):

```
.secure_key            mode 0600, next to the program (SecureKeyFile to move it)
X-Secure-Key: <value>  the header a request sends
```

A request whose header matches the file skips the form protection — reCAPTCHA
included. Whoever cannot read the file has nothing to send.

**It skips the protection, not the login.** The account and the password are
still required, and so is everything else: the retry count, the permissions.

```go
lb.SecureKeyFile("/etc/myapp/.secure_key")   // where it lives
key, err := lb.RenewSecureKey()              // write a fresh one now
err = lb.StartSecureKeyRenewal()             // write one and keep renewing it
err = lb.StartSecureKeyRenewal(appCron)      // ... on the application's cron
lb.StopSecureKeyRenewal()
```

| | |
| --- | --- |
| `RVQ_SECURE_KEY_RENEW` | how often the key is rewritten, as a cron rule — `"@daily"` by default, or `"0 */6 * * *"`, `"@every 12h"`, … (robfig/cron, the same one the worker uses). Empty turns the renewal off and leaves whatever key is on disk. |
| mode `0600` | the file is written with it, and a file with anything wider is **ignored**, with a warning: a key the whole machine can read is not a key. |
| no file | no bypass — the protection answers for every request. |

The value is 32 random bytes, base64. Renewal writes to a temporary file and
moves it into place, so a reader never sees half a key, and the old value stops
working the moment the new one lands.

## What it does NOT do

- **It is not a rate limiter.** A script that fetches the login page before each
  attempt gets a valid challenge every time. What limits attempts against one
  account is `Builder.MaxRetryCount` (the user is locked after N failures);
  limiting per IP is the proxy's job.
- **Several instances behind a balancer** each keep their own spent-challenge
  list in memory, so a captured challenge could be replayed once per instance.
  Every other check still stands.
- **It is not a captcha.** It stops automated submission of THIS form, not a
  human being paid to type.

## Tests

[`challenge_test.go`](../challenge_test.go): a rendered form passes; a POST from
nowhere, one signed with another secret, a replayed one, one submitted too fast
or too late, and a filled honeypot are refused; the protection turns itself off
under reCAPTCHA and when disabled; and the rendered fields carry no `<script>` —
the property that makes it work over plain HTTP.

[`secure_key_test.go`](../secure_key_test.go): the key holder goes past the
protection (and past reCAPTCHA); a wrong key, no key, an empty file and a
missing one do not; the file is written `0600` and one that others can read is
ignored; renewing writes a fresh value and retires the previous one; and the
renewal runs on schedule, on its own cron or on the application's, refusing a
rule that makes no sense.
