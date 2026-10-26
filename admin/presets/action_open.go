package presets

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
)

// ActionSuffix is the suffix of the address of an action: the address of
// what it is of — a record's detail, a listing, a page — and
// "/action:<name>" (RecordPathTarget.Href): "/admin/users/1/action:reset".
// The page of the action: its form in the body of a page, as in its dialog,
// its button doing it; done, back to the page it is of.
const ActionSuffix = "/action:"

const (
	// ParamActionPage says the event is of the page of an action
	ParamActionPage = "presets_action_page"
	// ActionPagePortalName is the portal of the form of the page of an
	// action: drawn again with its errors
	ActionPagePortalName = "presets_ActionPage"
)

type actionPageKey struct{}

// withActionSuffix is r for the page under an action's address: its path
// without "/action:<name>", the name kept in its context (ActionPageName).
// Other requests are r.
func withActionSuffix(r *http.Request) *http.Request {
	i := strings.LastIndex(r.URL.Path, ActionSuffix)
	if i < 0 || r.Method != http.MethodGet || r.URL.Query().Has(web.EventFuncIDName) {
		return r
	}
	name, err := url.PathUnescape(r.URL.Path[i+len(ActionSuffix):])
	if err != nil || name == "" || strings.Contains(name, "/") {
		return r
	}
	r = r.WithContext(context.WithValue(r.Context(), actionPageKey{}, name))
	u := *r.URL
	u.Path, u.RawPath = u.Path[:i], ""
	r.URL = &u
	r.RequestURI = u.RequestURI()
	return r
}

// ActionPageName is the action of the page asked — the name of its
// address's "/action:<name>" —; "" for none.
func ActionPageName(r *http.Request) string {
	name, _ := r.Context().Value(actionPageKey{}).(string)
	return name
}

// IsActionPage says r is an event of the page of an action.
func IsActionPage(r *http.Request) bool {
	return r.FormValue(ParamActionPage) != ""
}

// errNoActionPage is the error of an address of an action that is not one
// of what it is of.
var errNoActionPage = errors.New("no action")

// PageComponent is the form of the action in the body of a page — what its
// dialog shows —, its button doing it. Its fields are of the form of the
// page — no scope of their own —: the button, in the bar of the page, sends
// that one.
func (b *ActionBuilder) PageComponent(mb *ModelBuilder, id string, ctx *web.EventContext) (comp h.HTMLComponent, err error) {
	var cb *ContentComponentBuilder
	if cb, comp, err = b.content(mb, id, "", ctx); err != nil {
		return
	}
	if comp == nil {
		comp = cb.BuildPage()
	}
	return
}

// actionPage makes r the page of the action its address asks, of list,
// when it asks one (ok): id is the record's ("" for none), allowed says the
// permission of the action; back is the address of what it is of. The events
// of its form go to the address of what it is of, with the action and the
// record.
func actionPage(ctx *web.EventContext, r *web.PageResponse, mb *ModelBuilder, list []*ActionBuilder, id string,
	allowed func(a *ActionBuilder) bool) (ok bool, err error) {
	name := ActionPageName(ctx.R)
	if name == "" {
		return
	}
	ok = true
	a := getAction(list, name)
	if a == nil {
		return ok, errNoActionPage
	}
	if allowed != nil && !allowed(a) {
		return ok, perm.PermissionDenied
	}
	if a.enabledFunc != nil {
		var enabled bool
		if enabled, err = a.IsEnabled(id, ctx); err != nil {
			return
		}
		if !enabled {
			return ok, ErrActionNotAllowed
		}
	}

	// what is of the page (the toolbar, the breadcrumb) is of what the
	// action is of: back to it
	back := ctx.R.URL.Path
	if title := r.PageTitle; title != "" {
		GetOrInitBreadcrumbs(ctx.R).Append(&Breadcrumb{Label: title, URI: back})
	}

	q := ctx.R.URL.Query()
	q.Set(ParamAction, name)
	q.Set(ParamActionPage, "1")
	if id != "" {
		q.Set(ParamID, id)
	}
	u := *ctx.R.URL
	u.RawQuery = q.Encode()
	ctx.R.URL = &u
	ctx.R.RequestURI = u.RequestURI()

	var comp h.HTMLComponent
	if comp, err = a.PageComponent(mb, id, ctx); err != nil {
		return
	}
	r.PageTitle = a.RequestTitle(mb, ctx.Context())
	r.Body = web.Portal(comp).Name(ActionPagePortalName)
	return
}

// IsErrorFlash says the flash says something failed: a message, or one of
// the messages, of the color "error".
func IsErrorFlash(flash any) bool {
	switch m := flash.(type) {
	case *FlashMessage:
		return m != nil && m.Color == "error"
	case FlashMessages:
		for _, m := range m {
			if m != nil && m.Color == "error" {
				return true
			}
		}
	}
	return false
}
