# Permissões do grupo

Uma permissão no grupo vale para tudo o que está dentro dele — os grupos, as
partes e as páginas abaixo dele —, e mostra o grupo no menu quando uma delas é
permitida.

{%= admin.permissions() %}

O recurso acompanha o menu: é a cadeia dos grupos em que o grupo está. Veja
{%= admin.doc("guides/policies/01-permissions").link %}.

Os recursos começam por `admin:`, o **escopo** do admin — veja os escopos em {%= admin.doc("guides/policies/01-permissions").link %}.
