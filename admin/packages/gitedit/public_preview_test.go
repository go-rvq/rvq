package gitedit

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	"github.com/go-rvq/rvq/web/multipartestutils"
)

// The public preview: ana opens it on her page — its link served with no
// login, the site as her draft makes it, not indexed —; bia cannot end it,
// root (!drafts) can; ended, made again or expired, the link is not found.
func TestPublicPreview(t *testing.T) {
	a := newShareApp(t)
	ctx := context.Background()
	pub := a.b.PublicPreviewHandler()
	get := func(path string) (int, string, string) {
		w := httptest.NewRecorder()
		pub.ServeHTTP(w, httptest.NewRequest("GET", path, nil)) // no user: no login
		return w.Code, w.Body.String(), w.Header().Get("X-Robots-Tag")
	}
	act := func(user, action string, form ...[2]string) *httptest.ResponseRecorder {
		mb := multipartestutils.NewMultipartBuilder().PageURL("/admin/site-files").EventFunc(actions.DoAction).
			Query(presets.ParamAction, action)
		for _, f := range form {
			mb = mb.AddField(f[0], f[1])
		}
		r := mb.BuildEventFuncRequest()
		r.Header.Set("X-User", user)
		w := httptest.NewRecorder()
		a.h.ServeHTTP(w, r)
		return w
	}

	a.get("ana", "/admin/site-files") // ana's draft made
	if _, body := a.get("ana", "/admin/site-files"); !strings.Contains(body, "data-public-preview=&#39;off&#39;") {
		t.Fatal("ana's page does not offer the public preview")
	}
	if code, _, _ := get("/_preview/nothing/x"); code != 404 {
		t.Errorf("a token never made: %d", code)
	}

	act("ana", ActionPublicPreview)
	p := a.b.PublicPreviewOf(ctx, "ana")
	if p == nil || len(p.Token) < 40 {
		t.Fatalf("not opened: %+v", p)
	}
	code, body, robots := get("/_preview/" + p.Token + "/en/about")
	if code != 200 || body != "the draft of ana at /_preview/"+p.Token || robots != "noindex, nofollow" {
		t.Errorf("served %d %q (robots %q)", code, body, robots)
	}
	if _, body := a.get("ana", "/admin/site-files"); !strings.Contains(body, "/_preview/"+p.Token) {
		t.Error("ana's page does not show the link")
	}

	// bia ends nothing of ana's; root does
	act("bia", ActionPublicPreviewOff, [2]string{"u", "ana"})
	if a.b.PublicPreviewOf(ctx, "ana") == nil {
		t.Error("bia ended ana's public preview")
	}
	// made again: the old link stops
	act("ana", ActionPublicPreview)
	p2 := a.b.PublicPreviewOf(ctx, "ana")
	if p2 == nil || p2.Token == p.Token {
		t.Fatal("made again with the same token")
	}
	if code, _, _ := get("/_preview/" + p.Token + "/"); code != 404 {
		t.Errorf("the old link still served: %d", code)
	}
	act("ana", ActionPublicPreviewOff)
	if a.b.PublicPreviewOf(ctx, "ana") != nil {
		t.Error("ana did not end it")
	}
	if code, _, _ := get("/_preview/" + p2.Token + "/"); code != 404 {
		t.Errorf("an ended link served: %d", code)
	}

	// until a date: expired, not found
	past := time.Now().Add(-time.Minute)
	p3, err := a.b.EnablePublicPreview(ctx, "ana", "ana", &past)
	if err != nil {
		t.Fatal(err)
	}
	if code, _, _ := get("/_preview/" + p3.Token + "/"); code != 404 {
		t.Errorf("an expired link served: %d", code)
	}
	// root (!drafts) ends another's
	p4, _ := a.b.EnablePublicPreview(ctx, "bia", "bia", nil)
	r := multipartestutils.NewMultipartBuilder().PageURL("/admin/site-files/u/bia").EventFunc(actions.DoAction).
		Query(presets.ParamAction, ActionPublicPreviewOff).BuildEventFuncRequest()
	r.Header.Set("X-User", "root")
	a.h.ServeHTTP(httptest.NewRecorder(), r)
	if code, _, _ := get("/_preview/" + p4.Token + "/"); code != 404 {
		t.Error("root did not end bia's")
	}
}
