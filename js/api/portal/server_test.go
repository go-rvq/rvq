package portal

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/presets"
)

// forms is a source of the test: a content by its id, an action that echoes
// what it got (the body as the script sent it).
type forms struct{}

func (forms) Content(r *http.Request, key Key) (*Content, error) {
	if key.ID == "none" {
		return nil, nil
	}
	return &Content{Template: "<p>{{ name }}</p>", Data: map[string]string{"name": key.ID + "/" + key.Locale}}, nil
}

func (forms) ServeAction(w http.ResponseWriter, r *http.Request, key Key, action string) {
	_ = r.ParseMultipartForm(1 << 20)
	f, _, _ := r.FormFile("photo")
	var photo []byte
	if f != nil {
		photo, _ = io.ReadAll(f)
	}
	_ = json.NewEncoder(w).Encode(map[string]string{
		"action": action, "id": key.ID, "name": r.FormValue("name"), "photo": string(photo),
	})
}

func newServer() *Server {
	return NewServer(presets.NewHMACFormSigner([]byte("test-key"))).
		Register("form", forms{}).
		Register("plain", SourceFunc(func(r *http.Request, key Key) (*Content, error) {
			return &Content{Template: "<b>plain</b>"}, nil
		}))
}

func post(s *Server, token, action string, body io.Reader, contentType string) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodPost, DefaultPath, body)
	if token != "" {
		r.Header.Set(HeaderToken, token)
	}
	if action != "" {
		r.Header.Set(HeaderAction, action)
	}
	if contentType != "" {
		r.Header.Set("Content-Type", contentType)
	}
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	return w
}

// A token names its content, opaque — no address in it —, and comes back
// to the same key; tampered or past its time, it does not.
func TestToken(t *testing.T) {
	s := newServer()
	tok := s.Token("form", "quote", "pt-BR")
	if strings.Contains(tok, "/") || strings.Contains(tok, "quote") {
		t.Errorf("an address or the id in the token: %s", tok)
	}
	k, err := s.Key(tok)
	if err != nil || k.Kind != "form" || k.ID != "quote" || k.Locale != "pt-BR" {
		t.Fatalf("key %+v %v", k, err)
	}
	if _, err := s.Key(tok[:len(tok)-2] + "xx"); err != ErrBadToken {
		t.Errorf("tampered: %v", err)
	}
	if _, err := s.Key("nonsense"); err != ErrBadToken {
		t.Errorf("not a token: %v", err)
	}
	other := NewServer(presets.NewHMACFormSigner([]byte("another-key")))
	if _, err := other.Key(tok); err != ErrBadToken {
		t.Errorf("another server's: %v", err)
	}
	s.TTL(-time.Minute)
	if _, err := s.Key(s.Token("form", "quote", "")); err != ErrExpired {
		t.Errorf("expired: %v", err)
	}
}

// The server answers the content of a token, an action of its source with
// the body as sent (a multipart of files too), and an error code the
// runtime tells apart.
func TestServer(t *testing.T) {
	s := newServer()

	w := post(s, s.Token("form", "quote", "en"), "", nil, "")
	var c Content
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &c) != nil || c.Template != "<p>{{ name }}</p>" {
		t.Fatalf("content: %d %s", w.Code, w.Body)
	}
	if d, _ := c.Data.(map[string]any); d["name"] != "quote/en" {
		t.Errorf("its data: %+v", c.Data)
	}
	if w.Header().Get("Cache-Control") != "no-store" {
		t.Error("the content cached")
	}

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	_ = mw.WriteField("name", "Ana")
	fw, _ := mw.CreateFormFile("photo", "a.jpg")
	_, _ = fw.Write([]byte("JPEG"))
	_ = mw.Close()
	w = post(s, s.Token("form", "quote", ""), "submit", &body, mw.FormDataContentType())
	var got map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &got)
	if got["action"] != "submit" || got["id"] != "quote" || got["name"] != "Ana" || got["photo"] != "JPEG" {
		t.Errorf("action: %d %s", w.Code, w.Body)
	}

	for _, c := range []struct {
		name, token, action string
		code                int
		err                 string
	}{
		{"no token", "", "", 403, "token"},
		{"bad token", "x.y", "", 403, "token"},
		{"unknown kind", s.Token("other", "x", ""), "", 404, "kind"},
		{"no content", s.Token("form", "none", ""), "", 404, "content"},
		{"no actions", s.Token("plain", "x", ""), "submit", 404, "action"},
		{"expired", NewServer(presets.NewHMACFormSigner([]byte("test-key"))).TTL(-time.Minute).Token("form", "x", ""), "", 410, "expired"},
	} {
		w := post(s, c.token, c.action, nil, "")
		var e map[string]string
		_ = json.Unmarshal(w.Body.Bytes(), &e)
		if w.Code != c.code || e["error"] != c.err {
			t.Errorf("%s: %d %s, want %d %s", c.name, w.Code, w.Body, c.code, c.err)
		}
	}

	r := httptest.NewRequest(http.MethodGet, DefaultPath, nil)
	rw := httptest.NewRecorder()
	s.ServeHTTP(rw, r)
	if rw.Code != 405 {
		t.Errorf("GET: %d", rw.Code)
	}
}

// The bundle is served with the server's path in it, and cached by its
// ETag; its address for a page carries its version, the ETag's.
func TestScript(t *testing.T) {
	s := newServer()
	s.Path = "/x/portal"
	mux := http.NewServeMux()
	s.Mount(mux)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/x/portal.js", nil))
	body := w.Body.String()
	if w.Code != 200 || !strings.Contains(body, `"/x/portal"`) || strings.Contains(body, `"/_rvq/portal"`) {
		t.Fatalf("bundle: %d, the path written: %v", w.Code, strings.Contains(body, `"/x/portal"`))
	}
	w2 := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/x/portal.js", nil)
	r.Header.Set("If-None-Match", w.Header().Get("ETag"))
	mux.ServeHTTP(w2, r)
	if w2.Code != http.StatusNotModified {
		t.Errorf("cached: %d", w2.Code)
	}
	if u := s.ScriptURL(); u != "/x/portal.js?v="+strings.Trim(w.Header().Get("ETag"), `"`) {
		t.Errorf("the address of the bundle: %s", u)
	}
	w3 := httptest.NewRecorder()
	mux.ServeHTTP(w3, httptest.NewRequest(http.MethodPost, "/x/portal", nil))
	if w3.Code != http.StatusForbidden {
		t.Errorf("the server at its path: %d", w3.Code)
	}
}
