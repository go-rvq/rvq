package presets

import (
	"strings"

	"github.com/go-rvq/rvq/web"
)

func (m *ModelBuilder) isMenuItemActive(ctx *web.EventContext) bool {
	href := m.Info().ListingHref(ParentsModelID(ctx.R)...)
	if m.link != "" {
		href = m.link
	}
	path := strings.TrimSuffix(ctx.R.URL.Path, "/")
	if path == "" && href == "/" {
		return true
	}
	if path == href {
		return true
	}
	if href == m.p.prefix {
		return false
	}
	if href != "/" && strings.HasPrefix(path, href) {
		return true
	}

	return false
}
