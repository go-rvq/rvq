package gorm2op

import (
	"context"

	"gorm.io/gorm"
)

// AssociationAuditor records create/update/delete of association child records
// while a has-many field is saved by SaveHasManyAssociation.
//
// It is resolved from the request context and used only when present, so
// gorm2op stays decoupled from the activity module (no import of it). The
// activity module injects an implementation while it saves an *audited* parent
// model, which means children are logged only when the model that contains the
// association is itself audited.
//
// Each method receives the transaction *gorm.DB used for the save, so log rows
// are written inside the same transaction as the data change. A nil auditor
// (nothing injected) means "do not log".
type AssociationAuditor interface {
	// LogCreated logs the creation of a child (called after it has an id).
	LogCreated(db *gorm.DB, obj any) error
	// LogUpdated logs an edit of a child, given the previous and current rows.
	LogUpdated(db *gorm.DB, old, now any) error
	// LogDeleted logs the removal of a child, given its row before deletion.
	LogDeleted(db *gorm.DB, obj any) error
}

type assocAuditorKey struct{}

// ContextWithAssociationAuditor stores a into ctx so nested has-many saves can
// log their child changes.
func ContextWithAssociationAuditor(ctx context.Context, a AssociationAuditor) context.Context {
	return context.WithValue(ctx, assocAuditorKey{}, a)
}

// AssociationAuditorFromContext returns the auditor stored by
// ContextWithAssociationAuditor, or nil when none was injected.
func AssociationAuditorFromContext(ctx context.Context) AssociationAuditor {
	a, _ := ctx.Value(assocAuditorKey{}).(AssociationAuditor)
	return a
}
