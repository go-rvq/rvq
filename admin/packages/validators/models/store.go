package models

import (
	"errors"

	"gorm.io/gorm"
)

// Get loads the validator registered under name.
func Get(db *gorm.DB, name string) (*Validator, error) {
	var v Validator
	if err := db.Where("name = ?", name).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// Check loads the validator named name and runs it against value, returning a
// language-translated error when the value is invalid (see Validator.Check). A
// missing validator returns the underlying gorm error.
func Check(db *gorm.DB, name, value, lang string) error {
	v, err := Get(db, name)
	if err != nil {
		return err
	}
	return v.Check(value, lang)
}

// Seed upserts a registered validator by Name: it inserts it when missing and,
// when it already exists, refreshes the application-owned registration fields
// (Description, DefaultLang, Initial*) while preserving the user overrides
// (Value/Doc/Messages). Use it at application start to (re)register validators.
func Seed(db *gorm.DB, v *Validator) error {
	var existing Validator
	err := db.Where("name = ?", v.Name).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return db.Create(v).Error
	}
	if err != nil {
		return err
	}
	existing.Description = v.Description
	existing.DefaultLang = v.DefaultLang
	existing.InitialValue = v.InitialValue
	existing.InitialDoc = v.InitialDoc
	existing.InitialMessages = v.InitialMessages
	return db.Save(&existing).Error
}
