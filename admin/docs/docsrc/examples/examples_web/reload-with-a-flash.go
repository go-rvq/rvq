package examples_web

// @snippet_begin(ReloadWithFlashSample)
import (
	"fmt"
	"time"

	. "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples"
	"github.com/go-rvq/rvq/web"
)

var count int

func ReloadWithFlash(ctx *web.EventContext) (pr web.PageResponse, err error) {
	var msg HTMLComponent

	if d, ok := ctx.Flash.(*Data1); ok {
		msg = Div().Text(d.Msg).Style("border: 5px solid orange;")
	} else {
		count = 0
	}

	pr.Body = Div(
		H1("Whole Page Reload With a Flash"),
		msg,
		Div().Text(time.Now().Format(time.RFC3339Nano)),
		Button("Do Something").
			Attr("@click", web.POST().EventFunc("update2").Go()),
	)
	return
}

type Data1 struct {
	Msg string
}

func update2(ctx *web.EventContext) (er web.EventResponse, err error) {
	count++
	ctx.Flash = &Data1{Msg: fmt.Sprintf("The page is reloaded: %d", count)}
	er.Reload = true
	return
}

var ReloadWithFlashPB = web.Page(ReloadWithFlash).EventFunc("update2", update2)

var ReloadWithFlashPath = examples.URLPathByFunc(ReloadWithFlash)

// @snippet_end
