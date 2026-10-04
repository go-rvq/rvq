# Opções dos layouts

As opções que uma página, um tipo de post ou um post mostram no admin são
escritas em `config/layout_config.gad`: cada conjunto de opções é uma
**classe**, cada opção um **campo** dela. Salvar e publicar o arquivo muda os
formulários do admin — sem reiniciar nada.

```gad
class Banner {
    [label="Título", hint="Aparece sobre a imagem"]
    title str

    [label="Altura"]
    height int = 400
}
```

- `[label="…", hint="…"]` antes de um campo: o nome que o formulário mostra e
  a linha de ajuda abaixo dele.
- `nome tipo`: `str` (texto), `int`, `bool` (um interruptor), `float`;
  `= valor` é o padrão; `nome?` pode ficar vazio.

## De quem são as opções de cada classe

O arquivo termina nas classes que o admin lê:

- **Páginas** — `PageConfig`: o `Layout` dela é a escolha do layout da
  página, cada layout uma classe com as suas opções (`Layout Default |
  PostList | …`).
- **Tipos de postagem** — `PostTypeConfig`: o `PostLayout` dela lista as
  classes que um tipo pode escolher **para as suas postagens** (`PostLayout?
  ServiceOptions | PortfolioOptions`).
- **Postagens** — não têm classe própria aqui: **as opções de uma postagem são
  a classe que o tipo de postagem dela escolheu**. Em **Tipos de postagem**,
  as opções do tipo escolhem o layout das postagens dele; então toda postagem
  desse tipo mostra esses campos nas suas opções. Um tipo que não escolheu
  nenhum: as postagens dele não têm opções.

Por isso, para dar novas opções às postagens de um tipo, acrescente os campos
à classe que o tipo escolheu (ou escreva uma classe nova, acrescente-a ao
`PostLayout` e escolha-a no tipo). Trocar a classe que um tipo escolheu deixa
fora do formulário as opções que as postagens dele salvaram com a antiga.

## Uma escolha entre valores: um select

Um campo que guarda **um de uma lista fechada** de valores é um select. Dê a
ele os valores, e o que o formulário mostra para cada um, com `options`:

```gad
class Banner {
    [label="Alinhamento", options=(;left="À esquerda", center="Centralizado", right="À direita")]
    align str

    [label="Colunas", options=[[2, "Duas"], [3, "Três"], [4, "Quatro"]]]
    columns int = 3

    [label="Estilo", options=["light", "dark"]]
    style? str
}
```

`options` se escreve de um destes três jeitos, na ordem em que o select
mostra:

| Escrito como | Cada item | Exemplo |
|---|---|---|
| um key-value array `(;…)` | valor `=` rótulo | `(;left="À esquerda")` |
| um array de pares | `[valor, rótulo]` | `[[2, "Duas"]]` |
| um array de valores | o valor é o próprio rótulo | `["light", "dark"]` |

- O que se **salva** é o **valor** (`left`, `3`); o rótulo só aparece.
  Trocar um rótulo é seguro; trocar um valor deixa o que já foi salvo com o
  antigo fora da lista.
- O formulário só aceita um valor da lista. Um campo sem `?` precisa ser
  escolhido; com `?` pode ficar vazio.
- Funciona com qualquer tipo: `columns int` continua salvando um número.
- Numa lista de valores (`tags []str`), cada item é um select das opções.

Um `enum` declarado no arquivo também é um select — seus membros são os
valores, cada um mostrado pelo próprio nome:

```gad
enum Align { left, center, right }

class Banner {
    align Align
}
```

Use `options` quando os valores precisam de rótulos próprios, ou são números.

## Quando o arquivo está errado

Um valor repetido, um par sem o rótulo (`[["a"]]`), um `options` que não é
nenhum dos três jeitos, ou `options` num campo que é um grupo de campos
(`interface {…}`) é um erro: o admin o mostra no lugar dessas opções até o
arquivo ser corrigido. Confira o formulário antes de publicar.
