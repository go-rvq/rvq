package i18n_test

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// The convention (README.md, "Code Conventions"): what the end user reads is
// available in English AND in Brazilian Portuguese. A module's messages are
// stored whole, with no per-field fallback, so a module registered for only one
// of them leaves the other language reading someone else's language — or blanks.
//
// This walks the source instead of the i18n.Builder because registration happens
// all over the tree, in each package's own setup.

var registerRe = regexp.MustCompile(`RegisterForModules?\(\s*language\.(\w+)\s*,\s*([\w.]+)\s*,`)

// Upstream demos and doc templates: they are sample applications, not our
// interface, and several of them no longer build.
var skipDirs = []string{
	"admin/docs",
	"admin/example",
	"js/integration_tests",
}

func TestEveryMessagesModuleHasEnglishAndPortuguese(t *testing.T) {
	root := repoRoot(t)

	type site struct {
		langs map[string]bool
		file  string
	}
	modules := map[string]*site{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		if info.IsDir() {
			if info.Name() == "node_modules" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			for _, skip := range skipDirs {
				if rel == skip {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, m := range registerRe.FindAllStringSubmatch(string(b), -1) {
			lang, key := m[1], m[2]
			// the module key as written; two packages may both call it
			// MessagesKey, so keep the directory with it
			id := filepath.Dir(rel) + ":" + key
			s := modules[id]
			if s == nil {
				s = &site{langs: map[string]bool{}, file: rel}
				modules[id] = s
			}
			s.langs[lang] = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(modules) < 10 {
		t.Fatalf("only %d modules found — the scan is not seeing the source", len(modules))
	}

	var missing []string
	for id, s := range modules {
		english := s.langs["English"] || s.langs["AmericanEnglish"]
		portuguese := s.langs["BrazilianPortuguese"]
		switch {
		case english && portuguese:
		case !portuguese:
			missing = append(missing, id+" (no pt-BR) at "+s.file)
		default:
			missing = append(missing, id+" (no en) at "+s.file)
		}
	}
	sort.Strings(missing)

	if len(missing) > 0 {
		t.Errorf("these i18n modules are not in both languages — see \"Code Conventions\" in README.md:\n\t%s",
			strings.Join(missing, "\n\t"))
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found above " + dir)
		}
		dir = parent
	}
}
