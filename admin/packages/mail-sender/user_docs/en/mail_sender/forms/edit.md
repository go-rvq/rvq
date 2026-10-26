# Edit: {%= doc.model %}

The sending settings, in two tabs: **Sending** and **Answers**. A tab with
a field in error is marked, and is the one open after saving.

{%= admin.fields("edit") %}

## When it is available

- In the **Sending** tab, the fields **Gmail** and **SMTP** follow the
  **sending method**: **Gmail** when the method is `GMAIL` — the account
  authorized with Google —, **SMTP** when it is `SMTP` — the host, the port,
  the user and the password of the account. The account of the other method
  is kept, unused.
- In the **Answers** tab, the IMAP server, port, TLS, STARTTLS, user,
  password and mailbox are shown only when **Who reads the answers** is
  **The application, by IMAP**. With **People, in the team mailbox** the
  application reads nothing, and those fields are kept unused.
- **Same user and password as the SMTP**, marked, hides the IMAP user and
  password (in the form and in the detail): the reading uses the user and
  password of the SMTP account of the **Sending** tab. What they held is kept,
  unused, and counts again if the option is unmarked. With the sending
  method **Gmail** there is no SMTP account to use: leave it unmarked and give
  the IMAP account.
- The password saved is in its field, hidden (the eye icon shows it);
  clearing it removes it.

Check the account saved with
{%= admin.action("mail_sender", "TestSendMail").link %}.
