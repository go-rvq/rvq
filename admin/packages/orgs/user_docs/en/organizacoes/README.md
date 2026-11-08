# Organizations

An organization holds what is under it: each resource of the menu under the
organization — people, finance… — is of the organization chosen, and an
organization's records are never seen from another.

## The organization chosen

The organization you work in is at the **top of the menu**: its chip says
which one it is, its resources are under it. Click the chip to choose
another; the ✕ clears the choice. The one chosen is the one of the page you
are on, else the one you chose last in this browser, else your only one,
else **your default organization**.

Your default is yours alone — any organization you reach, yours or shared
with you: in its detail, {%= admin.action("organizacoes", "set_my_default").link %}
(and {%= admin.action("organizacoes", "unset_my_default").link %} to have
none). **List organizations** lists the organizations you reach — only
those —, **My default** marking yours.

## Who reaches an organization

An organization has an **owner**, who reaches everything in it, and is
**shared** with users and groups, each with a **role**:

| Role | May |
|---|---|
| **Viewer** | see |
| **Editor** | see and change |
| **Manager** | see, change and delete |
| **Administrator** | all of that, and share the organization |
| **None** | nothing: takes a record out of what is shared above it |

A role given on an organization is good for **what is under it** too; one
given on a record under it wins over it (the most specific wins). A user
reaches by their own accesses and by those of their **groups**: whoever
joins a group reaches what the group reaches at once, and whoever leaves it
loses it at once.

In the detail of an organization (its ⋮ menu):

- {%= admin.action("organizacoes", "share_access").link %}: give a user or a
  group a role;
- {%= admin.action("organizacoes", "accesses").link %}: who reaches it, and
  their roles — change or revoke them;
- {%= admin.action("organizacoes", "access_history").link %}: everything that
  happened to its sharing;
- {%= admin.action("organizacoes", "transfer_owner").link %}: give it another
  owner.

## The owner confirms by e-mail

Nothing is given without the **owner**: a new sharing, a role that may more
than before, a new owner — whoever asks, the owner gets an **e-mail** with a
link, good for 48 hours, to a page where they **Confirm** or **Decline** it,
logged in as themselves. Until then nothing changes. A role that may less
than before, a revocation and a sharing as **None** are at once.

## The invite

Confirmed by the owner, a sharing with a **user** is an **invite**: it
appears in their **notifications** (the bell, at the top), with **Accept**
and **Decline**. Accepted, the organization is in their list and in the
choice of the menu. A sharing with a **group** needs no invite: it is good
at once for its members.
