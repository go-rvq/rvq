# User groups

A group gathers users that work alike: give the group roles, or share a
record with it, once — and every member has it.

- **Members** — the users of the group. Whoever is added reaches at once
  what the group reaches; whoever is taken out loses it at once.
- **Roles** — the roles of the group, **in order** (drag to change it), as
  a user's (see {%= admin.model("users").link %}): in a member's
  permissions they come after the member's own roles, the groups in the
  member's order of their groups.

A record shared with a group (an organization, say) needs no invite: the
owner confirms it by e-mail, and from then on every member reaches it. What
is shared with a member themselves counts before what is shared with their
groups.

Deleting a group takes its members out of it — each losing what they had by
it.
