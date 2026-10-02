# Edição: {%= doc.model %}

O formulário que altera um registro. Abra-o com o lápis (**Editar**) do
detalhe dele, ou pelo menu da linha (**⋯**) na listagem; o que for alterado é
mantido quando o formulário é salvo.

{%= admin.fields("edit") %}

Um campo pode ser obrigatório, ou conferido quando o formulário é salvo: o
que está errado aparece sob o campo, e nada é salvo até ser corrigido. Um
campo que *aparece conforme o registro* surge só nos registros a que se
aplica.
