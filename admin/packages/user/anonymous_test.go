package user

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-rvq/rvq/admin/role"
	"github.com/go-rvq/rvq/x/login"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type anonUser struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	DeletedAt gorm.DeletedAt
	login.UserPass
}

func (u *anonUser) TableName() string               { return "users" }
func (u *anonUser) GetID() uuid.UUID                { return u.ID }
func (u *anonUser) SetID(v uuid.UUID)               { u.ID = v }
func (u *anonUser) Anonymous() bool                 { return IsAnonymous(u.ID) }
func (u *anonUser) GetName() string                 { return u.Name }
func (u *anonUser) SetName(v string)                { u.Name = v }
func (u *anonUser) SetEmail(v string)               { u.Account = v }
func (u *anonUser) SetRegistrationDate(v time.Time) {}
func (u *anonUser) GetStatus() string               { return "active" }
func (u *anonUser) GetRoles() role.Roles            { return nil }
func (u *anonUser) SetRoles(roles role.Roles)       {}

type anonLog struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uuid.UUID `gorm:"type:uuid"`
}

type anonNote struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	AuthorID uuid.UUID `gorm:"type:uuid"`
}

var anonSeq int64

// The anonymous user is no id: the row of before taken out, what pointed to
// it pointing to no one; the other users' kept; once done, nothing more.
func TestDropLegacyAnonymous(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:anon_%d?mode=memory&cache=shared", atomic.AddInt64(&anonSeq, 1))), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&anonUser{}, &anonLog{}, &anonNote{}); err != nil {
		t.Fatal(err)
	}
	ann := uuid.New()
	db.Create(&anonUser{ID: legacyAnonymousID, Name: "Anonymous"})
	db.Create(&anonUser{ID: ann, Name: "Ann"})
	db.Create(&anonLog{UserID: legacyAnonymousID})
	db.Create(&anonLog{UserID: ann})
	note := uuid.New()
	db.Create(&anonNote{ID: note, AuthorID: legacyAnonymousID})
	newModel := func() User { return &anonUser{} }

	for i := 0; i < 2; i++ {
		if err := DropLegacyAnonymous(db, newModel); err != nil {
			t.Fatalf("drop %d: %v", i, err)
		}
	}
	var users, anonLogs, nullLogs, annLogs, nullNotes int64
	db.Model(&anonUser{}).Count(&users)
	db.Model(&anonLog{}).Where("user_id = ?", legacyAnonymousID).Count(&anonLogs)
	db.Model(&anonLog{}).Where("user_id IS NULL").Count(&nullLogs)
	db.Model(&anonLog{}).Where("user_id = ?", ann).Count(&annLogs)
	db.Model(&anonNote{}).Where("id = ? AND author_id IS NULL", note).Count(&nullNotes)
	if users != 1 || anonLogs != 0 || nullLogs != 1 || annLogs != 1 || nullNotes != 1 {
		t.Fatalf("users %d, logs of the legacy %d, of none %d, of Ann %d, notes of none %d", users, anonLogs, nullLogs, annLogs, nullNotes)
	}
	if !IsAnonymous(uuid.Nil) || IsAnonymous(ann) || !NewAnonymous(newModel, "x").Anonymous() {
		t.Fatal("the anonymous is no id")
	}
}
