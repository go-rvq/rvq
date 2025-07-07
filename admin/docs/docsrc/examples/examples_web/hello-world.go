package examples_web

// @snippet_begin(HelloWorldSample)
import (
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples"
	"github.com/go-rvq/rvq/web"
	. "github.com/theplant/htmlgo"
)

func HelloWorld(ctx *web.EventContext) (pr web.PageResponse, err error) {
	pr.Body = H1("Hello World")
	return
}

var HelloWorldPB = web.Page(HelloWorld) // this is already a http.Handler

var HelloWorldPath = examples.URLPathByFunc(HelloWorld)

// @snippet_end
