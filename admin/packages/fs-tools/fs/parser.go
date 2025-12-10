package fs

import (
	"fmt"
	gos "os"
	"path/filepath"
	"strings"

	"github.com/hack-pad/hackpadfs"
	"github.com/hack-pad/hackpadfs/mem"
	"github.com/hack-pad/hackpadfs/os"

	"github.com/hack-pad/hackpadfs/mount"
)

// Parse Parse string of os directories map to MergedFS.
// Example:
//   - `w,/backup,_backup;r,/data,data`
//     two dirs:
//   - w,/backup,_backup
//     writable, mounted at /backup, local dir ./_backup
//   - r,/data,data
//     readonly, mounted at /data, local dir is ./data
func Parse(v string) (_ hackpadfs.FS, err error) {
	var (
		parts    = strings.Split(v, ";")
		existsMp = map[string]bool{}
		existsLp = map[string]bool{}
		osfs     = os.NewFS()
		memFS    *mem.FS
		mfs      *mount.FS
		cwd      string
	)

	if cwd, err = gos.Getwd(); err != nil {
		return
	}

	if memFS, err = mem.NewFS(); err != nil {
		return
	}

	if mfs, err = mount.NewFS(NewROfs(memFS)); err != nil {
		return
	}

	for i, part := range parts {
		part = strings.TrimSpace(part)
		v := strings.Split(part, ",")

		var ro bool
		if len(v) == 3 {
			switch v[0] {
			case "r":
				ro = true
			case "w":
			default:
				return nil, fmt.Errorf("Parse FS: item %d (%q): invalid format", i+1, v)
			}
			ro = v[0] == "r"
			v = v[1:]
		}

		if len(v) == 2 {
			v[0] = strings.Trim(v[0], "/")

			if existsMp[v[0]] {
				return nil, fmt.Errorf("Parse FS: item %d (%q): duplicated mount path", i+1, v)
			}

			if existsLp[v[1]] {
				return nil, fmt.Errorf("Parse FS: item %d (%q): duplicated local path", i+1, v)
			}

			existsMp[v[0]] = true
			existsLp[v[1]] = true

			var dfs hackpadfs.FS
			sub := v[1]
			if sub[0] != '/' {
				sub = filepath.Join(cwd, sub)
			}
			if dfs, err = osfs.Sub(strings.Trim(sub, "/")); err != nil {
				return nil, err
			}

			if ro {
				dfs = NewROfs(dfs.(*os.FS))
			}

			if err = memFS.MkdirAll(v[0], 0755); err != nil {
				return nil, err
			}

			if err = mfs.AddMount(v[0], dfs); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("Parse FS: item %d (%q): invalid format", i+1, v)
		}
	}

	return mfs, nil
}
