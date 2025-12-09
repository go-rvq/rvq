package media

import (
	"github.com/go-rvq/rvq/admin/media/base"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/gorm"
)

type Builder struct {
	db                  *gorm.DB
	permVerifier        *perm.Verifier
	mediaLibraryPerPage int

	base.WithConfigField
}

func New(db *gorm.DB) *Builder {
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
