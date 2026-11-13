# Conjuntos de permissões

Um conjunto é um nome para permissões dadas juntas. Uma linha de um papel o
usa com `%`: `admin:posts:%readers` lista as postagens e vê cada uma. Os
papéis estão em {%= admin.model("roles").link %}.

Aqui se criam os conjuntos do admin, ao lado dos do sistema (`readers`,
`authors`, `editors`…):

- **Nome**: como as linhas o escrevem depois de `%` — letras minúsculas,
  algarismos e `_`, começando por letra. Não pode ser o de um conjunto do
  sistema.
- **Definição**: o que ele dá por tipo de recurso, por exemplo
  `{model: {@list, <*>:@get}, singleton: @get}`.
- **Títulos** e **Descrições**: um em cada idioma do admin. É o que o campo
  de permissões mostra.

A linha guarda o nome. Mudar um conjunto muda o que todas as linhas dele
dão, no próximo acesso de cada usuário. Um conjunto excluído deixa as suas
linhas guardadas e ignoradas, mostradas em vermelho.

Tudo isso, com os conjuntos do sistema, em
{%= admin.doc("guides/policies/01-permissions/02-sets").link %}.
