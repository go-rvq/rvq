# Getting started

## The Builder

The admin is a `*presets.Builder`, created with an i18n builder:

```go
b := presets.New(i18n.New())
```

It holds the registered models, the i18n messages, the permission builder and the
default data operator. Configure cross-cutting concerns on it:

- `b.DataOperator(op)` — the default persistence (usually `gorm2op.DataOperator(db)`).
- `b.Permission(perm.New()...)` — the permission policies.
- `b.URIPrefix("/admin")`, `b.BrandFunc(...)`, `b.MenuOrder(...)`, layout hooks, etc.

The builder is an `http.Handler`; mount it on your router.

## Registering a model

```go
m := b.Model(&Post{})
```

`Model` returns a `*ModelBuilder` for the resource. Options can be passed:

```go
m := b.Model(&Post{}, presets.ModelWithID("post"))
```

Related helpers:

- `presets.NewModelBuilder(b, v, opts...)` — build a model without registering it
  on the menu (used for nested/child resources).
- `mb.ChildOf(parent)` — mark a model as a child of another.
- `mb.Singleton(true)` — a one-record resource (e.g. settings, profile).

Each model exposes the four sub-builders — `Listing`, `Editing`, `Creating`,
`Detailing` — described in [models.md](models.md).

## Plugins

Cross-cutting features install as plugins:

```go
b.Use(activityBuilder)          // Builder-level plugin (Install)
m.Use(activityBuilder)          // per-model plugin (ModelInstall)
```

A `Plugin` implements `Install(*Builder) error` and/or
`ModelInstall(*Builder, *ModelBuilder) error`. Examples: the `activity` audit
module, the `login`/`user`/`profile` helpers, and the packages under
`admin/packages`.

## Data operator

Persistence is abstracted by the `DataOperator` interface (search / fetch / save
/ delete). The GORM implementation lives in
[`gorm2op`](../gorm2op) (see its docs):

```go
b.DataOperator(gorm2op.DataOperator(db))
```

A model can override the operator (e.g. to add preloads or custom callbacks):

```go
m.DataOperator(gorm2op.DataOperator(db.Preload("Author")))
// or wrap the existing one:
m.UpdateDataOperator(func(op presets.DataOperator) presets.DataOperator { ... })
```
