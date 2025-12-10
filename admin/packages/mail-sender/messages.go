package mail_sender

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "rqv-admin/mail-sender"

type Messages struct {
	MailSender                        string
	TestMessageSubject                string
	TestMessageBody                   string
	MailSenderSender                  string
	MailSenderSubjectPrefix           string
	MailSenderGmail                   string
	MailSenderSMTP                    string
	MailSender_Action_TestSendMail    string
	GmailSenderConfiguredSuccessfully string
	GmailSenderLogOutSuccessfully     string
	GmailSenderCallbackURI            string
	GmailSenderCredentials            string
	GmailSenderCredentialsFile        string
	GmailSenderCredentialsFile_Hint   string
	SmtpSenderTLS                     string
	SmtpSenderServer                  string
	SmtpSenderPort                    string
	SmtpSenderFromMail                string
	SmtpSenderUser                    string
	SmtpSenderPassword                string
	SendMailTestFormTo                string
	SendMailTestFormSubject           string
	SendMailTestFormMessage           string
	SendMailTestFormSender            string
	ErrGmailSenderCredentialsInvalid  i18n.ErrorString
	SendMailSuccessfully              string
}

var (
	Messages_en_US = &Messages{
		MailSender:                        "Mail Sender",
		TestMessageSubject:                "Test Send Mail",
		TestMessageBody:                   "Test OK.",
		MailSenderSubjectPrefix:           "Subject Prefix",
		MailSenderSender:                  "Sender",
		MailSenderGmail:                   "GMAIL",
		MailSenderSMTP:                    "SMTP",
		MailSender_Action_TestSendMail:    "Send Test Mail",
		GmailSenderConfiguredSuccessfully: "Gmail Sender configured Successfully",
		GmailSenderLogOutSuccessfully:     "Gmail Log Out Successfully",
		GmailSenderCallbackURI:            "CallbackURI",
		GmailSenderCredentials:            "App Credentials",
		GmailSenderCredentialsFile:        "App Credentials File",
		GmailSenderCredentialsFile_Hint:   "The .json file, available for download by Desktop Application on https://console.cloud.google.com",
		SmtpSenderTLS:                     "TLS",
		SmtpSenderServer:                  "Server",
		SmtpSenderPort:                    "Port",
		SmtpSenderFromMail:                "FromMail",
		SmtpSenderUser:                    "User",
		SmtpSenderPassword:                "Password",
		SendMailTestFormTo:                "To",
		SendMailTestFormSubject:           "Subject",
		SendMailTestFormMessage:           "Message",
		SendMailTestFormSender:            "Sender",
		ErrGmailSenderCredentialsInvalid:  "The credentials is not a valid DESKTOP APPLICATION credentials",
		SendMailSuccessfully:              "Send mail successfully",
	}

	Messages_pt_BR = &Messages{
		MailSender:                        "Envio de Email",
		TestMessageSubject:                "Teste de envio de email",
		TestMessageBody:                   "Teste executado com sucesso.",
		MailSenderSubjectPrefix:           "Prefixo do Assunto",
		MailSenderSender:                  "Método de Envio",
		MailSenderGmail:                   "GMAIL",
		MailSenderSMTP:                    "SMTP",
		MailSender_Action_TestSendMail:    "Enviar email de Teste",
		GmailSenderConfiguredSuccessfully: "Envio por Gmail configurado com sucesso",
		GmailSenderLogOutSuccessfully:     "Desconectado do GMAIL com sucesso",
		GmailSenderCallbackURI:            "CallbackURI",
		GmailSenderCredentials:            "Credenciais de App",
		GmailSenderCredentialsFile:        "Arquivo de Credenciais de App",
		GmailSenderCredentialsFile_Hint:   "Arquivo .json com as credenciais de acesso. Este arquivo está disponível para Download da Aplicação Desktop em https://console.cloud.google.com",
		SmtpSenderTLS:                     "TLS",
		SmtpSenderServer:                  "Servidor",
		SmtpSenderPort:                    "Porta",
		SmtpSenderFromMail:                "Email DE",
		SmtpSenderUser:                    "Usuário",
		SmtpSenderPassword:                "Senha",
		SendMailTestFormTo:                "Para",
		SendMailTestFormSubject:           "Assunto",
		SendMailTestFormMessage:           "Mensagem",
		SendMailTestFormSender:            "Método de Envio",
		ErrGmailSenderCredentialsInvalid:  "Estas crendencias não são do tipo DESKTOP APPLICATION (Aplicação de Desktop).",
		SendMailSuccessfully:              "Email enviado com sucesso",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_pt_BR).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}
