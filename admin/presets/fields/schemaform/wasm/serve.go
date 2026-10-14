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

// Stater is a Checker that says how the fields drawn by what is answered
// are (schemaform.Decoder: State), and the values of a field of options of
// a function (Options) — a request fetched, by the browser.
type Stater interface {
	State(data map[string]any) map[string]schemaform.FieldState
	Options(path string, data map[string]any) ([]schemaform.EnumItem, error)
}

// Changer is a Checker that does what is done as a field's value changes
// (schemaform.Decoder: Changed) — a request fetched, by the browser —: the
// fields of its record it set.
type Changer interface {
	Changed(path string, data map[string]any) (map[string]any, error)
}

// Names are the page's names of a check: the function Check(source,
// check); State(source, {data}), how the fields drawn by what is answered
// are; Options(source, {path, data}), a promise of the values of a field;
// Changed(source, {path, data}), a promise of the fields of its record a
// field's change set; Ready, true once they are there; Loaded, a function of the
// page's called then.
type Names struct {
	Check, State, Options, Changed, Ready, Loaded string
}

// NamesOf are the names of a check of prefix: <prefix>Check, <prefix>State,
// <prefix>Options, <prefix>Ready, <prefix>Loaded.
func NamesOf(prefix string) Names {
	return Names{Check: prefix + "Check", State: prefix + "State", Options: prefix + "Options",
		Changed: prefix + "Changed", Ready: prefix + "Ready", Loaded: prefix + "Loaded"}
}

// stateArgs is what State and Options are asked with.
type stateArgs struct {
	Path string         `json:"path"`
	Data map[string]any `json:"data"`
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
	if n.State != "" {
		js.Global().Set(n.State, js.FuncOf(func(_ js.Value, args []js.Value) any {
			out := func(v any) any {
				b, _ := json.Marshal(v)
				return string(b)
			}
			if len(args) != 2 {
				return out(map[string]string{"error": n.State + "(source, {data})"})
			}
			c, err := get(args[0].String())
			if err != nil {
				return out(map[string]string{"error": err.Error()})
			}
			st, ok := c.(Stater)
			if !ok {
				return out(map[string]any{})
			}
			var a stateArgs
			if err := json.Unmarshal([]byte(args[1].String()), &a); err != nil {
				return out(map[string]string{"error": err.Error()})
			}
			return out(st.State(schemaform.Unflatten(a.Data)))
		}))
	}
	if n.Options != "" {
		// a promise: the request fetched by the browser (Go's HTTP there),
		// out of the page's call
		js.Global().Set(n.Options, js.FuncOf(func(_ js.Value, args []js.Value) any {
			var source, ask string
			if len(args) == 2 {
				source, ask = args[0].String(), args[1].String()
			}
			return js.Global().Get("Promise").New(js.FuncOf(func(_ js.Value, pa []js.Value) any {
				resolve, reject := pa[0], pa[1]
				go func() {
					fail := func(err error) { reject.Invoke(js.Global().Get("Error").New(err.Error())) }
					c, err := get(source)
					if err != nil {
						fail(err)
						return
					}
					st, ok := c.(Stater)
					if !ok {
						resolve.Invoke("[]")
						return
					}
					var a stateArgs
					if err := json.Unmarshal([]byte(ask), &a); err != nil {
						fail(err)
						return
					}
					items, err := st.Options(a.Path, schemaform.Unflatten(a.Data))
					if err != nil {
						fail(err)
						return
					}
					b, _ := json.Marshal(items)
					resolve.Invoke(string(b))
				}()
				return nil
			}))
		}))
	}
	if n.Changed != "" {
		js.Global().Set(n.Changed, js.FuncOf(func(_ js.Value, args []js.Value) any {
			var source, ask string
			if len(args) == 2 {
				source, ask = args[0].String(), args[1].String()
			}
			return promise(func() (any, error) {
				c, err := get(source)
				if err != nil {
					return nil, err
				}
				ch, ok := c.(Changer)
				if !ok {
					return map[string]any{}, nil
				}
				var a stateArgs
				if err := json.Unmarshal([]byte(ask), &a); err != nil {
					return nil, err
				}
				return ch.Changed(a.Path, schemaform.Unflatten(a.Data))
			})
		}))
	}
	js.Global().Set(n.Ready, js.ValueOf(true))
	if cb := js.Global().Get(n.Loaded); cb.Type() == js.TypeFunction {
		cb.Invoke()
	}
	select {}
}

// promise is a promise of the page's of what do answers, as JSON — done out
// of the page's call (a request fetched meanwhile).
func promise(do func() (any, error)) js.Value {
	return js.Global().Get("Promise").New(js.FuncOf(func(_ js.Value, pa []js.Value) any {
		resolve, reject := pa[0], pa[1]
		go func() {
			v, err := do()
			if err != nil {
				reject.Invoke(js.Global().Get("Error").New(err.Error()))
				return
			}
			b, _ := json.Marshal(v)
			resolve.Invoke(string(b))
		}()
		return nil
	}))
}
