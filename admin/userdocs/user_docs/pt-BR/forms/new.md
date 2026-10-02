# Cadastro: {%= doc.model %}

O formulário de um novo registro. Abra-o com **Novo** no topo da listagem; o
registro é criado quando o formulário é salvo, e o detalhe dele se abre.

{%= admin.fields("new") %}

Um campo pode ser obrigatório, ou conferido quando o formulário é salvo: o
que está errado aparece sob o campo, e nada é salvo até ser corrigido. Um
campo que *aparece conforme o registro* surge quando os campos de que ele
depende são preenchidos.
