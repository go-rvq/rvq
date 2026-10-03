package js

import (
	"embed"
	"io/fs"
)

//go:embed gadide/dist
var gadIDE embed.FS

// GadIDE are the files of the app of the gad IDE (gadide/dist, built by
// gadide/build.sh): index.html at the root. The editor of the files of a
// git repository (admin/packages/gitedit) serves it in a frame.
func GadIDE() fs.FS {
	sub, err := fs.Sub(gadIDE, "gadide/dist")
	if err != nil {
		panic(err)
	}
	return sub
}
