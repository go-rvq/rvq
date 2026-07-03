package tiptap

// @snippet_begin(TipTapPackrSample)
import (
	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
)

func JSComponentsPack() web.ComponentsPack {
	v, err := rvqjs.TipTap.ReadFile("tiptap/dist/tiptapjs.umd.cjs")
	if err != nil {
		panic(err)
	}

	return web.ComponentsPack(v)
}

func CSSComponentsPack() web.ComponentsPack {
	v, err := rvqjs.TipTap.ReadFile("tiptap/dist/tiptap.css")
	if err != nil {
		panic(err)
	}

	return web.ComponentsPack(v)
}

// @snippet_end
