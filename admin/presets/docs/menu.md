# Menu

The side menu is a single tree owned by the `Builder`. Every model, page and
group is one entry of it, and every entry is reachable by a key.

```go
b.MenuOrder(
    presets.ModelItem("dashboards"),
    b.MenuGroup("content").Icon("mdi-file-document").Add(
        presets.ModelItem("posts"),
        presets.ModelItem("pages"),
        presets.PageItem("/import"),
    ),
)
```

## The key

An entry is identified by its TYPE and its NAME — `model:posts`, `page:/import`,
`group:content`. The type is half of the key, so a model, a page and a group may
share a name without colliding.

| Type | Reference | Name |
| --- | --- | --- |
| model | `presets.ModelItem(id)` | the model's registration id — see below |
| page | `presets.PageItem(path)` | the page's path, with the leading `/` |
| group | `presets.GroupItem(name)` | the group's name |

A model's id is NOT its URI name and not what `ModelBuilder.ID()` returns (that
one is the URI name). It is the snake_case of the model's plural label — of its
label for a singleton — or whatever `presets.ModelWithID("…")` set when the
model was registered:

```go
b.Model(&BlogPost{})                                  // model:blog_posts
b.Model(&BlogPost{}).URIName("articles")              // model:blog_posts still
b.Model(&BlogPost{}, presets.ModelWithID("articles")) // model:articles
```

A reference (`MenuRef`) names an entry without requiring it to exist. That is
what makes the order of configuration irrelevant: a group may name a model
registered later in the file, or in a plugin installed afterwards.

## Registration and reservation

Two things can happen to a key:

- **Registration** gives it its value — `b.Model(&Post{})` registers
  `model:posts`, `b.PagesRegistrator().AddHttpPage(page)` registers
  `page:/import`, `b.MenuGroup("content")` registers `group:content`.
- **Reference** puts the key in place with no value yet: a *reservation*. It
  holds the position and renders nothing until the value arrives.

A key is unique. Registering a value over a value is reported as an error
(`RegisterMenuItem`) — two models with the same id, a page added twice — instead
of one silently replacing the other. Registering the *same* value again is the
same registration, not a second one.

Only what is IN the menu takes a key. A model's id is not unique by
construction — a plugin may register a second builder over the same type, told
apart by `URIName`, as publish does for its version-list dialog and pagebuilder
for its editor. Such a model declares `presets.ModelNotInMenu()` at
registration, takes no key, and so cannot collide with the model that does want
the entry:

```go
b.Model(&Post{})                                   // model:posts
b.Model(&Post{}, presets.ModelNotInMenu()).        // no key at all
    URIName("dialog-select-favor-posts")
```

`InMenu(false)` says the same thing, but too late: the key is claimed when the
model is registered.

Filling a reservation keeps the position the reservation gave it, so this

```go
b.MenuOrder(b.MenuGroup("content").Add(presets.ModelItem("posts")))
b.Model(&Post{})
```

and this

```go
b.Model(&Post{})
b.MenuOrder(b.MenuGroup("content").Add(presets.ModelItem("posts")))
```

end in the same tree.

## Placing and moving

Everything is done through the key, registered or not: there is one entry per
key, and that entry is what moves.

| Call | What it does |
| --- | --- |
| `b.MenuOrder(items…)` | puts the items at the root, in the order given |
| `group.Add(items…)` | puts them inside the group, in the order given |
| `b.MoveMenuItem(ref, dst, index…)` | puts one entry inside `dst` (`nil` = the root), at `index` when given |
| `mb.SetMenuGroup(group)` | moves the model's entry |
| `page.MenuGroup(name)` | moves the page's entry (the name waits until the page is registered) |

An entry is in exactly one place at a time: moving it removes it from where it
was, never duplicates it. The last call wins. An out-of-range or negative
`index` appends.

A group cannot be moved into itself or into one of its own descendants — the
tree would stop being one — and such a move is refused without changing
anything.

## The group path is also the URL

A group is a path segment, and an entry carries the WHOLE chain of groups it
sits under, not only the innermost one:

```go
b.Model(&SEOConfig{}, presets.ModelWithID("seo_config")).URIName("seo_config")
b.MenuGroup("site").Add(b.MenuGroup("seo").Add(presets.ModelItem("seo_config")))
// menu:       Site > SEO > Seo Config
// URL:        /admin/site/seo/seo_config
// permission: …:site:seo:seo_config
```

`group.Path()` is that chain as `"site/seo"`, `group.PathNames()` its segments,
and `mb.MenuGroupName()` the same for the group a model sits in. The model's
`URI()`, the page's `FullPath()`, the breadcrumb and the automatic path-based
permission all read from the tree, so moving an entry moves its URL with it.

Nothing is cached on the model or on the page: a copy would go stale at the
first move.

## Title and icon

A group's title defaults to a humanized version of its name; `Title` /
`TitleFunc` override it, and `Icon` sets its icon. A model's icon comes from
`MenuIcon`, and falls back to one guessed from its label. Icon names are the
Material Design ones (<https://fonts.google.com/icons>).

```go
b.MenuGroup("content").Title("Conteúdo").Icon("mdi-file-document")
b.Model(&Post{}).MenuIcon("mdi-post")
```

`InMenu(false)`, on a model or on a page, keeps it out of the menu while it
keeps its routes.

## Pages of a model are not menu entries

A page registered on a model — `mb.Listing().PagesRegistrator()` or
`mb.Detailing().PagesRegistrator()` — is a CHILD of the listing or of the
record, not an item of the side menu. It takes no key in the tree: it is mounted
under the model's route and appears in the model's own page menu (the listing's
action menu, or the record menu on the detail page). Two models are therefore
free to have a page of the same path.

Only pages of the builder itself — `b.PagesRegistrator()` — are side-menu
entries.

A **singleton** has one record and no listing, so the menu it shows is the
record's. A page registered on the listing of a singleton would be mounted
nowhere and appear in no menu, so the boot refuses it and says to register it on
`Detailing()` instead.

## Reading the tree

| Call | Returns |
| --- | --- |
| `b.MenuTree()` | the root sentinel — it has no entry of its own and never renders |
| `b.MenuItems()` | the key registry, `map[string]*MenuItem` |
| `b.MenuItemOf(ref)` | the entry a key names, creating the reservation if needed |
| `group.Items()` | the group's entries, in order |
| `item.Registered()` | whether the value has arrived |
| `item.Parent()` | the group the entry sits in, `nil` at the root |

## The older string API

`SubItems` and bare strings still work: a string starting with `/` is a page,
anything else a model.

```go
b.MenuOrder("books", b.MenuGroup("Media").SubItems("videos", "musics"))
```

A bare string cannot say *the group named x*, which is what its ambiguity costs.
Prefer `ModelItem` / `PageItem` / `GroupItem` with `Add`.
