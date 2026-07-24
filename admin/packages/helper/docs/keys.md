# Primary & foreign keys in `helper`

This note explains how `helper.NestedSlice` and `helper.ModelSelectorBuilder`
handle record keys, why the old code was integer-only, and how the current code
supports any key type and composite (multi-column) keys.

## The old limitation

The many-to-many machinery built raw SQL that assumed a single integer column
named `id`:

```sql
-- old: cast to a fixed type, hardcoded column
select id::BIGINT as f_id, ?::BIGINT as p_id from related where id in ?
```

and rejected composite keys outright:

```go
if len(FieldModel.Schema().PrimaryFields()) > 1 {
    panic("NestSlice doesn't supports ModelBuilder with many primary fields")
}
```

So a related model keyed by a UUID (a `string` column) or by two columns could
not be used.

## The approach now

Keys are never assumed: they come from the schema/relationship and are used as
their own Go/SQL type.

### Reading (filter)

For **listing** the related rows that belong to a parent, the WHERE clause is
built from the gorm relationship's `References` — one condition per key column:

- has-many (`hasManyParentFilter`): `fk_a = ? AND fk_b = ?`
- many-to-many (`m2mParentFilter`):
  `EXISTS (SELECT 1 FROM <join> j WHERE j.owner_a = ? AND j.rel_a = <related>.pk_a AND ...)`

The `?` placeholders are bound from the parent id in reference order
(`parentFilterArgs(id, ownerFields)`), where `id` is a `model.ID` that already
carries every parent primary-key value.

### Writing (link / unlink)

For **linking and unlinking** related rows, the code delegates to gorm's
Association API instead of raw SQL:

```go
db.Model(parent).Association(field).Append(records...) // link
db.Model(parent).Association(field).Delete(record)     // unlink (join row only)
```

gorm builds the join rows from each side's primary key(s), so integer, UUID and
composite keys, and duplicate-link avoidance, are all handled by gorm. `parent`
and `record` are bare structs carrying only their (possibly composite) key, set
via `model.ID.SetTo`.

### Belongs-to selector

`ModelSelectorBuilder` resolves the owner's foreign-key columns from the
relationship (`foreignKeyFieldsOf`) in related-primary-key order. When a record
is picked in the UI, its `model.ID` (parsed from the selected value string) is
mapped onto those columns:

```go
id.Related(ownerSchema, fkFields...).SetTo(owner)
```

so a composite foreign key is written in full; clearing the selection zeroes
every foreign-key column. A single key still uses the `<Field>ID` convention.

## Worked example — composite many-to-many

```go
type Item struct {            // composite primary key
    OrgID uint   `gorm:"primaryKey"`
    Code  string `gorm:"primaryKey"`
    Name  string
}

type Parent struct {
    ID    uint
    Items []*Item `gorm:"many2many:parent_items;"`
}
```

`m2mParentFilter` yields, for `Parent.Items`:

```sql
EXISTS (SELECT 1 FROM parent_items j
        WHERE j.parent_id = ?              -- owner key (bound to Parent.ID)
          AND j.item_org_id = items.org_id -- related key columns, correlated
          AND j.item_code   = items.code)
```

Linking `Append([]*Item{{OrgID:1,Code:"a"}, {OrgID:2,Code:"a"}})` inserts the two
join rows; `Delete(&Item{OrgID:1,Code:"a"})` removes just that join row. The
`items` rows themselves are never modified by link/unlink.

See `nestedslice_test.go` (`TestM2MParentFilter_CompositeRelatedKey`,
`TestM2MLinkQuery_NonIntegerKey`) and `model_select_test.go`
(`TestForeignKeyFieldsOf_Composite`).
