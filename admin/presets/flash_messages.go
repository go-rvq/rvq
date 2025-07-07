package presets

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/js"
	"github.com/go-rvq/rvq/web/vue"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"

	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

type FlashMessage = web.FlashMessage

type FlashMessages = web.FlashMessages

func NewFlashMessages(msg any, color ...string) FlashMessages {
	switch t := msg.(type) {
	case FlashMessages:
		return t
	case []*FlashMessage:
		return t
	case *FlashMessage:
		return FlashMessages{t}
	default:
		return FlashMessages{NewFlashMessage(msg, color...)}
	}
}

func NewFlashMessage(msg any, color ...string) (m *FlashMessage) {
	if m, _ = msg.(*FlashMessage); m != nil {
		return
	}
	m = &FlashMessage{}
	for _, m.Color = range color {
	}

	switch t := msg.(type) {
	case *web.ValidationErrors:
		gErr := t.GetGlobalError()
		if len(gErr) > 0 {
			m.Text = gErr
			if m.Color == "" {
				m.Color = "error"
			}
		}
	case error:
		if m.Color == "" {
			m.Color = "error"
		}
		m.Text = t.Error()
	case string:
		m.Text = t
	case h.RawHTML:
		m.HtmlText = string(t)
	case h.HTMLComponent:
		m.Text = "Bad flash message type. Contact your administrator."
		m.Color = "error"
	case FlashMessage:
		*m = t
	default:
		m.Text = fmt.Sprintf("%v", t)
	}
	return
}

func GetFlashMessages(ctx *web.EventContext) (v FlashMessages) {
	v, _ = ctx.ContextValue(ctxFlashMessages).(FlashMessages)
	return v
}

func AddFlashMessage(ctx *web.EventContext, message ...*FlashMessage) {
	if old := web.GetContextValuer(ctx.R.Context(), ctxFlashMessages); old != nil {
		old.Set(append(old.Get().(FlashMessages), message...))
	} else {
		ctx.WithContextValue(ctxFlashMessages, FlashMessages(message))
	}
}

func RespondFlash(r *web.EventResponse, msg any, color ...string) {
	comp := RenderFlash(msg, color...)
	if comp != nil {
		r.UpdatePortal(FlashPortalName, comp)
	}
}

func RenderFlash(msg any, color ...string) h.HTMLComponent {
	if msg == nil {
		return nil
	}

	switch t := msg.(type) {
	case []*FlashMessage:
		return RenderFlash(FlashMessages(t))
	case FlashMessages:
		if len(t) == 0 {
			return nil
		}
		var comps h.HTMLComponents
		for _, msg := range t {
			comps = append(comps, RenderFlash(msg))
		}
		return comps
	}

	var (
		m       = web.NewFlashMessage(msg, color...)
		comps   h.HTMLComponents
		actions h.HTMLComponents
		flash   = js.Object{"show": true}
	)

	if m.Text == "" {
		if len(m.HtmlText) == 0 {
			return nil
		}

	}

	if m.Color == "" {
		m.Color = "success"
	}

	snack := VSnackbar().Location("top").
		Color(m.Color).
		ZIndex(2000000).
		Attr("v-model", "flash.show")

	if m.Duration == 0 {
		m.Duration = 5
		if m.Color == "error" {
			m.Duration *= 2
		}
	}

	snack.Timeout(m.Duration * 1000)

	comps = append(comps, snack)

	if m.Detail != nil {
		flash["detail"] = false

		if d, _ := web.Unscoped(m.Detail).(*vx.VXDialogBuilder); d != nil {
			d.
				Expandable(true).
				Closable(true).
				VModel("flash.detail").
				Density(DensityCompact)
			comps = append(comps, m.Detail)
		} else {
			comps = append(comps, vx.VXDialog().
				SlotBody(m.Detail).
				VModel("flash.detail").
				Expandable(true).
				Closable(true).
				Width("400").
				Density(DensityCompact))
		}

		actions = append(actions, VBtn("").Variant("text").
			Attr("@click", "flash.detail = true").
			Children(VIcon("mdi-information")))
	}

	if !m.NotClosable {
		actions = append(actions, VBtn("").Variant("text").
			Attr("@click", "flash.show = false").
			Children(VIcon("mdi-close")))
	}

	if len(actions) > 0 {
		actions = h.HTMLComponents{web.Slot(actions...).Name("actions")}
	}

	var text h.HTMLComponent
	if len(m.Text) > 0 {
		text = h.Text(m.Text)
	} else {
		text = h.RawHTML(m.HtmlText)
	}

	snack.TagBuilder.Children(text, actions)
	return vue.UserComponent(append(comps, snack)...).Scope("flash", flash)
}
