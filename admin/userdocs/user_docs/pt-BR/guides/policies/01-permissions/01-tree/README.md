# 1.1 A árvore das permissões

Todas as permissões do admin, até as folhas: cada folha é uma permissão —
o que uma linha da pilha escreve, como `admin:posts:<*>:#Title:@edit`. Duas
visões:

- **Por verbo**: um modelo pelo que se pede dele — listar, criar, ver um
  registro, editar… —, cada verbo com o que pede: os campos do formulário,
  as seções (com os seus campos), as ações, as páginas, os modelos filhos.
- **Por tipo**: um modelo pelo que ele tem — permissões, ações, campos,
  seções, páginas, modelos filhos —, cada parte com os verbos pedidos dela.

Sem os grupos do menu (o padrão), cada permissão aparece pelo nome único,
que é o que decide primeiro; com eles, também pelo caminho dos grupos.
**Tipos** mostra só as linhas dos tipos escolhidos (e o que está acima
delas); **Buscar** procura no nome, na descrição e na permissão.

{%= admin.permissionsTree() %}

O que um usuário tem ativo, linha por linha, está na aba **Árvore** das
permissões dele (Usuários › o usuário › Permissões).
