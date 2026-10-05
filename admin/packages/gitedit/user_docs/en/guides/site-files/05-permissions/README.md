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
| Share your draft (and revoke it) | `!share` |
| See and reach any user's draft, revoke any sharing | `!drafts` |

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
| Clone and push by git (see the step **By git**) | `!git` | — |

Importing asks `!import` **and** what it writes: `!create` for each new file,
`!edit` for each one that is there (overwritten). The download from the
internet is made by the server, and only from public addresses: never from
the server itself nor from its network.

## Permissions of a path

Any of the permissions of the files can be given or taken **for a path and
everything under it**. The path goes in the resource, between `<` and `>`,
before the action; `*` is any text, so `<static/*>` is all of `static/`, in
every subfolder.

The page's permissions are declared two ways — **with its group**, where the
menu puts it now, and **out of the group**, by its unique name:

| Declared | Resource |
| --- | --- |
| with the group | `{%= admin.page("/site-files").resource %}` |
| out of the group | `{%= admin.page("/site-files").uniqueResource %}` |

The one **out of the group has priority**: with both, it decides. And it keeps
holding if the menu moves the page to another group. So the policies below
are declared out of the group:

| Policy | Effect |
| --- | --- |
| Deny `{%= admin.page("/site-files").uniqueResource %}<config/*>:!edit` | no one of the role edits under `config/` |
| Allow `{%= admin.page("/site-files").uniqueResource %}<static/*>:*` | the role does anything under `static/` |
| Allow `{%= admin.page("/site-files").uniqueResource %}<static/img/*>:!import` | the role imports images into `static/img/` |

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
scope and the page, out of the group (the same wherever the menu puts the
page). With the group, the same lines start by `{%= admin.page("/site-files").resource %}`. An action is `!create`, `!edit`, `!rename`, `!move`,
`!delete` or `!import`; `!*` is all of them, and `*` adds seeing (`@get`).

**Only the styles: changes `static/css/` and nothing else**

| Effect | Permission |
| --- | --- |
| Allow | `{%= admin.page("/site-files").uniqueResource %}@get` |
| Allow | `{%= admin.page("/site-files").uniqueResource %}<static/css/*>:!edit` |
| Allow | `{%= admin.page("/site-files").uniqueResource %}<static/css/*>:!create` |

**Everything, except `config/`, which is read only**

| Effect | Permission |
| --- | --- |
| Allow | `{%= admin.page("/site-files").uniqueResource %}*` |
| Deny | `{%= admin.page("/site-files").uniqueResource %}<config/*>:!*` |

**Images: imports into `static/img/` only**

| Effect | Permission |
| --- | --- |
| Allow | `{%= admin.page("/site-files").uniqueResource %}@get` |
| Allow | `{%= admin.page("/site-files").uniqueResource %}<static/img/*>:!import` |
| Allow | `{%= admin.page("/site-files").uniqueResource %}<static/img/*>:!create` |
| Allow | `{%= admin.page("/site-files").uniqueResource %}<static/img/*>:!edit` |

**The layouts are never deleted nor moved, whatever else the role may**

| Effect | Permission |
| --- | --- |
| Deny | `{%= admin.page("/site-files").uniqueResource %}<templates/layouts/*>:!delete` |
| Deny | `{%= admin.page("/site-files").uniqueResource %}<templates/layouts/*>:!move` |
| Deny | `{%= admin.page("/site-files").uniqueResource %}<templates/layouts/*>:!rename` |

**Rename and move only inside `static/`** — a move asks the permission on
both paths, so out of `static/` (or into it) is denied:

| Effect | Permission |
| --- | --- |
| Allow | `{%= admin.page("/site-files").uniqueResource %}<static/*>:!rename` |
| Allow | `{%= admin.page("/site-files").uniqueResource %}<static/*>:!move` |

**A single file**: the path whole, no `*` — `{%= admin.page("/site-files").uniqueResource %}<config/layout_config.gad>:!edit`.
