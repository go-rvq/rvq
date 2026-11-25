package media

import (
	"testing"

	"github.com/go-rvq/rvq/admin/media/media_library"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// anexo is a file attached to a record of an owner.
type anexo struct {
	ID        uint `gorm:"primaryKey"`
	DonoID    uint
	Posicao   int
	Midia     media_library.MediaBox `gorm:"type:text"`
	Descricao string
	Org       string
}

// The attachments of an owner saved: the ones of before replaced, each in its
// order, the owner and what set says set; another owner's untouched.
func TestAttachmentsSave(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+t.Name()+"?mode=memory&cache=shared"), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&anexo{}); err != nil {
		t.Fatal(err)
	}
	pb := presets.New(i18n.New()).URIPrefix("/admin")
	a := FileAttachments(pb, db, &anexo{}).Fields("Midia", "Descricao", "Posicao")
	db.Create(&anexo{DonoID: 1, Descricao: "velho"})
	db.Create(&anexo{DonoID: 2, Descricao: "do outro"})
	kept := []any{&anexo{Descricao: "b"}, &anexo{Descricao: "a"}}
	if err := a.Save(db, "DonoID", uint(1), kept, func(obj any) { obj.(*anexo).Org = "acme" }); err != nil {
		t.Fatal(err)
	}
	var got []anexo
	db.Where("dono_id = ?", 1).Order("posicao").Find(&got)
	if len(got) != 2 || got[0].Descricao != "b" || got[0].Posicao != 1 || got[1].Descricao != "a" || got[1].Posicao != 2 || got[0].Org != "acme" {
		t.Errorf("the owner's: %+v", got)
	}
	var n int64
	db.Model(&anexo{}).Where("dono_id = ?", 2).Count(&n)
	if n != 1 {
		t.Errorf("another owner's: %d", n)
	}
	// of a polymorphic owner: its id and its kind
	if err := a.SaveOf(db, map[string]any{"DonoID": uint(3), "Org": "beta"}, []any{&anexo{Descricao: "c"}}, nil); err != nil {
		t.Fatal(err)
	}
	var c anexo
	db.Where("dono_id = ? AND org = ?", 3, "beta").First(&c)
	if c.Descricao != "c" || c.Posicao != 1 {
		t.Errorf("of the owner of two fields: %+v", c)
	}
	// the model, its fields: the table of the files
	if a.Model() == nil || a.Nested() == nil {
		t.Fatal("no model")
	}
	if FileIcon("contrato.PDF") == "mdi-file" || FileIcon("") != "mdi-file-outline" {
		t.Errorf("the icons: %s %s", FileIcon("contrato.PDF"), FileIcon(""))
	}
}
