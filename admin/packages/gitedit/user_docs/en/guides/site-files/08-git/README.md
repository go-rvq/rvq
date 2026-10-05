# By git

The same files can be cloned and pushed with **git**, from your computer,
with your **login** and an **access key** (see below) — or your admin
password. There are two repositories — their
URLs are on the page {%= admin.page("/site-files").link %}, under **By git**:

| Repository | A push… | Asks |
| --- | --- | --- |
| **your draft** — `…/site-files-draft.git` | changes the draft: the editor and the **preview** show it at once; then **publish** from the page | `!git` (see below) |
| **the site** — `…/site-files.git` | **publishes** at once | `!git` and `!publish` |

```
git clone https://your-site/admin/site-files-draft.git
cd site-files-draft
# … edit, commit …
git push
```

## The permissions of each repository

They are those of the page {%= admin.page("/site-files").link %} (step
**Permissions**), the editor's own. Declared out of the group — they hold
wherever the menu puts the page —:

| To | Your draft | Another user's draft | The site |
| --- | --- | --- | --- |
| clone, `fetch`, `pull` | `{%= admin.page("/site-files").uniqueResource %}@get` of **each** file and `{%= admin.page("/site-files").uniqueResource %}!git` | the same, and the draft **shared with you** (or `{%= admin.page("/site-files").uniqueResource %}!drafts`) | `{%= admin.page("/site-files").uniqueResource %}@get` of each file and `{%= admin.page("/site-files").uniqueResource %}!git` |
| `push` | `{%= admin.page("/site-files").uniqueResource %}!git` and, of each file changed, what the editor asks (`!create`, `!edit`, `!rename`, `!move`, `!delete` — see **What a push checks**) | the same, and the draft shared with you (or `{%= admin.page("/site-files").uniqueResource %}!drafts`) | the same as the draft's, and `{%= admin.page("/site-files").uniqueResource %}!publish` |

- A permission can be **of a path** and all under it: `{%= admin.page("/site-files").uniqueResource %}<static/*>:!edit`
  changes `static/` only (step **Permissions**).
- In the draft, the push is refused while there are changes **not
  committed** made in the editor; in the site, the push checks what
  **Publish** does.
- An **access key** only restricts: it needs the permissions above **and**
  your role too. For your draft, a key with:

  ```
  {%= admin.page("/site-files").uniqueResource %}@get
  {%= admin.page("/site-files").uniqueResource %}!{git,create,edit,rename,move,delete,import}
  ```

  and, to push to the site too, `{%= admin.page("/site-files").uniqueResource %}!publish`.

## With an access key

Instead of your password, use an **access key** — the way for git: it works
even when you sign in to the admin with Google or a second factor, it can be
for the site files only, it expires and is revoked without touching your
password. Make yours in {%= admin.model("my_access_keys").link %}, with the
permissions of the site files (`admin:/site-files:@get`, `!git` and the
actions on the files you will change); its **code** shows only once — copy
it.

### Clone

git asks for the username and the password: the username is your **login**;
the **password**, the key's **code**.

```
git clone https://<the-site>/admin/site-files-draft.git
Username for 'https://<the-site>': <your-login>
Password for 'https://<your-login>@<the-site>': <the-key-code>
```

### Not asked again

With nothing else, git asks for the code at each `pull` and `push`. Keep it
in git's **credential helper**, once for the site:

- **In the file** `~/.git-credentials` (text, readable by you only):

  ```
  git config --global credential.https://<the-site>.helper store
  ```

  The next `git pull` (or `push`) asks for the username and the code and
  keeps them; the next ones do not ask. The file gets a line
  `https://<your-login>:<the-code>@<the-site>`.

- **In the system's keychain**, encrypted: instead of `store`,
  `osxkeychain` (macOS), `manager` (Windows, the Git Credential Manager) or
  `libsecret` (Linux, when installed).

Do not put the code **in the URL** (`https://login:code@…`): it is written in
the clone's `.git/config`, in plain sight, and goes along if the clone is
copied.

### When the key expires (or is revoked)

git starts answering `Authentication failed`. Make a **new key** in
{%= admin.model("my_access_keys").link %} and change the one kept:

1. Forget the old one:

   ```
   printf 'protocol=https\nhost=<the-site>\n\n' | git credential reject
   ```

   (with `store`, editing the site's line of `~/.git-credentials` also
   does.)

2. The next `git pull` asks for the username and the code: give the **new**
   one — it is kept in place of the old.

If the code was in the clone's URL, take it out:

```
git remote set-url origin https://<the-site>/admin/site-files-draft.git
```

Each use of a key is in its history and in your access log.

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

## Who made each commit

Each commit pushed must say **who made it on this site**: a line at the end
of its message, as git's `Co-authored-by` —

    Site-User: your-login <your-key@the-site's-address>

A push of a commit without it (or of another site, or of a user the site does
not know) is refused, and git shows the line missing. The **hook**
`commit-msg` adds the line to each commit; in **By git**, on the page, is the
command that installs it in your clone:

    curl -fsSL -u your-login …/site-files/commit-msg -o .git/hooks/commit-msg && chmod +x .git/hooks/commit-msg

To sign commits already made: `git rebase -x 'git commit --amend --no-edit'
<the commit before them>`.

## Another user's draft

A draft shared with you (step **Work together**) is cloned and pushed by git
too: `…/site-files-drafts/<the owner's key>.git` — the address is on the page
of the draft. Your commits stay yours.

## Cloning

A clone hands over **every** file: so it asks `@get` of each one. If a path
may not be seen (a denial of `@get` on it), the clone is refused.

## With the local development

The site's repository is the deploy's (the step **Sync with the local
development**): a push here is as a publishing — whoever develops brings what
was published before the next deploy.
