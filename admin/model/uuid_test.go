package model_test

import (
	"encoding/json"
	"testing"

	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/web/zeroer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// The zero UUID is none everywhere: NULL in the database (and NULL comes
// back as it), "" as text, null in JSON, zero for zeroer; a set one is its
// text.
func TestUUID(t *testing.T) {
	id := model.NewUUID()
	if id.IsZero() || !model.Nil.IsZero() || !zeroer.IsZero(model.Nil) || zeroer.IsZero(id) {
		t.Fatal("IsZero")
	}
	if v, err := model.Nil.Value(); v != nil || err != nil {
		t.Errorf("zero value: %v %v", v, err)
	}
	if v, _ := id.Value(); v != id.String() {
		t.Errorf("value: %v", v)
	}
	var back model.UUID = id
	if err := back.Scan(nil); err != nil || !back.IsZero() {
		t.Errorf("scan NULL: %v %v", back, err)
	}
	if err := back.Scan(id.String()); err != nil || back != id {
		t.Errorf("scan: %v %v", back, err)
	}
	if p, err := model.ParseUUID(""); err != nil || !p.IsZero() {
		t.Errorf("parse empty: %v %v", p, err)
	}
	if _, err := model.ParseUUID("x"); err == nil {
		t.Error("parse bad: no error")
	}

	var v struct {
		A model.UUID
		B model.UUID
	}
	v.A = id
	b, _ := json.Marshal(v)
	if string(b) != `{"A":"`+id.String()+`","B":null}` {
		t.Errorf("json: %s", b)
	}
	v.A, v.B = model.Nil, id
	if err := json.Unmarshal([]byte(`{"A":"`+id.String()+`","B":""}`), &v); err != nil || v.A != id || !v.B.IsZero() {
		t.Errorf("unjson: %+v %v", v, err)
	}
}

type uuidPage struct {
	ID   model.UUID `gorm:"primaryKey"`
	Name string
}

type uuidForm struct {
	ID     model.UUID `gorm:"primaryKey"`
	PageID model.UUID `gorm:"index"`
	Page   *uuidPage  `gorm:"foreignKey:PageID"`
	Name   string
}

// A foreign key left zero is NULL in its column — the constraint holds —
// and reads back as the zero; set, its value.
func TestUUIDColumn(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:uuidcol?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.Exec("PRAGMA foreign_keys = ON")
	if err = db.AutoMigrate(&uuidPage{}, &uuidForm{}); err != nil {
		t.Fatal(err)
	}
	f := &uuidForm{ID: model.NewUUID(), Name: "f"}
	if err = db.Create(f).Error; err != nil {
		t.Fatalf("create with no page: %v", err)
	}
	var null bool
	db.Raw("SELECT page_id IS NULL FROM uuid_forms WHERE id = ?", f.ID).Scan(&null)
	if !null {
		t.Error("the zero is not NULL in the column")
	}
	var got uuidForm
	if err = db.First(&got, "id = ?", f.ID).Error; err != nil || !got.PageID.IsZero() || got.ID != f.ID {
		t.Errorf("read: %+v %v", got, err)
	}
	p := &uuidPage{ID: model.NewUUID(), Name: "p"}
	db.Create(p)
	got.PageID = p.ID
	if err = db.Save(&got).Error; err != nil {
		t.Fatal(err)
	}
	got.PageID = model.Nil
	if err = db.Save(&got).Error; err != nil {
		t.Fatal(err)
	}
	db.Raw("SELECT page_id IS NULL FROM uuid_forms WHERE id = ?", f.ID).Scan(&null)
	if !null {
		t.Error("cleared, not NULL")
	}
}
