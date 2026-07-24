# Keys & relationships

This note explains how `thirdpart/gorm/utils` reads gorm keys/relationships so
that **any key type (int, UUID, string) and composite (multi-column) keys** work,
and why the helpers live here (a single source of truth).

## Why this package

The same key-reading logic used to be duplicated (and integer-only) in several
places — `admin/packages/helper`, `admin/presets/gorm2op`, and an application-side
`gormutils`. They assumed a single integer `id` column, e.g.:

```sql
select id::BIGINT as f_id ... where id in ?     -- broken for UUID / composite
```

These helpers replace that: read the key columns from the schema/relationship and
use them as their own type; iterate every reference so composite keys are handled.

## Reading a relationship

gorm exposes a relationship as `db.Model(x).Association(field).Relationship`, whose
`References []*schema.Reference` map each key column across the two sides:

```go
type Reference struct {
    PrimaryKey    *Field // the PK field (owner or related)
    ForeignKey    *Field // the FK field (owner column, or join-table column)
    OwnPrimaryKey bool   // true when PrimaryKey is the OWNER's
}
```

`ForeignKeyFields`, `HasManyParentFilter` and `M2MParentFilter` iterate the
references, so a belongs-to / has-many / m2m with two or more key columns is
matched in full.

## Filters

`HasManyParentFilter(rel)` → `WHERE fk_a = ? AND fk_b = ?` (one per reference),
plus the owner PK field names for `ParentFilterArgs`.

`M2MParentFilter(rel)` → an `EXISTS(...)` over the join table:

```sql
EXISTS (SELECT 1 FROM <join> j
        WHERE j.owner_a = ?                    -- owner side, bound to the parent id
          AND j.rel_a   = <related>.pk_a       -- related side, correlated
          AND j.rel_b   = <related>.pk_b)
```

Bind the `?` from the parent id in reference order:

```go
sql, ownerFields := utils.M2MParentFilter(rel)
db.Where(sql, utils.ParentFilterArgs(parentID, ownerFields)...) // parentID is a model.ID
```

## Writing (link / unlink)

Do NOT build raw join SQL. gorm's Association API handles any key type, composite
keys and duplicate links:

```go
db.Model(parent).Association(field).Append(records...) // link
db.Model(parent).Association(field).Delete(record)     // unlink (join row only)
```

`parent` and `record` need only carry their (possibly composite) key — set via
`model.ID.SetTo`.

## Primary-key matching (reconciliation)

When diffing a submitted slice against persisted rows (e.g.
`gorm2op.SaveHasManyAssociation`):

```go
pkNames := utils.PrimaryFieldNames(schema)      // every PK field
if utils.PKAllZero(item, pkNames) { /* new / keyless */ }
key := utils.PKMapKey(item, pkNames)            // stable composite map key
```

## Belongs-to foreign key

`ForeignKeyFields(db, owner, field)` returns the owner's FK columns in
related-primary-key order, so a selected record's `model.ID` maps onto them all:

```go
fks := utils.ForeignKeyFields(db, &Owner{}, "Target")   // e.g. [TargetOrgID TargetCode]
id.Related(ownerSchema, fks...).SetTo(owner)            // writes every FK column
```

See `relationship_test.go` for composite has-many, many-to-many and belongs-to
tests.
