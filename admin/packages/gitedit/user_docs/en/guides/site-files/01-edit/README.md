# Edit in the draft

Each person edits a **draft** of their own: a copy of the files of the site,
made the first time they open the page. What you save changes your draft
only — the site and the drafts of the others stay as they are.

![A file changed in the draft, its diff open](images/changes.png)

- **Files**, on the right, is the editor: the tree of the files at its left,
  each file open in a tab. Save with the tab's button or Ctrl+S.
- Create, rename and delete files and folders: from the tree's menu.
- The editor formats and checks the `.gad`/`.gadx` files and shows the
  language's documentation; it does not run code.
- **Changes not committed**, on the left, lists each file changed since the
  last commit; open one to see what changed. **Discard** drops its changes.
- **Start over** drops the whole draft (the commits not published too): a new
  one is made from the site.

The files are in folders: `templates/` (the pages, the layouts, the
components), `static/` (styles, scripts, images) and `config/` (the options
of the layouts).

## By WebDAV

The same draft can be opened from a computer, as a network drive, by the
admin's **WebDAV** (its address is in **Administrator → File System**), in
the folder `site-files`: edit the files in the program you like. The same
permissions as in the admin hold — seeing asks `@get`, changing asks
`!edit` —, and the `.git` folder is not shown. Committing and publishing
stay here, on this page.
