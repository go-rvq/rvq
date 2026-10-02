package userdocs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strconv"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/presets/fields/schemaform"
	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/vue"
	"github.com/go-rvq/rvq/x/i18n"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"gorm.io/gorm"
)

// Schema is the schema of the form of a package's documents: the path of
// each, shown, and its content.
const Schema = `[]{get path str; content text}`

// formSuffix is the key of the form of the documents beside the Value
// column's: the column is a string, the form a list.
const formSuffix = "__userdocs"

// Builder is the documentation of an admin: its sources, and the models that
// keep and show it.
type Builder struct {
	p          *presets.Builder
	db         *gorm.DB
	sources    []*Source
	assetsPath string
	// mb is the UserDoc's, pkg the UserDocPackage's.
	mb, pkg *presets.ModelBuilder
	// locale is the language a request's documentation is in — the
	// request's, chosen in the menu —; locales the languages there are.
	locale  func(ctx *web.EventContext) string
	locales func() []string
	// customs are the documents of no part of the admin (Custom)
	customs []Custom
}

// New is the documentation of p, kept in db. The words of its page are
// registered on p's i18n.
func New(p *presets.Builder, db *gorm.DB) *Builder {
	ConfigureMessages(p.I18n())
	b := &Builder{p: p, db: db, assetsPath: strings.TrimSuffix(p.GetURIPrefix(), "/") + "/user-docs/_assets"}
	b.locales = func() []string {
		var out []string
		for _, t := range p.I18n().GetSupportLanguages() {
			out = append(out, t.String())
		}
		return out
	}
	b.locale = func(ctx *web.EventContext) string {
		langs := b.locales()
		if d := i18n.DynaFromContext(ctx.Context()); d != nil && slices.Contains(langs, d.GetLanguage()) {
			return d.GetLanguage()
		}
		if len(langs) > 0 {
			return langs[0]
		}
		return "en"
	}
	// its own documents: the forms of every model
	b.Register(coreSource())
	return b
}

// Register adds the documentation of packages; the words of each are
// registered on the i18n — a module of their own, in the languages they say.
func (b *Builder) Register(sources ...*Source) *Builder {
	for _, src := range sources {
		b.sources = append(b.sources, src)
		for tag, msgs := range src.Messages {
			b.p.I18n().RegisterForModule(tag, src.MessagesKey, msgs)
		}
	}
	return b
}

// Sources are the documentation of the packages registered.
func (b *Builder) Sources() []*Source { return b.sources }

// Locales sets the languages there are (default: the i18n's).
func (b *Builder) Locales(f func() []string) *Builder {
	b.locales = f
	return b
}

// source is the source of the package pkg, or nil.
func (b *Builder) source(pkg string) *Source {
	for _, s := range b.sources {
		if s.Package == pkg {
			return s
		}
	}
	return nil
}

// AssetsPath is where the pictures of the documents are served
// (AssetsHandler).
func (b *Builder) AssetsPath() string { return b.assetsPath }

// AssetsHandler serves the pictures of the documents:
// AssetsPath/<package, escaped>/<path under the language's directory>.
func (b *Builder) AssetsHandler() http.Handler {
	return http.StripPrefix(b.assetsPath+"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.EscapedPath()
		esc, rest, _ := strings.Cut(p, "/")
		pkg, err := url.PathUnescape(esc)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		src := b.source(pkg)
		if src == nil || path.Ext(rest) == ".md" {
			http.NotFound(w, r)
			return
		}
		rest, _ = url.PathUnescape(rest)
		rest = path.Clean(rest)
		// the picture of the language asked for, else the source's
		lang, err := src.Asset(r.URL.Query().Get("locale"), rest)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		sub, err := fs.Sub(src.FS, lang)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		http.ServeFileFS(w, r, sub, rest)
	}))
}

// DocHref is the URL of the document of a node: in the language of the
// request that opens it (the one chosen in the menu).
func (b *Builder) DocHref(_ *web.EventContext, node string) string {
	return b.mb.Info().ListingHref() + "?" + url.Values{"doc": {node}}.Encode()
}

// packagesLocale is the language of the documents of the packages listed:
// the one asked for (?locale=, the link of a document to its package's), else
// the request's.
func (b *Builder) packagesLocale(ctx *web.EventContext) string {
	if l := ctx.R.FormValue("locale"); slices.Contains(b.locales(), l) {
		return l
	}
	return b.locale(ctx)
}

// Configure makes mb — of UserDoc, a singleton — the documentation page: the
// document of the node chosen, the tree of the nodes under its item of the
// menu (presets' MenuChildren); and pkg — of UserDocPackage — the documents of
// each package, edited.
func (b *Builder) Configure(mb, pkg *presets.ModelBuilder) {
	b.mb, b.pkg = mb, pkg

	fetch := func(obj interface{}, _ presets.ID, ctx *web.EventContext) error {
		d := obj.(*UserDoc)
		d.LocaleCode = b.locale(ctx)
		return nil
	}
	// the page shows; the documents are edited in the packages (pkg)
	mb.SetReadonly(true).SetDeletingDisabled(true)
	mb.Detailing().FetchFunc(fetch)
	mb.Editing().FetchFunc(fetch)
	mb.Detailing("Documentation").Field("Documentation").
		SetI18nLabel(func(context.Context) string { return "" })
	mb.Detailing().Field("Documentation").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		return b.page(ctx)
	})
	mb.MenuChildren(b.menuChildren)

	b.configurePackages(pkg)
}

// page is the document of the node chosen — the first when none is.
func (b *Builder) page(ctx *web.EventContext) h.HTMLComponent {
	tree := b.Tree(ctx)
	return h.Div(
		// with the component — a page reached without a reload has no head
		// of its own —; a <style> Vue would drop, a <component is="style">
		// it keeps
		h.Tag("component").Attr(":is", "'style'").Children(h.RawHTML(docStyle)),
		b.document(ctx, tree, b.shown(ctx, tree)),
	).Class("user-docs").Attr("@click", b.openScript())
}

// shown is the node the request shows: its doc, else the first.
func (b *Builder) shown(ctx *web.EventContext, tree []*Node) string {
	node := ctx.R.FormValue("doc")
	if node == "" && len(tree) > 0 {
		node = tree[0].ID
	}
	return node
}

// menuChildren are the nodes of the tree as the children of the item of the
// documentation in the menu, each its document's page; the one shown active.
func (b *Builder) menuChildren(ctx *web.EventContext) ([]*presets.MenuNode, string) {
	tree := b.Tree(ctx)
	var nodes func(tree []*Node) []*presets.MenuNode
	nodes = func(tree []*Node) (r []*presets.MenuNode) {
		for _, n := range tree {
			href := b.DocHref(ctx, n.ID)
			r = append(r, &presets.MenuNode{Title: n.Title, Value: href,
				Props: map[string]any{"href": href, "prependIcon": n.Icon}, Children: nodes(n.Children)})
		}
		return
	}
	return nodes(tree), b.DocHref(ctx, b.shown(ctx, tree))
}

// openScript is the click of the document: a link to another document opens
// it as the menu opens its items — its page without a reload, the menu
// following —; another link goes as links go. A link keeps its href, so it
// still opens in a new tab.
func (b *Builder) openScript() string {
	base := strconv.Quote(b.mb.Info().ListingHref())
	// a template of Vue sees none of the globals but a few (no window, URL):
	// the window is the element's
	return `(e) => {
	const a = e.target && e.target.closest && e.target.closest("a[href]");
	if (!a || e.ctrlKey || e.metaKey || e.shiftKey || e.button) return;
	const w = a.ownerDocument.defaultView;
	const u = new w.URL(a.getAttribute("href"), w.location.href);
	if (u.pathname !== ` + base + ` || !u.searchParams.get("doc")) return;
	e.preventDefault();
	plaid().vars(vars).pushStateURL(u.pathname + u.search).go();
}`
}

// docStyle is the look of a document: its pictures within the column,
// framed; its tables, code and headings spaced.
const docStyle = `
.user-doc img { max-width: 100%; height: auto; border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity)); border-radius: 6px; margin: 8px 0; }
.user-doc h1 { font-size: 1.6rem; margin: 0 0 8px; }
.user-doc h2 { font-size: 1.25rem; margin: 20px 0 8px; }
.user-doc p, .user-doc ul, .user-doc ol { margin-bottom: 10px; }
.user-doc ul, .user-doc ol { padding-left: 24px; }
.user-doc code { background: rgba(var(--v-theme-on-surface), 0.06); padding: 1px 4px; border-radius: 4px; }
.user-doc table { border-collapse: collapse; margin: 10px 0; }
.user-doc th, .user-doc td { border: 1px solid rgba(var(--v-border-color), var(--v-border-opacity)); padding: 4px 8px; }
`

// document is the document of node in the request's language: its file, as
// the package keeps it in that language; a page made of the node's parts
// where nothing was written.
func (b *Builder) document(ctx *web.EventContext, tree []*Node, node string) h.HTMLComponent {
	msgs := GetMessages(ctx.Context())
	locale := b.locale(ctx)
	r := &renderer{b: b, ctx: ctx, tree: tree, node: node}
	file, src, content, err := b.find(locale, node)
	if err != nil {
		return v.VAlert(h.Text(fmt.Sprintf(msgs.RenderError, err))).Type("error")
	}
	if src != nil {
		out, err := r.Render(src, file, content)
		if err != nil {
			return v.VAlert(h.Text(fmt.Sprintf(msgs.RenderError, err))).Type("error")
		}
		var edit h.HTMLComponent
		if b.pkg != nil && !b.pkg.Permissioner().ReqLister(ctx.R).Denied() {
			edit = h.A(h.Text(msgs.EditDocuments)).
				Href(b.pkg.Info().ListingHref() + "?" + url.Values{"locale": {locale}}.Encode()).
				Class("text-caption")
		}
		return h.Div(
			h.Div(h.RawHTML(out)).Class("user-doc").Attr("lang", locale),
			edit,
		)
	}
	// nothing written: the node's title, and its parts
	var title string
	var parts []h.HTMLComponent
	if n := findNode(tree, node); n != nil {
		title = n.Title
		for _, c := range n.Children {
			parts = append(parts, h.Li(h.A(h.Text(c.Title)).Href(b.DocHref(ctx, c.ID))))
		}
	}
	// in the frame of a document: its look (docStyle) the same
	out := h.Div(
		h.H1(title),
		h.P(h.Text(msgs.NothingWritten)).Class("text-medium-emphasis"),
	).Class("user-doc")
	if len(parts) > 0 {
		out.AppendChildren(h.H2(msgs.SeeAlso), h.Ul(parts...))
	}
	return out
}

// find is the document of node in locale: its file, the source that has it,
// and its content as kept in locale — nil when none has it. An action that
// its model's documents do not explain is the one of whatever package
// explains it in general: actions/NAME.md (restore, the trash's, is the
// trash's); a nested model, children/WORD/README.md, WORD the last word of
// its key (revisions).
func (b *Builder) find(locale, node string) (file string, src *Source, content string, err error) {
	files := DocFiles(node)
	for _, f := range files {
		for _, s := range b.sources {
			c, ok, err := b.content(locale, s, f)
			if err != nil {
				return "", nil, "", err
			}
			if ok {
				return f, s, c, nil
			}
		}
	}
	return "", nil, "", nil
}

// DocFiles are the files the document of node is looked for in, in order:
// its own (DocFile); for an action, actions/NAME.md — the one of whatever
// package explains it in general —; for a form, forms/NAME.md (the
// documentation's own, userdocs/user_docs); for a nested model,
// children/WORD/README.md, WORD the last word of its key (posts_revisions is
// children/revisions).
func DocFiles(node string) []string {
	files := []string{DocFile(node)}
	if i := strings.LastIndex(node, "/actions/"); i >= 0 && !strings.Contains(node[i+len("/actions/"):], "/") {
		files = append(files, "actions/"+node[i+len("/actions/"):]+".md")
	} else if i := strings.LastIndex(node, "/forms/"); i >= 0 && !strings.Contains(node[i+len("/forms/"):], "/") {
		files = append(files, "forms/"+node[i+len("/forms/"):]+".md")
	} else if i := strings.LastIndex(node, "/children/"); i >= 0 {
		key := node[i+len("/children/"):]
		if j := strings.LastIndex(key, "_"); j >= 0 {
			key = key[j+1:]
		}
		files = append(files, "children/"+key+"/README.md")
	}
	return files
}

// content is the file of src as kept in locale, and whether there is one.
func (b *Builder) content(locale string, src *Source, file string) (string, bool, error) {
	var row UserDocPackage
	err := b.db.Session(&gorm.Session{}).Select("value").
		Take(&row, "locale_code = ? AND id = ?", locale, src.Package).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// not synced yet: the source's own — the language's, where written
		if l := src.LangOf(locale); l != "" {
			written, err := src.FilesOf(l)
			if err != nil {
				return "", false, err
			}
			if c, ok := filesByPath(written)[file]; ok {
				return c, true, nil
			}
		}
		files, err := src.Files()
		if err != nil {
			return "", false, err
		}
		c, ok := filesByPath(files)[file]
		return c, ok, nil
	}
	if err != nil {
		return "", false, err
	}
	files, err := ParseFiles(row.Value)
	if err != nil {
		return "", false, err
	}
	c, ok := filesByPath(files)[file]
	return c, ok, nil
}

// configurePackages makes pkg the documents of the packages in a language:
// the listing of those of the request's language, and the form of each one's
// documents — their content edited, their path shown —, with a reset.
func (b *Builder) configurePackages(pkg *presets.ModelBuilder) {
	pkg.SetCreatingDisabled(true).SetDeletingDisabled(true)
	pkg.UpdateDataOperator(func(do presets.DataOperator) presets.DataOperator {
		return do.(*gorm2op.DataOperatorBuilder).
			WrapPrepare(func(old gorm2op.Preparer) gorm2op.Preparer {
				return func(db *gorm.DB, mode gorm2op.Mode, obj interface{}, id model.ID, params *presets.SearchParams, ctx *web.EventContext) *gorm.DB {
					if mode.Is(gorm2op.Search) {
						params.Where("locale_code = ?", b.packagesLocale(ctx))
					}
					return old(db, mode, obj, id, params, ctx)
				}
			})
	})
	pkg.Listing("ID", "SourceLocale", "Unused").OrderBy("id")
	pkg.Detailing("ID", "SourceLocale", "Unused", "Value")
	ed := pkg.Editing("Value")

	schema, err := schemaform.Parse(Schema)
	if err != nil {
		panic(err)
	}
	formBuilder := func(ctx *web.EventContext, row *UserDocPackage) *schemaform.Builder {
		fb := schemaform.New()
		if src := b.source(row.ID); src != nil {
			words := packageMessages(ctx.Context(), src)
			fb.FieldInfo(func(_ *web.EventContext, p string) schemaform.FieldInfo {
				switch p {
				case "path":
					return schemaform.FieldInfo{Label: words.Path, Hint: words.PathHint}
				case "content":
					return schemaform.FieldInfo{Label: words.Content, Hint: words.ContentHint}
				}
				return schemaform.FieldInfo{}
			})
		}
		return fb
	}
	value := func(fb *schemaform.Builder, stored string) any {
		if data, err := fb.DecodeValue(stored); err == nil && data != nil {
			if _, isList := data.([]any); isList {
				return schema.Filled(data)
			}
		}
		return schema.Filled(nil)
	}

	ed.Field("Value").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		row := field.Obj.(*UserDocPackage)
		fb := formBuilder(ctx, row)
		fc := *field
		fc.FormKey = field.FormKey + formSuffix
		reset := "form[" + strconv.Quote(fc.FormKey) + "] = " + string(h.JSONString(value(fb, row.InitialValue)))
		return h.Div(
			v.VBtn(GetMessages(ctx.Context()).ResetDocuments).PrependIcon("mdi-restore").Size(v.SizeSmall).
				Variant(v.VariantTonal).Class("mb-2").Attr("@click", reset),
			vue.UserComponent(fb.ComponentFunc(schema)(&fc, ctx)).
				Assign("form", fc.FormKey, value(fb, row.Value)),
		)
	})
	ed.Field("Value").SetterFunc(func(obj interface{}, field *presets.FieldContext, ctx *web.EventContext) error {
		row := obj.(*UserDocPackage)
		key := field.FormKey + formSuffix
		if !hasKey(ctx.R.Form, key) {
			return nil
		}
		fb := formBuilder(ctx, row)
		posted, err := fb.DecodeForm(ctx, schema, ctx.R.Form, key)
		if err != nil {
			return err
		}
		var stored []any
		if old, err := ParseFiles(storedValue(b.db, row)); err == nil {
			for _, f := range old {
				stored = append(stored, map[string]any{"path": f.Path, "content": f.Content})
			}
		}
		out, err := fb.EncodeValue(schema.KeepReadOnly(posted, stored))
		if err != nil {
			return err
		}
		row.Value = out
		forgetTitles()
		return nil
	})
	d := pkg.Detailing()
	d.Field("Value").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
		files, _ := ParseFiles(field.Obj.(*UserDocPackage).Value)
		items := make([]h.HTMLComponent, len(files))
		for i, f := range files {
			items[i] = h.Li(h.Code(f.Path))
		}
		return h.Div(h.Div(h.Text(field.Label)).Class("text-caption text-medium-emphasis"), h.Ul(items...))
	})
}

// storedValue is the value row has in the database: what a post may not
// change of it (the paths).
func storedValue(db *gorm.DB, row *UserDocPackage) string {
	var stored UserDocPackage
	if err := db.Session(&gorm.Session{}).Select("value").
		Take(&stored, "locale_code = ? AND id = ?", row.LocaleCode, row.ID).Error; err != nil {
		return ""
	}
	return stored.Value
}

func hasKey(values url.Values, key string) bool {
	for k := range values {
		if k == key || strings.HasPrefix(k, key+"[") || strings.HasPrefix(k, key+".") {
			return true
		}
	}
	return false
}
