package presets

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// A close the user asks for — ESC, a click outside, the close button — is a
// request (closer.requestClose): a form with changes not saved holds it and
// asks first (FormGuard). A close by the app — `closer.show = false`, once
// saved — is not one.

// CloserRequestClose is the expression of a close the user asks for, in the
// scope of a closer: guarded (closer.requestClose), else a plain close (a
// closer from before the guards).
const CloserRequestClose = `(closer.requestClose ? closer.requestClose() : (closer.show = false))`

// CloserModelUpdate is the handler of the update of an overlay's model bound
// to the closer (`:model-value="closer.show"`): opening opens; closing is a
// close the user asked for (CloserRequestClose).
const CloserModelUpdate = `(v) => v ? (closer.show = true) : ` + CloserRequestClose

// FormGuard is what a form in an overlay adds to it: the guard of its
// content (go-plaid-form-guard: the values of `form` against what they were
// when it opened) and the confirmation of a close it holds — keep editing,
// or discard the changes.
func FormGuard(ctx *web.EventContext) h.HTMLComponent {
	msgr := MustGetMessages(ctx.Context())
	return h.Components(
		h.Tag("go-plaid-form-guard").Attr(":form", "form"),
		v.VDialog(
			v.VCard(
				v.VCardTitle(h.Text(msgr.DiscardChangesTitle)),
				v.VCardText(h.Text(msgr.DiscardChangesText)),
				v.VCardActions(
					v.VSpacer(),
					v.VBtn(msgr.KeepEditing).Variant(v.VariantText).
						Attr("@click", "closer.confirmClose = false").
						Attr("data-form-guard-keep", "true"),
					v.VBtn(msgr.DiscardChanges).Variant(v.VariantFlat).Color("error").
						Attr("@click", "closer.show = false").
						Attr("data-form-guard-discard", "true"),
				),
			).Attr("data-form-guard-confirm", "true"),
		).Attr("v-model", "closer.confirmClose").MaxWidth("440").Persistent(true),
	)
}
