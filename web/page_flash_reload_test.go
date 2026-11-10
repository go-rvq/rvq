package web_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	. "github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/multipartestutils"
)

// The flash an event leaves in the cookie shows in the page that the reload
// of its push state renders: a request of its own, an event too.
func TestFlashCookieInReload(t *testing.T) {
	p := New().Page(func(ctx *EventContext) (pr PageResponse, err error) {
		pr.Body = h.Div(h.Text(fmt.Sprintf("flash=%v", ctx.Flash)))
		ctx.Flash = nil
		return
	}).EventFunc("say", func(ctx *EventContext) (r EventResponse, err error) {
		ctx.Flash = "done"
		RespondFlashCookie(FlashCookieName, "/", nil, ctx, &PageResponse{})
		r.PushState = Location(nil)
		return
	})

	w := httptest.NewRecorder()
	p.ServeHTTP(w, multipartestutils.NewMultipartBuilder().EventFunc("say").BuildEventFuncRequest())
	cookies := w.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Name != FlashCookieName {
		t.Fatalf("the event set no flash cookie: %v", cookies)
	}

	r := multipartestutils.NewMultipartBuilder().EventFunc("__reload__").BuildEventFuncRequest()
	r.AddCookie(cookies[0])
	w = httptest.NewRecorder()
	p.ServeHTTP(w, r)
	if body := w.Body.String(); !strings.Contains(body, "done") {
		t.Fatalf("the reload did not show the flash: %s", body)
	}
	// read, the cookie is removed
	if c := w.Result().Cookies(); len(c) != 1 || c[0].MaxAge >= 0 {
		t.Fatalf("the reload kept the flash cookie: %v", c)
	}
}

// A page that does not show the flash it read does not write it again: the
// cookie, written by another (the admin's messages have more fields), was
// written back and came again in every page.
func TestFlashCookieNotWrittenBack(t *testing.T) {
	p := New().Page(func(ctx *EventContext) (pr PageResponse, err error) {
		pr.Body = h.Div(h.Text("a page that shows no flash"))
		return
	})
	raw := `[{"text":"done","htmlText":"","detail":null,"more":1}]` + "\r/"
	for _, event := range []bool{false, true} {
		var r *http.Request
		if event {
			r = multipartestutils.NewMultipartBuilder().EventFunc("__reload__").BuildEventFuncRequest()
		} else {
			r = httptest.NewRequest("GET", "/", nil)
		}
		r.AddCookie(&http.Cookie{Name: FlashCookieName, Value: base64.StdEncoding.EncodeToString([]byte(raw))})
		w := httptest.NewRecorder()
		p.ServeHTTP(w, r)
		for _, c := range w.Result().Cookies() {
			if c.Name == FlashCookieName && c.MaxAge >= 0 {
				t.Errorf("event %v: the flash was written back: %+v", event, c)
			}
		}
	}
}
