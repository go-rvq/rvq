# Papéis

Um papel é um conjunto de permissões: o que seus usuários podem ver e fazer
no admin. Um usuário tem um ou mais papéis (veja
{%= admin.model("users").link %}).

![Os papéis](images/listing.png)

- O papel **Administrador** pode fazer tudo.
- Crie um papel para cada tipo de trabalho — editar o site, responder ao
  formulário de contato… — e dê a ele só o que esse trabalho precisa.

## As permissões

As permissões de um papel são uma **lista em ordem**, uma por linha
(arraste para mudar a ordem): um recurso, que **permite**, ou com um prefixo
— `deny ` **nega**, `allow ` força a permissão do que casa com ela, por
cima das negações anteriores, `reset ` tira das linhas anteriores o que casa
com ela. Um grupo é várias
linhas: `admin:users:{*, deny @delete}`. Tudo isso, com exemplos, em
{%= admin.doc("guides/policies/01-permissions").link %}.

Os papéis de um usuário também têm **ordem**: as linhas de um papel vêm
depois das do anterior, então ele pode negar o que o anterior permite, ou
permitir (`allow`) o que ele nega. A página **Permissões** do usuário mostra
o resultado.

Uma mudança num papel vale na hora para todos os usuários dele.

## Papéis do sistema

Alguns papéis são do **sistema**: a aplicação precisa deles, e os cria
sozinha quando faltam — a listagem e o detalhe dizem *Papel do sistema* e
para que cada um serve.

- Um papel do sistema **não é excluído** nem **renomeado**.
- As **permissões** dele podem ser mudadas como as de qualquer papel, e
  voltam às originais com a ação **Restaurar as permissões**, no detalhe.
- Um papel do sistema cujas permissões são as do próprio sistema — o
  **Administrador**, que pode tudo — não tem permissões para mudar.
- Um papel de mesmo nome que já existia passa a ser o do sistema: as
  permissões dele são mantidas, e as originais que lhe faltam, somadas.
