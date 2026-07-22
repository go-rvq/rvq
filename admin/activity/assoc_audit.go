package activity

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets/gorm2op"
	"gorm.io/gorm"
)

// assocAuditor adapts the activity Builder to gorm2op.AssociationAuditor so
// has-many children saved through gorm2op.SaveHasManyAssociation are logged.
//
// It is injected into the request context while an audited parent model is
// saved (see installModelBuilder). A child is only logged when its own type is
// also registered with activity; unregistered child types are silently skipped,
// so enabling audit on the parent never fails a save.
type assocAuditor struct {
	ab  *Builder
	ctx context.Context // carries the creator and request info of the parent action
}

func newAssocAuditor(ab *Builder, ctx context.Context) assocAuditor {
	return assocAuditor{ab: ab, ctx: ctx}
}

// dbCtx binds the transaction db of the current save to the base context so the
// context-based recording resolves the same db (and thus the same transaction).
func (a assocAuditor) dbCtx(db *gorm.DB) context.Context {
	return ContextWithDB(a.ctx, db)
}

func (a assocAuditor) LogCreated(db *gorm.DB, obj any) error {
	mb, ok := a.ab.GetModelBuilder(obj)
	if !ok {
		return nil
	}
	return mb.AddRecords(ActivityCreate, a.dbCtx(db), obj)
}

func (a assocAuditor) LogUpdated(db *gorm.DB, old, now any) error {
	mb, ok := a.ab.GetModelBuilder(now)
	if !ok {
		return nil
	}
	return mb.AddEditRecordWithOldCtx(a.dbCtx(db), old, now)
}

func (a assocAuditor) LogDeleted(db *gorm.DB, obj any) error {
	mb, ok := a.ab.GetModelBuilder(obj)
	if !ok {
		return nil
	}
	return mb.AddRecords(ActivityDelete, a.dbCtx(db), obj)
}

// injectAssociationAuditor stores an association auditor bound to r's origin
// (creator + request info) into r's context, returning the updated request.
func (ab *Builder) injectAssociationAuditor(ctx context.Context) context.Context {
	return gorm2op.ContextWithAssociationAuditor(ctx, newAssocAuditor(ab, ctx))
}
