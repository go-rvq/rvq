# 1. Permissões

O que um usuário pode fazer vem dos **papéis** dele: cada papel permite, parte
por parte do admin, listar, ver, criar, editar e excluir os registros, e
executar as ações.

- **O que não é permitido não aparece**: um item do menu lateral, um botão,
  uma ação, uma aba, uma coluna. Uma parte sem nada permitido não está no
  menu.
- **Campos e seções** podem ser permitidos para ver mas não para editar —
  aparecem somente leitura —, ou não ser permitidos — não aparecem.
- **O servidor confere também**: um endereço digitado à mão, ou um botão de
  uma página antiga, é recusado quando o papel não permite.
- A **lixeira** de uma parte só é vista pelos papéis com permissão para vê-la
  ({%= admin.doc("guides/policies/02-trash").link %}).

Peça a quem administra os papéis o que o seu trabalho precisa.
