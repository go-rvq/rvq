# Sync with the local development

The files of the site are a **git** repository, in two places: on the
**server** — what the site uses and what this editor changes — and on the
computer of whoever **develops** the site (the project's `public/`). The
commits made here stay on the server; the ones made on the computer reach the
server on deploy. For neither side to lose the other's, each brings the
other's before sending its own.

## What was done in the admin, to the computer

Before changing `public/` and, above all, before a deploy, the developer
brings the server's commits:

```
pub repo public pull server main
```

(`pub` runs the git of `public/` with the server as `server`). Then records in
the project the new commit of `public/` (the submodule):

```
git add public && git commit -m "public: what was published from the admin"
```

Without that, the deploy is refused: the server has commits the computer does
not, and the deploy does not overwrite what was published from the admin.

## What was done on the computer, to the admin

The deploy (`pub run`) publishes the commit of the project's `public/` on the
server: the site uses it. Whoever edits in the admin brings those commits into
their draft with **Update**.

## Files changed on the server

If a file of the site is changed right on the server, out of git, both the
admin's **Publish** and the deploy stop and say so, overwriting nothing:
whoever runs the server decides what stays (a commit with the change, or
undoing it).
