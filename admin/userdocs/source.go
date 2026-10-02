// Package userdocs is the documentation of the admin for whoever uses it: a
// page, "Documentation", whose tree is the side menu's — its groups, the
// models in them and their actions, the models nested in them, the pages —,
// each node opening the document that explains it.
//
// The documents come from the packages: each one has a user_docs directory, a
// Go package of its own (user_docs), with its documents in markdown under the
// directory of the language they are written in — en, the SOURCE language,
// and its translations beside it (pt-BR), each with its pictures:
//
//	PACKAGE/user_docs/
//	  docs.go                       //go:embed en; Source()
//	  messages.go                   the words of the package's documentation
//	  en/
//	    MODEL_ID/README.md          the model (MODEL_ID: its menu id)
//	    MODEL_ID/actions/ACTION.md  an action of it
//	    MODEL_ID/children/CHILD/…   a model nested in it, the same way,
//	                                recursively
//	    MODEL_ID/images/*.png       its pictures, referred to by the documents
//	    actions/ACTION.md           an action, of whatever model has it
//	                                (restore), when the model's documents do
//	                                not explain it
//	    children/WORD/README.md     a nested model, of whatever model has it,
//	                                by the last word of its key (revisions)
//	    pages/PAGE/README.md        a page of the admin (db-tools)
//	    groups/GROUP/README.md      a group of the menu
//	  pt-BR/                        the same, in Portuguese: a document
//	                                missing there is translated from en's,
//	                                a picture missing there is en's
//
// A document is a Gad template ({% … %} code, {%= expr %} a value), so it
// says where things are in the admin as the admin says it:
// {%= admin.model("admin_locales").link %} is a link to the listing of the
// model, by its label; admin.action, admin.page and admin.doc are the others.
//
// Documents of no part of the admin — a guide, a tutorial — go in the tree
// where they are registered (Builder.Custom): under a node (a group, a model;
// the top), first or last among its children, each node a document of its ID
// (ID/README.md) in a package, titled by its heading in the request's
// language:
//
//	docs.Custom(userdocs.Custom{First: true, Nodes: []*userdocs.CustomNode{
//		{ID: "guides/getting-started", Icon: "mdi-school-outline",
//			Children: []*userdocs.CustomNode{{ID: "guides/getting-started/01-post-types"}}},
//	}})
//
// The documents are kept in the database, per admin language (UserDoc, one per
// language, and its UserDocPackage, one per package): written from the code's
// on boot (Sync) — the initial value: the language's own documents, else the
// source's translated —, and editable there.
package userdocs

import (
	"fmt"
	"io/fs"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// Source is the documentation of a package: its user_docs directory.
type Source struct {
	// Package is the Go package the documentation is of: the key of its
	// UserDocPackage.
	Package string
	// FS is the user_docs directory: the directory of the source language,
	// holding the documents and their pictures.
	FS fs.FS
	// MessagesKey and Messages are the words of the package's documentation
	// (its user_docs/messages.go), per language: registered in the i18n as a
	// module of its own, they give the labels and hints of its form.
	MessagesKey i18n.ModuleKey
	Messages    map[language.Tag]i18n.Messages
}

// File is a document: its path under the language's directory, and its
// content.
type File struct {
	Path    string `yaml:"path"`
	Content string `yaml:"content"`
}

// Langs are the languages of the source: its directories whose name is a
// language (en, pt-BR), sorted.
func (s *Source) Langs() ([]string, error) {
	entries, err := fs.ReadDir(s.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("userdocs %s: %w", s.Package, err)
	}
	var langs []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := language.Parse(e.Name()); err == nil {
			langs = append(langs, e.Name())
		}
	}
	if len(langs) == 0 {
		return nil, fmt.Errorf("userdocs %s: no directory of a language", s.Package)
	}
	sort.Strings(langs)
	return langs, nil
}

// Lang is the source language: en when the source has it, else its only
// language — the one the others are translated from.
func (s *Source) Lang() (string, error) {
	langs, err := s.Langs()
	if err != nil {
		return "", err
	}
	if slices.Contains(langs, "en") {
		return "en", nil
	}
	if len(langs) != 1 {
		return "", fmt.Errorf("userdocs %s: no en among the languages %q: which is the source?", s.Package, langs)
	}
	return langs[0], nil
}

// LangOf is the language of the source written in locale — the same, else
// one of its base (pt for pt-BR) —, "" when none is.
func (s *Source) LangOf(locale string) string {
	langs, err := s.Langs()
	if err != nil {
		return ""
	}
	for _, l := range langs {
		if strings.EqualFold(l, locale) {
			return l
		}
	}
	for _, l := range langs {
		if sameLanguage(l, locale) {
			return l
		}
	}
	return ""
}

// Files are the documents of the source language, by path.
func (s *Source) Files() ([]File, error) {
	lang, err := s.Lang()
	if err != nil {
		return nil, err
	}
	return s.FilesOf(lang)
}

// FilesOf are the documents written in lang, by path.
func (s *Source) FilesOf(lang string) ([]File, error) {
	var files []File
	err := fs.WalkDir(s.FS, lang, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || path.Ext(p) != ".md" {
			return err
		}
		b, err := fs.ReadFile(s.FS, p)
		if err != nil {
			return err
		}
		files = append(files, File{Path: strings.TrimPrefix(p, lang+"/"), Content: string(b)})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("userdocs %s: %w", s.Package, err)
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

// Asset is the language whose directory has the file — a picture — for
// locale: locale's own, else the source's.
func (s *Source) Asset(locale, file string) (string, error) {
	if l := s.LangOf(locale); l != "" {
		if _, err := fs.Stat(s.FS, path.Join(l, file)); err == nil {
			return l, nil
		}
	}
	return s.Lang()
}
