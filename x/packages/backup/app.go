package backup

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// AppOptions of AppSource.
type AppOptions struct {
	// Binary is the program ("": the one running, os.Executable).
	Binary string
	// Config are the files of the configuration, of the root — those not
	// there are left out (default: app.yml, app.yaml, .env). The secure
	// key (.secure_key) is never one: the copy makes its own.
	Config []string
}

// AppSource is the source of the application itself: its binary — at its
// path of the root, or "app" when it is out of it — and its configuration.
// Meta: "binary", its path in the state; "config", the files, ",".
func AppSource(o AppOptions) SourceFunc {
	return func(ctx context.Context, st *State) error {
		bin := o.Binary
		if bin == "" {
			var err error
			if bin, err = os.Executable(); err != nil {
				return err
			}
		}
		bin, err := filepath.EvalSymlinks(bin)
		if err != nil {
			return err
		}
		at := "app"
		if rel, err := st.FS.rel(bin); err == nil {
			at = rel
		}
		if err := st.FS.AddFrom(at, bin); err != nil {
			return err
		}
		if err := st.FS.Meta("binary", at); err != nil {
			return err
		}
		config := o.Config
		if config == nil {
			config = []string{"app.yml", "app.yaml", ".env"}
		}
		var have []string
		for _, c := range config {
			if filepath.Base(c) == ".secure_key" {
				continue
			}
			if err := st.FS.Add(c); err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					continue
				}
				return err
			}
			have = append(have, filepath.ToSlash(c))
		}
		return st.FS.Meta("config", strings.Join(have, ","))
	}
}
