package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	xlogin "github.com/go-rvq/rvq/x/login"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// Logging in again without losing the page.
//
// A session dies while somebody has a form open. The next request — a save, a
// list reload — comes back from the middleware as a redirect to the login page,
// and the browser replaces everything: the form, and whatever was typed in it.
//
// So a request made BY THE PAGE (plaid — see web.PlaidRequestHeader) is answered
// differently. It gets 401, which is what the server log should show anyway, and
// a single instruction: web.EventResponse.LoginURI, the address of the dialog.
// Nothing on the page is touched.
//
// plaid() then asks for that URI, mounts what comes back in the layout's login
// portal (presets.LoginPortalName) — the login page in a dialog over everything
// — and hands the way back in its scope: onLoginSuccess. Logging in closes the
// dialog, calls it, and the request the user originally made runs again, so the
// click that hit the dead session finally goes through.
//
// The dialog holds the real login page in a frame, so everything that page does
// keeps working (the form protection, reCAPTCHA when configured, OAuth buttons,
// the error messages). The login lands on LoginDoneURI, a page whose whole job
// is to say "done" to the page around it.

// installLoginDialog wires the answer above into the login middleware.
func (b *Builder) installLoginDialog(pb *presets.Builder) {
	lb := b.Builder()
	uriPrefix := pb.GetURIPrefix()
	doneURL := uriPrefix + presets.LoginDoneURI
	dialogURL := uriPrefix + presets.LoginDialogURI

	pb.MuxSetup(func(prefix string, mux *http.ServeMux) {
		// the page the frame lands on when the login succeeds
		mux.Handle(prefix+presets.LoginDoneURI, LoginDonePage())

		// The dialog itself. EnsureLanguage because this handler is not a page —
		// nothing else would put the translations in its context.
		mux.Handle(prefix+presets.LoginDialogURI, pb.I18n().EnsureLanguage(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// where the login goes when it succeeds
			lb.SetContinueURL(w, doneURL)

			loginURL := lb.GetLoginPageURL()
			msgr := xlogin.GetMessages(r.Context())

			var res web.EventResponse
			res.UpdatePortal(presets.LoginPortalName, loginDialog(loginURL, msgr.LoginAgainTitle))
			web.AppendRunScripts(&res, loginDialogScript(loginURL))

			writeEventResponse(w, r, &res)
		})))
	})

	// Both are asked for exactly when there is no session — that is the whole
	// point of them — so nothing may turn them away.
	lb.WhiteList(dialogURL, doneURL)

	lb.SetUnauthorizedResponder(func(w http.ResponseWriter, r *http.Request) bool {
		// Only for the page's own requests. A plain navigation still goes to the
		// login page: there is no page to keep standing.
		if !web.IsPlaidRequest(r) {
			return false
		}

		res := web.EventResponse{LoginURI: dialogURL}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(&res); err != nil {
			panic(err)
		}
		return true
	})
}

// HiddenByLoginClass marks the overlays put out of the way while the login
// dialog is up. They are only made INVISIBLE: still mounted, still holding
// everything the user had typed, and back exactly as they were once the session
// returns.
const HiddenByLoginClass = "rvq-hidden-by-login"

// loginDialogScript runs in the scope of the request that asked for the dialog,
// which is where onLoginSuccess lives (see js/corejs/src/builder.ts, login()).
// It puts the dialogs already open out of sight, opens this one, and waits for
// the frame to report that the session is back — then undoes all of it and lets
// the interrupted request go through.
//
// The message comes from our own origin and says exactly one thing, so no other
// page can close this dialog.
func loginDialogScript(loginURL string) string {
	return fmt.Sprintf(strings.TrimSpace(`
(function () {
	const loginURL = %s;

	// Dialogs already open — the form the user was filling in — step aside so
	// only the login is on screen. Anything hidden by an earlier round stays as
	// it is, and so does anything already holding the login frame.
	const hidden = [];
	document.querySelectorAll(".v-overlay--active").forEach((el) => {
		if (el.classList.contains(%s)) { return }
		if (el.querySelector('iframe[src="' + loginURL + '"]')) { return }
		el.classList.add(%s);
		hidden.push(el);
	});

	const onMessage = (e) => {
		if (e.origin !== window.location.origin || e.data !== "rvq:login-done") { return }
		window.removeEventListener("message", onMessage);
		hidden.forEach((el) => el.classList.remove(%s));
		%s = false;
		// back to what the user was doing
		if (typeof onLoginSuccess === "function") { onLoginSuccess() }
	};

	window.addEventListener("message", onMessage);
	%s = true;
})()`),
		h.JSONString(loginURL),
		h.JSONString(HiddenByLoginClass),
		h.JSONString(HiddenByLoginClass),
		h.JSONString(HiddenByLoginClass),
		presets.LoginDialogVar,
		presets.LoginDialogVar)
}

// loginDialog is the dialog itself: the login page in a frame, under a title
// that says why it showed up — the session ended, sign in again.
func loginDialog(loginURL, title string) h.HTMLComponent {
	return h.Components(
		// invisible, not removed: `display:none` would drop the size the dialog
		// had, and the form inside it would come back measured from scratch
		h.Style("."+HiddenByLoginClass+" { visibility: hidden !important; pointer-events: none !important; }"),

		vx.VXDialog(
			h.Tag("iframe").
				Attr("src", loginURL).
				Attr("title", "login").
				// takes whatever room the dialog has, so expanding it gives the
				// login page the whole window
				Attr("style", "width:100%;height:100%;min-height:60vh;border:none"),
		).
			Attr("v-model", presets.LoginDialogVar).
			Title(title).
			Width("520").
			// the login page of an application can be tall (logo, OAuth buttons,
			// a captcha): expanding gives it the whole window
			Expandable(true).
			// no way out but logging in: the page underneath is waiting for the
			// session, and closing this would only lead to another redirect
			Attr("persistent", true),
	)
}

// writeEventResponse answers a plain event response — the components rendered,
// then the JSON the client expects.
func writeEventResponse(w http.ResponseWriter, r *http.Request, res *web.EventResponse) {
	ctx := r.Context()

	res.Body = h.RawHTML(h.MustString(res.Body, ctx))
	for _, up := range res.UpdatePortals {
		up.Body = h.RawHTML(h.MustString(up.Body, ctx))
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

// LoginDonePage is what the frame lands on after a successful login: it tells
// the page around it and shows nothing else.
func LoginDonePage() http.HandlerFunc {
	body := `<!doctype html><meta charset="utf8"><title>ok</title>` +
		`<script>try{parent.postMessage("rvq:login-done",window.location.origin)}catch(e){}</script>`
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write([]byte(body))
	}
}
