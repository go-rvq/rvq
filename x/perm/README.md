# perm

`perm` decides whether a request may do an action on a resource. Two layers:

1. **The policies** ([ladon](https://github.com/ory/ladon)): `PolicyFor(subjects).WhoAre(Allowed|Denied).ToDo(actions).On(resources)`,
   from the code (`Builder.Policies`, `CreatePolicies`) or from the database
   (`DBPolicy`, `DefaultDBPolicy`). Unchanged by the stacks: they are still
   registered and matched as always (`PathMatcher`: gobwas globs, `*` matches
   any characters, `:` too).
2. **The stacks** (`stack.go`, `stacks.go`): the permissions of a user as an
   ordered list of lines — the code's, the database's, then each of the user's
   roles, in the user's order — resolved into the lines that count (the
   *active* ones), which are registered as ordinary policies of one subject,
   the user's (`user:<id>`). A user's request is decided by them alone.

The stacks change only what is sent to the policies: the resolution happens
before, and the decision is the one the policies always made — **a resource
is allowed when an active line of no prefix matches it and no active `deny`
does**.

## Lines

A line is a resource pattern (a glob, as the policies') and maybe a prefix:
a word and a space.

| Line | Kind | What it does |
|---|---|---|
| `admin:pages:*` | `LineGrant` | allows what it matches, unless an active deny matches it too |
| `deny admin:users:*` | `LineDeny` | denies what it matches |
| `allow admin:users:@list` | `LineAllow` | the denies before it stop counting for what it matches — a wider one too; then it is a line of no prefix: it forces the allow of what it matches |
| `reset admin:users:*` | `LineReset` | the lines before it stop counting for what it matches — a wider one too; it allows nor denies anything |

- The prefixes: `PrefixDeny` (`"deny "`), `PrefixAllow` (`"allow "`),
  `PrefixReset` (`"reset "`). Without its space a word is part of the
  resource: `denyadmin:x` is a grant of the resource `denyadmin:x` — and a
  prefix alone (`deny`) is wrong.
- `ParseLine(s)` is the kind and the pattern; `CheckLine(s)` says what is
  wrong (`ErrLine`: a prefix with no resource, a group not closed); a blank
  line is nothing (not an error).
- Every line has a context: there is no `deny` alone, it is
  `deny admin:…` or something else.

### Groups

A group of alternatives `{a,b}` is several lines, in its order. A new line
inside a group separates alternatives as a comma does.

```
admin:users:{list,new}            → admin:users:list
                                    admin:users:new
admin:users:{*, deny @get}        → admin:users:*
                                    deny admin:users:@get
admin:users:{*, deny {get, set}}  → admin:users:*
                                    deny admin:users:get
                                    deny admin:users:set
admin:{pages,posts}:{*, deny @delete}
                                  → admin:pages:*
                                    deny admin:pages:@delete
                                    admin:posts:*
                                    deny admin:posts:@delete
```

An alternative may carry its own prefix (it applies to that alternative
alone); a prefix before the group applies to all of it
(`deny admin:{a,b}` → `deny admin:a`, `deny admin:b`); groups nest. This
works anywhere in the resource: models, fields (`!field`), actions (`@verb`).
`ExpandLine(s)` gives the lines; `Stack.Add` expands what it adds, so a stack
never holds a group.

## Resolution

`Stack.Add(origin, lines...)` appends lines (expanded, the blank and the
wrong ones left out), each with its `Origin` — `"code: Administrador"`,
`"database: user:…"`, `"role: Editor"` — and its `Index` (from 1, its place
in the stack). `Stack.Resolve()` sets `ReplacedBy` of each line: the index of
the **first later line that covers and replaces it** (0: none).

- **Covers**: a line covers an earlier one when its pattern matches the
  earlier one's pattern (as text): `admin:*` covers `admin:users:list`;
  `admin:users:list` does not cover `admin:users:*`. With actions (lines of
  the code that name them), the later one's actions must be any, or include
  the earlier one's. A line with conditions neither covers nor is covered.
- **Replaces** (when it covers):

  | Later line | Replaces an earlier |
  |---|---|
  | `deny …` | line of no prefix, `allow`, `deny` |
  | `allow …` | line of no prefix, `allow`, `deny` |
  | `reset …` | any line |
  | no prefix | line of no prefix, `allow` |

- **Excepts** (when the earlier one covers the later, and not the other way:
  the later is narrower): an `allow` after a `deny`, a `reset` after any line
  — the earlier one stays active, but no longer counts for what the later one
  matches: its `ExceptBy` (the indexes) and `Except` (the patterns). A line
  replaced has no exceptions.
- **Active** (`Line.Active`, `Stack.Active`): not replaced, and not a reset.
  The active lines are only of no prefix (an `allow` counts as one: it is
  registered as a plain allow) or `deny`, each with its exceptions.

### Examples

```
1 admin:*                         active
2 deny admin:users:*              active        → users denied
```

```
1 admin:*                         active
2 deny admin:users:*              replaced by 3
3 allow admin:users:*             active        → users allowed again
```

```
1 admin:*                         active
2 deny admin:users:*              replaced by 3
3 reset admin:users:*             (a removal)   → users allowed by 1 again
```

A narrower `allow` or `reset` after a wider line does not replace it: it
takes from it what it matches (an exception) — the `allow` forces the allow
of just that:

```
1 admin:*                         active
2 deny admin:users:*              active, but what 3 matches
3 allow admin:users:@list         active        → the users denied, but @list
```

```
1 admin:*                         active, but what 2 matches
2 reset admin:users:@get          (a removal)   → all but users:@get (nothing
                                                   allows it now)
```

A narrower line of no prefix does not take anything: a deny wins over it, as
in the policies:

```
1 deny admin:users:*              active
2 admin:users:@list               active        → @list still denied (by 1)
```

The exceptions are by the patterns as text, as the covering: two lines that
overlap, neither matching the other's pattern (`admin:*:@list` and
`admin:users:*`), do not except each other.

```
1 admin:users:{*, deny @get}  →  1 admin:users:*         active
                                 2 deny admin:users:@get active   → all but @get
```

```
1 admin:users:*                   replaced by 2
2 reset admin:*                   (a removal)
3 admin:pages:*                   active        → only pages
```

### The tree

`Tree(lines)` writes lines (the active ones: their order does not matter) as
one, indented, their common beginnings once, segment by segment (`:`), the
prefix before each alternative:

```
admin:{
	pages:*
	users:{
		list
		new
		deny @get
	}
}
```

It is a line itself: `ExpandLine(Tree(lines))` gives the lines back.

## The stacks of the users

```go
pb.Stacks(src) // src: a StackSource
```

```go
type StackSource interface {
	StackOf(r *http.Request) (subject, version string, build func() (*Stack, error))
}
```

For each request the builder asks the source for the user's stack:

- `subject` — what the active lines are registered as (`user:<id>`); `""`
  when the request has no user (the anonymous): **it is decided by the
  policies, as before**.
- `version` — what the stack is made of (its roles, their lines, their
  order). The stack is made again only when it changes, or when any policy
  of the builder changed (of the code or of the database: every stack is
  made again).
- `build` — makes the stack, not resolved: `Register` resolves it.

`Register(subject, version, build)` builds (outside of the lock: building
reads the policies), resolves, and puts the policies of the active lines
(`StackPolicies`: IDs `stack:<subject>:<index>`, a deny `Denied`, the others
`Allowed`, the actions of the line or any, its conditions and, of a line
with exceptions, an `ExceptCondition` — fulfilled by a request whose resource
matches none of them) into a ladon of the stacks, apart from the builder's, taking out those of the version before.
`Registered(subject)`, `Forget(subjects...)` (all of them when none: made
again at their next request).

The verifier then asks as that subject alone, in the stacks' ladon
(`ladonOf`): the policies of the code and of the database no longer decide a
user's request — they are what the stacks are made of. **The cache of the
permissions is made of the active lines only.**

### What a stack is made of

`Builder.CodeLines(subjects)` are the lines of the policies of the code whose
subjects match (globs) one of `subjects` — origin `code: <subjects>`;
`Builder.DBLines(subjects)` the database's (what is granted a record, a
share: `admin/packages/perms`) — origin `database: <subjects>`; then each
role's lines. `role.Stack(pb, subjects, roles)` puts them in this order;
`user.Builder.Stack(pb, u)` uses the user's subjects (`UserRoles(u)` and
`user:<id>`) and the user's roles in the user's order — the initial account
(`LOGIN_INITIAL_USER_*`) has the code's lines only (its role,
Administrador, allows all).

`user.Builder.StackSource(pb)` is the source of a `user` admin; its version is
`role.Generation()` and the IDs of the user's roles in their order.

### Roles

- A role's `Permissions` are lines (`pq.StringArray`, `text[]`), in order;
  the admin edits them as a sortable list of texts, each checked by
  `CheckLine`.
- A role saved or deleted: `role.Changed(pb)` (a new generation, every stack
  forgotten). In an application of several processes the other ones must be
  told: HermonCMS announces every change of the tables `default_db_policies`
  and `roles` on a Postgres channel, and each process reloads its policies
  and calls `role.Changed`.
- The roles of the system (`SystemRole.Lines`, or its `Policies` made lines
  with `SystemPolicy.Lines()`) keep their originals as lines; a reset puts
  them back.
- `role.MigrateLines(db)` (boot): the policies of a role of before (`ReferID`
  its id) become its lines, in the order they were made — a deny `deny `
  before it —, and are deleted; the old originals
  (`effect|actions|resources`) become lines.

### The order of a user's roles

`user_role_join.position` (`user.EnsureRolePositions`, added when missing);
the user form lists the roles as a sortable list; `OrderedRoles`,
`SetRolePositions`. A later role's lines come later in the stack: a `deny` of
it wins over what an earlier one allows; an `allow` of it forces what an
earlier one denies.

### The page of a user's permissions

The user's detail has a page *Permissões*: **Tudo** — the whole stack, a
table of Índice, Permissão, Origem, Situação (ativa; ativa, menos o que #N
casa — its exceptions —; substituída pela #N; or remoção for a reset) — and
**Ativas** — the active lines (an `allow` shown without its prefix: it is a
plain allow), each with what it no longer counts for ("menos …"), as a list
or as a tree (`Tree`, the exceptions under it).

## Access keys

A request by an access key (`admin/packages/accesskeys`) is its user's,
restricted: `WithRestriction(ctx, "key:<id>")`. The verifier decides as the
user, then, when allowed, decides again as the restriction's subjects — an
intersection.

- A key's `Permissions` accept only `deny` lines (each line of it, its groups
  broken: `admin:users:{deny @edit, @list}` is refused). They are applied on
  top of what the user may: the key's stack is `*`, its denies, then
  `ProtectedLines` (no key reaches the keys, theirs, their uses), registered
  as `key:<id>` (`AccessKey.Stack`, version its update time and its lines).
- No lines: everything the user may (but the keys).
- A key of before (lines of no prefix, or policies made lines by
  `accesskeys.MigrateLines`): its lines are what it reaches — kept until it
  is saved again with denies only.

## Anonymous requests and code without stacks

A builder with no `Stacks` decides every request by the policies, as always
(`SubjectsFunc`). With `Stacks`, a request whose source gives no subject
(the anonymous one) too.

## API summary

| | |
|---|---|
| `PrefixDeny`, `PrefixAllow`, `PrefixReset`, `LineKind` | the prefixes and the kinds |
| `ParseLine`, `CheckLine`, `ExpandLine`, `ErrLine` | a line |
| `Line{Index, Kind, Pattern, Origin, Actions, Conditions, ReplacedBy}`, `Text`, `Active`, `Denies` | a line of a stack |
| `Stack{Lines}`, `Add`, `AddLine`, `Resolve`, `Active` | a stack |
| `Builder.CodeLines`, `Builder.DBLines` | the policies as lines |
| `PolicyLines(effect, actions, resources)` | a policy of before as lines (migrations) |
| `Tree(lines)` | the lines as a tree |
| `StackSource`, `Builder.Stacks`, `GetStacks`, `Register`, `Registered`, `Forget`, `StackPolicies`, `ExceptCondition` | the stacks |
| `WithRestriction`, `RestrictionOf` | a request restricted (access keys) |
