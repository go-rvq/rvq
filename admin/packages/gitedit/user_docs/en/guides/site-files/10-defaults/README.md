# The default files

The application brings the **default files** of a site: the layouts, the
components and the configuration it starts with. Two actions use them, each
with its own permission. They are not part of what is given to whoever only
edits the files, because they replace the files of the site as a whole.

| Action | Permission |
| --- | --- |
| Start with the default files | `!init_defaults` |
| Restore the default files | `!reset_defaults` |

## The empty repository

The first time, with no repository in `public/`, the application makes one
as it starts. When `public/` is not there or is empty, the repository is
**empty**: no files and no commits. While it stays that way, the page
{%= admin.page("/site-files").link %} shows only the **Empty repository**,
with two ways to give it its first files:

- **Start with the default files** (`!init_defaults`): the default files are
  checked (the templates compile) and go in a first commit of yours, with the
  `Site-User` of whoever made it. The site has them at once.
- **By git** (`!git` and `!publish`): make a repository in the folder of the
  files, commit (with the hook that signs the commits) and push the branch.
  The push is checked as **Publish** does. The section has the commands ready
  to copy.

When `public/` has files and is no repository, those files go in a first
commit, `first commit`, without the hidden ones (`.gitignore` with `.*`). The
administrator of the server decides whether the application makes the
repository (`PUBLIC_GIT_INIT`).

## Restore the default files

**Restore the default files** (`!reset_defaults`) brings the default files of
the version of the application that is running into **your draft**:

- before you confirm, the panel shows what changes: the tree of the files (A
  new, M different, D removed) and, for each one opened, the draft's and the
  default side by side. The tree follows the choice of removing the extras;
- each default file replaces the draft's, and the other files stay. The
  draft's `.gitignore` stays as it is: the default goes in only when it has
  none;
- with **Remove the files the defaults do not have** checked, the files of
  the folders of the defaults (such as `templates/` and `config/`) that the
  defaults do not have go too;
- it all goes in a commit of the draft, and the history stays. It is not a
  *revert*: the default files are those of the application of now, which may
  have changed since the first commit.

The draft must have no changes out of a commit. When the files do not work (a
template does not compile), nothing changes. Then check the **preview** and
**publish**, as with any other commit.
