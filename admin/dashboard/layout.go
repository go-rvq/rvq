package dashboard

import (
	"cmp"
	"crypto/rand"
	"encoding/hex"
	"slices"
	"strings"
)

// The layout of a user's dashboard (Prefs.Sections): its sections in order —
// the groups (their ids), and the user's own (Custom: a title and an icon of
// theirs) —, each with its widgets in order. A widget is in one section, or
// removed (Prefs.Removed: out of the dashboard, offered to be added back to
// its group), and may be disabled (Prefs.Disabled: in its place, shown only
// while customizing).

// Section is a section of a user's layout.
type Section struct {
	ID string `json:"id"`
	// Custom: the user's own (its Title and Icon theirs), else a group's
	Custom bool `json:"custom,omitempty"`
	// Parent is the section it is in ("" none: at the top), a tree
	Parent  string   `json:"parent,omitempty"`
	Title   string   `json:"title,omitempty"`
	Icon    string   `json:"icon,omitempty"`
	Widgets []string `json:"widgets,omitempty"`
}

// customPrefix begins the id of a section of the user's own.
const customPrefix = "s-"

// newSectionID is the id of a new section of the user's own.
func newSectionID() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return customPrefix + hex.EncodeToString(b)
}

func (p *Prefs) removed(id string) bool {
	return slices.Contains(p.Removed, id) || slices.Contains(p.Hidden, id)
}
func (p *Prefs) disabled(id string) bool { return slices.Contains(p.Disabled, id) }

// section is the section id of the layout, or nil.
func section(ss []Section, id string) *Section {
	for i := range ss {
		if ss[i].ID == id {
			return &ss[i]
		}
	}
	return nil
}

// layout is the user's layout as it is now: their sections (none yet: the
// groups, their widgets in the old order), every group a section, every
// widget known in one — a new one at the end of its group's — but the
// removed; the ids no longer known out.
func (b *Builder) layout(p *Prefs) []Section {
	known := map[string]*Widget{}
	for _, w := range b.widgets {
		known[w.ID] = w
	}
	var out []Section
	if len(p.Sections) > 0 {
		for _, s := range p.Sections {
			if !s.Custom && b.groupOf(s.ID) == nil {
				continue
			}
			s.Widgets = slices.DeleteFunc(slices.Clone(s.Widgets), func(id string) bool { return known[id] == nil || p.removed(id) })
			out = append(out, s)
		}
	} else {
		for _, g := range b.sortedGroups() {
			out = append(out, Section{ID: g.ID})
		}
	}
	// every group a section
	for _, g := range b.sortedGroups() {
		if section(out, g.ID) == nil {
			out = append(out, Section{ID: g.ID})
		}
	}
	placed := map[string]bool{}
	for _, s := range out {
		for _, id := range s.Widgets {
			placed[id] = true
		}
	}
	// the widgets not placed: in their group's, in the old order, then theirs
	pos := func(id string) int {
		if i := slices.Index(p.Order, id); i >= 0 {
			return i
		}
		return len(p.Order)
	}
	rest := slices.DeleteFunc(slices.Clone(b.widgets), func(w *Widget) bool { return placed[w.ID] || p.removed(w.ID) })
	slices.SortStableFunc(rest, func(a, c *Widget) int {
		return cmp.Or(cmp.Compare(pos(a.ID), pos(c.ID)), cmp.Compare(a.Order, c.Order))
	})
	for _, w := range rest {
		if s := section(out, w.Group); s != nil {
			s.Widgets = append(s.Widgets, w.ID)
		}
	}
	// a parent not there, or in a loop: at the top
	for i := range out {
		if out[i].Parent != "" && (section(out, out[i].Parent) == nil || isUnder(out, out[i].Parent, out[i].ID)) {
			out[i].Parent = ""
		}
	}
	return out
}

// isUnder says whether the section id is under the section anc (or is it).
func isUnder(ss []Section, id, anc string) bool {
	for seen := 0; id != "" && seen <= len(ss); seen++ {
		if id == anc {
			return true
		}
		s := section(ss, id)
		if s == nil {
			return false
		}
		id = s.Parent
	}
	return false
}

// children are the sections in the section id ("": at the top), in order.
func children(ss []Section, id string) []Section {
	var out []Section
	for _, s := range ss {
		if s.Parent == id {
			out = append(out, s)
		}
	}
	return out
}

// groupOf is the group id, or nil.
func (b *Builder) groupOf(id string) *Group {
	for _, g := range b.groups {
		if g.ID == id {
			return g
		}
	}
	return nil
}

// remove takes the widgets ids out of the layout, removed.
func (b *Builder) remove(p *Prefs, ids ...string) {
	ss := b.layout(p)
	for i := range ss {
		ss[i].Widgets = slices.DeleteFunc(ss[i].Widgets, func(id string) bool { return slices.Contains(ids, id) })
	}
	p.Sections = ss
	for _, id := range ids {
		p.Removed = toggle(p.Removed, id, true)
	}
}

// add puts the removed widgets ids back, each at the end of its group's
// section (layout places a widget in no section there).
func (b *Builder) add(p *Prefs, ids ...string) {
	for _, id := range ids {
		p.Removed = toggle(p.Removed, id, false)
		p.Hidden = toggle(p.Hidden, id, false)
	}
	p.Sections = b.layout(p)
}

// arrange makes the layout the one the user made (by dragging): each
// section's widgets as given — only the widgets known, each once; those of
// the old layout the user does not see (not allowed) kept where they were —;
// the sections in the order given, their titles the layout's.
func (b *Builder) arrange(p *Prefs, given []Section) {
	old := b.layout(p)
	seen := map[string]bool{}
	var out []Section
	for _, g := range given {
		s := section(old, g.ID)
		if s == nil {
			continue
		}
		ns := *s
		ns.Widgets = nil
		for _, id := range g.Widgets {
			if b.widget(id) == nil || seen[id] || p.removed(id) {
				continue
			}
			seen[id] = true
			ns.Widgets = append(ns.Widgets, id)
		}
		out = append(out, ns)
	}
	// the sections not given, and the widgets not given: where they were
	for _, s := range old {
		ns := section(out, s.ID)
		if ns == nil {
			out = append(out, Section{ID: s.ID, Custom: s.Custom, Parent: s.Parent, Title: s.Title, Icon: s.Icon})
			ns = &out[len(out)-1]
		}
		for _, id := range s.Widgets {
			if !seen[id] {
				seen[id] = true
				ns.Widgets = append(ns.Widgets, id)
			}
		}
	}
	p.Sections = out
}

// saveSection makes a section of the user's own (id "": a new one, last),
// or changes its title and icon; a section (its own, or a group's) put in
// parent ("" at the top; not in itself nor under itself).
func (b *Builder) saveSection(p *Prefs, id, title, icon, parent string) string {
	ss := b.layout(p)
	title = strings.TrimSpace(title)
	if !strings.HasPrefix(icon, "mdi-") {
		icon = ""
	}
	s := section(ss, id)
	switch {
	case s != nil && s.Custom:
		s.Title, s.Icon = title, icon
	case s == nil:
		id = newSectionID()
		ss = append(ss, Section{ID: id, Custom: true, Title: title, Icon: icon})
		s = &ss[len(ss)-1]
	}
	if parent == "" || (section(ss, parent) != nil && !isUnder(ss, parent, id)) {
		s.Parent = parent
	}
	p.Sections = ss
	return id
}

// deleteSection deletes the section of the user's own id: its widgets back
// to their groups' sections, its sections in its parent.
func (b *Builder) deleteSection(p *Prefs, id string) {
	ss := b.layout(p)
	s := section(ss, id)
	if s == nil || !s.Custom {
		return
	}
	moved, parent := s.Widgets, s.Parent
	for i := range ss {
		if ss[i].Parent == id {
			ss[i].Parent = parent
		}
	}
	ss = slices.DeleteFunc(ss, func(x Section) bool { return x.ID == id })
	for _, wid := range moved {
		if w := b.widget(wid); w != nil {
			if g := section(ss, w.Group); g != nil {
				g.Widgets = append(g.Widgets, wid)
			}
		}
	}
	p.Sections = ss
}

// moveSection moves the section id by delta among its siblings (the
// sections of its parent).
func (b *Builder) moveSection(p *Prefs, id string, delta int) {
	ss := b.layout(p)
	defer func() { p.Sections = ss }()
	s := section(ss, id)
	if s == nil {
		return
	}
	var idx []int
	for i := range ss {
		if ss[i].Parent == s.Parent {
			idx = append(idx, i)
		}
	}
	k := slices.IndexFunc(idx, func(i int) bool { return ss[i].ID == id })
	if j := k + delta; k >= 0 && j >= 0 && j < len(idx) {
		ss[idx[k]], ss[idx[j]] = ss[idx[j]], ss[idx[k]]
	}
}

// moveWidget moves the widget id by delta in its section.
func (b *Builder) moveWidget(p *Prefs, id string, delta int) {
	ss := b.layout(p)
	for k := range ss {
		ws := ss[k].Widgets
		i := slices.Index(ws, id)
		if i < 0 {
			continue
		}
		if j := i + delta; j >= 0 && j < len(ws) {
			ws[i], ws[j] = ws[j], ws[i]
		}
	}
	p.Sections = ss
}
