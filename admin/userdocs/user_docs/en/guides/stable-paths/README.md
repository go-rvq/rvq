# Stable paths

Every address of the admin — what follows its prefix, `/admin` — is one of
these:

1. **A page of the menu**: a model (`/admin/<model>`, `/admin/<model>/<id>`), a
   group of the menu or a page. A model **hidden** from the menu is still a
   model: its address works, it is just not shown in the menu. These
   addresses follow the menu: when a model moves to another group, its
   address may change.
2. **A service**: an address no page has — of a program, a download, an
   API, a protocol, the assets of the pages. It has one of two forms:
   - **`/admin/!<name>/…`** — the service's name after a **`!`**:
     `/admin/!backup/latest`, `/admin/!webdav/`. No page begins with `!`: a
     service is never mistaken for a page, and its address **does not
     change** with the menu.
   - **an address whose form already tells its protocol** —
     `/admin/site-files.git` (git over HTTP: the `.git` at the end). No page
     has such an address, and the program on the other side (`git`) needs it
     as it is.

The addresses of the services are **stable**: a script, a `git remote`, a
WebDAV client, a link in a mail can keep them. An old address is not
redirected: when a service changes its address, the release notes say the
new one.

## How each is reached

- **public**: anyone (the static assets of the pages).
- **session**: whoever is signed in to the admin, with their account's
  permissions.
- **signed link**: only who was given the link, for a while.
- **access key**: a program, with an access key in the header
  `Authorization: Bearer <key>` — no admin session.
- **token**: a token the service itself gave (for instance, the session of a
  backup download).
- **HTTP Basic**: a login and a password — or an access key in place of the
  password — asked by the program itself (`git`, a WebDAV client).

A service with **its own login** is reached before the admin's login: a
client uses it without going through the login screen.

## The services of this admin

{%= admin.services() %}
