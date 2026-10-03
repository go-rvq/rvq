# 1. Permissions

What a user may do comes from their **roles**: each role allows or denies,
part by part of the admin, to list, see, create, edit and delete records, and
to run actions.

- **What is not allowed is not shown**: an item of the side menu, a button, an
  action, a tab, a column. A group with nothing allowed is not in the menu.
- **Fields and sections** may be allowed to see but not to edit — they are
  shown read only —, or not at all — they are not shown.
- **The server checks too**: an address typed by hand, or a button of an old
  page, is refused when the role does not allow it.
- The **trash** of a part is seen only by the roles allowed to see it
  ({%= admin.doc("guides/policies/02-trash").link %}).

## How a permission is written

A permission names a **resource** — its parts separated by `:` — and ends in
what is asked of it:

| Part | Example |
| --- | --- |
| a group of the menu: its name and `/` | `site/` |
| a part (model) | `seo_config` |
| a record: its id in `<…>` | `<7>` — `<*>` any |
| a field | `#Title` |
| a section of the detail | `$Main` |
| a page | `/report` |
| a permission: `@` and its name | `@list`, `@get`, `@create`, `@edit`, `@delete` |
| an action: `!` and its name | `!publish` |

A part is reached two ways — **through its groups** and **by its unique
name**, its own:

    presets:site/:seo/:seo_config:<7>:@edit     through the groups
    presets:seo_config:<7>:@edit                by the unique name

`*` stands for anything: `presets:site/:*` is everything inside the group
*site*; `presets:seo_config:*` everything of the part *seo_config*, wherever
the menu puts it.

## Which one decides

1. What is given by the **unique name** decides first: a **deny** — of any of
   the user's roles — denies; an allow allows.
2. Only when nothing is given by the unique name do the **groups** decide: an
   allow of any role allows.

So a group may be allowed while one part of it is denied by its name, or a
group denied while one part of it is allowed by its name.

## The permissions of the admin

Every part of the admin and what may be asked of it — open a node to see what
is inside it. Each part of the documentation also has its own
**Permissions** item.

{%= admin.permissionsTree() %}

Ask whoever administers the roles for what your work needs.
