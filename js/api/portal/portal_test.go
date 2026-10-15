package portal

import (
	"context"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
)

// The builder is <rvq-portal> with its attributes as the component reads
// them.
func TestBuilder(t *testing.T) {
	got := h.MustString(Portal().Tag("div").Model("vars.x").Loading("vars.loading").Data("{vars}").
		Name("p").Source("TOKEN").Slot("{show}", h.Text("S")), context.Background())
	for _, want := range []string{"<rvq-portal", `tag='div'`, `v-model='vars.x'`, `v-model:loading='vars.loading'`,
		`:data='{vars}'`, `name='p'`, `source='TOKEN'`, `<template #default='{show}'>S</template>`, "</rvq-portal>"} {
		if !strings.Contains(got, want) {
			t.Errorf("no %s in %s", want, got)
		}
	}
	if got := h.MustString(Portal().DefaultSlot(h.Text("T")), context.Background()); !strings.Contains(got, "<template v-pre>T</template>") {
		t.Errorf("default slot: %s", got)
	}
}
