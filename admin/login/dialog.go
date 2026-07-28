package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// Logging in again without losing the page.
//
// A session dies while somebody has a form open. The next request — a save, a
// list reload — comes back from the middleware as a redirect to the login page,
// and the browser replaces everything: the form, and whatever was typed in it.
//
// So an EVENT request whose session is gone is answered differently: the page
// stays exactly as it is and the login opens in a DIALOG over it, in a portal of
// its own at the layout root (presets.LoginPortalName). Logging in there closes
// the dialog and leaves the page — and the form — untouched, ready for the user
// to press save again.
//
// The dialog holds the real login page in a frame, so everything that page does
// keeps working (the form protection, reCAPTCHA when configured, OAuth buttons,
// the error messages). The login lands on LoginDoneURI, a page whose whole job
// is to say "done" to the page around it.

// installLoginDialog wires the answer above into the login middleware.
func (b *Builder) installLoginDialog(pb *presets.Builder) {
	lb := b.Builder()
	doneURL := pb.GetURIPrefix() + presets.LoginDoneURI

	// the page the frame lands on when the login succeeds
	pb.MuxSetup(func(prefix string, mux *http.ServeMux) {
		mux.Handle(prefix+presets.LoginDoneURI, LoginDonePage())
	})

	lb.SetUnauthorizedResponder(func(w http.ResponseWriter, r *http.Request) bool {
		// Only for the page's own requests. A plain navigation still goes to the
		// login page: there is no page to keep standing.
		if r.FormValue(web.EventFuncIDName) == "" {
			return false
		}

		// where the login goes when it succeeds
		lb.SetContinueURL(w, doneURL)

		var res web.EventResponse
		res.UpdatePortal(presets.LoginPortalName, loginDialog(lb.GetLoginPageURL(), doneURL))
		web.AppendRunScripts(&res, presets.LoginDialogVar+" = true")

		writeEventResponse(w, r, &res)
		return true
	})
}

// HiddenByLoginClass marks the overlays put out of the way while the login
// dialog is up. They are only made INVISIBLE: still mounted, still holding
// everything the user had typed, and back exactly as they were once the session
// returns.
const HiddenByLoginClass = "rvq-hidden-by-login"

// loginDialog is the dialog itself: the login page in a frame, the overlays
// that were already open put out of sight, and a listener that undoes both when
// that page reports the session is back.
func loginDialog(loginURL, doneURL string) h.HTMLComponent {
	// The message comes from our own origin and says exactly one thing, so no
	// other page can close this dialog.
	script := fmt.Sprintf(`
(scope) => {
	const loginURL = %s;

	// Dialogs already open — the form the user was filling in — step aside so
	// only the login is on screen. Anything holding the login frame is skipped,
	// and so is anything hidden by an earlier round.
	const hidden = [];
	document.querySelectorAll(".v-overlay--active").forEach((el) => {
		if (el.classList.contains(%s)) { return }
		if (el.querySelector('iframe[src="' + loginURL + '"]')) { return }
		el.classList.add(%s);
		hidden.push(el);
	});

	const restore = () => hidden.forEach((el) => el.classList.remove(%s));

	const onMessage = (e) => {
		if (e.origin !== window.location.origin || e.data !== "rvq:login-done") { return }
		window.removeEventListener("message", onMessage);
		restore();
		%s = false;
	};

	window.addEventListener("message", onMessage);
}`,
		h.JSONString(loginURL),
		h.JSONString(HiddenByLoginClass),
		h.JSONString(HiddenByLoginClass),
		h.JSONString(HiddenByLoginClass),
		presets.LoginDialogVar)

	return web.Scope(
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
			Width("520").
			// the login page of an application can be tall (logo, OAuth buttons,
			// a captcha): expanding gives it the whole window
			Expandable(true).
			// no way out but logging in: the page underneath is waiting for the
			// session, and closing this would only lead to another redirect
			Attr("persistent", true),

		web.RunScript(strings.TrimSpace(script)),
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
