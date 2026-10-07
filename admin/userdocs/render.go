package userdocs

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/gad-lang/gad"
	"github.com/gad-lang/gad/gadx"
	"github.com/gad-lang/gad/parser"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/web"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	gparser "github.com/yuin/goldmark/parser"
	grenderer "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// newMarkdown is the renderer of a document, its links and images resolved by
// rewrite: Gad's markdown (gadx.NewMarkdown — every bundled extension: GFM,
// typographer, definition lists, footnotes; auto heading ids; HTML kept, the
// documents being the code's and the administrators'). Its fences are the
// admin's code viewer (<vx-code>), its copy button saying labels.
func newMarkdown(rewrite func(dest string) string, labels codeLabels) goldmark.Markdown {
	md := gadx.NewMarkdown()
	md.Parser().AddOptions(gparser.WithASTTransformers(util.Prioritized(linkRewriter(rewrite), 100)))
	md.Renderer().AddOptions(grenderer.WithNodeRenderers(util.Prioritized(fenceRenderer{labels}, 100)))
	return md
}

// codeLabels are the words of a code block's copy button.
type codeLabels struct{ Copy, Copied, Failed string }

// fenceRenderer renders a fenced block of code — ```gad, ```json, … — as the
// admin's code viewer, <vx-code>: highlighted by Prism in its language, with
// its line numbers, its braces matched, its language and a copy button. The
// code is an attribute, so nothing in it is read as markup, nor by Vue.
type fenceRenderer struct{ labels codeLabels }

func (r fenceRenderer) RegisterFuncs(reg grenderer.NodeRendererFuncRegisterer) {
	reg.Register(ast.KindFencedCodeBlock, r.render)
}

func (r fenceRenderer) render(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkSkipChildren, nil
	}
	n := node.(*ast.FencedCodeBlock)
	var code strings.Builder
	for i := 0; i < n.Lines().Len(); i++ {
		seg := n.Lines().At(i)
		code.Write(seg.Value(source))
	}
	attr := func(name, value string) {
		_, _ = w.WriteString(" " + name + `="` + html.EscapeString(value) + `"`)
	}
	_, _ = w.WriteString("<vx-code")
	attr("code", strings.TrimSuffix(code.String(), "\n"))
	if lang := string(n.Language(source)); lang != "" {
		attr("language", lang)
	}
	attr("copy-text", r.labels.Copy)
	attr("copied-text", r.labels.Copied)
	attr("copy-error-text", r.labels.Failed)
	_, _ = w.WriteString("></vx-code>\n")
	return ast.WalkSkipChildren, nil
}

// linkRewriter rewrites the destination of every link and image of a
// document.
type linkRewriter func(dest string) string

func (f linkRewriter) Transform(doc *ast.Document, _ text.Reader, _ gparser.Context) {
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch t := n.(type) {
		case *ast.Link:
			t.Destination = []byte(f(string(t.Destination)))
		case *ast.Image:
			t.Destination = []byte(f(string(t.Destination)))
		}
		return ast.WalkContinue, nil
	})
}

type compiledTemplate struct {
	bc       *gad.Bytecode
	builtins *gad.Builtins
}

// templates holds a document's compilation by its source.
var templates sync.Map

// templateGlobals are the names a document's template sees.
var templateGlobals = []string{"admin", "doc"}

// RenderTemplate runs a document's template — {% … %} code, {%= expr %} a
// value — with globals in scope (templateGlobals): the markdown it writes.
func RenderTemplate(src string, globals gad.Dict) (string, error) {
	var c *compiledTemplate
	if v, ok := templates.Load(src); ok {
		c = v.(*compiledTemplate)
	} else {
		builtins := gad.NewBuiltins()
		st := gad.NewSymbolTable(builtins.NameSet)
		if _, err := st.DefineGlobals(templateGlobals); err != nil {
			return "", err
		}
		opts := gad.CompileOptions{}
		opts.ModuleFile = "doc.gadt"
		opts.ParserOptions.Mode |= parser.ParseMixed
		opts.ScannerOptions.Mode |= parser.ScanMixed | parser.ScanConfigDisabled
		opts.ScannerOptions.MixedDelimiter = parser.DefaultMixedDelimiter
		res, err := gad.Compile(st, []byte(src), opts)
		if err != nil {
			return "", err
		}
		c = &compiledTemplate{bc: res.Bytecode, builtins: builtins}
		templates.Store(src, c)
	}
	var out bytes.Buffer
	vm := gad.NewVM(c.builtins.Build(), c.bc)
	if _, err := vm.RunOpts(&gad.RunOpts{StdOut: &out, Globals: globals}); err != nil {
		return "", err
	}
	return out.String(), nil
}

// DocFile is the file of a node of the tree: README.md of a model, a group, a
// page of the admin; NAME.md of an action (MODEL/actions/NAME), of a form
// (MODEL/forms/NAME) and of a page of a model (MODEL/pages/NAME — not
// pages/NAME, a page of the admin).
func DocFile(node string) string {
	rest := node
	if i := strings.Index(node, "/"); i >= 0 {
		rest = node[i:]
	}
	// an action or a form is the last word under actions/ or forms/ — a
	// model may be called forms itself (forms/children/…)
	parts := strings.Split(node, "/")
	n := len(parts)
	if (n >= 2 && parts[n-2] == "actions") || (n >= 3 && parts[n-2] == "forms") ||
		strings.HasSuffix(node, "/permissions") || strings.Contains(rest, "/pages/") {
		return node + ".md"
	}
	return node + "/README.md"
}

// docNode is the node of a document's file (DocFile's inverse).
func docNode(file string) string {
	if n, ok := strings.CutSuffix(file, "/README.md"); ok {
		return n
	}
	return strings.TrimSuffix(file, ".md")
}

// renderer renders the documents of a request.
type renderer struct {
	b   *Builder
	ctx *web.EventContext
	// node is the node of the document rendered
	node string
	// tree is the request's tree, made once
	tree []*Node
}

// globals are the template's: admin.model, admin.action, admin.page and
// admin.doc — each a record of its label, its href and a markdown link; a
// page's also its menu, where it is in the main menu (PageMenuPath), its
// resource, of its permissions, by its groups, and its uniqueResource, by its
// unique name, out of the groups, which decides first —;
// of the model of the document (the one of its node): admin.fields(form) —
// the table of the fields of a form (FormNew, FormEdit, FormDetail) —,
// admin.menu() — the menu of its detail —, admin.permissions() — the
// permissions of the part of a permissions node, from the menu as it is —,
// admin.permissionsTree() — the tree of all the permissions, a VTreeview —,
// the application's functions (Builder.Func), and doc.model, its title.
func (r *renderer) globals() gad.Dict {
	link := func(label, href string) gad.Dict {
		return gad.Dict{
			"label": gad.Str(label),
			"href":  gad.Str(href),
			"link":  gad.Str("[" + label + "](" + href + ")"),
		}
	}
	arg := func(c gad.Call, i int, name string) (string, error) {
		if c.Args.Length() <= i {
			return "", fmt.Errorf("%s: want %d arguments", name, i+1)
		}
		return c.Args.Get(i).ToString(), nil
	}
	ctx := r.ctx.Context()
	mb := r.b.modelOf(r.node)
	var model gad.Object = gad.Str("")
	if mb != nil {
		model = gad.Str(mb.TTitleAuto(ctx))
	}
	admin := gad.Dict{
		"fields": gad.NewFunction("fields", func(c gad.Call) (gad.Object, error) {
			form, err := arg(c, 0, "admin.fields")
			if err != nil {
				return nil, err
			}
			if mb == nil {
				return nil, fmt.Errorf("admin.fields: %s is of no model", r.node)
			}
			return gad.Str(fieldsTable(mb, form, r.ctx)), nil
		}),
		"permissions": gad.NewFunction("permissions", func(c gad.Call) (gad.Object, error) {
			return gad.Str(r.permissions()), nil
		}),
		"permissionsTree": gad.NewFunction("permissionsTree", func(c gad.Call) (gad.Object, error) {
			return gad.Str(r.permissionsTree()), nil
		}),
		"menu": gad.NewFunction("menu", func(c gad.Call) (gad.Object, error) {
			if mb == nil {
				return nil, fmt.Errorf("admin.menu: %s is of no model", r.node)
			}
			return gad.Str(r.detailMenu(mb)), nil
		}),
		"model": gad.NewFunction("model", func(c gad.Call) (gad.Object, error) {
			id, err := arg(c, 0, "admin.model")
			if err != nil {
				return nil, err
			}
			mb := r.b.p.GetModelByID(id)
			if mb == nil {
				return nil, fmt.Errorf("admin.model: no model %q", id)
			}
			href := r.b.DocHref(r.ctx, r.b.modelNode(mb))
			if mb.Parent() == nil {
				href = mb.Info().ListingHref()
			}
			return link(mb.TTitleAuto(ctx), href), nil
		}),
		"action": gad.NewFunction("action", func(c gad.Call) (gad.Object, error) {
			id, err := arg(c, 0, "admin.action")
			if err != nil {
				return nil, err
			}
			name, err := arg(c, 1, "admin.action")
			if err != nil {
				return nil, err
			}
			mb := r.b.p.GetModelByID(id)
			if mb == nil {
				return nil, fmt.Errorf("admin.action: no model %q", id)
			}
			title, ok := actionTitle(mb, name, r.ctx)
			if !ok {
				return nil, fmt.Errorf("admin.action: the model %q has no action %q", id, name)
			}
			return link(title, r.b.DocHref(r.ctx, r.b.modelNode(mb)+"/actions/"+name)), nil
		}),
		"page": gad.NewFunction("page", func(c gad.Call) (gad.Object, error) {
			p, err := arg(c, 0, "admin.page")
			if err != nil {
				return nil, err
			}
			page := r.b.p.PagesRegistrator().GetHttpPage(p)
			if page == nil {
				return nil, fmt.Errorf("admin.page: no page %q", p)
			}
			l := link(page.TTitle(ctx), page.FullPath())
			l["menu"] = gad.Str(PageMenuPath(ctx, page))
			// the resource of its permissions, "admin:site/:/site-files:"
			// ("" with no permissions)
			// and by its unique name, out of the groups, "admin:/site-files:"
			res, unique := "", ""
			if r.ctx.R != nil {
				if v := page.ActionVerifier(r.ctx.R, ""); v != nil {
					res, unique = v.Resource(), v.PreferredResource()
				}
			}
			if unique == "" {
				unique = res // at the top of the menu: one and the same
			}
			l["resource"] = gad.Str(res)
			l["uniqueResource"] = gad.Str(unique)
			return l, nil
		}),
		"doc": gad.NewFunction("doc", func(c gad.Call) (gad.Object, error) {
			node, err := arg(c, 0, "admin.doc")
			if err != nil {
				return nil, err
			}
			title := node
			if r.tree == nil {
				r.tree = r.b.Tree(r.ctx)
			}
			if n := findNode(r.tree, node); n != nil {
				title = n.Title
			}
			return link(title, r.b.DocHref(r.ctx, node)), nil
		}),
	}
	// the application's (Func)
	for name, f := range r.b.funcs {
		admin[name] = gad.NewFunction(name, func(c gad.Call) (gad.Object, error) {
			args := make([]string, c.Args.Length())
			for i := range args {
				args[i] = c.Args.Get(i).ToString()
			}
			md, err := f(r.ctx, args...)
			if err != nil {
				return nil, fmt.Errorf("admin.%s: %w", name, err)
			}
			return gad.Str(md), nil
		})
	}
	return gad.Dict{"doc": gad.Dict{"model": model}, "admin": admin}
}

// actionTitle is the title of the action name of mb: of its detail, or of its
// listing (in bulk).
func actionTitle(mb *presets.ModelBuilder, name string, ctx *web.EventContext) (string, bool) {
	for _, a := range detailingActions(mb) {
		if a.Name() == name {
			return a.RequestTitle(mb, ctx.Context()), true
		}
	}
	for _, a := range mb.Listing().GetBulkActions() {
		if a.Name() == name {
			return a.RequestTitle(ctx.Context()), true
		}
	}
	for _, it := range rowMenuItems(mb) {
		if it.Name() == name {
			return it.TTitle(ctx.Context()), true
		}
	}
	return "", false
}

// Render is a document — the file of a node, its content as kept for the
// language — as HTML: its template run, its markdown rendered, its relative
// links pointing at the documentation (a .md) or at its package's files (an
// image).
func (r *renderer) Render(src *Source, file, content string) (string, error) {
	md, err := RenderTemplate(content, r.globals())
	if err != nil {
		return "", fmt.Errorf("%s: %w", file, err)
	}
	dir := path.Dir(file)
	rewrite := func(target string) string {
		if target == "" || strings.HasPrefix(target, "/") || strings.HasPrefix(target, "#") ||
			strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:") {
			return target
		}
		rel, frag, _ := strings.Cut(target, "#")
		rel = path.Join(dir, rel)
		if strings.HasSuffix(rel, ".md") {
			href := r.b.DocHref(r.ctx, docNode(rel))
			if frag != "" {
				href += "#" + frag
			}
			return href
		}
		return r.b.AssetHref(src, rel, r.b.locale(r.ctx))
	}
	msgs := GetMessages(r.ctx.Context())
	var out bytes.Buffer
	if err := newMarkdown(rewrite, codeLabels{msgs.CopyCode, msgs.CodeCopied, msgs.CopyCodeError}).Convert([]byte(md), &out); err != nil {
		return "", fmt.Errorf("%s: %w", file, err)
	}
	return out.String(), nil
}

// AssetHref is the URL of a file of a source (a picture), by its path under
// the language's directory: locale's, where it has it (AssetsHandler).
func (b *Builder) AssetHref(src *Source, file, locale string) string {
	return b.assetsPath + "/" + url.PathEscape(src.Package) + "/" + file + "?" + url.Values{"locale": {locale}}.Encode()
}

// PageMenuPath is where page is in the main menu, from the top: the titles of
// its groups and its own, "Site Settings → Site files" — wherever the
// application mounted it.
func PageMenuPath(ctx context.Context, page *presets.HttpPageBuilder) string {
	var parts []string
	for g := page.GetMenuGroupBuilder(); g != nil; g = g.Parent() {
		parts = append([]string{g.TTitle(ctx)}, parts...)
	}
	return strings.Join(append(parts, page.TTitle(ctx)), " → ")
}
