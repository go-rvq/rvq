# Developing a package

A *package* here is anything that installs models into a host application's
`presets.Builder` without being the host: everything under `admin/packages/`,
the features under `admin/` (worker, note, activity, seo, pagebuilder), and
whatever an application adds from outside this repository.

The host owns the obvious names. A package that takes one takes it from an
application that will want it — and the clash is not caught at compile time, it
is caught at boot, or not at all.

## Tables carry the package prefix

Every table a package creates is named after the package:

```go
func (Organizacao) TableName() string { return "org_organizacoes" }
func (ShareInvite) TableName() string { return "shared_invites" }
func (MailSender) TableName() string  { return "mail_sender_settings" }
```

Never rely on the name gorm derives from the type — that name is the bare one
(`Category` → `categories`), and it is exactly the one the application wants.

The prefix is the package name in snake_case, with a trailing underscore:
`orgs` → `org_`, `shared` → `shared_`, `mail-sender` → `mail_sender_`.

**What it costs when you skip it.** The pagebuilder's `Category` and an
application's own `Category` both wanted `categories`. The two models were told
apart by `URIName` alone, so they shared the table name, the registration id
and the permission resource. Giving the pagebuilder's an id and a table of its
own was a breaking change, made long after the fact.

`people` (`people`) and `validators` (`validators`) predate the rule and are
the exception, not the pattern. Do not copy them.

## The registration id is not the table and not the URI

A model has three names, and they are not the same thing:

| | where it comes from | what it is used for |
| --- | --- | --- |
| id | snake_case of the plural label, or `presets.ModelWithID` | the menu key, the permission resource |
| URI name | `presets.ModelBuilder.URIName`, defaulting to the kebab of the id | the URL |
| table | `TableName()`, defaulting to gorm's plural of the type | the database |

A package sets the id explicitly, with the same prefix, instead of inheriting
whatever the type is called:

```go
mb := Model(b, &models.Organizacao{}, presets.ModelWithID(OrganizacaoModelID))
```

The id is unique across the whole admin: registering two models under one id is
an error, because the menu cannot show both and the permission tree cannot tell
them apart.

## A model that is not in the menu takes no key

A second builder over the same type — a dialog, an editor, a form — is not a
menu entry, and says so **at registration**, where the key is claimed:

```go
mb := b.Model(pm.NewModel(), presets.ModelNotInMenu()).
    URIName(pm.Info().URI() + "-version-list-dialog")
```

`InMenu(false)` says the same thing one call too late.

## Renaming carries the database

`AutoMigrate` does not rename: it creates the new table empty, adds the new
column beside the old one, and leaves every row where it was. A package that
renames a table or a field renames it in the database too, before AutoMigrate:

```go
func renameLegacy(db *gorm.DB) error {
	m := db.Migrator()
	if m.HasTable("qor_notes") && !m.HasTable("notes") {
		return m.RenameTable("qor_notes", "notes")
	}
	return nil
}
```

Guarded both ways — the old one exists, the new one does not — it is safe to
run on every boot, forever.

Remember what else moves with a rename: the registration id changes the menu
name the host writes and the permission resource any stored policy names.

## Menu group names are URL segments

A group is a path segment and a permission segment, so its name is URL-safe and
the human title is set apart:

```go
b.MenuGroup("page-builder").Title("Page Builder").Icon("mdi-view-quilt")
```

A name with a space reaches `http.ServeMux` as `/Page Builder/…`, which it
refuses to register — it reads the first word as a method.

See [the menu reference](admin/presets/docs/menu.md) for the tree itself.

## Checklist for a new package

- [ ] Every `TableName()` starts with the package prefix.
- [ ] Every model registers with an explicit `presets.ModelWithID`, prefixed.
- [ ] Models that are not menu entries pass `presets.ModelNotInMenu()`.
- [ ] Menu group names are URL-safe, with `Title` for display.
- [ ] `AutoMigrate` runs behind a `renameLegacy` for anything you renamed.
- [ ] The i18n module ships `Messages_en_US` and `Messages_pt_BR`, both
      registered (see the [README](README.md) conventions).
- [ ] `make check` passes.
