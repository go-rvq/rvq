# Permissions of the group

A permission on the group holds for everything inside it — the groups, the
parts and the pages under it —, and shows the group in the menu when one of
them is allowed.

{%= admin.permissions() %}

The resource follows the menu: it is the chain of the groups the group sits
in. See {%= admin.doc("guides/policies/01-permissions").link %}.

The resources begin with `:`: the admin's **scope** is the default, empty — see the scopes in {%= admin.doc("guides/policies/01-permissions").link %}.
