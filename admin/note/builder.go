package note

import (
	"github.com/go-rvq/rvq/admin/presets"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type AfterCreateFunc func(db *gorm.DB) error

type Builder struct {
	db              *gorm.DB
	afterCreateFunc AfterCreateFunc
}

func New(db *gorm.DB) *Builder {
	b := &Builder{
		db: db,
	}
	return b
}

func (b *Builder) AfterCreate(f AfterCreateFunc) (r *Builder) {
	b.afterCreateFunc = f
	return b
}

func (b *Builder) Install(pb *presets.Builder) error {
	db := b.db
	if err := renameLegacy(db); err != nil {
		return err
	}
	if err := db.AutoMigrate(Note{}, UserNote{}); err != nil {
		return err
	}

	pb.I18n().
		RegisterForModule(language.English, I18nNoteKey, Messages_en_US).
		RegisterForModule(language.SimplifiedChinese, I18nNoteKey, Messages_zh_CN).
		RegisterForModule(language.Japanese, I18nNoteKey, Messages_ja_JP).
		RegisterForModule(language.BrazilianPortuguese, I18nNoteKey, Messages_pt_BR)
	return nil
}

func (b *Builder) ModelInstall(pb *presets.Builder, m *presets.ModelBuilder) error {
	db := b.db
	if m.Info().HasDetailing() {
		m.Detailing().AppendTabsPanelFunc(tabsPanel(db, m))
	}
	m.Editing().AppendTabsPanelFunc(tabsPanel(db, m))
	m.RegisterEventHandler(createNoteEvent, createNoteAction(b, m))
	m.RegisterEventHandler(updateUserNoteEvent, updateUserNoteAction(b, m))
	m.Listing().Field("Notes").ComponentFunc(noteFunc(db, m))
	return nil
}

// renameLegacy moves the table off the old "qor_" prefix, which the model
// carried until it was renamed. It runs before AutoMigrate, which would
// otherwise create the new table empty and leave the rows behind in the old
// one.
func renameLegacy(db *gorm.DB) error {
	m := db.Migrator()
	if m.HasTable("qor_notes") && !m.HasTable("notes") {
		return m.RenameTable("qor_notes", "notes")
	}
	return nil
}
