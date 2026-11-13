# 1.2 Sets of permissions

A **set** is a name for permissions given together: reading, writing,
removing. It goes at the end of a line, after `%`:

    admin:posts:%readers           lists the posts and sees each one
    deny admin:users:%editors      denies what editors gives to the users

The line keeps **the name**, not what it gives. What it gives is worked out
again each time a user's permissions are made. So:

- when the set **changes**, every line of it changes with it;
- when it is **no more**, or gives nothing to that kind of resource, the line
  is kept and **ignored**: it allows and denies nothing, and is shown red in
  the field of the permissions.

## What a set gives

A set says what it gives **by kind of resource**: `model` (a part with
records), `singleton` (a part of a single record), `page` (a page), `check`
(a permission of its own). For example:

    readers  {model: {@list, <*>:@get}, singleton: @get, page: @get}

Given to `admin:posts:`, which is a *model*, it gives `admin:posts:@list` and
`admin:posts:<*>:@get`. Given to a page, it gives that page's `@get`.

A `deny` line of a set denies everything it gives. A set may hold others
(`editors` is `{%authors, %removers}`).

A resource may have **its own version** of a name. The site files', for
example, says what *publishers* is there. It takes the place of the kind's
version.

## The default sets

| Set | What it gives |
| --- | --- |
| `%readers` — Reading | lists and sees; of a single record and of pages, sees |
| `%authors` — Authoring | lists, creates, sees and edits; does not delete |
| `%removers` — Removal | the trash, deleting and restoring records |
| `%editors` — Editing | authoring and removing together |
| `%managers` — Management | everything |

Packages and applications add their own, such as translation, comment
moderation and history.

## In the field of the permissions

The **+** of a row of the tree shows, under the parts, the section **Sets of
permissions**: those that give something to that resource, each with what it
gives there. Choosing one inserts the line of the set (`…:%readers`).

In the tree, the line is shown under its resource, with the chip of its name.
Hovering over the chip shows what it grants. What the set allows is active in
the rows of the tree, as any line's.

## Sets of the admin

In **Administrator › Permission sets** new sets are made:

- a **name**: lower case letters, digits and `_`, a letter first;
- the **definition**, by kind of resource, as above;
- a **title** and a **description** in each language of the admin.

The name may not be one of a set of the system. Saving or deleting a set
makes the users' permissions again at their next request.
