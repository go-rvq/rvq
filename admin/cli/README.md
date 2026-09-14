# admin/cli — `http_api`

`HttpApiCommand` adds an `http_api` subcommand that runs HTTP requests **against
the application handler in the same process**, without going over the network
and without the secure key. It is meant for a trusted operator on the server
host (migrations, seeding, scripted admin actions, smoke tests).

## Wiring

The command is application-agnostic; the app supplies two functions:

```go
mainCmd.Sub(admincli.HttpApiCommand(admincli.Config{
    // builds the app HTTP handler on the given mux; the request is served against it
    Handler:  func(mux *http.ServeMux) http.Handler { return app.Handler(mux) },
    // resolves an account name to the app user object (e.g. login.Builder)
    FindUser: loginBuilder.FindUserByAccount,
}))
```

- The request is tagged with `web.RequestSourceCLI` (context only, not a header,
  so it cannot be forged over HTTP); read it with `web.IsCLIRequest(r)`.
- With a login, the request runs as that user via `login.WithTrustedUser`, which
  the login middleware honours in place of a session cookie.

## Invocation

Simple form — a single positional URI (body from `--data` or piped stdin):

```sh
app http_api --http-method POST --user admin \
  --data '{"Title":"Hi"}' '/admin/posts?__execute_event__=presets_Update'
```

JSON form — spec(s) from stdin or file(s). A JSON **array** runs many requests,
a JSON **object** runs one; several files are several requests:

```sh
echo '{"login":"admin","method":"POST","uri":"/admin/...","body":{...}}' | app http_api
app http_api requests.json
app http_api a.json b.json
```

Flags (`--http-method`, `--user`, `--content-type`, `--data`, `--raw`) are
optional; in the JSON form they only provide defaults for omitted spec fields.

## Spec format

Each request is a JSON object. Every field is optional except `uri`; an omitted
field falls back to the matching flag.

| Field | Meaning |
|-------|---------|
| `login` | user login to run the request as (alias: `user`) |
| `method` | HTTP method (default `GET`) |
| `uri` | request URI (**required**) |
| `contentType` | request `Content-Type` (default `application/json`) |
| `body` | request body — see [Body flattening](#body-flattening) |
| `expectedResponse` | assertions on the result — see [Validation](#validation) |

### Body flattening

The admin event handlers read **form values**, not JSON, so a JSON `body` is
flattened into dot/bracket form fields and sent as `multipart/form-data`:

| JSON | Form field(s) |
|------|---------------|
| `{"Config":{"TypeID":2}}` | `Config.TypeID=2` |
| `{"Tags":["a","b"]}` | `Tags[0]=a`, `Tags[1]=b` |
| `{"file:Cover":"/img.png"}` | file field `Cover`, read from disk |

A key of the form `"file:KEY"` attaches the file at its path value as the
multipart **file** field `KEY`. The file is read straight from disk and **left
in place** (it is not a temporary upload to be removed after the request).

### Validation

`expectedResponse` asserts the result. Every field present must hold; the first
failure stops the run and exits non-zero **before the next request**.

| Field | Check |
|-------|-------|
| `status` | exact HTTP status code |
| `body.equal` | response body equals this string |
| `body.contains` | response body contains this string |
| `body.starts` | response body starts with this string |
| `body.ends` | response body ends with this string |
| `keys` | list of JSON dot-paths that must exist (e.g. `response.updatePortals`) |

## Output

One result object per request; a JSON array input yields a JSON array of
results:

```json
{ "status": 200, "flash": ["Saved successfully"], "response": { ... } }
```

The flash portal (`FlashPortalName`) is handled specially: its snackbar text is
stripped of HTML into `flash`, and that portal is **removed** from `response`.
Use `--raw` to print the untouched response body instead.

## Example spec

```json
{
  "login": "admin",
  "method": "POST",
  "uri": "/admin/posts?__execute_event__=presets_Update",
  "body": { "Title": "Oi", "file:Cover": "/tmp/capa.png" },
  "expectedResponse": {
    "status": 200,
    "body": { "contains": "salvo" },
    "keys": ["flash", "response.updatePortals"]
  }
}
```
