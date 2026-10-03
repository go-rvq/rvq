# Policies and permissions

The rules the admin follows for everybody: what each one may see and do, what
happens to what is deleted, and when a change is refused. Nothing here is
configured screen by screen — it holds in every part of the admin.

1. {%= admin.doc("guides/policies/01-permissions").link %}: what a role lets
   its users see and do.
2. {%= admin.doc("guides/policies/02-trash").link %}: a deleted record, kept
   read only until it is brought back.
3. {%= admin.doc("guides/policies/03-editing").link %}: when a change is
   refused, and what cannot be changed at all.
4. {%= admin.doc("guides/policies/04-deleting").link %}: what a deletion takes
   with it, and what cannot be deleted.

Each action that depends on the record says when it is available, in its
document and in the menu of the detail (see the **Detail** of each part).
