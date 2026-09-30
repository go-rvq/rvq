// Package softdelete is what a soft-deleted record keeps of its deletion:
// when (gorm's DeletedAt, which makes the delete soft), by whom and from
// where. A model embeds Deletion in place of a DeletedAt field.
package softdelete

import (
	"github.com/go-rvq/rvq/admin/origin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Deletion is the deletion of a record: when — DeletedAt, nil while it is not
// deleted —, by whom — DeletedByID, the user (a foreign key to the users, set
// up by the admin's trash) — and from where — DeletedOrigin, the columns
// deleted_origin_*. The admin's trash fills them in when it deletes a record,
// and clears them when it restores it.
type Deletion struct {
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	DeletedByID   *uuid.UUID     `gorm:"type:uuid"`
	DeletedOrigin origin.Origin  `gorm:"embedded;embeddedPrefix:deleted_origin_"`
}

// Columns are the values of the columns of a deletion — by, from o —, to set
// with its deleted_at (by nil: nobody known).
func Columns(by *uuid.UUID, o origin.Origin) map[string]any {
	cols := map[string]any{"deleted_by_id": by}
	cols["deleted_origin_ip"] = o.IP
	cols["deleted_origin_user_agent"] = o.UserAgent
	cols["deleted_origin_country"] = o.Country
	cols["deleted_origin_country_name"] = o.CountryName
	cols["deleted_origin_region"] = o.Region
	cols["deleted_origin_city"] = o.City
	cols["deleted_origin_latitude"] = o.Latitude
	cols["deleted_origin_longitude"] = o.Longitude
	return cols
}

// Restored are the values of the columns of a record that is not deleted
// anymore: no deletion, by nobody, from nowhere.
func Restored() map[string]any {
	cols := Columns(nil, origin.Origin{})
	cols["deleted_at"] = nil
	return cols
}
