package admin

import (
	"context"

	"github.com/go-rvq/rvq/admin/publish"
	"gorm.io/gorm"
)

// installPublishTag marks the record's current revision as the published one —
// the git "tag" — whenever it is published. It appends a publish callback
// (callbacks accumulate, so it does not clobber the app's own), and does nothing
// on models that are not publishable.
func (h *ModelHistory) installPublishTag() {
	publish.WithPublishCallback(h.mb, func(db *gorm.DB, ctx context.Context, obj interface{}) (func(err error) error, error) {
		tag := ""
		if v, ok := obj.(publish.VersionInterface); ok {
			tag = v.EmbedVersion().VersionName
		}
		return nil, h.markPublished(db, h.mb.MustRecordID(obj).String(), tag)
	})
}

// markPublished moves the "published" flag to the record's latest revision and
// stamps its tag. The current published revision of a record is the one with
// published = true (there is at most one per record).
func (h *ModelHistory) markPublished(db *gorm.DB, recordKey, tag string) error {
	if err := db.Table(h.table).
		Where("record_key = ? AND published", recordKey).
		Update("published", false).Error; err != nil {
		return err
	}
	// Use the publish transaction's db; scanLatestHash reads the bytea through
	// Hash's sql.Scanner (a gorm .Scan into a []byte would misread it).
	hash, err := scanLatestHash(db, h.table, recordKey)
	if err != nil {
		return err
	}
	if hash == nil {
		return nil
	}
	return db.Table(h.table).
		Where("record_key = ? AND hash = ?", recordKey, []byte(hash)).
		Updates(map[string]any{"published": true, "tag": tag}).Error
}
