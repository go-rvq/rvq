package perms

import (
	"testing"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"github.com/go-rvq/rvq/x/perm"
)

type pageTestModel struct {
	ID   uint
	Name string
}

// TestPageNodes proves the rule: a page with a custom (developer-defined)
// verifier has its own check mechanism and is NOT enumerated; an auto-perm page
// enters the tree using its path as the permission segment; a page with no perm
// is skipped.
func TestPageNodes(t *testing.T) {
	b := presets.New(i18n.New())
	mb := b.Model(&pageTestModel{})
	mb.Detailing("Name") // enable a detail page context

	reg := mb.Detailing().PagesRegistrator()
	reg.AddHttpPage(presets.HttpPage("/relatorio").AutoPerm())             // enters
	reg.AddHttpPage(presets.HttpPage("/custom").Perm(perm.PermVerifier())) // own mechanism, skipped
	reg.AddHttpPage(presets.HttpPage("/public"))                           // no perm, skipped

	nodes := PageNodes(mb)
	if len(nodes) != 1 {
		t.Fatalf("PageNodes = %+v, want exactly one (the auto-perm page)", nodes)
	}
	n := nodes[0]
	if n.Path != "/relatorio" || n.Mode != ModeDetail {
		t.Errorf("node = %+v, want detail page /relatorio", n)
	}
}
