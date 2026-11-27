// Package charts serves the code of the admin's interactive charts — Apache
// ECharts with what vx-chart (vuetifyx) draws (rvq js.ECharts) — as the
// service "charts" of the admin: GET /!charts/echarts.js, behind its login,
// cached a year (its address carries the version of its content). A page
// does not load it: vx-chart does, the first time it shows a chart, from
// Src.
package charts

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"net/http"
	"sync"

	"github.com/go-rvq/rvq/admin/presets"
	rvqjs "github.com/go-rvq/rvq/js"
)

// Name is the service's name: its path is /!charts.
const Name = "charts"

// script is the file served: ECharts, minified.
const script = "echarts/dist/echarts.min.js"

var (
	body    []byte
	version string
	once    sync.Once
)

func load() {
	once.Do(func() {
		b, err := fs.ReadFile(rvqjs.ECharts, script)
		if err != nil {
			panic("charts: " + err.Error())
		}
		sum := sha256.Sum256(b)
		body, version = b, hex.EncodeToString(sum[:6])
	})
}

// Setup registers the service on p (once: a second call is the same).
func Setup(p *presets.Builder) {
	load()
	m := presets.ServiceMessage
	s := p.Service(Name)
	if len(s.Routes()) > 0 {
		return
	}
	s.Title(m(func(m *presets.Messages) string { return m.ServiceChartsTitle })).
		Description(m(func(m *presets.Messages) string { return m.ServiceChartsDesc })).
		Route("GET", "/echarts.js", presets.ServiceAuthSession, m(func(m *presets.Messages) string { return m.RouteChartsScript })).
		Handle("GET /echarts.js", http.HandlerFunc(serve))
}

// Src is where vx-chart loads the charts' code from, in the admin p: the
// service's script, its version in the address.
func Src(p *presets.Builder) string {
	load()
	return p.ServicePath(Name) + "/echarts.js?v=" + version
}

func serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	if r.URL.Query().Get("v") == version {
		w.Header().Set("Cache-Control", "private, max-age=31536000, immutable")
	} else {
		w.Header().Set("Cache-Control", "no-cache")
	}
	_, _ = w.Write(body)
}
