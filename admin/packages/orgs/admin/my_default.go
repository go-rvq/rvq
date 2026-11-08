package admin

import (
	"context"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/packages/orgs/messages"
	"github.com/go-rvq/rvq/admin/packages/orgs/models"
	"github.com/go-rvq/rvq/admin/packages/user"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// The default organization of the user (models.UserDefault): the one the
// menu opens when none was chosen — any they reach, theirs or shared with
// them. Made in the detail of one (Make it my default), shown in the
// listing (My default).

const (
	ActionSetMyDefault   = "set_my_default"
	ActionUnsetMyDefault = "unset_my_default"
)

// currentUserID is the user of the request (uuid.Nil: none).
func currentUserID(ctx *web.EventContext) uuid.UUID {
	if u := user.GetCurrentUser(ctx.R); u != nil && !u.Anonymous() {
		return u.GetID()
	}
	return uuid.Nil
}

func configureMyDefault(mb *presets.ModelBuilder, db *gorm.DB) {
	isDefault := func(obj any, ctx *web.EventContext) bool {
		uid := currentUserID(ctx)
		o, _ := obj.(*models.Organizacao)
		return uid != uuid.Nil && o != nil && models.DefaultOf(db.WithContext(ctx.R.Context()), uid) == o.ID
	}

	mb.Listing().Field("MyDefault").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		if isDefault(field.Obj, ctx) {
			return h.Td(VIcon("mdi-check").Size(SizeSmall)).Attr("data-org-my-default", true)
		}
		return h.Td()
	})

	dt := mb.Detailing()
	dt.Action(ActionSetMyDefault).Icon("mdi-star-outline").
		SetI18nLabel(func(ctx context.Context) string { return messages.Get(ctx).SetMyDefault }).
		SetEnabledObj(func(obj any, _ string, ctx *web.EventContext) (bool, error) {
			return currentUserID(ctx) != uuid.Nil && !isDefault(obj, ctx), nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			org, err := uuid.Parse(id)
			if err != nil {
				return err
			}
			if err := models.SetDefault(db.WithContext(ctx.R.Context()), currentUserID(ctx), org); err != nil {
				return err
			}
			ctx.Flash = messages.Get(ctx.Context()).MyDefaultSet
			return nil
		})
	dt.Action(ActionUnsetMyDefault).Icon("mdi-star-off-outline").
		SetI18nLabel(func(ctx context.Context) string { return messages.Get(ctx).UnsetMyDefault }).
		SetEnabledObj(func(obj any, _ string, ctx *web.EventContext) (bool, error) {
			return isDefault(obj, ctx), nil
		}).
		UpdateFunc(func(id string, ctx *web.EventContext) error {
			if err := models.ClearDefault(db.WithContext(ctx.R.Context()), currentUserID(ctx)); err != nil {
				return err
			}
			ctx.Flash = messages.Get(ctx.Context()).MyDefaultUnset
			return nil
		})
}
