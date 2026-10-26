package presets

import (
	"errors"
	"fmt"
	"strings"
)

// A record path says "this record, under these parents" as text, to be kept
// in a field (the origin of a mail sent, of a notification) and turned into
// the admin's address later: the models from the root to the record, joined
// by "/", each "<model id>:<record id>", or "<model id>" for a singleton.
//
//	posts:01a0…_en-US
//	surveys:1/Places:2/Products:3
//	settings
//
// The model id is the model builder's (MenuID, what GetModelByID and
// GetChildByID find); the record id is ID.String, path-escaped: a "/" in it
// is "%2F", and a ":" in it is no separator, the model id ending at the
// first one.

// ErrRecordPath is the error of a record path that names no model, or a
// record its model cannot parse.
var ErrRecordPath = errors.New("record path")

// RecordPathSegment is a model of a record path and the id of its record
// (zero for a singleton).
type RecordPathSegment struct {
	Model *ModelBuilder
	ID    ID
}

// RecordPath is the path of the record id of the model under the records
// parents of its parents, root first — one for each parent, zero for a
// singleton parent, as ListingHref takes them; id is zero for a singleton.
func (i *ModelInfo) RecordPath(parents IDSlice, id ID) string {
	return i.mb.RecordPath(parents, id)
}

// RecordPath is ModelInfo.RecordPath.
func (mb *ModelBuilder) RecordPath(parents IDSlice, id ID) string {
	models := append(mb.Parents(), mb)
	segments := make([]string, len(models))
	for n, m := range models {
		var v ID
		switch {
		case n == len(models)-1:
			v = id
		case n < len(parents):
			v = parents[n]
		}
		segments[n] = m.id
		if !m.singleton {
			segments[n] += ":" + v.String()
		}
	}
	return strings.Join(segments, "/")
}

// ParseRecordPath is the models and the record ids of the path p, root
// first. The first model is a root one (GetModelByID among those without a
// parent), each next a child of the one before (GetChildByID); each id is
// parsed by its model (ParseRecordID), so a bad id is an error, not a broken
// link.
func (b *Builder) ParseRecordPath(p string) (segments []RecordPathSegment, err error) {
	if p == "" {
		return nil, fmt.Errorf("%w: empty", ErrRecordPath)
	}
	var parent *ModelBuilder
	for s := range strings.SplitSeq(p, "/") {
		name, rawID, hasID := strings.Cut(s, ":")
		var mb *ModelBuilder
		if parent == nil {
			if mb = b.GetModelByID(name); mb != nil && mb.parent != nil {
				mb = nil
			}
		} else {
			mb = parent.GetChildByID(name)
		}
		if mb == nil {
			return nil, fmt.Errorf("%w %q: no model %q", ErrRecordPath, p, name)
		}
		seg := RecordPathSegment{Model: mb}
		switch {
		case mb.singleton && hasID:
			return nil, fmt.Errorf("%w %q: %q is a singleton, with no record id", ErrRecordPath, p, name)
		case !mb.singleton && (!hasID || rawID == ""):
			return nil, fmt.Errorf("%w %q: %q has no record id", ErrRecordPath, p, name)
		case !mb.singleton:
			if seg.ID, err = mb.ParseRecordID(rawID); err != nil {
				return nil, fmt.Errorf("%w %q: the id of %q: %w", ErrRecordPath, p, name, err)
			}
		}
		segments = append(segments, seg)
		parent = mb
	}
	return
}

// RecordPathHref is the admin's address of the record of the path p: the
// detail of the last model's record, or the page of the singleton.
func (b *Builder) RecordPathHref(p string) (string, error) {
	segments, err := b.ParseRecordPath(p)
	if err != nil {
		return "", err
	}
	last := segments[len(segments)-1]
	parents := make(IDSlice, len(segments)-1)
	for n, s := range segments[:len(segments)-1] {
		parents[n] = s.ID
	}
	info := last.Model.Info()
	if last.Model.singleton {
		return info.ListingHref(parents...), nil
	}
	return info.DetailingHref(last.ID, parents...), nil
}
