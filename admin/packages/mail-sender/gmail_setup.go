package mail_sender

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"golang.org/x/oauth2"
)

const gmailConsoleCredentialsURL = "https://console.cloud.google.com/apis/credentials"

// gmailSetupComponent renders the help button that opens the dialog explaining
// how to create the Google OAuth client. There is no API that creates the
// client, so this walks the user through the Console by hand.
func gmailSetupComponent(msgr *Messages, callback string) h.HTMLComponent {
	steps := make(h.HTMLComponents, len(msgr.GmailSenderSetupSteps))
	for i, step := range msgr.GmailSenderSetupSteps {
		steps[i] = h.Li(h.Text(step)).Class("mb-2")
	}

	// Rendered from GmailScopes itself, so the dialog can never list a scope the
	// sender does not ask for.
	scopes := make(h.HTMLComponents, len(GmailScopes))
	for i, scope := range GmailScopes {
		scopes[i] = h.Div(h.Code(scope)).Class("mb-1")
	}

	scopeSteps := make(h.HTMLComponents, len(msgr.GmailSenderSetupScopesSteps))
	for i, step := range msgr.GmailSenderSetupScopesSteps {
		scopeSteps[i] = h.Li(h.Text(step)).Class("mb-2")
	}

	return web.Scope(
		v.VBtn(msgr.GmailSenderSetup).
			PrependIcon("mdi-help-circle-outline").
			Variant(v.VariantText).
			Size(v.SizeSmall).
			Color("info").
			Density(v.DensityCompact).
			Class("mb-2 px-1").
			Attr("@click", "locals.setupOpen = true"),
		vx.VXDialog().
			Title(msgr.GmailSenderSetup).
			SlotBody(
				h.Div(
					h.Div(h.Text(msgr.GmailSenderSetupTitle)).Class("mb-4"),
					h.Ol(steps...).Class("ms-4"),
					h.Div(
						h.Div(h.Strong(msgr.GmailSenderSetupScopes)).Class("mb-2"),
						h.Ol(scopeSteps...).Class("ms-4 mb-3"),
						h.Div(scopes...).Class("ps-2"),
						h.Div(h.Text(msgr.GmailSenderSetupScopesHint)).Class("text-caption mt-1"),
					).Class("mt-4"),
					v.VAlert(h.Text(msgr.GmailSenderSetupBlocked)).
						Type("warning").
						Variant(v.VariantTonal).
						Density(v.DensityCompact).
						Class("mt-4"),
					h.Div(
						h.Strong(msgr.GmailSenderSetupCallback+": "),
						h.Code(callback),
					).Class("mt-4"),
					h.Div(h.Text(msgr.GmailSenderSetupCallbackHint)).
						Class("text-caption mt-1"),
					v.VBtn(msgr.GmailSenderSetupConsoleBtn).
						Variant(v.VariantText).
						Color("info").
						Class("mt-4 px-1").
						PrependIcon("mdi-open-in-new").
						Attr("href", gmailConsoleCredentialsURL).
						Attr("target", "_blank"),
				).Class("pa-4 text-body-2"),
			).
			Closable(true).
			Width("680").
			VModel("locals.setupOpen"),
	).LocalsInit("{setupOpen: false}").Slot("{ locals }")
}

// gmailAuthCodeURL builds the authorization URL for the sign in.
//
// ApprovalForce (prompt=consent) is not optional here. Without it Google skips
// the consent screen for an account that already authorized this client, and
// reissues from the remembered grant: a scope added later never reaches the
// token, and the granular consent checkboxes are never shown to be ticked.
// It also governs the refresh token, which Google returns only when the
// consent screen is actually displayed — re-authorizing without it yields a
// token that cannot be refreshed once the access token expires.
func gmailAuthCodeURL(config *oauth2.Config) string {
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}
