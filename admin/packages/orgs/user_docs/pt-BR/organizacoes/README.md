# Organizações

Uma organização guarda o que fica sob ela: cada recurso do menu sob a
organização — pessoas, financeiro… — é da organização escolhida, e os
registros de uma organização nunca são vistos de outra.

## A organização escolhida

A organização em que você trabalha fica no **topo do menu**: o chip dela diz
qual é, e os recursos dela ficam abaixo. Clique no chip para escolher outra;
o ✕ limpa a escolha. A escolhida é a da página em que você está, senão a
última que você escolheu neste navegador, senão a sua única, senão a
organização padrão (**Padrão**).

**Listar organizações** lista as organizações que você alcança — só essas.

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

## O proprietário confirma por e-mail

Nada é dado sem o **proprietário**: um compartilhamento novo, um papel que
pode mais que antes, um novo proprietário — quem quer que peça, o
proprietário recebe um **e-mail** com um link, válido por 48 horas, para uma
página em que ele **Confirma** ou **Recusa**, entrando como ele mesmo. Até
lá, nada muda. Um papel que pode menos que antes, e uma revogação, valem na
hora.

## O convite

Confirmado pelo proprietário, o compartilhamento com um **usuário** é um
**convite**: aparece nas **notificações** dele (o sino, no topo), com
**Aceitar** e **Recusar**. Aceito, a organização entra na lista dele e na
escolha do menu. O compartilhamento com um **grupo** não precisa de convite:
vale na hora para os membros.
