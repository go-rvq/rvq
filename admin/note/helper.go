package note

import (
	"fmt"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/ui/vuetify"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"
)

type contextUserIDKey int

const (
	UserIDKey contextUserIDKey = iota
	UserKey
)

func GetUserData(ctx *web.EventContext) (userID uint, creator string) {
	if ctx.R.Context().Value(UserIDKey) != nil {
		userID = ctx.R.Context().Value(UserIDKey).(uint)
	}
	if ctx.R.Context().Value(UserKey) != nil {
		creator = ctx.R.Context().Value(UserKey).(string)
	}

	return
}

func getNotesTab(ctx *web.EventContext, db *gorm.DB, resourceType string, resourceId string) h.HTMLComponent {
	msgr := i18n.MustGetModuleMessages(ctx.Context(), I18nNoteKey, Messages_en_US).(*Messages)

	c := h.Div(
		web.Scope(
			VCardText(
				h.Text(msgr.NewNote),
				VRow(VCol(VTextField().Attr(web.VField("Content", "")...).Clearable(true))),
			),
			VCardActions(h.Components(
				VSpacer(),
				VBtn(presets.MustGetMessages(ctx.Context()).Create).
					Color("primary").
					Attr("@click", web.Plaid().
						EventFunc(createNoteEvent).
						Query("resource_id", resourceId).
						Query("resource_type", resourceType).
						Go()+";"+web.Plaid().
						EventFunc(actions.ReloadList).
						Go(),
					),
			)),
		).Slot("{form}"),
	)

	var notes []QorNote
	db.Where("resource_type = ? and resource_id = ?", resourceType, resourceId).
		Order("id DESC").Find(&notes)

	var panels []h.HTMLComponent
	for _, note := range notes {
		panels = append(panels, vuetify.VCard(
			vuetify.VCardText(
				h.Div(h.Text(fmt.Sprintf("%v - %v", note.Creator, note.CreatedAt.Format("2006-01-02 15:04:05 MST")))).
					Class("text-h6"),
				h.Text(note.Content)),
		))
	}
	c.AppendChildren(panels...).Class("p-2")
	return c
}
