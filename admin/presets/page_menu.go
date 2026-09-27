package presets

import (
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
