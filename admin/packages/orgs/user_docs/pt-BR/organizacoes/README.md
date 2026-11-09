# Organizações

Uma organização guarda o que fica sob ela: cada recurso do menu sob a
organização — pessoas, financeiro… — é da organização escolhida, e os
registros de uma organização nunca são vistos de outra.

## A organização escolhida

A organização em que você trabalha fica no **topo do menu**: o chip dela diz
qual é, e os recursos dela ficam abaixo. Clique no chip para escolher outra;
o ✕ limpa a escolha. A escolhida é a da página em que você está, senão a
última que você escolheu neste navegador, senão a sua única, senão **a sua
organização padrão**.

A sua padrão é só sua — qualquer organização que você alcança, sua ou
compartilhada com você: no detalhe dela,
{%= admin.action("organizacoes", "set_my_default").link %} (e
{%= admin.action("organizacoes", "unset_my_default").link %} para não ter
nenhuma). **Listar organizações** lista as organizações que você alcança —
só essas —, **Minha padrão** marcando a sua; uma de outro, compartilhada
com você, diz **compartilhada por** e o proprietário.

## Quem alcança uma organização

Uma organização tem um **proprietário**, que alcança tudo nela, e é
**compartilhada** com usuários e grupos, cada um com um **papel**:

| Papel | Pode |
|---|---|
| **Viewer** | ver |
| **Editor** | ver e alterar |
| **Manager** | ver, alterar e excluir |
| **Administrator** | tudo isso, e compartilhar a organização |
| **None** | nada: tira um registro do que está compartilhado acima dele |

Um papel dado numa organização vale também para **o que fica sob ela**; um
dado num registro abaixo dela vale mais que ele (o mais específico vence).
Um usuário alcança pelos próprios acessos e pelos dos seus **grupos**: quem
entra num grupo alcança na hora o que o grupo alcança, e quem sai dele perde
na hora.

No detalhe de uma organização (o menu ⋮ dela):

- {%= admin.action("organizacoes", "share_access").link %}: dá a um usuário
  ou grupo um papel;
- {%= admin.action("organizacoes", "accesses").link %}: quem a alcança, e os
  papéis — troque ou revogue;
- {%= admin.action("organizacoes", "access_history").link %}: tudo o que
  aconteceu com o compartilhamento dela;
- {%= admin.action("organizacoes", "transfer_owner").link %}: dá a ela outro
  proprietário.

## Compartilhar só um registro sob a organização

Um registro sob a organização — um projeto, por exemplo — pode ser
compartilhado sozinho, sem a organização. Quem o recebe o alcança (como o
papel diz) e, a não ser que **Ver o pai e os registros de apoio dele** tenha
sido desligado ao compartilhar, lê a organização: na lista e no menu dele
(**compartilhada por**), o nome e a descrição, e os registros **de apoio**
que a aplicação declara — nomes, códigos, o que o registro precisa (as
moedas, os bancos…), só alguns dos campos; dos ligados ao registro, só os
dele (as carteiras, os cartões dele). Os outros registros compartilhados
assim (os outros projetos) e o resto da organização ficam de fora.

## O proprietário confirma por e-mail

Nada é dado sem o **proprietário**: um compartilhamento novo, um papel que
pode mais que antes, um novo proprietário — quem quer que peça, o
proprietário recebe um **e-mail** com um link, válido por 48 horas, para uma
página em que ele **Confirma** ou **Recusa**, entrando como ele mesmo. Até
lá, nada muda. Um papel que pode menos que antes, uma revogação e um
compartilhamento como **None** valem na hora.

## O convite

Confirmado pelo proprietário, o compartilhamento com um **usuário** é um
**convite**: aparece nas **notificações** dele (o sino, no topo), com
**Aceitar** e **Recusar**. Aceito, a organização entra na lista dele e na
escolha do menu. O compartilhamento com um **grupo** não precisa de convite:
vale na hora para os membros.
