// Package shared provides reusable, resource-agnostic record sharing for rvq
// presets admins. Sharing a record with other subjects (users) grants them a
// set of verbs on that record, tracked as a "share" (perm.DefaultDBPolicy.
// SharedID) so the sharing UI can list and revoke it as a unit, separately from
// the ad-hoc grants made in the permission manager. It is a thin, user-facing
// layer over the perm API (see the perms package).
//
// A share is meant to be attached per resource (Install) — typically only on a
// few aggregate resources (e.g. an organization or a project) whose owner can
// always do everything and who decides who else may access it.
package shared

import (
	"github.com/go-rvq/rvq/admin/packages/perms"
	"github.com/go-rvq/rvq/x/perm"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Verb sets granted by a share (reused from the perms package so the sharing and
// permission-manager grants are interchangeable).
var (
	VerbView   = perms.VerbView
	VerbEdit   = perms.VerbEdit
	VerbDelete = perms.VerbDelete
)

// Create shares resource with each subject, granting them actions. Every grant
// is stamped with one new SharedID (returned), so the share can be listed and
// revoked as a unit. resource is the exact permission resource of the record
// being shared (e.g. perms.RecordResource(mb, id, parents...) + "*").
func Create(db *gorm.DB, resource string, subjects, actions []string) (uuid.UUID, error) {
	sharedID := uuid.New()
	for _, s := range subjects {
		if _, err := perms.GrantShared(db, ReferID(sharedID, s), s, resource, actions, sharedID); err != nil {
			return uuid.Nil, err
		}
	}
	return sharedID, nil
}

// AddSubjects grants more subjects the given actions within an existing share.
func AddSubjects(db *gorm.DB, sharedID uuid.UUID, resource string, subjects, actions []string) error {
	for _, s := range subjects {
		if _, err := perms.GrantShared(db, ReferID(sharedID, s), s, resource, actions, sharedID); err != nil {
			return err
		}
	}
	return nil
}

// Subjects returns the policies (one per subject) of a share.
func Subjects(db *gorm.DB, sharedID uuid.UUID) ([]perm.DefaultDBPolicy, error) {
	return perms.ListShared(db, sharedID)
}

// Revoke removes the whole share (every subject's grant).
func Revoke(db *gorm.DB, sharedID uuid.UUID) error {
	return perms.RevokeShared(db, sharedID)
}

// RevokeSubject removes a single subject from a share.
func RevokeSubject(db *gorm.DB, sharedID uuid.UUID, subject string) error {
	return perms.Revoke(db, ReferID(sharedID, subject))
}

// ForResource returns the share policies attached to resource (those carrying a
// SharedID), most-shared first. It is used by the sharing dialog to list who a
// record is currently shared with.
func ForResource(db *gorm.DB, resource string) ([]perm.DefaultDBPolicy, error) {
	pols, err := perms.List(db, resource)
	if err != nil {
		return nil, err
	}
	out := pols[:0:0]
	for _, p := range pols {
		if p.SharedID != nil {
			out = append(out, p)
		}
	}
	return out, nil
}

// ReferID is the stable refer id of a subject's grant within a share, so it can
// be updated or revoked individually.
func ReferID(sharedID uuid.UUID, subject string) string {
	return "share:" + sharedID.String() + ":" + subject
}
