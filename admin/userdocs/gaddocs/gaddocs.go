// Package gaddocs is the documentation of Gad — the language of the code an
// application writes in the admin: Gad, its templates in mixed mode (Gadt)
// and its HTML templates (Gadx) — in the user docs, under the node Gad (Node):
// the documentation of the Gad module of the program (docs/en/gad, copied by
// go generate: sync.go), as written (in English, not translated: Raw), under a
// page of its own (the version of Gad of the program, its three parts).
package gaddocs

//go:generate go run sync.go

import (
	"embed"
	"io/fs"
	"path"
	"runtime/debug"
	"slices"
	"strings"
	"testing/fstest"

	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

//go:embed docs
var docsFS embed.FS

// Package is the package documented: the key of its documents.
const Package = "github.com/gad-lang/gad"

// The nodes of the documentation: Gad, and its three parts.
const (
	Node     = "gad"
	NodeGad  = "gad/gad"
	NodeGadt = "gad/gadt"
	NodeGadx = "gad/gadx"
)

// Module is the module of Gad.
const Module = "github.com/gad-lang/gad"

// Version is the version of Gad of the program (its build information); ""
// when it does not say.
func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, d := range info.Deps {
		if d.Path == Module {
			return d.Version
		}
	}
	return ""
}

// MessagesKey is the module of the words of the documentation.
const MessagesKey i18n.ModuleKey = "rvq-admin/userdocs/gaddocs"

// Messages are the words of the documentation.
type Messages struct {
	userdocs.PackageMessages
}

var (
	Messages_en_US = &Messages{PackageMessages: userdocs.PackageMessagesEn("Gad Script",
		"The language of the code written in the admin: Gad, its text templates (Gadt) and its HTML templates (Gadx).")}
	Messages_pt_BR = &Messages{PackageMessages: userdocs.PackageMessagesPtBR("Script Gad",
		"A linguagem do código escrito no admin: Gad, os seus templates de texto (Gadt) e os de HTML (Gadx).")}
)

// overview is the page of the node Gad in a language: what Gad is, its parts,
// its version (Version) — the only page written in each language.
func overview(lang string) string {
	v := Version()
	if lang == "pt-BR" {
		if v == "" {
			v = "não informada"
		}
		return "# Script Gad\n\n" +
			"Gad é a linguagem do código escrito no admin: a estrutura dos formulários, as tarjas, " +
			"os templates do site e os textos. Versão deste programa: **" + v + "**.\n\n" +
			"- [Gad](gad/README.md): a linguagem — tipos, classes, funções, módulos e a biblioteca padrão.\n" +
			"- [Gadt](gadt/README.md): templates de texto (modo misto) — texto com código entre `{% … %}` " +
			"e valores em `{%= … %}`.\n" +
			"- [Gadx](gadx/README.md): templates de HTML — tags e componentes, como os templates do site e o " +
			"corpo de uma tarja.\n\n" +
			"A documentação de Gad está em inglês. Código-fonte: https://github.com/gad-lang/gad\n"
	}
	if v == "" {
		v = "unknown"
	}
	return "# Gad Script\n\n" +
		"Gad is the language of the code written in the admin: the structure of the forms, the banners, " +
		"the site's templates and its texts. Version of this program: **" + v + "**.\n\n" +
		"- [Gad](gad/README.md): the language — types, classes, functions, modules and the standard library.\n" +
		"- [Gadt](gadt/README.md): text templates (mixed mode) — text with code in `{% … %}` and values in " +
		"`{%= … %}`.\n" +
		"- [Gadx](gadx/README.md): HTML templates — tags and components, as the site's templates and the body " +
		"of a banner.\n\n" +
		"Source: https://github.com/gad-lang/gad\n"
}

// FS is the documents: the copy of the module's (en), and the page of the
// node Gad in each language.
func FS() fs.FS {
	sub, err := fs.Sub(docsFS, "docs")
	if err != nil {
		panic(err)
	}
	return &overlayFS{top: fstest.MapFS{
		"en/gad/README.md":    {Data: []byte(overview("en"))},
		"pt-BR/gad/README.md": {Data: []byte(overview("pt-BR"))},
	}, base: sub}
}

// Source is the documentation, as written (Raw: in English, not a template).
func Source() *userdocs.Source {
	return &userdocs.Source{
		Package:     Package,
		FS:          FS(),
		MessagesKey: MessagesKey,
		Messages: map[language.Tag]i18n.Messages{
			language.English:             Messages_en_US,
			language.BrazilianPortuguese: Messages_pt_BR,
		},
		Raw: true,
	}
}

// Register puts the documentation in b: its source, and its tree (Nodes) at
// the end of the documentation's — Gad, Gadt and Gadx with their icons
// (b.AssetHref: b configured before).
func Register(b *userdocs.Builder) *userdocs.Builder {
	src := Source()
	nodes := Nodes()
	for _, n := range nodes[0].Children {
		n.Image = b.AssetHref(src, path.Join("gad/icons", path.Base(n.ID)+".svg"), "en")
	}
	return b.Register(src).Custom(userdocs.Custom{Nodes: nodes})
}

// order is the order of the documents of a directory: these first, the
// others by name.
var order = []string{"gad", "gadt", "gadx", "getting-started", "tutorial", "syntax", "conventions",
	"builtins", "stdlib", "components-and-slots", "embedding", "reflect", "formatting", "workspace-config",
	"examples", "samples"}

// Nodes is the tree of the documentation: a node a document (DIR/README.md),
// its children the documents under it.
func Nodes() []*userdocs.CustomNode {
	var walk func(dir string) []*userdocs.CustomNode
	walk = func(dir string) (nodes []*userdocs.CustomNode) {
		entries, _ := fs.ReadDir(docsFS, path.Join("docs/en", dir))
		var names []string
		for _, e := range entries {
			if e.IsDir() {
				if _, err := fs.Stat(docsFS, path.Join("docs/en", dir, e.Name(), "README.md")); err == nil {
					names = append(names, e.Name())
				}
			}
		}
		slices.SortStableFunc(names, func(a, b string) int {
			ia, ib := slices.Index(order, a), slices.Index(order, b)
			switch {
			case ia >= 0 && ib >= 0:
				return ia - ib
			case ia >= 0:
				return -1
			case ib >= 0:
				return 1
			}
			return strings.Compare(a, b)
		})
		for _, n := range names {
			id := path.Join(dir, n)
			nodes = append(nodes, &userdocs.CustomNode{ID: id, Icon: "mdi-language-go", Children: walk(id)})
		}
		return
	}
	return []*userdocs.CustomNode{{ID: Node, Icon: "mdi-code-braces", Children: walk(Node)}}
}

// overlayFS is base with the files of top over it.
type overlayFS struct {
	top  fstest.MapFS
	base fs.FS
}

func (o *overlayFS) Open(name string) (fs.File, error) {
	if _, ok := o.top[name]; ok {
		return o.top.Open(name)
	}
	f, err := o.base.Open(name)
	if err == nil {
		return f, nil
	}
	// a directory only top has (pt-BR, pt-BR/gad)
	return o.top.Open(name)
}

func (o *overlayFS) ReadDir(name string) ([]fs.DirEntry, error) {
	a, errA := fs.ReadDir(o.base, name)
	b, errB := fs.ReadDir(o.top, name)
	if errA != nil && errB != nil {
		return nil, errA
	}
	seen := map[string]bool{}
	var out []fs.DirEntry
	for _, e := range append(b, a...) {
		if !seen[e.Name()] {
			seen[e.Name()] = true
			out = append(out, e)
		}
	}
	slices.SortFunc(out, func(x, y fs.DirEntry) int { return strings.Compare(x.Name(), y.Name()) })
	return out, nil
}
