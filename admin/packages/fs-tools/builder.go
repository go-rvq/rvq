package fs_tools

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/fs-tools/explorer"
	"github.com/go-rvq/rvq/admin/packages/fs-tools/fs"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/osenv"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/hack-pad/hackpadfs"
	"golang.org/x/net/webdav"
)

var EnvDirs = osenv.Get("RVQ_FSTOOLS_DIRS", "Directories available for webdav service. Values separated by "+
	"semi collon. Value format is MODE,MOUNT_POINT,LOCAL_PATH, when MODE=[r or w (default)]."+
	" Example: 'r,/media,data/media;w,/writable,/x/media'", "")

type Builder struct {
	lb      *login.Builder
	p       *presets.Builder
	mb      *presets.ModelBuilder
	davPath string
	fs      hackpadfs.FS
	log     *slog.Logger
	// page is the page of the files: its permissions are the WebDAV's
	page *presets.PageBuilder
	// mounts are the directories of each request (AddMount)
	mounts []*Mount
	// git serves the repositories added (AddGitRepo)
	git *gitServer
	// onAccess is told each request of the WebDAV and the git, its user
	// told (OnAccess)
	onAccess func(r *http.Request, kind string)
	// term is the terminal (mountTerminal)
	term *terminal
}

// The kinds of access OnAccess is told.
const (
	AccessWebDAV = "webdav"
	AccessGit    = "git"
	// AccessFiles: the explorer of the page; AccessTerminal, its terminal.
	AccessFiles    = "files"
	AccessTerminal = "terminal"
)

// OnAccess is told each request of the WebDAV (AccessWebDAV) and of the git
// (AccessGit), once its user is told (login.AuthOf: how): the log of its
// sessions.
func (b *Builder) OnAccess(f func(r *http.Request, kind string)) *Builder {
	b.onAccess = f
	return b
}

// recordAccess is h telling OnAccess its requests.
func (b *Builder) recordAccess(h http.Handler, kind string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if b.onAccess != nil {
			b.onAccess(r, kind)
		}
		h.ServeHTTP(w, r)
	})
}

func New(p *presets.Builder, ib *i18n.Builder, lb *login.Builder) (b *Builder, err error) {
	ConfigureMessages(ib)

	b = &Builder{
		lb:      lb,
		p:       p,
		davPath: p.ServicePath("webdav"),
		log:     web.NewLogger("fs-tools"),
	}

	if len(EnvDirs) > 0 {
		b.fs, err = fs.Parse(EnvDirs)
	} else {
		b.fs = fs.EmptyFS()
	}

	return
}

// explorerPath is where the explorer's API is: <prefix>/!fs-tools/api/fs/….
const explorerPath = "/" + presets.ServicePrefix + ServiceName + "/"

// pageFunc is the page of the files: the explorer (<vx-files>) over the
// file system the WebDAV serves; how to reach it by WebDAV, a drawer at the
// right opened by an icon of the explorer's bar.
func (b *Builder) pageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	m := GetMessages(ctx.Context())
	r.PageTitle = m.FileSystem
	term := b.terminalOn(ctx.R)
	r.Body = web.Scope(VContainer(
		h.Tag("vx-files").Attr("base", b.p.GetURIPrefix()+explorerPath).Attr("title", m.FileSystem).
			Attr(":fetch", "false").Attr(":extract", "true").
			Attr(":terminal", h.JSONString(term)).Attr("terminal-url", b.p.GetURIPrefix()+terminalPath).
			Attr("vendor", VendorPath(b.p)).Attr("term-root-var", "FILES_ROOT").
			// the folder of the address (?dir=…: the explorer keeps it there)
			Attr("start", ctx.R.URL.Query().Get("dir")).
			Children(
				h.Template(
					VBtn("").Icon("mdi-folder-network-outline").Size(SizeSmall).
						Attr(":variant", `locals.dav ? "tonal" : "text"`).
						Attr("title", m.WebDavProtocolTitle).Attr("aria-label", m.WebDavProtocolTitle).
						Attr("@click", "locals.dav = !locals.dav").Attr("data-fs-tools-webdav-open", true),
				).Attr("#toolbar", true),
			),
		vx.VXNavigationDrawer(h.Div(b.davInfo(ctx)).Class("pa-4")).
			Title(m.WebDavProtocolTitle).
			Location(LocationRight).
			Temporary(true).
			Closable(true).
			Expandable(true).
			Width("440").
			// its edge drags its width, kept in the browser
			Resizable("rvq.fstools.webdav.width").
			MinWidth(320).
			MaxWidth(900).
			Attr("v-model", "locals.dav").
			Attr("data-fs-tools-webdav", true),
	).Fluid(true).Class("pa-2")).LocalsInit(`{dav: false}`).Slot("{ locals }")
	return
}

// davInfo is how to reach the files by WebDAV: its address, the user and
// the password (an access key), the folders and whether each is written,
// the programs.
func (b *Builder) davInfo(ctx *web.EventContext) h.HTMLComponent {
	m := GetMessages(ctx.Context())
	scheme := "http"
	if ctx.R.TLS != nil || strings.EqualFold(ctx.R.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	url := fmt.Sprintf("%v://%v%v", scheme, ctx.R.Host, b.davPath)
	// the folders of RVQ_FSTOOLS_DIRS only read (r,…)
	var rows []h.HTMLComponent
	for _, d := range b.rootFolders(ctx.R) {
		mode := m.ReadOnly
		if !readOnlyDir("/"+d) && b.allowed(ctx.R, "/"+d, "", true) {
			mode = m.ReadWrite
		}
		rows = append(rows, h.Tr(h.Td(h.Code("/"+d)), h.Td(h.Text(mode))))
	}
	var folders h.HTMLComponent
	if len(rows) > 0 {
		folders = h.Div(
			h.Div(h.Strong(m.Folders)).Class("mt-3 mb-1"),
			VTable(h.Tbody(rows...)).Density(DensityCompact),
		)
	}
	var off h.HTMLComponent
	if err := b.terminalOff(ctx.R); err != nil {
		off = h.Div(h.Text(fmt.Sprintf(m.TerminalOff, err))).Class("mt-3 text-caption text-medium-emphasis").Attr("data-fs-tools-terminal-off", true)
	}
	return h.Div(
		h.Div(m.WebDavAccess(url)),
		h.Div(h.Text(m.WebDavCredentials)).Class("mt-2"),
		b.keyPermissions(ctx),
		folders,
		h.Div(m.WebDavProtocolSoftwareExample).Class("mt-2"),
		off,
	).Attr("data-fs-tools-webdav-info", true)
}

// keyPermissions is what an access key needs to reach the WebDAV: its
// lines (denies, on top of what its user may) must not deny the page's
// resource @get (reading) nor !write (writing).
func (b *Builder) keyPermissions(ctx *web.EventContext) h.HTMLComponent {
	if b.page == nil {
		return nil
	}
	ver := b.page.Page().ActionVerifier(ctx.R, "@get")
	if ver == nil {
		return nil
	}
	res := ver.PreferredResource()
	if res == "" {
		res = ver.Resource()
	}
	m := GetMessages(ctx.Context())
	return h.Div(
		h.Div(h.Strong(m.KeyPermissions)).Class("mb-1"),
		h.Div(h.Text(m.KeyPermissionsHint)),
		h.Div(h.Code(res)).Class("my-1"),
		VTable(h.Tbody(
			h.Tr(h.Td(h.Code("@get")), h.Td(h.Text(m.KeyPermRead))),
			h.Tr(h.Td(h.Code("!write")), h.Td(h.Text(m.KeyPermWrite))),
		)).Density(DensityCompact),
		h.Div(h.Text(m.KeyPermExample)).Class("mt-2"),
		h.Div(h.Code("deny "+res+"!write")).Class("mt-1"),
	).Class("mt-3").Attr("data-fs-tools-key-permissions", true)
}

// rootFolders are the names at the root of the file system of r: the
// folders of RVQ_FSTOOLS_DIRS, the mounts.
func (b *Builder) rootFolders(r *http.Request) (out []string) {
	fsys := b.fileSystem()
	ctx := context.WithValue(r.Context(), davRequestKey{}, r)
	f, err := fsys.OpenFile(ctx, "/", os.O_RDONLY, 0)
	if err != nil {
		return nil
	}
	defer f.Close()
	infos, _ := f.Readdir(-1)
	for _, fi := range infos {
		if fi.IsDir() {
			out = append(out, fi.Name())
		}
	}
	return
}

// readOnlyDir says p is in a folder of RVQ_FSTOOLS_DIRS only read (r,…).
func readOnlyDir(p string) bool {
	dirs, _ := fs.Dirs(EnvDirs)
	p = strings.Trim(p, "/")
	for _, d := range dirs {
		if d.ReadOnly && (p == d.Mount || strings.HasPrefix(p, d.Mount+"/")) {
			return true
		}
	}
	return false
}

// fileSystem is the file system of the WebDAV and of the explorer.
func (b *Builder) fileSystem() webdav.FileSystem {
	return &mountedFS{b: b, base: fs.NewWebDavFS(b.fs)}
}

// explorerHandler is the explorer's API (<prefix>/!fs-tools/api/fs/…) over
// the WebDAV's file system, each operation by its permissions (the page's,
// a mount's); what changes, told OnAccess (AccessFiles).
func (b *Builder) explorerHandler() http.Handler {
	x := explorer.New(b.fileSystem(), explorer.Options{
		Allow:    func(r *http.Request, p string, write bool) bool { return b.allowed(r, p, "", write) },
		ReadOnly: readOnlyDir,
		Extract:  true,
		Changed:  func(r *http.Request, op, p string) { b.log.Info("files", "op", op, "path", p, "user", whoOf(r)) },
	})
	prefix := b.p.GetURIPrefix() + strings.TrimSuffix(explorerPath, "/")
	h := http.StripPrefix(prefix, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// the read of the page: no session, nothing
		if !b.allowed(r, "/", "", false) {
			explorer.Fail(w, http.StatusForbidden, explorer.ErrNotAllowed)
			return
		}
		x.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), davRequestKey{}, r)))
	}))
	return b.recordAccess(h, AccessFiles)
}

func (b *Builder) Install(p *presets.Builder) (err error) {
	b.p = p
	page := p.PagesRegistrator().New(
		presets.HttpPage("/fs-tools").
			MenuIcon("mdi-file-cabinet").
			TitleFunc(func(ctx context.Context) string {
				return GetMessages(ctx).FileSystem
			})).
		Private().
		Layout(b.pageFunc)
	b.page = page
	// the explorer: its API, its components
	UseAssets(p)
	b.registerServices(p)
	// the WebDAV: reading asks the page's @get, writing !write
	page.Page().PermActions("!write", "!terminal").
		// its sets of permissions: to read, to write (titled as the general
		// ones)
		PermSet("readers", "@get", nil, nil).
		PermSet("authors", "{@get, !write}", nil, nil)

	defer page.Build()
	return
}

func (b *Builder) DavPath() string {
	return b.davPath
}

func (b *Builder) SetFS(fs hackpadfs.FS) *Builder {
	b.fs = fs
	return b
}

func (b *Builder) FS() hackpadfs.FS {
	return b.fs
}

func (b *Builder) Model() *presets.ModelBuilder {
	return b.mb
}

func (b *Builder) init() {
}

func (b *Builder) WebDavHandler() (h http.Handler) {
	return b.lb.BasichAuthMiddleware(b.recordAccess(b.davHandler(), AccessWebDAV))
}

// davHandler is the WebDAV of the user of the request (authenticated by
// WebDavHandler): its permissions verified.
func (b *Builder) davHandler() (h http.Handler) {
	log := b.log.With("handler", "webdav")

	h = &webdav.Handler{
		Prefix:     b.davPath,
		FileSystem: &mountedFS{b: b, base: fs.NewWebDavFS(b.fs)},
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			// We're totally abusing the logger here to update
			// the global lastOp, since it's called on every
			// request.

			// We do not count (or log) PROPFIND: on large
			// filesystems, it does a traversal of
			// everything, and tries to write ._<file>
			// properties files, and it's very spammy and
			// is unlikely to complete in a reasonable
			// time.

			if r.Method == "PROPFIND" {
				return
			}

			log := log.With("method", r.Method, "path", r.URL.Path)

			if err != nil {
				log.Error(err.Error())
			} else {
				log.Info("ok")
			}
		},
	}

	// each request by its permissions: the page's, a mount's
	return b.permissions(h)
}
