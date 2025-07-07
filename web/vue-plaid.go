package web

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-rvq/rvq/web/js"
	h "github.com/theplant/htmlgo"
)

type JsCall struct {
	Method string
	Args   []interface{}
	Raw    string
}

func (j JsCall) String() string {
	if j.Raw != "" {
		return j.Raw
	}
	return "#" + j.Method + "(...)"
}

func (j *JsCall) UnmarshalJSON(bytes []byte) error {
	var data = make([]any, 0)
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}
	if len(data) == 3 {
		j.Method = data[0].(string)
		if data[1] != nil {
			j.Args = data[1].([]any)
		}
		j.Raw = data[2].(string)
	}
	return nil
}

func (j JsCall) MarshalJSON() ([]byte, error) {
	return json.Marshal([]any{j.Method, j.Raw, j.Args})
}

type VueEventTagBuilder struct {
	beforeScript string
	Calls        []JsCall
	afterScript  string
	thenScript   string
}

func Plaid() (r *VueEventTagBuilder) {
	r = &VueEventTagBuilder{
		Calls: []JsCall{
			{
				Method: "plaid",
			},
		},
	}
	r.Vars(Var("vars")).
		Locals(Var("locals")).
		Form(Var("form")).
		Closer(Var("closer")).
		Scope(js.Object{
			"presetsListing": js.Raw("presetsListing"),
			// "presetsDetailing": js.Raw("presetsDetailing"),
			// "presetsCreating":  js.Raw("presetsCreating"),
			// "presetsEditing":   js.Raw("presetsEditing"),
		})
	return
}

func GET() (r *VueEventTagBuilder) {
	return Plaid().Method("GET")
}
func DELETE() (r *VueEventTagBuilder) {
	return Plaid().Method("POST").Query(RequestMethodParam, "DELETE")
}

func POST() (r *VueEventTagBuilder) {
	return Plaid().Method("POST")
}

func (b VueEventTagBuilder) Clone() *VueEventTagBuilder {
	return &b
}

// URL is request page url without push state
func (b *VueEventTagBuilder) URL(url interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "url",
		Args:   []interface{}{url},
	})
	return b
}

func (b *VueEventTagBuilder) ParseURL(url string) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "parseUrl",
		Args:   []interface{}{url},
	})
	return b
}

func (b *VueEventTagBuilder) Parent(index int, v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "parent",
		Args:   []interface{}{index, v},
	})
	return b
}

func (b *VueEventTagBuilder) EventFunc(id interface{}) (r *VueEventTagBuilder) {
	c := JsCall{
		Method: "eventFunc",
		Args:   []interface{}{id},
	}
	b.Calls = append(b.Calls, c)
	return b
}

func (b *VueEventTagBuilder) NoCache() (r *VueEventTagBuilder) {
	c := JsCall{
		Method: "noCache",
	}
	b.Calls = append(b.Calls, c)
	return b
}

func (b *VueEventTagBuilder) Method(v interface{}) (r *VueEventTagBuilder) {
	c := JsCall{
		Method: "method",
		Args:   []interface{}{v},
	}
	b.Calls = append(b.Calls, c)
	return b
}

func (b *VueEventTagBuilder) Reload() (r *VueEventTagBuilder) {
	b.Raw("reload()")
	return b
}

func (b *VueEventTagBuilder) Vars(v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "vars",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) Locals(v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "locals",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) Scope(v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "scope",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) MergeQuery(v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "mergeQuery",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) Query(key interface{}, vs interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "query",
		Args:   []interface{}{key, vs},
	})
	return b
}

func (b *VueEventTagBuilder) ValidQuery(key, vs interface{}) (r *VueEventTagBuilder) {
	if vs != nil {
		if vs := fmt.Sprint(vs); vs != "" {
			b.Query(key, vs)
		}
	}
	return b
}

func (b *VueEventTagBuilder) QueryIf(key interface{}, vs interface{}, add bool) (r *VueEventTagBuilder) {
	if !add {
		return b
	}
	b.Calls = append(b.Calls, JsCall{
		Method: "query",
		Args:   []interface{}{key, vs},
	})
	return b
}

// ClearMergeQuery param v use interface{} because you can not only pass []string,
// but also pass in javascript variables by using web.Var("$event")
func (b *VueEventTagBuilder) ClearMergeQuery(v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "clearMergeQuery",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) StringQuery(v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "stringQuery",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) PushState(v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "pushState",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) Location(v *LocationBuilder) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "location",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) Queries(v url.Values) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "queries",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) PushStateURL(v string) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "pushStateURL",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) Form(v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "form",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) Closer(v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "closer",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) FieldValue(name interface{}, v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "fieldValue",
		Args:   []interface{}{name, v},
	})
	return b
}

func (b *VueEventTagBuilder) PopState(v interface{}) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "popstate",
		Args:   []interface{}{v},
	})
	return b
}

func (b *VueEventTagBuilder) Run(script string) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "run",
		Args:   []interface{}{script},
	})
	return b
}

func (b *VueEventTagBuilder) Raw(script string) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Raw: script,
	})
	return b
}

func (b *VueEventTagBuilder) PreFetch(v string) (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "preFetch",
		Args:   []interface{}{Var(v)},
	})
	return b
}

func (b *VueEventTagBuilder) Go() (r string) {
	b.Raw("go()")
	return b.String()
}

func (b *VueEventTagBuilder) JSON() (r *VueEventTagBuilder) {
	b.Calls = append(b.Calls, JsCall{
		Method: "json",
	})
	return b
}

func (b *VueEventTagBuilder) RunPushState() (r string) {
	b.Raw("runPushState()")
	return b.String()
}

func (b *VueEventTagBuilder) BeforeScript(script string) (r *VueEventTagBuilder) {
	b.beforeScript = script
	return b
}

func (b *VueEventTagBuilder) AfterScript(script string) (r *VueEventTagBuilder) {
	b.afterScript = script
	return b
}

func (b *VueEventTagBuilder) ThenScript(script string) (r *VueEventTagBuilder) {
	b.thenScript = script
	return b
}

func (b *VueEventTagBuilder) String() string {
	var cs []string
	for _, c := range b.Calls {
		if len(c.Raw) > 0 {
			cs = append(cs, c.Raw)
			continue
		}

		if len(c.Args) == 0 {
			cs = append(cs, fmt.Sprintf("%s()", c.Method))
			continue
		}

		if len(c.Args) == 1 {
			cs = append(cs, fmt.Sprintf("%s(%s)", c.Method, toJsValue(c.Args[0])))
			continue
		}

		var args []string
		for _, arg := range c.Args {
			args = append(args, toJsValue(arg))
		}
		cs = append(cs, fmt.Sprintf("%s(%s)", c.Method, strings.Join(args, ", ")))
	}

	if len(b.thenScript) > 0 {
		cs = append(cs, fmt.Sprintf("then(function(r){ %s })", b.thenScript))
	}

	var sems []string
	if len(b.beforeScript) > 0 {
		sems = append(sems, b.beforeScript)
	}
	sems = append(sems, strings.Join(cs, "."))
	if len(b.afterScript) > 0 {
		sems = append(sems, b.afterScript)
	}
	return strings.Join(sems, "; ")
}

func (b *VueEventTagBuilder) Encode() string {
	var buf bytes.Buffer
	w, _ := flate.NewWriter(&buf, flate.BestCompression)
	json.NewEncoder(w).Encode(b.Calls)
	w.Close()
	return buf.String()
}

func (b *VueEventTagBuilder) Decode(s string) {
	r := flate.NewReader(bytes.NewBufferString(s))
	json.NewDecoder(r).Decode(&b.Calls)
}

type VueEventTagBuilderSlice []*VueEventTagBuilder

func (v VueEventTagBuilderSlice) Encode() string {
	var elems = make([][]JsCall, len(v))
	for i, builder := range v {
		elems[i] = builder.Calls
	}

	var buf bytes.Buffer
	w, _ := flate.NewWriter(&buf, flate.BestCompression)
	json.NewEncoder(w).Encode(elems)
	w.Close()
	return buf.String()
}

func (v *VueEventTagBuilderSlice) Decode(s string) {
	r := flate.NewReader(bytes.NewBufferString(s))
	var elems = make([][]JsCall, 0)
	json.NewDecoder(r).Decode(elems)
	*v = make([]*VueEventTagBuilder, len(elems))
	for i, elem := range elems {
		(*v)[i] = &VueEventTagBuilder{Calls: elem}
	}
}

func toJsValue(v interface{}) string {
	switch t := v.(type) {
	case js.Raw:
		return string(t)
	case js.Object:
		return t.String()
	case js.RawSlice:
		return t.String()
	case Var:
		return string(t)
	default:
		return h.JSONString(v)
	}
}

func (b *VueEventTagBuilder) MarshalJSON() ([]byte, error) {
	panic(fmt.Sprintf("call .Go() at the end, value: %s", b.String()))
}
