# Edit: {%= doc.model %}

The form of the mail sender. Choose the **sender** first: the fields of its
account follow it.

{%= admin.fields("edit") %}

## When it is available

The fields **Gmail** and **SMTP** are shown depending on the sender chosen:
**Gmail** when the sender is `GMAIL` — the account authorized with Google —,
**SMTP** when it is `SMTP` — the host, the port, the user and the password of
the account. The account of the other kind is kept, unused.

Check the sender saved with
{%= admin.action("mail_sender", "TestSendMail").link %}.
