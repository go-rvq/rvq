package examples_web

// @snippet_begin(HelloWorldReloadSample)
import (
	"time"

	. "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples"
	"github.com/go-rvq/rvq/web"
)

func HelloWorldReload(ctx *web.EventContext) (pr web.PageResponse, err error) {
	pr.Body = Div(
		H1("Hello World"),
		Text(time.Now().Format(time.RFC3339Nano)),
		Button("Reload Page").Attr("@click", web.GET().
			EventFunc(reloadEvent).
			Go()),
	)
	return
}

func update(ctx *web.EventContext) (er web.EventResponse, err error) {
	er.Reload = true
	return
}

const reloadEvent = "reload"

var HelloWorldReloadPB = web.Page(HelloWorldReload).
	EventFunc(reloadEvent, update)

var HelloWorldReloadPath = examples.URLPathByFunc(HelloWorldReload)

// @snippet_end
