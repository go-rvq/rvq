package fs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDirs(t *testing.T) {
	cwd, _ := os.Getwd()
	ds, err := Dirs("r,/media,data/media; w,/x/y,/abs")
	if err != nil || len(ds) != 2 {
		t.Fatal(ds, err)
	}
	if ds[0] != (Dir{Mount: "media", Local: filepath.Join(cwd, "data/media"), ReadOnly: true}) || ds[1] != (Dir{Mount: "x/y", Local: "/abs"}) {
		t.Error(ds)
	}
	if _, err := Dirs("z,/a,b"); err == nil {
		t.Error("z accepted")
	}
}

func TestEmptyFS(t *testing.T) {
	f, err := EmptyFS().Open(".")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if fi, err := f.Stat(); err != nil || !fi.IsDir() {
		t.Error(fi, err)
	}
}
