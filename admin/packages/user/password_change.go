package user

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
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
