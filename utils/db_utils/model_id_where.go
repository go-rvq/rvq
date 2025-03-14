package db_utils

import (
	"fmt"
	"slices"
	"strings"

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

func ModelIdsWhere(db *gorm.DB, ids model.IDSlice, withoutKeys ...string) *gorm.DB {
	var (
		allq []string
		args []any
	)

	for _, mid := range ids {
		var (
			q []string
		)

		for i, field := range mid.Fields {
			if !slices.Contains(withoutKeys, field.Name()) {
				q = append(q, fmt.Sprintf("%s = ?", field.QuotedFullDBName()))
				args = append(args, mid.Values[i])
			}
		}

		allq = append(allq, "("+strings.Join(q, " AND ")+")")
	}

	return db.Where("( "+strings.Join(allq, " OR ")+" )", args...)
}
