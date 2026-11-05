package media

import (
	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/web"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
)

// FileDrop is a place to drop files on (or to click, to choose them): a dashed
// area showing placeholder. The files are uploaded to the media library — as
// cfg says: its type (another type's are left out) and sizes
// (UploadMediaBoxesEvent) — and then runs, with boxes the media boxes of the
// files uploaded, in their order (an array of objects, as a MediaBox's JSON):
//
//	media.FileDrop(cfg, "Drop the files here", "$listEditorAddItems(boxes.map(m => ({Media: m})))")
func FileDrop(cfg *media_library.MediaBoxConfig, placeholder string, then string) h.HTMLComponent {
	if cfg == nil {
		cfg = &media_library.MediaBoxConfig{}
	}
	accept := ""
	if cfg.AllowType == media_library.ALLOW_TYPE_IMAGE {
		accept = "image/*"
	}
	upload := web.Plaid().
		EventFunc(UploadMediaBoxesEvent).
		Form(web.Var("{}")).
		FieldValue("NewFiles", web.Var("files")).
		FieldValue("cfg", h.JSONString(cfg)).
		ThenScript("locals.uploading = 0; const boxes = (r && r.data) || []; " + then).
		Go()
	send := "((files) => { if (!files.length) return; locals.uploading = files.length; " + upload + " })"
	// statements first: a handler that starts as an arrow function is taken by
	// Vue for the function itself
	input := h.Input("").Type("file").Attr("multiple", true).Style("display: none").
		Attr("@change", "const files = [...$event.target.files]; $event.target.value = ''; "+send+"(files)")
	if accept != "" {
		input.Attr("accept", accept)
	}
	return web.Scope(
		input,
		h.Div(
			v.VIcon("mdi-cloud-upload-outline").Class("mb-1"),
			h.Div(h.Text(placeholder)).Attr("v-if", "!locals.uploading").Class("text-body-2"),
			v.VProgressLinear().Indeterminate(true).Attr("v-else", true).Style("max-width: 160px").Color("primary"),
		).
			Class("d-flex flex-column align-center justify-center pa-3 my-2 text-medium-emphasis rounded").
			Attr(":style", `{border: "2px dashed", borderColor: locals.over ? "rgb(var(--v-theme-primary))" : "rgba(var(--v-border-color), var(--v-border-opacity))", background: locals.over ? "rgba(var(--v-theme-primary), 0.06)" : "", cursor: "pointer", minHeight: "64px"}`).
			Attr("@click", "$event.currentTarget.previousElementSibling.click()").
			Attr("@dragover.prevent", "locals.over = true").
			Attr("@dragleave.prevent", "locals.over = false").
			Attr("@drop.prevent", "locals.over = false; "+send+"([...$event.dataTransfer.files])").
			Attr("data-file-drop", true),
	).Slot("{ locals }").LocalsInit("{over: false, uploading: 0}")
}
