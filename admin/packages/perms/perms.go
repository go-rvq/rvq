// Package perms provides reusable, resource-agnostic permission tooling for
// rvq presets admins: a per-record permission manager (a shortcut to the perm
// API, backed by perm.DefaultDBPolicy) and a soft-delete trash view gated by a
// 'trash' listing permission. Both work with any ModelBuilder.
package perms

import (
	"github.com/go-rvq/rvq/admin/model"
	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
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

// GrantShared is like Grant but stamps the policy with sharedID, marking it as
// part of a resource share (organization/project). On update the SharedID is
// (re)applied so a policy adopted by a share is tracked.
func GrantShared(db *gorm.DB, referID, subject, resource string, actions []string, sharedID uuid.UUID) (*perm.DefaultDBPolicy, error) {
	p, err := Grant(db, referID, subject, resource, actions)
	if err != nil {
		return nil, err
	}
	p.SharedID = &sharedID
	return p, db.Model(p).Update("shared_id", sharedID).Error
}

// Revoke soft-deletes the policy identified by referID.
func Revoke(db *gorm.DB, referID string) error {
	return db.Where("refer_id = ?", referID).Delete(&perm.DefaultDBPolicy{}).Error
}

// RevokeShared soft-deletes every policy belonging to the given share.
func RevokeShared(db *gorm.DB, sharedID uuid.UUID) error {
	return db.Where("shared_id = ?", sharedID).Delete(&perm.DefaultDBPolicy{}).Error
}

// ListShared returns the active policies belonging to the given share, ordered
// by subject.
func ListShared(db *gorm.DB, sharedID uuid.UUID) (out []perm.DefaultDBPolicy, err error) {
	err = db.Where("shared_id = ?", sharedID).Order("subject").Find(&out).Error
	return
}

// IsShared reports whether the policy identified by referID carries a SharedID
// (and therefore may only be removed through the sharing UI, not the permission
// manager).
func IsShared(db *gorm.DB, referID string) (bool, error) {
	var p perm.DefaultDBPolicy
	if err := db.Select("shared_id").Where("refer_id = ?", referID).First(&p).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return p.SharedID != nil, nil
}

// List returns the active policies whose resources contain resource.
func List(db *gorm.DB, resource string) (out []perm.DefaultDBPolicy, err error) {
	err = db.Where("? = ANY(resources)", resource).Order("subject").Find(&out).Error
	return
}
