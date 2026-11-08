// Package shared is the sharing of the records of any model (Service): a
// subject — a user or a group — given a role on a record (Access), inherited
// by the records under it; every grant and every change of the owner
// confirmed by the owner by e-mail (Request), a user's then accepted by them
// (AccessInvite); all of it in the history (Event). In the admin
// (Service.Install) the accesses are kept as policies of the permissions
// (perms.GrantShared), the reads scoped to the records reached.
package shared

import "github.com/go-rvq/rvq/admin/packages/perms"

// Verb sets of the roles (the perms package's, so the sharing and the
// permission manager speak the same verbs).
var (
	VerbView   = perms.VerbView
	VerbEdit   = perms.VerbEdit
	VerbDelete = perms.VerbDelete
)
