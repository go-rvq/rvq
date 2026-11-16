package uploadparts

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const id = "sending-0001"

func shaOf(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// Parts out of order, one sent twice: the file, its size and hash checked.
func TestPartsMakeTheFile(t *testing.T) {
	b := New(t.TempDir())
	content := "abcdefghij"
	for _, p := range []struct {
		off  int64
		text string
	}{{5, "fghij"}, {0, "abcde"}, {5, "fghij"}} {
		if err := b.WritePart(id, p.off, int64(len(content)), strings.NewReader(p.text)); err != nil {
			t.Fatal(err)
		}
	}
	tok, err := b.Done(id, int64(len(content)), "../dir/name.txt", "not trusted")
	if err != nil {
		t.Fatal(err)
	}
	if tok.Size != 10 || tok.Hash != shaOf(content) || tok.Name != "name.txt" {
		t.Errorf("the token: %+v", tok)
	}
	f, err := b.Open(id)
	if err != nil {
		t.Fatal(err)
	}
	got, _ := io.ReadAll(f)
	f.Close()
	if string(got) != content {
		t.Errorf("the file: %q", got)
	}
	b.Remove(id)
	if _, err := b.Open(id); !errors.Is(err, ErrNotDone) {
		t.Errorf("removed: %v", err)
	}
}

// What is refused: a part past the size, a size other than the parts', an
// id that is a path; an empty file is one.
func TestPartsRefused(t *testing.T) {
	b := New(t.TempDir())
	if err := b.WritePart(id, 3, 5, strings.NewReader("abc")); !errors.Is(err, ErrPart) {
		t.Errorf("past the size: %v", err)
	}
	b2 := New(t.TempDir())
	b2.WritePart(id, 0, 5, strings.NewReader("abc"))
	if _, err := b2.Done(id, 5, "x", ""); !errors.Is(err, ErrSize) {
		t.Errorf("the size: %v", err)
	}
	if err := b.WritePart("../../etc/passwd", 0, 1, strings.NewReader("x")); !errors.Is(err, ErrID) {
		t.Errorf("an id of a path: %v", err)
	}
	if tok, err := New(t.TempDir()).Done("empty-0001", 0, "e", ""); err != nil || tok.Hash != shaOf("") {
		t.Errorf("an empty file: %+v %v", tok, err)
	}
}

// A sending left for MaxAge goes, at the next part.
func TestPartsLeftAreRemoved(t *testing.T) {
	dir := t.TempDir()
	b := New(dir)
	b.WritePart("old-sending", 0, 3, strings.NewReader("abc"))
	old := time.Now().Add(-MaxAge - time.Hour)
	os.Chtimes(filepath.Join(dir, "old-sending"+partsSuffix), old, old)
	b.WritePart(id, 0, 3, strings.NewReader("abc"))
	if _, err := os.Stat(filepath.Join(dir, "old-sending"+partsSuffix)); !os.IsNotExist(err) {
		t.Errorf("the sending left is there: %v", err)
	}
}
