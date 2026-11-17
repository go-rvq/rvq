# File system

The files of the server the site keeps: browse the folders, upload files and
folders (drag them onto the list), create, edit, rename, copy and delete;
download a file, or a folder as zip; extract an archive (zip, tar, gz).

The address follows the folder open (`?dir=/folder`): a link opens the page
there.

## WebDAV

The same files may be reached from a computer through **WebDAV**, with the
address and the account shown on the page — as a network drive: the last
icon of the explorer's bar opens the **WebDAV protocol** panel, at the right
(expandable; its width dragged by its edge). There are the address, the
account (your e-mail of the admin; the password, better an access key of
your profile, which can be revoked), the permissions the key needs, the
folders — and which are read-only — and programs to use.

### The access key

A key's permissions only **deny**, on top of what you may: a key with no line
reaches what you reach. For the WebDAV, its lines must not deny this page's
resource (the panel shows its name) its `@get` (read and download) nor, to
write, its `!write` (upload, change, move and delete). A key that only reads:

```
deny admin:/fs-tools:!write
```

## Terminal

The terminal button of the bar (or **F4**) opens a shell below the list, in
the open folder, and follows it when the folder changes. The bar between the
list and the terminal is dragged to change its height (a double click: the
default height; ↑ ↓ by the keyboard), kept in the browser; the terminal's
expand button puts it over the whole window, and gives it back.

Beside the state of the connection, a chip says how the shell is kept apart
(**Isolated · Landlock**, or **Not isolated**); a click shows the details of
the session — the isolation, the folders and whether each is written, the
network (on, or off: to use it, ask the server's administrator), the
limits, when it closes unused.

It runs apart from the rest of the server:

- it sees only the folders of the file system, in `$FILES_ROOT`
  (`cd "$FILES_ROOT"/media`), each read or written as you may;
- its `HOME` and `/tmp` are a folder of the session, removed when it closes;
- it sees neither the server's configuration and secrets nor the network
  (unless the configuration allows it);
- it has limits of CPU time, memory and file size;
- it closes by itself after 30 minutes unused, and when the page closes.

## Permissions

This page's permissions hold for the explorer, the WebDAV and the terminal:

| Permission | What it allows |
| --- | --- |
| `@get` | see and download the files |
| `!write` | upload, create, change, move and delete |
| `!terminal` | open the terminal — in no set of permissions: give it only to who needs it |

Beside the folders of the server are folders of each user: the draft of the
**Site files** (`site-files`), with the permissions of editing it in the
admin. A folder configured read-only is never written.

## For the server's administrator

The folders come from `RVQ_FSTOOLS_DIRS`, in the `env:` of the `app.yml`:
`;` apart, each `MODE,POINT,FOLDER` — `r` read-only, `w` read and write:

```yaml
env:
  RVQ_FSTOOLS_DIRS: "w,/media,data/media;r,/logs,/var/log/app"
```

The terminal is a feature, `rvq.fstools.terminal`, **on by default**. To turn
it off, in the `app.yml` (and restart):

```yaml
features:
  rvq.fstools.terminal: false
```

A program built with the tag `rvq.fstools.noterminal` has none.

Its isolation and limits, in `fstools.terminal` of the `app.yml` (every key
optional; here, the defaults):

```yaml
fstools:
  terminal:
    sandbox: auto    # auto, bwrap, landlock or off (no isolation)
    network: false   # the shell reaches the network (below)
    shell: ""        # default: /bin/bash, else /bin/sh
    limits:
      cpu: 10m
      memory: 2G
      file_size: 1G
      processes: 0   # counts every process of the server's user
    idle: 30m
```

`auto` uses bubblewrap when it runs on the server, else the Linux kernel's
Landlock; with neither, the terminal is off and the WebDAV panel says why.
With Landlock before Linux 6.12, the shell may still signal the processes of
the server's user: run the server as a user of its own, or turn the terminal
off where that matters.

### The terminal's network

`network: false` (the default): the shell reaches no other machine — with
Landlock, no TCP connection and no port opened (Linux 6.7 or newer; before
it Landlock does not block TCP); with bubblewrap, a network of its own,
empty. The chip shows the network off and tells the user to ask the
administrator.

`network: true`: the shell reaches what the server reaches — the database
for a `pg_dump`, say. The server's passwords stay out of the shell's
environment: the tool needs its own (a `.pgpass` in a folder of the user,
or typed).

```yaml
fstools:
  terminal:
    network: true
```
