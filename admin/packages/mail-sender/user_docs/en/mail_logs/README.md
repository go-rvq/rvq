# Mails sent

Every e-mail the site sends, or tries to, is recorded here: the messages of
the contact forms, the confirmations, the notices of comments, the password
resets, the test e-mail. It is only to look at: nobody edits, creates or
deletes a record (the old ones go by the retention, below).

**The listing** shows the newest first, with the status (**Sent** or
**Failed**), the kind, whom to, the subject and the method. The **Failed** tab
gathers the ones that did not go out. The filters choose by kind, period,
recipient and status; the search finds by the subject, the recipient, the
reference of the record of origin and the Message-ID.

**The detail** of each e-mail shows:

- when, the status and, when it failed, **why** (a wrong password, the limit
  per address, sending not configured…);
- **Origin**: the button **Open origin** goes to the record the e-mail is
  of — the contact message, the answer, the user;
- whom to, reply-to, the subject and the method (Gmail, SMTP, or the test
  server in development);
- the Message-ID and the message it answers: the **conversation**. Below the
  body, the other e-mails of the same conversation;
- who made it be sent, from which IP and browser, and the language;
- the **body**, as it was sent, isolated from the page.

## What is never kept

Sensitive content is hidden before it is written and shows as `••••••`: the
values the sender marks as sensitive (a password reset link, a phone) and, in
any link, the value of parameters such as `token`, `key`, `code` and
`password`. Some e-mails do not keep the body at all (**[content hidden]**).
An e-mail with anything hidden cannot be sent again: what it had is not here
any more.

## For how long

The settings of {%= admin.model("mail_sender").link %}, tab **Sending**, say
for how many days each e-mail is kept (365 by default; `0`, for ever). The
older ones are deleted once a day.
