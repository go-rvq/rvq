package cli

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

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
