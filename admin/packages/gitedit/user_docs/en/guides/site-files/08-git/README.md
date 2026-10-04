# By git

The same files can be cloned and pushed with **git**, from your computer,
with your admin **login and password**. There are two repositories — their
URLs are on the page {%= admin.page("/site-files").link %}, under **By git**:

| Repository | A push… | Asks |
| --- | --- | --- |
| **your draft** — `…/site-files-draft.git` | changes the draft: the editor and the **preview** show it at once; then **publish** from the page | `!git` |
| **the site** — `…/site-files.git` | **publishes** at once | `!git` and `!publish` |

```
git clone https://your-site/admin/site-files-draft.git
cd site-files-draft
# … edit, commit …
git push
```

## What a push checks

Each file changed asks **what the editor asks** — the same permissions, of
each path (see the step **Permissions**):

| In the push | Asks |
| --- | --- |
| a new file | `!create` |
| a file changed | `!edit` |
| a file deleted | `!delete` |
| renamed in its folder | `!rename` (from and to) |
| moved to another folder | `!move` (from and to) |

One denial refuses the whole push — nothing written — and git shows the
files: `remote: config/layout_config.gad: no permission (!edit)`.

And:

- **Only the site's branch**, and only forward (fast-forward): if others
  published, `git pull` first.
- **To the draft**, a push is refused while the editor has changes not
  committed: commit or discard them first.
- **To the site**, a push checks what **Publish** checks: the templates
  compile, and the site's files were not changed out of git. A push and a
  publishing from the page never happen at once.

## Cloning

A clone hands over **every** file: so it asks `@get` of each one. If a path
may not be seen (a denial of `@get` on it), the clone is refused.

## With the local development

The site's repository is the deploy's (the step **Sync with the local
development**): a push here is as a publishing — whoever develops brings what
was published before the next deploy.
