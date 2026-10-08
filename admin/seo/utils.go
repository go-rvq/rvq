package seo

import (
	"fmt"
	"strconv"

	"github.com/go-rvq/rvq/admin/presets/fields/schemaform"
	"gopkg.in/yaml.v3"
)

// configDataYAML renders the SEO config's Data payload as YAML (empty → "{}").
func configDataYAML(c *SEOConfig) string {
	if len(c.Data) == 0 {
		return "{}"
	}
	b, err := yaml.Marshal(map[string]any(c.Data))
	if err != nil {
		return ""
	}
	return string(b)
}

// configSchema is the schema of the settings of the SEOConfig (its Data), each
// labelled in msgr's words: what its detail shows, read only.
func configSchema(msgr *Messages) (*schemaform.Schema, error) {
	return schemaform.Parse(fmt.Sprintf("interface {\n\t[label=%s, hint=%s]\n\t%s? str\n}",
		strconv.Quote(msgr.MapsKey), strconv.Quote(msgr.MapsKeyHint), SEOConfigGoogleMapsKey))
}

// configDataOthers are the settings of c the schema does not know, as YAML;
// "" when none.
func configDataOthers(c *SEOConfig, schema *schemaform.Schema) string {
	known := map[string]bool{}
	for _, f := range schema.Fields {
		known[f.Name] = true
	}
	rest := map[string]any{}
	for k, v := range c.Data {
		if !known[k] {
			rest[k] = v
		}
	}
	if len(rest) == 0 {
		return ""
	}
	b, err := yaml.Marshal(rest)
	if err != nil {
		return ""
	}
	return string(b)
}
