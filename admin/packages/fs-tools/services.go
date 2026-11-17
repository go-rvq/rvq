package fs_tools

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-rvq/rvq/admin/presets"
)

// registerServices registers the services of the file system — no page of
// the menu has their paths (presets.Service) —: fs-tools (the explorer's
// API, its terminal, its components' vendor), webdav, git.
func (b *Builder) registerServices(p *presets.Builder) {
	m := func(f func(m *Messages) string) func(context.Context) string {
		return func(ctx context.Context) string { return f(GetMessages(ctx)) }
	}
	p.Service(ServiceName).
		Title(m(func(m *Messages) string { return m.ServiceFsTitle })).
		Description(m(func(m *Messages) string { return m.ServiceFsDesc })).
		Route("GET", "/api/fs/list?path=…", presets.ServiceAuthSession, m(func(m *Messages) string { return m.RouteFsList })).
		Route("GET", "/api/fs/download?path=…", presets.ServiceAuthSession, m(func(m *Messages) string { return m.RouteFsDownload })).
		Route("POST", "/api/fs/{op}?path=…", presets.ServiceAuthSession, m(func(m *Messages) string { return m.RouteFsOp })).
		Route("GET", "/terminal?dir=…", presets.ServiceAuthSession+", !terminal", m(func(m *Messages) string { return m.RouteFsTerminal })).
		Route("GET", "/vendor/{file}", presets.ServiceAuthPublic, m(func(m *Messages) string { return m.RouteFsVendor })).
		Handle("/api/fs/", b.explorerHandler()).
		Handle("/terminal", b.terminalHandler())

	p.Service("webdav").
		Title(m(func(m *Messages) string { return m.ServiceDavTitle })).
		Description(m(func(m *Messages) string { return m.ServiceDavDesc })).
		Route("*", "/{path}", presets.ServiceAuthBasic, m(func(m *Messages) string { return m.RouteDav })).
		Outer(func(admin http.Handler) http.Handler { return b.WebDavHandler() })

	// served at the repositories' paths alone (<name>.git, Mount)
	p.Service(GitServiceName).MountOnly().
		Title(m(func(m *Messages) string { return m.ServiceGitTitle })).
		Description(m(func(m *Messages) string { return m.ServiceGitDesc })).
		Route("GET, POST", "<prefix>/{name}.git/…", presets.ServiceAuthBasic, m(func(m *Messages) string { return m.RouteGit })).
		Outer(func(admin http.Handler) http.Handler {
			if b.git == nil {
				return admin
			}
			return b.GitHandler()
		})
}

// gitMountPath is the path of the repository name in git's protocol:
// <prefix>/<name>.git/ (one of each key: <prefix>/<base>/, its keys'
// <key>.git under it).
func (b *Builder) gitMountPath(name string) string {
	base := strings.TrimSuffix(b.p.GetURIPrefix(), "/")
	if dir, ok := strings.CutSuffix(name, "/"+GitKeyPattern); ok {
		return base + "/" + dir + "/"
	}
	return base + "/" + name + ".git/"
}
