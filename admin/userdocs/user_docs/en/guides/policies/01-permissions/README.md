# 1. Permissions

What a user may do comes from their **roles**: each role allows, part by part
of the admin, to list, see, create, edit and delete its records, and to run
its actions.

- **What is not allowed is not shown**: an item of the side menu, a button, an
  action, a tab, a column. A part with nothing allowed is not in the menu.
- **Fields and sections** may be allowed to see but not to edit — they are
  shown read only —, or not at all — they are not shown.
- **The server checks too**: an address typed by hand, or a button of an old
  page, is refused when the role does not allow it.
- The **trash** of a part is seen only by the roles allowed to see it
  ({%= admin.doc("guides/policies/02-trash").link %}).

Ask whoever administers the roles for what your work needs.
