package login

import (
	"fmt"
	"net/http"

	. "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/login"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

var DefaultViewCommon = &ViewCommon{
	WrapperClass: "d-flex pt-16 flex-column mx-auto",
	WrapperStyle: "max-width: 28rem;",
	TitleClass:   "text-h5 mb-6 font-weight-bold",
	LabelClass:   "d-block mb-1 grey--text text--darken-2 text-sm-body-2",
}

type ViewCommon struct {
	WrapperClass string
	WrapperStyle string
	TitleClass   string
	LabelClass   string
}

func (vc *ViewCommon) Notice(vh *login.ViewHelper, msgr *login.Messages, w http.ResponseWriter, r *http.Request) HTMLComponent {
	var nn HTMLComponent
	if n := vh.GetNoticeFlash(w, r); n != nil && n.Message != "" {
		switch n.Level {
		case login.NoticeLevel_Info:
			nn = vc.InfoNotice(n.Message)
		case login.NoticeLevel_Warn:
			nn = vc.WarnNotice(n.Message)
		case login.NoticeLevel_Error:
			nn = vc.ErrNotice(n.Message)
		}
	}
	return Components(
		vc.ErrNotice(vh.GetFailFlashMessage(msgr, w, r)),
		vc.WarnNotice(vh.GetWarnFlashMessage(msgr, w, r)),
		vc.InfoNotice(vh.GetInfoFlashMessage(msgr, w, r)),
		nn,
	)
}

func (vc *ViewCommon) ErrNotice(msg string) HTMLComponent {
	if msg == "" {
		return nil
	}

	return VAlert(Text(msg)).
		Density(DensityCompact).
		Class("text-center").
		Icon(false).
		Type("error")
}

func (vc *ViewCommon) WarnNotice(msg string) HTMLComponent {
	if msg == "" {
		return nil
	}

	return VAlert(Text(msg)).
		Density(DensityCompact).
		Class("text-center").
		Icon(false).
		Type("warning")
}

func (vc *ViewCommon) InfoNotice(msg string) HTMLComponent {
	if msg == "" {
		return nil
	}

	return VAlert(Text(msg)).
		Density(DensityCompact).
		Class("text-center").
		Icon(false).
		Type("info")
}

func (vc *ViewCommon) ErrorBody(msg string) HTMLComponent {
	return Div(
		Text(msg),
	)
}

func (vc *ViewCommon) Input(
	id string,
	placeholder string,
	val string,
) *VTextFieldBuilder {
	return VTextField().
		Attr("name", id).
		Id(id).
		Placeholder(placeholder).
		ModelValue(val).
		Variant(VariantOutlined).
		HideDetails(true).
		Density(DensityCompact)
}

func (vc *ViewCommon) PasswordInput(
	id string,
	placeholder string,
	val string,
	canReveal bool,
) *VTextFieldBuilder {
	in := vc.Input(id, placeholder, val)
	if canReveal {
		varName := fmt.Sprintf(`show_%s`, id)
		in.Attr(":append-inner-icon", fmt.Sprintf(`vars.%s ? "mdi-eye-off" : "mdi-eye"`, varName)).
			Attr(":type", fmt.Sprintf(`vars.%s ? "text" : "password"`, varName)).
			Attr("@click:append-inner", fmt.Sprintf(`vars.%s = !vars.%s`, varName, varName)).
			Attr(web.VAssign("vars", fmt.Sprintf(`{%s: false}`, varName))...)
	}

	return in
}

// need to import zxcvbn.js
func (vc *ViewCommon) PasswordInputWithStrengthMeter(in *VTextFieldBuilder, id string, val string) HTMLComponent {
	in.Attr("v-model", fmt.Sprintf(`form.%s`, id))
	return Components(
		in,
		Div(
			web.Scope(
				VProgressLinear().
					Class("mt-2").
					Attr(":model-value", fmt.Sprintf(`(vars.meter_score?vars.meter_score(form.%s):0) * 20`, id)).
					// TODO reset color
					Attr(":color", fmt.Sprintf(`["secondary", "error-darken-1", "error", "warning", "warning-lighten-1", "success"][(vars.meter_score?vars.meter_score(form.%s):0)]`, id)),
			).Slot("{ locals }").
				LocalsInit(fmt.Sprintf(`{ meter_score:  0 }`)),
		).ID(fmt.Sprintf("password_%s", id)).
			Attr("v-show", fmt.Sprintf("!!form.%s", id)),
	)
}

func (vc *ViewCommon) FormSubmitBtn(
	label string,
) *VBtnBuilder {
	return VBtn(label).
		Color("primary").
		Block(true).
		Size(SizeLarge).
		Attr("type", "submit").
		Class("mt-6")
}

// requirements:
// - submit button
//   - add class `g-recaptcha`
//   - add attr `data-sitekey=<key>`
//   - add attr `data-callback=onSubmit`
//
// - add token field like `Input("token").Id("token").Type("hidden")`
func (vc *ViewCommon) InjectRecaptchaAssets(ctx *web.EventContext, formID string, tokenFieldID string) {
	ctx.Injector.HeadHTML(`
<style>
.grecaptcha-badge { visibility: hidden; }
</style>
    `)
	ctx.Injector.HeadHTML(fmt.Sprintf(`
<script>
function onSubmit(token) {
	document.getElementById("%s").value = token;
	document.getElementById("%s").submit();
}
</script>
    `, tokenFieldID, formID))
	ctx.Injector.TailHTML(`
<script src="https://www.google.com/recaptcha/api.js"></script>
    `)
}

func (vc *ViewCommon) InjectZxcvbn(ctx *web.EventContext) {
	ctx.Injector.HeadHTML(fmt.Sprintf(`
<script src="%s"></script>
    `, login.ZxcvbnJSURL))
}
