package examples_admin

import (
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/richeditor"
	"github.com/go-rvq/rvq/web"
)

func layoutPieces() {
	var presetsBuilder *presets.Builder

	// @snippet_begin(CustomizeVuetifyOptions)
	presetsBuilder.VuetifyOptions(`
        {
            icons: {
                iconfont: 'md',
            },
            theme: {
                themes: {
                    light: {
                        primary: "#673ab7",
                        secondary: "#009688",
                        accent: "#ff5722",
                        error: "#f44336",
                        warning: "#ff9800",
                        info: "#8bc34a",
                        success: "#4caf50"
                    },
                },
            },
        }
    `)
	// @snippet_end

	// @snippet_begin(InjectAssetViaExtraAsset)
	presetsBuilder.ExtraAsset("/redactor.js", "text/javascript", richeditor.JSComponentsPack())
	presetsBuilder.ExtraAsset("/redactor.css", "text/css", richeditor.CSSComponentsPack())
	// @snippet_end

	// @snippet_begin(InjectAssetViaAssetFunc)
	presetsBuilder.AssetFunc(func(ctx *web.EventContext) {
		ctx.Injector.Meta(web.MetaKey("charset"), "charset", "utf8")
		ctx.Injector.HeadHTML(`<script src="https://cdn.example.com/hello.js"></script>`)
	})
	// @snippet_end
}
