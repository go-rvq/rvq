package db_tools

import "context"

// BackupMethod is a way the backups are made, as it is offered: whether it
// may be used now and, when not, why.
type BackupMethod struct {
	Name      string
	Label     string
	Available bool
	// Reason is why it is not available, in the language of the context.
	Reason string
	// Version is what makes it ("pg_dump (PostgreSQL) 17.11"), when known.
	Version string
}

// BackupMethods are the methods of a BackupController that has more than one
// way to back the database up: the user chooses one (its Name, kept in the
// configuration of the backups; "" is the controller's default).
type BackupMethods struct {
	List func(ctx context.Context) []BackupMethod
	// Get is the method backups are made with now.
	Get func() string
	// Set chooses the method backups are made with.
	Set func(name string) error
}

// ContextCreator is a BackupController whose backup stops when its context
// is done (the job that makes it, aborted).
type ContextCreator interface {
	CreateContext(ctx context.Context, auto bool, message string) (Backuper, error)
}

// CreateBackup is a backup by c: with ctx, when c takes it (ContextCreator).
func CreateBackup(ctx context.Context, c BackupController, auto bool, message string) (Backuper, error) {
	if cc, ok := c.(ContextCreator); ok {
		return cc.CreateContext(ctx, auto, message)
	}
	return c.Create(auto, message)
}
