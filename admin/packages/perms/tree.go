package perms

import (
	"strings"

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

// ActionNode is a detail action of a resource, grantable per record.
type ActionNode struct {
	Name string // action name (label)
	Verb string // perm verb (PermName) checked at runtime
}

// ActionNodes enumerates the record-level (detail) actions of a resource. Each
// is grantable by allowing its verb on the record resource.
func ActionNodes(mb *presets.ModelBuilder) (out []ActionNode) {
	if !mb.HasDetailing() {
		return
	}
	for _, a := range mb.Detailing().GetActions() {
		out = append(out, ActionNode{Name: a.Name(), Verb: a.PermName()})
	}
	return
}

// PageNode is a listing/detail page grantable per record (detail) or per
// listing (list). Its Path is used as the permission segment ("verb").
type PageNode struct {
	Mode string // ModeList or ModeDetail
	Name string // display label (the path)
	Path string // path segment used as the permission "verb"
}

// PageNodes enumerates the resource's listing and detail pages that take part
// in the permission tree. A page that has its own (developer-defined) verifier
// has its own check mechanism and is skipped; a page with automatic
// path-based permission (AutoPerm) enters the tree using its path as the
// permission segment (e.g. …:carteiras:<id>:/relatorio). Pages with no perm at
// all are also skipped.
func PageNodes(mb *presets.ModelBuilder) (out []PageNode) {
	collect := func(mode string, pages []*presets.HttpPageBuilder) {
		for _, p := range pages {
			if p.HasCustomVerifier() || !p.AutoPermEnabled() {
				continue
			}
			out = append(out, PageNode{Mode: mode, Name: p.Path(), Path: p.Path()})
		}
	}
	collect(ModeList, mb.Listing().PagesRegistrator().HttpPages())
	if mb.HasDetailing() {
		collect(ModeDetail, mb.Detailing().PagesRegistrator().HttpPages())
	}
	return
}

// PageResource returns the permission resource of a page. For a detail page it
// is the record resource with the page path appended; for a listing page it is
// the listing resource with the page path appended. This mirrors the automatic
// path-based permission the runtime builds for AutoPerm pages.
func PageResource(mb *presets.ModelBuilder, mode string, id model.ID, pagePath string) string {
	var base string
	if mode == ModeDetail {
		base = uniqueResource(mb.Permissioner().Verifier(id).On(pagePath))
	} else {
		base = uniqueResource(mb.Permissioner().ListVerifier().On(pagePath))
	}
	return base
}

// FieldResource returns the exact permission resource of a record field (or a
// nested field, given the field path), computed by the presets permissioner —
// the same computation the runtime field-permission check uses, so a policy
// stored with this resource is matched for that field.
func FieldResource(mb *presets.ModelBuilder, id model.ID, path ...string) string {
	return uniqueResource(mb.Permissioner().Verifier(id).On(mb.Permissioner().FieldPermParts(strings.Join(path, "."))...))
}
