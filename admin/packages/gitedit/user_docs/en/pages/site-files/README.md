# Site files

The files the site is made of — its templates, styles, scripts and images —
edited here, in a **draft** of your own: the site changes only when you
**publish**.

- **Files** is the editor: the tree of the files, and each file open in a tab.
  Saving changes your draft only.
- **Your draft** says where it stands: commits to publish, commits others
  published that it does not have.
- **Changes not committed** lists each file changed since the last commit;
  open it to see what changed. **Discard** drops the changes of a file.
- **Commit** records the changes, with a message saying what they are; the
  commit carries your name.
- **Update** brings into your draft what others published: your commits go
  on top of theirs, your changes are kept.
- **Publish** puts your commits on the site.
- **Preview the site** opens the site as your draft makes it: navigate it,
  nothing of it is seen by the visitors.
- **Start over** drops your draft: a new one is made from the site.

## When it is available

- **Commit** and **Publish** check the draft first: if the site, as the
  draft makes it, shows an error — a template that does not compile —, nothing
  is committed nor published, and the error is shown.
- **Publish** needs every change committed, and the draft updated with what
  others published.
- If a file of the site was changed out of this editor, **Publish** overwrites
  nothing and says so.

## Permissions

Each action is a permission of the page: seeing it (`@get`), changing the
files of the draft (`!edit`), `!commit`, `!update`, `!publish`, `!discard`,
`!reset` and `!preview` — publishing can be given to fewer people than
editing.
