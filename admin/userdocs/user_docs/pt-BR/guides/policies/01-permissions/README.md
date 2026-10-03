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
| o escopo: quem pede — sempre o primeiro; o do admin é `admin` (veja abaixo) | `admin:…` |
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

    admin:site/:seo/:seo_config:<7>:@edit     pelos grupos
    admin:seo_config:<7>:@edit                pelo nome único

`*` vale por qualquer coisa: `admin:site/:*` é tudo o que está no grupo
*site*; `admin:seo_config:*`, tudo da parte *seo_config*, onde quer que o
menu a ponha.

## Os escopos

A primeira parte de um recurso é o **escopo**: a parte do sistema que pede a
permissão. O do admin é **`admin`** — por isso os recursos dele começam
por `admin:` (`admin:content/:posts:<12>:@edit`) —, e contém tudo
o que está no menu: os grupos, as partes, os registros delas, campos, seções,
ações e páginas. A biblioteca de mídia e as tarefas também são partes dele:

| O quê | Recurso |
| --- | --- |
| enviar um arquivo, num campo de imagem ou de arquivo | `admin:media_libraries:@create` |
| excluir um arquivo / editar a descrição dele | `admin:media_libraries:<5>:@delete` / `admin:media_libraries:<5>:@edit` |
| executar uma tarefa de um tipo (criar, rodar de novo, abortar) | `admin:jobs:!upload_posts` |

Um sistema feito sobre o admin pode ter escopos próprios, um nome antes do
primeiro `:`; a árvore abaixo mostra cada escopo, o do admin como *Administração*.

## Qual decide

1. O que é dado pelo **nome único** decide primeiro: **negar** — em qualquer
   papel do usuário — nega; permitir permite.
2. Só quando nada é dado pelo nome único os **grupos** decidem: permitir em
   qualquer papel permite.

Assim, um grupo pode ser permitido com uma parte dele negada pelo nome, ou um
grupo negado com uma parte dele permitida pelo nome.

## As permissões do admin

Cada parte do admin e o que pode ser pedido dela — abra um nó para ver o que
está dentro dele. Dentro de um modelo — e dos registros dele — o que ele
tem fica em grupos: permissões, ações, campos, seções, páginas, modelos
filhos e permissões próprias, cada um com o seu recurso. Cada parte da
documentação tem também o seu item **Permissões**.

{%= admin.permissionsTree() %}

Peça a quem administra os papéis o que o seu trabalho precisa.
