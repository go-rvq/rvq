package cli

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"context"
	"net/http"

	"github.com/go-rvq/rvq/cli"

	"github.com/go-rvq/rvq/web"
)

func TestHttpApiHelp(t *testing.T) {
	cmd := HttpApiCommand(Config{})
	if cmd.Help == nil {
		t.Fatal("command has no Help handler")
	}
	var buf bytes.Buffer
	if err := cmd.Help(&cli.CommandContext{Err: &buf}); err != nil {
		t.Fatal(err)
	}
	s := buf.String()
	for _, want := range []string{"login", "method", "uri", "body", "file:", "expectedResponse", "keys"} {
		if !strings.Contains(s, want) {
			t.Errorf("help text missing %q", want)
		}
	}
}

func TestStripHTML(t *testing.T) {
	in := `<div class="v-snackbar"><div class="v-snackbar__content">Post &amp; salvo</div></div>`
	if got := stripHTML(in); got != "Post & salvo" {
		t.Fatalf("stripHTML = %q, want %q", got, "Post & salvo")
	}
}

func TestProcessResponseLiftsAndRemovesFlash(t *testing.T) {
	body := `{
		"updatePortals": [
			{"name": "flash", "body": "<div><span>Saved successfully</span></div>"},
			{"name": "presets_ListingDialog", "body": "<div>list</div>"}
		]
	}`
	res := processResponse(200, []byte(body))

	if len(res.Flash) != 1 || res.Flash[0] != "Saved successfully" {
		t.Fatalf("flash = %#v, want [\"Saved successfully\"]", res.Flash)
	}

	resp, ok := res.Response.(map[string]any)
	if !ok {
		t.Fatalf("response is %T, want map", res.Response)
	}
	ups, _ := resp["updatePortals"].([]any)
	if len(ups) != 1 {
		t.Fatalf("updatePortals kept %d, want 1 (flash removed)", len(ups))
	}
	if name := ups[0].(map[string]any)["name"]; name != "presets_ListingDialog" {
		t.Fatalf("kept portal %v, want presets_ListingDialog", name)
	}
}

func TestProcessResponseOnlyFlashDropsPortalsKey(t *testing.T) {
	body := `{"updatePortals":[{"name":"flash","body":"<p>Done</p>"}]}`
	res := processResponse(200, []byte(body))
	resp := res.Response.(map[string]any)
	if _, ok := resp["updatePortals"]; ok {
		t.Fatalf("updatePortals should be removed when only flash was present")
	}
	if len(res.Flash) != 1 || res.Flash[0] != "Done" {
		t.Fatalf("flash = %#v, want [\"Done\"]", res.Flash)
	}
}

func TestFlattenJSON(t *testing.T) {
	var decoded any
	json.Unmarshal([]byte(`{
		"Config": {"TypeID": 2, "Enabled": true},
		"Tags": ["a", "b"],
		"Title": "Hello"
	}`), &decoded)

	values := url.Values{}
	var files []formFile
	flattenJSON("", decoded, values, &files)

	want := map[string]string{
		"Config.TypeID":  "2",
		"Config.Enabled": "true",
		"Tags[0]":        "a",
		"Tags[1]":        "b",
		"Title":          "Hello",
	}
	for k, v := range want {
		if got := values.Get(k); got != v {
			t.Errorf("field %q = %q, want %q", k, got, v)
		}
	}
	if len(files) != 0 {
		t.Errorf("unexpected files: %#v", files)
	}
}

func TestFlattenJSONFileKey(t *testing.T) {
	var decoded any
	json.Unmarshal([]byte(`{"Title":"x","file:Cover":"/tmp/pic.png","Config":{"file:Doc":"/tmp/d.pdf"}}`), &decoded)

	values := url.Values{}
	var files []formFile
	flattenJSON("", decoded, values, &files)

	if values.Get("Cover") != "" || values.Get("file:Cover") != "" {
		t.Errorf("file field must not appear as a form value: %#v", values)
	}
	want := map[string]string{"Cover": "/tmp/pic.png", "Config.Doc": "/tmp/d.pdf"}
	got := map[string]string{}
	for _, f := range files {
		got[f.Field] = f.Path
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("files = %#v, want %#v", got, want)
	}
}

func TestEncodeMultipartWritesFileFromDisk(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pic.txt")
	if err := os.WriteFile(path, []byte("PNGDATA"), 0o600); err != nil {
		t.Fatal(err)
	}

	body, ct, err := encodeMultipart(url.Values{"Title": {"x"}}, []formFile{{Field: "Cover", Path: path}})
	if err != nil {
		t.Fatal(err)
	}
	if ct == "" {
		t.Fatal("empty content type")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("source file must remain on disk: %v", err)
	}
	if !containsAll(string(body), `name="Title"`, `name="Cover"`, "pic.txt", "PNGDATA") {
		t.Fatalf("multipart body missing expected parts:\n%s", body)
	}
}

func TestParseSpecsMode(t *testing.T) {
	specs, array, err := parseSpecsMode([]byte(`[{"uri":"/a"},{"uri":"/b"}]`))
	if err != nil || !array || len(specs) != 2 {
		t.Fatalf("array parse: specs=%d array=%v err=%v", len(specs), array, err)
	}
	specs, array, err = parseSpecsMode([]byte(`{"uri":"/a","method":"POST","login":"admin"}`))
	if err != nil || array || len(specs) != 1 || specs[0].Method != "POST" || specs[0].account() != "admin" {
		t.Fatalf("object parse: specs=%d array=%v login=%q err=%v", len(specs), array, specs[0].account(), err)
	}
}

func strptr(s string) *string { return &s }
func intptr(i int) *int       { return &i }

func TestValidateResponseStatus(t *testing.T) {
	if err := validateResponse(&ExpectedResponse{Status: intptr(200)}, 200, nil); err != nil {
		t.Fatalf("status match should pass: %v", err)
	}
	if err := validateResponse(&ExpectedResponse{Status: intptr(200)}, 500, nil); err == nil {
		t.Fatal("status mismatch should fail")
	}
	if err := validateResponse(nil, 500, nil); err != nil {
		t.Fatalf("nil expected should always pass: %v", err)
	}
}

func TestValidateResponseBody(t *testing.T) {
	body := []byte("hello world")
	cases := []struct {
		name string
		m    BodyMatch
		ok   bool
	}{
		{"equal ok", BodyMatch{Equal: strptr("hello world")}, true},
		{"equal fail", BodyMatch{Equal: strptr("nope")}, false},
		{"contains ok", BodyMatch{Contains: strptr("lo wo")}, true},
		{"contains fail", BodyMatch{Contains: strptr("xyz")}, false},
		{"starts ok", BodyMatch{Starts: strptr("hello")}, true},
		{"starts fail", BodyMatch{Starts: strptr("world")}, false},
		{"ends ok", BodyMatch{Ends: strptr("world")}, true},
		{"ends fail", BodyMatch{Ends: strptr("hello")}, false},
	}
	for _, c := range cases {
		err := validateResponse(&ExpectedResponse{Body: &c.m}, 200, body)
		if (err == nil) != c.ok {
			t.Errorf("%s: err=%v want ok=%v", c.name, err, c.ok)
		}
	}
}

func TestValidateResponseKeys(t *testing.T) {
	body := []byte(`{"response":{"updatePortals":[]},"flash":[]}`)
	if err := validateResponse(&ExpectedResponse{Keys: []string{"flash", "response.updatePortals"}}, 200, body); err != nil {
		t.Fatalf("existing keys should pass: %v", err)
	}
	if err := validateResponse(&ExpectedResponse{Keys: []string{"response.missing"}}, 200, body); err == nil {
		t.Fatal("missing key should fail")
	}
	if err := validateResponse(&ExpectedResponse{Keys: []string{"x"}}, 200, []byte("not json")); err == nil {
		t.Fatal("non-JSON response with key check should fail")
	}
}

func TestExpectedResponseParsesFromSpec(t *testing.T) {
	specs, _, err := parseSpecsMode([]byte(`{"uri":"/a","expectedResponse":{"status":200,"body":{"contains":"ok"},"keys":["flash"]}}`))
	if err != nil || len(specs) != 1 {
		t.Fatalf("parse: %v", err)
	}
	exp := specs[0].ExpectedResponse
	if exp == nil || exp.Status == nil || *exp.Status != 200 {
		t.Fatalf("status not parsed: %#v", exp)
	}
	if exp.Body == nil || exp.Body.Contains == nil || *exp.Body.Contains != "ok" {
		t.Fatalf("body.contains not parsed: %#v", exp)
	}
	if len(exp.Keys) != 1 || exp.Keys[0] != "flash" {
		t.Fatalf("keys not parsed: %#v", exp.Keys)
	}
}

func TestRequestSpecAccountAlias(t *testing.T) {
	if got := (requestSpec{Login: "a"}).account(); got != "a" {
		t.Fatalf("login field: got %q", got)
	}
	if got := (requestSpec{User: "b"}).account(); got != "b" {
		t.Fatalf("user alias: got %q", got)
	}
	if got := (requestSpec{Login: "a", User: "b"}).account(); got != "a" {
		t.Fatalf("login should win over user: got %q", got)
	}
}

func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		found := false
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// The opt-out travels in the request record, so a spec file can carry it per
// request — and it defaults to OFF, which is to say the stamp is required.
func TestParseSpecsSkipFormSign(t *testing.T) {
	specs, _, err := parseSpecsMode([]byte(`{"uri":"/a","skipFormSign":true}`))
	if err != nil || len(specs) != 1 {
		t.Fatalf("parse: specs=%d err=%v", len(specs), err)
	}
	if !specs[0].SkipFormSign {
		t.Error("skipFormSign do JSON não chegou ao spec")
	}

	specs, _, err = parseSpecsMode([]byte(`[{"uri":"/a"},{"uri":"/b","skipFormSign":true}]`))
	if err != nil || len(specs) != 2 {
		t.Fatalf("parse: specs=%d err=%v", len(specs), err)
	}
	if specs[0].SkipFormSign {
		t.Error("sem o campo, o padrão devia ser exigir o stamp")
	}
	if !specs[1].SkipFormSign {
		t.Error("skipFormSign do segundo request não chegou")
	}
}

// And what it does is put the flag in the request's CONTEXT, which is where the
// admin reads it (web.SkipFormSign) — a network request can never set it.
func TestServeCarriesSkipFormSignInTheContext(t *testing.T) {
	for _, skip := range []bool{false, true} {
		var got, served bool

		cfg := Config{
			Handler: func(mux *http.ServeMux) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					got, served = web.SkipFormSign(r), true
					w.WriteHeader(http.StatusOK)
				})
			},
		}

		if _, err := Serve(context.Background(), cfg, Dispatch{
			Method:       http.MethodPost,
			URI:          "/admin/things",
			ContentType:  "application/json",
			RawBody:      []byte(`{}`),
			SkipFormSign: skip,
		}); err != nil {
			t.Fatal(err)
		}

		if !served {
			t.Fatal("o handler não foi chamado")
		}
		if got != skip {
			t.Errorf("web.SkipFormSign = %v, want %v", got, skip)
		}
	}
}
