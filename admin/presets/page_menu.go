package presets

import (
	"regexp"
	"strings"

	"github.com/go-rvq/rvq/web"
)

func (m *HttpPageBuilder) isMenuItemActive(ctx *web.EventContext) bool {
	href := m.pageHandler.path

	path := strings.TrimSuffix(ctx.R.URL.Path, "/")

	if path == "" && href == "/" {
		return true
	}

	if path == href {
		return true
	}

	if href != "/" && strings.HasPrefix(path, href) {
		return true
	}

	return false
}

// pathParam is a parameter of a pattern of http.ServeMux ("{user}",
// "{rest...}").
var pathParam = regexp.MustCompile(`\{([^{}.]+)(\.\.\.)?\}`)

// trail appends to bc where the page is: the groups of the menu it is in
// (their titles), and the pages it is under (Sub) — their titles, linked,
// the parameters of their paths the request's.
func (m *HttpPageBuilder) trail(ctx *web.EventContext, bc *BreadcrumbsBuilder) {
	if g := m.GetMenuGroupBuilder(); g != nil {
		for _, g := range append(g.Ancestors(), g) {
			bc.Append(&Breadcrumb{Label: g.TTitle(ctx.Context())})
		}
	}
	var parents []*HttpPageBuilder
	for p := m.parent; p != nil; p = p.parent {
		parents = append([]*HttpPageBuilder{p}, parents...)
	}
	for _, p := range parents {
		href := pathParam.ReplaceAllStringFunc(p.fullPath, func(s string) string {
			return ctx.R.PathValue(pathParam.FindStringSubmatch(s)[1])
		})
		bc.Append(&Breadcrumb{Label: p.TTitle(ctx.Context()), URI: href})
	}
}
