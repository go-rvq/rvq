package gitedit

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/perm"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// The actions of the default files (Builder.Defaults), each its permission
// ("!init_defaults", "!reset_defaults"): out of what an editor of the files
// is given — they replace files of the site as a whole.
const (
	// ActionInitDefaults is giving the empty repository its first files,
	// the default ones, in a first commit: the site's at once.
	ActionInitDefaults = "InitDefaults"
	// ActionResetDefaults is bringing the default files into the draft, over
	// the ones there, in a commit: published as any other.
	ActionResetDefaults = "ResetDefaults"
)

// ResetDefaultsForm is the form of ResetDefaults: whether the files out of
// the defaults, in the folders of the defaults (Builder.DefaultsDirs), go.
type ResetDefaultsForm struct {
	RemoveExtras bool
}

// InitFiles gives the empty repository (ErrNotEmpty) its first commit: the
// files of files, in a directory of their own, committed by commit (which
// may refuse them: a validation) — the site's files then those. One at a
// time with the publishings.
func (r *Repo) InitFiles(ctx context.Context, files fs.FS, commit func(d *Draft) error) error {
	if empty, err := r.Empty(ctx); err != nil {
		return err
	} else if !empty {
		return ErrNotEmpty
	}
	if err := os.MkdirAll(r.DraftsDir, 0o755); err != nil {
		return err
	}
	// a dot directory: never a draft's key
	dir, err := os.MkdirTemp(r.DraftsDir, ".init-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	if _, err := git(ctx, dir, "init", "--quiet", "--initial-branch", r.branch()); err != nil {
		return err
	}
	if err := writeFiles(dir, files); err != nil {
		return err
	}
	d := &Draft{Dir: dir, Key: filepath.Base(dir), repo: r}
	if err := commit(d); err != nil {
		return err
	}
	// a branch made meanwhile (a push): checkout refuses to move it
	return r.checkout(ctx, dir)
}

// Overlay writes the files of files in the draft, over the ones there; with
// dirs, the files under them that files has not go (a folder of the draft
// dirs names, "templates": the files it has are the ones of files).
func (d *Draft) Overlay(files fs.FS, dirs ...string) error {
	if err := writeFiles(d.Dir, files); err != nil {
		return err
	}
	extras, err := d.defaultsExtras(files, dirs)
	if err != nil {
		return err
	}
	for _, p := range extras {
		if err := os.Remove(filepath.Join(d.Dir, filepath.FromSlash(p))); err != nil {
			return err
		}
	}
	return nil
}

// writeFiles writes the files of files under dir, over the ones there.
func writeFiles(dir string, files fs.FS) error {
	return fs.WalkDir(files, ".", func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == "." || inGit("/"+p) {
			if e.IsDir() && p != "." {
				return fs.SkipDir
			}
			return nil
		}
		dst := filepath.Join(dir, filepath.FromSlash(p))
		if e.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}
		b, err := fs.ReadFile(files, p)
		if err != nil {
			return err
		}
		return os.WriteFile(dst, b, 0o644)
	})
}

// installDefaults names the permissions of the actions of the default files.
func (b *Builder) installDefaults() {
	if b.Defaults == nil {
		return
	}
	for action, words := range map[string]func(*Messages) (string, string){
		ActionInitDefaults:  func(m *Messages) (string, string) { return m.InitDefaultsAction, m.InitDefaultsAction_Desc },
		ActionResetDefaults: func(m *Messages) (string, string) { return m.ResetDefaultsAction, m.ResetDefaultsAction_Desc },
	} {
		b.page.Page().PermActionInfo(presets.ActionPerm(action),
			func(ctx context.Context) string { t, _ := words(GetMessages(ctx)); return t },
			func(ctx context.Context) string { _, d := words(GetMessages(ctx)); return d })
	}
}

// setupDefaultsActions are the actions of the default files of the page pg.
func (b *Builder) setupDefaultsActions(p *presets.Builder, pg *presets.PageBuilder,
	action func(name, asks string) *presets.ActionBuilder, label func(f func(m *Messages) string) func(ctx context.Context) string) {
	if b.Defaults == nil {
		return
	}
	// the first files of the empty repository: validated, committed by the
	// user, the site's at once
	action(ActionInitDefaults, ActionInitDefaults).Icon("mdi-folder-download-outline").
		SetI18nLabel(label(func(m *Messages) string { return m.InitDefaultsAction })).
		UpdateFunc(func(_ string, ctx *web.EventContext) error {
			m := GetMessages(ctx.Context())
			_, author := b.Identity(ctx.R)
			var hash string
			err := b.Repo.InitFiles(ctx.Context(), b.Defaults, func(d *Draft) (err error) {
				if err := b.validate(ctx.Context(), d, m); err != nil {
					return err
				}
				hash, err = d.Commit(ctx.Context(), b.panelMessage(ctx.R, d, m.InitDefaultsCommit), author)
				return err
			})
			switch {
			case errors.Is(err, ErrNotEmpty):
				ctx.Flash = m.ErrNotEmpty
			case errors.Is(err, ErrSiteChanged):
				ctx.Flash = m.ErrSiteChanged
			case err != nil:
				return err
			default:
				ctx.Flash = fmt.Sprintf(m.InitDefaultsDone, hash)
			}
			return nil
		})

	// the default files over the draft's, in a commit of it: seen in the
	// preview, published as any other (the history kept)
	reset := action(ActionResetDefaults, ActionResetDefaults).Icon("mdi-file-restore-outline").
		SetI18nLabel(label(func(m *Messages) string { return m.ResetDefaultsAction }))
	form := presets.NewModelBuilder(p, &ResetDefaultsForm{}, presets.ModelConfig().SetModuleKey(MessagesKey))
	ed := form.Editing("RemoveExtras")
	ed.Field("RemoveExtras").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		m := GetMessages(ctx.Context())
		hint := m.ResetDefaultsFormRemoveExtrasHint
		if len(b.DefaultsDirs) > 0 {
			hint += " " + strings.Join(b.DefaultsDirs, ", ")
		}
		// the changes of the restoring, browsed before it is done: drawn
		// again as the choice changes
		panel, err := b.defaultsPanel(ctx, false)
		if err != nil {
			panel = v.VAlert(h.Text(err.Error())).Type("error").Variant(v.VariantTonal)
		}
		return h.Div(
			h.P(h.Text(m.ResetDefaultsWhat)).Class("mb-3"),
			v.VCheckbox().Label(field.Label).Hint(hint).PersistentHint(true).
				Density(v.DensityCompact).Attr(web.VField(field.FormKey, false)...).
				Attr("@update:model-value", web.Plaid().EventFunc(eventDefaultsPanel).
					Query("extras", web.Var("$event ? '1' : ''")).Go()).
				Attr("data-reset-defaults-extras", true).Class("mb-3"),
			web.Portal(panel).Name(defaultsPanelPortal),
		)
	})
	reset.SetDialogWidth("95vw")
	presets.ActionForm[*ResetDefaultsForm](reset, ed, func(c *presets.ActionFormContext[*ResetDefaultsForm]) error {
		ctx := c.Context
		return b.do(ctx, func(d *Draft, author Author, m *Messages) (string, error) {
			changes, err := d.Status(ctx.Context())
			if err != nil {
				return "", err
			}
			if len(changes) > 0 {
				return m.ResetDefaultsUncommitted, nil
			}
			var dirs []string
			if c.Form.RemoveExtras {
				dirs = b.DefaultsDirs
			}
			if err := d.Overlay(b.resetFiles(d), dirs...); err != nil {
				return "", err
			}
			if changes, err = d.Status(ctx.Context()); err != nil {
				return "", err
			}
			if len(changes) == 0 {
				return m.ResetDefaultsNothing, nil
			}
			if err := b.validate(ctx.Context(), d, m); err != nil {
				// the draft as it was: nothing half done
				_, _ = git(ctx.Context(), d.Dir, "reset", "--quiet", "--hard")
				_, _ = git(ctx.Context(), d.Dir, "clean", "--quiet", "-d", "--force")
				return "", err
			}
			hash, err := d.Commit(ctx.Context(), b.panelMessage(ctx.R, d, m.ResetDefaultsCommit), author)
			return fmt.Sprintf(m.ResetDefaultsDone, len(changes), hash), err
		})
	}).Build()
}

// emptyView is the page of the empty repository: no draft is made of it
// (ErrEmpty) — its first files come by a push (git) or are the default ones
// (InitDefaults).
func (b *Builder) emptyView(ctx *web.EventContext, m *Messages) h.HTMLComponent {
	const portal = "gitEditAction"
	var push h.HTMLComponent
	if b.hasGit(ctx) {
		url := b.GitURL(ctx.R, "")
		branch := b.Repo.branch()
		cmd := func(label, value, attr string) h.HTMLComponent {
			return v.VTextField().Label(label).ModelValue(value).Readonly(true).
				Variant(v.FieldVariantOutlined).Density(v.DensityCompact).HideDetails(true).Class("mb-3").
				Attr("append-inner-icon", "mdi-content-copy").Attr(attr, value).
				Attr("@click:append-inner", "copyToClipboard("+h.JSONString(value)+")")
		}
		push = v.VCard(
			v.VCardTitle(h.Text(m.EmptyPushTitle)),
			v.VCardSubtitle(h.Text(m.EmptyPushWhat)),
			v.VCardText(
				cmd(m.GitURLLabel, url, "data-git-url"),
				cmd(m.EmptyPushNew, "git init -b "+branch+" && git add -A && git commit -m "+shellQuote(m.InitDefaultsCommit), "data-git-init"),
				cmd(m.EmptyPushCommand, "git push "+shellQuote(url)+" HEAD:"+branch, "data-git-push"),
				h.P(h.Text(m.GitHookHint)).Class("text-caption mb-2"),
				cmd(m.GitHook, fmt.Sprintf("curl -fsSL -u %s %s -o .git/hooks/commit-msg && chmod +x .git/hooks/commit-msg",
					shellQuote(b.actor(ctx.R).Login), absURL(ctx.R, b.hookPath())), "data-git-hook"),
			),
		).Variant(v.VariantOutlined).Class("mb-4").Attr("data-empty-push", true)
	}
	var init h.HTMLComponent
	if b.Defaults != nil {
		if ver := b.page.Page().ActionVerifier(ctx.R, presets.ActionPerm(ActionInitDefaults)); ver == nil || ver.Allowed() {
			init = v.VCard(
				v.VCardTitle(h.Text(m.InitDefaultsAction)),
				v.VCardSubtitle(h.Text(m.InitDefaultsWhat)),
				v.VCardText(
					v.VBtn(m.InitDefaultsAction).PrependIcon("mdi-folder-download-outline").Color("primary").
						Variant(v.VariantTonal).Attr("data-init-defaults", true).
						Attr("@click", b.actionOnClick(ctx, ActionInitDefaults, portal)),
				),
			).Variant(v.VariantOutlined).Class("mb-4")
		}
	}
	return h.Div(
		web.Portal().Name(portal),
		v.VAlert(h.Text(m.EmptyHint)).Title(m.EmptyTitle).Type("info").Variant(v.VariantTonal).
			Class("mb-4").Attr("data-empty-repo", true),
		init,
		push,
	).Class("pa-4")
}

// defaultsPanelPortal is the portal of the changes ResetDefaults makes,
// drawn again (eventDefaultsPanel) when the choice of the extras changes.
const defaultsPanelPortal = "gitEditDefaultsPanel"

// eventDefaultsPanel draws the changes of ResetDefaults again (query
// extras: "1" the extras removed).
const eventDefaultsPanel = "gitEditDefaultsPanel"

// eventDefaultsDiff is the diff of a file ResetDefaults changes (the
// browser's Load; query extras as eventDefaultsPanel's).
const eventDefaultsDiff = "gitEditDefaultsDiff"

// defaultsDraft is the request's draft, for the events of ResetDefaults:
// its permission asked.
func (b *Builder) defaultsDraft(ctx *web.EventContext) (*Draft, error) {
	if b.Defaults == nil || !b.allowedAction(ctx.R, ActionResetDefaults) {
		return nil, perm.PermissionDenied
	}
	d, _, err := b.Draft(ctx.R)
	if errors.Is(err, ErrNoAccess) {
		return nil, perm.PermissionDenied
	}
	return d, err
}

// defaultsDirs are the folders whose extras go: none unless asked.
func (b *Builder) defaultsDirs(extras bool) []string {
	if extras {
		return b.DefaultsDirs
	}
	return nil
}

// defaultsPanel are the changes ResetDefaults makes in the request's draft
// — with the extras removed or not —, browsed as the draft's: the tree of
// the files, each opened in a tab, the draft's and the default side by side.
func (b *Builder) defaultsPanel(ctx *web.EventContext, extras bool) (h.HTMLComponent, error) {
	m := GetMessages(ctx.Context())
	d, err := b.defaultsDraft(ctx)
	if err != nil {
		return nil, err
	}
	var warn h.HTMLComponent
	if changes, err := d.Status(ctx.Context()); err != nil {
		return nil, err
	} else if len(changes) > 0 {
		warn = v.VAlert(h.Text(m.ResetDefaultsUncommitted)).Type("warning").Variant(v.VariantTonal).Class("mb-3").
			Attr("data-defaults-uncommitted", true)
	}
	changes, err := d.DefaultsChanges(b.resetFiles(d), b.defaultsDirs(extras)...)
	if err != nil {
		return nil, err
	}
	if len(changes) == 0 {
		return h.Div(warn, h.P(h.Text(m.ResetDefaultsNothing)).Class("text-medium-emphasis").
			Attr("data-defaults-nothing", true)), nil
	}
	files := make([]vx.DiffFile, 0, len(changes))
	for _, ch := range changes {
		f := vx.DiffSummary(ch.Path, ch.Status, "", CodeLanguage(ch.Path))
		f.ReadOnly = true
		// an extra the restoring removes: to look at before confirming
		f.Warning = ch.Status == "D"
		files = append(files, f)
	}
	ex := ""
	if extras {
		ex = "1"
	}
	return h.Div(
		warn,
		h.P(h.Text(fmt.Sprintf(m.ResetDefaultsCount, len(changes)))).Class("text-caption mb-2"),
		vx.VXDiffBrowser(files...).Height("60vh").
			Labels(m.DiffFiles, m.ResetDefaultsDiffOld, m.ResetDefaultsDiffNew, m.DiffBinary, m.ResetDefaultsNothing).
			Load(web.Plaid().EventFunc(eventDefaultsDiff).Query("extras", ex)).
			LoadTexts(m.DiffLoading, m.DiffLoadError).
			NavTexts(m.DiffPrev, m.DiffNext),
	).Attr("data-defaults-changes", len(changes)), nil
}

// defaultsPanelEvent is eventDefaultsPanel's handler.
func (b *Builder) defaultsPanelEvent(ctx *web.EventContext) (r web.EventResponse, err error) {
	panel, err := b.defaultsPanel(ctx, ctx.R.FormValue("extras") != "")
	if err != nil {
		return
	}
	r.UpdatePortal(defaultsPanelPortal, panel)
	return
}

// defaultsDiffEvent is eventDefaultsDiff's handler: the draft's file and the
// default one, compared.
func (b *Builder) defaultsDiffEvent(ctx *web.EventContext) (r web.EventResponse, err error) {
	d, err := b.defaultsDraft(ctx)
	if err != nil {
		return
	}
	m := GetMessages(ctx.Context())
	f, err := d.DefaultsDiff(b.resetFiles(d), ctx.R.FormValue("path"), b.defaultsDirs(ctx.R.FormValue("extras") != "")...)
	switch {
	case err != nil:
		r.RunScript = vx.DiffContentErrorScript(err.Error())
	case f == nil:
		r.RunScript = vx.DiffContentErrorScript(m.DiffGone)
	default:
		r.RunScript = vx.DiffContentScript(*f)
	}
	return r, nil
}

// DefaultsChanges are the changes Overlay would make in the draft — its
// files as they are against the defaults, the extras of dirs deleted —, by
// path: A a default file the draft has not, M one it has otherwise, D an
// extra.
func (d *Draft) DefaultsChanges(files fs.FS, dirs ...string) (changes []Change, err error) {
	err = fs.WalkDir(files, ".", func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if p == "." || inGit("/"+p) {
			if e.IsDir() && p != "." {
				return fs.SkipDir
			}
			return nil
		}
		if e.IsDir() {
			return nil
		}
		def, err := fs.ReadFile(files, p)
		if err != nil {
			return err
		}
		cur, err := os.ReadFile(filepath.Join(d.Dir, filepath.FromSlash(p)))
		switch {
		case errors.Is(err, fs.ErrNotExist):
			changes = append(changes, Change{Status: "A", Path: p})
		case err != nil:
			return err
		case !bytes.Equal(cur, def):
			changes = append(changes, Change{Status: "M", Path: p})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	extras, err := d.defaultsExtras(files, dirs)
	if err != nil {
		return nil, err
	}
	for _, p := range extras {
		changes = append(changes, Change{Status: "D", Path: p})
	}
	sort.Slice(changes, func(i, j int) bool { return changes[i].Path < changes[j].Path })
	return changes, nil
}

// defaultsExtras are the files of the draft under dirs that files has not.
func (d *Draft) defaultsExtras(files fs.FS, dirs []string) (extras []string, err error) {
	for _, dir := range dirs {
		dir = strings.Trim(filepath.ToSlash(filepath.Clean(dir)), "/")
		if err := d.check(dir); err != nil {
			return nil, err
		}
		root := filepath.Join(d.Dir, filepath.FromSlash(dir))
		err := filepath.WalkDir(root, func(p string, e fs.DirEntry, err error) error {
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					return nil
				}
				return err
			}
			if e.IsDir() {
				return nil
			}
			rel, err := filepath.Rel(d.Dir, p)
			if err != nil {
				return err
			}
			rel = filepath.ToSlash(rel)
			if _, err := fs.Stat(files, rel); errors.Is(err, fs.ErrNotExist) {
				extras = append(extras, rel)
			} else if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return extras, nil
}

// DefaultsDiff is the diff of the change of DefaultsChanges at path: the
// draft's file (old) and the default one (new); nil when path is none of
// them.
func (d *Draft) DefaultsDiff(files fs.FS, path string, dirs ...string) (*vx.DiffFile, error) {
	changes, err := d.DefaultsChanges(files, dirs...)
	if err != nil {
		return nil, err
	}
	for _, ch := range changes {
		if ch.Path != path {
			continue
		}
		var old, cur []byte
		if ch.Status != "A" {
			if old, err = os.ReadFile(filepath.Join(d.Dir, filepath.FromSlash(path))); err != nil {
				return nil, err
			}
		}
		if ch.Status != "D" {
			if cur, err = fs.ReadFile(files, path); err != nil {
				return nil, err
			}
		}
		f := vx.NewDiffFile(path, ch.Status, CodeLanguage(path), string(old), string(cur))
		f.ReadOnly = true
		return &f, nil
	}
	return nil, nil
}

// resetFiles are the defaults ResetDefaults brings into d: those of
// DefaultsKeep that d has, its own.
func (b *Builder) resetFiles(d *Draft) fs.FS {
	if len(b.DefaultsKeep) == 0 {
		return b.Defaults
	}
	return keepFS{b.Defaults, os.DirFS(d.Dir), b.DefaultsKeep}
}

// keepFS is files, but the paths of keep that own has: own's.
type keepFS struct {
	files, own fs.FS
	keep       []string
}

func (k keepFS) Open(name string) (fs.File, error) {
	if slices.Contains(k.keep, name) {
		if f, err := k.own.Open(name); err == nil {
			return f, nil
		}
	}
	return k.files.Open(name)
}

func (k keepFS) ReadDir(name string) ([]fs.DirEntry, error) { return fs.ReadDir(k.files, name) }
