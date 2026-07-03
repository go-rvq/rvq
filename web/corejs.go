package web

// @snippet_begin(PackrSample)
import (
	"io/fs"

	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/theplant/osenv"
)

func JSComponentsPack() ComponentsPack {
	return ComponentsPackFromFile(rvqjs.CoreJS, "corejs/dist/index.js")
}

var webVueDebug = osenv.GetBool("WEB_VUE_DEBUG", "Use dev vue.js javascript source code to debug vue components", false)

func JSVueComponentsPack() ComponentsPack {
	return ComponentsPackBuilder(func(ctx *ComponentsPackBuilderContext) {
		name := "vue.global.prod.js"
		if webVueDebug {
			name = "vue.global.dev.js"
		}
		box, _ := fs.Sub(rvqjs.CoreJS, "corejs/dist")
		ctx.AppendFile(box, name)
		ctx.AppendFile(box, "vue-i18n.js")
	})
}

// @snippet_end
