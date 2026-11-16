package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Dir is a folder of the value Parse reads: where it is shown (Mount, "media"
// or "a/b"), the folder of the OS (Local, absolute) and whether it is only
// read.
type Dir struct {
	Mount, Local string
	ReadOnly     bool
}

// Dirs are the folders of v, the value of Parse.
func Dirs(v string) (out []Dir, err error) {
	if strings.TrimSpace(v) == "" {
		return nil, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for i, part := range strings.Split(v, ";") {
		f := strings.Split(strings.TrimSpace(part), ",")
		var d Dir
		if len(f) == 3 {
			switch f[0] {
			case "r":
				d.ReadOnly = true
			case "w":
			default:
				return nil, fmt.Errorf("Parse FS: item %d (%q): invalid format", i+1, f)
			}
			f = f[1:]
		}
		if len(f) != 2 || f[1] == "" {
			return nil, fmt.Errorf("Parse FS: item %d (%q): invalid format", i+1, f)
		}
		d.Mount, d.Local = strings.Trim(f[0], "/"), f[1]
		if !filepath.IsAbs(d.Local) {
			d.Local = filepath.Join(cwd, d.Local)
		}
		out = append(out, d)
	}
	return
}
