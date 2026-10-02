package userdocs

import (
	"bytes"
	"fmt"
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
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// newMarkdown is the renderer of a document, its links and images resolved by
// rewrite: Gad's markdown (gadx.NewMarkdown — every bundled extension: GFM,
// typographer, definition lists, footnotes; auto heading ids; HTML kept, the
// documents being the code's and the administrators').
func newMarkdown(rewrite func(dest string) string) goldmark.Markdown {
	md := gadx.NewMarkdown()
	md.Parser().AddOptions(gparser.WithASTTransformers(util.Prioritized(linkRewriter(rewrite), 100)))
	return md
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
var templateGlobals = []string{"admin"}

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
// page; NAME.md of an action.
func DocFile(node string) string {
	if strings.Contains("/"+node, "/actions/") {
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
}

// globals are the template's: admin.model, admin.action, admin.page and
// admin.doc — each a record of its label, its href and a markdown link.
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
	return gad.Dict{"admin": gad.Dict{
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
			return link(mb.TTitlePlural(ctx), href), nil
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
			return link(page.TTitle(ctx), page.FullPath()), nil
		}),
		"doc": gad.NewFunction("doc", func(c gad.Call) (gad.Object, error) {
			node, err := arg(c, 0, "admin.doc")
			if err != nil {
				return nil, err
			}
			title := node
			if n := findNode(r.b.Tree(r.ctx), node); n != nil {
				title = n.Title
			}
			return link(title, r.b.DocHref(r.ctx, node)), nil
		}),
	}}
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
		return r.b.AssetHref(src, rel)
	}
	var out bytes.Buffer
	if err := newMarkdown(rewrite).Convert([]byte(md), &out); err != nil {
		return "", fmt.Errorf("%s: %w", file, err)
	}
	return out.String(), nil
}

// AssetHref is the URL of a file of a source (a picture), by its path under
// the language's directory.
func (b *Builder) AssetHref(src *Source, file string) string {
	return b.assetsPath + "/" + url.PathEscape(src.Package) + "/" + file
}
