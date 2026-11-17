package presets

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

// A path of the admin (after its prefix) is of a model, of a group of the
// menu or of a page — shown in the menu or not (a model hidden is still a
// model) —, or of a service: a path no page has, of a program, a download,
// an API, a protocol, its assets. A service is at ServicePrefix and its
// name, "/!backup/…", "/!webdav/…": a path that never moves with the menu
// (a client keeps it), and never meets a page's. A path whose form tells
// its protocol — "<name>.git/…", git over HTTP — is told apart from a
// page's by itself: a service may be there too (Mount), registered as the
// others.
//
// A package registers its service (Service) — never a path of its own out
// of a model, a group or a page —, says its routes (Route, Handle) and
// documents them in its user documentation, in a section of its stable
// paths.

// ServicePrefix marks a service's path: "/!<name>".
const ServicePrefix = "!"

var serviceNameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

// ServiceBuilder is a service of the admin: its name, what it is, its
// routes; served under the admin's login (Handle) or by an authentication
// of its own, before it (Outer).
type ServiceBuilder struct {
	b           *Builder
	name        string
	title, desc func(ctx context.Context) string
	routes      []*ServiceRoute
	handles     []serviceHandle
	outer       func(admin http.Handler) http.Handler
	// mounts are other paths its outer handler serves (Mount): a protocol's
	// (<name>.git)
	mounts []string
	// mountOnly: served at its mounts alone, not at its /!<name>
	mountOnly bool
}

// ServiceRoute is a route of a service, for its documentation: its method,
// its path under the service ("/{id}/manifest"), what it does, how it is
// authenticated.
type ServiceRoute struct {
	Method      string
	Path        string
	Description func(ctx context.Context) string
	// Auth says how it is reached: the ways (ServiceAuth*) and the
	// permissions ("!terminal") it asks, separated by ", " — written as
	// they are, the ways translated (GetAuth); the route's description
	// says whether one is enough or all are asked.
	Auth string
}

// The ways a service's route is reached (ServiceRoute.Auth), translated by
// GetAuth.
const (
	ServiceAuthPublic     = "public"      // anyone
	ServiceAuthSession    = "session"     // the admin's login
	ServiceAuthSignedLink = "signed link" // a link signed for the user
	ServiceAuthAccessKey  = "access key"  // an access key (Authorization: Bearer)
	ServiceAuthToken      = "token"       // a token the service gave
	ServiceAuthBasic      = "HTTP Basic"  // login and password, or an access key
)

// GetAuth is Auth in the language of ctx.
func (r *ServiceRoute) GetAuth(ctx context.Context) string {
	m := MustGetMessages(ctx)
	ways := map[string]string{
		ServiceAuthPublic:     m.ServiceAuthPublic,
		ServiceAuthSession:    m.ServiceAuthSession,
		ServiceAuthSignedLink: m.ServiceAuthSignedLink,
		ServiceAuthAccessKey:  m.ServiceAuthAccessKey,
		ServiceAuthToken:      m.ServiceAuthToken,
		ServiceAuthBasic:      m.ServiceAuthBasic,
	}
	parts := strings.Split(r.Auth, ", ")
	for i, w := range parts {
		if t, ok := ways[w]; ok && t != "" {
			parts[i] = t
		}
	}
	return strings.Join(parts, ", ")
}

// ServiceMessage is a text of the presets' messages in the language of
// ctx: a title, a description of a service of rvq.
func ServiceMessage(f func(m *Messages) string) func(ctx context.Context) string {
	return func(ctx context.Context) string { return f(MustGetMessages(ctx)) }
}

type serviceHandle struct {
	pattern string
	h       http.Handler
}

// Service is the service name of the admin, made once (a name registered
// again is the same service). The name is [a-z0-9][a-z0-9._-]*.
func (b *Builder) Service(name string) *ServiceBuilder {
	if !serviceNameRe.MatchString(name) {
		panic(fmt.Sprintf("presets: invalid service name %q", name))
	}
	for _, s := range b.services {
		if s.name == name {
			return s
		}
	}
	s := &ServiceBuilder{b: b, name: name}
	b.services = append(b.services, s)
	return s
}

// Services are the services registered, by name.
func (b *Builder) Services() []*ServiceBuilder {
	out := append([]*ServiceBuilder(nil), b.services...)
	sort.Slice(out, func(i, j int) bool { return out[i].name < out[j].name })
	return out
}

// ServicePath is the path of the service name: <prefix>/!<name>.
func (b *Builder) ServicePath(name string) string {
	return strings.TrimSuffix(b.prefix, "/") + "/" + ServicePrefix + name
}

// IsServicePath says p (a path of the admin's) is a service's.
func (b *Builder) IsServicePath(p string) bool {
	rest, ok := strings.CutPrefix(p, strings.TrimSuffix(b.prefix, "/")+"/"+ServicePrefix)
	return ok && rest != ""
}

// Name is the service's name.
func (s *ServiceBuilder) Name() string { return s.name }

// Path is where the service is: <prefix>/!<name>.
func (s *ServiceBuilder) Path() string { return s.b.ServicePath(s.name) }

// Title and Description say what the service is (its documentation).
func (s *ServiceBuilder) Title(f func(ctx context.Context) string) *ServiceBuilder {
	s.title = f
	return s
}

func (s *ServiceBuilder) Description(f func(ctx context.Context) string) *ServiceBuilder {
	s.desc = f
	return s
}

// GetTitle and GetDescription in ctx's language ("" when none).
func (s *ServiceBuilder) GetTitle(ctx context.Context) string {
	if s.title == nil {
		return s.name
	}
	return s.title(ctx)
}

func (s *ServiceBuilder) GetDescription(ctx context.Context) string {
	if s.desc == nil {
		return ""
	}
	return s.desc(ctx)
}

// Route documents a route of the service (Routes): path under the service
// ("/{id}/manifest"), or, of a path of a protocol (Mount), whole from the
// admin's prefix — "<prefix>/{name}.git/…" (FullPath).
func (s *ServiceBuilder) Route(method, path, auth string, desc func(ctx context.Context) string) *ServiceBuilder {
	s.routes = append(s.routes, &ServiceRoute{Method: method, Path: path, Auth: auth, Description: desc})
	return s
}

// Routes are the routes documented.
func (s *ServiceBuilder) Routes() []*ServiceRoute { return s.routes }

// RoutePath is the full path of the route r of the service: its path under
// the service, or "<prefix>" replaced by the admin's.
func (s *ServiceBuilder) RoutePath(r *ServiceRoute) string {
	if rest, ok := strings.CutPrefix(r.Path, "<prefix>"); ok {
		return strings.TrimSuffix(s.b.GetURIPrefix(), "/") + rest
	}
	return s.Path() + r.Path
}

// Handle serves h at pattern under the service — "GET /files/{id}",
// "/api/" — in the admin's mux, under its login.
func (s *ServiceBuilder) Handle(pattern string, h http.Handler) *ServiceBuilder {
	s.handles = append(s.handles, serviceHandle{pattern, h})
	return s
}

// Outer serves the service by an authentication of its own, before the
// admin's login (an access key, a token, HTTP Basic): f is given the
// admin's handler, for what is not the service's to answer (a request
// with no credential of it: the admin's login).
func (s *ServiceBuilder) Outer(f func(admin http.Handler) http.Handler) *ServiceBuilder {
	s.outer = f
	return s
}

// Mount serves paths (absolute; one ending in "/" a prefix) by the
// service's outer handler too (Outer), besides its /!<name>: a path whose
// form tells its protocol — "<prefix>/site-files.git/" (git over HTTP):
// no page has it, the client's tool needs it as it is. A path of before a
// service is not kept (no redirect): the services' paths are its own.
func (s *ServiceBuilder) Mount(paths ...string) *ServiceBuilder {
	s.mounts = append(s.mounts, paths...)
	return s
}

// MountOnly serves the service at its mounts alone (Mount), not at its
// /!<name>: a service of a protocol's paths (the repositories of git,
// <name>.git), whose name only registers and documents it.
func (s *ServiceBuilder) MountOnly() *ServiceBuilder {
	s.mountOnly = true
	return s
}

// IsMountOnly says the service is served at its mounts alone (MountOnly).
func (s *ServiceBuilder) IsMountOnly() bool { return s.mountOnly }

// Mounts are the paths of Mount.
func (s *ServiceBuilder) Mounts() []string { return s.mounts }

// IsOuter says the service has an authentication of its own (Outer).
func (s *ServiceBuilder) IsOuter() bool { return s.outer != nil }

// pattern is p (a pattern under the service, a method maybe before it) at
// the service's path.
func (s *ServiceBuilder) pattern(p string) string {
	method, path, ok := strings.Cut(p, " ")
	if !ok {
		method, path = "", p
	}
	full := s.Path() + path
	if method != "" {
		return method + " " + full
	}
	return full
}

// setupServices serves the services' handles in the admin's mux.
func (b *Builder) setupServices(mux *http.ServeMux) {
	for _, s := range b.services {
		for _, hd := range s.handles {
			mux.Handle(s.pattern(hd.pattern), hd.h)
		}
	}
}

// OuterServices is admin (the admin's handler, its login before it) with
// the services of an authentication of their own (Outer) in front: each at
// its path, the rest the admin's.
func (b *Builder) OuterServices(admin http.Handler) http.Handler {
	var outer []*ServiceBuilder
	for _, s := range b.services {
		if s.outer != nil {
			outer = append(outer, s)
		}
	}
	if len(outer) == 0 {
		return admin
	}
	mux := http.NewServeMux()
	for _, s := range outer {
		h := s.outer(admin)
		if !s.mountOnly {
			mux.Handle(s.Path()+"/", h)
			mux.Handle(s.Path(), h)
		}
		for _, l := range s.mounts {
			mux.Handle(l, h)
		}
	}
	mux.Handle("/", admin)
	return mux
}

// AssetsService is the service of the admin's own assets: its scripts and
// styles (main.js, vue.js, main.css), the icons (vuetify/assets/…), the
// extra assets of the packages (extra/…, ExtraAsset).
const AssetsService = "assets"

// registerAssetsService documents the assets' service.
func (b *Builder) registerAssetsService() {
	m := ServiceMessage
	b.Service(AssetsService).
		Title(m(func(m *Messages) string { return m.ServiceAssetsTitle })).
		Description(m(func(m *Messages) string { return m.ServiceAssetsDesc })).
		Route("GET", "/main.js", ServiceAuthPublic, m(func(m *Messages) string { return m.RouteAssetsMainJS })).
		Route("GET", "/vue.js", ServiceAuthPublic, m(func(m *Messages) string { return m.RouteAssetsVueJS })).
		Route("GET", "/main.css", ServiceAuthPublic, m(func(m *Messages) string { return m.RouteAssetsMainCSS })).
		Route("GET", "/vuetify/assets/{file}", ServiceAuthPublic, m(func(m *Messages) string { return m.RouteAssetsVuetify })).
		Route("GET", "/extra/{path}", ServiceAuthPublic, m(func(m *Messages) string { return m.RouteAssetsExtra }))
}
