package perms

import (
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
)

// Modes of a resource. Each maps to a presets fields/actions context and to a
// permission verb.
const (
	ModeList   = "list"
	ModeDetail = "detail"
	ModeNew    = "new"
	ModeEdit   = "edit"
)

// FieldNode is a field in a mode, possibly with nested children (SPEC/user:
// nested fields also enter the permission manager).
type FieldNode struct {
	Name     string
	Children []FieldNode
}

// modeFieldBuilders returns the field builders of mb for a mode.
func modeFieldBuilders(mb *presets.ModelBuilder, mode string) presets.FieldBuilders {
	switch mode {
	case ModeList:
		return mb.ListFieldBuilders()
	case ModeDetail:
		return mb.DetailFieldBuilders()
	case ModeNew, ModeEdit:
		return mb.WriteFieldBuilders()
	}
	return nil
}

// FieldNodes enumerates the fields of a mode as nodes, recursing into nested
// fields.
func FieldNodes(mb *presets.ModelBuilder, mode string) (nodes []FieldNode) {
	return fieldNodesSlice(modeFieldBuilders(mb, mode))
}

func fieldNodesSlice(fbs presets.FieldBuilders) (nodes []FieldNode) {
	for _, f := range fbs {
		node := FieldNode{Name: f.Name()}
		if nested := f.GetNested(); nested != nil {
			node.Children = fieldNodesFromBuilder(nested.FieldsBuilder())
		}
		nodes = append(nodes, node)
	}
	return
}

func fieldNodesFromBuilder(fb *presets.FieldsBuilder) (nodes []FieldNode) {
	if fb == nil {
		return
	}
	for _, name := range fb.FieldNames() {
		n, ok := name.(string)
		if !ok || n == "" {
			continue
		}
		node := FieldNode{Name: n}
		if f := fb.GetField(n); f != nil {
			if nested := f.GetNested(); nested != nil {
				node.Children = fieldNodesFromBuilder(nested.FieldsBuilder())
			}
		}
		nodes = append(nodes, node)
	}
	return
}

// FieldResource returns the exact permission resource of a record field (or a
// nested field, given the field path), computed by the presets permissioner —
// the same computation the runtime field-permission check uses, so a policy
// stored with this resource is matched for that field.
func FieldResource(mb *presets.ModelBuilder, id model.ID, path ...string) string {
	v := mb.Permissioner().Verifier(id)
	for _, seg := range path {
		v = v.SnakeOn(presets.FieldPerm(seg))
	}
	return v.Resource()
}
