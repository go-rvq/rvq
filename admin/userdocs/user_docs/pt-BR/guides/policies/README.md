# Políticas e permissões

As regras que o admin segue para todos: o que cada um pode ver e fazer, o que
acontece com o que é excluído, e quando uma alteração é recusada. Nada aqui é
configurado tela por tela — vale em todas as partes do admin.

1. {%= admin.doc("guides/policies/01-permissions").link %}: o que um papel
   deixa seus usuários verem e fazerem.
2. {%= admin.doc("guides/policies/02-trash").link %}: um registro excluído,
   mantido somente leitura até ser restaurado.
3. {%= admin.doc("guides/policies/03-editing").link %}: quando uma alteração
   é recusada, e o que não pode ser alterado.
4. {%= admin.doc("guides/policies/04-deleting").link %}: o que uma exclusão
   leva junto, e o que não pode ser excluído.

Cada ação que depende do registro diz quando está disponível, no documento
dela e no menu do detalhe (veja o **Detalhe** de cada parte).
