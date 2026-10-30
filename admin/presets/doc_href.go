package presets

// The address of a document of the documentation of whoever uses the admin
// (rvq userdocs, when the application has it): what a field links its help
// to — the guide of its editor —. No documentation, no address.

// DocHrefFunc sets what makes the address of the document of a node
// ("guides/text-editors/01-rich-text"): the documentation's (userdocs
// Configure sets it).
func (b *Builder) DocHrefFunc(f func(node string) string) *Builder {
	b.docHref = f
	return b
}

// DocHref is the address of the document of node; "", no documentation.
func (b *Builder) DocHref(node string) string {
	if b.docHref == nil {
		return ""
	}
	return b.docHref(node)
}
