# 1.1 The tree of the permissions

Every permission of the admin, down to the leaves: each leaf is one
permission — what a line of the stack writes, as
`admin:posts:<*>:#Title:@edit`. Two views:

- **By verb**: a model by what is asked of it — listing, creating, seeing a
  record, editing… —, each verb with what it asks: the fields of the form,
  the sections (with their fields), the actions, the pages, the nested
  models.
- **By kind**: a model by what it holds — permissions, actions, fields,
  sections, pages, nested models —, each part with the verbs asked of it.

Without the groups of the menu (the default), each permission is by its
unique name, which decides first; with them, by the path of the groups too.
**Kinds** shows only the rows of the kinds chosen (and what is above them);
**Search** looks in the name, the description and the permission.

{%= admin.permissionsTree() %}

What a user has active, row by row, is in the **Tree** tab of their
permissions (Users › the user › Permissions).
