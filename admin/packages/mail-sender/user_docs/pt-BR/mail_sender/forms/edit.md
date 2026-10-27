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
- Na aba **Respostas**, **O mesmo do SMTP** esconde os campos que ele aproveita
  da conta SMTP (veja abaixo).
- A senha salva vem no campo, oculta (o ícone do olho a mostra); apagá-la
  a remove.

## A aba Respostas

**Quem lê as respostas** decide o resto da aba:

- **Pessoas, na caixa da equipe** — a aplicação não lê nada; as respostas
  ficam na caixa de e-mail da equipe. Nenhum outro campo aparece.
- **A aplicação, por IMAP** — a aplicação lê a caixa da equipe pelo IMAP e
  liga cada resposta à conversa de que ela é. Aparecem os campos abaixo.

### O mesmo do SMTP

Chips — marque os que a caixa que recebe tem em comum com a conta que envia
(aba **Envio**, método **SMTP**). Cada chip marcado faz a leitura usar o valor
da conta SMTP e esconde o campo dele no IMAP, no formulário e no detalhe:

| Chip | A leitura usa | Some do formulário |
| --- | --- | --- |
| **Usuário** | o usuário da conta SMTP | **Usuário** |
| **Senha** | a senha da conta SMTP (no Gmail, a senha de app) | **Senha** |
| **Servidor** | o servidor da conta SMTP — quando o mesmo servidor atende o IMAP | **Servidor IMAP** |

- O `?` à direita de cada chip, ao passar o mouse ou clicar, diz o que ele faz;
  clicar no `?` não marca nem desmarca o chip.
- O que estava num campo escondido fica guardado, sem uso, e volta a valer se
  o chip for desmarcado.
- A porta, o TLS, o STARTTLS e a caixa de entrada são sempre do IMAP: o IMAP
  usa outra porta que o SMTP.
- Com o método de envio **Gmail** não há conta SMTP a aproveitar: deixe os
  chips desmarcados e informe a conta do IMAP.
- No detalhe os três chips aparecem, os marcados com um ✓.

### Os campos do IMAP

- **Servidor IMAP** — o servidor da caixa, como `imap.gmail.com`.
- **Porta** — `0` é a de costume: 993 com TLS, 143 sem.
- **TLS** / **STARTTLS** — a criptografia da conexão: TLS desde o início (porta
  993, o que a maioria dos servidores pede) ou ligada depois de conectar
  (porta 143).
- **Usuário** e **Senha** — a conta da caixa da equipe (no Gmail, uma senha de
  app). A senha nunca é mostrada de novo depois de salva.
- **Caixa de entrada** — a pasta lida; vazia, `INBOX`.

### Exemplos

- **A mesma conta envia e recebe** (um servidor de e-mail próprio, SMTP e IMAP
  no mesmo host): marque **Usuário**, **Senha** e **Servidor**; informe só a
  porta e o TLS do IMAP.
- **Mesmo usuário e senha, servidores diferentes** (`smtp.exemplo.com` e
  `imap.exemplo.com`): marque **Usuário** e **Senha**; informe o
  **Servidor IMAP**.
- **Contas diferentes**: nenhum chip; preencha todos os campos do IMAP.

Confira a conta salva com
{%= admin.action("mail_sender", "TestSendMail").link %}.
