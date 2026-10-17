# Os arquivos padrão

A aplicação traz os **arquivos padrão** de um site: os layouts, os
componentes e a configuração com que ele começa. Duas ações os usam, cada uma
com a sua permissão. Elas não estão no que se dá a quem só edita os arquivos,
porque trocam os arquivos do site como um todo.

| Ação | Permissão |
| --- | --- |
| Começar com os arquivos padrão | `!init_defaults` |
| Restaurar os arquivos padrão | `!reset_defaults` |

## O repositório vazio

Na primeira vez, sem repositório em `public/`, a aplicação cria um quando
sobe. Se `public/` não existe ou está vazio, o repositório fica **vazio**: sem
arquivos e sem commits. Enquanto ele estiver assim, a página
{%= admin.page("/site-files").link %} mostra só o **Repositório vazio**, com
dois caminhos para os primeiros arquivos:

- **Começar com os arquivos padrão** (`!init_defaults`): os arquivos padrão
  são conferidos (os templates compilam) e entram num primeiro commit seu,
  com o `Site-User` de quem o fez. O site passa a tê-los na hora.
- **Pelo git** (`!git` e `!publish`): faça um repositório na pasta dos
  arquivos, faça o commit (com o hook que assina os commits) e envie o branch.
  O push é conferido como em **Publicar**. A seção traz os comandos prontos
  para copiar.

Se `public/` já tem arquivos e não é um repositório, eles mesmos entram num
primeiro commit, `first commit`, sem os ocultos (`.gitignore` com `.*`). Quem
decide se a aplicação cria o repositório é o administrador do servidor
(`PUBLIC_GIT_INIT`).

## Restaurar os arquivos padrão

**Restaurar os arquivos padrão** (`!reset_defaults`) leva os arquivos padrão
da versão da aplicação que está rodando para o **seu rascunho**:

- antes de confirmar, o painel mostra o que vai mudar: a árvore dos arquivos
  (A novo, M diferente, D removido) e, para cada um aberto, o do rascunho e
  o padrão lado a lado. A árvore muda junto quando você marca ou desmarca a
  remoção;
- cada arquivo padrão substitui o do rascunho, e os outros arquivos ficam. O
  `.gitignore` do rascunho fica como está: o padrão só entra se ele não tiver
  um;
- marcando **Remover os arquivos que o padrão não tem**, saem também os
  arquivos das pastas do padrão (como `templates/` e `config/`) que ele não
  tem;
- tudo entra num commit do rascunho, e o histórico continua. Não é um
  *revert*: os arquivos padrão são os da aplicação de agora, que podem ter
  mudado desde o primeiro commit.

O rascunho precisa estar sem alterações fora de um commit. Se os arquivos não
funcionam (um template não compila), nada muda. Depois, confira no
**preview** e **publique**, como qualquer outro commit.
