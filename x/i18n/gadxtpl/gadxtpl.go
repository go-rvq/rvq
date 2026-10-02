// Package gadxtpl is a message written in Gadx: a text of the messages of a
// module that is HTML with logic — rendered where it is used, edited as Gadx
// (the schemaform's CodeMirror) in the admin, where its source is the
// message's value.
package gadxtpl

import (
	"bytes"
	"fmt"
	"html"
	"sort"
	"strings"
	"sync"

	"github.com/gad-lang/gad"
	"github.com/gad-lang/gad/gadx"
	h "github.com/go-rvq/htmlgo"
)

// Template is the Gadx source of a message. Empty, it renders nothing.
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
	builtins := gadx.AppendBuiltins(gad.NewBuiltins())
	st := gad.NewSymbolTable(builtins.NameSet)
	if _, err := st.DefineGlobals(names); err != nil {
		return nil, nil, err
	}
	opts := gad.CompileOptions{}
	opts.ModuleFile = "message.gadx"
	res, err := gad.Compile(st, []byte(t), opts)
	if err != nil {
		return nil, nil, err
	}
	bc := res.Bytecode
	cache.Store(key, &compiled{bc: bc, builtins: builtins})
	return bc, builtins, nil
}

// Render writes the template's HTML, data its globals.
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
	ret, err := vm.RunOpts(&gad.RunOpts{StdOut: &out, Globals: data})
	if err != nil {
		return "", err
	}
	if el, ok := ret.(gadx.Element); ok {
		if _, err := el.WriteTo(vm, &out); err != nil {
			return "", err
		}
	}
	return out.String(), nil
}

// Component is the template rendered, as a component — an error in it shown
// where the HTML would be, escaped, rather than breaking the page.
func (t Template) Component(data gad.Dict) h.HTMLComponent {
	out, err := t.Render(data)
	if err != nil {
		return h.RawHTML(fmt.Sprintf(`<pre class="text-error">%s</pre>`, html.EscapeString(err.Error())))
	}
	return h.RawHTML(out)
}
