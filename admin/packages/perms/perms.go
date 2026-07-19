// Package perms provides reusable, resource-agnostic permission tooling for
// rvq presets admins: a per-record permission manager (a shortcut to the perm
// API, backed by perm.DefaultDBPolicy) and a soft-delete trash view gated by a
// 'trash' listing permission. Both work with any ModelBuilder.
package perms

import (
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/perm"
	"gorm.io/gorm"
)

// Permission verbs granted per record. They map to the presets action names
// checked by the verifier.
var (
	VerbView   = []string{presets.PermList, presets.PermGet}
	VerbEdit   = []string{presets.PermUpdate}
	VerbDelete = []string{presets.PermDelete}
)

// PermTrash is the listing permission verb that gates the trash view and the
// restore action.
const PermTrash = "trash"

// PermManage is the perm verb (and action name) that guards the permission
// manager itself.
const PermManage = "managePermissions"

// RecordResource returns the exact permission resource string of a record, as
// computed by the presets permissioner (the module "presets" is already seeded
// into Resource(), so it is not added again).
func RecordResource(mb *presets.ModelBuilder, id model.ID) string {
	return mb.Permissioner().Verifier(id).Resource()
}

// ListResource returns the permission resource of the whole listing.
func ListResource(mb *presets.ModelBuilder) string {
	return mb.Permissioner().ListVerifier().Resource()
}

// Grant creates or updates a DB policy granting subject the actions on
// resource, keyed by referID (so it can be updated or revoked later).
func Grant(db *gorm.DB, referID, subject, resource string, actions []string) (*perm.DefaultDBPolicy, error) {
	p := &perm.DefaultDBPolicy{}
	err := db.Where("refer_id = ?", referID).First(p).Error
	switch err {
	case nil:
		p.Subject = subject
		p.Effect = perm.Allowed
		p.Actions = actions
		p.Resources = []string{resource}
		return p, db.Save(p).Error
	case gorm.ErrRecordNotFound:
		p = &perm.DefaultDBPolicy{
			ReferID:   referID,
			Subject:   subject,
			Effect:    perm.Allowed,
			Actions:   actions,
			Resources: []string{resource},
		}
		return p, db.Create(p).Error
	default:
		return nil, err
	}
}

// Revoke soft-deletes the policy identified by referID.
func Revoke(db *gorm.DB, referID string) error {
	return db.Where("refer_id = ?", referID).Delete(&perm.DefaultDBPolicy{}).Error
}

// List returns the active policies whose resources contain resource.
func List(db *gorm.DB, resource string) (out []perm.DefaultDBPolicy, err error) {
	err = db.Where("? = ANY(resources)", resource).Order("subject").Find(&out).Error
	return
}
