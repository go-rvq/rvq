package cropper

import (
	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
)

func JSComponentsPack() web.ComponentsPack {
	v, err := rvqjs.Cropper.ReadFile("cropper/dist/cropperjs.umd.cjs")
	if err != nil {
		panic(err)
	}

	return web.ComponentsPack(v)
}

func CSSComponentsPack() web.ComponentsPack {
	v, err := rvqjs.Cropper.ReadFile("cropper/dist/cropperjs.css")
	if err != nil {
		panic(err)
	}

	return web.ComponentsPack(v)
}
