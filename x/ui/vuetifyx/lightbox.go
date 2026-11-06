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
