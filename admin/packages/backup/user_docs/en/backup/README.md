# Backup

A copy of the **whole application** — the database, the site's files, the
uploads, the program itself, its configuration —, each one a **state**. A
client downloads the state into a folder, and the application starts there.

A state is made on the server itself, the application running: the files
are linked (hard links), not copied — a state takes on the disk only what
changed since —, and a **manifest** lists each file with its sha256. The
client compares the manifest with what it already has and downloads only
what changed.

## A new state

**New** asks:

- **Description**: why this state ("before the 2.0 migration").
- **Sources**: what goes in. Without one, the state is **partial**: it
  recovers a part, but does not start the application by itself.
- **Keep**: the retention never removes it.

It is **queued** and **built** in the background; the list and the detail
show the progress of each source. It ends **ready**, or **failed** (the
error in the detail), or **canceled**.

A state is also made by itself, every day (origin **scheduled**). Of the
ready states, the 7 newest stay (and as many of the failed or canceled
ones); those marked **Keep** stay always.

## The detail

- **Status**, **Origin**, **Created by**, **Started** and **Finished**,
  **Files** and **Size**, **Manifest (sha256)**.
- **By source**: the files, the size, the time and the warnings of each
  source.
- **Download**: once ready, the link of the manifest and the **client
  routes**, with the id of the state.
- **Downloads**: the download sessions opened for it.

A state being built cannot be deleted: cancel it first. A state whose files
are gone (its folder removed on the server) is marked failed when the
application restarts.

## Downloading a state

A program downloads a state with an **access key** (your profile › **Access
keys**) that does not deny the permission **Download** of this page. The key
opens a **session**, which gives a **token**; the token downloads the
manifest and the files — and renews the session — for an hour.

```
ADMIN=https://your-site/admin
KEY=hck_…                                   # the access key

# the newest ready, complete state
curl -fsS -H "Authorization: Bearer $KEY" "$ADMIN/!backup/latest"
# → {"id":"…","created":"…","files":…,"bytes":…}

# a session: the token
curl -fsS -X POST -H "Authorization: Bearer $KEY" "$ADMIN/!backup/<id>/session"
# → {"token":"…","expires":"…","state":"<id>","manifest":"…","files":"…"}

# the manifest and each file, by the token
curl -fsS -H "Authorization: Bearer $TOKEN" "$ADMIN/!backup/<id>/manifest" -o manifest.gz
curl -fsS -H "Authorization: Bearer $TOKEN" "$ADMIN/!backup/<id>/files/<path>" -o <path>
```

A file takes **Range** (an interrupted download goes on where it stopped)
and **If-Match** with its sha256 (the ETag): if the file changed, the answer
is 412 instead of a mixed file.

The manifest is a compressed text (gzip), a line per entry, apart by tabs:
the sha256 (base64url, no `=`), the mode, the time, the size, the kind (`f`
file, `d` folder, `l` link) and the path.

## Permissions

- **View** the page: the list, the detail, the manifest through the admin.
- **Create**, **Edit** (the description, **Keep**), **Delete**, **Cancel**.
- **Download** (`!download`): the manifest and the files of a state — all
  the data and secrets of the application. Grant it with care.

## Stable paths

The client uses the service `!backup`, out of the pages of the menu and of
the admin's login (see the guide **Stable paths**). The addresses do not
change with the menu. `<admin>` is the address of the admin.

| Method | Path | Access | What it gives |
| --- | --- | --- | --- |
| GET | `<admin>/!backup/latest` | access key | the id of the newest ready, complete state |
| POST | `<admin>/!backup/{id}/session` | access key, or the token (renews) | a session: its token, the addresses of the manifest and the files |
| GET | `<admin>/!backup/{id}/manifest` | token, or the admin's session | the manifest |
| GET | `<admin>/!backup/{id}/files/{path}` | token | a file of the manifest (Range, If-Match) |

With no `Authorization` header a request goes to the admin: by the admin's
session, only the manifest (the link of the detail). Errors are JSON
(`{"error": "…"}`): 401 with no key or a refused one, 403 when the key may
not download, 404 when the state does not exist, 409 when it is not ready
yet.
