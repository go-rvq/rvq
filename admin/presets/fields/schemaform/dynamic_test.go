package schemaform

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

// citiesServer answers the cities of a state (?state=…) as {data: […]}; of
// the state "XX", a failure. It counts what it is asked.
func citiesServer(t *testing.T) (*httptest.Server, *int32) {
	t.Helper()
	var asked int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&asked, 1)
		switch r.URL.Query().Get("state") {
		case "Massachusetts":
			json.NewEncoder(w).Encode(map[string]any{"data": []string{"Boston", "Salem"}})
		case "Maine":
			json.NewEncoder(w).Encode(map[string]any{"data": []string{"Portland"}})
		default:
			http.Error(w, "no", http.StatusInternalServerError)
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &asked
}

func dynamicSchema(t *testing.T, url, onError string) *Schema {
	t.Helper()
	requestCache.Range(func(k, _ any) bool { requestCache.Delete(k); return true })
	src := `const states = (;MA="Massachusetts", ME="Maine", XX="Nowhere")
export class Form {
	kind str
	lead class {
		use WithWarnings
		[options=states]
		state? str
		[disabled=(this) => !this.state, optionsError="` + onError + `",
			options=(this, parents) => this.state ? Request("` + url + `"; query={state: states.get(this.state)}).then((r) => r.data) : nil]
		city str
		[if=(this, parents) => parents[0].fieldName == "lead" && parents[0].self.kind == "house"]
		rooms int
		[show=(this) => this.state == "ME"]
		harbor str
		[readonly=true]
		code str
		[class=(this) => ({wide: this.state == "MA"})]
		note? str
	}
}`
	s, err := Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func issuesText(issues []Issue) string {
	var out []string
	for _, e := range issues {
		out = append(out, e.Path+":"+e.Code)
	}
	return strings.Join(out, " ")
}

// A field of a function of what is answered: disabled, not there or hidden
// it is not asked — not required, its value dropped —; read-only, not
// required. Its options, of a request: the value one of them.
func TestDynamicFieldsDecode(t *testing.T) {
	srv, _ := citiesServer(t)
	s := dynamicSchema(t, srv.URL, OptionsErrorRefuse)

	// no state: the city disabled; a flat, the rooms not there; the harbor
	// hidden; the code read-only
	rec, issues := (&Decoder{Schema: s}).Decode(map[string]any{"kind": "flat",
		"lead": map[string]any{"city": "Boston"}})
	if len(issues) != 0 {
		t.Fatalf("no state: %s", issuesText(issues))
	}
	lead := rec.Map()["lead"].(Record).Map()
	if lead["city"] != nil || lead["rooms"] != nil {
		t.Errorf("a field not asked kept its value: %v", lead)
	}

	// a house, of Massachusetts: the rooms asked, the city one of its
	data := map[string]any{"kind": "house", "lead": map[string]any{"state": "MA", "city": "Boston", "rooms": "3"}}
	rec, issues = (&Decoder{Schema: s}).Decode(data)
	if len(issues) != 0 {
		t.Fatalf("Boston, MA: %s", issuesText(issues))
	}
	if got := rec.Map()["lead"].(Record).Map()["city"]; got != "Boston" {
		t.Errorf("the city: %v", got)
	}
	data["lead"].(map[string]any)["city"] = "Portland"
	if _, issues = (&Decoder{Schema: s}).Decode(data); issuesText(issues) != "lead.city:choice" {
		t.Errorf("Portland, MA: %s", issuesText(issues))
	}
	delete(data["lead"].(map[string]any), "rooms")
	data["lead"].(map[string]any)["city"] = "Boston"
	if _, issues = (&Decoder{Schema: s}).Decode(data); issuesText(issues) != "lead.rooms:required" {
		t.Errorf("a house, no rooms: %s", issuesText(issues))
	}
	// Maine: the harbor shown, asked
	if _, issues = (&Decoder{Schema: s}).Decode(map[string]any{"kind": "flat",
		"lead": map[string]any{"state": "ME", "city": "Portland"}}); issuesText(issues) != "lead.harbor:required" {
		t.Errorf("Maine, no harbor: %s", issuesText(issues))
	}
}

// Options that could not be had: refused (the default), or — warn — the
// value kept, the warning in the record's $warnings (WithWarnings).
func TestDynamicOptionsError(t *testing.T) {
	srv, _ := citiesServer(t)
	data := map[string]any{"kind": "flat", "lead": map[string]any{"state": "XX", "city": "Gotham"}}

	if _, issues := (&Decoder{Schema: dynamicSchema(t, srv.URL, OptionsErrorRefuse)}).Decode(data); issuesText(issues) != "lead.city:items" {
		t.Errorf("refused: %s", issuesText(issues))
	}

	d := &Decoder{Schema: dynamicSchema(t, srv.URL, OptionsErrorWarn)}
	rec, issues := d.Decode(data)
	if len(issues) != 0 {
		t.Fatalf("warned: %s", issuesText(issues))
	}
	lead := rec.Map()["lead"].(Record).Map()
	ws, _ := lead[WarningsField].([]any)
	if lead["city"] != "Gotham" || len(ws) != 1 || len(d.Warnings()) != 1 {
		t.Fatalf("the value kept, a warning: %v; %v", lead, d.Warnings())
	}
	if w := ws[0].(map[string]any); w["fieldName"] != "city" || !strings.Contains(w["message"].(string), "500") {
		t.Errorf("the warning: %v", w)
	}
}

// A record of a field that keeps a value with a warning uses WithWarnings.
func TestDynamicWarnNeedsWithWarnings(t *testing.T) {
	_, err := Parse(`export class Form { lead class { [optionsError="warn", options=(this) => nil] city str } }`)
	if err == nil || !strings.Contains(err.Error(), WithWarningsClass) {
		t.Errorf("no WithWarnings: %v", err)
	}
	if _, err := Parse(`export class Form { lead class { [optionsError="maybe"] city str } }`); err == nil {
		t.Error("an optionsError of no meaning taken")
	}
}

// What the page asks: how each field is (State) — a request's key, not
// fetched —, the values of a field (Options), fetched once (cached).
func TestDynamicStateAndOptions(t *testing.T) {
	srv, asked := citiesServer(t)
	s := dynamicSchema(t, srv.URL, OptionsErrorRefuse)

	st := (&Decoder{Schema: s}).State(map[string]any{"kind": "house"})
	city := st["lead.city"]
	if !city.Disabled || !city.NoOptions || !city.Dynamic || !st["lead.rooms"].If || st["lead.harbor"].Show {
		t.Errorf("no state: %+v", st)
	}
	if v, ok := st["lead.note"].Class.(Meta); !ok || v["wide"] != false {
		t.Errorf("the class: %#v", st["lead.note"].Class)
	}
	data := map[string]any{"kind": "flat", "lead": map[string]any{"state": "MA"}}
	st = (&Decoder{Schema: s}).State(data)
	if city = st["lead.city"]; city.Disabled || !strings.Contains(city.Request, "state=Massachusetts") || st["lead.rooms"].If {
		t.Errorf("Massachusetts: %+v", st)
	}
	if atomic.LoadInt32(asked) != 0 {
		t.Error("State fetched")
	}
	for i := 0; i < 2; i++ {
		items, err := (&Decoder{Schema: s}).Options("lead.city", data)
		if err != nil || len(items) != 2 || items[0].Name != "Boston" {
			t.Fatalf("Options: %+v %v", items, err)
		}
	}
	if n := atomic.LoadInt32(asked); n != 1 {
		t.Errorf("asked %d times, not once", n)
	}
}

// The page's check: by the options it has; of a field it has none of (a
// request it has not answered), any value — nothing fetched.
func TestDynamicCheckByPageOptions(t *testing.T) {
	srv, asked := citiesServer(t)
	s := dynamicSchema(t, srv.URL, OptionsErrorRefuse)
	data := map[string]any{"kind": "flat", "lead.state": "MA", "lead.city": "Salem"}
	res := (&Decoder{Schema: s}).Check(Check{Data: data, Options: map[string][]EnumItem{"lead.city": {{Name: "Boston"}}}})
	if len(res.Errors) != 1 || res.Errors[0].Path != "lead.city" {
		t.Errorf("by the page's options: %+v", res)
	}
	if res = (&Decoder{Schema: s}).Check(Check{Data: data}); len(res.Errors) != 0 {
		t.Errorf("no options of the page: %+v", res)
	}
	if atomic.LoadInt32(asked) != 0 {
		t.Error("the page's check fetched")
	}
}

// Request: only http(s); its query; what a then answers.
func TestRequest(t *testing.T) {
	srv, _ := citiesServer(t)
	requestCache.Range(func(k, _ any) bool { requestCache.Delete(k); return true })
	s, err := Parse(`export class Form {
	[options=(this) => Request("` + srv.URL + `"; query={state: "Maine"}).then((r) => r.data).then((l) => [[l[0], "City of " + l[0]]])]
	city str
}`)
	if err != nil {
		t.Fatal(err)
	}
	items, err := (&Decoder{Schema: s}).Options("city", map[string]any{})
	if err != nil || len(items) != 1 || items[0].Name != "Portland" || items[0].Label != "City of Portland" {
		t.Errorf("then: %+v %v", items, err)
	}
	bad, _ := Parse(`export class Form {
	[options=(this) => Request("file:///etc/passwd")]
	city str
}`)
	if _, err := (&Decoder{Schema: bad}).Options("city", map[string]any{}); err == nil || !strings.Contains(err.Error(), "http") {
		t.Errorf("a file address: %v", err)
	}
}

// A handler's this is the instance of the class of the field's record — its
// values set: this.state —; parents, one item a record above: the field of
// the parent it is (fieldName) and the parent's instance (self).
func TestDynamicThisAndParents(t *testing.T) {
	s, err := Parse(`export class Form {
	kind str
	lead class {
		state str
		[options=(this, parents) => [
			typeName(this), this.state,
			str(len(parents)), parents[0].fieldName, typeName(parents[0].self), parents[0].self.kind,
		]]
		city? str
	}
}`)
	if err != nil {
		t.Fatal(err)
	}
	items, err := (&Decoder{Schema: s}).Options("lead.city", map[string]any{"kind": "house", "lead": map[string]any{"state": "Maine"}})
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, it := range items {
		got = append(got, it.Name)
	}
	if strings.Join(got, "|") != "Form.lead|Maine|1|lead|Form|house" {
		t.Errorf("this and parents: %v", got)
	}
	// nothing answered yet: still the instances
	s2, err := Parse(`export class Form {
	lead class {
		state str
		[options=(this, parents) => [typeName(this), typeName(parents[0].self)]]
		city? str
	}
}`)
	if err != nil {
		t.Fatal(err)
	}
	items, err = (&Decoder{Schema: s2}).Options("lead.city", map[string]any{})
	if err != nil || len(items) != 2 || items[0].Name != "Form.lead" || items[1].Name != "Form" {
		t.Errorf("nothing answered: %+v %v", items, err)
	}
}

// A field's [changed=…]: the fields of its record it set (this's) — in a
// request's then —; nothing set, none; a changed that is not a function,
// refused.
func TestDynamicChanged(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/us/02108" {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"places": []any{map[string]any{"place name": "Boston", "state": "Massachusetts"}}})
	}))
	defer srv.Close()
	requestCache.Range(func(k, _ any) bool { requestCache.Delete(k); return true })
	s, err := Parse(`export class Form {
	lead class {
		[changed=(this, parents) => len(this.zip) >= 5 ? Request("` + srv.URL + `/us/" + this.zip[:5]).then(func(r) {
			this.state = r.places[0].state
			this.city = r.places[0]["place name"]
		}) : nil]
		zip str
		state? str
		city? str
		unit? str
	}
}`)
	if err != nil {
		t.Fatal(err)
	}
	got, err := (&Decoder{Schema: s}).Changed("lead.zip", map[string]any{"lead": map[string]any{"zip": "02108-1234", "unit": "4B"}})
	if err != nil || len(got) != 2 || got["state"] != "Massachusetts" || got["city"] != "Boston" {
		t.Errorf("changed: %v %v", got, err)
	}
	if got, err = (&Decoder{Schema: s}).Changed("lead.zip", map[string]any{"lead": map[string]any{"zip": "021"}}); err != nil || got != nil {
		t.Errorf("a zip not whole: %v %v", got, err)
	}
	if _, err = (&Decoder{Schema: s}).Changed("lead.zip", map[string]any{"lead": map[string]any{"zip": "99999"}}); err == nil {
		t.Error("a zip not found: no error")
	}
	if _, err := Parse(`export class Form { [changed="x"] zip str }`); err == nil {
		t.Error("a changed not a function taken")
	}
}
