package vuetify

func (b *VToolbarBuilder) AutoHeight(v bool) (r *VToolbarBuilder) {
	if v {
		b.Attr(":height", `"auto"`)
	} else {
		b.Attr(":height", ``)
	}
	return b
}
