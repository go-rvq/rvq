package user

import (
	"context"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/login"
)

type ChangePassword struct {
	OldPassword     string `admin:"required"`
	NewPassword     string `admin:"required"`
	ConfirmPassword string `admin:"required"`
}

func ConfigureChangePasswordAction(lb *login.Builder, currentPasswordCheck bool, d *presets.DetailingBuilder, getUser func(ctx *web.EventContext) (u User, err error)) {
	p := d.ModelBuilder().Builder()
	cpfe := NewModel(p, &ChangePassword{}, presets.ModelConfig().SetSingleton(true)).
		ChildOf(d.ModelBuilder()).
		Editing()

	if currentPasswordCheck {
		cpfe.Field("OldPassword").
			ComponentFunc(presets.PasswordFieldComponentFunc).
			SetI18nLabel(func(ctx context.Context) string {
				return login.GetMessages(ctx).ChangePasswordOldLabel
			})
	} else {
		cpfe.Field("OldPassword").SetEnabled(func(ctx *presets.FieldContext) bool {
			return currentPasswordCheck
		})
	}

	cpfe.Field("NewPassword").
		ComponentFunc(presets.PasswordFieldComponentFunc).
		SetI18nLabel(func(ctx context.Context) string {
			return login.GetMessages(ctx).ChangePasswordNewLabel
		})

	cpfe.Field("ConfirmPassword").
		ComponentFunc(presets.PasswordFieldComponentFunc).
		SetI18nLabel(func(ctx context.Context) string {
			return login.GetMessages(ctx).ChangePasswordNewConfirmLabel
		})

	presets.ActionForm(d.Action("change_password").
		Icon("mdi-form-textbox-password").
		SetI18nLabel(func(ctx context.Context) string {
			return login.GetMessages(ctx).ChangePasswordTitle
		}),
		cpfe,
		func(ctx *presets.ActionFormContext[*ChangePassword]) (err error) {
			var u User
			if u, err = getUser(ctx.Context); err != nil {
				return
			}
			if u == nil {
				return presets.ErrRecordNotFound
			}

			return lb.ChangePasswordT(u, currentPasswordCheck, ctx.Context.R, ctx.Form.OldPassword, ctx.Form.NewPassword, ctx.Form.ConfirmPassword, "")
		}).Build()
}
