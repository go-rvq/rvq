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
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/hack-pad/hackpadfs"
	"github.com/theplant/osenv"
	"golang.org/x/net/webdav"
)

var EnvDirs = osenv.Get("IRVQ_FSTOOLS_DIRS", "Directories available for webdav service. Values separated by "+
	"semi collon. Value format is MODE,MOUNT_POINT,LOCAL_PATH, when MODE=[r or w (default)]."+
	" Example: 'r,/media,data/media;w,/writable,/x/media'", "")

type Builder struct {
	lb      *login.Builder
	p       *presets.Builder
	mb      *presets.ModelBuilder
	davPath string
	fs      hackpadfs.FS
	log     *slog.Logger
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
	return
}

func (b *Builder) WebDavHandler() (h http.Handler) {
	log := b.log.With("handler", "webdav")

	h = &webdav.Handler{
		Prefix:     b.davPath,
		FileSystem: fs.NewWebDavFS(b.fs),
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

	h = b.lb.BasichAuthMiddleware(h)
	return
}
