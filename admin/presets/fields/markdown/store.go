package markdown

import (
	"context"

	"github.com/go-rvq/rvq/admin/media"
	"github.com/go-rvq/rvq/admin/presets/fields/tiptap"
	"gorm.io/gorm"
)

type dbKey struct{}

// MediaLibraryStorer stores the pictures in the media library of b (hidden
// images, as the TipTap fields'), in db.
func MediaLibraryStorer(b *media.Builder, db *gorm.DB) Storer {
	s := tiptap.MediaLibraryStorer(b, dbKey{})
	return StoreFunc(func(ctx context.Context, format string, data []byte, alt string) (string, error) {
		return s.Store(context.WithValue(ctx, dbKey{}, db), format, data, alt)
	})
}
