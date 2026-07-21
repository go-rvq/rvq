package models

import (
	"errors"

	"gorm.io/gorm"
)

// DocumentType identifies the document kind of a Person.
type DocumentType string

const (
	DocumentCPF   DocumentType = "cpf"
	DocumentCNPJ  DocumentType = "cnpj"
	DocumentOther DocumentType = "other"
)

// DocumentTypes lists every valid value, in display order.
var DocumentTypes = []DocumentType{DocumentCPF, DocumentCNPJ, DocumentOther}

// Person is an individual, company or other identification, owned by an
// organization. Its financial receiving data (PIX keys, bank accounts, other
// forms) lives in the finance package (PeopleData), mounted as a nested
// sub-resource of this Person.
type Person struct {
	Base

	Name         string       `gorm:"size:255;not null"`
	Address      string       `gorm:""`
	DocumentType DocumentType `gorm:"size:10;not null;default:other"`
	Document     string       `gorm:"size:50"`
	Notes        string       `gorm:""`
	Active       bool         `gorm:"default:true"`
}

func (Person) TableName() string { return "people" }

func (p *Person) String() string { return p.Name }

var ErrInvalidDocumentType = errors.New("invalid document type")

// BeforeSave validates the document type is a known value. The document content
// itself (CPF/CNPJ) is validated through the validators registry — see the
// document validators (models/validators.go) applied as a field validator in the
// admin (FieldBuilder.Validators), instead of hardcoded Go functions.
func (p *Person) BeforeSave(*gorm.DB) error {
	switch p.DocumentType {
	case DocumentCPF, DocumentCNPJ, DocumentOther:
		return nil
	}
	return ErrInvalidDocumentType
}
