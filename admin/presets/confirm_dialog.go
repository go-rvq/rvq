package presets

import (
	"fmt"
	"net/http"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/js"
	"github.com/go-rvq/rvq/web/vue"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	h "github.com/theplant/htmlgo"
)

const (
	EventOpenConfirmDialog = "presets_ConfirmDialog"
)

func (b *Builder) openConfirmDialog(ctx *web.EventContext) (er web.EventResponse, err error) {
	var (
		msgr = MustGetMessages(ctx.Context())
		cb   = &ConfirmDialogOpenerBuilder{
			headerColor: ColorWarning,
			alertType:   TypeWarning,
			okText:      msgr.OK,
			okColor:     ColorError,
		}
	)

	cb.ParseRequest(ctx.R)

	if confirmEvent := ctx.R.FormValue(ConfirmDialogConfirmEvent); confirmEvent != "" {
		cb.handler = js.Raw(confirmEvent)
	}

	if cb.handler == "" {
		ShowMessage(&er, "confirm event is empty", "error")
		return
	}

	if cb.prompt == "" {
		if v := ctx.R.FormValue(ConfirmDialogPromptText); v != "" {
			cb.prompt = h.RawHTML(v)
		} else {
			cb.prompt = h.RawHTML(msgr.ConfirmDialogPromptText)
		}
	}

	if cb.portal == "" {
		if v := ctx.R.FormValue(ConfirmDialogDialogPortalName); v != "" {
			cb.portal = v
		} else if v := ctx.R.FormValue(ParamTargetPortal); v != "" {
			cb.portal = v
		} else {
			cb.portal = DefaultConfirmDialogPortalName
		}
	}

	var prompt h.HTMLComponent = cb.prompt

	if len(cb.promptHandler) > 0 {
		prompt = web.Portal().LoaderString(string(cb.promptHandler))
	}

	b.Dialog().
		SetTargetPortal(cb.portal).
		Respond(ctx, &er, vue.UserComponent(
			vx.VXDialog().
				Density(DensityCompact).
				Style("max-width: 600px").
				SlotBody(
					VAlert(
						prompt,
					).Type(cb.alertType).
						Variant(VariantTonal),
				).
				ToolbarProps(fmt.Sprintf(`{color:%q}`, cb.headerColor)).
				ScopedSlotBottom("{isActive}", h.HTMLComponents{
					VSpacer(),
					VBtn(msgr.Cancel).
						Variant(VariantFlat).
						Class("ml-2").
						On("click", "isActive.value = false"),
					VBtn(cb.okText).
						Color(cb.okColor).
						Variant(VariantFlat).
						Theme(ThemeDark).
						Attr("@click", fmt.Sprintf("%s; isActive.value = false", cb.handler)),
				}).
				Title(msgr.ConfirmDialogPromptTitle),
		).ScopeVar("loading", "{value: false}"))
	return
}

type ConfirmDialogOpenerBuilder struct {
	alertType     string
	headerColor   string
	title         string
	prompt        h.RawHTML
	okText        string
	handler       js.Raw
	portal        string
	promptHandler js.Raw
	okColor       string
}

func (b *ConfirmDialogOpenerBuilder) AlertType(v string) *ConfirmDialogOpenerBuilder {
	b.alertType = v
	return b
}

func (b *ConfirmDialogOpenerBuilder) HeaderColor(color string) *ConfirmDialogOpenerBuilder {
	b.headerColor = color
	return b
}

func (b *ConfirmDialogOpenerBuilder) Title(v string) *ConfirmDialogOpenerBuilder {
	b.title = v
	return b
}

func (b *ConfirmDialogOpenerBuilder) Prompt(v h.RawHTML) *ConfirmDialogOpenerBuilder {
	b.prompt = v
	return b
}

func (b *ConfirmDialogOpenerBuilder) OkText(v string) *ConfirmDialogOpenerBuilder {
	b.okText = v
	return b
}

func (b *ConfirmDialogOpenerBuilder) Handler(v js.Raw) *ConfirmDialogOpenerBuilder {
	b.handler = v
	return b
}

func (b *ConfirmDialogOpenerBuilder) Portal(v string) *ConfirmDialogOpenerBuilder {
	b.portal = v
	return b
}

func (b *ConfirmDialogOpenerBuilder) PromptHandler(v js.Raw) *ConfirmDialogOpenerBuilder {
	b.promptHandler = v
	return b
}

func (b *ConfirmDialogOpenerBuilder) OkColor(v string) *ConfirmDialogOpenerBuilder {
	b.okColor = v
	return b
}

func (b *ConfirmDialogOpenerBuilder) ParseRequest(r *http.Request) {
	if v := r.FormValue(ParamTargetPortal); v != "" {
		b.portal = v
	}
	if v := r.FormValue("alertType"); v != "" {
		b.alertType = v
	}
	if v := r.FormValue("headerColor"); v != "" {
		b.headerColor = v
	}
	if v := r.FormValue("title"); v != "" {
		b.title = v
	}
	if v := r.FormValue("prompt"); v != "" {
		b.prompt = h.RawHTML(v)
	}
	if v := r.FormValue("okText"); v != "" {
		b.okText = v
	}
	if v := r.FormValue("handler"); v != "" {
		b.handler = js.Raw(v)
	}
	if v := r.FormValue("promptHandler"); v != "" {
		b.promptHandler = js.Raw(v)
	}
	if v := r.FormValue("okColor"); v != "" {
		b.okColor = v
	}
}

func (b *ConfirmDialogOpenerBuilder) Build() *web.VueEventTagBuilder {
	return web.Plaid().
		EventFunc(EventOpenConfirmDialog).
		ValidQuery(ParamTargetPortal, b.portal).
		ValidQuery("headerColor", b.headerColor).
		ValidQuery("alertType", b.alertType).
		ValidQuery("title", b.title).
		ValidQuery("prompt", b.prompt).
		ValidQuery("okText", b.okText).
		ValidQuery("handler", b.handler).
		ValidQuery("promptHandler", b.promptHandler).
		ValidQuery("okColor", b.okColor)
}

func OpenConfirmDialog() *ConfirmDialogOpenerBuilder {
	return &ConfirmDialogOpenerBuilder{}
}
