# The trash: the `TRASH` mode

A model whose records are deleted softly keeps them in a **trash**: listed by a
place of its listing (a tab), each one opened in its detail. In the trash a
record is **read only**: it is not edited nor deleted again, and none of its
actions runs but those made for the trash. `presets` knows nothing of how
records are deleted; whoever keeps them (`perms.SetupTrash`, say) tells it
where the trash is, and `presets` does the rest.

## The mode

`TRASH` is a `FieldMode` bit that **qualifies** `LIST` and `DETAIL`:

| Where | Mode |
|---|---|
| a listing of the trash | `LIST \| TRASH` |
| the detail of a deleted record | `DETAIL \| TRASH` |

`Is` compares the whole mode (`(LIST|TRASH).Is(LIST)` is false); test a part
of it with `Has`, or the helpers made of it: `IsList()`, `IsDetail()` — true
in the trash too — and `IsTrash()`. A field component that draws a listing's
or a detail's cell keeps working in the trash; one that offers an action
hides it when `field.Mode.IsTrash()`.

The mode of the request is in its context: `presets.WithFieldMode(ctx, m)`
sets it, `presets.GetFieldMode(ctx)` reads it, `presets.InTrash(ctx)` says it
is of the trash. The listing sets `LIST` or `LIST | TRASH` (a listing in a
detail keeps the detail's in the context); the detail, `DETAIL` or
`DETAIL | TRASH`. The fields get it in `FieldContext.Mode`.

## Saying where the trash is

```go
// the listing of a request is of the trash
mb.Listing().SetInTrashFunc(func(ctx *web.EventContext) bool {
    return ctx.R.URL.Query().Get(presets.ActiveFilterTabQueryKey) == "trash"
})
// a record is deleted
mb.SetDeletedFunc(func(obj any) bool { return obj.(*Post).DeletedAt.Valid })
```

`ListingBuilder.Mode(ctx)` and `ModelBuilder.IsDeleted(obj)` answer from
them.

## Where an action is available: `TrashPolicy`

An action — `ActionBuilder`, `BulkActionBuilder`, `RowMenuItemBuilder` — says
where it is available with `SetTrash`:

| Policy | Out of the trash | In the trash |
|---|---|---|
| `TrashOut` (the default) | yes | no |
| `TrashToo` | yes | yes |
| `TrashOnly` | no | yes |

```go
lb.BulkAction("restore").SetTrash(presets.TrashOnly)  // brings records back
d.Action("deleted_origin").SetTrash(presets.TrashOnly) // where it was deleted from
d.Action("activityLog").SetTrash(presets.TrashToo)     // reading is harmless
```

## What the trash offers

In a listing of the trash:

- the bar has only the bulk actions and the actions available there — none of
  the buttons added to it (`PrependListButtons`, `AppendListButtons`), its
  pages, its new button; out of the trash, a `TrashOnly` action is not there;
- a row has its checkbox only while a bulk action is available;
- the menu of a row (**⋯**) keeps the models nested in the record and its
  pages, and the actions available there (an item that opens a nested model is
  there whatever its policy).

In the detail of a deleted record:

- no edit button, no edit of a section;
- its menu, as a row's in the trash.

## What is refused

The server refuses what the trash does not allow, whatever the page shows:

- editing a deleted record — the edit form, its save, a section's save —:
  the model's `EditingRestriction` refuses it, and every builder's inherits it;
- deleting it again: the `DeletingRestriction`, likewise (the detail's and the
  listing's inherit the model's deleting restriction);
- running an action not available where it is asked for: `ActionBuilder.Do`
  and `View` (on a deleted record fetched by its id, or in a listing of the
  trash) and `BulkActionBuilder.Do` and `View` answer `ErrActionNotAllowed`.

A package with actions of its own outside these builders refuses a deleted
record itself (`publish` does: publishing, unpublishing, scheduling, renaming
a version of a deleted record answer `ErrUpdateRecordNotAllowed`, and its bar
draws no button of them in the trash).
