# Permissions: {%= doc.model %}

What a role may do here is given on these resources, in the permissions of
the role ({%= admin.doc("guides/policies/01-permissions").link %}). A resource
ending in `*` holds everything under it; `<id>` is a record.

{%= admin.permissions() %}

A permission by the **unique name** decides before one through the
**groups**: a deny by the unique name — of any role — denies, an allow allows;
only when none is given by it do the groups decide. The groups follow the
menu: moving this part to another group changes its resource by the groups,
never its unique name.
