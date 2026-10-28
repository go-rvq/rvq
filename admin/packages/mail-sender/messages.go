package mail_sender

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "rvq-admin/mail-sender"

type Messages struct {
	ModuleDescription                   string           `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	MailSender                          string           `i18n:"hint='Name of the mail sender settings (menu, title).'"`
	TestMessageSubject                  string           `i18n:"label='Test message: subject', hint='Subject of the test email.'"`
	TestMessageBody                     string           `i18n:"label='Test message: body', hint='Body of the test email.'"`
	MailSenderSender                    string           `i18n:"label='Sender', hint='Label of the kind of sender (Gmail, SMTP).'"`
	MailSenderSubjectPrefix             string           `i18n:"label='Subject prefix', hint='Label of the text put before the subject of every email.'"`
	MailSenderIDFooter                  string           `i18n:"label='IDs in the body', hint='Label of the option that puts the ids of a message at the end of its body.'"`
	MailSenderIDFooter_Hint             string           `i18n:"label='IDs in the body: hint', hint='Hint of the option of the ids in the body.'"`
	MailSenderSubjectRules              string           `i18n:"label='Subject parts', hint='Label of the table of the parts of the subject by type of message and to whom.'"`
	MailSenderSubjectRules_Hint         string           `i18n:"label='Subject parts: hint', hint='Hint of the table of the parts of the subject.'"`
	SubjectRulesFamily                  string           `i18n:"label='Subject parts: type column', hint='Column of the type of message.'"`
	SubjectRulesTeam                    string           `i18n:"label='Subject parts: to the team', hint='Column of the messages to the team.'"`
	SubjectRulesExternal                string           `i18n:"label='Subject parts: outside', hint='Column of the messages to outside the installation (a visitor, an author).'"`
	SubjectPartSite                     string           `i18n:"label='Subject part: site', hint='The part of the site name.'"`
	SubjectPartTopic                    string           `i18n:"label='Subject part: thing', hint='The part of the thing the message is of.'"`
	SubjectPartRef                      string           `i18n:"label='Subject part: record', hint='The part of the record (#12).'"`
	SubjectFamilyOther                  string           `i18n:"label='Subject parts: others', hint='The row of the types with no row of their own.'"`
	MailSenderGmail                     string           `i18n:"label='Gmail', hint='Name of the Gmail sender.'"`
	MailSenderSMTP                      string           `i18n:"label='SMTP', hint='Name of the SMTP sender.'"`
	MailSender_Action_TestSendMail      string           `i18n:"label='Send test mail', hint='Action that sends a test email.'"`
	GmailSenderConfiguredSuccessfully   string           `i18n:"label='Gmail configured', hint='Shown once the Gmail sender was authorized.'"`
	GmailSenderLogOutSuccessfully       string           `i18n:"label='Gmail logged out', hint='Shown once the Gmail sender was signed out.'"`
	GmailSenderConnectPage              string           `i18n:"label='Gmail: connect', hint='Page that sends to Google to authorize the Gmail sender (a permission of its own).'"`
	GmailSenderLogoutPage               string           `i18n:"label='Gmail: disconnect', hint='Page that signs the Gmail sender out (a permission of its own).'"`
	GmailSenderCallbackPage             string           `i18n:"label='Gmail: authorization return', hint='Page Google returns to after the authorization (a permission of its own).'"`
	GmailSenderSuccessPage              string           `i18n:"label='Gmail: authorized', hint='Page shown once the Gmail sender was authorized (a permission of its own).'"`
	GmailSenderCallbackURI              string           `i18n:"label='Gmail: callback URI', hint='Label of the address Google returns to after the authorization.'"`
	GmailSenderCredentials              string           `i18n:"label='Gmail: app credentials', hint='Label of the OAuth credentials of the Gmail sender.'"`
	GmailSenderCredentialsFile          string           `i18n:"label='Gmail: app credentials file', hint='Label of the field that uploads the credentials file.'"`
	GmailSenderCredentialsFile_Hint     string           `i18n:"label='Gmail: app credentials file hint', hint='Hint of the credentials file: the .json of a Desktop application from the Google Cloud Console.'"`
	GmailSenderSetup                    string           `i18n:"label='Gmail: setup', hint='Button that opens how to create the Gmail credentials.'"`
	GmailSenderSetupTitle               string           `i18n:"label='Gmail: setup title', hint='First paragraph of the Gmail setup: why the OAuth client is created by hand.'"`
	GmailSenderSetupSteps               []string         `i18n:"label='Gmail: setup steps', hint='The steps to create the Gmail credentials, one per item.'"`
	GmailSenderSetupScopes              string           `i18n:"label='Gmail: scopes', hint='Title of the scopes to declare on the consent screen.'"`
	GmailSenderSetupScopesHint          string           `i18n:"label='Gmail: scopes hint', hint='Why the scopes listed are the ones needed.'"`
	GmailSenderSetupScopesSteps         []string         `i18n:"label='Gmail: scopes steps', hint='The steps to declare the scopes, one per item.'"`
	GmailSenderSetupBlocked             string           `i18n:"label='Gmail: access blocked', hint='What to do when Google blocks the access because the app is not verified.'"`
	GmailSenderSetupCallback            string           `i18n:"label='Gmail: setup callback', hint='Title of the callback address of this installation.'"`
	GmailSenderSetupCallbackHint        string           `i18n:"label='Gmail: setup callback hint', hint='Why the callback address need not be registered.'"`
	GmailSenderSetupConsoleBtn          string           `i18n:"label='Gmail: open console', hint='Button that opens the Google Cloud Console.'"`
	SmtpSenderTLS                       string           `i18n:"label='SMTP: TLS', hint='Label of whether the SMTP connection uses TLS.'"`
	SmtpSenderServer                    string           `i18n:"label='SMTP: server', hint='Label of the SMTP server.'"`
	SmtpSenderPort                      string           `i18n:"label='SMTP: port', hint='Label of the SMTP port.'"`
	SmtpSenderFromMail                  string           `i18n:"label='SMTP: from', hint='Label of the address the emails are sent from.'"`
	SmtpSenderUser                      string           `i18n:"label='SMTP: user', hint='Label of the SMTP user.'"`
	SmtpSenderPassword                  string           `i18n:"label='SMTP: password', hint='Label of the SMTP password.'"`
	MailSenderMaxPerRecipientHour       string           `i18n:"label='Limit per address: hour', hint='Label of the most emails an address outside the team (a contact) gets in an hour; 0, no limit.'"`
	MailSenderMaxPerRecipientDay        string           `i18n:"label='Limit per address: day', hint='Label of the most emails an address outside the team (a contact) gets in a day; 0, no limit.'"`
	MailSenderInbound                   string           `i18n:"label='Answers', hint='Label of the group that says where the answers to the emails are read from.'"`
	InboundSettingsMode                 string           `i18n:"label='Answers: read by', hint='Label of who reads the answers: people, in the team mailbox, or the application, by IMAP.'"`
	InboundSettingsServer               string           `i18n:"label='Answers: IMAP server', hint='Label of the IMAP server the answers are read from.'"`
	InboundSettingsPort                 string           `i18n:"label='Answers: IMAP port', hint='Label of the IMAP port; 0, the default (993 with TLS, 143 without).'"`
	InboundSettingsTLS                  string           `i18n:"label='Answers: TLS', hint='Label of whether the IMAP connection uses TLS (port 993).'"`
	InboundSettingsStartTLS             string           `i18n:"label='Answers: STARTTLS', hint='Label of whether the IMAP connection uses STARTTLS.'"`
	InboundSettingsReuse                string           `i18n:"label='Answers: reuse of the SMTP', hint='Label of what the reading of the answers takes from the SMTP sender (user, password, host).'"`
	InboundSettingsUser                 string           `i18n:"label='Answers: user', hint='Label of the IMAP user (the team mailbox address).'"`
	InboundSettingsPassword             string           `i18n:"label='Answers: password', hint='Label of the IMAP password.'"`
	InboundSettingsAfter                string           `i18n:"label='Answers: after importing', hint='Label of what is done with a message imported by IMAP.'"`
	InboundSettingsAfter_Hint           string           `i18n:"label='Answers: after importing hint', hint='Hint of what is done with a message imported.'"`
	InboundSettingsMoveTo               string           `i18n:"label='Answers: folder', hint='Label of the folder the answers imported are moved to.'"`
	InboundSettingsMoveTo_Hint          string           `i18n:"label='Answers: folder hint', hint='Hint of the folder the answers imported are moved to.'"`
	InboundAfterRead                    string           `i18n:"label='After importing: read', hint='Choice: the message imported is marked read.'"`
	InboundAfterDelete                  string           `i18n:"label='After importing: delete', hint='Choice: the message imported is deleted from the server.'"`
	InboundAfterMove                    string           `i18n:"label='After importing: move', hint='Choice: the message imported is moved to a folder.'"`
	InboundSettingsMailbox              string           `i18n:"label='Answers: mailbox', hint='Label of the IMAP mailbox read; empty, INBOX.'"`
	InboundManual                       string           `i18n:"label='Answers: by people', hint='Choice: the answers are read by people, in the team mailbox.'"`
	InboundIMAP                         string           `i18n:"label='Answers: by the application', hint='Choice: the application reads the answers by IMAP.'"`
	MailSenderTabSending                string           `i18n:"label='Tab: sending', hint='Tab of how the emails are sent: the sender, its account, the limits.'"`
	MailSenderTabAnswers                string           `i18n:"label='Tab: answers', hint='Tab of how the answers to the emails are read.'"`
	MailSenderSender_Hint               string           `i18n:"label='Mail sender: sending method hint', hint='Hint of the sender: Gmail or SMTP; the fields of its account follow the choice.'"`
	MailSenderSubjectPrefix_Hint        string           `i18n:"label='Mail sender: subject prefix hint', hint='Hint of the text put before the subject of every email.'"`
	MailSenderMaxPerRecipientHour_Hint  string           `i18n:"label='Mail sender: limit per hour hint', hint='Hint of the hourly limit per address outside the team.'"`
	MailSenderMaxPerRecipientDay_Hint   string           `i18n:"label='Mail sender: limit per day hint', hint='Hint of the daily limit per address outside the team.'"`
	InboundSettingsMode_Hint            string           `i18n:"label='Answers: who reads them hint', hint='Hint of who reads the answers to the emails.'"`
	InboundSettingsServer_Hint          string           `i18n:"label='Answers: IMAP server hint', hint='Hint of the IMAP server.'"`
	InboundSettingsPort_Hint            string           `i18n:"label='Answers: IMAP port hint', hint='Hint of the IMAP port.'"`
	InboundSettingsTLS_Hint             string           `i18n:"label='Answers: IMAP TLS hint', hint='Hint of TLS of the IMAP connection.'"`
	InboundSettingsStartTLS_Hint        string           `i18n:"label='Answers: IMAP STARTTLS hint', hint='Hint of STARTTLS of the IMAP connection.'"`
	InboundSettingsReuse_Hint           string           `i18n:"label='Answers: reuse of the SMTP hint', hint='Hint of what the reading of the answers takes from the SMTP sender.'"`
	InboundSettingsUser_Hint            string           `i18n:"label='Answers: IMAP user hint', hint='Hint of the IMAP user.'"`
	InboundSettingsPassword_Hint        string           `i18n:"label='Answers: IMAP password hint', hint='Hint of the IMAP password.'"`
	InboundSettingsMailbox_Hint         string           `i18n:"label='Answers: IMAP mailbox hint', hint='Hint of the IMAP mailbox read.'"`
	SendMailTestFormTo                  string           `i18n:"label='Test: to', hint='Label of who the test email goes to.'"`
	SendMailTestFormSubject             string           `i18n:"label='Test: subject', hint='Label of the subject of the test email.'"`
	SendMailTestFormMessage             string           `i18n:"label='Test: message', hint='Label of the body of the test email.'"`
	SendMailTestFormSender              string           `i18n:"label='Test: sender', hint='Label of the sender the test email goes through.'"`
	ErrGmailSenderCredentialsInvalid    i18n.ErrorString `i18n:"hint='Error when the credentials file is not of a Desktop application.'"`
	ErrGmailSenderScopeNotGranted       string           `i18n:"hint='Error when Google gave a token without the scope needed.', fields=(;'%s'='the scope')"`
	SendMailSuccessfully                string           `i18n:"label='Mail sent', hint='Shown once an email was sent.'"`
	GmailSenderToken                    string           `i18n:"hint='Label of the token field of the Gmail sender.'"`
	GmailSenderUser                     string           `i18n:"hint='Label of the user field of the Gmail sender.'"`
	MailLog                             string           `i18n:"label='Central: name', hint='Name of a message of the central of the mails sent.'"`
	MailLogs                            string           `i18n:"label='Central: plural', hint='Name of the central of the mails sent (menu, title).'"`
	MailLogCreatedAt                    string           `i18n:"label='Central: when', hint='Label of when the message was sent.'"`
	MailLogStatus                       string           `i18n:"label='Central: status', hint='Label of how the sending went.'"`
	MailLogStatusSent                   string           `i18n:"label='Central: sent', hint='A message sent.'"`
	MailLogStatusFailed                 string           `i18n:"label='Central: failed', hint='A message that failed to be sent.'"`
	MailLogKind                         string           `i18n:"label='Central: kind', hint='Label of what the message is (contact.received…).'"`
	MailLogOrigin                       string           `i18n:"label='Central: origin', hint='Label of the record the message is of.'"`
	MailLogOpenOrigin                   string           `i18n:"label='Central: open origin', hint='Button that opens the record the message is of.'"`
	MailLogRecipients                   string           `i18n:"label='Central: to', hint='Label of whom the message went to.'"`
	MailLogReplyTo                      string           `i18n:"label='Central: reply-to', hint='Label of the address the answers go to.'"`
	MailLogSubject                      string           `i18n:"label='Central: subject', hint='Label of the subject of the message.'"`
	MailLogMethod                       string           `i18n:"label='Central: method', hint='Label of the sender the message went by.'"`
	MailLogError                        string           `i18n:"label='Central: error', hint='Label of why the message failed.'"`
	MailLogMessageID                    string           `i18n:"label='Central: Message-ID', hint='Label of the Message-ID of the message.'"`
	MailLogInReplyTo                    string           `i18n:"label='Central: In-Reply-To', hint='Label of the Message-ID the message answers.'"`
	MailLogUserID                       string           `i18n:"label='Central: user', hint='Label of who made the message be sent.'"`
	MailLogRemoteIP                     string           `i18n:"label='Central: IP', hint='Label of the IP of the request that made the message be sent.'"`
	MailLogUserAgent                    string           `i18n:"label='Central: browser', hint='Label of the browser of the request that made the message be sent.'"`
	MailLogLocaleCode                   string           `i18n:"label='Central: language', hint='Label of the language of the message.'"`
	MailLogBody                         string           `i18n:"label='Central: body', hint='Label of the body of the message (its sensitive content hidden).'"`
	MailLogBodyHidden                   string           `i18n:"label='Central: body hidden', hint='Shown in place of a body that was not kept.'"`
	MailLogConversation                 string           `i18n:"label='Central: conversation', hint='Label of the other messages of the conversation.'"`
	MailLogTabAll                       string           `i18n:"label='Central: tab all', hint='Tab of every message.'"`
	MailLogTabFailed                    string           `i18n:"label='Central: tab failed', hint='Tab of the messages that failed.'"`
	MailLog_Action_Resend               string           `i18n:"label='Central: send again', hint='Action that sends a message again, as it was.'"`
	MailSenderMailLogRetentionDays      string           `i18n:"label='Mail sender: retention', hint='Label of how long the central keeps the messages.'"`
	MailSenderMailLogRetentionDays_Hint string           `i18n:"label='Mail sender: retention hint', hint='Hint of how long the central keeps the messages.'"`
	SMTPReuseUser                       string           `i18n:"label='Reuse of the SMTP: user', hint='Label of the option that takes the user of the SMTP.'"`
	SMTPReuseUser_Desc                  string           `i18n:"label='Reuse of the SMTP: user, description', hint='What the option that takes the user of the SMTP does.'"`
	SMTPReusePassword                   string           `i18n:"label='Reuse of the SMTP: password', hint='Label of the option that takes the password of the SMTP.'"`
	SMTPReusePassword_Desc              string           `i18n:"label='Reuse of the SMTP: password, description', hint='What the option that takes the password of the SMTP does.'"`
	SMTPReuseHost                       string           `i18n:"label='Reuse of the SMTP: host', hint='Label of the option that takes the host of the SMTP.'"`
	SMTPReuseHost_Desc                  string           `i18n:"label='Reuse of the SMTP: host, description', hint='What the option that takes the host of the SMTP does.'"`
}

var (
	Messages_en_US = &Messages{
		ModuleDescription:                 "The mail senders: the accounts the e-mails are sent through.",
		MailSender:                        "Mail Sender",
		TestMessageSubject:                "Test Send Mail",
		TestMessageBody:                   "Test OK.",
		MailSenderSubjectPrefix:           "Site name in the subject",
		MailSenderIDFooter:                "IDs in the body",
		MailSenderIDFooter_Hint:           "At the end of every message of a record, its ids — the record (hermon:issue:…), the Message-ID, the kind —: what finds it in a mailbox many things write to, and what an answer that lost its headers still quotes.",
		MailSenderSubjectRules:            "Subject parts",
		MailSenderSubjectRules_Hint:       "Which parts each subject has, by type of message and to whom it goes: the site's name, the thing it is of (an issue's type, a form), the record (#12). An answer always begins with Re:. A type with no row of its own follows Others.",
		SubjectRulesFamily:                "Type of message",
		SubjectRulesTeam:                  "To the team",
		SubjectRulesExternal:              "To outside",
		SubjectPartSite:                   "Site",
		SubjectPartTopic:                  "Thing",
		SubjectPartRef:                    "Record",
		SubjectFamilyOther:                "Others",
		MailSenderSender:                  "Sender",
		MailSenderGmail:                   "GMAIL",
		MailSenderSMTP:                    "SMTP",
		MailSender_Action_TestSendMail:    "Send Test Mail",
		GmailSenderConfiguredSuccessfully: "Gmail Sender configured Successfully",
		GmailSenderLogOutSuccessfully:     "Gmail Log Out Successfully",
		GmailSenderConnectPage:            "Connect Gmail",
		GmailSenderLogoutPage:             "Disconnect Gmail",
		GmailSenderCallbackPage:           "Gmail authorization return",
		GmailSenderSuccessPage:            "Gmail authorized",
		GmailSenderCallbackURI:            "CallbackURI",
		GmailSenderCredentials:            "App Credentials",
		GmailSenderCredentialsFile:        "App Credentials File",
		GmailSenderCredentialsFile_Hint:   "The .json file, available for download by Desktop Application on https://console.cloud.google.com",
		GmailSenderSetup:                  "How to create the credentials",
		GmailSenderSetupTitle:             "Gmail requires an OAuth client created by you on the Google Cloud Console. There is no API that creates it: the steps below are done once, by hand.",
		GmailSenderSetupSteps: []string{
			"Open the Google Cloud Console and create a project, or select an existing one.",
			"Enable the Gmail API: APIs & Services > Library > search for \"Gmail API\" > Enable.",
			"Configure the OAuth consent screen. While the app stays in \"Testing\", add the sending account under \"Test users\", otherwise Google blocks the sign in.",
			"Declare the scopes on the consent screen, following \"How to declare the scopes\" below.",
			"Create the client: APIs & Services > Credentials > Create credentials > OAuth client ID > Application type: Desktop app.",
			"Download the client JSON and upload it in the field below, then save.",
			"Back on this page, click \"Sign In\" and authorise with the account that will send the mail.",
		},
		GmailSenderSetupScopes:     "Scopes to declare",
		GmailSenderSetupScopesHint: "These are exactly what this sender requests. gmail.send is enough because sending is the only Gmail call it makes.",
		GmailSenderSetupScopesSteps: []string{
			"In the Console, open APIs & Services > OAuth consent screen (on the new Google Auth Platform it is called \"Data access\").",
			"Click \"Add or remove scopes\". A side panel opens listing the scopes of the APIs you enabled.",
			"Tick .../auth/userinfo.email and .../auth/userinfo.profile, which are the \"email\" and \"profile\" below.",
			"gmail.send only shows up in the list once the Gmail API is enabled. If it is missing, paste it into \"Manually add scopes\" at the bottom of the panel and click \"Add to table\".",
			"Click \"Update\" to close the panel, then \"Save\" at the bottom of the page. gmail.send will be filed under \"Sensitive scopes\": that is expected and needs no security assessment.",
			"Add nothing else. Each extra Gmail scope (mail.google.com, gmail.modify, gmail.compose) moves the app to the \"restricted\" tier, whose verification demands a security assessment.",
			"If this sender already has a token, click Sign Out and sign in again: an existing token keeps the scopes it was issued with.",
		},
		GmailSenderSetupBlocked:             "Seeing \"Access blocked: ... has not completed the Google verification process\"? The consent screen is published. Switch it back to Testing and add the sending account under Test users. On the \"Google hasn't verified this app\" warning, click Advanced and continue. Note that in Testing the refresh token expires after 7 days.",
		GmailSenderSetupCallback:            "Callback URL of this installation",
		GmailSenderSetupCallbackHint:        "A Desktop app client accepts any loopback address, so this URL does not need to be registered on the Console.",
		GmailSenderSetupConsoleBtn:          "Open the Google Cloud Console",
		SmtpSenderTLS:                       "TLS",
		SmtpSenderServer:                    "Server",
		SmtpSenderPort:                      "Port",
		SmtpSenderFromMail:                  "FromMail",
		SmtpSenderUser:                      "User",
		SmtpSenderPassword:                  "Password",
		MailSenderMaxPerRecipientHour:       "Emails per address, per hour",
		MailSenderMaxPerRecipientDay:        "Emails per address, per day",
		MailSenderInbound:                   "Answers to the emails",
		InboundSettingsMode:                 "Who reads the answers",
		InboundSettingsServer:               "IMAP server",
		InboundSettingsPort:                 "Port",
		InboundSettingsTLS:                  "TLS",
		InboundSettingsStartTLS:             "STARTTLS",
		InboundSettingsReuse:                "The same as the SMTP",
		InboundSettingsUser:                 "User",
		InboundSettingsPassword:             "Password",
		InboundSettingsAfter:                "After importing",
		InboundSettingsAfter_Hint:           "What is done with an answer imported (a comment now): marked read, deleted from the server or moved to a folder. Deleting needs UIDPLUS and moving MOVE (or UIDPLUS): a server without them only marks it read — anything else in the mailbox is never touched.",
		InboundSettingsMoveTo:               "Folder",
		InboundSettingsMoveTo_Hint:          "Where the answers imported go; empty, the domain of the sending's address (pac.com). Made when it is not there.",
		InboundAfterRead:                    "Mark read",
		InboundAfterDelete:                  "Delete from the server",
		InboundAfterMove:                    "Move to a folder",
		InboundSettingsMailbox:              "Mailbox",
		InboundManual:                       "People, in the team mailbox",
		InboundIMAP:                         "The application, by IMAP",
		MailSenderTabSending:                "Sending",
		MailSenderTabAnswers:                "Answers",
		MailSenderSender_Hint:               "Which account sends the emails: Gmail (authorized with Google) or SMTP (a server and its account). The fields of the account follow the choice.",
		MailSenderSubjectPrefix_Hint:        "The site's name, in brackets in the subjects (\"[PAC]\") as the subject parts say; empty, the site's name of its settings.",
		MailSenderMaxPerRecipientHour_Hint:  "The most emails one address outside the team (a contact, a customer) gets in an hour; past it, the email is refused. The team's own emails do not count. 0: no limit.",
		MailSenderMaxPerRecipientDay_Hint:   "The same, in a day. Together they keep a form of the site from being used to flood somebody's mailbox.",
		InboundSettingsMode_Hint:            "People: the answers stay in the team mailbox, read there. The application: it reads them by IMAP, and an answer to a conversation it knows becomes part of it.",
		InboundSettingsServer_Hint:          "The IMAP server of the team mailbox, such as imap.gmail.com.",
		InboundSettingsPort_Hint:            "0: the usual one, 993 with TLS, 143 without.",
		InboundSettingsTLS_Hint:             "An encrypted connection from the start (port 993): what most servers ask for.",
		InboundSettingsStartTLS_Hint:        "Encryption turned on after connecting (port 143), for servers that use it instead of TLS.",
		InboundSettingsReuse_Hint:           "What the reading of the answers takes from the SMTP account (Sending tab): its user, its password, its host — the mailbox that sends is the one that receives. A chip marked hides that field below; it is not used.",
		InboundSettingsUser_Hint:            "The account of the team mailbox, usually its address.",
		InboundSettingsPassword_Hint:        "Its password (with Gmail, an app password). Never shown again after saved.",
		InboundSettingsMailbox_Hint:         "The folder read; empty, the inbox (INBOX).",
		SendMailTestFormTo:                  "To",
		SendMailTestFormSubject:             "Subject",
		SendMailTestFormMessage:             "Message",
		SendMailTestFormSender:              "Sender",
		ErrGmailSenderScopeNotGranted:       "Google issued the token without the %s scope, so this sender cannot send mail. Declare that scope on the consent screen under \"Data access\", then Sign Out and sign in again, ticking the permission that lets the app send email on your behalf.",
		ErrGmailSenderCredentialsInvalid:    "The credentials is not a valid DESKTOP APPLICATION credentials",
		SendMailSuccessfully:                "Send mail successfully",
		GmailSenderToken:                    "Token",
		GmailSenderUser:                     "User",
		MailLog:                             "Mail sent",
		MailLogs:                            "Mails Sent",
		MailLogCreatedAt:                    "Sent at",
		MailLogStatus:                       "Status",
		MailLogStatusSent:                   "Sent",
		MailLogStatusFailed:                 "Failed",
		MailLogKind:                         "Kind",
		MailLogOrigin:                       "Origin",
		MailLogOpenOrigin:                   "Open origin",
		MailLogRecipients:                   "To",
		MailLogReplyTo:                      "Reply-To",
		MailLogSubject:                      "Subject",
		MailLogMethod:                       "Method",
		MailLogError:                        "Error",
		MailLogMessageID:                    "Message-ID",
		MailLogInReplyTo:                    "In-Reply-To",
		MailLogUserID:                       "User",
		MailLogRemoteIP:                     "IP",
		MailLogUserAgent:                    "Browser",
		MailLogLocaleCode:                   "Language",
		MailLogBody:                         "Body",
		MailLogBodyHidden:                   "[content hidden]",
		MailLogConversation:                 "Conversation",
		MailLogTabAll:                       "All",
		MailLogTabFailed:                    "Failed",
		MailLog_Action_Resend:               "Send again",
		MailSenderMailLogRetentionDays:      "Keep the mails sent for (days)",
		MailSenderMailLogRetentionDays_Hint: "How many days the central of the mails sent keeps each message; older ones are deleted every day. 0: for ever.",
		SMTPReuseUser:                       "User",
		SMTPReuseUser_Desc:                  "Reads the answers with the user of the SMTP account.",
		SMTPReusePassword:                   "Password",
		SMTPReusePassword_Desc:              "Reads the answers with the password of the SMTP account (with Gmail, its app password).",
		SMTPReuseHost:                       "Host",
		SMTPReuseHost_Desc:                  "Reads the answers from the host of the SMTP account — when the same host serves IMAP.",
	}

	Messages_pt_BR = &Messages{
		ModuleDescription:                 "Os remetentes de e-mail: as contas pelas quais os e-mails são enviados.",
		MailSender:                        "Envio de Email",
		TestMessageSubject:                "Teste de envio de email",
		TestMessageBody:                   "Teste executado com sucesso.",
		MailSenderSubjectPrefix:           "Nome do site no assunto",
		MailSenderIDFooter:                "IDs no corpo",
		MailSenderIDFooter_Hint:           "Ao fim de toda mensagem de um registro, os seus IDs — o registro (hermon:issue:…), o Message-ID, o tipo —: o que a acha numa caixa de entrada usada por várias coisas, e o que uma resposta que perdeu os cabeçalhos ainda cita.",
		MailSenderSubjectRules:            "Partes do assunto",
		MailSenderSubjectRules_Hint:       "Que partes cada assunto tem, por tipo de mensagem e a quem vai: o nome do site, a coisa (o tipo do chamado, o formulário), o registro (#12). Uma resposta começa sempre com Re:. Um tipo sem linha própria segue Outros.",
		SubjectRulesFamily:                "Tipo de mensagem",
		SubjectRulesTeam:                  "Para a equipe",
		SubjectRulesExternal:              "Para fora",
		SubjectPartSite:                   "Site",
		SubjectPartTopic:                  "Coisa",
		SubjectPartRef:                    "Registro",
		SubjectFamilyOther:                "Outros",
		MailSenderSender:                  "Método de Envio",
		MailSenderGmail:                   "GMAIL",
		MailSenderSMTP:                    "SMTP",
		MailSender_Action_TestSendMail:    "Enviar email de Teste",
		GmailSenderConfiguredSuccessfully: "Envio por Gmail configurado com sucesso",
		GmailSenderLogOutSuccessfully:     "Desconectado do GMAIL com sucesso",
		GmailSenderConnectPage:            "Conectar ao Gmail",
		GmailSenderLogoutPage:             "Desconectar do Gmail",
		GmailSenderCallbackPage:           "Retorno da autorização do Gmail",
		GmailSenderSuccessPage:            "Gmail autorizado",
		GmailSenderCallbackURI:            "CallbackURI",
		GmailSenderCredentials:            "Credenciais de App",
		GmailSenderCredentialsFile:        "Arquivo de Credenciais de App",
		GmailSenderCredentialsFile_Hint:   "Arquivo .json com as credenciais de acesso. Este arquivo está disponível para Download da Aplicação Desktop em https://console.cloud.google.com",
		GmailSenderSetup:                  "Como criar as credenciais",
		GmailSenderSetupTitle:             "O Gmail exige um cliente OAuth criado por você no Google Cloud Console. Não existe API que o crie: os passos abaixo são feitos uma única vez, à mão.",
		GmailSenderSetupSteps: []string{
			"Abra o Google Cloud Console e crie um projeto, ou selecione um existente.",
			"Ative a Gmail API: APIs e serviços > Biblioteca > procure por \"Gmail API\" > Ativar.",
			"Configure a tela de consentimento OAuth. Enquanto o app estiver em \"Testing\", inclua a conta remetente em \"Usuários de teste\", senão o Google bloqueia o login.",
			"Declare os escopos na tela de consentimento, seguindo \"Como declarar os escopos\" abaixo.",
			"Crie o cliente: APIs e serviços > Credenciais > Criar credenciais > ID do cliente OAuth > Tipo de aplicativo: App para computador (Desktop app).",
			"Baixe o JSON do cliente e envie-o no campo abaixo, depois salve.",
			"De volta a esta página, clique em \"Sign In\" e autorize com a conta que vai enviar os emails.",
		},
		GmailSenderSetupScopes:     "Escopos a declarar",
		GmailSenderSetupScopesHint: "É exatamente o que este remetente solicita. gmail.send basta porque enviar é a única chamada que ele faz ao Gmail.",
		GmailSenderSetupScopesSteps: []string{
			"No Console, abra APIs e serviços > Tela de permissão OAuth (na nova Google Auth Platform ela se chama \"Acesso a dados\").",
			"Clique em \"Adicionar ou remover escopos\". Abre um painel lateral com os escopos das APIs que você ativou.",
			"Marque .../auth/userinfo.email e .../auth/userinfo.profile, que são o \"email\" e o \"profile\" listados abaixo.",
			"O gmail.send só aparece na lista depois que a Gmail API está ativada. Se não estiver lá, cole-o em \"Adicionar escopos manualmente\", no rodapé do painel, e clique em \"Adicionar à tabela\".",
			"Clique em \"Atualizar\" para fechar o painel e depois em \"Salvar\", no rodapé da página. O gmail.send fica em \"Escopos sensíveis\": é o esperado e não exige avaliação de segurança.",
			"Não inclua mais nada. Cada escopo extra do Gmail (mail.google.com, gmail.modify, gmail.compose) joga o app na faixa \"restricted\", cuja verificação exige avaliação de segurança.",
			"Se este remetente já tiver token, clique em Sign Out e entre de novo: um token existente mantém os escopos com que foi emitido.",
		},
		GmailSenderSetupBlocked:             "Apareceu \"Acesso bloqueado: ... não concluiu o processo de verificação do Google\"? A tela de consentimento está publicada. Volte-a para Testing e inclua a conta remetente em Usuários de teste. No aviso \"O Google não verificou este app\", clique em Avançado e prossiga. Atenção: em Testing o refresh token expira em 7 dias.",
		GmailSenderSetupCallback:            "URL de callback desta instalação",
		GmailSenderSetupCallbackHint:        "Um cliente do tipo App para computador aceita qualquer endereço de loopback, então esta URL não precisa ser registrada no Console.",
		GmailSenderSetupConsoleBtn:          "Abrir o Google Cloud Console",
		SmtpSenderTLS:                       "TLS",
		SmtpSenderServer:                    "Servidor",
		SmtpSenderPort:                      "Porta",
		SmtpSenderFromMail:                  "Email DE",
		SmtpSenderUser:                      "Usuário",
		SmtpSenderPassword:                  "Senha",
		MailSenderMaxPerRecipientHour:       "E-mails por endereço, por hora",
		MailSenderMaxPerRecipientDay:        "E-mails por endereço, por dia",
		MailSenderInbound:                   "Respostas aos e-mails",
		InboundSettingsMode:                 "Quem lê as respostas",
		InboundSettingsServer:               "Servidor IMAP",
		InboundSettingsPort:                 "Porta",
		InboundSettingsTLS:                  "TLS",
		InboundSettingsStartTLS:             "STARTTLS",
		InboundSettingsReuse:                "O mesmo do SMTP",
		InboundSettingsUser:                 "Usuário",
		InboundSettingsPassword:             "Senha",
		InboundSettingsAfter:                "Depois de importar",
		InboundSettingsAfter_Hint:           "O que é feito com uma resposta importada (agora um comentário): marcada como lida, excluída do servidor ou movida para uma pasta. Excluir pede UIDPLUS e mover pede MOVE (ou UIDPLUS): um servidor sem eles só a marca como lida — nada mais da caixa é tocado.",
		InboundSettingsMoveTo:               "Pasta",
		InboundSettingsMoveTo_Hint:          "Para onde vão as respostas importadas; vazio, o domínio do endereço de envio (pac.com). Criada se não existir.",
		InboundAfterRead:                    "Marcar como lida",
		InboundAfterDelete:                  "Excluir do servidor",
		InboundAfterMove:                    "Mover para uma pasta",
		InboundSettingsMailbox:              "Caixa de entrada",
		InboundManual:                       "Pessoas, na caixa da equipe",
		InboundIMAP:                         "A aplicação, por IMAP",
		MailSenderTabSending:                "Envio",
		MailSenderTabAnswers:                "Respostas",
		MailSenderSender_Hint:               "Qual conta envia os e-mails: Gmail (autorizada com o Google) ou SMTP (um servidor e a sua conta). Os campos da conta acompanham a escolha.",
		MailSenderSubjectPrefix_Hint:        "O nome do site, entre colchetes nos assuntos (\"[PAC]\") como dizem as partes do assunto; vazio, o nome do site das suas configurações.",
		MailSenderMaxPerRecipientHour_Hint:  "O máximo de e-mails que um endereço de fora da equipe (um contato, um cliente) recebe em uma hora; acima, o e-mail é recusado. Os e-mails da própria equipe não contam. 0: sem limite.",
		MailSenderMaxPerRecipientDay_Hint:   "O mesmo, em um dia. Juntos, impedem que um formulário do site seja usado para encher a caixa de alguém.",
		InboundSettingsMode_Hint:            "Pessoas: as respostas ficam na caixa da equipe, lidas lá. A aplicação: ela as lê por IMAP, e a resposta a uma conversa que ela conhece passa a fazer parte dela.",
		InboundSettingsServer_Hint:          "O servidor IMAP da caixa da equipe, como imap.gmail.com.",
		InboundSettingsPort_Hint:            "0: a de costume, 993 com TLS, 143 sem.",
		InboundSettingsTLS_Hint:             "Conexão criptografada desde o início (porta 993): o que a maioria dos servidores pede.",
		InboundSettingsStartTLS_Hint:        "Criptografia ligada depois de conectar (porta 143), nos servidores que a usam no lugar do TLS.",
		InboundSettingsReuse_Hint:           "O que a leitura das respostas aproveita da conta SMTP (aba Envio): o usuário, a senha, o servidor — a caixa que envia é a que recebe. Um chip marcado esconde esse campo abaixo; ele não é usado.",
		InboundSettingsUser_Hint:            "A conta da caixa da equipe, em geral o seu endereço.",
		InboundSettingsPassword_Hint:        "A senha dela (no Gmail, uma senha de app). Nunca é mostrada de novo depois de salva.",
		InboundSettingsMailbox_Hint:         "A pasta lida; vazio, a caixa de entrada (INBOX).",
		SendMailTestFormTo:                  "Para",
		SendMailTestFormSubject:             "Assunto",
		SendMailTestFormMessage:             "Mensagem",
		SendMailTestFormSender:              "Método de Envio",
		ErrGmailSenderScopeNotGranted:       "O Google emitiu o token sem o escopo %s, então este remetente não consegue enviar. Declare esse escopo na tela de consentimento, em \"Acesso a dados\", depois clique em Sign Out e entre de novo, marcando a permissão que autoriza o app a enviar email em seu nome.",
		ErrGmailSenderCredentialsInvalid:    "Estas crendencias não são do tipo DESKTOP APPLICATION (Aplicação de Desktop).",
		SendMailSuccessfully:                "Email enviado com sucesso",
		GmailSenderToken:                    "Token",
		GmailSenderUser:                     "Usuário",
		MailLog:                             "E-mail enviado",
		MailLogs:                            "Central de E-mails",
		MailLogCreatedAt:                    "Enviado em",
		MailLogStatus:                       "Situação",
		MailLogStatusSent:                   "Enviado",
		MailLogStatusFailed:                 "Falhou",
		MailLogKind:                         "Tipo",
		MailLogOrigin:                       "Origem",
		MailLogOpenOrigin:                   "Abrir origem",
		MailLogRecipients:                   "Para",
		MailLogReplyTo:                      "Responder para",
		MailLogSubject:                      "Assunto",
		MailLogMethod:                       "Método",
		MailLogError:                        "Erro",
		MailLogMessageID:                    "Message-ID",
		MailLogInReplyTo:                    "Em resposta a",
		MailLogUserID:                       "Usuário",
		MailLogRemoteIP:                     "IP",
		MailLogUserAgent:                    "Navegador",
		MailLogLocaleCode:                   "Idioma",
		MailLogBody:                         "Corpo",
		MailLogBodyHidden:                   "[conteúdo oculto]",
		MailLogConversation:                 "Conversa",
		MailLogTabAll:                       "Todos",
		MailLogTabFailed:                    "Falhas",
		MailLog_Action_Resend:               "Enviar de novo",
		MailSenderMailLogRetentionDays:      "Guardar os e-mails enviados por (dias)",
		MailSenderMailLogRetentionDays_Hint: "Quantos dias a central de e-mails guarda cada mensagem; as mais antigas são apagadas todo dia. 0: para sempre.",
		SMTPReuseUser:                       "Usuário",
		SMTPReuseUser_Desc:                  "Lê as respostas com o usuário da conta SMTP.",
		SMTPReusePassword:                   "Senha",
		SMTPReusePassword_Desc:              "Lê as respostas com a senha da conta SMTP (no Gmail, a senha de app).",
		SMTPReuseHost:                       "Servidor",
		SMTPReuseHost_Desc:                  "Lê as respostas do servidor da conta SMTP — quando o mesmo servidor atende o IMAP.",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}
