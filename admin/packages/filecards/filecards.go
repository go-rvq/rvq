// Package filecards shows files as cards: an image, its thumbnail — opened
// in the lightbox (vuetifyx's) among the images of the cards —;
// another file, the icon of its type — opened in a tab —; each with its name,
// its size and the button that downloads it. The attachments of a comment,
// the files of an answer to a form.
package filecards

import (
	"fmt"
	"path/filepath"
	"strings"

	h "github.com/go-rvq/htmlgo"
	v "github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
	"github.com/rs/xid"
)

// File is a file of the cards.
type File struct {
	// Key identifies it in the page (the value of its attribute, Attr)
	Key  string
	Name string
	Size int64
	// Type is its media type (image/png…): an image of ImageTypes is seen
	Type string
	URL  string
}

// ImageTypes are the images a card shows (the browser shows them, and they
// run nothing: no SVG).
var ImageTypes = map[string]bool{"image/png": true, "image/jpeg": true, "image/gif": true, "image/webp": true}

// Attr is the attribute of each card, its value the file's Key
// ("data-file"; an app's own: "data-answer-file").
type Attr string

// Cards are the files as cards (nil: none).
func Cards(files []File, attr ...Attr) h.HTMLComponent {
	if len(files) == 0 {
		return nil
	}
	a := "data-file"
	if len(attr) > 0 && attr[0] != "" {
		a = string(attr[0])
	}
	// the images of these cards, a group of the lightbox (vuetifyx's): the
	// previous, the next of these only
	group := "filecards-" + xid.New().String()
	var cards h.HTMLComponents
	for _, f := range files {
		image := ImageTypes[strings.ToLower(f.Type)]
		label := h.Div(
			h.Div(h.Text(f.Name)).Class("text-body-2 text-truncate").Attr("title", f.Name),
		).Class("pa-2").Style("min-width:0")
		var card *v.VCardBuilder
		if image {
			// on the thumbnail its resolution and size (a chip, below and
			// centered; the resolution read when it loads)
			card = v.VCard(
				h.Div(
					h.Img(f.URL).Attr("alt", f.Name).Attr("loading", "lazy").Attr(vx.LightboxFile(group, f.URL, f.Name, f.Size, f.Name)...).
						Attr("@load", `$event.target.nextElementSibling.querySelector("[data-file-resolution]").textContent = `+
							`$event.target.naturalWidth + "\u00a0×\u00a0" + $event.target.naturalHeight + "\u00a0·\u00a0"`).
						Style("display:block; width:100%; height:96px; object-fit:cover; cursor:zoom-in"),
					v.VChip(h.Span("").Attr("data-file-resolution", true), h.Text(FormatSize(f.Size))).
						Size(v.SizeXSmall).Variant(v.VariantFlat).Color("surface").Style(chipStyle),
				).Style("position:relative"),
				label,
			)
		} else {
			card = v.VCard(
				// its size on it, as an image's (a chip, below and centered)
				h.Div(
					v.VIcon(Icon(f.Name, f.Type)).Size(48).Color("primary"),
					v.VChip(h.Text(FormatSize(f.Size))).Size(v.SizeXSmall).Variant(v.VariantFlat).Color("surface").
						Attr("data-file-size", true).Style(chipStyle),
				).Class("d-flex align-center justify-center bg-surface-light").Style("height:96px; position:relative"),
				label,
			).Href(f.URL).Attr("target", "_blank").Attr("rel", "noopener")
		}
		cards = append(cards, h.Div(
			card.Variant(v.VariantOutlined).Width(160).Attr(a, f.Key).Attr("data-file-card", true),
			download(f),
		).Style("position:relative; flex:none"))
	}
	return h.Div(cards...).Class("d-flex flex-wrap ga-3")
}

// chipStyle puts the chip of the resolution and the size on the thumbnail:
// below, centered.
const chipStyle = "position:absolute; bottom:4px; left:50%; transform:translateX(-50%); opacity:0.9; white-space:nowrap"

// download is the button that downloads f, at the corner of its card.
func download(f File) h.HTMLComponent {
	return v.VBtn("").Icon("mdi-download").Size(v.SizeXSmall).Variant(v.VariantFlat).
		Href(f.URL).Attr("download", f.Name).Attr("title", f.Name).Attr("@click.stop", "").
		Attr("data-file-download", true).
		Style("position:absolute; top:4px; right:4px")
}

// Icon is the icon of a file of name and type.
func Icon(name, typ string) string {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
	switch {
	case typ == "application/pdf" || ext == "pdf":
		return "mdi-file-pdf-box"
	case strings.HasPrefix(typ, "image/"):
		return "mdi-file-image-outline"
	case strings.HasPrefix(typ, "video/"):
		return "mdi-file-video-outline"
	case strings.HasPrefix(typ, "audio/"):
		return "mdi-file-music-outline"
	}
	switch ext {
	case "doc", "docx", "odt", "rtf":
		return "mdi-file-word-outline"
	case "xls", "xlsx", "ods", "csv":
		return "mdi-file-excel-outline"
	case "ppt", "pptx", "odp":
		return "mdi-file-powerpoint-outline"
	case "zip", "rar", "7z", "gz", "tar":
		return "mdi-folder-zip-outline"
	case "txt", "md":
		return "mdi-file-document-outline"
	}
	return "mdi-file-outline"
}

// FormatSize is n bytes as people read it.
func FormatSize(n int64) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.0f KB", float64(n)/(1<<10))
	}
	return fmt.Sprintf("%d B", n)
}
