# Roles

A role is a set of permissions: what its users may see and do in the admin.
A user has one or more roles (see {%= admin.model("users").link %}).

![The roles](images/listing.png)

- The role **Administrator** may do everything.
- Make a role for each kind of work — editing the site, answering the
  contact form… — and give it only what that work needs.

## System roles

Some roles are the **system**'s: the application needs them, and makes them
by itself when missing — the listing and the detail say *System role* and
what each one is for.

- A system role is **not deleted** nor **renamed**.
- Its **permissions** can be changed as any role's, and go back to the
  originals with the action **Reset the permissions**, in the detail.
- A system role whose permissions are the system's own — the
  **Administrator**, who may do everything — has no permissions to change.
- A role of the same name that existed already becomes the system's: its
  permissions are kept, and the originals it lacks, added.
