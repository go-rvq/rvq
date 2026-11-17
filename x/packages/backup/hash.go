package backup

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// CacheKey is what tells a file has not changed since its sum was taken:
// its path, size, modification time and, where the system has them, its
// device and inode (a file replaced by a rename is another).
type CacheKey struct {
	Path       string
	Size       int64
	MTime      int64 // unix ns
	Dev, Inode uint64
}

func keyOf(path string, fi fs.FileInfo) CacheKey {
	dev, ino := devIno(fi)
	return CacheKey{Path: path, Size: fi.Size(), MTime: fi.ModTime().UnixNano(), Dev: dev, Inode: ino}
}

// HashCache keeps the sums of the files: a state does not hash again what
// did not change.
type HashCache interface {
	Lookup(k CacheKey) (Sum, bool)
	Store(k CacheKey, s Sum)
}

// FileCache is a HashCache in a file (gzip, a line a file): loaded by
// OpenFileCache, written by Save with only the sums used since — those of
// files gone are dropped.
type FileCache struct {
	path string
	mu   sync.Mutex
	old  map[CacheKey]Sum
	used map[CacheKey]Sum
}

// OpenFileCache is the cache of path: empty when there is none (or it is
// not readable — a cache is never a reason to fail).
func OpenFileCache(path string) *FileCache {
	c := &FileCache{path: path, old: map[CacheKey]Sum{}, used: map[CacheKey]Sum{}}
	f, err := os.Open(path)
	if err != nil {
		return c
	}
	defer f.Close()
	zr, err := gzip.NewReader(f)
	if err != nil {
		return c
	}
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 64<<10), 1<<20)
	for sc.Scan() {
		if k, s, ok := parseCacheLine(sc.Text()); ok {
			c.old[k] = s
		}
	}
	return c
}

func (c *FileCache) Lookup(k CacheKey) (Sum, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	s, ok := c.used[k]
	if !ok {
		if s, ok = c.old[k]; ok {
			c.used[k] = s
		}
	}
	return s, ok
}

func (c *FileCache) Store(k CacheKey, s Sum) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.used[k] = s
}

// Save writes the sums used (whole or nothing: a file beside, renamed).
func (c *FileCache) Save() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(c.path), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(c.path), ".hash-cache-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	zw := gzip.NewWriter(tmp)
	bw := bufio.NewWriter(zw)
	for k, s := range c.used {
		fmt.Fprintf(bw, "%s\t%d\t%d\t%d\t%d\t%s\n", s, k.Size, k.MTime, k.Dev, k.Inode, strconv.Quote(k.Path))
	}
	err = errors.Join(bw.Flush(), zw.Close(), tmp.Chmod(0o600), tmp.Close())
	if err != nil {
		return err
	}
	return os.Rename(tmp.Name(), c.path)
}

func parseCacheLine(line string) (k CacheKey, s Sum, ok bool) {
	f := strings.SplitN(line, "\t", 6)
	if len(f) != 6 {
		return
	}
	var err error
	if s, err = ParseSum(f[0]); err != nil {
		return
	}
	n := make([]uint64, 4)
	for i := range n {
		if n[i], err = strconv.ParseUint(f[1+i], 10, 64); err != nil {
			return
		}
	}
	if k.Path, err = strconv.Unquote(f[5]); err != nil {
		return
	}
	k.Size, k.MTime, k.Dev, k.Inode = int64(n[0]), int64(n[1]), n[2], n[3]
	return k, s, true
}

// hashFile is the sha256 of the file at path, which must still be as fi
// says at the end (a file written while hashed is an error: the sum would
// be of no version of it).
func hashFile(path string, fi fs.FileInfo) (Sum, error) {
	var s Sum
	f, err := os.Open(path)
	if err != nil {
		return s, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return s, err
	}
	after, err := f.Stat()
	if err != nil {
		return s, err
	}
	if after.Size() != fi.Size() || !after.ModTime().Equal(fi.ModTime()) {
		return s, fmt.Errorf("%s: %w", path, ErrChanged)
	}
	copy(s[:], h.Sum(nil))
	return s, nil
}

// ErrChanged: a file changed while it was read.
var ErrChanged = errors.New("changed while read")

// hasher is a sha256 being written.
type hasher struct {
	h interface {
		io.Writer
		Sum([]byte) []byte
	}
}

func newHasher() *hasher { return &hasher{sha256.New()} }

func (h *hasher) Write(p []byte) (int, error) { return h.h.Write(p) }

func (h *hasher) sum() (s Sum) {
	copy(s[:], h.h.Sum(nil))
	return
}
