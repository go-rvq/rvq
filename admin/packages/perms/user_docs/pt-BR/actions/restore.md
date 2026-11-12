# Restaurar

Traz um registro excluído de volta da lixeira, como estava quando foi
excluído.

Um registro excluído de uma listagem não some: vai para a aba **Lixeira** da
listagem, com quem o excluiu, quando e de onde. Para trazê-lo de volta:

1. Abra a listagem do registro e escolha a aba **Lixeira**.
2. Abra o menu da linha (**⋯**) do registro — ou o detalhe dele.
3. Escolha **Restaurar da lixeira**.

Para vários de uma vez, marque-os na lixeira e escolha **Restaurar**, na barra:
se algum deles não puder ser restaurado, nenhum é, e a mensagem diz quais.

O registro volta para as outras abas da listagem, e a exclusão fica guardada
no histórico de revisões dele.

Veja também: {%= admin.doc("actions/deleted_origin").link %}.

## Quando está disponível

Só na aba Lixeira de uma listagem, para quem pode vê-la (`@trash` da parte), e
de cada registro para quem pode restaurá-lo (`@restore` dele:
`admin:posts:<7>:@restore`).
