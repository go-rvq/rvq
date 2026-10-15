# vx-portal

A place in the page whose content is a Vue template given at run time. The
template is compiled and shown there, and compiled again each time it changes.

It belongs to vuetifyx (`Portal.vue`). It is not corejs's `go-plaid-portal`
(Go `web.Portal`), which is a different component: that one is named, loaded
and updated by the events of plaid (`r.UpdatePortal`, `reloadPortals`). In
Go, use `vuetifyx.VXPortal` (`x/ui/vuetifyx/portal.go`).

## Props

| Prop | What it is |
|---|---|
| `v-model` (`modelValue`) | The content (see below). |
| `v-model:loading` | When `true`, the slot `loading` is shown in place of the content. Setting a content makes it `false`. |
| `data` | The variables of the content's templates, which are also the props of the default slot. |
| `tag` | The element the content is put in, e.g. `div`. It takes `id`, `class`, `style` and `attrs`. Without it, the portal has no element of its own. |
| `name` | Registers the portal in `app.config.globalProperties.vx$portals[name]`. |
| `id`, `class`, `style`, `attrs` | Applied to the `tag` element. |

## The content

| Value | What is shown |
|---|---|
| `"<p>{{ x }}</p>"` | That template. |
| `{template, options: {data, setup, mounted, unmounted, resetScroll}}` | The template. `options.data` adds variables to `data`'s. The callbacks receive `{data, locals}` (`component-by-template.ts`). |
| `""` or `{template: ""}` | Nothing. |
| `null` or `undefined` | The initial content again. |

**The initial content** is the model's when it is a template; otherwise it
is the default slot (or nothing, when there is no slot). That is what the
portal shows once mounted.

## Scope: what a template sees

A template of the content sees only **its own variables**: those of `data`,
those of `options.data`, `Vue`, `window`, `vx$scope` (all of them) and
`vx$portal` (the portal). It does **not** see the scope of the page where the
portal is (`vars`, `locals`, `form`, `plaid`…).

There are two ways to give it that scope:

1. **`data`**: pass the variables on by name, e.g.
   `:data="{vars, locals}"`.
2. **The scoped default slot**: write the slot in the page. It sees the whole
   scope of the page, and its props are the portal's variables (`data`):

   ```html
   <vx-portal tag="div" v-model="content" :data="{given}">
     <template #default="{given}">
       <span>{{ pageVar }} / {{ given }}</span>
     </template>
   </vx-portal>
   ```

   The slot is the initial content. **It also stays the default slot of any
   content given**: a template of the content shows it wherever it says
   `<slot></slot>`. That is how dynamic content and the page's scope are
   combined:

   ```js
   content = '<h3>{{ given }}</h3><slot></slot>'
   ```

   Here `<h3>` sees `given`; the slot, written in the page, sees all of the
   page.

A default slot wrapped in `<template v-pre>` (Go: `DefaultSlot`) is instead
taken as a template of the portal. It then sees only the portal's variables,
as a content given would.

## Setting the content from an event (plaid)

Ask for the response's body with `loadPortalBody(true)` (otherwise plaid would
take the body as the page's) and put it in the model's expression:

```js
plaid().vars(vars).locals(locals).eventFunc("my_Event").query("x", 1)
  .loadPortalBody(true).go()
  .then(r => { vars.myContent = r.body })
```

In Go:

```go
vx.VXPortal().
  Tag("div").
  Model("vars.myContent").
  Data("{vars, locals}").
  Slot("{}", h.Text("…"))
```

## Events

Each event receives an `UpdateEvent` (`types.ts`):
`{portal, contentArg, content?, data?, locals?, added?, removed?}`. Its
`portal.elements` is what the portal shows at that moment.

| Event | When |
|---|---|
| `preContentUpdate` | A content is set, before anything changes. |
| `setup` | The content's component is created. |
| `mounted` | The content is in the page. `portal.elements` are its root elements. |
| `postContentUpdate`, `change` | The content set is shown (after `mounted`). When the content is emptied (`""`), they come after the old content's `unmounted`, with no element. |
| `unmounted` | A content goes out (replaced or emptied). `portal.elements` are still its own. |
| `elements` | The root elements changed while the content is shown, e.g. a `v-if` or a `v-for` at its root. `added` and `removed` are the nodes that came in and went out, and `portal.elements` the new ones. A change inside a root element is not this event. |

When a content replaces another, the order is: `preContentUpdate`, `setup`
(the new one), `unmounted` (the old one), `mounted`, `postContentUpdate`,
`change`.

## The elements

`portal.elements` are the **root nodes of the content**:

- its element;
- with several roots, all of them in order, including the items of a `v-for`
  and the roots of a component at the root;
- texts with something written.

It does **not** include:

- comments, such as the place a `v-if` leaves;
- the empty anchors of a fragment;
- anything of the page around the portal.

They are kept up to date by the `elements` event.

## The portal

Inside the content, the portal is `vx$portal` (or `usePortal()`), and it is
also the component's exposed object:
`{name, content (the model's ref), elements, containerElement}`. With `name`,
it is also `vx$portals[name]`, so `vx$portals[name].content.value = "…"`
changes the content.

Tests: `src/lib/__tests__/Portal.spec.ts` (scope, slot) and
`PortalEvents.spec.ts` (events and elements, with and without `tag`).
