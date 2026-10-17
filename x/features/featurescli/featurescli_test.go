package featurescli

import (
	"os"
	"strings"
	"testing"

	"github.com/go-rvq/rvq/x/features"
)

// The features on alone; all, every one compiled, why one is off; one not
// compiled — and one depending on it — nowhere.
func TestPrint(t *testing.T) {
	site := features.Register("zzcli.site", "Site", "").CompiledIf(true)
	forms := features.Register("zzcli.forms", "Forms", "").CompiledIf(true).DependsOn(site)
	drafts := features.Register("zzcli.drafts", "Drafts", "").CompiledIf(true).DependsOn(forms)
	hidden := features.Register("zzcli.hidden", "Hidden", "")
	features.Register("zzcli.needshidden", "NeedsHidden", "").CompiledIf(true).DependsOn(hidden)
	features.Set(site.Name, true)
	features.Set(forms.Name, true)
	features.Set(hidden.Name, true)

	var b strings.Builder
	if err := PrintOn(&b); err != nil {
		t.Fatal(err)
	}
	if out := b.String(); !strings.Contains(out, "zzcli.site") || !strings.Contains(out, "zzcli.forms") ||
		strings.Contains(out, "zzcli.drafts") || strings.Contains(out, "idden") {
		t.Errorf("on:\n%s", out)
	}
	b.Reset()
	if err := PrintAll(&b); err != nil {
		t.Fatal(err)
	}
	out := b.String()
	if !strings.Contains(out, "zzcli.drafts") || !strings.Contains(out, "zzcli.forms") {
		t.Errorf("all:\n%s", out)
	}
	if strings.Contains(out, "Hidden") || strings.Contains(out, "zzcli.needshidden") || strings.Contains(out, "compiled") {
		t.Errorf("all shows what is not compiled:\n%s", out)
	}
	features.Set(drafts.Name, true)
	features.Set(forms.Name, false)
	b.Reset()
	_ = PrintAll(&b)
	if !strings.Contains(b.String(), "no (needs zzcli.forms)") {
		t.Errorf("all, a dependency off:\n%s", b.String())
	}
	if handled, _ := RunIfAsked([]string{"app", "version"}); handled {
		t.Error("version handled")
	}
	null, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	defer null.Close()
	stdout := os.Stdout
	os.Stdout = null
	handled, code := RunIfAsked([]string{"app", "features", "all"})
	os.Stdout = stdout
	if !handled || code != 0 {
		t.Errorf("features all: handled=%v code=%d", handled, code)
	}
}
