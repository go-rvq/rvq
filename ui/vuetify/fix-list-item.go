package vuetify

func (b *VListItemBuilder) Slot(v string) (r *VListItemBuilder) {
	b.Attr("slot", v)
	return b
}
