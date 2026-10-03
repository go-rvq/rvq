# Permissões: {%= doc.model %}

O que um papel pode fazer aqui é dado nestes recursos, nas permissões do
papel ({%= admin.doc("guides/policies/01-permissions").link %}). Um recurso
terminado em `*` vale para tudo abaixo dele; `<id>` é um registro.

{%= admin.permissions() %}

Uma permissão pelo **nome único** decide antes de uma dada pelos **grupos**:
negar pelo nome único — em qualquer papel — nega, permitir permite; só quando
nada é dado por ele os grupos decidem. Os grupos acompanham o menu: mover esta
parte para outro grupo muda o recurso pelos grupos, nunca o nome único.

Os recursos começam por `admin:`, o **escopo** do admin — veja os escopos em {%= admin.doc("guides/policies/01-permissions").link %}.
