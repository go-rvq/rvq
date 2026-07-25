# Form host — opening and destroying edit/create overlays

An edit or create form is an **overlay** (dialog or drawer) that some page opens:
a listing opens *New* and the per-row *Edit*/*Detail*; a detailing opens *Edit*.

The page is the **host** of that overlay. It declares one reactive variable that
*is* the overlay's `closer`, and guards the overlay's block with it:

```html
<user-component :assign='[[vars,{...{"$presetsEditing": {show:false}}}]]'>

  … the page (its Edit button only does `vars.$presetsEditing.show = true`) …

  <go-plaid-scope v-if='vars.$presetsEditing?.show' v-slot='{ form }' :form='[{}]'>
    <go-plaid-portal portal-name='<uid>' :scope='{"closer": vars.$presetsEditing}'/>
    <go-plaid-run-script :script='plaid()…eventFunc("presets_Edit")
        .query("target_portal","<uid>")
        .scope({closer: vars.$presetsEditing})
        .query("presets_closer_provided","true").go()'/>
  </go-plaid-scope>

</user-component>
```

Everything follows from that one variable:

| `VAR.show` | Effect |
| --- | --- |
| `= true` | the guarded block mounts → its run-script loads the form into the host portal → the overlay appears. **One request.** |
| `= false` | the block unmounts → the overlay, the form and its `form` scope are destroyed. No state survives. |
| `= true` again | the block mounts again → a **fresh** form is loaded. |

Built by [`FormHost`](../form_host.go). The scope variables:

| Constant | Variable | Hosted by | Opens |
| --- | --- | --- | --- |
| `DetailingEditScope` | `vars.$presetsEditing` | detailing (event and page) | the edit form of the shown record |
| `ListingNewScope` | `vars.$presetsCreating` | listing | the create form |
| `ListingItemEditScope` | `vars.$presetsItemEditing` | listing | the edit form of a row (`.id`) |
| `ListingItemDetailScope` | `vars.$presetsItemDetailing` | listing | the detailing of a row (`.id`) |

Two design points that are easy to get wrong:

- **The state lives on `vars`, not on a slot scope.** On a page the primary
  action (the Edit button) is rendered by the layout into the app bar — *outside*
  the host component — where a slot-scoped variable does not exist. `vars` is the
  app-wide reactive object, reachable from anywhere. The `$presets` prefix keeps
  these out of the application's own namespace.
- **The overlay must not create its own closer.** The host seeds `closer` in the
  portal's content scope *and* sends `presets_closer_provided=true`; the
  `Dialog`/`Drawer` responder then binds to that closer instead of wrapping the
  content in a new one (`CloserScope`). Without this, `closer.show = false` on
  save would close a *child* closer, leaving the host on — the form would never
  be destroyed and the button would not re-open it.

## Host lifecycle

```mermaid
stateDiagram-v2
    direction LR
    [*] --> Closed : host renders<br/>(VAR = {show:false})

    Closed --> Loading : button sets VAR.show = true<br/>(guarded block mounts)
    Loading --> Open : run-script → event → overlay<br/>rendered into host portal
    Open --> Closed : VAR.show = false<br/>(block unmounts, form destroyed)
    Closed --> Loading : re-open → fresh form

    state Open {
        [*] --> Idle
        Idle --> Idle : list editor add/remove,<br/>validation re-render<br/>(inner portal only)
        Idle --> Saved : save succeeds
    }

    Saved --> Closed : runScript "closer.show = false"<br/>+ post-change callback
```

`closer` inside the overlay **is** `VAR`: the dialog binds `v-model='closer.show'`
and the save response runs `closer.show = false`, which is what unmounts the
block. Closing by the ✕ button does the same through the same binding.

## NEW, from a listing

The listing renders the `$presetsCreating` host around itself
(`ListingComponentBuilder.hostForms`); the New button carries no plaid.

```mermaid
stateDiagram-v2
    direction TB
    [*] --> ListingShown : listing rendered<br/>vars.$presetsCreating = {show:false}

    ListingShown --> Loading : click New<br/>vars.$presetsCreating.show = true
    Loading --> FormOpen : GET presets_New<br/>(target_portal=host, overlay=Dialog|RightDrawer,<br/>closer=vars.$presetsCreating, closer_provided=true)

    state FormOpen {
        [*] --> Editing
        Editing --> Invalid : POST presets_Create → validation fails
        Invalid --> Editing : form re-rendered in the inner portal<br/>(same `form` scope, errors shown)
        Editing --> Editing : list editor add row (POST listEditor_addRowEvent)<br/>remove row (client-side only)
    }

    FormOpen --> Closed : POST presets_Create OK →<br/>runScript "closer.show = false" + reload callback
    FormOpen --> Closed : ✕ / Esc → closer.show = false
    Closed --> Loading : click New again → fresh form
    Closed --> [*]
```

The `post_change_callback` given to the host is the listing's reload callback, so
a successful create refreshes the table in the same round-trip that closes the
overlay.

## DETAIL / EDIT, from a listing row

The listing renders **one** host per action for the whole table — not one per row.
Each host has an extra `id` variable; the row points it at its record:

```html
<td @click.self='vars.$presetsItemDetailing.id = "42";
                 vars.$presetsItemDetailing.show = true; …'>
```

and the host's run-script loads whatever the variable holds:

```js
plaid()…eventFunc("presets_Detailing").query("id", vars.$presetsItemDetailing.id)…
```

Which host a row uses: `$presetsItemDetailing` when the model has a detailing,
otherwise `$presetsItemEditing` (the row menu's *Edit* always uses the latter).
Outside a listing — a row menu rendered elsewhere — there is no host and the call
site falls back to the self-hosting `presets_EditForm` event (see below).

```mermaid
stateDiagram-v2
    direction TB
    [*] --> ListingShown : hosts rendered:<br/>$presetsItemDetailing / $presetsItemEditing = {show:false, id:null}

    ListingShown --> Loading : click row (or row menu Edit)<br/>VAR.id = "<id>"; VAR.show = true
    Loading --> OverlayOpen : GET presets_Detailing | presets_Edit<br/>query id = VAR.id

    OverlayOpen --> ListingShown : VAR.show = false<br/>(overlay destroyed)
    ListingShown --> Loading : click another row → same host, new id

    state OverlayOpen {
        [*] --> Detail
        Detail --> EditOpen : the DETAILING hosts its own edit form<br/>(see next section)
    }
```

Reusing one host per action means the DOM holds a single overlay portal no matter
how many rows the table has, and opening a second row simply reloads it.

## EDIT, from a detailing

Both detailing entry points host the edit form the same way — the event (overlay)
and the page (`defaultPageFunc`), through `DetailingBuilder.hostedComponent`:

```mermaid
stateDiagram-v2
    direction TB
    [*] --> DetailShown : detailing rendered<br/>vars.$presetsEditing = {show:false}

    DetailShown --> Loading : click Edit (app bar)<br/>vars.$presetsEditing.show = true
    Loading --> EditOpen : GET presets_Edit (id fixed = shown record)<br/>overlay = Dialog on an overlayed detailing,<br/>RightDrawer on a page

    state EditOpen {
        [*] --> Editing
        Editing --> Invalid : POST presets_Update → validation fails
        Invalid --> Editing : re-render in the inner portal
        Editing --> Editing : list editor add/remove
    }

    EditOpen --> DetailShown : save OK → "closer.show = false"<br/>+ post-change callback reloads the detailing
    EditOpen --> DetailShown : ✕ → closer.show = false
    DetailShown --> Loading : Edit again → fresh form
```

Because the button only flips a variable, **any** other button on the detailing
can open the same form with `vars.$presetsEditing.show = true`.

## Round-trip in detail

```mermaid
sequenceDiagram
    participant U as User
    participant P as Page (host)
    participant S as Server
    participant O as Host portal

    U->>P: click (New / row / Edit)
    P->>P: vars.$presetsX.show = true<br/>(+ .id for row hosts)
    Note over P: v-if mounts the block:<br/>portal + run-script + child `form` scope
    P->>S: presets_New | presets_Edit | presets_Detailing<br/>target_portal, overlay, closer_provided=true
    S-->>O: updatePortal → dialog/drawer bound to v-model='closer.show'
    Note over O: `closer` is vars.$presetsX<br/>(seeded by the portal's :scope)

    U->>O: edit fields, add/remove list-editor rows
    U->>S: POST presets_Update | presets_Create
    alt validation fails
        S-->>O: updatePortal → same form re-rendered with errors
        Note over O: the `form` scope is NOT recreated
    else success
        S-->>P: runScript "closer.show = false" + post-change callback
        Note over P: v-if unmounts → form destroyed
    end
```

## Re-renders inside an open overlay

While the overlay is open, the server re-renders **only the inner portal**. The
host's `form` scope is created once, by the host, so it survives every re-render —
which is what keeps the reactive form intact:

| Interaction | Request | What is re-rendered |
| --- | --- | --- |
| Validation failure on save | `presets_Update` / `presets_Create` | the whole form, with `Errors` on the offending fields; the reactive `form` keeps the submitted values |
| List editor **add row** | `listEditor_addRowEvent` | the form; the appended row is flagged `__new` and re-seeded by the component setup |
| List editor **remove row** | *(none)* | client-side only: `form["<field>[i].__deleted"] = true`; the row renders as removed (or vanishes, if it was `__new`) and is reconciled on the next submit |
| List editor **sort** | `listEditor_sortEvent` | the form, with the new `__pos` order |

See [list editor & tables](list-editor.md) for the per-item metadata
(`__pos`, `__index`, `__new`, `__deleted`, `__present`).

That is why the `form` scope belongs to the host and not to the form: an
add/remove/validation round-trip must not reset the fields the user already
filled in.

## Your own buttons

Anything rendered **inside the host** — and on a page that includes the app bar,
because the state lives on `vars` — can drive the overlay by writing to the
variable. No plaid, no portal, no event: just `VAR.show = true|false`.

### Open the edit form from an extra button on a detailing

```go
d := mb.Detailing("Name", "Total")

d.Field("Total").ComponentFunc(func(field *presets.FieldContext, ctx *web.EventContext) h.HTMLComponent {
    return h.Div(
        h.Text(fmt.Sprint(field.Value())),
        // a second way into the same edit form
        VBtn("Corrigir total").
            Variant(VariantText).
            Attr("@click", "vars."+presets.DetailingEditScope+".show = true"),
    )
})
```

### Close it from inside the form

The form is rendered with `closer` bound to the host's variable, so inside the
overlay both spellings work:

```go
// inside an editing field / action of the same model
VBtn("Cancelar").Attr("@click", "closer.show = false")
// … identical to:
VBtn("Cancelar").Attr("@click", "vars."+presets.DetailingEditScope+".show = false")
```

### Open a row's edit form from a custom listing action

Per-item hosts also take the record id. Use the published hosts so the expression
matches whatever the listing rendered (and degrade gracefully when there is none):

```go
l.RowMenu().SetRowMenuItem("Editar total").ComponentFunc(
    func(rctx *presets.RecordMenuItemContext) h.HTMLComponent {
        open := presets.GetItemFormHosts(rctx.Ctx).OpenEditExpr(rctx.ID)
        if open == "" { // rendered outside a listing: no host to drive
            open = web.Plaid().
                EventFunc(actions.EditForm).
                Query(presets.ParamID, rctx.ID).
                Go()
        }

        return VListItem(VListItemTitle(h.Text("Editar total"))).Attr("@click", open)
    })
```

`OpenEditExpr("42")` produces
`vars.$presetsItemEditing.id = "42"; vars.$presetsItemEditing.show = true`.

### Open the create form from anywhere in a listing

```go
l.NewButtonFunc(func(ctx *web.EventContext) h.HTMLComponent {
    return VBtn("Novo produto").
        Color("primary").
        Attr("@click", "vars."+presets.ListingNewScope+".show = true")
})
```

### Toggle, or react to the state

The variable is plain reactive state, so it can also be read:

```go
// toggle
VBtn("").Attr("@click", "vars.$presetsEditing.show = !vars.$presetsEditing.show")

// disable something while the form is open
VBtn("Excluir").Attr(":disabled", "vars.$presetsEditing?.show")
```

Use `?.` when reading in a place that may render before the host assigns its
state; writing (`… .show = true`) happens on click, when the state already exists.

## Callers that cannot host the form

A row menu rendered outside a listing, an action button somewhere else, and other
subsystems (publish, l10n, model_select, pagebuilder) cannot render a host around
themselves. They use the **self-hosting events** `presets_EditForm` /
`presets_NewForm`, which respond with a `FormHost` that is already on
(`Show(true)`): it mounts, loads the form immediately, and behaves exactly like a
page-hosted one from then on. The only difference is one extra round-trip — the
wrapper response — which is precisely what page hosts avoid.

## API

```go
// A page hosts a form:
comp := presets.FormHost(presets.DetailingEditScope, portal,
        web.Plaid().URL(url).EventFunc(actions.Edit).Query(presets.ParamID, id)…).
    Children(pageContent).
    Component()

// A button opens it:
btn.Attr("@click", "vars."+presets.DetailingEditScope+".show = true")
```

| Method | Purpose |
| --- | --- |
| `Show(bool)` | initial state; `true` makes the host self-opening |
| `Var(name, init)` | extra reactive variable (e.g. `id`), readable by the load event |
| `Ref()` | the JS expression for the host state (`vars.$presetsEditing`) |
| `ShowExpr()` / `OpenExpr(vars)` | what a button runs to open it |
| `Children(...)` | the page content rendered inside the host |

Listings publish their per-item hosts on the context so rows can find them:
`WithItemFormHosts` / `GetItemFormHosts`, whose `OpenEditExpr(id)` /
`OpenDetailExpr(id)` return the expression for a row (or `""` when there is no
host, so the caller can fall back).

Related: [`ParamCloserProvided`](../const.go), `DialogBuilder.SetCloserProvided`,
`Drawer.SetCloserProvided`.
