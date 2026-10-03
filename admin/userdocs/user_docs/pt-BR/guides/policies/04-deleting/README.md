# 4. Exclusão

- Uma exclusão é **confirmada** antes: o diálogo diz o que é excluído.
- **Registros que dependem dele**: o diálogo os lista (**Ver itens
  relacionados**), e a exclusão pode levá-los junto (**Excluir items
  relacionados**).
- Onde a listagem tem lixeira, o registro vai para lá e pode ser restaurado
  ({%= admin.doc("guides/policies/02-trash").link %}); nos demais casos a
  exclusão é definitiva.
- O que não pode ser excluído não é oferecido: as configurações, um registro
  guardado como histórico, um registro excluído.

Veja {%= admin.doc("actions/Delete").link %}.
