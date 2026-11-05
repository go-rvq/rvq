# 1. Permissions

What a user may do comes from their **roles**, in order: each role allows or
denies, part by part of the admin, to list, see, create, edit and delete
records, and to run actions (see *The stack of permissions*, below).

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
| the scope: who asks — first, always; the admin's is `admin` (see below) | `admin:…` |
| a group of the menu: its name and `/` | `site/` |
| a part (model) | `seo_config` |
| a record: its id in `<…>` | `<7>` — `<*>` any |
| a field | `#Title` |
| a record edited in place, in a field: `&` and the field; its fields under it | `&Config`, `&Config:#Title` |
| a section of the detail | `$Main` |
| a page | `/report` |
| what is asked of a page: `@` and its HTTP method, one for each the page answers | `/report:@get`, `/report:@post` |
| a permission: `@` and its name | `@list`, `@get`, `@create`, `@edit`, `@delete` |
| an action: `!` and its name | `!publish` |

A part is reached two ways — **through its groups** and **by its unique
name**, its own:

    admin:site/:seo/:seo_config:<7>:@edit     through the groups
    admin:seo_config:<7>:@edit                by the unique name

A **page** too: its unique name is its path, so a permission of it may be
declared **out of its group** — with no group at all:

    admin:site/:/site-files:!edit             through the group site
    admin:/site-files:!edit                   by the unique name (no group)

The unique name goes on, part by part, to whatever is under it: a record, a
field, an action, a path (`admin:/site-files:<static/*>:!edit`).

`*` stands for anything: `admin:site/:*` is everything inside the group
*site*; `admin:seo_config:*` everything of the part *seo_config*, wherever
the menu puts it.

## The scopes

The first part of a resource is its **scope**: the part of the system that
asks for the permission. The admin's is **`admin`** — so its resources
begin with `admin:` (`admin:content/:posts:<12>:@edit`) —, and it holds
everything of the menu: the groups, the parts, their records, fields,
sections, actions and pages. The media library and the jobs are parts of it
too:

| What | Resource |
| --- | --- |
| sending a file, in a field of pictures or files | `admin:media_libraries:@create` |
| deleting a file / editing its description | `admin:media_libraries:<5>:@delete` / `admin:media_libraries:<5>:@edit` |
| doing a job of a kind (creating, rerunning, aborting it) | `admin:jobs:!upload_posts` |

A system built on the admin may have scopes of its own, a name before the
first `:`; the tree below shows each scope, the admin's as *Admin*.

## Which one decides

1. What is given by the **unique name** decides first: a **deny** — of any of
   the user's roles — denies; an allow allows.
2. Only when nothing is given by the unique name do the **groups** decide: an
   allow of any role allows.

So a group may be allowed while one part of it is denied by its name, or a
group denied while one part of it is allowed by its name. A permission
declared **out of the group** (`admin:/site-files:!edit`) has priority over
the one declared with the group (`admin:site/:/site-files:!edit`): with both,
the one out of the group decides. And it keeps holding when the menu moves the
part to another group — the one with the group no longer matches.

## The stack of permissions

The permissions of a role are an **ordered list**, one a line. A user's are a
**stack**: the system's first, then those granted them on a record (a share),
then those of **each of their roles, in the order of the roles** — the order
they have in the user's form (drag to change it).

Each line is a resource, and may begin with a **prefix** — the word **and a
space**:

| Line | What it does |
| --- | --- |
| `admin:pages:*` | **allows** what matches it, if no deny matches |
| `deny admin:users:*` | **denies** what matches it |
| `allow admin:users:@list` | the **earlier** denies stop counting for what matches it — a wider one too —, and it allows: it **forces** the allow of just what matches it |
| `reset admin:users:*` | the **earlier** lines stop counting for what matches it — a wider one too —; it allows nor denies anything |

A line **covers** another when its resource matches the other's: `admin:*`
covers `admin:users:@list`; `admin:users:@list` does not cover
`admin:users:*`. A line after one it covers **replaces** it:

- `deny` and `allow` replace the earlier ones of no prefix, `allow` and `deny`;
- `reset` replaces any;
- one of no prefix replaces the earlier ones of no prefix (and `allow`).

An `allow` or a `reset` **narrower** than an earlier line — which covers it —
does not replace it: it takes from it just what it matches. After
`deny admin:users:*`, an `allow admin:users:@list` frees just the list; the
deny still counts for the rest of the users. A narrower line of no prefix
takes nothing: the deny wins over it.

What is left are the **active** permissions: those of no prefix (an `allow`
becomes one) and the `deny` ones. They decide as always: a resource is
allowed when an active one of no prefix matches it and no active `deny` does.

**Groups** are several lines, in their order: `admin:users:{@list,@create}` is
`admin:users:@list` and `admin:users:@create`. An alternative may have its own
prefix — `admin:users:{*, deny @delete}` is `admin:users:*` and then
`deny admin:users:@delete`; `admin:users:{*, deny {@edit,@delete}}` denies
just those two —, anywhere in the resource: parts, fields, actions. A new line
inside a group separates as a comma does.

Examples (a user's stack, from top to bottom):

| Stack | May | May not |
| --- | --- | --- |
| `admin:*`<br>`deny admin:users:*` | all but the users | the users |
| `admin:*`<br>`deny admin:users:*`<br>`allow admin:users:*` | all: the `allow` undoes the `deny` | — |
| `admin:*`<br>`deny admin:users:*`<br>`reset admin:users:*` | all: the `reset` takes the `deny` out, and `admin:*` counts | — |
| `admin:users:*`<br>`reset admin:*`<br>`admin:pages:*` | only the pages: the `reset` took the users out | the rest |
| `admin:*`<br>`deny admin:users:*`<br>`allow admin:users:@list` | all, and the list of the users: the `allow` forces just it | the users, but the list |
| `admin:*`<br>`deny admin:users:*`<br>`admin:users:@list` | all but the users: without the `allow`, the deny wins | the users, their list too |
| `admin:*`<br>`reset admin:users:@delete` | all but deleting users: the `reset` took it from `admin:*` | deleting users |
| `admin:users:{*, deny @delete}` | the users, but deleting | deleting users |

So a later role may deny what an earlier one allows (`deny`), or allow what an
earlier one denies (`allow`).

The page **Permissions** of each user's detail shows their stack: in the tab
**All**, each line with its index, where it comes from (the system, a role)
and whether it is active — or active but what a later line matches —, or was
replaced — and by which —; in the tab **Active**, only the ones that count,
each with what it no longer covers ("but …"), as a list or as a tree.

## The access keys

The permissions of an **access key** accept only **denies**
(`deny admin:users:*`), applied on top of what its user may: the key may all
its user may, but what it denies. With none, all its user may. The examples
are in {%= admin.model("my_access_keys").link %}.

## The permissions of the admin

Every part of the admin and what may be asked of it — open a node to see what
is inside it. Inside a model — its records too — what it has is in groups:
its permissions, actions, fields, sections, pages, nested models and
permissions of its own, each with its resource. Each part of the
documentation also has its own **Permissions** item.

{%= admin.permissionsTree() %}

Ask whoever administers the roles for what your work needs.
