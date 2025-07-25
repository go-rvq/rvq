package media

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
)

const (
	mediaLibraryListField = "media-library-list"
)

func configList(b *presets.Builder, mb *Builder) {
	mm := b.Model(&media_library.MediaLibrary{},
		presets.ModelConfig().
			SetModuleKey(I18nMediaLibraryKey)).
		MenuIcon("mdi-multimedia").
		URIName("media-library")

	mm.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.PageTitle = mm.TTitlePlural(ctx.Context())
		keyword := ctx.R.FormValue("keyword")
		ctx.R.Form.Set(searchKeywordName(mediaLibraryListField), keyword)
		r.Body = h.Components(
			web.Portal().Name(deleteConfirmPortalName(mediaLibraryListField)),
			web.Portal(
				h.Input("").
					Type("hidden").
					Attr(web.VField(searchKeywordName(mediaLibraryListField), keyword)...),
				fileChooserDialogContent(mb, mediaLibraryListField, ctx, &media_library.MediaBoxConfig{}),
			).Name(dialogContentPortalName(mediaLibraryListField)),
		)
		return
	})
}
