# 2. A lixeira

Um registro excluído de uma listagem não some: vai para a **lixeira** da
listagem — a aba com a lixeira, à direita —, com quem o excluiu, quando e de
onde.

![A lixeira](images/trash.png)

Na lixeira um registro é **somente leitura**:

- não é editado — o detalhe dele não tem botão de editar, as seções não têm
  edição — nem excluído de novo;
- nenhuma ação dele é executada, exceto as feitas para a lixeira:
  {%= admin.doc("actions/restore").link %}, que o traz de volta, e
  {%= admin.doc("actions/deleted_origin").link %};
- a barra da lixeira oferece só **Restaurar** — nenhum novo registro, nenhuma
  outra ação —, e o menu de uma linha (**⋯**) só o que está dentro do registro
  (revisões, comentários…) e as páginas dele, para consultar.

O servidor recusa o mesmo, mostre a página o que mostrar. Para alterar um
registro excluído, restaure-o antes.

## Quando está disponível

A lixeira aparece para os papéis com a permissão **`@trash`** da parte
(`admin:posts:@trash`); os outros não veem a aba. Restaurar um registro pede a
permissão **`@restore`** dele (`admin:posts:<7>:@restore`, ou
`admin:posts:<*>:@restore` para todos): quem pode excluir um registro
compartilhado pode também restaurá-lo.
