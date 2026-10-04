# Permissões

Cada ação da página é uma permissão dela, dada aos papéis em **Papéis de
Usuário** (veja as políticas e permissões). A página é
{%= admin.page("/site-files").link %}, no menu em
**{%= admin.page("/site-files").menu %}**.

## As ações da página

| Ação | Permissão |
| --- | --- |
| Ver a página e os arquivos (só leitura) | `@get` |
| Commit | `!commit` |
| Atualizar | `!update` |
| Publicar | `!publish` |
| Descartar | `!discard` |
| Recomeçar | `!reset` |
| Pré-visualizar | `!preview` |

Assim, publicar pode ficar com menos gente que editar: quem edita faz o
commit, e outra pessoa, que revisa, publica.

## Os arquivos

Cada jeito de mudar os arquivos do rascunho é uma permissão própria — no
editor e no **WebDAV** igualmente. Quem tem só `@get` vê os arquivos e não
muda nenhum: o editor os abre só para leitura, e os botões do que o usuário
não pode fazer não aparecem.

| Ação | Permissão | No WebDAV |
| --- | --- | --- |
| Criar um arquivo ou uma pasta que não existe | `!create` | PUT de um arquivo novo, MKCOL, COPY |
| Mudar um arquivo que existe | `!edit` | PUT de um arquivo existente; locks, propriedades |
| Renomear um arquivo na sua pasta | `!rename` | MOVE na mesma pasta |
| Mover um arquivo para outra pasta | `!move` | MOVE para outra pasta |
| Excluir um arquivo ou uma pasta | `!delete` | DELETE |
| Importar: upload do computador, download da internet | `!import` | — |
| Clonar e enviar (push) pelo git (veja o passo **Pelo git**) | `!git` | — |

Importar pede `!import` **e** o que ele grava: `!create` para cada arquivo
novo, `!edit` para cada um que já existe (sobrescrito). O download da
internet é feito pelo servidor, e só de endereços públicos: nunca do próprio
servidor nem da rede dele.

## Permissões de um caminho

Qualquer permissão dos arquivos pode ser dada ou tirada **para um caminho e
tudo abaixo dele**. O caminho vai no recurso, entre `<` e `>`, antes da ação;
`*` é qualquer texto, então `<static/*>` é todo o `static/`, em todas as
subpastas.

As permissões da página são declaradas de dois jeitos — **com o grupo**, onde o
menu a põe agora, e **fora do grupo**, pelo nome único dela:

| Declarada | Recurso |
| --- | --- |
| com o grupo | `{%= admin.page("/site-files").resource %}` |
| fora do grupo | `{%= admin.page("/site-files").uniqueResource %}` |

A **de fora do grupo tem prioridade**: com as duas, decide ela. E continua
valendo se o menu mudar a página de grupo. Por isso as políticas abaixo são
declaradas fora do grupo:

| Política | Efeito |
| --- | --- |
| Negar `{%= admin.page("/site-files").uniqueResource %}<config/*>:!edit` | ninguém do papel edita em `config/` |
| Permitir `{%= admin.page("/site-files").uniqueResource %}<static/*>:*` | o papel faz qualquer coisa em `static/` |
| Permitir `{%= admin.page("/site-files").uniqueResource %}<static/img/*>:!import` | o papel importa imagens em `static/img/` |

Como decide, nesta ordem:

1. uma **negação do caminho** nega;
2. uma **permissão da página** (`…:!edit`) permite;
3. uma **negação da página** nega — diga o caminho o que disser;
4. sem nada dito da página, uma **permissão do caminho** permite.

Assim, para um papel que só pode mudar os estilos: dê `@get` na página e
`<static/css/*>:!edit` (e `!create`, se ele cria arquivos lá), e nada mais
dos arquivos.

## Exemplos por caminho

Cada linha é uma política de um papel: o efeito e a permissão, inteira — o
escopo e a página, fora do grupo (a mesma onde quer que o menu ponha a
página). Com o grupo, as mesmas linhas começam por `{%= admin.page("/site-files").resource %}`. Uma ação é `!create`, `!edit`, `!rename`, `!move`, `!delete` ou
`!import`; `!*` é todas elas, e `*` inclui ver (`@get`).

**Só os estilos: muda `static/css/` e mais nada**

| Efeito | Permissão |
| --- | --- |
| Permitir | `{%= admin.page("/site-files").uniqueResource %}@get` |
| Permitir | `{%= admin.page("/site-files").uniqueResource %}<static/css/*>:!edit` |
| Permitir | `{%= admin.page("/site-files").uniqueResource %}<static/css/*>:!create` |

**Tudo, menos `config/`, que fica só para leitura**

| Efeito | Permissão |
| --- | --- |
| Permitir | `{%= admin.page("/site-files").uniqueResource %}*` |
| Negar | `{%= admin.page("/site-files").uniqueResource %}<config/*>:!*` |

**Imagens: importa só em `static/img/`**

| Efeito | Permissão |
| --- | --- |
| Permitir | `{%= admin.page("/site-files").uniqueResource %}@get` |
| Permitir | `{%= admin.page("/site-files").uniqueResource %}<static/img/*>:!import` |
| Permitir | `{%= admin.page("/site-files").uniqueResource %}<static/img/*>:!create` |
| Permitir | `{%= admin.page("/site-files").uniqueResource %}<static/img/*>:!edit` |

**Os layouts nunca são excluídos nem movidos, pode o papel o que puder**

| Efeito | Permissão |
| --- | --- |
| Negar | `{%= admin.page("/site-files").uniqueResource %}<templates/layouts/*>:!delete` |
| Negar | `{%= admin.page("/site-files").uniqueResource %}<templates/layouts/*>:!move` |
| Negar | `{%= admin.page("/site-files").uniqueResource %}<templates/layouts/*>:!rename` |

**Renomear e mover só dentro de `static/`** — mover pede a permissão nos dois
caminhos, então para fora de `static/` (ou para dentro) é negado:

| Efeito | Permissão |
| --- | --- |
| Permitir | `{%= admin.page("/site-files").uniqueResource %}<static/*>:!rename` |
| Permitir | `{%= admin.page("/site-files").uniqueResource %}<static/*>:!move` |

**Um arquivo só**: o caminho inteiro, sem `*` — `{%= admin.page("/site-files").uniqueResource %}<config/layout_config.gad>:!edit`.
