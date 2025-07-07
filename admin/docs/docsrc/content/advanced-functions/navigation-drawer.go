package advanced_functions

import (
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_vuetify"
	"github.com/go-rvq/rvq/admin/docs/docsrc/generated"
	"github.com/go-rvq/rvq/admin/docs/docsrc/utils"
	. "github.com/theplant/docgo"
	"github.com/theplant/docgo/ch"
)

var NavigationDrawer = Doc(
	Markdown(`
Vuetify navigation drawer provide a popup layer that show on the side of the window.

Here is one example:
`),
	ch.Code(generated.VuetifyNavigationDrawerSample).Language("go"),
	utils.DemoWithSnippetLocation("Vuetify Navigation Drawer", examples_vuetify.VuetifyNavigationDrawerPath, generated.VuetifyNavigationDrawerSampleLocation),
).Title("Navigation Drawer").
	Slug("vuetify-components/navigation-drawer")
