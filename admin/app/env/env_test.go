package env

import (
	"os"
	"strings"
	"testing"
)

// The env: of an app.yml: strings with ${VAR} expanded, numbers and bools as
// written, a list joined as PATH; its features: as they are.
func TestParseAppYAML(t *testing.T) {
	t.Setenv("ENV_TEST_HOME", "/home/x")
	env, feats, err := ParseAppYAML([]byte(`
env:
  NAME: site
  DIR: ${ENV_TEST_HOME}/data
  PORT: 9000
  DEBUG: true
  EMPTY:
  PATHS: [bin, "${ENV_TEST_HOME}/bin", 3]
features:
  hermon.forms.drafts: true
  x:
    y: false
`))
	if err != nil {
		t.Fatal(err)
	}
	sep := string(os.PathListSeparator)
	want := map[string]string{
		"NAME":  "site",
		"DIR":   "/home/x/data",
		"PORT":  "9000",
		"DEBUG": "true",
		"EMPTY": "",
		"PATHS": strings.Join([]string{"bin", "/home/x/bin", "3"}, sep),
	}
	for k, v := range want {
		if env[k] != v {
			t.Errorf("%s = %q, want %q", k, env[k], v)
		}
	}
	if len(env) != len(want) {
		t.Errorf("env: %v", env)
	}
	if feats["hermon.forms.drafts"] != true {
		t.Errorf("features: %v", feats)
	}
}

// What is not a value is refused, saying the key.
func TestParseAppYAMLErrors(t *testing.T) {
	for src, want := range map[string]string{
		"env:\n  A: {b: 1}\n": "env.A: a map is not a value",
		"env:\n  A: [[1]]\n":  "env.A: a list in a list",
		"env: [1]\n":          "error parsing app.yml",
		"env:\n  A: ${B\n":    "env.A",
	} {
		if _, _, err := ParseAppYAML([]byte(src)); err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%q: %v, want %q", src, err, want)
		}
	}
}

func TestIsYAML(t *testing.T) {
	for name, want := range map[string]bool{"app.yml": true, "x/APP.YAML": true, ".env": false, "a.env": false} {
		if IsYAML(name) != want {
			t.Errorf("%s: %v", name, !want)
		}
	}
}
