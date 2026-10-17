package gadcode

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// gad, gadx and gadt formatted as `gad fmt` does — a gadt of the delimiters
// it is told —; what does not parse, and another language, as they are.
func TestFormat(t *testing.T) {
	for _, c := range []struct {
		code Code
		want string
		ok   bool
	}{
		{Code{Source: "export class default {a str}", Language: Gad}, "export class default {\n\ta str\n}", true},
		{Code{Source: "div.a\n  span {= x }\n", Language: Gadx}, "div.a\n\tspan {= x }\n", true},
		{Code{Source: "Hi, {name}!", Language: Gadt, TemplateStart: "{", TemplateEnd: "}"}, "Hi, { name }!", true},
		{Code{Source: "export class default {", Language: Gad}, "export class default {", false},
		{Code{Source: "a  b", Language: "js"}, "a  b", false},
	} {
		got, ok := Format(c.code)
		if got != c.want || ok != c.ok {
			t.Errorf("%s %q: got %q %v, want %q %v", c.code.Language, c.code.Source, got, ok, c.want, c.ok)
		}
	}
}

// The structure: the declarations, a class's fields under it, with where
// each is; of a file being written, what parsed.
func TestOutline(t *testing.T) {
	syms := Outline(Code{Language: Gad, Source: "enum Kind { A, B }\nclass Link { href str }\nexport class default {\n\tlinks []Link\n\tkind Kind\n}"})
	var names []string
	for _, s := range syms {
		names = append(names, s.Kind+":"+s.Name)
	}
	if got := strings.Join(names, " "); !strings.Contains(got, "enum:Kind") || !strings.Contains(got, "class:Link") || !strings.Contains(got, "class:default") {
		t.Fatalf("outline: %s", got)
	}
	last := syms[len(syms)-1]
	if len(last.Children) != 2 || last.Children[0].Name != "links" || last.Children[0].Line != 4 {
		t.Errorf("the fields of default: %+v", last.Children)
	}
	if got := Outline(Code{Language: Gad, Source: "class A { a str }\nclass B {"}); len(got) == 0 || got[0].Name != "A" {
		t.Errorf("a file being written: %+v", got)
	}
}

// The operations, posted as JSON: format, outline, doc; another is none.
func TestServe(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /admin"+Path+"{op}", serve)
	post := func(op string, c Code) (*httptest.ResponseRecorder, map[string]any) {
		b, _ := json.Marshal(c)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, httptest.NewRequest("POST", "/admin"+Path+op, strings.NewReader(string(b))))
		var m map[string]any
		_ = json.Unmarshal(w.Body.Bytes(), &m)
		return w, m
	}
	code := Code{Language: Gad, Source: "// The class.\nexport class default {a str}"}
	if _, m := post("format", code); m["ok"] != true || !strings.Contains(m["source"].(string), "\n\ta str\n") {
		t.Errorf("format: %v", m)
	}
	w, _ := post("outline", code)
	if !strings.Contains(w.Body.String(), `"name":"default"`) {
		t.Errorf("outline: %s", w.Body.String())
	}
	if _, m := post("doc", code); m["html"] == "" || m["error"] != nil {
		t.Errorf("doc: %v", m)
	}
	if w, _ := post("nothing", code); w.Code != http.StatusNotFound {
		t.Errorf("another op: %d", w.Code)
	}
}
