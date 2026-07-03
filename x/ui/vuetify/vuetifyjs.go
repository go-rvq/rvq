package vuetify

import (
	"io/fs"
	"net/http"
	"strings"

	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
	"github.com/theplant/osenv"
)

var customizeVuetifyCSS = osenv.GetBool("CUSTOMIZE_VUETIFY_CSS", "Use customized styles for vuetify", true)

func JSComponentsPack() web.ComponentsPack {
	return web.ComponentsPackBuilder(func(ctx *web.ComponentsPackBuilderContext) {
		FS, _ := fs.Sub(rvqjs.VuetifyRuntime, "vuetify/runtime-dist")
		ctx.AppendFile(FS, "vuetify.min.js")
		ctx.AppendFile(FS, "vuetify-labs.min.js")
		ctx.WriteString("window.Vuetify.__locale_messages = ")
		ctx.AppendFile(FS, "locale.json")
	})
}

func CSSComponentsPack() web.ComponentsPack {
	var v []byte
	var err error
	if customizeVuetifyCSS {
		v, err = rvqjs.Vuetify.ReadFile("vuetify/dist/vuetify/assets/index.css")
	} else {
		v, err = rvqjs.VuetifyRuntime.ReadFile("vuetify/runtime-dist/vuetify.min.css")
		// v2, err := assetsbox.ReadFile("dist/vuetify-labs.min.css")
		// if err != nil {
		//	panic(err)
		// }
		// v = append(v, v2...)
	}
	v = append(v, ';')
	if err != nil {
		panic(err)
	}

	return web.ComponentsPack(v)
}

func fontEot() web.ComponentsPack {
	v, err := rvqjs.Vuetify.ReadFile("vuetify/dist/vuetify/assets/materialdesignicons-webfont.eot")
	if err != nil {
		panic(err)
	}
	return web.ComponentsPack(v)
}

func fontTtf() web.ComponentsPack {
	v, err := rvqjs.Vuetify.ReadFile("vuetify/dist/vuetify/assets/materialdesignicons-webfont.ttf")
	if err != nil {
		panic(err)
	}
	return web.ComponentsPack(v)
}

func fontWoff() web.ComponentsPack {
	v, err := rvqjs.Vuetify.ReadFile("vuetify/dist/vuetify/assets/materialdesignicons-webfont.woff")
	if err != nil {
		panic(err)
	}
	return web.ComponentsPack(v)
}

func fontWoff2() web.ComponentsPack {
	v, err := rvqjs.Vuetify.ReadFile("vuetify/dist/vuetify/assets/materialdesignicons-webfont.woff2")
	if err != nil {
		panic(err)
	}
	return web.ComponentsPack(v)
}

type muxer interface {
	Handle(pattern string, handler http.Handler)
}

func HandleMaterialDesignIcons(prefix string, mux muxer) {
	mux.Handle(prefix+"/vuetify/assets/index.css", web.PacksHandler(
		"text/css",
		web.ComponentsPack(
			strings.ReplaceAll(string(CSSComponentsPack()), "/vuetify/assets/materialdesignicons", prefix+"/vuetify/assets/materialdesignicons")),
	))
	mux.Handle(prefix+"/vuetify/assets/materialdesignicons-webfont.eot", web.PacksHandler("application/vnd.ms-fontobject",
		fontEot()))
	mux.Handle(prefix+"/vuetify/assets/materialdesignicons-webfont.ttf", web.PacksHandler("font/ttf", fontTtf()))
	mux.Handle(prefix+"/vuetify/assets/materialdesignicons-webfont.woff", web.PacksHandler("font/woff", fontWoff()))
	mux.Handle(prefix+"/vuetify/assets/materialdesignicons-webfont.woff2", web.PacksHandler("font/woff2", fontWoff()))
}

const initVuetify = `(app, vueOptions) => app.use(Vuetify.createVuetify({{vuetifyOpts}}))`

const defaultVuetifyOpts = `{
  icons: {
    // defaultSet: 'md', // 'mdi' || 'mdiSvg' || 'md' || 'fa' || 'fa4'
  },
  locale: {
    messages: {
		'pt-BR': window.Vuetify.__locale_messages.pt,
		...window.Vuetify.__locale_messages
	}
  },
  theme: {
	themes: {
	  qor5: {
		dark: false,
		colors: {
		  primary:   "#3E63DD",
		  secondary: "#5B6471",
		  accent:    "#82B1FF",
		  error:     "#82B1FF",
		  info:      "#0091FF",
		  success:   "#30A46C",
		  warning:   "#F76808",
		}
	  },
	},
  },
}`

var vuetifyOpts string

func ChangeVuetifyOpts(opts string) {
	vuetifyOpts = opts
}

func Vuetify() string {
	if vuetifyOpts == "" {
		vuetifyOpts = defaultVuetifyOpts
	}
	return strings.NewReplacer("{{vuetifyOpts}}", vuetifyOpts).
		Replace(initVuetify)
}
