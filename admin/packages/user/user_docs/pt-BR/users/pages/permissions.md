# Permissões

A pilha de permissões do usuário, aberta no detalhe dele (**Permissões**): as
do sistema primeiro, depois as dadas a ele num registro, depois as de cada
papel dele, **na ordem dos papéis** — os grupos já quebrados, uma permissão
por linha (`admin:users:{@list,@create}` são duas).

- **Tudo**: cada linha, com o **Índice** (o lugar dela na pilha), a
  **Permissão**, a **Origem** (o papel de onde veio — `role: Editor` —, o
  sistema — `code: …` —, um registro — `database: …`) e a **Situação**:
  *ativa*, *ativa, menos o que #N casa* (um `allow` ou `reset` mais estreito
  depois dela), *substituída pela #N* (a linha que a substituiu) ou
  *remoção* (um `reset`, que só tira das anteriores).
- **Ativas**: só as que valem — as sem prefixo (um `allow` aparece sem ele:
  permite) e as `deny`, cada uma com o que não vale mais para ela
  ("menos admin:users:@list"). **Em árvore** as mostra agrupadas pelo começo comum:

  ```
  admin:users:{
  	@list
  	@create
  }
  ```

Como uma linha substitui outra está em
{%= admin.doc("guides/policies/01-permissions").link %}. O usuário inicial
pode tudo, pelo sistema: os papéis dele não contam.
