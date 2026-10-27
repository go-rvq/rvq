# Envio de e-mails

Uma tela só, de configuração: a conta pela qual saem todos os e-mails do
site — as mensagens dos formulários de contato, as redefinições de senha, as
notificações — e quem lê as respostas a eles. Não há uma lista de
remetentes: o que se configura aqui vale para tudo.

![A configuração do envio de e-mails](images/listing.png)

A configuração tem duas abas, no detalhe e na edição:

- **Envio** — como os e-mails saem:
  - o **Método de envio**: **Gmail**, uma conta do Google autorizada pelo
    próprio painel, ou **SMTP**, um servidor e a sua conta (servidor, porta,
    usuário, senha). Só a conta do método escolhido é usada; a do outro fica
    guardada;
  - o **prefixo do assunto**, posto antes do assunto de todo e-mail;
  - os **limites por endereço**: quantos e-mails um mesmo endereço de fora
    da equipe (quem escreveu pelo site, um cliente) recebe por hora e por
    dia. Acima disso o e-mail é recusado e o erro fica registrado; os
    e-mails para a própria equipe não contam. Os limites impedem que um
    formulário do site seja usado para encher a caixa de alguém. `0`: sem
    limite;
  - por quantos dias a {%= admin.model("mail_logs").link %} guarda cada
    e-mail enviado (padrão 365; `0`, para sempre).
- **Respostas** — quem lê o que as pessoas respondem aos e-mails:
  - **Pessoas, na caixa da equipe**: as respostas chegam à caixa de e-mail
    da equipe e são lidas lá, como qualquer e-mail;
  - **A aplicação, por IMAP**: a aplicação lê a caixa pelo IMAP; a resposta
    a uma conversa que ela conhece (uma mensagem de contato, um comentário)
    passa a fazer parte dela, e a mensagem é marcada como lida na caixa. Os
    campos do servidor IMAP só aparecem neste modo. Quando a caixa que recebe
    é a mesma que envia, marque em **O mesmo do SMTP** o que a
    leitura aproveita da conta da aba **Envio** — o usuário, a senha, o
    servidor —; cada um marcado some do formulário (uma senha só para
    trocar). As opções e os exemplos estão na {%= admin.doc("mail_sender/forms/edit").link %}.

As senhas (da conta SMTP e da caixa IMAP) nunca são mostradas no detalhe.
Na edição a senha salva vem no campo, oculta (o ícone do olho a mostra):
apagá-la a remove.

Todo e-mail enviado, ou que falhou, fica na {%= admin.model("mail_logs").link %}, com o conteúdo sensível oculto.

Antes de o site depender da conta, confira-a com
{%= admin.action("mail_sender", "TestSendMail").link %}.

## Em desenvolvimento e nos testes

Com a variável de ambiente `RVQ_MAIL_SMTP` (por exemplo
`127.0.0.1:2525`), todo e-mail vai para esse servidor SMTP, seja qual for o
método configurado aqui — e nada sai da máquina. O servidor de testes do
rvq (`make mail-dev`) recebe ali, mostra cada mensagem que chega e a deixa
legível por IMAP (`127.0.0.1:1143`, qualquer endereço, senha `test`): é com
ele que o modo **A aplicação, por IMAP** é testado. Em produção a variável
não existe.
