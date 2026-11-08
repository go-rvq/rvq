package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserDefault is the default organization of a user: the one the menu opens
// when none was chosen. One per user — its key is the user's id (a foreign
// key to the users, one to one) —, of an organization the user reaches.
type UserDefault struct {
	UserID        uuid.UUID    `gorm:"type:uuid;primaryKey"`
	OrganizacaoID uuid.UUID    `gorm:"type:uuid;not null;index"`
	Organizacao   *Organizacao `gorm:"constraint:OnDelete:CASCADE"`
	UpdatedAt     time.Time
}

func (UserDefault) TableName() string { return "org_user_defaults" }

// userFK is the name of the foreign key of UserDefault.UserID.
const userFK = "fk_org_user_defaults_user"

// DefaultOf is the default organization of the user (uuid.Nil: none).
func DefaultOf(db *gorm.DB, user uuid.UUID) uuid.UUID {
	var d UserDefault
	if db.Limit(1).Find(&d, "user_id = ?", user).RowsAffected == 0 {
		return uuid.Nil
	}
	return d.OrganizacaoID
}

// SetDefault makes org the default organization of the user.
func SetDefault(db *gorm.DB, user, org uuid.UUID) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"organizacao_id", "updated_at"}),
	}).Create(&UserDefault{UserID: user, OrganizacaoID: org, UpdatedAt: time.Now()}).Error
}

// ClearDefault leaves the user with no default organization.
func ClearDefault(db *gorm.DB, user uuid.UUID) error {
	return db.Where("user_id = ?", user).Delete(&UserDefault{}).Error
}

// migrateUserDefaults makes the table of the defaults — its key a foreign
// key to the users of usersTable (where the database can add one) — and the
// defaults of before (Organizacao.Padrao, its owner's) rows of it, the
// column gone.
func migrateUserDefaults(db *gorm.DB, usersTable string) error {
	if err := db.AutoMigrate(&UserDefault{}); err != nil {
		return err
	}
	m := db.Migrator()
	if usersTable != "" && db.Dialector.Name() == "postgres" && m.HasTable(usersTable) && !m.HasConstraint(&UserDefault{}, userFK) {
		if err := db.Exec(`ALTER TABLE org_user_defaults ADD CONSTRAINT `+userFK+` FOREIGN KEY (user_id) REFERENCES ?(id) ON DELETE CASCADE`,
			clause.Table{Name: usersTable}).Error; err != nil {
			return fmt.Errorf("orgs: the foreign key of the defaults to %s: %w", usersTable, err)
		}
	}
	if hasColumn(db, "padrao") {
		if err := db.Exec(`INSERT INTO org_user_defaults (user_id, organizacao_id, updated_at)
			SELECT owner_id, id, CURRENT_TIMESTAMP FROM org_organizacoes WHERE padrao AND deleted_at IS NULL
			ON CONFLICT (user_id) DO NOTHING`).Error; err != nil {
			return fmt.Errorf("orgs: the defaults of before: %w", err)
		}
		if err := db.Exec(`ALTER TABLE org_organizacoes DROP COLUMN padrao`).Error; err != nil {
			return err
		}
	}
	return nil
}

// hasColumn says the table of the organizations has the column name.
func hasColumn(db *gorm.DB, name string) bool {
	cols, err := db.Migrator().ColumnTypes(&Organizacao{})
	if err != nil {
		return false
	}
	for _, c := range cols {
		if c.Name() == name {
			return true
		}
	}
	return false
}
