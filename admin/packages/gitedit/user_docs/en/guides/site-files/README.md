# Step by step

How to edit the files the site is made of — templates, styles, scripts,
images — from the admin, in **{%= admin.page("/site-files").menu %}**, leaving the site
as it is until you publish; and how that goes along with whoever develops the
site on their own computer.

Step by step:

1. **Edit in the draft** — yours, only yours.
2. **Preview** — the site as the draft makes it.
3. **Commit** — record the changes, with a message.
4. **Publish** — put the commits on the site.
5. **Permissions** — who may do each thing.
6. **Sync with the local development** — the commits made here and the ones
   made on the developer's computer.
7. **Options of the layouts** — the fields of their forms, and a select of
   values with `options`.
8. **By git** — clone and push the draft, or the site, from your computer.
9. **Work together** — share your draft.
10. **The default files** — start an empty repository with them, or restore
    them.

**The first time**, with no repository in `public/`, the application makes
one as it starts. When `public/` is not there or is empty, the repository is
empty, and the page shows how to give it its first files: the default ones,
or a push by git (see **The default files**). When `public/` has files, they
go in a first commit, `first commit`, without the hidden ones (`.gitignore`
with `.*`). The administrator of the server decides it (`PUBLIC_GIT_INIT`).
