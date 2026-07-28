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

// loginDialog is the dialog itself: the login page in a frame, and a listener
// that closes it when that page reports the session is back.
func loginDialog(loginURL, doneURL string) h.HTMLComponent {
	// The frame reports back with postMessage; the message carries the URL it
	// landed on, so a page other than ours cannot close the dialog.
	listener := fmt.Sprintf(`
(scope) => {
	const done = %s, onMessage = (e) => {
		if (e.origin !== window.location.origin || e.data !== "rvq:login-done") { return }
		window.removeEventListener("message", onMessage);
		%s = false;
	};
	window.addEventListener("message", onMessage);
}`, h.JSONString(doneURL), presets.LoginDialogVar)

	return web.Scope(
		vx.VXDialog(
			h.Tag("iframe").
				Attr("src", loginURL).
				Attr("title", "login").
				Attr("style", "width:100%;height:70vh;border:none"),
		).
			Attr("v-model", presets.LoginDialogVar).
			Attr("width", "520").
			// no way out but logging in: the page underneath is waiting for the
			// session, and closing this would only lead to another redirect
			Attr("persistent", true),
		web.RunScript(strings.TrimSpace(listener)),
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
