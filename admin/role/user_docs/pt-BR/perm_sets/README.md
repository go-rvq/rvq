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
- **Título** e **Descrição**: no idioma do registro. É o que o campo de
  permissões mostra.

Um conjunto é **localizado pelos idiomas do admin**, como as páginas e
postagens pelos do site: um registro em cada idioma, cada um com o seu título
e a sua descrição. A listagem mostra os do idioma do admin; **Localizar**
copia um conjunto para outro idioma e **Localizações** mostra os que ele já
tem. Nome e Definição são do conjunto: salvos num idioma, valem em todos. Um
idioma do admin criado a partir de outro recebe os conjuntos já traduzidos,
pelo tradutor automático, quando ele está ligado.

Os conjuntos **do sistema** (`readers`, `authors`…) são registrados aqui
sozinhos, marcados *Do sistema*: Nome e Definição são os do código, só
leitura, e eles não são excluídos. Título e Descrição podem ser traduzidos.

O detalhe mostra quem criou e quem alterou por último, e quando. A aba de
histórico guarda cada mudança.

A linha guarda o nome. Mudar um conjunto muda o que todas as linhas dele
dão, no próximo acesso de cada usuário. Um conjunto excluído deixa as suas
linhas guardadas e ignoradas, mostradas em vermelho.

Tudo isso, com os conjuntos do sistema, em
{%= admin.doc("guides/policies/01-permissions/02-sets").link %}.
