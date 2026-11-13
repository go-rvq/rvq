package userdocs

import (
	"context"
	"sort"
	"strings"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

// The parts of a permission, for a field that writes one
// (PermissionInput, vuetifyx's vx-permission-input): what may follow a
// beginning of it (PermChoicesOf) — by the tree of the permissions of the
// admin (presets.Builder.BuildPermissions) —, the records of a model
// (PermRecordsOf), and the labels of the records of a line, as they are now
// (ResolveLabels).
//
// A beginning is the parts of a permission up to a ":" — "admin:posts:<7>:"
// —, by the unique names ("admin:posts:") or by the groups of the menu
// ("admin:site/:content/:posts:"); "" is the root, whose parts are the
// scopes ("admin:").

// More kinds of the parts offered (PermChoice.Kind), beside the tree's.
const (
	PermKindRecord = "record" // a record, <7>; all of them, <*>
	PermKindAll    = "all"    // everything under it, *
)

// PermChoice is a part that may follow a beginning of a permission.
type PermChoice struct {
	// Segment is the part as written: "posts:", "@list", "<*>:", "*"; a
	// ":" at its end when parts may follow it
	Segment     string `json:"segment"`
	Kind        string `json:"kind"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	// Children: parts may follow it
	Children bool `json:"children,omitempty"`
	// Active: of a stack (PermChoicesOf's allows), whether it allows what is
	// under the part — PermActiveAll, PermActiveSome, PermActiveNone
	Active string `json:"active,omitempty"`
	// Resource is what choosing it writes: the beginning and the part, a
	// resource by its unique name, out of the groups of the menu
	// (perm.UngroupPattern) — a group by the groups
	Resource string `json:"resource"`
	// Lines: of a set (PermKindSet, "%readers"), the lines it gives there,
	// after the resource ("@list", "<*>:@get")
	Lines []string `json:"lines,omitempty"`
}

// PermChoices are the parts that may follow a beginning, and, of a model's
// listing, Records: its records may be chosen (PermRecordsOf of the same
// beginning).
type PermChoices struct {
	Choices []PermChoice `json:"choices"`
	Records bool         `json:"records,omitempty"`
}

// permIndex is the tree of the permissions by the beginnings of its nodes,
// by the groups and by the unique names.
type permIndex struct {
	root   *presets.PermNode
	byName map[string]*presets.PermNode
	scopes []string
}

func newPermIndex(p *presets.Builder) *permIndex {
	ix := &permIndex{root: p.BuildPermissions().Tree(), byName: map[string]*presets.PermNode{}}
	seen := map[string]bool{}
	ix.root.Walk(func(_ []*presets.PermNode, n *presets.PermNode) {
		if n.Name != "" {
			ix.byName[n.Name] = n
			if s := firstSegment(n.Name); s != "" && !seen[s] {
				seen[s] = true
				ix.scopes = append(ix.scopes, s)
			}
		}
		if n.Unique != "" {
			ix.byName[n.Unique] = n
		}
	})
	return ix
}

// firstSegment is the first part of a resource, with its ":".
func firstSegment(s string) string {
	if i := strings.Index(s, ":"); i > 0 {
		return s[:i+1]
	}
	return ""
}

// permPath is a beginning read: its key in the index (its records <*>), and
// the records in it, each with its model.
type permPath struct {
	key     string
	records []pathRecord
}

type pathRecord struct {
	model *presets.ModelBuilder
	id    string
}

// readPath reads the beginning path (its labels none of it).
func (ix *permIndex) readPath(path string) (pp permPath, err error) {
	e, err := perm.ParseExpr(perm.StripLabels(path))
	if err != nil {
		return pp, err
	}
	var b strings.Builder
	for _, n := range e.Nodes {
		switch n := n.(type) {
		case *perm.ExprText:
			b.WriteString(n.Text)
		case *perm.ExprRecord:
			if node := ix.byName[b.String()]; node != nil && node.Model != nil && n.ID != "*" {
				pp.records = append(pp.records, pathRecord{node.Model, n.ID})
			}
			b.WriteString("<*>")
		default:
			// a group: no beginning of one part
			return pp, &perm.ExprError{Pos: n.Position().Start, Msg: "a beginning of no group"}
		}
	}
	pp.key = b.String()
	return
}

// parents are the ids of the records of the path, the outermost first, of
// the model of each.
func (pp permPath) parents() (ids model.IDSlice) {
	for _, r := range pp.records {
		id, err := r.model.ParseRecordID(r.id)
		if err != nil {
			return nil
		}
		ids = append(ids, id)
	}
	return
}

// segment is the part of the child c after its parent's beginning prefix,
// by the unique names when the prefix is one.
func segment(c *presets.PermNode, prefix string, byUnique bool) string {
	if byUnique && c.Unique != "" && strings.HasPrefix(c.Unique, prefix) {
		return c.Unique[len(prefix):]
	}
	if strings.HasPrefix(c.Name, prefix) {
		return c.Name[len(prefix):]
	}
	return ""
}

// PermChoicesOf are the parts that may follow the beginning path, in the
// language of ctx; allows, when not nil, gives what a stack allows: the parts
// of which it allows nothing left out (an access key: what its user has).
// groups: the groups of the menu offered, each resource in its group (a
// scope: its groups and the resources out of any group); else none, the
// resources of a scope by their unique names.
func PermChoicesOf(ctx context.Context, p *presets.Builder, path string, groups bool, allows func(PermCheck) bool) (r PermChoices, err error) {
	defer func() {
		for i := range r.Choices {
			r.Choices[i].Resource = perm.UngroupPattern(path + r.Choices[i].Segment)
		}
	}()
	ix := newPermIndex(p)
	t := &permTree{ctx: ctx, msgs: GetMessages(ctx)}
	if strings.TrimSpace(path) == "" {
		for _, s := range ix.scopes {
			r.Choices = append(r.Choices, ix.scopeChoice(t, s, allows))
		}
		return
	}
	pp, err := ix.readPath(path)
	if err != nil {
		return r, err
	}
	var nodes []*presets.PermNode
	var verbs []*presets.PermNodeAction
	var owner *presets.PermNode
	byUnique := false
	if n := ix.byName[pp.key]; n != nil {
		owner, nodes, verbs = n, n.Children, n.Actions
		byUnique = n.Unique != "" && n.Unique == pp.key && n.Unique != n.Name
		r.Records = n.Model != nil && (n.Kind == presets.PermNodeModel)
	} else if containsString(ix.scopes, pp.key) && groups {
		// the scope, by the menu: its groups and what is out of any
		for _, c := range ix.root.Children {
			if strings.HasPrefix(c.Name, pp.key) {
				nodes = append(nodes, c)
			}
		}
	} else if containsString(ix.scopes, pp.key) {
		// the scope, no group: its models and pages by their unique names
		ix.root.Walk(func(_ []*presets.PermNode, n *presets.PermNode) {
			if n.Kind == presets.PermNodeGroup || n.Unique == "" || !strings.HasPrefix(n.Unique, pp.key) {
				return
			}
			if rest := strings.TrimSuffix(n.Unique[len(pp.key):], ":"); rest != "" && !strings.Contains(rest, ":") {
				nodes = append(nodes, n)
			}
		})
		byUnique = true
	} else {
		return r, nil
	}
	label := ""
	if owner != nil {
		label = t.title(owner)
	}
	for _, a := range verbs {
		c := PermChoice{Segment: a.Name, Kind: verbKind(a.Name)}
		if owner.Kind == presets.PermNodePage {
			c.Title = pageVerbLabel(t.msgs, a, ctx)
		} else {
			c.Title = verbLabel(t.msgs, a, ctx)
		}
		if a.Description != nil {
			c.Description = a.Description(ctx)
		}
		if allows != nil {
			c.Active = activeOf([]PermCheck{t.check(owner, a.Name)}, allows)
		}
		r.Choices = append(r.Choices, c)
	}
	seen := map[string]bool{}
	for _, n := range nodes {
		prefix := pp.key
		if owner != nil {
			prefix = owner.Name
			if byUnique {
				prefix = owner.Unique
			}
		}
		seg := segment(n, prefix, byUnique)
		if seg == "" || seen[seg] {
			continue
		}
		seen[seg] = true
		title := t.title(n)
		c := PermChoice{Segment: seg, Kind: kindOf(n), Title: title, Description: t.describe(n, title),
			Children: len(n.Children) > 0 || len(n.Actions) > 0}
		if n.Kind == presets.PermNodeRecord {
			c.Kind, c.Title = PermKindRecord, t.msgs.PermInputAllRecords
			c.Description = strings.TrimSpace(t.msgs.PermInputAllRecords + " · " + label)
		}
		if allows != nil {
			c.Active = activeOf(checksUnder(t, n), allows)
		}
		r.Choices = append(r.Choices, c)
	}
	if owner == nil || len(owner.Children) > 0 || len(owner.Actions) > 0 {
		c := PermChoice{Segment: "*", Kind: PermKindAll, Title: t.msgs.PermInputEverything,
			Description: t.msgs.PermInputEverythingDesc}
		if allows != nil {
			var checks []PermCheck
			if owner != nil {
				checks = checksUnder(t, owner)
			} else {
				for _, n := range nodes {
					checks = append(checks, checksUnder(t, n)...)
				}
			}
			c.Active = activeOf(checks, allows)
		}
		r.Choices = append(r.Choices, c)
	}
	// the sets that give something there: their name, what they give
	if owner != nil {
		res := owner.Unique
		if res == "" {
			res = owner.Name
		}
		for _, o := range p.PermSetsOf(ctx, owner) {
			c := PermChoice{Segment: perm.SetPrefix + o.Name, Kind: PermKindSet, Title: o.Title, Description: o.Description,
				Lines: relativeLines(res, o.Lines)}
			if allows != nil {
				// of what it gives there, what the stack allows
				given := (&perm.Stack{}).Add("set", o.Lines...).Resolve()
				var checks []PermCheck
				for _, ch := range checksUnder(t, owner) {
					if given.Allows(ch.Resource, ch.Preferred, ch.Action) {
						checks = append(checks, ch)
					}
				}
				if c.Active = activeOf(checks, allows); c.Active == "" {
					c.Active = PermActiveNone
				}
			}
			r.Choices = append(r.Choices, c)
		}
	}
	if allows != nil {
		out := r.Choices[:0]
		for _, c := range r.Choices {
			if c.Active != PermActiveNone {
				out = append(out, c)
			}
		}
		r.Choices = out
	}
	return
}

// relativeLines are the lines, each after res when it is under it ("deny
// <*>:#Secret:*" of "deny admin:posts:<*>:#Secret:*").
func relativeLines(res string, lines []string) []string {
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		k, p := perm.ParseLine(l)
		out = append(out, k.Prefix()+strings.TrimPrefix(p, res))
	}
	return out
}

// scopeChoice is the part of the scope s ("admin:").
func (ix *permIndex) scopeChoice(t *permTree, s string, allows func(PermCheck) bool) PermChoice {
	c := PermChoice{Segment: s, Kind: PermKindScope, Title: strings.TrimSuffix(s, ":"), Children: true}
	if s == presets.PermModule+":" {
		c.Title, c.Description = t.msgs.PermScopeAdmin, t.msgs.PermScopeAdminDesc
	}
	if allows != nil {
		var checks []PermCheck
		for _, n := range ix.root.Children {
			if strings.HasPrefix(n.Name, s) {
				checks = append(checks, checksUnder(t, n)...)
			}
		}
		c.Active = activeOf(checks, allows)
	}
	return c
}

// checksUnder are the permissions asked of n and of what is under it.
func checksUnder(t *permTree, n *presets.PermNode) (out []PermCheck) {
	for _, a := range n.Actions {
		out = append(out, t.check(n, a.Name))
	}
	if len(n.Actions) == 0 && n.Kind != presets.PermNodeGroup && n.Kind != presets.PermNodeRecord {
		out = append(out, t.check(n, ""))
	}
	for _, c := range n.Children {
		out = append(out, checksUnder(t, c)...)
	}
	return
}

// activeOf is what of checks allows allows: PermActiveAll, PermActiveSome,
// PermActiveNone ("" none to check).
func activeOf(checks []PermCheck, allows func(PermCheck) bool) string {
	var yes, no int
	for _, c := range checks {
		if allows(c) {
			yes++
		} else {
			no++
		}
	}
	switch {
	case yes > 0 && no == 0:
		return PermActiveAll
	case yes > 0:
		return PermActiveSome
	case no > 0:
		return PermActiveNone
	}
	return ""
}

func containsString(l []string, s string) bool {
	for _, e := range l {
		if e == s {
			return true
		}
	}
	return false
}

// PermRecordsLimit is how many records PermRecordsOf answers.
const PermRecordsLimit = 20

// PermRecordsOf are the records of the model of the beginning path (its
// listing, "admin:posts:"), of the records in it (a nested model), whose
// title or id has q: those the user of ctx may read.
func PermRecordsOf(ctx *web.EventContext, p *presets.Builder, path, q string) ([]presets.PermRecord, error) {
	ix := newPermIndex(p)
	pp, err := ix.readPath(path)
	if err != nil {
		return nil, err
	}
	n := ix.byName[pp.key]
	if n == nil || n.Model == nil || n.Kind != presets.PermNodeModel {
		return nil, nil
	}
	return n.Model.PermRecords(ctx, pp.parents(), strings.TrimSpace(q), PermRecordsLimit)
}

// PermLabelError is a record of a line not found: where it is, its id.
type PermLabelError struct {
	Pos perm.Pos `json:"pos"`
	ID  string   `json:"id"`
}

// ResolveLabels is the line with the labels of its records as they are now
// — "<7>" or "<7<an old title>>" is "<7<its title>>" —, in the language of
// ctx, the records of nested models too; the records not found (or the user
// may not read) are as they are, in errs. A line that cannot be read is as
// it is.
func ResolveLabels(ctx *web.EventContext, p *presets.Builder, line string) (out string, errs []PermLabelError) {
	e, err := perm.ParseExpr(line)
	if err != nil {
		return line, nil
	}
	ix := newPermIndex(p)
	e.Records(func(path []perm.ExprNode, r *perm.ExprRecord) {
		if r.ID == "*" || strings.ContainsAny(r.ID, "*?[") {
			return
		}
		pp, err := ix.readPath(perm.NodesText(path))
		if err != nil {
			return
		}
		n := ix.byName[pp.key]
		if n == nil || n.Model == nil {
			errs = append(errs, PermLabelError{Pos: r.Pos, ID: r.ID})
			return
		}
		label, ok := n.Model.PermRecordLabel(ctx, pp.parents(), r.ID)
		if !ok {
			errs = append(errs, PermLabelError{Pos: r.Pos, ID: r.ID})
			return
		}
		r.Label = label
	})
	sort.SliceStable(errs, func(i, j int) bool { return errs[i].Pos.Start < errs[j].Pos.Start })
	return e.String(), errs
}
