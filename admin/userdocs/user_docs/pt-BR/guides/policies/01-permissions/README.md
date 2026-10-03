# 1. Permissões

O que um usuário pode fazer vem dos **papéis** dele: cada papel permite ou
nega, parte por parte do admin, listar, ver, criar, editar e excluir
registros, e executar ações.

- **O que não é permitido não aparece**: um item do menu lateral, um botão,
  uma ação, uma aba, uma coluna. Um grupo sem nada permitido não está no
  menu.
- **Campos e seções** podem ser permitidos para ver mas não para editar —
  aparecem somente leitura —, ou não ser permitidos — não aparecem.
- **O servidor confere também**: um endereço digitado à mão, ou um botão de
  uma página antiga, é recusado quando o papel não permite.
- A **lixeira** de uma parte só é vista pelos papéis com permissão para vê-la
  ({%= admin.doc("guides/policies/02-trash").link %}).

## Como uma permissão é escrita

Uma permissão nomeia um **recurso** — as partes dele separadas por `:` — e
termina no que é pedido dele:

| Parte | Exemplo |
| --- | --- |
| o escopo: quem pede — sempre o primeiro (veja abaixo) | `presets` |
| um grupo do menu: o nome e `/` | `site/` |
| uma parte (model) | `seo_config` |
| um registro: o id entre `<…>` | `<7>` — `<*>` qualquer um |
| um campo | `#Title` |
| uma seção do detalhe | `$Main` |
| uma página | `/report` |
| uma permissão: `@` e o nome | `@list`, `@get`, `@create`, `@edit`, `@delete` |
| uma ação: `!` e o nome | `!publish` |

Uma parte é alcançada de dois jeitos — **pelos grupos** e **pelo nome
único**, o dela:

    presets:site/:seo/:seo_config:<7>:@edit     pelos grupos
    presets:seo_config:<7>:@edit                pelo nome único

`*` vale por qualquer coisa: `presets:site/:*` é tudo o que está no grupo
*site*; `presets:seo_config:*`, tudo da parte *seo_config*, onde quer que o
menu a ponha.

## Os escopos

A primeira parte de um recurso é o **escopo**: a parte do sistema que pede a
permissão. Uma política o nomeia, ou usa `*` para qualquer um:

| Escopo | O que protege | Exemplo |
| --- | --- | --- |
| `presets` | o admin: cada item do menu — os grupos, as partes, os registros delas, campos, seções, ações e páginas | `presets:content/:posts:<12>:@edit` |
| `media_library` | a biblioteca de mídia como os campos de imagem e de arquivo a usam, num formulário: enviar um arquivo, excluir um, editar a descrição dele | `media_library:media_libraries:@upload`, `media_library:media_libraries:<5>:@delete`, `…:<5>:@update_desc` |
| `workers` | as tarefas: editar uma — uma tarefa que recebe argumentos —, pelo nome dela | `workers:upload_posts:@edit` |

A listagem da biblioteca de mídia no menu é uma parte de `presets`
(`presets:content/:media/:media_libraries:…`); `media_library` é só o que os
formulários fazem com os arquivos.

## Qual decide

1. O que é dado pelo **nome único** decide primeiro: **negar** — em qualquer
   papel do usuário — nega; permitir permite.
2. Só quando nada é dado pelo nome único os **grupos** decidem: permitir em
   qualquer papel permite.

Assim, um grupo pode ser permitido com uma parte dele negada pelo nome, ou um
grupo negado com uma parte dele permitida pelo nome.

## As permissões do admin

Cada parte do admin e o que pode ser pedido dela — abra um nó para ver o que
está dentro dele. Cada parte da documentação tem também o seu item
**Permissões**.

{%= admin.permissionsTree() %}

Peça a quem administra os papéis o que o seu trabalho precisa.
