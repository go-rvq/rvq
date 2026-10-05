package js

import (
	"embed"
	"io/fs"
)

//go:embed gadide/dist
var gadIDE embed.FS

// GadIDE are the files of the gad IDE (gadide/dist, built by gadide/build.sh):
// gadide.js and gadide.css, the component the admin's <vx-gad-ide> loads and
// renders with the page's Vue and Vuetify. The editor of the files of a git
// repository (admin/packages/gitedit) serves them.
func GadIDE() fs.FS {
	sub, err := fs.Sub(gadIDE, "gadide/dist")
	if err != nil {
		panic(err)
	}
	return sub
}
