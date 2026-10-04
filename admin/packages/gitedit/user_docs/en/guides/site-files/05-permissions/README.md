# Permissions

Each action of the page is a permission of it, given to the roles in **User
Roles** (see the policies and permissions). The page is
{%= admin.page("/site-files").link %}, in the menu at
**{%= admin.page("/site-files").menu %}**.

## The page's actions

| Action | Permission |
| --- | --- |
| See the page and the files (read only) | `@get` |
| Commit | `!commit` |
| Update | `!update` |
| Publish | `!publish` |
| Discard | `!discard` |
| Start over | `!reset` |
| Preview | `!preview` |

So publishing can stay with fewer people than editing: whoever edits commits,
and someone else, who reviews, publishes.

## The files

Each way of changing the files of the draft is a permission of its own — in
the editor and in the **WebDAV** alike. Whoever has only `@get` sees the files
and changes none: the editor opens them read only, and the buttons of what the
user may not do are not shown.

| Action | Permission | In the WebDAV |
| --- | --- | --- |
| Create a file or a folder that is not there | `!create` | PUT of a new file, MKCOL, COPY |
| Change a file that is there | `!edit` | PUT of a file there; locks, properties |
| Rename a file in its folder | `!rename` | MOVE in the same folder |
| Move a file to another folder | `!move` | MOVE to another folder |
| Delete a file or a folder | `!delete` | DELETE |
| Import: upload from the computer, download from the internet | `!import` | — |

Importing asks `!import` **and** what it writes: `!create` for each new file,
`!edit` for each one that is there (overwritten). The download from the
internet is made by the server, and only from public addresses: never from
the server itself nor from its network.

## Permissions of a path

Any of the permissions of the files can be given or taken **for a path and
everything under it**. The path goes in the resource, between `<` and `>`,
before the action; `*` is any text, so `<static/*>` is all of `static/`, in
every subfolder. The resource of this page is
`{%= admin.page("/site-files").resource %}`:

| Policy | Effect |
| --- | --- |
| Deny `{%= admin.page("/site-files").resource %}<config/*>:!edit` | no one of the role edits under `config/` |
| Allow `{%= admin.page("/site-files").resource %}<static/*>:*` | the role does anything under `static/` |
| Allow `{%= admin.page("/site-files").resource %}<static/img/*>:!import` | the role imports images into `static/img/` |

How it decides, in order:

1. a **deny of the path** denies;
2. an **allow of the page** (`…:!edit`) allows;
3. a **deny of the page** denies — whatever the paths say;
4. with nothing said of the page, an **allow of the path** allows.

So, for a role that may only change the styles: give `@get` on the page and
`<static/css/*>:!edit` (and `!create`, if it creates files there), and nothing
else of the files.

## Examples by path

Each line is a policy of a role: its effect and its permission, whole — the
scope and the page as the admin makes them now (they change with the menu, and
so do these lines). An action is `!create`, `!edit`, `!rename`, `!move`,
`!delete` or `!import`; `!*` is all of them, and `*` adds seeing (`@get`).

**Only the styles: changes `static/css/` and nothing else**

| Effect | Permission |
| --- | --- |
| Allow | `{%= admin.page("/site-files").resource %}@get` |
| Allow | `{%= admin.page("/site-files").resource %}<static/css/*>:!edit` |
| Allow | `{%= admin.page("/site-files").resource %}<static/css/*>:!create` |

**Everything, except `config/`, which is read only**

| Effect | Permission |
| --- | --- |
| Allow | `{%= admin.page("/site-files").resource %}*` |
| Deny | `{%= admin.page("/site-files").resource %}<config/*>:!*` |

**Images: imports into `static/img/` only**

| Effect | Permission |
| --- | --- |
| Allow | `{%= admin.page("/site-files").resource %}@get` |
| Allow | `{%= admin.page("/site-files").resource %}<static/img/*>:!import` |
| Allow | `{%= admin.page("/site-files").resource %}<static/img/*>:!create` |
| Allow | `{%= admin.page("/site-files").resource %}<static/img/*>:!edit` |

**The layouts are never deleted nor moved, whatever else the role may**

| Effect | Permission |
| --- | --- |
| Deny | `{%= admin.page("/site-files").resource %}<templates/layouts/*>:!delete` |
| Deny | `{%= admin.page("/site-files").resource %}<templates/layouts/*>:!move` |
| Deny | `{%= admin.page("/site-files").resource %}<templates/layouts/*>:!rename` |

**Rename and move only inside `static/`** — a move asks the permission on
both paths, so out of `static/` (or into it) is denied:

| Effect | Permission |
| --- | --- |
| Allow | `{%= admin.page("/site-files").resource %}<static/*>:!rename` |
| Allow | `{%= admin.page("/site-files").resource %}<static/*>:!move` |

**A single file**: the path whole, no `*` — `{%= admin.page("/site-files").resource %}<config/layout_config.gad>:!edit`.
