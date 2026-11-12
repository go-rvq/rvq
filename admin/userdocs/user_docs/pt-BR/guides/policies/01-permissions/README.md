# 1. Permissões

O que um usuário pode fazer vem dos **papéis** dele, em ordem: cada papel
permite ou nega, parte por parte do admin, listar, ver, criar, editar e
excluir registros, e executar ações (veja *A pilha de permissões*, abaixo).

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
| um registro: o id entre `<…>` — um rótulo depois do id só o mostra, não decide nada | `<7>`, `<7<Primeiro post>>` — `<*>` qualquer um |
| um campo | `#Title` |
| um registro editado no lugar, num campo: `&` e o campo; os campos dele embaixo | `&Config`, `&Config:#Title` |
| uma seção do detalhe | `$Main` |
| uma página | `/report` |
| o que se pede de uma página: `@` e o método HTTP, um para cada que ela atende | `/report:@get`, `/report:@post` |
| uma permissão: `@` e o nome | `@list`, `@get`, `@create`, `@edit`, `@delete`; numa parte com lixeira, `@trash` (navegar nela) e, no registro, `@restore` (trazê-lo de volta) |
| uma ação: `!` e o nome | `!publish` |

Uma parte é alcançada de dois jeitos — **pelos grupos** e **pelo nome
único**, o dela:

    admin:site/:seo/:seo_config:<7>:@edit     pelos grupos
    admin:seo_config:<7>:@edit                pelo nome único

Uma **página** também: o nome único dela é o seu caminho, então uma permissão
dela pode ser declarada **fora do grupo** — sem grupo nenhum:

    admin:site/:/site-files:!edit             pelo grupo site
    admin:/site-files:!edit                   pelo nome único (sem grupo)

O nome único segue, parte por parte, para o que está abaixo dele: um
registro, um campo, uma ação, um caminho (`admin:/site-files:<static/*>:!edit`).

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
grupo negado com uma parte dele permitida pelo nome. Uma permissão declarada
**fora do grupo** (`admin:/site-files:!edit`) tem prioridade sobre a declarada
junto com o grupo (`admin:site/:/site-files:!edit`): com as duas, decide a de
fora do grupo. E ela continua valendo quando o menu muda a parte de grupo — a
declarada com o grupo deixa de casar.

Por isso o **campo de permissões** (abaixo) grava sempre **fora do grupo** o
que é dado a uma parte — um modelo, uma página, e o que está neles —, mesmo
escolhida dentro de um grupo: `admin:content/:posts:*` é gravado
`admin:posts:*`. O que é dado **ao grupo** fica com o grupo e vale para o que
está nele: `admin:content/:*`, `admin:content/:sub/:*`. As permissões gravadas
antes dentro de grupos são regravadas assim quando a aplicação sobe.

## A pilha de permissões

As permissões de um papel são uma **lista em ordem**, uma por linha. As de um
usuário são uma **pilha**: as do sistema primeiro, depois as dadas a ele num
registro (um compartilhamento), depois as de **cada papel dele, na ordem dos
papéis** — a ordem em que estão no cadastro do usuário (arraste para mudar).

Cada linha é um recurso, e pode começar com um **prefixo** — a palavra **e
um espaço**:

| Linha | O que faz |
| --- | --- |
| `admin:pages:*` | **permite** o que casa com ela, se nenhuma negação casar |
| `deny admin:users:*` | **nega** o que casa com ela |
| `allow admin:users:@list` | as negações **anteriores** deixam de valer para o que casa com ela — mesmo uma mais larga —, e ela permite: **força** a permissão só do que casa com ela |
| `reset admin:users:*` | as linhas **anteriores** deixam de valer para o que casa com ela — mesmo uma mais larga —; ela não permite nem nega nada |

Uma linha **cobre** outra quando o recurso dela casa com o da outra:
`admin:*` cobre `admin:users:@list`; `admin:users:@list` não cobre
`admin:users:*`. Uma linha depois de outra que ela cobre a **substitui**:

- `deny` e `allow` substituem as anteriores sem prefixo, `allow` e `deny`;
- `reset` substitui qualquer uma;
- uma sem prefixo substitui as anteriores sem prefixo (e `allow`).

Um `allow` ou um `reset` **mais estreito** que uma linha anterior — que ela
cobre — não a substitui: tira dela só o que ele casa. Depois de
`deny admin:users:*`, um `allow admin:users:@list` libera só a lista; a
negação continua valendo para o resto dos usuários. Uma linha sem prefixo
mais estreita não tira nada: a negação ganha dela.

O que sobra são as permissões **ativas**: as sem prefixo (um `allow` vira
uma) e as `deny`. Elas decidem como sempre: um recurso é permitido quando
uma ativa sem prefixo casa com ele e nenhuma `deny` ativa casa.

**Grupos** são várias linhas, na ordem deles: `admin:users:{@list,@create}`
é `admin:users:@list` e `admin:users:@create`. Uma alternativa pode ter o
seu prefixo — `admin:users:{*, deny @delete}` é `admin:users:*` e depois
`deny admin:users:@delete`; `admin:users:{*, deny {@edit,@delete}}` nega
só as duas —, e isso vale em qualquer parte do recurso: partes, campos,
ações. Uma quebra de linha dentro de um grupo separa como a vírgula.
Várias linhas na raiz são uma lista — `admin:posts:*, site:*` é
`{admin:posts:*, site:*}` —, e um espaço entre o `:` e o grupo é permitido
(`admin: {…}`).

Exemplos (a pilha de um usuário, de cima para baixo):

| Pilha | Pode | Não pode |
| --- | --- | --- |
| `admin:*`<br>`deny admin:users:*` | tudo, menos os usuários | os usuários |
| `admin:*`<br>`deny admin:users:*`<br>`allow admin:users:*` | tudo: o `allow` desfaz o `deny` | — |
| `admin:*`<br>`deny admin:users:*`<br>`reset admin:users:*` | tudo: o `reset` tira o `deny`, e o `admin:*` vale | — |
| `admin:users:*`<br>`reset admin:*`<br>`admin:pages:*` | só as páginas: o `reset` tirou os usuários | o resto |
| `admin:*`<br>`deny admin:users:*`<br>`allow admin:users:@list` | tudo, e a lista dos usuários: o `allow` força só ela | os usuários, fora a lista |
| `admin:*`<br>`deny admin:users:*`<br>`admin:users:@list` | tudo, menos os usuários: sem o `allow`, a negação ganha | os usuários, a lista também |
| `admin:*`<br>`reset admin:users:@delete` | tudo, menos excluir usuários: o `reset` tirou isso do `admin:*` | excluir usuários |
| `admin:users:{*, deny @delete}` | os usuários, menos excluir | excluir usuários |

Assim, um papel que vem depois pode negar o que um anterior permite (`deny`)
ou permitir o que um anterior nega (`allow`).

A página **Permissões** do detalhe de cada usuário mostra a pilha dele: na
aba **Tudo**, cada linha com o seu índice, de onde veio (o sistema, um papel)
e se está ativa — ou ativa menos o que uma linha depois casa —, ou foi
substituída — e por qual —; na aba **Ativas**, só as que valem, cada uma
com o que ela não cobre mais ("menos …"), em lista ou em árvore.

## O campo de permissões

As permissões de um papel, ou de uma chave de acesso, são editadas num campo
só: a **árvore das permissões do admin** reduzida ao que elas tocam — o que
permitem (coluna **Ativa**), o que tem uma linha própria, e o que está acima
disso —, com as mesmas abas (**Por verbo**, **Por tipo**), a busca e o filtro
de **Tipos** da árvore do usuário. Na barra, os ícones:

- **grupos do menu**: ligados, os grupos aparecem na tabela e no diálogo, cada
  parte dentro do seu grupo; desligados, as partes ficam direto sob o escopo —
  o que se grava é o mesmo (fora do grupo, veja *Qual decide*);
- **só as definidas**: só as linhas próprias e o que está acima delas, e não
  tudo o que permitem;
- **abrir tudo** / **fechar tudo**.

Cada linha da tabela:

- o **prefixo**, num seletor: **—** (não tem linha própria; escolher um cria a
  linha), **CONCEDE** (sem prefixo), **NEGA** (`deny`), **LIBERA** (`allow`),
  **REDEFINE** (`reset`). Num ramo (uma parte, um registro, um grupo), vale
  para tudo o que está nele (`admin:posts:*`);
- **+**: o que pode vir depois dela — e, numa parte, os registros, buscados
  pelo id ou pelo título;
- a **lixeira**: tira a linha dela e as de tudo o que está abaixo;
- o botão de abrir fica na **cor primária** quando há algo definido nela ou
  abaixo dela.

Uma linha que só nomeia uma parte sem dizer o quê (`admin:users:`) não decide
nada: aparece como incompleta enquanto se edita, e **não é gravada**.

**Incluir a partir de texto**: permissões escritas em qualquer forma — uma por
linha, por vírgulas, indentadas como no painel — juntadas às do campo; uma de
mesmo recurso que uma existente toma o lugar dela. `unset ` antes de uma tira
as linhas anteriores que ela cobre, qualquer que seja o prefixo delas
(`unset admin:posts:*`); nunca é gravado. Uma linha mais larga não é tirada
(`admin:*` não sai com `unset admin:posts:*`: use `deny`) — o campo avisa.

O painel **Permissão em texto**, ao lado da tabela (abaixo dela numa tela
estreita), mostra as permissões indentadas, para ler e copiar; o ícone dele o
recolhe.

## As chaves de acesso

As permissões de uma **chave de acesso** só aceitam **negações**
(`deny admin:users:*`), aplicadas por cima do que o usuário dela pode: a
chave pode tudo o que o usuário pode, menos o que ela nega. Sem nenhuma, tudo
o que o usuário pode. Os exemplos estão em
{%= admin.model("my_access_keys").link %}.

## As permissões do admin

Todas as permissões do admin, parte por parte e até as folhas — cada uma
como uma linha da pilha a escreve —, estão em
{%= admin.doc("guides/policies/01-permissions/01-tree").link %}. Cada parte
da documentação tem também o seu item **Permissões**.

Peça a quem administra os papéis o que o seu trabalho precisa.
