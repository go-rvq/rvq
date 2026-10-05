package gitedit

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"sync"
	"time"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/ory/ladon"
	"gorm.io/gorm"
)

// The actions of the page, each its permission ("!commit", ActionPerm): the
// page's actions pass by the permissions with no more said.
const (
	ActionCommit  = "Commit"
	ActionUpdate  = "Update"
	ActionPublish = "Publish"
	ActionDiscard = "Discard"
	ActionReset   = "Reset"
	// ActionCreate is creating a file or a folder in the draft.
	ActionCreate = "Create"
	// ActionEdit is changing a file of the draft that is there.
	ActionEdit = "Edit"
	// ActionRename is renaming a file or a folder in its folder.
	ActionRename = "Rename"
	// ActionMove is moving a file or a folder to another folder.
	ActionMove = "Move"
	// ActionDelete is deleting a file or a folder of the draft.
	ActionDelete = "Delete"
	// ActionImport is bringing files into the draft: uploaded from the
	// computer, or downloaded from a URL by the server. What it writes asks
	// ActionCreate or ActionEdit too.
	ActionImport = "Import"
	// ActionPreview is seeing the site as the draft makes it.
	ActionPreview = "Preview"
)

// CommitForm is the form of a commit.
type CommitForm struct {
	Message string
}

// Builder is the page of the admin that edits the files of Repo: the IDE on
// the editor's draft, its changes, history and actions.
type Builder struct {
	Repo *Repo
	// Path is the page's ("/site-files").
	Path string
	// Identity is the draft and the author of the request (its user).
	Identity func(r *http.Request) (key string, author Author)
	// Validate says whether the files of a draft work (the templates
	// compile): asked before a commit and before publishing.
	Validate func(ctx context.Context, dir string) error
	// PreviewPath is where the site of the draft is seen, under the admin
	// ("/site-preview"): Preview serves it. "" is no preview.
	PreviewPath string
	// Preview serves r — its path the site's ("/en-us/about"), of the
	// draft d — as d makes the site; the URIs of the site written under
	// prefix ("/admin/site-preview").
	Preview func(w http.ResponseWriter, r *http.Request, d *Draft, prefix string)
	// PreviewURL is where the site of the draft is seen (PreviewPath's,
	// when not set).
	PreviewURL string
	// PublicPreviewPath, when set, is where the site serves the public
	// previews ("/_preview"): the application mounts PublicPreviewHandler
	// there, out of the admin. "" is no public preview.
	PublicPreviewPath string
	// HelpURL, when set, is the documentation of the editor, opened by the
	// page's Help button.
	HelpURL string
	// Assets are the files of the IDE — gadide.js and gadide.css, the
	// component the admin's <vx-gad-ide> renders —; the gad IDE's (rvq/js
	// GadIDE) when not set.
	Assets fs.FS
	// Brand is what the header of the editor (the page that has only the
	// IDE, opened in a tab of the browser of its own) shows of the admin:
	// its name and the URL of its logo ("" none). The page's title when not
	// set.
	Brand func(ctx *web.EventContext) (title, logo string)
	// LogoutURL, when set, is where the header of the editor signs the user
	// out ("" none): asked each time — a login knows its URLs once installed.
	LogoutURL func() string
	// GitURL, when set, is the URL by git of the files of the site (draft
	// "": a push publishes, GitRepo) or of the draft of the user of key draft
	// (DraftGitRepo: the user of r's; DraftsGitRepo: another's), shown on the
	// page to whoever may (!git).
	GitURL func(r *http.Request, draft string) string

	// Users are the users, as the application knows them: their logins and
	// names, whom to share a draft with. Without, a user is their key.
	Users Users
	// DB keeps the sharings of the drafts (DraftShare); none without.
	DB *gorm.DB
	// Host is the site's host of the request, what its commits say they were
	// made on (Site-User); the request's when not set.
	Host func(r *http.Request) string

	ide           *IDE
	page          *presets.PageBuilder
	userPage      *presets.PageBuilder
	adminURL      string
	previewPrefix string
	locks         sync.Map
}

// Page is the page of the editor.
func (b *Builder) Page() *presets.PageBuilder { return b.page }

// opPerms are the permissions of the operations on the files.
var opPerms = map[IdeOp]string{
	IdeRead:   presets.PermGet,
	IdeCreate: presets.ActionPerm(ActionCreate),
	IdeEdit:   presets.ActionPerm(ActionEdit),
	IdeRename: presets.ActionPerm(ActionRename),
	IdeMove:   presets.ActionPerm(ActionMove),
	IdeDelete: presets.ActionPerm(ActionDelete),
	IdeImport: presets.ActionPerm(ActionImport),
}

// AllowedOp says whether r may do op on the file (or folder) at p, a path of
// the draft ("static/css/a.css"; "" the files as a whole) — what the IDE asks,
// and whoever else serves the draft (a WebDAV, DavAllowed).
//
// The path is a part of the resource, between "<" and ">":
// "admin:…/site-files:<static/css/a.css>:!edit"; its policies are globs, so
// "<static/*>" is all under static/, recursively. It decides with the page's
// permission ("admin:…/site-files:!edit"):
//
//   - a deny of the path denies;
//   - an allow of the page allows;
//   - a deny of the page denies — whatever the path —;
//   - no policy of the page: an allow of the path allows (a role that may
//     edit only under static/: "…:<static/*>:!edit", nothing of the page).
func (b *Builder) AllowedOp(r *http.Request, op IdeOp, p string) bool {
	perm := opPerms[op]
	ver := b.page.Page().ActionVerifier(r, perm)
	if ver == nil {
		return true
	}
	var pathErr error = ladon.ErrRequestDenied // no path: no policy of it
	if p = strings.Trim(p, "/"); p != "" {
		pathErr = b.page.Page().ActionVerifier(r, perm).On("<" + p + ">").IsAllowed()
		if errors.Is(pathErr, ladon.ErrRequestForcefullyDenied) {
			return false
		}
	}
	switch pageErr := ver.IsAllowed(); {
	case pageErr == nil:
		return true
	case errors.Is(pageErr, ladon.ErrRequestForcefullyDenied):
		return false
	}
	return pathErr == nil
}

// DavAllowed says whether r, a request of a WebDAV on the draft, may do its
// method on p (to dest, of a COPY or a MOVE: "" otherwise), paths of the draft
// ("/static/a.css"): reading asks @get; PUT !create or !edit (the file is
// there or not); MKCOL and COPY !create; DELETE !delete; MOVE !rename or
// !move — on each path —; the rest of the writes (locks, properties) !edit.
func (b *Builder) DavAllowed(r *http.Request, p, dest string, write bool) bool {
	if !write {
		return b.AllowedOp(r, IdeRead, p)
	}
	ok := b.davWriteAllowed(r, p, dest)
	if ok {
		b.noteEditor(r) // a co-author of the next commit
	}
	return ok
}

func (b *Builder) davWriteAllowed(r *http.Request, p, dest string) bool {
	ok := func(op IdeOp, paths ...string) bool {
		for _, q := range paths {
			if !b.AllowedOp(r, op, q) {
				return false
			}
		}
		return true
	}
	switch r.Method {
	case http.MethodPut:
		d, _, err := b.Draft(r)
		if err != nil {
			return false
		}
		return ok(writeNeeds(d.Dir, p)[0].Op, p)
	case "MKCOL":
		return ok(IdeCreate, p)
	case "COPY":
		return ok(IdeRead, p) && ok(IdeCreate, dest)
	case http.MethodDelete:
		return ok(IdeDelete, p)
	case "MOVE":
		return ok(renameOp(p, dest), p, dest)
	}
	return ok(IdeEdit, p)
}

// Hidden says whether p, a path in a draft ("/.git/config"), is kept from
// whoever serves it.
func Hidden(p string) bool { return inGit(p) }

// Draft is the draft the request reaches — its own, or another's of its URL
// (…/u/{user}) it may reach (mayReach: ErrNoAccess) —, and the author of
// what it does there: its user's.
func (b *Builder) Draft(r *http.Request) (*Draft, Author, error) {
	owner, _ := b.ownerKey(r)
	if !b.mayReach(r, owner) {
		return nil, Author{}, ErrNoAccess
	}
	_, author := b.Identity(r)
	d, err := b.Repo.Draft(r.Context(), owner)
	return d, author, err
}

func (b *Builder) Install(p *presets.Builder) error {
	ConfigureMessages(p.I18n())
	if b.Path == "" {
		b.Path = "/site-files"
	}
	if b.Assets == nil {
		b.Assets = rvqjs.GadIDE()
	}
	b.page = p.PagesRegistrator().New(
		presets.HttpPage(b.Path).
			MenuIcon("mdi-file-code-outline").
			TitleFunc(func(ctx context.Context) string { return GetMessages(ctx).Title })).
		Private().
		Layout(b.pageFunc)
	defer b.page.Build()
	// the actions on the files, each its permission, named in the
	// permissions as the messages say
	for action, words := range map[string]func(*Messages) (string, string){
		ActionCreate:  func(m *Messages) (string, string) { return m.CreateAction, m.CreateAction_Desc },
		ActionEdit:    func(m *Messages) (string, string) { return m.EditAction, m.EditAction_Desc },
		ActionRename:  func(m *Messages) (string, string) { return m.RenameAction, m.RenameAction_Desc },
		ActionMove:    func(m *Messages) (string, string) { return m.MoveAction, m.MoveAction_Desc },
		ActionDelete:  func(m *Messages) (string, string) { return m.DeleteAction, m.DeleteAction_Desc },
		ActionImport:  func(m *Messages) (string, string) { return m.ImportAction, m.ImportAction_Desc },
		ActionGit:     func(m *Messages) (string, string) { return m.GitAction, m.GitAction_Desc },
		ActionCommit:  func(m *Messages) (string, string) { return m.CommitAction, m.CommitAction_Desc },
		ActionUpdate:  func(m *Messages) (string, string) { return m.UpdateAction, m.UpdateAction_Desc },
		ActionPublish: func(m *Messages) (string, string) { return m.PublishAction, m.PublishAction_Desc },
		ActionDiscard: func(m *Messages) (string, string) { return m.DiscardAction, m.DiscardAction_Desc },
		ActionReset:   func(m *Messages) (string, string) { return m.ResetAction, m.ResetAction_Desc },
	} {
		b.page.Page().PermActionInfo(presets.ActionPerm(action),
			func(ctx context.Context) string { t, _ := words(GetMessages(ctx)); return t },
			func(ctx context.Context) string { _, d := words(GetMessages(ctx)); return d })
	}

	for action, words := range map[string]func(*Messages) (string, string){
		ActionShare:  func(m *Messages) (string, string) { return m.ShareAction, m.ShareAction_Desc },
		ActionDrafts: func(m *Messages) (string, string) { return m.DraftsAction, m.DraftsAction_Desc },
		ActionPublicPreview: func(m *Messages) (string, string) {
			return m.PublicPreviewAction, m.PublicPreviewAction_Desc
		},
	} {
		if action == ActionPublicPreview && b.PublicPreviewPath == "" {
			continue
		}
		b.page.Page().PermActionInfo(presets.ActionPerm(action),
			func(ctx context.Context) string { t, _ := words(GetMessages(ctx)); return t },
			func(ctx context.Context) string { _, d := words(GetMessages(ctx)); return d })
	}
	if b.DB != nil {
		if err := b.DB.AutoMigrate(&DraftShare{}, &PublicPreview{}); err != nil {
			return err
		}
	}

	// the draft of another user (shared with the request's, or any with
	// !drafts): the same page under its key
	b.userPage = b.page.SubPage("/u/{" + userParam + "}").Layout(b.pageFunc)
	defer b.userPage.Build()

	b.setupActions(p, b.page)
	b.setupActions(p, b.userPage)

	b.ide = NewIDE(b.Repo)
	b.ide.DraftKey = func(r *http.Request) string { key, _ := b.ownerKey(r); return key }
	b.ide.Allowed = b.AllowedOp
	// a write of the IDE: its user a co-author of the next commit
	b.ide.OnWrite = b.noteEditor
	b.adminURL = p.GetURIPrefix()
	if b.adminURL == "" {
		b.adminURL = "/"
	}
	// the IDE: its code (Assets) and its API, the API asking the page's
	// permissions — of the own draft and of another's
	for _, sub := range []string{"/ide/{rest...}", "/u/{" + userParam + "}/ide/{rest...}"} {
		p.PagesRegistrator().AddHttpPage(b.page.Page().Sub(sub).Handler(http.HandlerFunc(b.serveIDE)))
	}
	// the editor: the IDE alone, the whole window, under a header — the page
	// "Open the editor" opens in a tab of the browser —; seeing it is seeing
	// the page (no permission of its own)
	for _, sub := range []string{"/editor", "/u/{" + userParam + "}/editor"} {
		defer b.page.SubPage(sub).Raw(p.PlainLayout(b.editorFunc)).Build()
	}
	// the hook commit-msg of the user: their commits say who made them
	p.PagesRegistrator().AddHttpPage(b.page.Page().Sub("/commit-msg").Handler(http.HandlerFunc(b.serveHook)))

	if b.PreviewPath != "" {
		b.page.Page().PermActionInfo(presets.ActionPerm(ActionPreview),
			func(ctx context.Context) string { return GetMessages(ctx).PreviewAction },
			func(ctx context.Context) string { return GetMessages(ctx).PreviewAction_Desc })
		prefix := strings.TrimSuffix(p.GetURIPrefix(), "/") + b.PreviewPath
		b.previewPrefix = prefix
		if b.PreviewURL == "" {
			// no "/" at the end: the admin sends such a URL to the one without
			b.PreviewURL = prefix
		}
		own := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { b.servePreview(w, r, prefix) })
		// another's draft: under its key
		another := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			b.servePreview(w, r, prefix+"/u/"+r.PathValue(userParam))
		})
		for path, h := range map[string]http.Handler{
			b.PreviewPath: own, b.PreviewPath + "/{rest...}": own,
			b.PreviewPath + "/u/{" + userParam + "}": another, b.PreviewPath + "/u/{" + userParam + "}/{rest...}": another,
		} {
			p.PagesRegistrator().AddHttpPage(presets.HttpPage(path).InMenu(false).Handler(h))
		}
	}
	return nil
}

// noteEditor notes the user of r as an editor of the draft it reaches (a
// co-author of its next commit).
func (b *Builder) noteEditor(r *http.Request) {
	if d, _, err := b.Draft(r); err == nil {
		key, _ := b.Identity(r)
		_ = d.NoteEditor(key)
	}
}

// servePreview serves the site as the draft of the request makes it.
func (b *Builder) servePreview(w http.ResponseWriter, r *http.Request, prefix string) {
	if ver := b.page.Page().ActionVerifier(r, presets.ActionPerm(ActionPreview)); ver != nil && ver.Denied() {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}
	if b.Preview == nil {
		http.NotFound(w, r)
		return
	}
	d, _, err := b.Draft(r)
	if errors.Is(err, ErrNoAccess) {
		http.Error(w, "permission denied", http.StatusForbidden)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r2 := r.Clone(r.Context())
	r2.URL.Path = "/" + strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, prefix), "/")
	r2.URL.RawPath = ""
	b.Preview(w, r2, d, prefix)
}

func (b *Builder) serveIDE(w http.ResponseWriter, r *http.Request) {
	rest := "/" + r.PathValue("rest")
	if strings.HasPrefix(rest, "/api/ide/") {
		if owner, _ := b.ownerKey(r); !b.mayReach(r, owner) {
			http.Error(w, "permission denied", http.StatusForbidden)
			return
		}
		r2 := *r
		u := *r.URL
		u.Path, u.RawPath = rest, ""
		r2.URL = &u
		b.ide.ServeHTTP(w, &r2)
		return
	}
	// the IDE's code (gadide.js, gadide.css) — the same for everybody —,
	// served as the static files it is
	if b.Assets == nil {
		http.Error(w, "the IDE is not built", http.StatusNotFound)
		return
	}
	r2 := *r
	u := *r.URL
	u.Path, u.RawPath = rest, ""
	r2.URL = &u
	w.Header().Set("Cache-Control", "no-cache")
	http.FileServerFS(b.Assets).ServeHTTP(w, &r2)
}

// ideBase is where the IDE of the draft r reaches is — its files and its
// API ("/admin/site/site-files/ide/", "…/site-files/u/<key>/ide/").
func (b *Builder) ideBase(r *http.Request) string { return b.scopeURL(r) + "/ide/" }

// EditorURL is the page of the editor of the draft r reaches: the IDE alone,
// under a header.
func (b *Builder) EditorURL(r *http.Request) string { return b.scopeURL(r) + "/editor" }

// ideComponent is the IDE (<vx-gad-ide>), height tall.
func (b *Builder) ideComponent(r *http.Request, height string) h.HTMLComponent {
	base := b.ideBase(r)
	return h.Tag("vx-gad-ide").Attr("src", base+"gadide.js?"+web.ExeMTime).Attr("css", base+"gadide.css?"+web.ExeMTime).
		Attr("base", base).Attr("height", height)
}

// editorFunc is the page of the editor: a header — the admin's brand, the
// title, the user, the way back to the admin, signing out, the theme — and
// the IDE under it, the rest of the window.
func (b *Builder) editorFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	if ver := b.page.Page().ActionVerifier(ctx.R, presets.PermGet); ver != nil && ver.Denied() {
		return r, perm.PermissionDenied
	}
	// another's: shared with the user, or every draft (!drafts)
	if owner, _ := b.ownerKey(ctx.R); !b.mayReach(ctx.R, owner) {
		return r, perm.PermissionDenied
	}
	m := GetMessages(ctx.Context())
	title, logo := m.Title, ""
	if b.Brand != nil {
		title, logo = b.Brand(ctx)
	}
	_, author := b.Identity(ctx.R)
	user := author.Email
	if user == "" {
		user = author.Name
	}
	var logout h.HTMLComponent
	if b.LogoutURL != nil {
		if u := b.LogoutURL(); u != "" {
			logout = v.VBtn(m.SignOut).PrependIcon("mdi-logout").Variant(v.VariantText).Href(u).
				Attr("data-editor-logout", true)
		}
	}
	r.PageTitle = m.Title
	r.Body = h.Div(
		v.VToolbar(
			h.A(
				h.If(logo != "", h.Img(logo).Attr("alt", "").Style("height: 36px; max-width: 160px; object-fit: contain")),
				h.Span(title).Class("text-h6 ms-2"),
			).Href(b.adminURL).Class("d-flex align-center ms-4 text-decoration-none text-high-emphasis").
				Attr("data-editor-brand", true),
			v.VDivider().Vertical(true).Class("mx-4 my-3"),
			h.Span(m.Title).Class("text-subtitle-1 text-medium-emphasis text-truncate"),
			v.VSpacer(),
			v.VChip(h.Text(user)).PrependIcon("mdi-account-circle").Variant(v.VariantText).
				Attr("title", author.Name).Attr("data-editor-user", true).Class("d-none d-sm-flex"),
			v.VBtn(m.BackToAdmin).PrependIcon("mdi-view-dashboard-outline").Variant(v.VariantText).
				Href(b.adminURL).Attr("data-editor-admin", true),
			logout,
			h.Tag("vx-theme-toggle").Attr("storage-key", "rvq.gitedit.theme").
				Attr("light-title", m.ThemeLight).Attr("dark-title", m.ThemeDark),
		).Density(v.DensityCompact).Height(56).Class("border-b").Flat(true),
		b.ideComponent(ctx.R, "calc(100vh - 57px)"),
	)
	return
}

// do runs f on the draft of the request — held (lockDraft): another user on
// it waits —: the flash says how it went.
func (b *Builder) do(ctx *web.EventContext, f func(d *Draft, author Author, m *Messages) (string, error)) error {
	m := GetMessages(ctx.Context())
	d, author, err := b.Draft(ctx.R)
	if err != nil {
		return err
	}
	unlock := b.lockDraft(d.Key)
	msg, err := f(d, author, m)
	unlock()
	switch {
	case errors.Is(err, ErrUncommitted):
		msg, err = m.ErrUncommitted, nil
	case errors.Is(err, ErrNotFastForward):
		msg, err = m.ErrBehind, nil
	case errors.Is(err, ErrSiteChanged):
		msg, err = m.ErrSiteChanged, nil
	}
	if err != nil {
		return err
	}
	ctx.Flash = msg
	return nil
}

func (b *Builder) validate(ctx context.Context, d *Draft, m *Messages) error {
	if b.Validate == nil {
		return nil
	}
	if err := b.Validate(ctx, d.Dir); err != nil {
		return web.NewValidationErrors().GlobalError(fmt.Sprintf(m.ErrInvalid, err.Error()))
	}
	return nil
}

// ShareForm is the form of a sharing: with whom, until when ("" until
// revoked; a date, "2026-12-31", to its end).
type ShareForm struct {
	User      string
	ExpiresAt string
}

// ActionRevoke is revoking a sharing: its owner (!share), or whoever may
// see every draft (!drafts).
const ActionRevoke = "Revoke"

// setupActions are the actions of the page pg — the own draft's, or
// another's —, each asking the permission of the page of the editor.
func (b *Builder) setupActions(p *presets.Builder, pg *presets.PageBuilder) {
	label := func(f func(m *Messages) string) func(ctx context.Context) string {
		return func(ctx context.Context) string { return f(GetMessages(ctx)) }
	}
	action := func(name, asks string) *presets.ActionBuilder {
		return pg.Action(name).SetVerifier(func(ctx *web.EventContext) *perm.Verifier {
			return b.page.Page().ActionVerifier(ctx.R, presets.ActionPerm(asks))
		})
	}

	commit := action(ActionCommit, ActionCommit).Icon("mdi-source-commit").
		SetI18nLabel(label(func(m *Messages) string { return m.CommitAction }))
	form := presets.NewModelBuilder(p, &CommitForm{}, presets.ModelConfig().SetModuleKey(MessagesKey))
	presets.ActionForm[*CommitForm](commit, form.Editing("Message"), func(c *presets.ActionFormContext[*CommitForm]) error {
		return b.do(c.Context, func(d *Draft, author Author, m *Messages) (string, error) {
			changes, err := d.Status(c.Context.Context())
			if err != nil {
				return "", err
			}
			if len(changes) == 0 {
				return m.ErrNothing, nil
			}
			if err := b.validate(c.Context.Context(), d, m); err != nil {
				return "", err
			}
			// its co-authors, who made it and where (Site-User)
			hash, err := d.Commit(c.Context.Context(), b.panelMessage(c.Context.R, d, c.Form.Message), author)
			return fmt.Sprintf(m.Committed, hash), err
		})
	}).Build()

	action(ActionUpdate, ActionUpdate).Icon("mdi-source-pull").
		SetI18nLabel(label(func(m *Messages) string { return m.UpdateAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			return b.do(ctx, func(d *Draft, _ Author, m *Messages) (string, error) {
				return m.Updated, d.Update(ctx.Context())
			})
		})

	action(ActionPublish, ActionPublish).Icon("mdi-publish").
		SetI18nLabel(label(func(m *Messages) string { return m.PublishAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			return b.do(ctx, func(d *Draft, _ Author, m *Messages) (string, error) {
				if err := b.validate(ctx.Context(), d, m); err != nil {
					return "", err
				}
				return m.Published, d.Publish(ctx.Context())
			})
		})

	action(ActionDiscard, ActionDiscard).Icon("mdi-undo").
		SetI18nLabel(label(func(m *Messages) string { return m.DiscardAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			path := ctx.R.FormValue("path")
			return b.do(ctx, func(d *Draft, _ Author, m *Messages) (string, error) {
				return fmt.Sprintf(m.Discarded, path), d.Discard(ctx.Context(), path)
			})
		})

	// starting over: the owner's, or whoever may see every draft
	action(ActionReset, ActionReset).Icon("mdi-restore").
		SetI18nLabel(label(func(m *Messages) string { return m.ResetAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			if !b.isOwnerOrAdmin(ctx.R) {
				return perm.PermissionDenied
			}
			return b.do(ctx, func(d *Draft, _ Author, m *Messages) (string, error) {
				b.ide.Forget(d.Key)
				return m.Reset, d.Reset()
			})
		})

	// sharing the own draft: with whom, until when
	share := action(ActionShare, ActionShare).Icon("mdi-account-multiple-plus-outline").
		SetI18nLabel(label(func(m *Messages) string { return m.ShareAction }))
	shareForm := presets.NewModelBuilder(p, &ShareForm{}, presets.ModelConfig().SetModuleKey(MessagesKey))
	shareEd := shareForm.Editing("User", "ExpiresAt")
	shareEd.Field("User").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return b.userSelect(ctx, field)
	})
	shareEd.Field("ExpiresAt").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		m := GetMessages(ctx.Context())
		return v.VTextField().Type("date").Label(field.Label).Hint(m.ShareExpiresHint).PersistentHint(true).
			Variant(v.FieldVariantOutlined).Density(v.DensityCompact).
			Attr(web.VField(field.FormKey, "")...).ErrorMessages(field.Errors...)
	})
	presets.ActionForm[*ShareForm](share, shareEd, func(c *presets.ActionFormContext[*ShareForm]) error {
		r := c.Context.R
		m := GetMessages(r.Context())
		owner, own := b.ownerKey(r)
		if !own {
			return perm.PermissionDenied
		}
		if c.Form.User == "" {
			return web.NewValidationErrors().FieldError("User", m.ErrShareUser)
		}
		var expires *time.Time
		if c.Form.ExpiresAt != "" {
			day, err := time.ParseInLocation("2006-01-02", c.Form.ExpiresAt, time.Local)
			if err != nil {
				return web.NewValidationErrors().FieldError("ExpiresAt", err.Error())
			}
			end := day.AddDate(0, 0, 1) // to the end of the day
			if !end.After(time.Now()) {
				return web.NewValidationErrors().FieldError("ExpiresAt", m.ErrShareExpired)
			}
			expires = &end
		}
		if _, err := b.Share(r.Context(), owner, c.Form.User, owner, expires); errors.Is(err, ErrShareSelf) {
			return web.NewValidationErrors().FieldError("User", m.ErrShareSelf)
		} else if err != nil {
			return err
		}
		c.Context.Flash = fmt.Sprintf(m.Shared, b.user(r, c.Form.User).Label())
		return nil
	}).Build()

	if b.PublicPreviewPath != "" {
		b.setupPublicPreviewActions(p, pg, action, label)
	}

	// revoking: the owner's (!share), any (!drafts)
	pg.Action(ActionRevoke).Icon("mdi-account-cancel-outline").
		SetVerifier(func(ctx *web.EventContext) *perm.Verifier {
			return b.page.Page().ActionVerifier(ctx.R, presets.PermGet)
		}).
		SetI18nLabel(label(func(m *Messages) string { return m.RevokeAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			r := ctx.R
			s, err := b.share(r.Context(), r.FormValue("id"))
			if err != nil {
				return err
			}
			actor, _ := b.Identity(r)
			if !(s.Owner == actor && b.allowedAction(r, ActionShare)) && !b.allowedAction(r, ActionDrafts) {
				return perm.PermissionDenied
			}
			if err := b.Revoke(r.Context(), s.ID, actor); err != nil {
				return err
			}
			ctx.Flash = fmt.Sprintf(GetMessages(r.Context()).Revoked, b.user(r, s.User).Label())
			return nil
		})
}

// userSelect is the choice of whom to share the draft with: the users the
// application knows, but its owner.
func (b *Builder) userSelect(ctx *web.EventContext, field *presets.FieldContext) h.HTMLComponent {
	type item struct {
		Title string `json:"title"`
		Value string `json:"value"`
	}
	owner, _ := b.ownerKey(ctx.R)
	var items []item
	if b.Users != nil {
		users, _ := b.Users.Search(ctx.R, "", 1000)
		for _, u := range users {
			if u.Key != owner {
				items = append(items, item{u.Label(), u.Key})
			}
		}
	}
	return v.VAutocomplete().Label(field.Label).Items(items).ItemTitle("title").ItemValue("value").
		Variant(v.FieldVariantOutlined).Density(v.DensityCompact).Attr("data-share-user", true).
		Attr(web.VField(field.FormKey, "")...).ErrorMessages(field.Errors...)
}

// actionOnClick opens the action name (its dialog, its confirmation).
func (b *Builder) actionOnClick(ctx *web.EventContext, name, portal string, query ...string) string {
	g := web.GET().URL(ctx.R.URL.Path).EventFunc(actions.Action).
		Query(presets.ParamAction, name).Query(presets.ParamTargetPortal, portal)
	for i := 0; i+1 < len(query); i += 2 {
		g = g.Query(query[i], query[i+1])
	}
	return g.Go()
}

func (b *Builder) pageFunc(ctx *web.EventContext) (r web.PageResponse, err error) {
	m := GetMessages(ctx.Context())
	// another's draft: the page's permission (it has none of its own) and
	// the draft's (mayReach)
	if ver := b.page.Page().ActionVerifier(ctx.R, presets.PermGet); ver != nil && ver.Denied() {
		return r, perm.PermissionDenied
	}
	d, _, err := b.Draft(ctx.R)
	if errors.Is(err, ErrNoAccess) {
		return r, perm.PermissionDenied
	}
	if err != nil {
		return
	}
	owner, own := b.ownerKey(ctx.R)
	c := ctx.Context()
	changes, err := d.Status(c)
	if err != nil {
		return
	}
	sync, err := d.Sync(c)
	if err != nil {
		return
	}
	log, err := d.Log(c, 10)
	if err != nil {
		return
	}
	const portal = "gitEditAction"
	allowed := func(action string) bool {
		ver := b.page.Page().ActionVerifier(ctx.R, presets.ActionPerm(action))
		return ver == nil || ver.Allowed()
	}
	btn := func(action, text, icon string, enabled bool, query ...string) h.HTMLComponent {
		if !allowed(action) {
			return nil
		}
		return v.VBtn(text).PrependIcon(icon).Variant(v.VariantTonal).Size(v.SizeSmall).Class("me-2 mb-2").
			Disabled(!enabled).Attr("@click", b.actionOnClick(ctx, action, portal, query...))
	}

	state := m.UpToDate
	if sync.Ahead > 0 || sync.Behind > 0 {
		var parts []string
		if sync.Ahead > 0 {
			parts = append(parts, fmt.Sprintf(m.Ahead, sync.Ahead))
		}
		if sync.Behind > 0 {
			parts = append(parts, fmt.Sprintf(m.Behind, sync.Behind))
		}
		state = strings.Join(parts, " · ")
	}

	var changeItems []h.HTMLComponent
	for _, ch := range changes {
		diff, _ := d.Diff(c, ch.Path)
		var discard h.HTMLComponent
		if allowed(ActionDiscard) {
			discard = v.VBtn("").Icon("mdi-undo").Variant(v.VariantText).Size(v.SizeSmall).
				Attr("title", m.DiscardAction).
				Attr("@click.stop", b.actionOnClick(ctx, ActionDiscard, portal, "path", ch.Path))
		}
		changeItems = append(changeItems, v.VExpansionPanel(
			v.VExpansionPanelTitle(
				v.VChip(h.Text(ch.Status)).Size(v.SizeXSmall).Class("me-2"),
				h.Code(ch.Path), v.VSpacer(), discard,
			),
			v.VExpansionPanelText(h.Pre(diff).Class("text-caption").Style("white-space: pre-wrap")),
		))
	}
	var changesComp h.HTMLComponent = h.P(h.Text(m.NoChanges)).Class("text-medium-emphasis")
	if len(changeItems) > 0 {
		changesComp = v.VExpansionPanels(changeItems...).Variant("accordion")
	}

	var logItems []h.HTMLComponent
	for _, cm := range log {
		logItems = append(logItems, v.VListItem(
			v.VListItemTitle(h.Text(cm.Subject)),
			v.VListItemSubtitle(h.Text(cm.Short+" · "+cm.Author+" · "+cm.Date.Format("2006-01-02 15:04"))),
		).Density(v.DensityCompact))
	}

	var help h.HTMLComponent
	if b.HelpURL != "" {
		help = v.VBtn(m.Help).PrependIcon("mdi-help-circle-outline").Variant(v.VariantText).Size(v.SizeSmall).
			Class("me-2 mb-2").Href(b.HelpURL).Attr("target", "_blank")
	}
	var preview h.HTMLComponent
	if u := b.previewURL(ctx.R); u != "" {
		preview = v.VBtn(m.Preview).PrependIcon("mdi-eye-outline").Variant(v.VariantTonal).Size(v.SizeSmall).
			Class("me-2 mb-2").Href(u).Attr("target", "_blank").Attr("data-preview", true)
	}
	publicPreview := b.publicPreviewView(ctx, m, portal, owner, own)
	var reset h.HTMLComponent
	if b.isOwnerOrAdmin(ctx.R) {
		reset = btn(ActionReset, m.ResetAction, "mdi-restore", true)
	}

	// the own draft's words, or another's
	draftTab, draftHint := m.Draft, m.DraftHint
	if !own {
		draftTab, draftHint = m.DraftOther, m.DraftHintOther
	}
	draft := v.VCard(
		v.VCardSubtitle(h.Text(draftHint)).Class("pt-4"),
		v.VCardText(
			h.P(h.Text(state)).Class("mb-3"),
			btn(ActionCommit, m.CommitAction, "mdi-source-commit", len(changes) > 0),
			btn(ActionUpdate, m.UpdateAction, "mdi-source-pull", sync.Behind > 0),
			btn(ActionPublish, m.PublishAction, "mdi-publish", sync.Ahead > 0 && len(changes) == 0),
			preview,
			reset,
			help,
			publicPreview,
			h.H4(m.Changes).Class("mt-4 mb-2"),
			changesComp,
			h.H4(m.History).Class("mt-4 mb-2"),
			v.VList(logItems...).Density(v.DensityCompact),
		),
	).Variant(v.VariantOutlined)

	// the draft and the IDE in tabs; at their right, a tab that opens the IDE
	// in a tab of the browser of its own, the whole window — a link: it is
	// never the tab shown (the update of the model leaves it out)
	tabItems := []h.HTMLComponent{
		v.VTab(h.Text(draftTab)).Value("draft").PrependIcon("mdi-source-branch"),
		v.VTab(h.Text(m.Files)).Value("ide").PrependIcon("mdi-file-code-outline"),
	}
	windows := []h.HTMLComponent{
		v.VTabsWindowItem(draft).Value("draft"),
		v.VTabsWindowItem(b.ideComponent(ctx.R, "calc(100vh - 220px)")).Value("ide"),
	}
	// the sharings of this draft; the drafts the user reaches
	if b.DB != nil {
		tabItems = append(tabItems, v.VTab(h.Text(m.SharesTab)).Value("shares").PrependIcon("mdi-account-multiple-outline"))
		windows = append(windows, v.VTabsWindowItem(b.sharesView(ctx, m, portal, owner, own)).Value("shares"))
	}
	if b.DB != nil || allowed(ActionDrafts) {
		tabItems = append(tabItems, v.VTab(h.Text(m.DraftsTab)).Value("drafts").PrependIcon("mdi-folder-account-outline"))
		windows = append(windows, v.VTabsWindowItem(b.draftsView(ctx, m, allowed(ActionDrafts))).Value("drafts"))
	}
	// the repositories by git: the draft's and the site's
	if b.hasGit(ctx) {
		tabItems = append(tabItems, v.VTab(h.Text(m.GitTab)).Value("git").PrependIcon("mdi-git"))
		windows = append(windows, v.VTabsWindowItem(b.gitView(ctx, m)).Value("git"))
	}
	tabItems = append(tabItems, v.VTab(h.Text(m.OpenIDE)).Value("open").PrependIcon("mdi-open-in-new").
		Attr("href", b.EditorURL(ctx.R)).Attr("target", "_blank").Attr("data-open-ide", true))

	// the draft and the IDE in tabs, and the sharings; at their right, a tab
	// that opens the IDE in a tab of the browser of its own, the whole window
	// — a link: it is never the tab shown (the update of the model leaves it
	// out)
	tabs := web.Scope(
		v.VTabs(tabItems...).Attr(":model-value", "locals.tab").
			Attr("@update:model-value", "(t) => { if (t !== 'open') locals.tab = t }").
			Color("primary").Class("mb-4"),
		v.VTabsWindow(windows...).Attr("v-model", "locals.tab"),
	).Slot("{ locals }").LocalsInit(`{tab: "draft"}`)

	r.Body = h.Div(
		web.Portal().Name(portal),
		b.whoseDraft(ctx, m, owner, own),
		tabs,
	).Class("pa-4")
	return
}
