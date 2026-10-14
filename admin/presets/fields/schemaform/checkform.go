package schemaform

import "strings"

// Check is what a page asks of a form before it posts — in the browser
// (schemaform/wasm), the same code the server checks with —: its values by
// their path ("a.name", as it posts them; Unflatten) and the files it holds
// — by the name of their part, what the browser says of each (the
// application reads them: Decoder.Files) —; Step, the step being passed (of
// a form of steps), "" the form whole.
type Check struct {
	Step  string              `json:"step"`
	Data  map[string]any      `json:"data"`
	Files map[string]FileInfo `json:"files"`
	// Options are the values the page has of the fields of options of a
	// function, by path: what their values are checked by (a request not
	// fetched again; a field the page has none of, any value).
	Options map[string][]EnumItem `json:"options"`
}

// FileInfo is a file as the browser says it: its name, its size, its type.
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

// CheckResult is what is wrong, field by field, and, of a step passed, the
// next (Next; "" after the last, Last); Error, the form's own error (it
// could not be read).
type CheckResult struct {
	Errors []Issue `json:"errors,omitempty"`
	Next   string  `json:"next,omitempty"`
	Last   bool    `json:"last,omitempty"`
	Error  string  `json:"error,omitempty"`
}

// Check checks c as a post is: of a step, its fields and its next; else the
// form whole (of a form of steps, the steps its way goes through).
func (d *Decoder) Check(c Check) CheckResult {
	data := Unflatten(c.Data)
	d.PageOptions, d.NoFetch = c.Options, true
	if d.PageOptions == nil {
		d.PageOptions = map[string][]EnumItem{}
	}
	fl := d.Schema.Flow()
	var res CheckResult
	switch {
	case c.Step != "" && fl != nil:
		if _, errs := d.DecodeStep(fl, c.Step, data); len(errs) > 0 {
			res.Errors = errs
			return res
		}
		next, err := fl.Next(c.Step, d.Typed(data).Map())
		if err != nil {
			res.Errors = []Issue{issueOf(err)}
			return res
		}
		res.Next, res.Last = next, next == ""
	case fl != nil:
		_, res.Errors = d.DecodeFlow(fl, data)
	default:
		_, res.Errors = d.Decode(data)
	}
	return res
}

// Unflatten is data — a field of a field by its path ("address.street") —
// as the form's records hold it.
func Unflatten(data map[string]any) map[string]any {
	out := map[string]any{}
	for k, v := range data {
		setNested(out, strings.Split(k, "."), v)
	}
	return out
}
