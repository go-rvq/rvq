package examples_presets

import (
	"github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

type confirmDialog struct{}

func PresetsConfirmDialog(b *presets.Builder, db *gorm.DB) (
	mb *presets.ModelBuilder,
	cl *presets.ListingBuilder,
	ce *presets.EditingBuilder,
	dp *presets.DetailingBuilder,
) {
	_ = []interface{}{
		// @snippet_begin(OpenConfirmDialog)
		presets.EventOpenConfirmDialog,
		// @snippet_end
		// @snippet_begin(ConfirmDialogConfirmEvent)
		presets.ConfirmDialogConfirmEvent,
		// @snippet_end
		// @snippet_begin(ConfirmDialogPromptText)
		presets.ConfirmDialogPromptText,
		// @snippet_end
		// @snippet_begin(ConfirmDialogDialogPortalName)
		presets.ConfirmDialogDialogPortalName,
		// @snippet_end
	}

	b.DataOperator(gorm2op.DataOperator(db))

	mb = b.Model(&confirmDialog{}).
		URIName("confirm-dialog").
		Label("Confirm Dialog")

	mb.Listing().PageFunc(func(ctx *web.EventContext) (r web.PageResponse, err error) {
		r.Body = htmlgo.Div(
			// @snippet_begin(ConfirmDialogSample)
			vuetify.VBtn("Delete File").
				Attr("@click",
					web.Plaid().
						EventFunc(presets.EventOpenConfirmDialog).
						Query(presets.ConfirmDialogConfirmEvent,
							`alert("file deleted")`,
						).
						Go(),
				),
			// @snippet_end
		).Class("ma-8")
		return r, nil
	})
	return
}
