# Caminhos estáveis

Todo endereço do admin — o que vem depois do seu prefixo, `/admin` — é de um
destes:

1. **Uma página do menu**: um modelo (`/admin/<modelo>`, `/admin/<modelo>/<id>`), um
   grupo do menu ou uma página. Um modelo **oculto** do menu continua sendo
   um modelo: o seu endereço funciona, só não aparece no menu. Esses
   endereços acompanham o menu: se um modelo muda de grupo, o seu endereço
   pode mudar.
2. **Um serviço**: um endereço que nenhuma página tem — de um programa, de um
   download, de uma API, de um protocolo, dos recursos das páginas. Ele tem
   uma de duas formas:
   - **`/admin/!<nome>/…`** — o nome do serviço depois de um **`!`**:
     `/admin/!backup/latest`, `/admin/!webdav/`. Nenhuma página começa por
     `!`: um serviço nunca se confunde com uma página, e o seu endereço
     **não muda** com o menu.
   - **um endereço cuja forma já diz o protocolo** — `/admin/site-files.git`
     (git por HTTP: o `.git` no fim). Nenhuma página tem um endereço assim, e
     o programa do outro lado (o `git`) precisa dele como ele é.

Os endereços dos serviços são **estáveis**: um script, um `git remote`, um
cliente de WebDAV, um link num e-mail podem guardá-los. Um endereço antigo
não é redirecionado: quando um serviço muda de endereço, a nota da versão
diz qual é o novo.

## Como cada um é acessado

- **público**: qualquer um (os recursos estáticos das páginas).
- **sessão**: quem está conectado ao admin, com as permissões da sua conta.
- **link assinado**: só quem recebeu o link, por um tempo.
- **chave de acesso**: um programa, com uma chave de acesso no cabeçalho
  `Authorization: Bearer <chave>` — sem a sessão do admin.
- **token**: um token que o próprio serviço deu (por exemplo, a sessão de
  um download de backup).
- **HTTP Basic**: login e senha — ou uma chave de acesso no lugar da senha —
  pedidos pelo próprio programa (`git`, um cliente de WebDAV).

Um serviço com **login próprio** é alcançado antes do login do admin: um
cliente o usa sem passar pela tela de login.

## Os serviços deste admin

{%= admin.services() %}
