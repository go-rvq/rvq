# Edição: {%= doc.model %}

O formulário do remetente de e-mail. Escolha primeiro o **remetente**: os
campos da conta dele acompanham a escolha.

{%= admin.fields("edit") %}

## Quando está disponível

Os campos **Gmail** e **SMTP** aparecem conforme o remetente escolhido:
**Gmail** quando o remetente é `GMAIL` — a conta autorizada com o Google —,
**SMTP** quando é `SMTP` — o servidor, a porta, o usuário e a senha da conta.
A conta do outro tipo é mantida, sem uso.

Confira o remetente salvo com
{%= admin.action("mail_sender", "TestSendMail").link %}.
