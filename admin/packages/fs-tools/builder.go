package fs_tools

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/fs-tools/fs"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/login"
	"github.com/go-rvq/rvq/x/osenv"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
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
}

// The kinds of access OnAccess is told.
const (
	AccessWebDAV = "webdav"
	AccessGit    = "git"
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
		davPath: p.GetURIPrefix() + "/!webdav",
		log:     web.NewLogger("fs-tools"),
	}

	if len(EnvDirs) > 0 {
		b.fs, err = fs.Parse(EnvDirs)
	} else {
		b.fs = fs.EmptyFS()
	}

	return
}

func (b *Builder) pageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	scheme := "http"
	if ctx.R.TLS != nil {
		scheme = "https"
	}
	m := GetMessages(ctx.Context())
	r.Body = VContainer(h.Div(
		VCard(
			VCardText(
				m.WebDavAccess(fmt.Sprintf("%v://%v%v", scheme, ctx.R.Host, b.davPath)),
				m.WebDavProtocolSoftwareExample,
			),
		).Title(m.WebDavProtocolTitle),
	))
	return
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
	// the WebDAV: reading asks the page's @get, writing !write
	page.Page().PermActions("!write")

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
