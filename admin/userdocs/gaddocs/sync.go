//go:build ignore

// sync copies the documentation of the Gad module this one requires (go list
// -m github.com/gad-lang/gad) into docs/en/gad, as the user docs read it: a
// document a node, NAME.md as NAME/README.md, the links between them made
// again for where they went; the other files (pictures, samples) as they are.
//
//	gad/gad/…   doc/ (the language, the stdlib, the samples)
//	gad/gadt/   doc/samples/template.md (templates in mixed mode)
//	gad/gadx/…  gadx/README.md and gadx/docs/
//	gad/icons/  the icons of Gad, Gadt, Gadx (assets/identity)
//
// Run by go generate (gaddocs.go).
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const out = "docs/en/gad"

// refPlaceholder is gaddocs.RefPlaceholder (this program is not the package).
const refPlaceholder = "GAD_REF"

// mapped is where a file of the module goes, under out; "" left out.
func mapped(rel string) string {
	dirOf := func(prefix, into string) (string, bool) {
		r, ok := strings.CutPrefix(rel, prefix)
		if !ok {
			return "", false
		}
		return path.Join(into, node(r)), true
	}
	if rel == "doc/samples/template.md" {
		return "gadt/README.md"
	}
	if rel == "gadx/README.md" {
		return "gadx/README.md"
	}
	if d, ok := dirOf("gadx/docs/", "gadx"); ok {
		return d
	}
	if d, ok := dirOf("doc/", "gad"); ok {
		return d
	}
	return ""
}

// node is the place of a file of a directory: NAME.md as NAME/README.md (a
// node of its own); README.md and any other file as they are.
func node(rel string) string {
	if strings.HasSuffix(rel, ".md") && path.Base(rel) != "README.md" {
		return strings.TrimSuffix(rel, ".md") + "/README.md"
	}
	return rel
}

// partTitles are the titles of the documents of the parts of Gad.
var partTitles = map[string]string{"gad/README.md": "Gad", "gadt/README.md": "Gadt", "gadx/README.md": "Gadx"}

var heading = regexp.MustCompile(`(?m)^# +(.+)$`)

// retitle is the markdown b with the title (its first # heading) title, the
// one it had below it, in italics, when another.
func retitle(b []byte, title string) []byte {
	m := heading.FindSubmatchIndex(b)
	if m == nil {
		return append([]byte("# "+title+"\n\n"), b...)
	}
	old := strings.TrimSpace(string(b[m[2]:m[3]]))
	head := "# " + title
	if !strings.EqualFold(old, title) && old != "doc" {
		head += "\n\n*" + old + "*"
	}
	return append(append(append([]byte{}, b[:m[0]]...), head...), b[m[1]:]...)
}

var link = regexp.MustCompile(`(\]\(|src=")([^)"#\s]+)([#][^)"\s]*)?`)

func main() {
	dirB, err := exec.Command("go", "list", "-m", "-f", "{{.Dir}} {{.Version}}", "github.com/gad-lang/gad").Output()
	if err != nil {
		panic(err)
	}
	src, _, _ := strings.Cut(strings.TrimSpace(string(dirB)), " ")
	// the repository at the version of Gad of the program that serves the
	// documents: GAD_REF, put by gaddocs.FS (Ref) when they are read — never
	// the version this copy was made from
	repo := "https://github.com/gad-lang/gad/blob/" + refPlaceholder + "/"
	if err := os.RemoveAll(out); err != nil {
		panic(err)
	}
	var files []string
	for _, root := range []string{"doc", "gadx/README.md", "gadx/docs"} {
		_ = filepath.WalkDir(filepath.Join(src, root), func(p string, d os.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return err
			}
			rel, _ := filepath.Rel(src, p)
			if mapped(filepath.ToSlash(rel)) != "" {
				files = append(files, filepath.ToSlash(rel))
			}
			return nil
		})
	}
	copied := map[string]bool{}
	for _, rel := range files {
		copied[rel] = true
	}
	for _, rel := range files {
		to := mapped(rel)
		b, err := os.ReadFile(filepath.Join(src, rel))
		if err != nil {
			panic(err)
		}
		if strings.HasSuffix(rel, ".md") {
			// the links: where their targets went, from where this one went
			b = link.ReplaceAllFunc(b, func(m []byte) []byte {
				g := link.FindSubmatch(m)
				target := string(g[2])
				if strings.Contains(target, "://") || strings.HasPrefix(target, "/") || strings.HasPrefix(target, "mailto:") {
					return m
				}
				abs := path.Clean(path.Join(path.Dir(rel), target))
				dest := ""
				if copied[abs] {
					dest = mapped(abs)
				}
				if dest == "" {
					// a file not copied (code, of another directory, or not
					// there): the repository's
					if strings.HasPrefix(abs, "..") {
						return m
					}
					return []byte(string(g[1]) + repo + abs + string(g[3]))
				}
				r, err := filepath.Rel(path.Dir(to), dest)
				if err != nil {
					return m
				}
				return []byte(string(g[1]) + filepath.ToSlash(r) + string(g[3]))
			})
		}
		// the parts' titles: their names (the menu's), the module's title
		// below
		if title, ok := partTitles[to]; ok {
			b = retitle(b, title)
		}
		dst := filepath.Join(out, filepath.FromSlash(to))
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			panic(err)
		}
		if err := os.WriteFile(dst, b, 0o644); err != nil {
			panic(err)
		}
	}
	// a directory of documents without its own (stdlib): a page of them, so
	// the tree opens them
	_ = filepath.WalkDir(out, func(p string, d os.DirEntry, err error) error {
		if err != nil || !d.IsDir() || p == out {
			return err
		}
		if _, err := os.Stat(filepath.Join(p, "README.md")); err == nil {
			return nil
		}
		entries, _ := os.ReadDir(p)
		var items []string
		for _, e := range entries {
			if _, err := os.Stat(filepath.Join(p, e.Name(), "README.md")); e.IsDir() && err == nil {
				items = append(items, "- ["+e.Name()+"]("+e.Name()+"/README.md)\n")
			}
		}
		if len(items) > 0 {
			b := "# " + filepath.Base(p) + "\n\n" + strings.Join(items, "")
			if err := os.WriteFile(filepath.Join(p, "README.md"), []byte(b), 0o644); err != nil {
				panic(err)
			}
		}
		return nil
	})
	// the icons of Gad, Gadt and Gadx (assets/identity): of their nodes
	for _, name := range []string{"gad", "gadt", "gadx"} {
		b, err := os.ReadFile(filepath.Join(src, "assets/identity", name+"-32.svg"))
		if err != nil {
			panic(err)
		}
		dst := filepath.Join(out, "icons", name+".svg")
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			panic(err)
		}
		if err := os.WriteFile(dst, b, 0o644); err != nil {
			panic(err)
		}
	}
	fmt.Printf("gaddocs: %d files of %s\n", len(files), src)
}
