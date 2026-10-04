# My keys

An **access key** is a code an automation — **git**, a **WebDAV** program,
a script — signs in with in your name, with no password and no session: each
request carries the key.

## Making a key

1. In **My keys**, make a key: give it a **name** (what it is for), how long
   it **lasts** (90 days when left empty; 1 year at most) and its
   **permissions**.
2. When saved, its **code** is shown **only once**: `hck_…`. Copy it and
   keep it safe. It is not kept anywhere — only a digest (hash) of it —; lost,
   make another key.
3. To stop a key at once: untick **Enabled**, or delete the key.

## The permissions

A key's permissions **only restrict**: a request by the key is allowed when
**you may** and **the key allows**. A key never gives more than you have —
and, with no permission, it does nothing.

They are written as a role's are (see
{%= admin.doc("guides/policies/01-permissions").link %}), but only allow:

| Resources | The key may |
| --- | --- |
| `admin:*` | all you may |
| `admin:/site-files:!git`, `admin:/site-files:@get` and the files' actions | use the site's files by git |
| `admin:/site-files:<static/*>:!edit` | only change the files of `static/` |

A key never makes nor changes keys — not even yours.

## Using it

- **git:** the password is the code (the user may be anything):

  ```
  git clone https://<the-site>/admin/site-files-draft.git
  # user: your login; password: hck_…
  ```

- **WebDAV:** the same — the user and, as the password, the code.
- **A script:** the header `Authorization: Bearer <code>`:

  ```
  curl -H "Authorization: Bearer hck_…" https://<the-site>/admin/…
  ```

## The history

Each request by the key is in its **History**: when, from where (address
and place), by which program, what (git, webdav, api; the path) and how it
went. It is also in your **access log** (the sessions of your profile), as
"Access key: *name*" — one line each half hour by address and program, not
one a request.
