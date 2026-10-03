package media

import (
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/admin/utils/uuidkey"
	"gorm.io/gorm"
)

type Builder struct {
	db                  *gorm.DB
	mediaLibraryPerPage int
	// model is the media library in the admin: its permissions are the
	// ones of what the forms do with the files (perm.go)
	model *presets.ModelBuilder

	base.WithConfigField
}

func New(db *gorm.DB) *Builder {
	uuidkey.MustRegister(db) // its records have UUID keys
	b := &Builder{}
	b.db = db
	b.mediaLibraryPerPage = 39
	return b
}

func (b *Builder) DB() *gorm.DB {
	return b.db
}

func (b *Builder) MediaLibraryPerPage(v int) *Builder {
	b.mediaLibraryPerPage = v
	return b
}

func (b *Builder) Install(pb *presets.Builder) error {
	configure(pb, b, b.db)
	return nil
}
