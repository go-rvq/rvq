package vuetifyx

import (
	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
)

func JSComponentsPack() web.ComponentsPack {
	v, err := rvqjs.VuetifyX.ReadFile("vuetifyx/dist/vuetifyxjs.umd.js")
	if err != nil {
		panic(err)
	}

	return web.ComponentsPack(v)
}

func CSSComponentsPack() web.ComponentsPack {
	v, err := rvqjs.VuetifyX.ReadFile("vuetifyx/dist/vuetifyxjs.css")
	if err != nil {
		panic(err)
	}

	return web.ComponentsPack(v)
}
