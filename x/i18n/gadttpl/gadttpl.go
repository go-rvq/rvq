// Package gadttpl is a message written as a Gad template (gadt) with the
// delimiters of the SEO fields: text with { … } code and {= expr } writing the
// value of expr (the SEO fields write { expr }; a message has statements —
// { if x begin }…{ else }…{ end }). Rendered where it is used, edited as gadt
// (the schemaform's CodeMirror) in the admin, where its source is the
// message's value. A message of HTML with markup of its own is a Gadx one
// (gadxtpl); a gadt is text, HTML when the message is.
package gadttpl

import (
	"bytes"
	"fmt"
	"html"
	"sort"
	"strings"
	"sync"

	"github.com/gad-lang/gad"
	"github.com/gad-lang/gad/parser"
	h "github.com/go-rvq/htmlgo"
)

// Delimiter is that of the code of a template, as in the SEO fields: { … }.
var Delimiter = parser.MixedDelimiter{Start: []rune("{"), End: []rune("}")}

// Template is the gadt source of a message. Empty, it renders nothing.
type Template string

type compiled struct {
	bc       *gad.Bytecode
	builtins *gad.Builtins
}

// cache holds a template's compilation per source and names of its data.
var cache sync.Map

// Compile is the template compiled for data of the given names.
func (t Template) Compile(names ...string) (*gad.Bytecode, *gad.Builtins, error) {
	names = append([]string(nil), names...)
	sort.Strings(names)
	key := string(t) + "\x00" + strings.Join(names, ",")
	if c, ok := cache.Load(key); ok {
		c := c.(*compiled)
		return c.bc, c.builtins, nil
	}
	builtins := gad.NewBuiltins()
	st := gad.NewSymbolTable(builtins.NameSet)
	if _, err := st.DefineGlobals(names); err != nil {
		return nil, nil, err
	}
	opts := gad.CompileOptions{}
	opts.ModuleFile = "message.gadt"
	// The main source is not told by its name, as a module is: a gadt is
	// parsed in the mixed mode, { … } its code and {= expr } a value
	// (Delimiter).
	opts.ParserOptions.Mode |= parser.ParseMixed
	opts.ScannerOptions.Mode |= parser.ScanMixed | parser.ScanConfigDisabled
	opts.ScannerOptions.MixedDelimiter = Delimiter
	res, err := gad.Compile(st, []byte(t), opts)
	if err != nil {
		return nil, nil, err
	}
	bc := res.Bytecode
	cache.Store(key, &compiled{bc: bc, builtins: builtins})
	return bc, builtins, nil
}

// Render writes the template's text, data its globals.
func (t Template) Render(data gad.Dict) (string, error) {
	if strings.TrimSpace(string(t)) == "" {
		return "", nil
	}
	names := make([]string, 0, len(data))
	for name := range data {
		names = append(names, name)
	}
	bc, builtins, err := t.Compile(names...)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	vm := gad.NewVM(builtins.Build(), bc)
	if _, err = vm.RunOpts(&gad.RunOpts{StdOut: &out, Globals: data}); err != nil {
		return "", err
	}
	return out.String(), nil
}

// HTML is the template rendered, a message of HTML — an error in it shown
// where the HTML would be, escaped, rather than breaking the page.
func (t Template) HTML(data gad.Dict) h.RawHTML {
	out, err := t.Render(data)
	if err != nil {
		return h.RawHTML(fmt.Sprintf(`<pre class="text-error">%s</pre>`, html.EscapeString(err.Error())))
	}
	return h.RawHTML(out)
}

// JoinAnd is join_and(sep, lastSep, elem…), for the data of a template: the
// elements not empty, sep between them and lastSep before the last
// ("a, b and c").
var JoinAnd = gad.NewFunction("join_and", func(c gad.Call) (gad.Object, error) {
	if c.Args.Length() < 2 {
		return nil, fmt.Errorf("join_and: want sep, lastSep, elem…")
	}
	sep, last := c.Args.Get(0).ToString(), c.Args.Get(1).ToString()
	var elems []string
	for i := 2; i < c.Args.Length(); i++ {
		if s := c.Args.Get(i).ToString(); s != "" {
			elems = append(elems, s)
		}
	}
	if len(elems) < 2 {
		return gad.Str(strings.Join(elems, "")), nil
	}
	return gad.Str(strings.Join(elems[:len(elems)-1], sep) + last + elems[len(elems)-1]), nil
})
