package seo

import (
	"encoding/json"
)

// SettingVarType is a registered custom-variable type: an id, a label, and the
// Vue component that edits a variable of this type. The application registers
// types (RegisterVarType); a variable whose Type is unregistered (or empty) is
// plain text — the editor then only toggles between a single- and multi-line
// input.
//
// The editor component is a custom element (tag name) that receives, per row,
// `v-model:data` (its own state — SettingVar.Data) and `v-model:value` (the
// final SettingVar.Value). A component that computes its value server-side posts
// to Action and lets the handler set the value.
type SettingVarType struct {
	// Name is the type id stored on SettingVar.Type.
	Name string
	// Label is what the user sees when picking the type.
	Label string
	// Description is a short subtitle for the type.
	Description string
	// Component is the Vue custom-element tag that edits a variable of this type
	// (e.g. "vx-seo-var-zipcodes"). Empty means the plain text editor.
	Component string
	// Action is the SEO-model event name the component posts to (optional); its
	// callback computes the final value. Empty when the component needs no server
	// round-trip.
	Action string
}

// RegisterVarType registers a custom-variable type. Registering the same Name
// twice replaces the earlier one.
func (b *Builder) RegisterVarType(t SettingVarType) *Builder {
	if b.varTypes == nil {
		b.varTypes = map[string]SettingVarType{}
	}
	b.varTypes[t.Name] = t
	return b
}

// VarType returns a registered type by name, or false.
func (b *Builder) VarType(name string) (SettingVarType, bool) {
	t, ok := b.varTypes[name]
	return t, ok
}

// VarTypes returns the registered types, in no particular order.
func (b *Builder) VarTypes() []SettingVarType {
	out := make([]SettingVarType, 0, len(b.varTypes))
	for _, t := range b.varTypes {
		out = append(out, t)
	}
	return out
}

// varTypesJSON is the registered types as a client-side map
// {typeName: {label, description, component, action}}, so the editor component
// can pick each row's editor by its Type. Types with no Component are omitted
// (they fall back to the plain text editor on the client).
func (b *Builder) varTypesJSON() string {
	m := map[string]map[string]string{}
	for name, t := range b.varTypes {
		if t.Component == "" {
			continue
		}
		m[name] = map[string]string{
			"label":       t.Label,
			"description": t.Description,
			"component":   t.Component,
			"action":      t.Action,
		}
	}
	b1, _ := json.Marshal(m)
	return string(b1)
}

// MergeVars overlays higher-priority variables onto lower-priority ones by Name:
// a var in high replaces a var of the same Name in low; the rest of low is kept.
// This is how a descendant SEO's vars override an ancestor's (and add new ones),
// while inheriting the ancestor's other vars. Order follows low, then any names
// only in high, so inherited vars keep a stable position.
func MergeVars(low, high []SettingVar) []SettingVar {
	idx := map[string]int{}
	out := make([]SettingVar, 0, len(low)+len(high))
	for _, v := range low {
		idx[v.Name] = len(out)
		out = append(out, v)
	}
	for _, v := range high {
		if i, ok := idx[v.Name]; ok {
			out[i] = v
		} else {
			idx[v.Name] = len(out)
			out = append(out, v)
		}
	}
	return out
}

// VarsValues maps each variable's Name to its resolved Value — what a template
// reads as {Name}. A later var of the same Name (already merged) wins.
func VarsValues(vars []SettingVar) map[string]string {
	m := make(map[string]string, len(vars))
	for _, v := range vars {
		if v.Name != "" {
			m[v.Name] = v.Value
		}
	}
	return m
}
