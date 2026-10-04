package gitedit

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/actions"
	rvqjs "github.com/go-rvq/rvq/js"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/ory/ladon"
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
	// HelpURL, when set, is the documentation of the editor, opened by the
	// page's Help button.
	HelpURL string
	// Assets are the files of the IDE's app (its index.html at the root); the
	// gad IDE's (rvq/js GadIDE) when not set.
	Assets fs.FS
	// GitURL, when set, is the URL by git of the files of the site (draft
	// false: a push publishes, GitRepo) or of the draft of the user of r
	// (DraftGitRepo), shown on the page to whoever may (!git).
	GitURL func(r *http.Request, draft bool) string

	ide     *IDE
	page    *presets.PageBuilder
	idePage *presets.HttpPageBuilder
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

// Draft is the draft of the request.
func (b *Builder) Draft(r *http.Request) (*Draft, Author, error) {
	key, author := b.Identity(r)
	d, err := b.Repo.Draft(r.Context(), key)
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

	b.setupActions(p)

	b.ide = NewIDE(b.Repo)
	b.ide.DraftKey = func(r *http.Request) string { key, _ := b.Identity(r); return key }
	b.ide.Allowed = b.AllowedOp
	// the IDE, in a frame: its app and its API, asking the page's permissions
	b.idePage = presets.HttpPage(b.Path + "/ide/{rest...}").InMenu(false).Handler(http.HandlerFunc(b.serveIDE))
	p.PagesRegistrator().AddHttpPage(b.idePage)
	// under the page wherever the menu puts it (its groups are in its URL)
	b.page.Page().PostBuild(func(*presets.PageHandler) {
		if g := b.page.Page().GetMenuGroupBuilder(); g != nil {
			b.idePage.SetMenuGroup(g)
		}
	})

	if b.PreviewPath != "" {
		b.page.Page().PermActionInfo(presets.ActionPerm(ActionPreview),
			func(ctx context.Context) string { return GetMessages(ctx).PreviewAction },
			func(ctx context.Context) string { return GetMessages(ctx).PreviewAction_Desc })
		prefix := strings.TrimSuffix(p.GetURIPrefix(), "/") + b.PreviewPath
		if b.PreviewURL == "" {
			b.PreviewURL = prefix + "/"
		}
		p.PagesRegistrator().AddHttpPage(presets.HttpPage(b.PreviewPath + "/{rest...}").InMenu(false).
			Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { b.servePreview(w, r, prefix) })))
	}
	return nil
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
	rest := r.URL.Path[strings.Index(r.URL.Path, b.Path+"/ide/")+len(b.Path+"/ide"):]
	if strings.HasPrefix(rest, "/api/ide/") {
		r2 := *r
		u := *r.URL
		u.Path, u.RawPath = rest, ""
		r2.URL = &u
		b.ide.ServeHTTP(w, &r2)
		return
	}
	// the app's assets — its code, the same for everybody — are served as
	// the static files they are: the login lets them by with no session
	// (no user to ask a permission of); its page and its API ask the page's
	if !strings.HasPrefix(rest, "/assets/") {
		if ver := b.page.Page().ActionVerifier(r, presets.PermGet); ver != nil && ver.Denied() {
			http.Error(w, "permission denied", http.StatusForbidden)
			return
		}
	}
	if b.Assets == nil {
		http.Error(w, "the IDE is not built", http.StatusNotFound)
		return
	}
	if rest == "/" || rest == "/index.html" {
		// itself: a server that drops the trailing slash of a URL would
		// send "…/ide/" away from its assets ("./assets/…"), and the files
		// server sends "index.html" to "./"
		index, err := fs.ReadFile(b.Assets, "index.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(index)
		return
	}
	r2 := *r
	u := *r.URL
	u.Path, u.RawPath = rest, ""
	r2.URL = &u
	http.FileServerFS(b.Assets).ServeHTTP(w, &r2)
}

// do runs f on the draft of the request: the flash says how it went.
func (b *Builder) do(ctx *web.EventContext, f func(d *Draft, author Author, m *Messages) (string, error)) error {
	m := GetMessages(ctx.Context())
	d, author, err := b.Draft(ctx.R)
	if err != nil {
		return err
	}
	msg, err := f(d, author, m)
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

func (b *Builder) setupActions(p *presets.Builder) {
	label := func(f func(m *Messages) string) func(ctx context.Context) string {
		return func(ctx context.Context) string { return f(GetMessages(ctx)) }
	}

	commit := b.page.Action(ActionCommit).Icon("mdi-source-commit").
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
			hash, err := d.Commit(c.Context.Context(), c.Form.Message, author)
			return fmt.Sprintf(m.Committed, hash), err
		})
	}).Build()

	b.page.Action(ActionUpdate).Icon("mdi-source-pull").
		SetI18nLabel(label(func(m *Messages) string { return m.UpdateAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			return b.do(ctx, func(d *Draft, _ Author, m *Messages) (string, error) {
				return m.Updated, d.Update(ctx.Context())
			})
		})

	b.page.Action(ActionPublish).Icon("mdi-publish").
		SetI18nLabel(label(func(m *Messages) string { return m.PublishAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			return b.do(ctx, func(d *Draft, _ Author, m *Messages) (string, error) {
				if err := b.validate(ctx.Context(), d, m); err != nil {
					return "", err
				}
				return m.Published, d.Publish(ctx.Context())
			})
		})

	b.page.Action(ActionDiscard).Icon("mdi-undo").
		SetI18nLabel(label(func(m *Messages) string { return m.DiscardAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			path := ctx.R.FormValue("path")
			return b.do(ctx, func(d *Draft, _ Author, m *Messages) (string, error) {
				return fmt.Sprintf(m.Discarded, path), d.Discard(ctx.Context(), path)
			})
		})

	b.page.Action(ActionReset).Icon("mdi-restore").
		SetI18nLabel(label(func(m *Messages) string { return m.ResetAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			return b.do(ctx, func(d *Draft, _ Author, m *Messages) (string, error) {
				key, _ := b.Identity(ctx.R)
				b.ide.Forget(key)
				return m.Reset, d.Reset()
			})
		})
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
	d, _, err := b.Draft(ctx.R)
	if err != nil {
		return
	}
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

	ideURL := strings.TrimSuffix(b.idePage.FullPath(), "{rest...}") + "index.html"
	// the IDE in a tab of its own, the whole window
	openIDE := v.VBtn(m.OpenIDE).PrependIcon("mdi-open-in-new").Variant(v.VariantTonal).Color("primary").
		Size(v.SizeSmall).Class("me-2 mb-2").Href(ideURL).Attr("target", "_blank")
	var help h.HTMLComponent
	if b.HelpURL != "" {
		help = v.VBtn(m.Help).PrependIcon("mdi-help-circle-outline").Variant(v.VariantText).Size(v.SizeSmall).
			Class("me-2 mb-2").Href(b.HelpURL).Attr("target", "_blank")
	}
	var preview h.HTMLComponent
	if b.PreviewURL != "" {
		preview = v.VBtn(m.Preview).PrependIcon("mdi-eye-outline").Variant(v.VariantTonal).Size(v.SizeSmall).
			Class("me-2 mb-2").Href(b.PreviewURL).Attr("target", "_blank")
	}

	r.Body = h.Div(
		web.Portal().Name(portal),
		v.VRow(
			v.VCol(
				v.VCard(
					v.VCardTitle(h.Text(m.Draft)),
					v.VCardSubtitle(h.Text(m.DraftHint)),
					v.VCardText(
						h.P(h.Text(state)).Class("mb-3"),
						openIDE,
						btn(ActionCommit, m.CommitAction, "mdi-source-commit", len(changes) > 0),
						btn(ActionUpdate, m.UpdateAction, "mdi-source-pull", sync.Behind > 0),
						btn(ActionPublish, m.PublishAction, "mdi-publish", sync.Ahead > 0 && len(changes) == 0),
						preview,
						btn(ActionReset, m.ResetAction, "mdi-restore", true),
						help,
						h.H4(m.Changes).Class("mt-4 mb-2"),
						changesComp,
						h.H4(m.History).Class("mt-4 mb-2"),
						v.VList(logItems...).Density(v.DensityCompact),
						b.gitURLs(ctx, m, allowed),
					),
				).Variant(v.VariantOutlined),
			).Cols(12).Lg(4),
			v.VCol(
				h.Iframe().Src(ideURL).Attr("title", m.Files).
					Style("width: 100%; height: calc(100vh - 160px); border: 0"),
			).Cols(12).Lg(8),
		),
	)
	return
}

// gitURLs are the URLs by git (GitURL) of the draft and — to whoever may
// publish — of the site, each with a button that copies it; nothing to whoever
// may not reach them by git (!git).
func (b *Builder) gitURLs(ctx *web.EventContext, m *Messages, allowed func(string) bool) h.HTMLComponent {
	if b.GitURL == nil || !allowed(ActionGit) {
		return nil
	}
	field := func(label, url string) h.HTMLComponent {
		return v.VTextField().Label(label).ModelValue(url).Readonly(true).
			Variant(v.FieldVariantOutlined).Density(v.DensityCompact).HideDetails(true).Class("mb-2").
			Attr("append-inner-icon", "mdi-content-copy").Attr("data-git-url", url).
			Attr("@click:append-inner", fmt.Sprintf("navigator.clipboard.writeText(%q)", url))
	}
	items := []h.HTMLComponent{
		h.H4(m.ByGit).Class("mt-4 mb-1"),
		h.P(h.Text(m.ByGitHint)).Class("text-caption mb-2"),
		field(m.GitDraftURL, b.GitURL(ctx.R, true)),
	}
	if allowed(ActionPublish) {
		items = append(items, field(m.GitSiteURL, b.GitURL(ctx.R, false)))
	}
	return h.Div(items...)
}
