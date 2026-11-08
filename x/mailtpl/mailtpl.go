// Package mailtpl is the mails written in Gadx, one template each language:
// the packages register their default templates (Register), an application
// gives its files over them (SetFS) — a file of its own wins, the defaults
// still there for what it lacks —, and Render makes a mail of a template in
// the language of whom it goes to.
//
// The templates live in the folders of the files of an application:
//
//	templates/mail/<package>/<name>.<language>.gadx
//
// e.g. templates/mail/shared/confirm.pt-BR.gadx. A template's <title> is the
// mail's subject (taken out of its body); its text is made of its HTML.
package mailtpl

import (
	"errors"
	"fmt"
	"html"
	"io/fs"
	"regexp"
	"slices"
	"sort"
	"strings"
	"sync"

	"github.com/gad-lang/gad"
	"github.com/go-rvq/rvq/x/i18n/gadxtpl"
	"golang.org/x/text/language"
)

// Dir is the folder of the templates of mails in the files.
const Dir = "templates/mail"

// Ext is the extension of a template.
const Ext = ".gadx"

// ErrNotFound: no template of the name in any language.
var ErrNotFound = errors.New("mailtpl: no template")

// Renderer renders the source of a template (its path, for its imports and
// its errors) with data.
type Renderer interface {
	Render(path string, src []byte, data gad.Dict) (string, error)
}

// RendererFunc is a function as a Renderer.
type RendererFunc func(path string, src []byte, data gad.Dict) (string, error)

func (f RendererFunc) Render(path string, src []byte, data gad.Dict) (string, error) {
	return f(path, src, data)
}

// standalone renders a template alone (gadxtpl): no imports.
var standalone = RendererFunc(func(_ string, src []byte, data gad.Dict) (string, error) {
	return gadxtpl.Template(src).Render(data)
})

type layer struct {
	name string
	fs   fs.FS
}

var (
	mu       sync.RWMutex
	layers   []layer
	appFS    fs.FS
	renderer Renderer = standalone
	fallback          = "en-US"
)

// Register adds the default templates of a package — fsys a tree with the
// folders of the files (templates/mail/…) —, from its init. Of a path in
// several, the one registered last.
func Register(name string, fsys fs.FS) {
	mu.Lock()
	defer mu.Unlock()
	for i, l := range layers {
		if l.name == name {
			layers[i].fs = fsys
			return
		}
	}
	layers = append(layers, layer{name, fsys})
}

// SetFS gives the files of the application: above the defaults. nil: the
// defaults alone.
func SetFS(fsys fs.FS) {
	mu.Lock()
	defer mu.Unlock()
	appFS = fsys
}

// SetRenderer gives the application's renderer (e.g. one that lets a
// template import a layout of its files). nil: the templates alone.
func SetRenderer(r Renderer) {
	mu.Lock()
	defer mu.Unlock()
	if r == nil {
		r = standalone
	}
	renderer = r
}

// SetFallback is the language of the templates when the mail's has none
// (default en-US).
func SetFallback(lang string) {
	mu.Lock()
	defer mu.Unlock()
	fallback = lang
}

// Defaults are the default templates of every package, as one tree.
func Defaults() fs.FS {
	mu.RLock()
	defer mu.RUnlock()
	fss := make([]fs.FS, 0, len(layers))
	for i := len(layers) - 1; i >= 0; i-- {
		fss = append(fss, layers[i].fs)
	}
	return Union(fss...)
}

// Missing are the paths of the default templates that app lacks — to bring
// into it the ones a new version adds, never touching its own.
func Missing(app fs.FS) ([]string, error) {
	var out []string
	err := fs.WalkDir(Defaults(), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			return err
		}
		if d.IsDir() {
			return nil
		}
		if _, err := fs.Stat(app, p); errors.Is(err, fs.ErrNotExist) {
			out = append(out, p)
		}
		return nil
	})
	return out, err
}

// Mail is a mail made of a template: its subject, its HTML and its text,
// the language and the path of the template.
type Mail struct {
	Subject, HTML, Text string
	Language, Path      string
}

// Languages are the languages a template of lang is looked for in: lang, its
// base, the fallback, its base.
func Languages(lang string) []string {
	mu.RLock()
	fb := fallback
	mu.RUnlock()
	var out []string
	add := func(s string) {
		if s != "" && !slices.Contains(out, s) {
			out = append(out, s)
		}
	}
	for _, l := range []string{lang, fb} {
		if t, err := language.Parse(l); err == nil {
			add(t.String())
			b, _ := t.Base()
			add(b.String())
		} else {
			add(l)
		}
	}
	return out
}

// Find is the template of name (package/name) in lang: its path and
// source, from the application's files, else the defaults; in the first of
// Languages(lang) there is.
func Find(name, lang string) (path string, src []byte, err error) {
	mu.RLock()
	app := appFS
	mu.RUnlock()
	fss := []fs.FS{Defaults()}
	if app != nil {
		fss = append([]fs.FS{app}, fss...)
	}
	for _, l := range Languages(lang) {
		p := Dir + "/" + name + "." + l + Ext
		for _, f := range fss {
			if b, err := fs.ReadFile(f, p); err == nil {
				return p, b, nil
			}
		}
	}
	return "", nil, fmt.Errorf("%w %s (%s)", ErrNotFound, name, lang)
}

var titleRe = regexp.MustCompile(`(?is)<title[^>]*>(.*?)</title>`)

// Render is the mail of the template name (package/name) in lang, data its
// globals.
func Render(name, lang string, data gad.Dict) (*Mail, error) {
	path, src, err := Find(name, lang)
	if err != nil {
		return nil, err
	}
	mu.RLock()
	r := renderer
	mu.RUnlock()
	out, err := r.Render(path, src, data)
	if err != nil {
		return nil, fmt.Errorf("mailtpl %s: %w", path, err)
	}
	m := &Mail{Path: path, Language: langOf(path)}
	if t := titleRe.FindStringSubmatch(out); t != nil {
		m.Subject = strings.TrimSpace(html.UnescapeString(stripTags(t[1])))
		out = titleRe.ReplaceAllString(out, "")
	}
	m.HTML = strings.TrimSpace(out)
	m.Text = Text(m.HTML)
	return m, nil
}

// langOf is the language of the path of a template.
func langOf(path string) string {
	base := strings.TrimSuffix(path, Ext)
	if i := strings.LastIndex(base, "."); i >= 0 {
		return base[i+1:]
	}
	return ""
}

var (
	tagRe   = regexp.MustCompile(`(?s)<[^>]*>`)
	breakRe = regexp.MustCompile(`(?i)<\s*(br|/p|/div|/h[1-6]|/li|/tr)\s*/?>`)
	liRe    = regexp.MustCompile(`(?i)<\s*li[^>]*>`)
	aRe     = regexp.MustCompile(`(?is)<a\s[^>]*href="([^"]*)"[^>]*>(.*?)</a>`)
	spaceRe = regexp.MustCompile(`[ \t]+`)
	linesRe = regexp.MustCompile(`\n{3,}`)
)

func stripTags(s string) string { return tagRe.ReplaceAllString(s, "") }

// Text is the text of a mail of HTML: its lines, its links with their
// addresses, no tags.
func Text(s string) string {
	s = aRe.ReplaceAllStringFunc(s, func(a string) string {
		m := aRe.FindStringSubmatch(a)
		label := strings.TrimSpace(stripTags(m[2]))
		if label == "" || label == m[1] {
			return m[1]
		}
		return label + " (" + m[1] + ")"
	})
	s = liRe.ReplaceAllString(s, "- ")
	s = breakRe.ReplaceAllString(s, "\n")
	s = html.UnescapeString(stripTags(s))
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimSpace(spaceRe.ReplaceAllString(l, " "))
	}
	return strings.TrimSpace(linesRe.ReplaceAllString(strings.Join(lines, "\n"), "\n\n"))
}

// Union is the trees as one, file by file: of a path in several, the first
// tree's.
func Union(fss ...fs.FS) fs.FS { return union(fss) }

type union []fs.FS

func (u union) Open(name string) (fs.File, error) {
	var dirs []fs.FS
	for _, f := range u {
		info, err := fs.Stat(f, name)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			return f.Open(name)
		}
		dirs = append(dirs, f)
	}
	if len(dirs) == 0 {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return dirs[0].Open(name)
}

func (u union) ReadDir(name string) ([]fs.DirEntry, error) {
	seen := map[string]bool{}
	var out []fs.DirEntry
	found := false
	for _, f := range u {
		entries, err := fs.ReadDir(f, name)
		if err != nil {
			continue
		}
		found = true
		for _, e := range entries {
			if !seen[e.Name()] {
				seen[e.Name()] = true
				out = append(out, e)
			}
		}
	}
	if !found {
		return nil, &fs.PathError{Op: "readdir", Path: name, Err: fs.ErrNotExist}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name() < out[j].Name() })
	return out, nil
}

func (u union) Stat(name string) (fs.FileInfo, error) {
	for _, f := range u {
		if info, err := fs.Stat(f, name); err == nil {
			return info, nil
		}
	}
	return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrNotExist}
}

func (u union) ReadFile(name string) ([]byte, error) {
	for _, f := range u {
		if b, err := fs.ReadFile(f, name); err == nil {
			return b, nil
		}
	}
	return nil, &fs.PathError{Op: "read", Path: name, Err: fs.ErrNotExist}
}
