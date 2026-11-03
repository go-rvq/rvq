package models

import (
	"errors"
	"strings"

	"github.com/go-rvq/rvq/admin/model"

	"gorm.io/gorm"
)

// Organizacao is an organization (company, entity, etc.) owned by a user. It is
// the container every domain module (e.g. finance) scopes its records to. The
// owner can always do everything; access for other users is granted through
// shares (see the admin sharing manager).
type Organizacao struct {
	Base

	Nome      string `gorm:"size:255;not null"`
	Descricao string

	// OwnerID is the owning user (a hermon/rvq user UUID). The owner is never
	// removed from the organization and always has full access.
	OwnerID model.UUID `gorm:"type:uuid;not null;index"`

	// Padrao: the default organization of its owner — the one the menu opens
	// when none was chosen (one per owner: marking one unmarks the others).
	Padrao bool `gorm:"not null;default:false"`
}

func (Organizacao) TableName() string { return "org_organizacoes" }

// ErrNomeObrigatorio is returned when an organization has no name.
var ErrNomeObrigatorio = errors.New("o nome da organização é obrigatório")

// ErrProprietarioObrigatorio is returned when an organization has no owner.
var ErrProprietarioObrigatorio = errors.New("a organização deve ter um proprietário")

// BeforeSave validates the required fields (SPEC: name and owner).
func (o *Organizacao) BeforeSave(*gorm.DB) error {
	if strings.TrimSpace(o.Nome) == "" {
		return ErrNomeObrigatorio
	}
	if o.OwnerID.IsZero() {
		return ErrProprietarioObrigatorio
	}
	return nil
}

// AfterSave keeps one default organization per owner: this one marked, the
// others of its owner unmarked.
func (o *Organizacao) AfterSave(tx *gorm.DB) error {
	if !o.Padrao {
		return nil
	}
	return tx.Session(&gorm.Session{NewDB: true}).Model(&Organizacao{}).
		Where("owner_id = ? AND id <> ? AND padrao", o.OwnerID, o.ID).
		UpdateColumn("padrao", false).Error
}

func (o *Organizacao) String() string { return o.Nome }
