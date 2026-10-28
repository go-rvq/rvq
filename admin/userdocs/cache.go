package userdocs

import (
	"sync"
	"time"
)

// parsedFiles keeps, per package in a language, the files by path of its
// value as of its last change (UpdatedAt): a value is YAML of every document
// of the package — pictures in it, up to megabytes —, and a page of the
// documentation reads several of them (its document, the titles of the
// tree). Every write of a value changes its UpdatedAt (the sync's Create and
// Updates, the edit's Save), so only the time is read while it is the same.
type parsedFiles struct {
	mu   sync.Mutex
	byID map[string]*parsedValue
}

type parsedValue struct {
	at    time.Time
	files map[string]string
}

func parsedKey(locale, id string) string { return locale + "\x00" + id }

// get is the files of the package id in locale as of at, when kept.
func (c *parsedFiles) get(locale, id string, at time.Time) (map[string]string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if p := c.byID[parsedKey(locale, id)]; p != nil && p.at.Equal(at) {
		return p.files, true
	}
	return nil, false
}

// parse is the files of value, the value of the package id in locale as of
// at, kept for get.
func (c *parsedFiles) parse(locale, id string, at time.Time, value string) (map[string]string, error) {
	files, err := ParseFiles(value)
	if err != nil {
		return nil, err
	}
	byPath := filesByPath(files)
	c.mu.Lock()
	if c.byID == nil {
		c.byID = map[string]*parsedValue{}
	}
	c.byID[parsedKey(locale, id)] = &parsedValue{at: at, files: byPath}
	c.mu.Unlock()
	return byPath, nil
}
