# Permissões da página

Quem pode abrir esta página é dado nestes recursos, nas permissões do papel.
O nome único decide antes dos grupos (negar prevalece).

{%= admin.permissions() %}

Veja {%= admin.doc("guides/policies/01-permissions").link %}.

Os recursos começam por `:`: o **escopo** do admin é o padrão, vazio — veja os escopos em {%= admin.doc("guides/policies/01-permissions").link %}.
