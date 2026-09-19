package seo

import (
	"encoding/json"
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web"
	. "github.com/go-rvq/rvq/x/ui/vuetify"
)

// varsEditor renders the custom-variables editor for a Setting: the
// `<vx-seo-vars>` component, its value the setting's Vars as a JSON string bound
// to `form["<fieldPrefix>.VarsJSON"]` (parsed back by EditSetterFunc). It is
// given the registered types (so it picks each row's editor component) and the
// Google Maps key (for the built-in zipcodes type). Returns nil when the field
// prefix is empty.
func (b *Builder) varsEditor(fieldPrefix string, setting *Setting) h.HTMLComponent {
	if fieldPrefix == "" {
		return nil
	}
	varsJSON, _ := json.Marshal(setting.Vars)
	if len(setting.Vars) == 0 {
		varsJSON = []byte("[]")
	}
	tag := h.Tag("vx-seo-vars").
		Attr(web.VField(fmt.Sprintf("%s.VarsJSON", fieldPrefix), string(varsJSON))...).
		Attr(":types", b.varTypesJSON())
	if key := b.GoogleMapsAPIKey(); key != "" {
		tag.Attr("maps-key", key)
	}
	return VCard(
		VCardText(tag),
	).Variant(VariantOutlined).Flat(true)
}

// customVarsMenuGroup builds the "+ Variável" menu group for a setting's own
// custom variables (name + description). It is nil when there are none. The
// group is added to the variables menu so a field can insert `{Name}`.
func customVarsMenuGroup(setting *Setting) *VariableGroup {
	if setting == nil || len(setting.Vars) == 0 {
		return nil
	}
	g := &VariableGroup{Name: "CustomVars", Label: "Variáveis", Description: "Variáveis personalizadas deste SEO"}
	for _, v := range setting.Vars {
		if v.Name == "" {
			continue
		}
		g.Items = append(g.Items, VariableItem{Label: v.Name, Expr: v.Name})
	}
	if len(g.Items) == 0 {
		return nil
	}
	return g
}
