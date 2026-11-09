package filecards

import (
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web/vue"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	"github.com/rs/xid"
)

// Drop is a place to drop files on (or to click, to choose them) — a dashed
// area showing placeholder, as the media library's (media.FileDrop) — whose
// files go to model (a JavaScript array: "locals.files", `form["Files"]`) —
// labels, the words of the buttons of a file: edit (an image), remove —,
// sent with the form, not uploaded; under it, the files chosen as cards (as
// Cards'), each removed by its button.
func Drop(model, placeholder string, labels ...string) h.HTMLComponent {
	// the images chosen, a group of the lightbox
	group := "filedrop-" + xid.New().String()
	edit, remove := "Edit", "Remove"
	if len(labels) > 1 {
		edit, remove = labels[0], labels[1]
	}
	// an image its preview (its address in the page: URL by the event's
	// window — a template sees no globals)
	add := fmt.Sprintf(`const URL_ = $event.target.ownerDocument.defaultView.URL; `+
		`files.forEach(f => { if (/^image\/(png|jpeg|gif|webp)$/.test(f.type)) f.preview = URL_.createObjectURL(f) }); `+
		`%[1]s = [...(%[1]s || []), ...files]`, model)
	input := h.Input("").Type("file").Attr("multiple", true).Style("display: none").
		Attr("@change", "const files = [...$event.target.files]; $event.target.value = ''; "+add)
	return vue.UserComponent(
		input,
		h.Div(
			v.VIcon("mdi-paperclip").Class("me-2"),
			h.Div(h.Text(placeholder)).Class("text-body-2"),
		).
			Class("d-flex align-center justify-center pa-3 my-2 text-medium-emphasis rounded text-center").
			Style("white-space: normal").
			Attr(":style", `{border: "2px dashed", borderColor: drop.over ? "rgb(var(--v-theme-primary))" : "rgba(var(--v-border-color), var(--v-border-opacity))", background: drop.over ? "rgba(var(--v-theme-primary), 0.06)" : "", cursor: "pointer", minHeight: "56px"}`).
			Attr("@click", "$event.currentTarget.previousElementSibling.click()").
			Attr("@dragover.prevent", "drop.over = true").
			Attr("@dragleave.prevent", "drop.over = false").
			Attr("@drop.prevent", "drop.over = false; const files = [...$event.dataTransfer.files]; "+add).
			Attr("data-file-drop", true),
		// the files chosen, as cards (Cards): an image its thumbnail, another
		// the icon of its type; each removed by its button
		h.Div(
			h.Div(
				v.VCard(
					// an image: its preview, on it its resolution and size (a chip,
					// below and centered; the resolution read when it loads)
					h.Div(
						// clicked, seen in the lightbox (vuetifyx's) among the
						// images chosen
						h.Img("").Attr(":src", "f.preview").Attr(":alt", "f.name").
							Attr("data-lightbox", group).Attr(":data-lightbox-src", "f.preview").
							Attr(":data-lightbox-caption", "f.name").Attr(":data-lightbox-size", "f.size").
							Attr(":data-lightbox-name", "f.name").
							Attr("@load", "drop.res[drop.key(f)] = $event.target.naturalWidth + '\u00a0×\u00a0' + $event.target.naturalHeight").
							Style("display:block; width:100%; height:96px; object-fit:cover; cursor:zoom-in"),
						v.VChip(h.Text("{{ drop.res[drop.key(f)] ? drop.res[drop.key(f)] + '\u00a0·\u00a0' : '' }}{{ drop.size(f.size) }}")).
							Size(v.SizeXSmall).Variant(v.VariantFlat).Color("surface").
							Attr("data-file-resolution", true).Style(chipStyle),
					).Attr("v-if", "f.preview").Style("position:relative"),
					// another file: the icon of its type, its size on it
					h.Div(
						v.VIcon("").Attr(":icon", "drop.icon(f)").Size(48).Color("primary"),
						v.VChip(h.Text("{{ drop.size(f.size) }}")).Size(v.SizeXSmall).Variant(v.VariantFlat).Color("surface").
							Attr("data-file-size", true).Style(chipStyle),
					).Attr("v-else", true).
						Class("d-flex align-center justify-center bg-surface-light").Style("height:96px; position:relative"),
					h.Div(
						h.Div(h.Text("{{ f.name }}")).Class("text-body-2 text-truncate").Attr(":title", "f.name"),
					).Class("pa-2").Style("min-width:0"),
				).Variant(v.VariantOutlined).Width(160),
				h.Div(
					// an image edited before it is sent (rvq's image editor):
					// the edit in its place
					v.VBtn("").Icon("mdi-pencil").Size(v.SizeXSmall).Variant(v.VariantFlat).
						Attr("v-if", "f.preview").Attr(":title", "drop.editLabel").
						Attr("@click.stop", fmt.Sprintf(`const w = $event.target.ownerDocument.defaultView, ed = w.rvqImageEditor; `+
							`if (ed) ed.editUploads([f]).then(list => { const nf = list && list !== 'retake' && list[0] && list[0].file; `+
							`if (nf && nf !== f) { nf.preview = w.URL.createObjectURL(nf); %s.splice(i, 1, nf) } })`, model)).
						Attr("data-file-edit", true),
					v.VBtn("").Icon("mdi-close").Size(v.SizeXSmall).Variant(v.VariantFlat).
						Attr(":title", "drop.removeLabel").
						Attr("@click.stop", model+".splice(i, 1)").Attr("data-file-remove", true),
				).Class("d-flex ga-1").Style("position:absolute; top:4px; right:4px"),
			).Attr("v-for", "(f, i) in ("+model+" || [])").Attr(":key", "f.name + i").
				Attr("data-file-chosen", true).Style("position:relative; flex:none"),
		).Class("d-flex flex-wrap ga-3"),
		// its state its own (drop): the locals of the page (model) not hidden
	).Scope("drop", vue.Var(fmt.Sprintf(dropState, h.JSONString(edit), h.JSONString(remove))))
}

// dropState is the state of a Drop: over, a file dragged over it; res, the
// resolutions of the images read (by key: a File is not reactive); icon, the
// icon of the type of a file chosen (Icon's); size, its size as people read
// it (FormatSize's).
const dropState = `{over: false, editLabel: %s, removeLabel: %s,
	res: {},
	key(f) { return f.name + ":" + f.size + ":" + f.lastModified + ":" + (f.preview || "") },
	icon(f) {
		const t = f.type || "", e = (f.name.split(".").pop() || "").toLowerCase()
		if (t === "application/pdf" || e === "pdf") return "mdi-file-pdf-box"
		if (t.startsWith("image/")) return "mdi-file-image-outline"
		if (t.startsWith("video/")) return "mdi-file-video-outline"
		if (t.startsWith("audio/")) return "mdi-file-music-outline"
		if (["doc", "docx", "odt", "rtf"].includes(e)) return "mdi-file-word-outline"
		if (["xls", "xlsx", "ods", "csv"].includes(e)) return "mdi-file-excel-outline"
		if (["ppt", "pptx", "odp"].includes(e)) return "mdi-file-powerpoint-outline"
		if (["zip", "rar", "7z", "gz", "tar"].includes(e)) return "mdi-folder-zip-outline"
		if (["txt", "md"].includes(e)) return "mdi-file-document-outline"
		return "mdi-file-outline"
	},
	size(n) { return n >= 1048576 ? (n / 1048576).toFixed(1) + " MB" : n >= 1024 ? Math.round(n / 1024) + " KB" : n + " B" }}`
