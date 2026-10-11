package schemaform

import (
	"fmt"
	"strings"
)

// WidthUnits are the units a row of a group is divided in: a field's
// `[width=N]` is N of them, 1 to WidthUnits — 1 is 20% —; the fields that
// say none share what is left, equally.
const WidthUnits = 5

// The kinds of a group, by its metadata: `[tabs] { … }` each field a tab,
// `[steps] { … }` each field a step (with next and back; `[steps,
// vertical]` down the page), none a row of columns.
const (
	GroupRow   = "row"
	GroupTabs  = "tabs"
	GroupSteps = "steps"
)

// Group is a group of a record: `{ … }` in its body, the field Name (`$N`).
type Group struct {
	Name string
	Kind string
	// Vertical lays steps down the page (`[steps, vertical]`); across, else.
	Vertical bool
	// Meta is the group's own metadata (written before its block).
	Meta Meta
}

// groupOf is the group f — the field `$N` — is, by its metadata.
func groupOf(f *Field) *Group {
	g := &Group{Name: f.Name, Kind: GroupRow, Meta: f.Meta}
	switch {
	case flag(f.Meta, GroupTabs):
		g.Kind = GroupTabs
	case flag(f.Meta, GroupSteps):
		g.Kind = GroupSteps
		g.Vertical = flag(f.Meta, "vertical")
	}
	return g
}

// flag says the metadata m says key, and not false: `[tabs]`, `[tabs=true]`.
func flag(m Meta, key string) bool {
	v, ok := m[key]
	if !ok {
		return false
	}
	if b, isBool := v.(bool); isBool {
		return b
	}
	return v != nil
}

// appendField adds f to s. A group — `{ … }` in the body, the field `$N` of
// a record — is not a field of the form: its fields are, in its place, each
// a column of one row (Field.Row) at its share of it (Field.Percent); what
// they hold is the record's own (the group adds no name to the value). A
// group inside a group is refused: a row has columns, not rows.
func (r *reader) appendField(s *Schema, f *Field) error {
	if !isGroupName(f.Name) || f.Schema == nil || f.Schema.Slice || f.Schema.Choice {
		s.Fields = append(s.Fields, f)
		return nil
	}
	for _, gf := range f.Schema.Fields {
		if gf.Row != 0 {
			return fmt.Errorf("schemaform: the group %s has a group inside: a group is a row of columns, one level", f.Name)
		}
	}
	g := groupOf(f)
	// the widths share a row; a tab, a step, has the whole of it
	if g.Kind == GroupRow {
		if err := groupWidths(f.Name, f.Schema.Fields); err != nil {
			return err
		}
	}
	row := maxRow(s.Fields) + 1
	for _, gf := range f.Schema.Fields {
		gf.Row, gf.Group = row, g
		s.Fields = append(s.Fields, gf)
	}
	return nil
}

// isGroupName says name is a group's: `$N` (gad's `{ … }` in a body).
func isGroupName(name string) bool {
	if len(name) < 2 || name[0] != '$' {
		return false
	}
	for _, r := range name[1:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// Rows are fields as they are drawn: a field alone a row of its own, the
// fields of a group (Field.Row) one row of them.
func Rows(fields []*Field) (rows [][]*Field) {
	for _, f := range fields {
		if n := len(rows); f.Row != 0 && n > 0 && rows[n-1][0].Row == f.Row {
			rows[n-1] = append(rows[n-1], f)
			continue
		}
		rows = append(rows, []*Field{f})
	}
	return
}

// appendRows adds fields — a parent's, read on their own — to s, their rows
// numbered after the ones s has.
func appendRows(s *Schema, fields []*Field) {
	offset := maxRow(s.Fields)
	for _, f := range fields {
		if f.Row != 0 {
			f.Row += offset
		}
		s.Fields = append(s.Fields, f)
	}
}

func maxRow(fields []*Field) (n int) {
	for _, f := range fields {
		if f.Row > n {
			n = f.Row
		}
	}
	return
}

// groupWidths gives each field of the group its share of the row: its
// `[width=N]`, N of WidthUnits; the ones that say none, what is left,
// equally. The widths given are whole units, 1 to WidthUnits, and the row is
// whole: they add up to WidthUnits, every field with a share.
func groupWidths(group string, fields []*Field) error {
	given, rest := 0, 0
	for _, f := range fields {
		if _, ok := f.Meta["width"]; !ok {
			rest++
			continue
		}
		n, ok, err := f.Meta.Int("width")
		if err != nil || !ok || n < 1 || n > WidthUnits {
			return fmt.Errorf("schemaform: group %s, field %s: width is a whole number from 1 to %d", group, f.Name, WidthUnits)
		}
		given += n
	}
	names := func() string {
		out := make([]string, len(fields))
		for i, f := range fields {
			out[i] = f.Name
		}
		return strings.Join(out, ", ")
	}
	switch {
	case given > WidthUnits:
		return fmt.Errorf("schemaform: group %s (%s): the widths add up to %d, more than %d", group, names(), given, WidthUnits)
	case rest == 0 && given != WidthUnits:
		return fmt.Errorf("schemaform: group %s (%s): the widths add up to %d, not %d", group, names(), given, WidthUnits)
	case rest > 0 && given == WidthUnits:
		return fmt.Errorf("schemaform: group %s (%s): the widths given take the whole row, and %d field(s) have none left", group, names(), rest)
	}
	share := float64(WidthUnits-given) / float64(max(rest, 1))
	for _, f := range fields {
		units := share
		if n, ok, _ := f.Meta.Int("width"); ok {
			units = float64(n)
		}
		f.Percent = units * 100 / WidthUnits
	}
	return nil
}
