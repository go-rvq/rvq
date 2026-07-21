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

var (
	ErrInvalidCPF          = errors.New("invalid CPF")
	ErrInvalidCNPJ         = errors.New("invalid CNPJ")
	ErrInvalidDocumentType = errors.New("invalid document type")
)

// BeforeSave validates the document against its type (CPF/CNPJ when applicable).
func (p *Person) BeforeSave(*gorm.DB) error {
	switch p.DocumentType {
	case DocumentCPF:
		if p.Document != "" && !ValidCPF(p.Document) {
			return ErrInvalidCPF
		}
	case DocumentCNPJ:
		if p.Document != "" && !ValidCNPJ(p.Document) {
			return ErrInvalidCNPJ
		}
	case DocumentOther:
	default:
		return ErrInvalidDocumentType
	}
	return nil
}
