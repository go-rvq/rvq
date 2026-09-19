package basics

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_presets"
	"github.com/go-rvq/rvq/admin/docs/docsrc/generated"
	"github.com/go-rvq/rvq/admin/docs/docsrc/utils"
	. "github.com/theplant/docgo"
	"github.com/theplant/docgo/ch"
)

var ManageMenu = Doc(
	Markdown(`
Menu refers to the list on the left side of the page, such as the menu of the Demo below contains Customers and Companies.
`),
	h.Br(),
	utils.Demo("", examples_presets.PresetsDetailPageCardsPath+"/customers", ""),
	Markdown(`
## Menu order
The menu is a tree, and every model, page and group in it is one entry with a
key of its own. ~MenuOrder~ puts entries at the root, in the order given.

Name an entry with ~presets.ModelItem~, ~presets.PageItem~ or
~presets.GroupItem~ — the type is half of the key, so a model, a page and a
group may share a name without colliding. A model's name is its registration id
(the snake_case of its plural label, or what ~presets.ModelWithID~ set), NOT its
URI name.

The entry need not exist yet: naming it puts it in place and it waits there for
its registration, so the order in which models, pages and groups are configured
does not matter. A bare string still works — a page when it starts with "/", a
model otherwise — but it cannot say *the group named x*.
`),
	ch.Code(generated.MenuOrderSample).Language("go"),
	utils.DemoWithSnippetLocation("Menu Order", examples.URLPathByFunc(examples_presets.PresetsOrderMenu)+"/books", generated.MenuOrderSampleLocation),
	Markdown(`
## Menu group and icon
~MenuGroup~ finds or creates a group, and ~Add~ puts entries inside it. Groups
nest: a group may be added to another group.

An entry is in exactly one place at a time — adding it elsewhere moves it,
never duplicates it — and the last call wins.

A group is also a URL segment, and an entry carries the whole chain of groups it
sits under: a model in ~site~ > ~seo~ is served at ~/admin/site/seo/<uri>~ and
its permission reads ~…:site:seo:<uri>~.

Use ~MenuIcon~ on ~ModelBuilder~ to set the item icon, and ~Icon~ on the group
for the group's own. Icon strings can be found at <https://fonts.google.com/icons>.

The full reference lives in the presets package, at ~admin/presets/docs/menu.md~.
`),
	ch.Code(generated.MenuGroupSample).Language("go"),
	utils.DemoWithSnippetLocation("Menu Group", examples.URLPathByFunc(examples_presets.PresetsGroupMenu)+"/videos", generated.MenuGroupSampleLocation),
).Title("Menu").
	Slug("basics/menu")
