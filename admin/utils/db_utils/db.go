package db_utils

import (
	"fmt"
	"log"
	"runtime/debug"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Transact(db *gorm.DB, f func(tx *gorm.DB) error) (err error) {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			if er, ok := r.(error); ok {
				err = er
			} else {
				err = fmt.Errorf("%+v", r)
			}
			log.Println(err)
			debug.PrintStack()
		}

		if err != nil {
			if e := tx.Rollback().Error; e != nil {
				log.Println("Rollback Error:", e)
			}
		} else {
			err = tx.Commit().Error
		}
	}()

	err = f(tx)
	return
}

func FieldsNames(fields []*schema.Field) (names []string) {
	names = make([]string, len(fields))
	for i, field := range fields {
		names[i] = field.Name
	}
	return
}
