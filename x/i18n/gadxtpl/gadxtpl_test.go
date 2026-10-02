package gadxtpl

import (
	"context"
	"strings"
	"testing"

	"github.com/gad-lang/gad"
	h "github.com/go-rvq/htmlgo"
)

func TestTemplateRenders(t *testing.T) {
	tpl := Template("@main\np.usage Hello {= Name}\n")
	out, err := tpl.Render(gad.Dict{"Name": gad.Str("<Gadx>")})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `<p class="usage">Hello &lt;Gadx&gt;</p>`) {
		t.Errorf("rendered %q", out)
	}
	// the compilation is cached: the same source renders the same
	if again, _ := tpl.Render(gad.Dict{"Name": gad.Str("x")}); !strings.Contains(again, "Hello x") {
		t.Errorf("second render %q", again)
	}
	if out, err := Template("").Render(nil); err != nil || out != "" {
		t.Errorf("empty template: %q %v", out, err)
	}
}

func TestTemplateErrorIsShown(t *testing.T) {
	c := Template("@main\np {= notDefined }\n").Component(nil)
	if s := h.MustString(c, context.Background()); !strings.Contains(s, "<pre") {
		t.Errorf("an error should render as a <pre>: %s", s)
	}
}
