# Test send mail

Sends a test e-mail, through the method chosen in the action's form (Gmail
or the SMTP set), to the address given, with the subject and message given:
a quick way to check the account before the site relies on it. An error
says why the e-mail could not be sent (a wrong password, a blocked port…).

The test e-mail does not count toward the limits per address. With
`RVQ_MAIL_SMTP` set (development), it goes to the test server, not to the
real address.
