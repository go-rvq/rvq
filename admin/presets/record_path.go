package presets

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-rvq/rvq/admin/model"
)

// A record path says "this record, under these parents" — or this listing,
// this page, this action — as text, to be kept in a field (the origin of a
// mail sent, of a notification) and turned into the admin's address later.
// It is the resource of the permissions by the unique names, the groups of
// the menu left out (perm.Verifier.PreferredResource), made by the same
// generator (ModelPermissioner), with the whole id of each record in its
// "<…>" (ID.String, path-escaped: a ">" in it is "%3E"):
//
//	admin:posts:                          the listing
//	admin:posts:<01a0…_en-US>:            a record
//	admin:surveys:<1>:places:<2>:         a record under its parent's
//	admin:mail_sender:                    a singleton
//	admin:posts:<01a0…_en-US>:/preview:   a page of the record
//	admin:posts:/import:                  a page of the listing
//	admin:/site-files:                    a page of the admin
//	admin:mail_sender:!test_send_mail     an action (ActionPerm)
//
// An action is of the record (or singleton) — its detail's —, of the
// listing, or of the page (PageBuilder.Action) it follows. Its address is the
// page's with the suffix "/action:<name>" (ActionSuffix): the page of the
// action, its form in the body (actionPage).

// ErrRecordPath is the error of a record path that names no model, page or
// action, or a record its model cannot parse.
var ErrRecordPath = errors.New("record path")

// RecordPathSegment is a model of a record path and the id of its record
// (zero for a singleton, or for the listing).
type RecordPathSegment struct {
	Model *ModelBuilder
	ID    ID
}

// RecordPathTarget is what a record path names: the models from the root
// (none for a page of the admin), the page, the action.
type RecordPathTarget struct {
	Models []RecordPathSegment
	// Listing says the last model is named without a record: its listing
	Listing bool
	Page    *HttpPageBuilder
	Action  *ActionBuilder
}

// pathID is id as the generator of the permissions takes it — the value of
// "ID" — holding the whole id, path-escaped.
func pathID(id ID) ID {
	if id.IsZero() {
		return ID{}
	}
	return ID{Fields: []model.Field{model.SingleField("ID")}, Values: []any{id.String()}}
}

// RecordPath is the path of the record id of the model under the records
// parents of its parents, root first — one for each parent, zero for a
// singleton parent, as ListingHref takes them —; of the singleton, id zero.
// "" with no permissions builder (Builder.Permission): there is no generator.
func (i *ModelInfo) RecordPath(parents IDSlice, id ID) string {
	return i.mb.RecordPath(parents, id)
}

// ListingPath is the path of the listing of the model under the records
// parents of its parents.
func (i *ModelInfo) ListingPath(parents IDSlice) string {
	return i.mb.ListingPath(parents)
}

// RecordPath is ModelInfo.RecordPath.
func (mb *ModelBuilder) RecordPath(parents IDSlice, id ID) string {
	return mb.permissioner.Verifier(pathID(id), mb.permParents(parents)...).PreferredResource()
}

// ListingPath is ModelInfo.ListingPath.
func (mb *ModelBuilder) ListingPath(parents IDSlice) string {
	return mb.permissioner.ListVerifier(mb.permParents(parents)...).PreferredResource()
}

// permParents are the ids of the parents as the permissioner takes them:
// only the ones of the parents that are not singletons, the nearest first.
func (mb *ModelBuilder) permParents(parents IDSlice) (r IDSlice) {
	for n, p := range mb.Parents() {
		if p.singleton {
			continue
		}
		var id ID
		if n < len(parents) {
			id = parents[n]
		}
		r = append(IDSlice{pathID(id)}, r...)
	}
	return
}

// PagePath is the path of the page pagePath ("/preview") of what base is
// the path of: a record, a listing.
func PagePath(base, pagePath string) string {
	return base + pagePath + ":"
}

// ActionPath is the path of the action name of what base is the path of: a
// record, a listing, a page.
func ActionPath(base, name string) string {
	return base + ActionPerm(name)
}

// PagePath is the path of the page of the admin pagePath ("/site-files").
func (b *Builder) PagePath(pagePath string) string {
	return b.verifier.Spawn().Prefer(PermModule, pagePath).PreferredResource()
}

// splitRecordPath is the parts of p — the ":" in a "<…>" is no separator —
// and the action at its end ("!test_send_mail"), if any.
func splitRecordPath(p string) (parts []string, action string, err error) {
	if i := strings.LastIndex(p, ":"); i >= 0 && strings.HasPrefix(p[i+1:], "!") {
		p, action = p[:i+1], p[i+1:]
	}
	if !strings.HasSuffix(p, ":") {
		return nil, "", fmt.Errorf("%w %q: no final \":\"", ErrRecordPath, p)
	}
	for p != "" {
		var part string
		if strings.HasPrefix(p, "<") {
			end := strings.Index(p, ">:")
			if end < 0 {
				return nil, "", fmt.Errorf("%w %q: an id with no end", ErrRecordPath, p)
			}
			part, p = p[:end+1], p[end+2:]
		} else {
			end := strings.Index(p, ":")
			part, p = p[:end], p[end+1:]
		}
		parts = append(parts, part)
	}
	return
}

// ParseRecordPath is what the path p names. Each model is found by its
// unique name (UniquePermName) — a root one, or a child of the one before —,
// each id parsed by its model (ParseRecordID), so a bad id is an error, not
// a broken link.
func (b *Builder) ParseRecordPath(p string) (t *RecordPathTarget, err error) {
	parts, action, err := splitRecordPath(p)
	if err != nil {
		return nil, err
	}
	bad := func(format string, args ...any) error {
		return fmt.Errorf("%w %q: %s", ErrRecordPath, p, fmt.Sprintf(format, args...))
	}
	if len(parts) < 2 || parts[0] != PermModule {
		return nil, bad("not of the admin (%s:)", PermModule)
	}
	t = &RecordPathTarget{}
	var (
		last  *RecordPathSegment
		pages *PageBuilder
	)
	for _, part := range parts[1:] {
		switch {
		case t.Page != nil:
			return nil, bad("%q after the page", part)
		case strings.HasPrefix(part, "<"):
			if last == nil || last.Model.singleton || !last.ID.IsZero() {
				return nil, bad("the id %s of no model", part)
			}
			if last.ID, err = last.Model.ParseRecordID(part[1 : len(part)-1]); err != nil {
				return nil, bad("the id of %s: %v", last.Model.UniquePermName(), err)
			}
		case strings.HasPrefix(part, "/"):
			var pr *PagesRegistrator
			switch {
			case last == nil:
				pr = b.pagesRegistrator
			case !last.Model.singleton && last.ID.IsZero():
				pr = last.Model.listing.pagesRegistrator
			default:
				pr = last.Model.detailing.pagesRegistrator
			}
			if pr != nil {
				t.Page = pr.GetHttpPage(part)
				pages = pr.Get(part)
			}
			if t.Page == nil {
				return nil, bad("no page %q", part)
			}
		default:
			var mb *ModelBuilder
			switch {
			case last == nil:
				for _, m := range b.models {
					if m.parent == nil && m.UniquePermName() == part {
						mb = m
					}
				}
			case last.Model.singleton || !last.ID.IsZero():
				for _, m := range last.Model.children {
					if m.UniquePermName() == part {
						mb = m
					}
				}
			}
			if mb == nil {
				return nil, bad("no model %q there", part)
			}
			t.Models = append(t.Models, RecordPathSegment{Model: mb})
			last = &t.Models[len(t.Models)-1]
		}
	}
	if last == nil && t.Page == nil {
		return nil, bad("names nothing")
	}
	t.Listing = last != nil && !last.Model.singleton && last.ID.IsZero()
	if action != "" {
		var list []*ActionBuilder
		switch {
		case t.Page != nil:
			if pages != nil {
				list = pages.actions
			}
		case t.Listing:
			list = last.Model.listing.actions
		default:
			list = last.Model.detailing.actions
		}
		for _, a := range list {
			if a.PermName() == action {
				t.Action = a
			}
		}
		if t.Action == nil {
			return nil, bad("no action %q there", action)
		}
	}
	return
}

// Href is the admin's address of the target: the detail of the record, the
// listing, the singleton's page, the page; an action's is its page's, with
// the suffix "/action:<name>".
func (t *RecordPathTarget) Href(b *Builder) (href string) {
	if len(t.Models) == 0 {
		href = t.Page.FullPath()
	} else {
		last := t.Models[len(t.Models)-1]
		parents := make(IDSlice, len(t.Models)-1)
		for n, s := range t.Models[:len(t.Models)-1] {
			parents[n] = s.ID
		}
		info := last.Model.Info()
		if last.Model.singleton || t.Listing {
			href = info.ListingHref(parents...)
		} else {
			href = info.DetailingHref(last.ID, parents...)
		}
		if t.Page != nil {
			href += t.Page.Path()
		}
	}
	if t.Action != nil {
		href += ActionSuffix + url.PathEscape(t.Action.name)
	}
	return
}

// RecordPathHref is the admin's address of what the path p names
// (RecordPathTarget.Href).
func (b *Builder) RecordPathHref(p string) (string, error) {
	t, err := b.ParseRecordPath(p)
	if err != nil {
		return "", err
	}
	return t.Href(b), nil
}
