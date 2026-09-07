package filesystem

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/go-rvq/rvq/admin/media/base"
)

// An empty url resolves to the storage directory, and os.Open succeeds on a
// directory — which is how cropping a media with no file ended up reporting
// "copy_file_range: is a directory" against the destination it was writing.
func TestRetrieveEmptyURLOpensADirectory(t *testing.T) {
	dir := t.TempDir()
	opt := &base.Option{}
	opt.Set("path", dir)

	f := FileSystem{}
	full, err := f.GetFullPath("", opt)
	if err != nil {
		t.Fatal(err)
	}
	if full != filepath.Clean(dir) {
		t.Fatalf("empty url resolved to %q, want the storage directory %q", full, dir)
	}

	src, err := os.Open(full)
	if err != nil {
		t.Fatalf("opening the directory failed, so the trap is gone: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join(dir, "file."))
	if err != nil {
		t.Fatal(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err == nil {
		t.Fatal("copying from a directory succeeded")
	}
	if !errors.Is(err, syscall.EISDIR) {
		t.Fatalf("copy failed with %v, want EISDIR", err)
	}
}
