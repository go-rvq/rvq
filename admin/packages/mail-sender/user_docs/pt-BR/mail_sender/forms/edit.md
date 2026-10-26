# Edição: {%= doc.model %}

A configuração do envio, em duas abas: **Envio** e **Respostas**. Uma aba
com um campo em erro fica marcada e é a que se abre ao salvar.

{%= admin.fields("edit") %}

## Quando está disponível

- Na aba **Envio**, os campos **Gmail** e **SMTP** acompanham o **Método de
  envio**: **Gmail** quando o método é `GMAIL` — a conta autorizada com o
  Google —, **SMTP** quando é `SMTP` — o servidor, a porta, o usuário e a
  senha da conta. A conta do outro método é mantida, sem uso.
- Na aba **Respostas**, o servidor IMAP, a porta, o TLS, o STARTTLS, o
  usuário, a senha e a caixa de entrada só aparecem quando **Quem lê as
  respostas** é **A aplicação, por IMAP**. Em **Pessoas, na caixa da
  equipe** a aplicação não lê nada, e esses campos ficam guardados sem uso.
- **Mesmo usuário e senha do SMTP**, marcado, esconde o usuário e a senha do
  IMAP (no formulário e no detalhe): a leitura usa o usuário e a senha da
  conta SMTP da aba **Envio**. O que estava neles fica guardado, sem uso, e
  volta a valer se a opção for desmarcada. Com o método de envio **Gmail**
  não há conta SMTP a usar: deixe desmarcado e informe a conta do IMAP.
- A senha salva vem no campo, oculta (o ícone do olho a mostra); apagá-la
  a remove.

Confira a conta salva com
{%= admin.action("mail_sender", "TestSendMail").link %}.
