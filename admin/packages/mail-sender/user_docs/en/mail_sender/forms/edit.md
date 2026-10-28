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
- In the **Answers** tab, **The same as the SMTP** hides the fields it takes
  from the SMTP account (see below).
- The password saved is in its field, hidden (the eye icon shows it);
  clearing it removes it.

## The subject parts

A row by type of message and **Others** (what has no row of its own); in each
column — **To the team**, **To outside** —, the boxes of the parts:

| Part | In the subject | Where it comes from |
| --- | --- | --- |
| **Site** | `[PAC]` | **Site name in the subject**; empty, the name of the site's settings |
| **Thing** | `[ISSUE - Issues]` | the type of message: the issue's type, the form — translatable in the site's messages (`mail_topic.*`) |
| **Record** | `[#12]` | the issue's number, the contact message's reference |

The title always goes, and an answer always begins with `Re:`. An issue with
everything marked: `Re: [PAC] [ISSUE - Issues] [#12] The page does not open`.

## The Answers tab

**Who reads the answers** decides the rest of the tab:

- **People, in the team mailbox** — the application reads nothing; the answers
  stay in the team's mailbox. No other field is shown.
- **The application, by IMAP** — the application reads the team's mailbox by
  IMAP and ties each answer to the conversation it is of. The fields below are
  shown.

### The same as the SMTP

Chips — mark what the mailbox that receives has in common with the account
that sends (tab **Sending**, method **SMTP**). Each chip marked makes the
reading use the value of the SMTP account and hides its IMAP field, in the
form and in the detail:

| Chip | The reading uses | Leaves the form |
| --- | --- | --- |
| **User** | the user of the SMTP account | **User** |
| **Password** | the password of the SMTP account (with Gmail, its app password) | **Password** |
| **Host** | the host of the SMTP account — when the same host serves IMAP | **IMAP server** |

- The `?` at the right of each chip, on hover or click, says what it does; a
  click on the `?` does not mark or unmark the chip.
- What a hidden field held is kept, unused, and counts again if the chip is
  unmarked.
- The port, TLS, STARTTLS and the mailbox are always the IMAP's: IMAP uses
  another port than SMTP.
- With the sending method **Gmail** there is no SMTP account to take from:
  leave the chips unmarked and give the IMAP account.
- In the detail the three chips are shown, the marked ones with a ✓.

### The IMAP fields

- **IMAP server** — the server of the mailbox, such as `imap.gmail.com`.
- **Port** — `0` is the usual one: 993 with TLS, 143 without.
- **TLS** / **STARTTLS** — the encryption of the connection: TLS from the start
  (port 993, what most servers ask for) or turned on after connecting (port
  143).
- **User** and **Password** — the account of the team mailbox (with Gmail, an
  app password). The password is never shown again once saved.
- **Mailbox** — the folder read; empty, `INBOX`.
- **After importing** — what is done with the answer imported (a comment now):
  **Mark read**, **Delete from the server** or **Move to a folder**. Only it is
  touched: deleting needs the server's UIDPLUS extension (`UID EXPUNGE` of it
  alone) and moving MOVE (or UIDPLUS); a server without them only marks it
  read, and the log says why — a plain `EXPUNGE` would also delete what
  another program marked deleted in the same mailbox.
- **Folder** — only when moving: where the answers imported go; empty, the
  domain of the sending's address (`pac.com`), made when it is not there.

### Examples

- **The same account sends and receives** (a mail server of your own, SMTP and
  IMAP on the same host): mark **User**, **Password** and **Host**; give only
  the IMAP's port and TLS.
- **Same user and password, different hosts** (`smtp.example.com` and
  `imap.example.com`): mark **User** and **Password**; give the **IMAP
  server**.
- **Different accounts**: no chip; fill every IMAP field.

Check the account saved with
{%= admin.action("mail_sender", "TestSendMail").link %}.
