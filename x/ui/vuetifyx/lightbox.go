package vuetifyx

// Lightbox are the attributes that make an element open src in the lightbox
// (vuetifyx's, over the page) when clicked: the elements of the same group —
// the images of an item, say — in their order on the page, gone through by
// the arrows; caption is said under the image. group "", the image alone.
//
//	h.Div(img).Attr(vuetifyx.Lightbox("item-1", url, name)...)
func Lightbox(group, src, caption string) []any {
	attrs := []any{"data-lightbox-src", src, "style", "cursor: zoom-in"}
	if group != "" {
		attrs = append(attrs, "data-lightbox", group)
	}
	if caption != "" {
		attrs = append(attrs, "data-lightbox-caption", caption)
	}
	return attrs
}

// LightboxFile is Lightbox of a file: its size (bytes, > 0) and the name it
// is downloaded as, said in the chip of the lightbox (with its dimensions
// and its download).
func LightboxFile(group, src, caption string, size int64, name string) []any {
	attrs := Lightbox(group, src, caption)
	if size > 0 {
		attrs = append(attrs, "data-lightbox-size", size)
	}
	if name != "" {
		attrs = append(attrs, "data-lightbox-name", name)
	}
	return attrs
}

// LightboxDescription are the attributes that say description below the
// image in the lightbox (hidden by its button); "" none.
func LightboxDescription(description string) []any {
	if description == "" {
		return nil
	}
	return []any{"data-lightbox-description", description}
}
