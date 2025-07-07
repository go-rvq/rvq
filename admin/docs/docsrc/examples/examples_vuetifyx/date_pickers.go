package examples_vuetifyx

// @snippet_begin(VuetifyxDatetimePickers)

import (
	"context"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples"
	"github.com/go-rvq/rvq/web"
)

type VXDateBuilder struct {
	tag *h.HTMLTagBuilder
}

func Vxdatepicker(children ...h.HTMLComponent) (r *VXDateBuilder) {
	r = &VXDateBuilder{
		tag: h.Tag("vx-datepicker").Children(children...),
	}
	return
}

func (b *VXDateBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	return b.tag.MarshalHTML(ctx)
}

func VuetifyxDatePickers(ctx *web.EventContext) (pr web.PageResponse, err error) {
	pr.Body = h.Div(
		Vxdatepicker(),
	)
	return
}

var DatePickersPB = web.Page(VuetifyxDatePickers)

var DatePickersPath = examples.URLPathByFunc(VuetifyxDatePickers)

// @snippet_end
