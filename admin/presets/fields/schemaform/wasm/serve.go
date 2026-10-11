//go:build js && wasm

package wasm

import (
	"encoding/json"
	"sync"
	"syscall/js"

	"github.com/go-rvq/rvq/admin/presets/fields/schemaform"
)

// Checker checks what a page asks of a form (schemaform.Decoder: its Check).
type Checker interface {
	Check(c schemaform.Check) schemaform.CheckResult
}

// Names are the page's names of a check: the function Check(source,
// check); Ready, true once it is there; Loaded, a function of the page's
// called then.
type Names struct {
	Check, Ready, Loaded string
}

// NamesOf are the names of a check of prefix: <prefix>Check, <prefix>Ready,
// <prefix>Loaded.
func NamesOf(prefix string) Names {
	return Names{Check: prefix + "Check", Ready: prefix + "Ready", Loaded: prefix + "Loaded"}
}

// Serve is the command of a check in the browser: the function of the page
// — n.Check(source, check), the form as the application gives it (load
// reads it; each read once) and what to check (schemaform.Check, as JSON) —
// answers schemaform.CheckResult, as JSON. Then n.Ready is true, and
// n.Loaded (a function of the page's, when there is one) is called. It does
// not return.
func Serve(n Names, load func(source string) (Checker, error)) {
	var (
		mu    sync.Mutex
		forms = map[string]Checker{}
	)
	get := func(source string) (Checker, error) {
		mu.Lock()
		defer mu.Unlock()
		if c, ok := forms[source]; ok {
			return c, nil
		}
		c, err := load(source)
		if err != nil {
			return nil, err
		}
		forms[source] = c
		return c, nil
	}
	js.Global().Set(n.Check, js.FuncOf(func(_ js.Value, args []js.Value) any {
		out := func(v any) any {
			b, _ := json.Marshal(v)
			return string(b)
		}
		if len(args) != 2 {
			return out(schemaform.CheckResult{Error: n.Check + "(source, check)"})
		}
		c, err := get(args[0].String())
		if err != nil {
			return out(schemaform.CheckResult{Error: err.Error()})
		}
		var ck schemaform.Check
		if err := json.Unmarshal([]byte(args[1].String()), &ck); err != nil {
			return out(schemaform.CheckResult{Error: err.Error()})
		}
		return out(c.Check(ck))
	}))
	js.Global().Set(n.Ready, js.ValueOf(true))
	if cb := js.Global().Get(n.Loaded); cb.Type() == js.TypeFunction {
		cb.Invoke()
	}
	select {}
}
