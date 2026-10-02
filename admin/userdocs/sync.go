package userdocs

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/text/language"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Translator translates the documents into the languages that are not their
// source's.
type Translator interface {
	// Enabled says the translator may be used now.
	Enabled() bool
	// Code is the translator's code of a language (a locale).
	Code(locale string) (string, error)
	// Translate is texts from src into dst (codes of Code).
	Translate(src, dst string, texts []string) ([]string, error)
}

// Job is what Sync left to translate: the documents of a package, in a
// language, from its source's.
type Job struct {
	Locale, Package, SourceLocale string
	// Files are the source's files to translate.
	Files []File
}

// Migrate makes the tables of the documentation.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&UserDoc{}, &UserDocPackage{})
}

// Sync writes the documents of every source in every language of locales:
//
//   - the initial value is the documents of the language, where the source
//     has them; else the source language's, translated (a Job: the source
//     language's, till done);
//   - the value is the initial one for a new package, and for a document the
//     user has not changed; a document changed keeps the change;
//   - a file whose source did not change is not translated again;
//   - a package the code has no more is marked Unused.
func Sync(db *gorm.DB, sources []*Source, locales []string) (jobs []*Job, err error) {
	if err = Migrate(db); err != nil {
		return nil, err
	}
	db = db.Session(&gorm.Session{})
	for _, locale := range locales {
		if err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(&UserDoc{LocaleCode: locale}).Error; err != nil {
			return nil, err
		}
		var rows []*UserDocPackage
		if err = db.Where("locale_code = ?", locale).Find(&rows).Error; err != nil {
			return nil, err
		}
		have := map[string]*UserDocPackage{}
		for _, r := range rows {
			have[r.ID] = r
		}
		seen := map[string]bool{}
		for _, src := range sources {
			seen[src.Package] = true
			job, err := syncPackage(db, src, locale, have[src.Package])
			if err != nil {
				return nil, err
			}
			if job != nil {
				jobs = append(jobs, job)
			}
		}
		for _, r := range rows {
			if !seen[r.ID] && !r.Unused {
				if err = db.Model(r).UpdateColumn("unused", true).Error; err != nil {
					return nil, err
				}
			}
		}
	}
	return jobs, nil
}

// syncPackage writes the documents of src in locale over row (nil: new).
func syncPackage(db *gorm.DB, src *Source, locale string, row *UserDocPackage) (*Job, error) {
	lang, err := src.Lang()
	if err != nil {
		return nil, err
	}
	files, err := src.Files()
	if err != nil {
		return nil, err
	}
	same := sameLanguage(lang, locale)
	// the documents written in locale: as they are, not translated
	own := map[string]string{}
	if l := src.LangOf(locale); l != "" && !same {
		written, err := src.FilesOf(l)
		if err != nil {
			return nil, err
		}
		own = filesByPath(written)
		for _, f := range written {
			if !slices.ContainsFunc(files, func(s File) bool { return s.Path == f.Path }) {
				files = append(files, File{Path: f.Path, Content: f.Content})
			}
		}
	}

	var oldInitial, oldValue map[string]string
	hashes := Hashes{}
	if row != nil {
		ini, err := ParseFiles(row.InitialValue)
		if err != nil {
			return nil, err
		}
		val, err := ParseFiles(row.Value)
		if err != nil {
			return nil, err
		}
		oldInitial, oldValue = filesByPath(ini), filesByPath(val)
		for k, v := range row.SourceHashes {
			hashes[k] = v
		}
	}

	job := &Job{Locale: locale, Package: src.Package, SourceLocale: lang}
	initial := make([]File, len(files))
	value := make([]File, len(files))
	for i, f := range files {
		h := hash(f.Content)
		content, translated := f.Content, true
		if c, ok := own[f.Path]; ok {
			content = c
		} else if !same {
			if prev, ok := oldInitial[f.Path]; ok && hashes[f.Path] == h {
				content = prev
			} else {
				translated = false
				job.Files = append(job.Files, f)
			}
		}
		if translated {
			hashes[f.Path] = h
		}
		initial[i] = File{Path: f.Path, Content: content}
		// a document the user changed keeps the change
		v := content
		if cur, ok := oldValue[f.Path]; ok && cur != oldInitial[f.Path] {
			v = cur
		}
		value[i] = File{Path: f.Path, Content: v}
	}
	for p := range hashes {
		if !slices.ContainsFunc(files, func(f File) bool { return f.Path == p }) {
			delete(hashes, p)
		}
	}

	ini, err := FormatFiles(initial)
	if err != nil {
		return nil, err
	}
	val, err := FormatFiles(value)
	if err != nil {
		return nil, err
	}
	next := &UserDocPackage{LocaleCode: locale, ID: src.Package, SourceLocale: lang,
		Value: val, InitialValue: ini, SourceHashes: hashes}
	if row == nil {
		err = db.Create(next).Error
	} else if row.Value != val || row.InitialValue != ini || row.SourceLocale != lang || row.Unused ||
		!equalHashes(row.SourceHashes, hashes) {
		err = db.Model(row).Select("value", "initial_value", "source_locale", "source_hashes", "unused").
			Updates(&UserDocPackage{Value: val, InitialValue: ini, SourceLocale: lang, SourceHashes: hashes}).Error
	}
	if err != nil {
		return nil, fmt.Errorf("userdocs %s %s: %w", locale, src.Package, err)
	}
	if len(job.Files) == 0 {
		return nil, nil
	}
	return job, nil
}

func equalHashes(a, b Hashes) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// sameLanguage says the locales are of the same language (en, en-US).
func sameLanguage(a, b string) bool {
	ba, _ := language.Make(a).Base()
	bb, _ := language.Make(b).Base()
	return ba == bb
}

// TranslateBatch is how many documents go to the translator at once.
const TranslateBatch = 20

// Translate does a Job: the documents translated from the source's language,
// written as the initial value — and as the value, where the user had not
// changed it —, their source's hashes with them.
func Translate(db *gorm.DB, tr Translator, job *Job) error {
	if tr == nil || !tr.Enabled() {
		return nil
	}
	src, err := tr.Code(job.SourceLocale)
	if err != nil {
		return err
	}
	dst, err := tr.Code(job.Locale)
	if err != nil {
		return err
	}
	out := make([]string, 0, len(job.Files))
	for i := 0; i < len(job.Files); i += TranslateBatch {
		batch := job.Files[i:min(i+TranslateBatch, len(job.Files))]
		texts := make([]string, len(batch))
		masks := make([][]string, len(batch))
		for j, f := range batch {
			texts[j], masks[j] = Mask(f.Content)
		}
		res, err := tr.Translate(src, dst, texts)
		if err != nil {
			return err
		}
		if len(res) != len(batch) {
			return fmt.Errorf("userdocs: translated %d documents into %d", len(batch), len(res))
		}
		for j := range res {
			out = append(out, Unmask(res[j], masks[j]))
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var row UserDocPackage
		if err := tx.Take(&row, "locale_code = ? AND id = ?", job.Locale, job.Package).Error; err != nil {
			return err
		}
		ini, err := ParseFiles(row.InitialValue)
		if err != nil {
			return err
		}
		val, err := ParseFiles(row.Value)
		if err != nil {
			return err
		}
		hashes := Hashes{}
		for k, v := range row.SourceHashes {
			hashes[k] = v
		}
		for i, f := range job.Files {
			for j := range ini {
				if ini[j].Path != f.Path {
					continue
				}
				// the value follows the initial one where the user did not
				// change it
				for k := range val {
					if val[k].Path == f.Path && val[k].Content == ini[j].Content {
						val[k].Content = out[i]
					}
				}
				ini[j].Content = out[i]
			}
			hashes[f.Path] = hash(f.Content)
		}
		iniS, err := FormatFiles(ini)
		if err != nil {
			return err
		}
		valS, err := FormatFiles(val)
		if err != nil {
			return err
		}
		return tx.Model(&row).Select("value", "initial_value", "source_hashes").
			Updates(&UserDocPackage{Value: valS, InitialValue: iniS, SourceHashes: hashes}).Error
	})
}

// maskRe is what a translator must not touch in a document: fenced code,
// inline code, template tags, the target of a link or an image, an HTML tag,
// and the markers that open a line (a heading, a list item, a quote, a table).
var maskRe = regexp.MustCompile("(?s)```.*?```|`[^`\\n]*`|\\{%.*?%\\}|\\]\\([^)\\n]*\\)|<[^>\\n]+>|(?m)^[ \\t]*(?:#{1,6}|[-*+]|\\d+\\.|>|\\|)[ \\t]")

// Mask is text with what a translator must not touch (maskRe) as tokens —
// ⟦0⟧, ⟦1⟧… —, and what they stand for.
func Mask(text string) (string, []string) {
	var masks []string
	out := maskRe.ReplaceAllStringFunc(text, func(m string) string {
		masks = append(masks, m)
		return "⟦" + strconv.Itoa(len(masks)-1) + "⟧"
	})
	return out, masks
}

// Unmask puts back what Mask took out.
func Unmask(text string, masks []string) string {
	for i := len(masks) - 1; i >= 0; i-- {
		text = strings.ReplaceAll(text, "⟦"+strconv.Itoa(i)+"⟧", masks[i])
	}
	return text
}
