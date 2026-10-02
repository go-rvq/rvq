package userdocs

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// UserDoc is the documentation in a language of the admin: a singleton per
// language, its packages' documents in Packages.
type UserDoc struct {
	// LocaleCode is the language (an admin locale).
	LocaleCode string `gorm:"primaryKey"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Packages []*UserDocPackage `gorm:"foreignKey:LocaleCode;references:LocaleCode;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// UserDocPackage is the documentation of a package in a language: its
// documents, as YAML — a list of {path, content}, the paths those of the
// source's files.
type UserDocPackage struct {
	LocaleCode string `admin:"-" gorm:"primaryKey"`
	// ID is the Go package.
	ID        string `admin:"readonly" gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// SourceLocale is the language the package's documents are written in.
	SourceLocale string `admin:"readonly" gorm:"not null;default:''"`
	// Value is the documents shown, and edited (YAML).
	Value string `gorm:"type:text"`
	// InitialValue is the code's: written on every boot from the source's
	// files — translated, in another language —; a reset brings it back.
	InitialValue string `admin:"readonly" gorm:"type:text"`
	// SourceHashes are the hashes of the source's files the initial value was
	// made from, by path: a file whose hash changed is translated again.
	SourceHashes Hashes `admin:"-" gorm:"type:jsonb"`
	// Unused says the code registers the package no more.
	Unused bool `admin:"readonly" gorm:"not null;default:false"`
}

// PrimarySlug is the key of the package as URLs carry it: the language and
// the package, "pt-BR_github.com/go-rvq/rvq/admin/packages/user"
// (path-escaped by the presets).
func (p *UserDocPackage) PrimarySlug() string {
	return p.LocaleCode + "_" + p.ID
}

// PrimaryColumnValuesBySlug is the key of PrimarySlug: the language up to the
// first "_" (a BCP-47 code has none), the package after it.
func (p *UserDocPackage) PrimaryColumnValuesBySlug(slug string) map[string]string {
	lang, pkg, _ := strings.Cut(slug, "_")
	return map[string]string{"LocaleCode": lang, "ID": pkg}
}

func (p *UserDocPackage) String() string {
	return p.ID
}

// Hashes are hashes by path, stored as a JSON object — NULL for none.
type Hashes map[string]string

// Value stores the hashes as JSON, and none as NULL.
func (h Hashes) Value() (driver.Value, error) {
	if len(h) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(map[string]string(h))
	return string(b), err
}

// Scan reads the hashes back from JSON.
func (h *Hashes) Scan(v any) error {
	var b []byte
	switch t := v.(type) {
	case nil:
		*h = nil
		return nil
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		return fmt.Errorf("userdocs.Hashes: cannot scan %T", v)
	}
	out := map[string]string{}
	if err := json.Unmarshal(b, &out); err != nil {
		return err
	}
	*h = out
	return nil
}

// hash is the hash of a source's file.
func hash(content string) string {
	s := sha256.Sum256([]byte(content))
	return hex.EncodeToString(s[:])
}

// ParseFiles reads documents written as YAML (UserDocPackage.Value).
func ParseFiles(value string) ([]File, error) {
	var files []File
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	if err := yaml.Unmarshal([]byte(value), &files); err != nil {
		return nil, fmt.Errorf("userdocs: %w", err)
	}
	return files, nil
}

// FormatFiles writes documents as YAML (UserDocPackage.Value).
func FormatFiles(files []File) (string, error) {
	if len(files) == 0 {
		return "", nil
	}
	b, err := yaml.Marshal(files)
	return string(b), err
}

// filesByPath is files by their path.
func filesByPath(files []File) map[string]string {
	m := make(map[string]string, len(files))
	for _, f := range files {
		m[f.Path] = f.Content
	}
	return m
}
