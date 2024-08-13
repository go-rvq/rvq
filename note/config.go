package note

import (
	"fmt"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/i18n"
	. "github.com/qor5/x/v3/ui/vuetify"
	"github.com/sunfmin/reflectutils"
	h "github.com/theplant/htmlgo"
	"gorm.io/gorm"
)

const (
	I18nNoteKey i18n.ModuleKey = "I18nNoteKey"

	createNoteEvent     = "note_CreateNoteEvent"
	updateUserNoteEvent = "note_UpdateUserNoteEvent"
)

func tabsPanel(db *gorm.DB, mb *presets.ModelBuilder) presets.TabComponentFunc {
	return func(obj interface{}, ctx *web.EventContext) (tab h.HTMLComponent, content h.HTMLComponent) {
		id := ctx.Param(presets.ParamID)
		if len(id) == 0 {
			return
		}

		tn := mb.Info().Label()

		notesSection := getNotesTab(ctx, db, tn, id)

		msgr := i18n.MustGetModuleMessages(ctx.R, I18nNoteKey, Messages_en_US).(*Messages)

		userID, _ := GetUserData(ctx)
		count := GetUnreadNotesCount(db, userID, tn, id)

		notesTab := VBadge(h.Text(msgr.Notes)).
			Attr(":content", "locals.unreadNotesCount").
			Attr(":value", "locals.unreadNotesCount").
			Color("red")

		clickEvent := web.Plaid().
			BeforeScript("locals.unreadNotesCount=0").
			EventFunc(updateUserNoteEvent).
			Query("resource_id", id).
			Query("resource_type", tn).
			Go() + ";" + web.Plaid().
			EventFunc(actions.ReloadList).
			Go()
		if count == 0 {
			clickEvent = ""
		}

		const tabKey = "notesTab"

		tab = VTab(notesTab).
			Attr(web.VAssign("locals", fmt.Sprintf("{unreadNotesCount:%v}", count))...).
			Attr("@click", clickEvent).Value(tabKey)
		content = VTabsWindowItem(
			web.Portal().Name("notes-section").Children(notesSection),
		).Value(tabKey)

		return
	}
}

func noteFunc(db *gorm.DB, mb *presets.ModelBuilder) presets.FieldComponentFunc {
	return func(field *presets.FieldContext, ctx *web.EventContext) (c h.HTMLComponent) {
		tn := mb.Info().Label()
		obj := field.Obj

		id := fmt.Sprint(reflectutils.MustGet(obj, "ID"))
		if ps, ok := obj.(interface {
			PrimarySlug() string
		}); ok {
			id = ps.PrimarySlug()
		}

		latestNote := QorNote{}
		db.Model(&QorNote{}).Where("resource_type = ? AND resource_id = ?", tn, id).Order("created_at DESC").First(&latestNote)

		content := []rune(latestNote.Content)
		result := string(content[:])
		if len(content) > 60 {
			result = string(content[0:60]) + "..."
		}
		userID, _ := GetUserData(ctx)
		count := GetUnreadNotesCount(db, userID, tn, id)
		return h.Td(
			h.If(count > 0,
				VBadge(h.Text(result)).Content(count).Color("red"),
			).Else(
				h.Text(fmt.Sprint(result)),
			),
		)
	}
}
