# The paths of the admin

Every path of the admin — what follows its prefix (`/admin`) — is one of:

1. **A model's, a menu group's or a page's**: what presets makes of the
   models (`p.Model`), the groups of the menu and the pages
   (`PagesRegistrator`): `/admin/<model>`, `/admin/<model>/<id>`,
   `/admin/<page>`. A model hidden from the menu (`InMenu(false)`)
   is still a model: its paths are of this kind.
   A menu group's path (`/admin/<group>`, `/admin/orgs/{id}/<group>`
   under a record) is its page: the page it is made of
   (`PageMenuGroup`: the group's name, title and icon the page's), else its
   index — the items the user sees, a tree with their icons and
   descriptions —; a page already at that path keeps it (go-rvq/rvq#19).
2. **A service's**: a path no page has — of a program, a download, an API,
   a protocol, the assets of the pages —, **registered** with
   `presets.Builder.Service`, in one of two forms:
   - **`/!<name>/…`**: the service's name after `!`
     (`presets.ServicePrefix`): `/admin/!backup/latest`,
     `/admin/!fs-tools/api/fs/list`, `/admin/!webdav/…`. The `!` never
     begins a model's or a page's path: the two never meet, whatever the
     menu becomes.
   - **a path whose form tells its protocol**: `/admin/site-files.git/…`
     (git over HTTP: the suffix `.git`). No page has such a path; the
     client's tool needs it as it is. It is mounted on its service
     (`ServiceBuilder.Mount`) and registered as the others.

Nothing else: a package never serves a path of its own out of these (no
`mux.Handle(prefix+"/my-thing/…")`).

## Why

- **Stable**: a service's path is the same whatever the menu (a model moved
  to another group changes its pages' paths, never a service's). A client —
  a script, a remote of git, a mail with a link — keeps it.
- **Told apart**: a path is a page's or a service's at a glance; a service
  never shadows a page, a page never a service.
- **Known**: the services are registered — listed, documented (the user
  documentation's page of the stable paths), checked.

## Registering a service

```go
svc := p.Service("backup").                       // /admin/!backup
	Title(func(ctx context.Context) string { … }).
	Description(func(ctx context.Context) string { … }).
	// its routes, documented: method, path under the service, how it is
	// authenticated, what it does
	Route("GET", "/latest", "access key", func(ctx context.Context) string { … }).
	Route("GET", "/{id}/manifest", "token, or the admin's session", …)

// served under the admin's login (its session, its permissions)
svc.Handle("GET /{id}/manifest", http.HandlerFunc(serveManifest))

// served by a login of its own (HTTP Basic, an access key, a token), before
// the admin's: the application mounts presets.Builder.OuterServices(admin)
// in front of the admin's handler; what the service does not answer (a
// request with no credential of it) it gives to admin
svc.Outer(func(admin http.Handler) http.Handler { … })

// a path of a protocol, on the same outer handler; MountOnly: served there
// alone, not at /!git (its name only registers and documents it)
p.Service("git").MountOnly().Mount("/admin/site-files.git/")
```

- `Service(name)` — the name is `[a-z0-9][a-z0-9._-]*`; registered once
  (the same name is the same service).
- `ServicePath(name)` is `<prefix>/!<name>`; `IsServicePath(p)` tells one.
- `Services()` lists them (the documentation's table, the checks).
- `AssetsService` (`/!assets`) is presets' own: `main.js`, `vue.js`,
  `main.css`, `vuetify/assets/…`, `extra/…` (`ExtraAsset`).

## The services of rvq

| Service | Package | Login | What |
| --- | --- | --- | --- |
| `/!assets` | presets | public (static) | the scripts, styles, icons of the pages; the packages' extra assets |
| `/!login` | admin/login | public | the login dialog over a page whose session ended |
| `/!gadcode` | admin/gadcode | session | the operations of the editor of Gad code |
| `/!notifications` | packages/notifications | session | the attachments of the notifications |
| `/!fs-tools` | packages/fs-tools | session | the explorer's API, its terminal, its vendor |
| `/!webdav` | packages/fs-tools | HTTP Basic (own) | the file system by WebDAV |
| `<name>.git` | packages/fs-tools (git) | HTTP Basic (own) | the repositories by git over HTTP |
| `/!user-docs` | userdocs | session | the pictures of the documents of the user documentation |
| `/!backup` | packages/backup | access key, token (own); session | the backup's states to a client |

## Documenting

Each package documents its services in its **user documentation**, in a
section of its own — **"Stable paths"** (pt-BR: **"Caminhos estáveis"**) —:
every route (method, full path), how it is authenticated, what it gives,
an example (`curl`). The guide `guides/stable-paths` of userdocs
(`userdocs.StablePathsNode`) explains the rule to the user and lists every
service from the registry (`admin.services()`, `userdocs.ServicesMarkdown`):
the application adds it to its tree.

A route's access (`Route`'s `auth`) is written with the constants
`ServiceAuthPublic`, `ServiceAuthSession`, `ServiceAuthSignedLink`,
`ServiceAuthAccessKey`, `ServiceAuthToken`, `ServiceAuthBasic` (translated
by `ServiceRoute.GetAuth`), and the permissions it asks (`!terminal`), `, `
apart. A route of a protocol's path is written whole from the prefix:
`<prefix>/{name}.git/…` (`RoutePath`).

## Checking

`presets.Builder.IsServicePath` tells a service's path; a test of the
application may walk its routes and require each to be a model's, a
group's, a page's or a service's.
