package basics

import (
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples"
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_presets"
	"github.com/go-rvq/rvq/admin/docs/docsrc/generated"
	"github.com/go-rvq/rvq/admin/docs/docsrc/utils"
	. "github.com/theplant/docgo"
	"github.com/theplant/docgo/ch"
)

var PresetsInstantCRUD = Doc(
	Markdown(`
Presets let you config generalized data management UI interface for database.
It's not a scaffolding to generate source code. But provide more abstract and
flexible API to enrich features along the way.

`),
	ch.Code(generated.PresetHelloWorldSample).Language("go"),
	utils.DemoWithSnippetLocation("Presets Hello World", examples.URLPathByFunc(examples_presets.PresetsHelloWorld)+"/customers", generated.PresetHelloWorldSampleLocation),
	Markdown(`
And this ~*presets.Builder~ instance is actually also a ~http.Handler~, So that we can mount it
to the http serve mux directly
`),
	Markdown(`
With ~b.Model(&Customer{})~:

- It setup the global layout with the left navigation menu
- It setup the listing page with a data table
- It add the new button to create a new record
- It setup the editing and creating form as a right side drawer
- It setup each row of data have a operation menu that you have edit and delete operations
- It setup the global search box, can search the model's all string columns
`),
).Title("presets, Instant CRUD").
	Slug("basics/presets-instant-crud")
