package basics

import (
	"github.com/go-rvq/rvq/admin/docs/docsrc/examples/examples_vuetify"
	"github.com/go-rvq/rvq/admin/docs/docsrc/generated"
	"github.com/go-rvq/rvq/admin/docs/docsrc/utils"
	. "github.com/theplant/docgo"
	"github.com/theplant/docgo/ch"
)

var BasicInputs = Doc(
	Markdown(`
Vuetify provides many form basic inputs, and also with error messages display on fields.

Here is one example:
`),
	ch.Code(generated.VuetifyBasicInputsSample).Language("go"),
	utils.DemoWithSnippetLocation("Vuetify Basic Inputs", examples_vuetify.VuetifyBasicInputsPath, generated.VuetifyBasicInputsSampleLocation),
).Title("Basic Inputs").
	Slug("vuetify-components/basic-inputs")
