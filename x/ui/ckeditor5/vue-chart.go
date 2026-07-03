package ckeditor5

import (
	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
)

func JSComponentsPack() web.ComponentsPack {
	v, err := rvqjs.CKEditor5.ReadFile("ckeditor5/dist/vue-chart.umd.cjs")
	if err != nil {
		panic(err)
	}

	return web.ComponentsPack(v)
}
