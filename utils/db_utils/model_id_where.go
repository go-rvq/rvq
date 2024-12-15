package db_utils

import (
	"fmt"
	"slices"

	"github.com/qor5/admin/v3/model"
	"gorm.io/gorm"
)

func ModelIdWhere(db *gorm.DB, obj any, id model.ID, withoutKeys ...string) *gorm.DB {
	if obj != nil {
		db = db.Model(obj)
	}

	if id.IsZero() {
		return db
	}

	for i, field := range id.Fields {
		if !slices.Contains(withoutKeys, field.Name()) {
			db = db.Where(fmt.Sprintf("%s = ?", field.QuotedFullDBName()), id.Values[i])
		}
	}

	return db
}
