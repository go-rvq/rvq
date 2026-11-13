package presets

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/go-rvq/rvq/x/perm"
	"github.com/gobwas/glob"
)

// The sets of permissions (x/perm sets.go): a name for permissions given to
// a resource together, "admin:posts:%readers". A set has, by kind of
// resource (the kind of its node in the tree of the permissions: "model",
// "singleton", "page", "check"…), what it gives there; a resource may have
// its own for a name (HttpPageBuilder.PermSet, ModelBuilder.PermSet), which
// takes the place of the kind's — the site files' "publishers". Sets come
// from the code (RegisterPermSets: the defaults, the packages', the
// application's) and from the admin (PermSetSource: the role package's
// model); their titles are translated (their Title, Description: of the
// request's language).

// PermSet is a set of permissions.
type PermSet struct {
	Name string
	// Title and Description: in the language of the request.
	Title       func(ctx context.Context) string
	Description func(ctx context.Context) string
	// Blocks are what it gives by kind of resource ("model": "{@list,
	// <*>:@get}"), relative to the resource (perm.SetBlocks).
	Blocks map[string]string
	// System: of the code — not edited in the admin.
	System bool
}

// TitleOf is the set's title in the language of ctx: its name when it has
// none.
func (s *PermSet) TitleOf(ctx context.Context) string {
	if s.Title != nil {
		if t := s.Title(ctx); t != "" {
			return t
		}
	}
	return s.Name
}

// Definition is the set as it is written: "{model: {@list, <*>:@get},
// singleton: @get}" — its kinds in order of name.
func (s *PermSet) Definition() string {
	kinds := make([]string, 0, len(s.Blocks))
	for k := range s.Blocks {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	parts := make([]string, len(kinds))
	for i, k := range kinds {
		parts[i] = k + ": " + s.Blocks[k]
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

// DescriptionOf is the set's description in the language of ctx; "" none.
func (s *PermSet) DescriptionOf(ctx context.Context) string {
	if s.Description != nil {
		return s.Description(ctx)
	}
	return ""
}

// NewPermSet is the set name of the definition def ("{model: {@list,
// <*>:@get}, singleton: @get}"), of the code; a definition that cannot be
// read panics — it is the code's.
func NewPermSet(name, def string, title, description func(ctx context.Context) string) *PermSet {
	if !perm.ValidSetName(name) {
		panic(fmt.Sprintf("permission set %q: not a name", name))
	}
	blocks, err := perm.SetBlocks(def)
	if err != nil {
		panic(fmt.Sprintf("permission set %q: %v", name, err))
	}
	return &PermSet{Name: name, Title: title, Description: description, Blocks: blocks, System: true}
}

// resourceSet is what a set gives to one resource of its own.
type resourceSet struct {
	body               string
	title, description func(ctx context.Context) string
}

// permSets are the sets of a Builder.
type permSets struct {
	mu     sync.RWMutex
	code   []*PermSet
	source func() []*PermSet
}

// RegisterPermSets registers sets of the code, in their order — a name
// twice panics.
func (b *Builder) RegisterPermSets(sets ...*PermSet) *Builder {
	b.permSets.mu.Lock()
	defer b.permSets.mu.Unlock()
	for _, s := range sets {
		for _, o := range b.permSets.code {
			if o.Name == s.Name {
				panic(fmt.Sprintf("permission set %q registered twice", s.Name))
			}
		}
		s.System = true
		b.permSets.code = append(b.permSets.code, s)
	}
	return b
}

// PermSetSource sets where the sets of the admin come from (the role
// package's model): asked each time the sets are; one of the name of a set of
// the code is left out.
func (b *Builder) PermSetSource(f func() []*PermSet) *Builder {
	b.permSets.mu.Lock()
	defer b.permSets.mu.Unlock()
	b.permSets.source = f
	return b
}

// PermSets are the sets: the code's, in their order, then the admin's, by
// name.
func (b *Builder) PermSets() []*PermSet {
	b.permSets.mu.RLock()
	code, source := b.permSets.code, b.permSets.source
	b.permSets.mu.RUnlock()
	out := append([]*PermSet{}, code...)
	if source == nil {
		return out
	}
	names := map[string]int{}
	for i, s := range code {
		names[s.Name] = i
	}
	more := source()
	sort.SliceStable(more, func(i, j int) bool { return more[i].Name < more[j].Name })
	for _, s := range more {
		if s == nil {
			continue
		}
		if i, ok := names[s.Name]; ok {
			// a set of the code, as the admin keeps it (role): what it gives is
			// the code's; its title and description, the admin's translations
			if i < len(code) {
				merged := *code[i]
				merged.Title = firstText(s.Title, code[i].Title)
				merged.Description = firstText(s.Description, code[i].Description)
				out[i] = &merged
			}
			continue
		}
		names[s.Name] = len(out)
		out = append(out, s)
	}
	return out
}

// CodePermSets are the sets of the code (RegisterPermSets), in their order.
func (b *Builder) CodePermSets() []*PermSet {
	b.permSets.mu.RLock()
	defer b.permSets.mu.RUnlock()
	return append([]*PermSet{}, b.permSets.code...)
}

// firstText is a text of first, else of then.
func firstText(first, then func(ctx context.Context) string) func(ctx context.Context) string {
	if first == nil {
		return then
	}
	return func(ctx context.Context) string {
		if t := first(ctx); t != "" {
			return t
		}
		if then != nil {
			return then(ctx)
		}
		return ""
	}
}

// PermSet is the set name; nil none.
func (b *Builder) PermSet(name string) *PermSet {
	for _, s := range b.PermSets() {
		if s.Name == name {
			return s
		}
	}
	return nil
}

// PermSet gives name, for this page, what body says (relative to it:
// "!{create, edit}", "{%authors, !publish}"), in the place of what the set
// gives a page; title and description, when the name is none of the sets':
// its own.
func (b *HttpPageBuilder) PermSet(name, body string, title, description func(ctx context.Context) string) *HttpPageBuilder {
	b.permSets = addResourceSet(b.permSets, name, body, title, description)
	return b
}

// PermSet gives name, for this model (its listing), what body says, as
// HttpPageBuilder.PermSet.
func (mb *ModelBuilder) PermSet(name, body string, title, description func(ctx context.Context) string) *ModelBuilder {
	mb.permSets = addResourceSet(mb.permSets, name, body, title, description)
	return mb
}

func addResourceSet(m map[string]*resourceSet, name, body string, title, description func(ctx context.Context) string) map[string]*resourceSet {
	if !perm.ValidSetName(name) {
		panic(fmt.Sprintf("permission set %q: not a name", name))
	}
	if err := perm.CheckSetBlock(body); err != nil {
		panic(fmt.Sprintf("permission set %q: %v", name, err))
	}
	if m == nil {
		m = map[string]*resourceSet{}
	}
	m[name] = &resourceSet{body: body, title: title, description: description}
	return m
}

// PermSetOffer is a set that gives something to a resource: its name, its
// title and description, and the lines it gives there.
type PermSetOffer struct {
	Name        string
	Title       string
	Description string
	Lines       []string
	// Own: the resource's own (HttpPageBuilder.PermSet), not the set's of
	// its kind.
	Own bool
}

// PermSetsOf are the sets that give something to the node n — its own, then
// the others by their order —, the lines of each there (n's resource:
// Unique, else Name), in the language of ctx.
func (b *Builder) PermSetsOf(ctx context.Context, n *PermNode) (out []*PermSetOffer) {
	if n == nil {
		return nil
	}
	res := n.Unique
	if res == "" {
		res = n.Name
	}
	seen := map[string]bool{}
	add := func(name string, title, description func(ctx context.Context) string, own bool) {
		if seen[name] {
			return
		}
		lines, ok := b.setLinesOf(n, res, name, map[string]bool{})
		if !ok {
			return
		}
		seen[name] = true
		o := &PermSetOffer{Name: name, Lines: lines, Own: own}
		if title != nil {
			o.Title = title(ctx)
		}
		if description != nil {
			o.Description = description(ctx)
		}
		if s := b.PermSet(name); s != nil {
			if o.Title == "" {
				o.Title = s.TitleOf(ctx)
			}
			if o.Description == "" {
				o.Description = s.DescriptionOf(ctx)
			}
		}
		if o.Title == "" {
			o.Title = name
		}
		out = append(out, o)
	}
	names := make([]string, 0, len(n.Sets))
	for name := range n.Sets {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		add(name, n.Sets[name].title, n.Sets[name].description, true)
	}
	for _, s := range b.PermSets() {
		add(s.Name, s.Title, s.Description, false)
	}
	return
}

// permIndex is the tree of the permissions by resource: by the groups and by
// the unique names.
type permSetIndex struct {
	byName map[string]*PermNode
	all    []*PermNode
}

func newPermSetIndex(root *PermNode) *permSetIndex {
	ix := &permSetIndex{byName: map[string]*PermNode{}}
	if root == nil {
		return ix
	}
	root.Walk(func(_ []*PermNode, n *PermNode) {
		ix.all = append(ix.all, n)
		if n.Name != "" {
			ix.byName[n.Name] = n
		}
		if n.Unique != "" {
			ix.byName[n.Unique] = n
		}
	})
	return ix
}

// nodes are the nodes of resource (a pattern with its ":" at its end), each
// with the resource the lines are given to: the one written, for a resource
// named; each node's (its unique name, else its name), for a pattern of a
// glob.
func (ix *permSetIndex) nodes(resource string) (nodes []*PermNode, resources []string) {
	if n := ix.byName[resource]; n != nil {
		return []*PermNode{n}, []string{resource}
	}
	if !strings.ContainsAny(resource, "*?[{") {
		return nil, nil
	}
	g, err := glob.Compile(resource)
	if err != nil {
		return nil, nil
	}
	seen := map[*PermNode]bool{}
	for _, n := range ix.all {
		if seen[n] || n.Kind == PermNodeGroup {
			continue
		}
		res := n.Unique
		if res == "" {
			res = n.Name
		}
		if g.Match(n.Name) || (n.Unique != "" && g.Match(n.Unique)) {
			seen[n] = true
			nodes = append(nodes, n)
			resources = append(resources, res)
		}
	}
	return
}

// expandSet is the x/perm expander of the sets (perm.SetSetExpander): what
// the set name gives to resource — to each node of it, by its own block for
// the name, else the set's of its kind; nothing, not ok.
func (b *Builder) expandSet(resource, name string) ([]string, bool) {
	ix := b.permSetIx
	if ix == nil {
		return nil, false
	}
	nodes, resources := ix.nodes(resource)
	var out []string
	ok := false
	for i, n := range nodes {
		if lines, got := b.setLinesOf(n, resources[i], name, map[string]bool{}); got {
			ok = true
			out = append(out, lines...)
		}
	}
	return out, ok
}

// maxSetDepth is how deep sets in sets go.
const maxSetDepth = 8

// setLinesOf are the lines the set name gives to the node n, its resource
// res: its own block for the name, else the set's of its kind; seen, the
// sets on the way (a set of itself gives nothing more).
func (b *Builder) setLinesOf(n *PermNode, res, name string, seen map[string]bool) ([]string, bool) {
	if seen[name] || len(seen) >= maxSetDepth {
		return nil, false
	}
	var body string
	if rs := n.Sets[name]; rs != nil {
		body = rs.body
	} else if s := b.PermSet(name); s != nil {
		body = s.Blocks[string(n.Kind)]
	}
	if body == "" {
		return nil, false
	}
	next := map[string]bool{name: true}
	for k := range seen {
		next[k] = true
	}
	lines, err := perm.SetBlockLines(res, body, func(sub string) []string {
		l, _ := b.setLinesOf(n, res, sub, next)
		return l
	})
	if err != nil {
		return nil, false
	}
	return lines, len(lines) > 0
}

// defaultPermSets are the sets every admin has: to see, to author, to remove,
// to edit (both), to manage (all).
func defaultPermSets() []*PermSet {
	msg := func(f func(m *Messages) string) func(ctx context.Context) string {
		return func(ctx context.Context) string { return f(MustGetMessages(ctx)) }
	}
	return []*PermSet{
		NewPermSet("readers", "{model: {@list, <*>:@get}, singleton: @get, page: @get}",
			msg(func(m *Messages) string { return m.PermSetReaders }), msg(func(m *Messages) string { return m.PermSetReadersDesc })),
		NewPermSet("authors", "{model: {@list, @create, <*>:{@get, @edit}}, singleton: {@get, @edit}}",
			msg(func(m *Messages) string { return m.PermSetAuthors }), msg(func(m *Messages) string { return m.PermSetAuthorsDesc })),
		NewPermSet("removers", "{model: {@trash, <*>:{@delete, @restore}}}",
			msg(func(m *Messages) string { return m.PermSetRemovers }), msg(func(m *Messages) string { return m.PermSetRemoversDesc })),
		NewPermSet("editors", "{model: {%authors, %removers}, singleton: %authors}",
			msg(func(m *Messages) string { return m.PermSetEditors }), msg(func(m *Messages) string { return m.PermSetEditorsDesc })),
		NewPermSet("managers", "{model: *, singleton: *, page: *, check: *}",
			msg(func(m *Messages) string { return m.PermSetManagers }), msg(func(m *Messages) string { return m.PermSetManagersDesc })),
	}
}
