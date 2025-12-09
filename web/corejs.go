package web

// @snippet_begin(PackrSample)
import (
	"embed"
	"io/fs"

	"github.com/theplant/osenv"
)

//go:embed corejs/dist/*.js
var box embed.FS

func JSComponentsPack() ComponentsPack {
	return ComponentsPackFromFile(box, "corejs/dist/index.js")
}

var webVueDebug = osenv.GetBool("WEB_VUE_DEBUG", "Use dev vue.js javascript source code to debug vue components", false)

func JSVueComponentsPack() ComponentsPack {
	return ComponentsPackBuilder(func(ctx *ComponentsPackBuilderContext) {
		name := "vue.global.prod.js"
		if webVueDebug {
			name = "vue.global.dev.js"
		}
		box, _ := fs.Sub(box, "corejs/dist")
		ctx.AppendFile(box, name)
		ctx.AppendFile(box, "vue-i18n.js")
	})
}

// @snippet_end
