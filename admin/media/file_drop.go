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
	return fileDrop(cfg, placeholder, then, fileDropOptions{multiple: true})
}

// fileDropOptions: click, what a click does (the file input's choice, else);
// multiple, several files at once (else the first only); minHeight, the
// area's.
type fileDropOptions struct {
	click     string
	multiple  bool
	minHeight string
	icon      string
}

func fileDrop(cfg *media_library.MediaBoxConfig, placeholder string, then string, o fileDropOptions) h.HTMLComponent {
	if o.minHeight == "" {
		o.minHeight = "64px"
	}
	if o.icon == "" {
		o.icon = "mdi-cloud-upload-outline"
	}
	if cfg == nil {
		cfg = &media_library.MediaBoxConfig{}
	}
	accept := ""
	if cfg.AllowType == media_library.ALLOW_TYPE_IMAGE {
		accept = "image/*"
	}
	// the images through the editor first when cfg says (EditOnUpload)
	// the spinner from the sending on (an edit cancelled sends nothing)
	upload := uploadThroughEditor(cfg, "files", "locals.uploading = list.length; const fd = {}; "+uploadFormScript(cfg, "fd")+" "+web.Plaid().
		EventFunc(UploadMediaBoxesEvent).
		Form(web.Var("fd")).
		FieldValue("cfg", h.JSONString(cfg)).
		ThenScript("locals.uploading = 0; const boxes = (r && r.data) || []; "+then).
		Go(), false, "")
	first := ""
	if !o.multiple {
		// one file: the first dropped
		first = "files = files.slice(0, 1); "
	}
	send := "((files) => { " + first + "if (!files.length) return; " + upload + " })"
	// statements first: a handler that starts as an arrow function is taken by
	// Vue for the function itself
	input := h.Input("").Type("file").Attr("multiple", o.multiple).Style("display: none").
		Attr("@change", "const files = [...$event.target.files]; $event.target.value = ''; "+send+"(files)")
	if accept != "" {
		input.Attr("accept", accept)
	}
	click := o.click
	if click == "" {
		click = "$event.currentTarget.previousElementSibling.click()"
	}
	return web.Scope(
		input,
		h.Div(
			v.VIcon(o.icon).Class("mb-1"),
			h.Div(h.Text(placeholder)).Attr("v-if", "!locals.uploading").Class("text-body-2"),
			v.VProgressLinear().Indeterminate(true).Attr("v-else", true).Style("max-width: 160px").Color("primary"),
		).
			// its words wrap: a narrow column (a table's) keeps narrow
			Class("d-flex flex-column align-center justify-center pa-3 my-2 text-medium-emphasis rounded text-center").
			Style("white-space: normal; min-width: 140px").
			Attr(":style", `{border: "2px dashed", borderColor: locals.over ? "rgb(var(--v-theme-primary))" : "rgba(var(--v-border-color), var(--v-border-opacity))", background: locals.over ? "rgba(var(--v-theme-primary), 0.06)" : "", cursor: "pointer", minHeight: "`+o.minHeight+`"}`).
			Attr("@click", click).
			Attr("@dragover.prevent", "locals.over = true").
			Attr("@dragleave.prevent", "locals.over = false").
			Attr("@drop.prevent", "locals.over = false; "+send+"([...$event.dataTransfer.files])").
			Attr("data-file-drop", true),
	).Slot("{ locals }").LocalsInit("{over: false, uploading: 0}")
}
