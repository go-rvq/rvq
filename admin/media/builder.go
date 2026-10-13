package media

import (
	"sync"

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
	// usageModels are the models of the fields of images (RegisterUsages),
	// by their keys (*presets.ModelBuilder): what a recrop finds their
	// records by, what the admin names them by
	usageModels sync.Map

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

// BackfillUses makes the uses of the files again (media.BackfillUses) from
// the records of the models of pb with media boxes — after they were all
// registered. It answers how many uses it wrote.
func (b *Builder) BackfillUses(pb *presets.Builder) (int, error) {
	var models []any
	for _, m := range pb.Models() {
		models = append(models, m.NewModel())
	}
	return BackfillUses(b.db, models...)
}

func (b *Builder) Install(pb *presets.Builder) error {
	configure(pb, b, b.db)
	return nil
}
