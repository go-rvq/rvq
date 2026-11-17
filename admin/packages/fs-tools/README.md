# FS Tools

The file system of the admin (**Administration › File system**): folders of
the server, reached by the page's explorer, by **WebDAV** and by a
**terminal** in a sandbox.

- [The folders](#the-folders)
- [Permissions](#permissions)
- [The page](#the-page)
- [WebDAV](#webdav)
- [The terminal](#the-terminal) — a feature, configured in the app.yml
- [The sandbox](#the-sandbox)
- [Reusing the explorer](#reusing-the-explorer)
- [Tests](#tests)

<style>
figure {
  border: thin silver solid;
  display: flex;
  flex-flow: column;
  padding: 5px;
}

figcaption {
  background-color: #222222;
  color: white;
  font: italic smaller sans-serif;
  padding: 3px;
  text-align: center;
}
</style>

<figure>
    <img src="./docs/page.png" alt="File system page" />
    <figcaption>File system page</figcaption>
</figure>

## Installing

```go
fsb, err := fs_tools.New(presets, i18nBuilder, loginBuilder)
// a directory of each request (optional): see Mounts
fsb.AddMount(&fs_tools.Mount{Name: "site-files", Dir: …, Allowed: …})
// told each access, for a log of sessions (optional)
fsb.OnAccess(func(r *http.Request, kind string) { … }) // webdav, git, files, terminal
err = fsb.Install(presets)
```

`Install` adds the page, its explorer's API (`<prefix>/!fs-tools/api/fs/…`),
its terminal (`<prefix>/!fs-tools/terminal`), the components (`UseAssets`) and
the WebDAV (`<prefix>/!webdav`).

These are services (`presets.Service`, `admin/PATHS.md`): `fs-tools`,
`webdav` and `git` — whose repositories are at git's protocol path,
`<prefix>/<name>.git/`, mounted on it. The user documentation lists them,
in its section **Stable paths**.

## The folders

**`RVQ_FSTOOLS_DIRS`** (the `env:` of the app.yml): folders of the server,
`;` apart, each `MODE,MOUNT_POINT,LOCAL_PATH` — `MODE` `r` (only read) or `w`
(read and written, the default); a relative `LOCAL_PATH` is of the working
directory.

```yaml
env:
  RVQ_FSTOOLS_DIRS: "w,/media,data/media;r,/logs,/var/log/app"
```

Unset, the root is empty — no error.

**Mounts** (`AddMount`): a directory of each request, beside them at the root
(the draft of a user's site files, say): `Dir(r)` is its folder, `Allowed(r,
MountAccess)` what the request may do in it, `Hidden(p)` the paths never
served (`.git`).

## Permissions

The page's permissions are those of everything — explorer, WebDAV, terminal:

| Permission | What |
| --- | --- |
| `@get` | see and download the files |
| `!write` | upload, create, change, move, delete |
| `!terminal` | open the terminal — in no set of permissions: give it by hand |

The page's sets of permissions: `readers` (`@get`), `authors` (`@get`,
`!write`). A mount answers for itself (`Allowed`); a folder `r,` of
`RVQ_FSTOOLS_DIRS` is never written, whatever the permissions.

## The page

The explorer (`<vx-files>`) over the WebDAV's files: list, tree and icons;
upload by dragging files and folders (large files in parts); edit a
text, rename, copy, delete; download a file, or a folder or several items as tar.gz; extract an
archive (zip, tar, gz). `?dir=/media/x` opens that folder; the address
follows the folder open.

The last icon of the explorer's bar opens **WebDAV protocol**, a drawer at
the right (expandable, its width dragged): the address, the account, what an
access key needs, the folders and whether each is only read, the programs to
use. When the terminal is off, the reason is there too, to who has
`!terminal`.

## WebDAV

`<prefix>/!webdav`, by HTTP Basic: the account of the admin, the password its
password or — better — an access key (Profile › Access keys: revocable). The
same folders and permissions as the page.

A key's permissions only deny, on top of what its user may (rvq
`accesskeys`): a key with no line reaches what its user does. For the
WebDAV its lines must not deny the page's resource (`admin:/fs-tools:`, the
drawer shows it as this installation names it) its `@get` (reading) nor —
to write — its `!write`. A key that only reads:

```
deny admin:/fs-tools:!write
```

## The terminal

**F4** (or the button of the bar) opens a shell below the explorer, in the
folder open; it follows the folder when it changes. The bar between the
files and the terminal is dragged to change its height (a double click: the
default height; ↑ ↓ when it has the focus), kept in the browser; the
terminal's expand button puts it over the whole window, and gives it back. Whoever opens it needs
`!terminal`, and the page needs a folder of the system (`RVQ_FSTOOLS_DIRS`
or a mount).

### The feature `rvq.fstools.terminal`

The terminal is a feature (`x/features`): **compiled and on by default**.

| To | Do |
| --- | --- |
| turn it off in an installation | `features: {rvq.fstools.terminal: false}` in the app.yml, and restart |
| leave it out of the binary | build with the tag `rvq.fstools.noterminal` (`go build -tags rvq.fstools.noterminal`) |
| see whether it is on | `app features` (on), `app features all` (compiled) — HermonCMS: Server operations' documentation |

Off or not compiled, the button is not shown and `<prefix>/!fs-tools/terminal`
answers 403. In code: `fs_tools.TerminalFeature.Enabled()`.

### Its configuration: `fstools.terminal` of the app.yml

```yaml
fstools:
  terminal:
    sandbox: auto    # auto | bwrap | landlock | off (The sandbox)
    network: false   # the shell reaches the network (below)
    shell: ""        # its program: default /bin/bash, else /bin/sh
    limits:          # 0 or "": none
      cpu: 10m       # time of CPU (a Go duration)
      memory: 2G     # virtual memory (K, M, G, T of 1024; bytes alone)
      file_size: 1G  # size of a file written
      processes: 0   # RLIMIT_NPROC: counts every process of the server's
                     # user, the server's threads too — high, or 0
    idle: 30m        # with no input, closed
```

Every key optional: the values above are the defaults. Read at the start (a
change needs a restart); a wrong value turns the terminal off, the reason in
the log and in the WebDAV drawer.

| Key | Default | What |
| --- | --- | --- |
| `sandbox` | `auto` | how the shell is kept apart: [The sandbox](#the-sandbox) |
| `network` | `false` | the shell reaches the network — below |
| `shell` | `""` | its program: `/bin/bash`, else `/bin/sh` |
| `limits.cpu` | `10m` | CPU time of the session's processes (a Go duration: `90s`, `10m`, `1h`); `0`: none |
| `limits.memory` | `2G` | virtual memory of each process (`512M`, `2G`; K, M, G, T of 1024; bytes alone); `0`: none |
| `limits.file_size` | `1G` | the largest file written; `0`: none |
| `limits.processes` | `0` | RLIMIT_NPROC — it counts every process and thread of the server's user, the server's too: high, or `0` (none) |
| `idle` | `30m` | with no input for so long, the session closes |

#### `network`

Off (`false`, the default), the shell reaches no other machine:

- `landlock`: no TCP connection and no port opened (Linux ≥ 6.7, Landlock
  ABI 4; below it Landlock cannot block TCP — use `bwrap`, or keep the
  terminal off). UDP is not covered (DNS answers, but nothing reaches a
  service over TCP). A Unix socket reached by its path is not covered by
  ABI 4: a database's local socket in a folder the shell sees answers.
- `bwrap`: a network of its own, with nothing in it.

On (`true`), the shell reaches what the server reaches: every TCP host and
port (`landlock`), the server's network (`bwrap`). Its environment still
has none of the server's secrets — a `pg_dump` to the database needs its
own credentials (a `.pgpass` in a folder of the user, typed).

The chip of the terminal (beside its state) says whether the network is on;
off, it tells the user to ask the server's administrator.

### What the shell sees

- `$FILES_ROOT`: a folder of links, one per folder the user reaches at its
  place (`$FILES_ROOT/media`, `$FILES_ROOT/site-files`), each read or written
  as the user may (`!write`, the mount's `Allowed`, `r,`);
- `HOME` and `TMPDIR`: a temporary folder of the session, removed at its end;
- `PATH=/usr/local/bin:/usr/bin:/bin`, `TERM`, `LANG`, `TZ`; nothing else of
  the server's environment (no `DB_PARAMS`, no secrets);
- the system's programs and libraries (`/usr`, `/bin`, `/lib*`, `/etc`,
  `/opt`), only read; not the server's working directory, not `/proc`
  (Landlock), not the other users' folders.

The chip beside the terminal's state says how it is kept apart
(**Isolated · Landlock**, **Not isolated**); a click shows the details of
the session: the sandbox and what it covers, the folders and whether each
is written, the network (on, or off — ask the server's administrator), the
limits, when it closes idle.

Closing the page hangs the shell up (SIGHUP, then SIGKILL) and kills what it
left in its session (a program that made a session of its own escapes
it: the limits still hold). The log says who opened it, where, with which sandbox, and how it
ended (`terminal opened`, `terminal closed`); `OnAccess` is told
`terminal`.

## The sandbox

Package `sandbox`, usable by any program:

| Mode | How | Needs |
| --- | --- | --- |
| `bwrap` | bubblewrap: namespaces of the mounts, processes, network, IPC — the shell sees only what is bound | `bwrap` installed, and unprivileged user namespaces (often refused: Ubuntu's AppArmor `apparmor_restrict_unprivileged_userns=1`, Docker's default seccomp) |
| `landlock` | the kernel's Landlock, by the server's process: files (Linux 5.13), TCP (6.7, ABI 4), signals and abstract sockets (6.12, ABI 6) | Linux with Landlock (`/sys/kernel/security/lsm` lists `landlock`); works in Docker |
| `auto` | `bwrap` if it runs, else `landlock`, else none: the terminal off | — |
| `off` | no sandbox: the shell reaches whatever the server's user does | an explicit choice |

All modes: the limits (`setrlimit`, set before the shell starts), a session
of its own (`setsid`), `no_new_privs`, killed with the server.

**Below Linux 6.12** Landlock does not isolate the signals: the shell may
signal the processes of the server's user (`kill -9 -1`). Where that
matters, run the server as a user of its own, prefer `bwrap`, or turn the
terminal off. The log says what the sandbox covers at the start:
`terminal on … sandbox="landlock ABI 4: files, TCP (signals not isolated: Linux < 6.12)"`.

```go
mode, err := sandbox.Detect(sandbox.ModeAuto)  // err: why none works here
p, err := sandbox.Start(mode, sandbox.Spec{
	Argv:   []string{"/bin/sh", "-c", "make"},
	Dir:    dir,
	Env:    []string{"PATH=/usr/bin:/bin"},
	Paths:  []sandbox.Path{{Path: dir, Write: true}},
	Limits: sandbox.Limits{CPU: time.Minute, Memory: 1 << 30},
	TTY:    f, // stdin, stdout, stderr
})
err = p.Wait()
```

## Reusing the explorer

`UseAssets(presets)` serves `vx-files.js`, `vx-terminal.js`, `vx-files.css`
and xterm (`VendorPath(presets)`); senir uses them for the files of a user of
the OS.

- **The API** — package `explorer`, over any `webdav.FileSystem`:
  `explorer.New(fsys, explorer.Options{Allow, Decorate, Writable, IsRepo,
  Extract, Fetch, HTTPClient, Symlink, Changed})`, served at `…/api/fs/<op>`
  (the operations: the package's documentation).
- **`<vx-files>`** — `base` (the API's prefix), `start`, `title`,
  `terminal`, `terminal-url`, `term-root-var` (`FILES_ROOT`),
  `term-recorded`, `vendor`, `fetch`, `extract`, `ide-url`, `git-base`; a
  slot `toolbar` for buttons of its bar.
- **`<vx-terminal>`** — `ws-url`, `vendor`, `root-var`, `recorded`, `keep`,
  `fill`. Its websocket: binary messages are the bytes of the session; it
  sends `{"resize":[cols,rows]}`; the server ends with `{"exit":"…"}`; `cols`
  and `rows` in the query of the first request. A text message
  `{"info": {label, icon, color, details: [{label, values: [{text, warn}]}]}}`
  — fs-tools sends it first, in the user's language — is the chip beside the
  state and its details; a server that sends none has no chip.
- **Their words** are the namespace `vxFiles` of the page's vue-i18n
  (`FilesMessages`, en and pt-BR, registered by `UseAssets`): in the user's
  language, translatable as any module's messages.

## Tests

`go test ./admin/packages/fs-tools/...`: the explorer, the folders, the
sandbox (each mode that works here: read, write, a folder only read, one out
of reach, TCP) and the terminal by its websocket, behind a writer of a log.
