package userdocs

import "sync"

// parsedFiles keeps, per package in a language, the last value read and its
// files by path: a value is YAML of every document of the package — pictures
// in it, up to megabytes —, and a page of the documentation reads several
// of them (its document, the titles of the tree). The value read is compared
// with the one kept, so whatever writes it (an edit, a sync of the boot) is
// seen at once; only an equal value skips the parse.
type parsedFiles struct {
	mu   sync.Mutex
	byID map[string]*parsedValue
}

type parsedValue struct {
	value string
	files map[string]string
}

// of is the files by path of value, the value of the package id in locale.
func (c *parsedFiles) of(locale, id, value string) (map[string]string, error) {
	key := locale + "\x00" + id
	c.mu.Lock()
	if p := c.byID[key]; p != nil && p.value == value {
		c.mu.Unlock()
		return p.files, nil
	}
	c.mu.Unlock()

	files, err := ParseFiles(value)
	if err != nil {
		return nil, err
	}
	byPath := filesByPath(files)
	c.mu.Lock()
	if c.byID == nil {
		c.byID = map[string]*parsedValue{}
	}
	c.byID[key] = &parsedValue{value: value, files: byPath}
	c.mu.Unlock()
	return byPath, nil
}
