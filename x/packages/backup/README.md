# backup

The states of an application — its database, its files, its binary, its
configuration —: what a client syncs into a directory with the
application's structure, where the application starts. The core: no admin
here (the model, the routes and the client come over it).

## The Builder

```go
b, err := backup.New(backup.Options{
	Root: ".",                // the application's directory: the paths of the states
	Dir:  "data/backup",      // the states, a folder each
	// Cache: nil → a FileCache at Dir/.hash-cache.gz
})
b.Source("db", func(ctx context.Context, st *backup.State) error {
	return st.FS.AddGenerated("data/db.pg", 0o600, func(w io.Writer) error { return dump(ctx, w) })
})
b.Source("site-files", func(ctx context.Context, st *backup.State) error {
	return st.FS.AddTree("public/.git")
})
b.Source("uploads", func(ctx context.Context, st *backup.State) error {
	return st.FS.AddHashed("data/system/files/a.jpg", sum, nil) // its sha256 known: not read
})

st := &backup.State{ID: "20261117-120000", Sources: nil /* all */}
err = b.BuildState(ctx, st) // st.Result: files, bytes, by source, the manifest's size and sum
```

A source registers through `st.FS`:

| Method | What |
| --- | --- |
| `Add(p)` | a file, a folder or a link (not followed) at `p` of the root; a file hashed now, or by the cache |
| `AddFrom(p, src)` | `src`, out of the root, at `p` of the state |
| `AddHashed(p, sum, fi)` | a file whose sha256 is known: not read |
| `AddGenerated(p, mode, write)` | content made now (a dump): kept in the state alone |
| `AddTree(p, exclude...)` / `AddTreeFrom(p, src, exclude...)` | `p` and all under it; `exclude`, patterns of `path.Match` of the paths under it (one with no `/` is of a name too: `.git`, `*.tmp`) |

A path registered twice the same is once; differently, `ErrDuplicate`.
Sockets, devices and pipes are left out; so are the states' folder and
`.backup` (the client's).

## A state

`BuildState` runs the sources (all, or `State.Sources`: a state without
some is **partial**), then writes the manifest. It is made in a folder
beside, renamed to `Dir/<id>` whole: a state there is complete. A source's
error, or the context's, fails it, and what it made is removed.

`Dir/<id>/`:

- `manifest.gz` — the manifest;
- `files/…` — the snapshot: a **hard link** of each file (the inode the
  manifest's sum is of: a file replaced by a rename leaves it as it was;
  no space taken), the generated files;
- `origins` — the files that could not be linked (another file system):
  where they are.

`OpenFile(id, entry)` opens a file of the state: its snapshot, else its
origin; one whose size or time is not the manifest's is `ErrChanged` (a
file written in place after the state). `Manifest(id)`, `Remove(id)`.

## The manifest

gzip; a line an entry, sorted by path, `\t` apart:

```
#backup-manifest  1  <id>  <created RFC3339Nano>  <sources, ",">  <partial 0|1>
<sha256 | ->  <mode octal>  <mtime unix ns>  <size>  <f|d|l>  <path>  [<link target>]
```

- the sha256 in base64url with no padding (43 characters; the ETag's);
- the mode as Unix's octal (`644`, `4755`: setuid 4000, setgid 2000,
  sticky 1000);
- the path relative to the root, `/` apart, clean, never `..` nor in
  `.backup`; a path (or a target) with a tab, a line break, a backslash, a
  control character, or starting with `"`, is quoted (Go's
  `strconv.Quote`).

`WriteManifest`, `ReadManifest` (refuses an unknown version, a path out of
the directory, entries out of order), `Manifest.Find` (by path).

## The cache of the sums

`CacheKey`: the path, the size, the time, the device and the inode — a
file not changed is not hashed again. `FileCache` keeps them in a file
(0600), saved after each state with only the sums used then (those of
files gone are dropped). A file written while hashed is `ErrChanged`.
