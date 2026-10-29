package js

import "embed"

// MdKit is the bundle of mdkit for a page (mdkit/dist: mdkit.js, its
// window.RvqMdKit): the editors of a text and their builders, served by
// x/ui/mdkit's JSComponentsPack.
//
//go:embed mdkit/dist
var MdKit embed.FS
