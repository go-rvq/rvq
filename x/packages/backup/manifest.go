package backup

import (
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ManifestVersion is the version of the format of the manifest.
const ManifestVersion = 1

// manifestMagic opens the first line of a manifest; metaMagic, a line of
// the sources' meta.
const (
	manifestMagic = "#backup-manifest"
	metaMagic     = "#meta"
)

// metaKeyRe is a key of the meta: "<source>.<name>".
var metaKeyRe = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

// Type is the kind of an entry: a file, a folder, a link.
type Type byte

const (
	TypeFile Type = 'f'
	TypeDir  Type = 'd'
	TypeLink Type = 'l'
)

// Entry is a line of the manifest: a file, a folder or a link of the
// application's directory, at Path (relative, "/" apart).
type Entry struct {
	Path string
	Type Type
	// Mode is the permission bits and setuid, setgid, sticky (0644, 04755).
	Mode fs.FileMode
	// MTime is the modification time.
	MTime time.Time
	// Size and Sum (sha256) are a file's.
	Size int64
	Sum  Sum
	// Target is a link's.
	Target string
}

// Sum is a sha256.
type Sum [sha256.Size]byte

// String is s in base64url, no padding (43 characters): the manifest's
// and the ETag's form.
func (s Sum) String() string { return base64.RawURLEncoding.EncodeToString(s[:]) }

// IsZero says s is no sum (a folder, a link).
func (s Sum) IsZero() bool { return s == Sum{} }

// SumFromHex reads s, the sha256 in hex (as sha256sum writes it).
func SumFromHex(s string) (sum Sum, err error) {
	b, err := hex.DecodeString(s)
	if err != nil || len(b) != len(sum) {
		return sum, fmt.Errorf("invalid sha256 %q", s)
	}
	copy(sum[:], b)
	return sum, nil
}

// Hex is s in hex.
func (s Sum) Hex() string { return hex.EncodeToString(s[:]) }

// ParseSum reads s, base64url with no padding.
func ParseSum(s string) (sum Sum, err error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil || len(b) != len(sum) {
		return sum, fmt.Errorf("invalid sha256 %q", s)
	}
	copy(sum[:], b)
	return sum, nil
}

// Header is the first line of the manifest: of which state, when, its
// sources and whether some were left out (Partial) — and what the sources
// tell the client (Meta: "site-files.git" the repository, "db.dump" the
// dump), a line each after it.
type Header struct {
	Version int
	ID      string
	Created time.Time
	Sources []string
	Partial bool
	Meta    map[string]string
}

// Manifest is a state's files: its header and its entries, by path.
type Manifest struct {
	Header  Header
	Entries []Entry
}

// Find is the entry of p, by a binary search (Entries sorted).
func (m *Manifest) Find(p string) (Entry, bool) {
	i := sort.Search(len(m.Entries), func(i int) bool { return m.Entries[i].Path >= p })
	if i < len(m.Entries) && m.Entries[i].Path == p {
		return m.Entries[i], true
	}
	return Entry{}, false
}

// The manifest, in gzip, line by line, "\t" apart:
//
//	#backup-manifest <version> <id> <created RFC3339Nano> <sources, ","> <partial 0|1>
//	#meta <key> <value>                (sorted by key; the value quoted as a path)
//	<sha256 | -> <mode octal> <mtime unix ns> <size> <f|d|l> <path> [<link target>]
//
// sorted by path. A path or a target with a tab, a line break, a
// backslash, a control character or starting with a quote is written
// quoted (Go's strconv.Quote).

// WriteManifest writes m (its entries sorted by path) to w, in gzip.
func WriteManifest(w io.Writer, m *Manifest) error {
	sort.Slice(m.Entries, func(i, j int) bool { return m.Entries[i].Path < m.Entries[j].Path })
	zw, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		return err
	}
	bw := bufio.NewWriter(zw)
	h := m.Header
	partial := "0"
	if h.Partial {
		partial = "1"
	}
	fmt.Fprintf(bw, "%s\t%d\t%s\t%s\t%s\t%s\n", manifestMagic, ManifestVersion, h.ID,
		h.Created.UTC().Format(time.RFC3339Nano), strings.Join(h.Sources, ","), partial)
	keys := make([]string, 0, len(h.Meta))
	for k := range h.Meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if !metaKeyRe.MatchString(k) {
			return fmt.Errorf("manifest: invalid meta key %q", k)
		}
		fmt.Fprintf(bw, "%s\t%s\t%s\n", metaMagic, k, quotePath(h.Meta[k]))
	}
	for _, e := range m.Entries {
		sum := "-"
		if e.Type == TypeFile {
			sum = e.Sum.String()
		}
		fmt.Fprintf(bw, "%s\t%o\t%d\t%d\t%c\t%s", sum, unixMode(e.Mode), e.MTime.UnixNano(), e.Size, e.Type, quotePath(e.Path))
		if e.Type == TypeLink {
			bw.WriteString("\t" + quotePath(e.Target))
		}
		bw.WriteByte('\n')
	}
	if err := bw.Flush(); err != nil {
		return err
	}
	return zw.Close()
}

// ReadManifest reads a manifest WriteManifest wrote.
func ReadManifest(r io.Reader) (*Manifest, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("manifest: %w", err)
	}
	defer zr.Close()
	sc := bufio.NewScanner(zr)
	sc.Buffer(make([]byte, 64<<10), 1<<20)
	m := &Manifest{}
	if !sc.Scan() {
		return nil, errors.New("manifest: empty")
	}
	if m.Header, err = parseHeader(sc.Text()); err != nil {
		return nil, err
	}
	for n := 2; sc.Scan(); n++ {
		if line := sc.Text(); strings.HasPrefix(line, metaMagic+"\t") {
			if len(m.Entries) > 0 {
				return nil, fmt.Errorf("manifest: line %d: meta after the entries", n)
			}
			f := strings.SplitN(line, "\t", 3)
			if len(f) != 3 || !metaKeyRe.MatchString(f[1]) {
				return nil, fmt.Errorf("manifest: line %d: meta", n)
			}
			v, err := unquotePath(f[2])
			if err != nil {
				return nil, fmt.Errorf("manifest: line %d: %w", n, err)
			}
			if m.Header.Meta == nil {
				m.Header.Meta = map[string]string{}
			}
			m.Header.Meta[f[1]] = v
			continue
		}
		e, err := parseEntry(sc.Text())
		if err != nil {
			return nil, fmt.Errorf("manifest: line %d: %w", n, err)
		}
		if k := len(m.Entries); k > 0 && m.Entries[k-1].Path >= e.Path {
			return nil, fmt.Errorf("manifest: line %d: %q out of order", n, e.Path)
		}
		m.Entries = append(m.Entries, e)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("manifest: %w", err)
	}
	return m, nil
}

func parseHeader(line string) (h Header, err error) {
	f := strings.Split(line, "\t")
	if len(f) != 6 || f[0] != manifestMagic {
		return h, errors.New("manifest: not a manifest")
	}
	if h.Version, err = strconv.Atoi(f[1]); err != nil || h.Version != ManifestVersion {
		return h, fmt.Errorf("manifest: version %q not known", f[1])
	}
	h.ID = f[2]
	if h.Created, err = time.Parse(time.RFC3339Nano, f[3]); err != nil {
		return h, fmt.Errorf("manifest: created: %w", err)
	}
	if f[4] != "" {
		h.Sources = strings.Split(f[4], ",")
	}
	h.Partial = f[5] == "1"
	return h, nil
}

func parseEntry(line string) (e Entry, err error) {
	f := strings.Split(line, "\t")
	if len(f) < 6 {
		return e, errors.New("too few fields")
	}
	e.Type = Type(f[4][0])
	switch {
	case len(f[4]) != 1:
		return e, fmt.Errorf("type %q", f[4])
	case e.Type == TypeFile:
		if e.Sum, err = ParseSum(f[0]); err != nil {
			return e, err
		}
	case e.Type == TypeDir, e.Type == TypeLink:
		if f[0] != "-" {
			return e, fmt.Errorf("a sum of a %c", e.Type)
		}
	default:
		return e, fmt.Errorf("type %q", f[4])
	}
	mode, err := strconv.ParseUint(f[1], 8, 32)
	if err != nil {
		return e, fmt.Errorf("mode %q", f[1])
	}
	e.Mode = fileMode(uint32(mode))
	ns, err := strconv.ParseInt(f[2], 10, 64)
	if err != nil {
		return e, fmt.Errorf("mtime %q", f[2])
	}
	e.MTime = time.Unix(0, ns)
	if e.Size, err = strconv.ParseInt(f[3], 10, 64); err != nil || e.Size < 0 {
		return e, fmt.Errorf("size %q", f[3])
	}
	if e.Path, err = unquotePath(f[5]); err != nil {
		return e, err
	}
	if err := CheckPath(e.Path); err != nil {
		return e, err
	}
	switch {
	case e.Type == TypeLink && len(f) == 7:
		if e.Target, err = unquotePath(f[6]); err != nil {
			return e, err
		}
	case len(f) != 6:
		return e, errors.New("fields")
	}
	return e, nil
}

// CheckPath says p is a path of the manifest: relative, "/" apart, clean,
// never out of the directory ("..") nor in .backup (the client's).
func CheckPath(p string) error {
	switch {
	case p == "" || p == "." || strings.HasPrefix(p, "/") || strings.HasSuffix(p, "/"):
		return fmt.Errorf("invalid path %q", p)
	case p == ReservedDir || strings.HasPrefix(p, ReservedDir+"/"):
		return fmt.Errorf("path %q in %s", p, ReservedDir)
	}
	for _, s := range strings.Split(p, "/") {
		if s == "" || s == "." || s == ".." {
			return fmt.Errorf("invalid path %q", p)
		}
	}
	return nil
}

// ReservedDir is the client's, in the directory: no entry is there.
const ReservedDir = ".backup"

// quotePath is p as the manifest writes it.
func quotePath(p string) string {
	for _, r := range p {
		if r < 0x20 || r == 0x7f || r == '\\' {
			return strconv.Quote(p)
		}
	}
	if strings.HasPrefix(p, `"`) {
		return strconv.Quote(p)
	}
	return p
}

func unquotePath(s string) (string, error) {
	if strings.HasPrefix(s, `"`) {
		p, err := strconv.Unquote(s)
		if err != nil {
			return "", fmt.Errorf("path %s: %w", s, err)
		}
		return p, nil
	}
	return s, nil
}

// unixMode is m as the octal of Unix: the permissions, 04000 setuid, 02000
// setgid, 01000 sticky.
func unixMode(m fs.FileMode) uint32 {
	u := uint32(m.Perm())
	if m&fs.ModeSetuid != 0 {
		u |= 0o4000
	}
	if m&fs.ModeSetgid != 0 {
		u |= 0o2000
	}
	if m&fs.ModeSticky != 0 {
		u |= 0o1000
	}
	return u
}

// fileMode is unixMode's inverse.
func fileMode(u uint32) fs.FileMode {
	m := fs.FileMode(u & 0o777)
	if u&0o4000 != 0 {
		m |= fs.ModeSetuid
	}
	if u&0o2000 != 0 {
		m |= fs.ModeSetgid
	}
	if u&0o1000 != 0 {
		m |= fs.ModeSticky
	}
	return m
}
