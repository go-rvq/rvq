package basics

import (
	. "github.com/go-rvq/rvq/admin/docs/docgo"
	"github.com/go-rvq/rvq/admin/docs/docgo/ch"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_presets"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_vuetifyx"
	"github.com/go-rvq/rvq/admin/docs/docsrc/generated"
	"github.com/go-rvq/rvq/admin/docs/docsrc/utils"
)

var LinkageSelect = Doc(
	Markdown(`
LinkageSelect is a component for multi-level linkage select.
    `),
	ch.Code(generated.VuetifyComponentsLinkageSelect).Language("go"),
	utils.DemoWithSnippetLocation("Vuetify LinkageSelect", examples_vuetifyx.VuetifyComponentsLinkageSelectPath, generated.VuetifyComponentsLinkageSelectLocation),
	Markdown(`
### Filter intergation
    `),
	ch.Code(generated.LinkageSelectFilterItem).Language("go"),
	utils.DemoWithSnippetLocation("LinkageSelect Filter Item", examples.URLPathByFunc(examples_presets.PresetsLinkageSelectFilterItem)+"/addresses", generated.LinkageSelectFilterItemLocation),
).Title("Linkage Select").
	Slug("vuetify-components/linkage-select")
