# Permissions

The user's stack of permissions, opened in their detail (**Permissions**): the
system's first, then those granted them on a record, then those of each of
their roles, **in the order of the roles** — the groups broken, one
permission a line (`admin:users:{@list,@create}` are two).

- **All**: each line, with its **Index** (its place in the stack), the
  **Permission**, the **Origin** (the role it comes from — `role: Editor` —,
  the system — `code: …` —, a record — `database: …`) and the **State**:
  *active*, *active, but what #N matches* (a narrower `allow` or `reset`
  after it), *replaced by #N* (the line that replaced it) or *removal* (a
  `reset`, which only takes from the earlier ones).
- **Active**: only the ones that count — those of no prefix (an `allow` is
  shown without it: it allows) and the `deny` ones, each with what it no
  longer counts for ("but admin:users:@list"). **As a tree** shows them
  grouped by their common beginning:

  ```
  admin:users:{
  	@list
  	@create
  }
  ```

How a line replaces another is in
{%= admin.doc("guides/policies/01-permissions").link %}. The initial user may
do everything, by the system: their roles do not count.
