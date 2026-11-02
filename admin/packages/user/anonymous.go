package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AnonymousID is the id of the anonymous user: none (uuid.Nil). Anonymous is
// a request with no one logged in — never a row of the users: what it does is
// recorded with no author (model.UUID: NULL), and User.Anonymous() is a plain
// comparison against it.
var AnonymousID = uuid.Nil

// IsAnonymous reports whether id is the anonymous user's id: none.
func IsAnonymous(id uuid.UUID) bool { return id == uuid.Nil }

// legacyAnonymousID is the id the anonymous user had when it was a row of the
// users (a fixed UUID): DropLegacyAnonymous takes it out of the database.
var legacyAnonymousID = uuid.MustParse("00000000-0000-4000-a000-000000000001")

// anonymousFactory builds the request's anonymous user: the app's concrete User
// model, with no id and a name translated for the request. It is set by the
// user Builder; GetCurrentUser falls back to it when the session has no user,
// so a request always has a user (possibly Anonymous()), in memory only.
var anonymousFactory func(r *http.Request) User

// NewAnonymous builds an anonymous user of the concrete model type, with no
// id and the given (already translated) name.
func NewAnonymous(newModel func() User, name string) User {
	u := newModel()
	u.SetID(AnonymousID)
	u.SetName(name)
	return u
}

// DropLegacyAnonymous takes out of the database the anonymous user of before
// (a row of the users, legacyAnonymousID): what pointed to it — any uuid
// column of any table (the authors of the activity log, of the revisions, of
// the notes…) — points to no one (NULL), and the row is deleted. Nothing when
// there is no such row (a new database, or one done already).
func DropLegacyAnonymous(db *gorm.DB, newModel func() User) error {
	var count int64
	if err := db.Unscoped().Model(newModel()).Where("id = ?", legacyAnonymousID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return nil
	}
	usersTable, err := tableOf(db, newModel())
	if err != nil {
		return err
	}
	return db.Transaction(func(tx *gorm.DB) error {
		tables, err := tx.Migrator().GetTables()
		if err != nil {
			return err
		}
		for _, t := range tables {
			cols, err := tx.Migrator().ColumnTypes(t)
			if err != nil {
				return err
			}
			for _, c := range cols {
				if t == usersTable && c.Name() == "id" {
					continue
				}
				if pk, ok := c.PrimaryKey(); ok && pk {
					continue
				}
				if !isUUIDColumn(c.DatabaseTypeName()) {
					continue
				}
				// a savepoint each: a column that takes no NULL (none of an
				// author) is left as it is, the others go on
				q := fmt.Sprintf("UPDATE %s SET %s = NULL WHERE %s = ?", quote(t), quote(c.Name()), quote(c.Name()))
				if err := tx.Transaction(func(sp *gorm.DB) error { return sp.Exec(q, legacyAnonymousID).Error }); err != nil {
					log.Printf("user: the legacy anonymous user in %s.%s kept: %v", t, c.Name(), err)
				}
			}
		}
		return tx.Unscoped().Where("id = ?", legacyAnonymousID).Delete(newModel()).Error
	})
}

// tableOf is the table of the model m.
func tableOf(db *gorm.DB, m any) (string, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(m); err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}

func isUUIDColumn(t string) bool {
	switch t {
	case "uuid", "UUID":
		return true
	}
	return false
}

func quote(name string) string { return `"` + name + `"` }
