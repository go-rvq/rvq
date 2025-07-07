package login

import (
	"fmt"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/admin/v3/presets/actions"
	"github.com/qor5/web/v3"
	"github.com/qor5/x/v3/login"
)

const (
	OpenChangePasswordDialogEvent = "login_openChangePasswordDialog"
)

type Builder struct {
	b  *login.Builder
	vh *login.ViewHelper
}

func (b *Builder) Builder() *login.Builder {
	return b.vh.Builder()
}

func (b *Builder) Install(pb *presets.Builder) (err error) {
	vh := b.ViewHelper()
	r := vh.Builder()
	r.LoginPageFunc(defaultLoginPage(vh, pb))
	r.ForgetPasswordPageFunc(defaultForgetPasswordPage(vh, pb))
	r.ResetPasswordLinkSentPageFunc(defaultResetPasswordLinkSentPage(vh, pb))
	r.ResetPasswordPageFunc(defaultResetPasswordPage(vh, pb))
	r.ChangePasswordPageFunc(defaultChangePasswordPage(vh, pb))
	r.TOTPSetupPageFunc(defaultTOTPSetupPage(vh, pb))
	r.TOTPValidatePageFunc(defaultTOTPValidatePage(vh, pb))

	registerChangePasswordEvents(r, pb)
	return
}

func (b *Builder) ViewHelper() *login.ViewHelper {
	return b.vh
}

func New(b *login.Builder) *Builder {
	return &Builder{
		b:  b,
		vh: b.ViewHelper(),
	}
}

func registerChangePasswordEvents(b *login.Builder, pb *presets.Builder) {
	vh := b.ViewHelper()

	showVar := "showChangePasswordDialog"
	pb.GetWebBuilder().RegisterEventFunc(OpenChangePasswordDialogEvent, func(ctx *web.EventContext) (r web.EventResponse, err error) {
		r.UpdatePortal(
			actions.Dialog.PortalName(),
			changePasswordDialog(vh, ctx, showVar, defaultChangePasswordDialogContent(vh, pb)(ctx)),
		)

		web.AppendRunScripts(&r, fmt.Sprintf(`
(function(){
var tag = document.createElement("script");
tag.src = "%s";
tag.onload= function(){
	vars.meter_score = function(x){return zxcvbn(x).score+1};
}
document.getElementsByTagName("head")[0].appendChild(tag);
})()
        `, login.ZxcvbnJSURL))
		return
	})

	pb.GetWebBuilder().RegisterEventFunc("login_changePassword", func(ctx *web.EventContext) (r web.EventResponse, err error) {
		oldPassword := ctx.R.FormValue("old_password")
		password := ctx.R.FormValue("password")
		confirmPassword := ctx.R.FormValue("confirm_password")
		otp := ctx.R.FormValue("otp")

		msgr := login.GetMessages(ctx.Context())
		err = b.ChangePasswordT(login.GetCurrentUser(ctx.R).(login.UserPasser), true, ctx.R, oldPassword, password, confirmPassword, otp)
		if err != nil {
			msg := msgr.ErrorSystemError
			var color string
			if ne, ok := err.(*login.NoticeError); ok {
				msg = ne.Message
				switch ne.Level {
				case login.NoticeLevel_Info:
					color = "info"
				case login.NoticeLevel_Warn:
					color = "warning"
				case login.NoticeLevel_Error:
					color = "error"
				}
			} else {
				color = "error"
			}

			presets.ShowMessage(&r, msg, color)
			return r, nil
		}

		presets.ShowMessage(&r, msgr.InfoPasswordSuccessfullyChanged, "info")
		web.AppendRunScripts(&r, fmt.Sprintf("vars.%s = false", showVar))
		return r, nil
	})
}
