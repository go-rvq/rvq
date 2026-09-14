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
	GmailSenderSetup                  string
	GmailSenderSetupTitle             string
	GmailSenderSetupSteps             []string
	GmailSenderSetupScopes            string
	GmailSenderSetupScopesHint        string
	GmailSenderSetupScopesSteps       []string
	GmailSenderSetupBlocked           string
	GmailSenderSetupCallback          string
	GmailSenderSetupCallbackHint      string
	GmailSenderSetupConsoleBtn        string
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
	ErrGmailSenderScopeNotGranted     string
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
		GmailSenderSetupBlocked:          "Seeing \"Access blocked: ... has not completed the Google verification process\"? The consent screen is published. Switch it back to Testing and add the sending account under Test users. On the \"Google hasn't verified this app\" warning, click Advanced and continue. Note that in Testing the refresh token expires after 7 days.",
		GmailSenderSetupCallback:         "Callback URL of this installation",
		GmailSenderSetupCallbackHint:     "A Desktop app client accepts any loopback address, so this URL does not need to be registered on the Console.",
		GmailSenderSetupConsoleBtn:       "Open the Google Cloud Console",
		SmtpSenderTLS:                    "TLS",
		SmtpSenderServer:                 "Server",
		SmtpSenderPort:                   "Port",
		SmtpSenderFromMail:               "FromMail",
		SmtpSenderUser:                   "User",
		SmtpSenderPassword:               "Password",
		SendMailTestFormTo:               "To",
		SendMailTestFormSubject:          "Subject",
		SendMailTestFormMessage:          "Message",
		SendMailTestFormSender:           "Sender",
		ErrGmailSenderScopeNotGranted:    "Google issued the token without the %s scope, so this sender cannot send mail. Declare that scope on the consent screen under \"Data access\", then Sign Out and sign in again, ticking the permission that lets the app send email on your behalf.",
		ErrGmailSenderCredentialsInvalid: "The credentials is not a valid DESKTOP APPLICATION credentials",
		SendMailSuccessfully:             "Send mail successfully",
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
		GmailSenderSetupBlocked:          "Apareceu \"Acesso bloqueado: ... não concluiu o processo de verificação do Google\"? A tela de consentimento está publicada. Volte-a para Testing e inclua a conta remetente em Usuários de teste. No aviso \"O Google não verificou este app\", clique em Avançado e prossiga. Atenção: em Testing o refresh token expira em 7 dias.",
		GmailSenderSetupCallback:         "URL de callback desta instalação",
		GmailSenderSetupCallbackHint:     "Um cliente do tipo App para computador aceita qualquer endereço de loopback, então esta URL não precisa ser registrada no Console.",
		GmailSenderSetupConsoleBtn:       "Abrir o Google Cloud Console",
		SmtpSenderTLS:                    "TLS",
		SmtpSenderServer:                 "Servidor",
		SmtpSenderPort:                   "Porta",
		SmtpSenderFromMail:               "Email DE",
		SmtpSenderUser:                   "Usuário",
		SmtpSenderPassword:               "Senha",
		SendMailTestFormTo:               "Para",
		SendMailTestFormSubject:          "Assunto",
		SendMailTestFormMessage:          "Mensagem",
		SendMailTestFormSender:           "Método de Envio",
		ErrGmailSenderScopeNotGranted:    "O Google emitiu o token sem o escopo %s, então este remetente não consegue enviar. Declare esse escopo na tela de consentimento, em \"Acesso a dados\", depois clique em Sign Out e entre de novo, marcando a permissão que autoriza o app a enviar email em seu nome.",
		ErrGmailSenderCredentialsInvalid: "Estas crendencias não são do tipo DESKTOP APPLICATION (Aplicação de Desktop).",
		SendMailSuccessfully:             "Email enviado com sucesso",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}
