package seo

import (
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
