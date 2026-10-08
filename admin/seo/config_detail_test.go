package seo

import (
	"context"
	"strings"
	"testing"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/presets/fields/schemaform"
	"github.com/go-rvq/rvq/web"
)

// The SEOConfig's detail: its settings by their labels, read only (a schema
// form); what the schema does not know, as YAML.
func TestConfigDetail(t *testing.T) {
	schema, err := configSchema(Messages_en_US)
	if err != nil {
		t.Fatal(err)
	}
	c := &SEOConfig{Data: SEOConfigData{SEOConfigGoogleMapsKey: "AIza-test", "other": 1}}
	ctx := &web.EventContext{}
	out, err := h.Marshal(schemaform.New().DetailComponent(ctx, schema, map[string]any(c.Data)), context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if s := string(out); !strings.Contains(s, Messages_en_US.MapsKey) || !strings.Contains(s, "AIza-test") || strings.Contains(s, "<input") {
		t.Errorf("detail:\n%s", s)
	}
	if rest := configDataOthers(c, schema); strings.TrimSpace(rest) != "other: 1" {
		t.Errorf("others: %q", rest)
	}
}
