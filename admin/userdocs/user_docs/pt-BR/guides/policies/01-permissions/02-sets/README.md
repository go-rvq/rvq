# 1.2 Conjuntos de permissões

Um **conjunto** é um nome para permissões dadas juntas: ler, escrever,
excluir. Ele entra no fim de uma linha, depois de `%`:

    admin:posts:%readers           lista as postagens e vê cada uma
    deny admin:users:%editors      nega o que editors dá aos usuários

A linha guarda **o nome**, não o que ele dá. O que ele dá é calculado de novo
cada vez que as permissões de um usuário são montadas. Por isso:

- se o conjunto **muda**, toda linha dele muda junto;
- se ele **deixa de existir**, ou não dá nada àquele tipo de recurso, a linha
  fica guardada e é **ignorada**: não permite nem nega nada, e aparece em
  vermelho no campo de permissões.

## O que um conjunto dá

Um conjunto diz o que dá **por tipo de recurso**: `model` (uma parte com
registros), `singleton` (uma parte de um registro só), `page` (uma página),
`check` (uma permissão própria). Por exemplo:

    readers  {model: {@list, <*>:@get}, singleton: @get, page: @get}

Dado a `admin:posts:`, que é um *model*, ele dá `admin:posts:@list` e
`admin:posts:<*>:@get`. Dado a uma página, dá o `@get` dela.

Uma linha `deny` de um conjunto nega tudo o que ele dá. Um conjunto pode
conter outros (`editors` é `{%authors, %removers}`).

Um recurso pode ter a **sua própria versão** de um nome. A dos arquivos do
site, por exemplo, diz o que *publishers* é ali. Ela vale no lugar da versão
do tipo.

## Os conjuntos padrão

| Conjunto | O que dá |
| --- | --- |
| `%readers` — Leitura | lista e vê; no registro único e nas páginas, vê |
| `%authors` — Autoria | lista, cria, vê e edita; não exclui |
| `%removers` — Exclusão | a lixeira, excluir e restaurar registros |
| `%editors` — Edição | autoria e exclusão juntas |
| `%managers` — Gestão | tudo |

Pacotes e aplicações acrescentam os seus, como a tradução, a moderação de
comentários e o histórico.

## No campo de permissões

O **+** de uma linha da árvore mostra, abaixo das partes, a seção
**Conjuntos de permissões**: os que dão algo àquele recurso, cada um com o que
dá ali. Escolher um insere a linha do conjunto (`…:%readers`).

Na árvore, a linha aparece sob o seu recurso, com o chip do nome. Passar o
mouse sobre o chip mostra o que ele concede. O que o conjunto permite fica
ativo nas linhas da árvore, como o de qualquer linha.

## Conjuntos do admin

Em **Administrador › Conjuntos de permissões** criam-se conjuntos novos:

- um **nome**: letras minúsculas, algarismos e `_`, começando por letra;
- a **definição**, por tipo de recurso, como acima;
- um **título** e uma **descrição** em cada idioma do admin.

O nome não pode ser o de um conjunto do sistema. Salvar ou excluir um
conjunto refaz as permissões dos usuários no próximo acesso.
