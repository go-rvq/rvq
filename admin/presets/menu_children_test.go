package presets

import (
	"encoding/json"
	"strings"
	"testing"
)

// A node of children to load goes with an empty array of them — what makes
// VTreeview load them —; another, without.
func TestMenuNodeJSON(t *testing.T) {
	lazy, _ := json.Marshal(&MenuNode{Title: "Docs", Value: "/docs", lazy: true})
	if !strings.Contains(string(lazy), `"children":[]`) {
		t.Errorf("lazy %s", lazy)
	}
	leaf, _ := json.Marshal(&MenuNode{Title: "Posts", Value: "/posts"})
	if strings.Contains(string(leaf), "children") {
		t.Errorf("leaf %s", leaf)
	}
	loaded, _ := json.Marshal([]*MenuNode{{Value: "/docs", lazy: true, Children: []*MenuNode{{Value: "/docs/a"}}}})
	if !strings.Contains(string(loaded), `"children":[{"title":"","value":"/docs/a"}]`) {
		t.Errorf("loaded %s", loaded)
	}
}

func TestMenuPath(t *testing.T) {
	tree := []*MenuNode{{Value: "a", Children: []*MenuNode{{Value: "b", Children: []*MenuNode{{Value: "c"}}}}}}
	if p := menuPath(tree, "c"); strings.Join(p, ",") != "a,b" {
		t.Errorf("path %q", p)
	}
	if p := menuPath(tree, "a"); p == nil || len(p) != 0 {
		t.Errorf("path of the top %q", p)
	}
	if menuPath(tree, "x") != nil || !tree[0].has("c") || tree[0].has("x") {
		t.Error("a value not there was found")
	}
}
