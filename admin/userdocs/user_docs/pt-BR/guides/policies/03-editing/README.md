# 3. Edição

Uma alteração só é salva quando pode ser:

- **Outra pessoa salvou antes**: o formulário lembra o registro como estava
  quando foi aberto. Se alguém o salvou nesse meio-tempo, salvar é recusado,
  dizendo quem e quando — assim o trabalho de ninguém é sobrescrito. Recarregue
  o formulário e refaça a alteração.
- **Um campo obrigatório ou inválido**: o que está errado aparece sob o campo,
  e nada é salvo até ser corrigido.
- **Um registro excluído** não é editado: veja
  {%= admin.doc("guides/policies/02-trash").link %}.

O que nunca é editado:

- **Registros guardados como histórico** — revisões, registros de
  atividade —: só são consultados.
- **Configurações** têm um único registro: são editadas, nunca criadas nem
  excluídas.
- Uma página ou postagem **publicada** continua mostrando no site o que foi
  publicado: uma alteração vai para lá só quando é publicada de novo
  ({%= admin.doc("actions/Publisher").link %}).
