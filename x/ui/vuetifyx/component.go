package vuetifyx

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
)

// ComponentDefinitionCallbacks is a list of callback functions.
// Each callback function haves a one argument as object with two properties:
//   - `data`: `{window, Vue, [key: string]: any}`, is a component context data.
//   - `locals` is interface {
//     // the component definition
//     componentDefinition: ComponentDefinition
//     // the component element
//     element: Node | null
//     // the component instance
//     componentInstance:ComponentInternalInstance | null
//     // list of component nodes
//     nodes: Node[]
//     // skip calls unmounted handlers
//     skipUnmountHandlers: boolean
//     // other data set by any callbacks
//     [key: string]: any
//     }
type ComponentDefinitionCallbacks []string

// AppendFunc add function definition.
// Example: cb.AppendFunc("{locals, data}", "console.log(locals)")
func (c *ComponentDefinitionCallbacks) AppendFunc(arg, body string) {
	*c = append(*c, fmt.Sprintf("(%s) => %s", arg, body))
}

// Append append raw function definition.
// Example: cb.Append("({locals, data}) => console.log(locals)")
func (c *ComponentDefinitionCallbacks) Append(fn string) {
	*c = append(*c, fn)
}

type ComponentDefinitionOptions struct {
	Data        any                          `json:"data,omitempty"`
	Setup       ComponentDefinitionCallbacks `json:"setup,omitempty"`
	Mounted     ComponentDefinitionCallbacks `json:"mounted,omitempty"`
	Unmounted   ComponentDefinitionCallbacks `json:"unmounted,omitempty"`
	ResetScroll bool                         `json:"resetScroll,omitempty"`
}

type ComponentDefinition struct {
	Template h.HTMLComponent            `json:"template,omitempty"`
	Options  ComponentDefinitionOptions `json:"options,omitempty"`
}
