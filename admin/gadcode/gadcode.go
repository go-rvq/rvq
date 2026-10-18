// Package gadcode serves what the admin's editor of gad code (vx-gad-code)
// asks of the code in it: the code formatted as `gad fmt` formats it, its
// structure — the declarations, as a tree (langsym.Outline, what `gad
// symbols` gives an IDE) — and its documentation (`gad doc`).
//
// The code is of a language: "gad" (plain), "gadx" (a template of
// indentation) or "gadt" (a mixed template; its delimiters `{%`/`%}`, or
// the ones the editor says, as the admin's texts' `{`/`}`).
package gadcode

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gad-lang/gad/gadx/node"
	gadxparser "github.com/gad-lang/gad/gadx/parser"
	"github.com/gad-lang/gad/langsym"
	"github.com/gad-lang/gad/parser"
	gadnode "github.com/gad-lang/gad/parser/node"
	"github.com/gad-lang/gad/parser/source"
	"github.com/gad-lang/gad/web/gadbridge"
	"github.com/go-rvq/rvq/admin/presets"
)

// The languages of the code.
const (
	Gad  = "gad"
	Gadx = "gadx"
	Gadt = "gadt"
)

// Path is where its operations are served, under the admin's prefix:
// POST <prefix>/__gadcode/{format,outline,doc}.
const Path = "/__gadcode/"

// Code is the code an operation is asked of.
type Code struct {
	Source   string `json:"source"`
	Language string `json:"language"`
	// TemplateStart and TemplateEnd are the delimiters of a gadt's code;
	// `{%` / `%}` when empty.
	TemplateStart string `json:"templateStart,omitempty"`
	TemplateEnd   string `json:"templateEnd,omitempty"`
}

func (c Code) delimiter() parser.MixedDelimiter {
	d := parser.DefaultMixedDelimiter
	if c.TemplateStart != "" && c.TemplateEnd != "" {
		d = parser.MixedDelimiter{Start: []rune(c.TemplateStart), End: []rune(c.TemplateEnd)}
	}
	return d
}

// Format is the code as `gad fmt` formats it, and whether it could: what
// does not parse — or is of no gad language — is answered as it is.
// Formatted code keeps whether it ended in a new line.
func Format(c Code) (string, bool) {
	src := c.Source
	if strings.TrimSpace(src) == "" {
		return src, false
	}
	var out string
	switch c.Language {
	case Gad, Gadx:
		r := gadbridge.FormatSource(src, c.Language)
		if !r.OK {
			return src, false
		}
		out = r.Source
	case Gadt:
		file, sf, err := c.parse()
		if err != nil || file == nil {
			return src, false
		}
		out = gadnode.Code(file.Stmts,
			gadnode.CodeWithFlags(gadnode.CodeWriteContextFlagFormat),
			gadnode.CodeWithComments(sf, file.Comments),
		)
	default:
		return src, false
	}
	if !strings.HasSuffix(src, "\n") {
		out = strings.TrimRight(out, "\n")
	}
	return out, true
}

// Outline is the structure of the code: its declarations, as a tree — what
// parsed of it, while it is being written.
func Outline(c Code) []langsym.OutlineSym {
	file, sf, _ := c.parse()
	var syms []langsym.OutlineSym
	if file != nil {
		// an exported declaration (`export class default { … }`) is the
		// declaration: what it exports is in its prelude
		stmts := make([]gadnode.Stmt, 0, len(file.Stmts))
		for _, st := range file.Stmts {
			if e, ok := st.(*gadnode.ExportStmt); ok && e.Prelude != nil {
				st = e.Prelude
			}
			stmts = append(stmts, st)
		}
		f := *file
		f.Stmts = stmts
		syms = langsym.Outline(&f, sf)
	}
	if syms == nil {
		syms = []langsym.OutlineSym{}
	}
	return syms
}

// Doc is the documentation of the code (`gad doc`), as HTML.
func Doc(c Code) (string, error) {
	typ := c.Language
	if typ == Gadt {
		typ = "gadTemplate"
	}
	return gadbridge.DocHTML(c.Source, typ)
}

// parse reads the code as its language does — tolerant: a file being
// written answers what parsed of it.
func (c Code) parse() (*parser.File, *source.File, error) {
	fs := source.NewFileSet()
	sf := fs.AddFileData("(editor)."+c.Language, -1, []byte(c.Source))
	if c.Language == Gadx {
		gf, err := gadxparser.NewParser(sf).ParseFile()
		if gf == nil {
			return nil, sf, err
		}
		return &parser.File{InputFile: sf, Stmts: node.Convert(gf.Stmts)}, sf, err
	}
	po := &parser.ParserOptions{Mode: parser.ParseComments}
	var so *parser.ScannerOptions
	if c.Language == Gadt {
		po.Mode |= parser.ParseMixed
		so = &parser.ScannerOptions{Mode: parser.ScanMixed | parser.ScanConfigDisabled, MixedDelimiter: c.delimiter()}
	}
	file, err := parser.NewParserWithOptions(sf, po, so).ParseFileTolerant()
	return file, sf, err
}

// Setup serves the operations under the admin's prefix (Path) — behind its
// login, as the rest of it — and registers the editor's words (Messages).
func Setup(p *presets.Builder) {
	ConfigureMessages(p.I18n())
	p.MuxSetup(func(prefix string, mux *http.ServeMux) {
		mux.HandleFunc("POST "+prefix+Path+"{op}", serve)
	})
}

// API is where the editor asks: the admin's prefix and Path.
func API(prefix string) string { return prefix + Path }

func serve(w http.ResponseWriter, r *http.Request) {
	var c Code
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 2<<20)).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var out any
	switch r.PathValue("op") {
	case "format":
		src, ok := Format(c)
		out = map[string]any{"ok": ok, "source": src}
	case "outline":
		out = Outline(c)
	case "doc":
		html, err := Doc(c)
		res := map[string]any{"html": html}
		if err != nil {
			res["error"] = err.Error()
		}
		out = res
	default:
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(out)
}
