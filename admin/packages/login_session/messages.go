package login_session

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "admin/packages/login_session"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}
func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

type Messages struct {
	SignOutAllOtherSessions    string `i18n:"label='Sign out all other sessions', hint='Button that ends every session of the user but the current one.'"`
	SignOutAllSuccessfullyTips string `i18n:"label='Signed out of the others', hint='Shown once the user\\'s other sessions were ended.'"`
	ChangePassword             string `i18n:"hint='Button that opens the change of the user\\'s password.'"`
	LoginSessions              string `i18n:"hint='Title of the list of places where the user is logged in.'"`
	LoginSessionsTips          string `i18n:"label='Login sessions: hint', hint='Text under the title of the login sessions list.'"`
	Expired                    string `i18n:"hint='Status of a session that has ended.'"`
	Active                     string `i18n:"hint='Status of a session still open.'"`
	CurrentSession             string `i18n:"hint='Mark of the session the user is using now.'"`
	Time                       string `i18n:"hint='Column with when a session started.'"`
	Device                     string `i18n:"hint='Column with the browser and system of a session.'"`
	IPAddress                  string `i18n:"label='IP address', hint='Column with the address a session came from.'"`
	HideIPTips                 string `i18n:"label='Hidden IP: hint', hint='Shown in place of the address of a session when it is hidden.'"`
	Status                     string `i18n:"hint='Column with whether a session is active or expired.'"`
}

var (
	Messages_en_US = &Messages{
		SignOutAllOtherSessions:    "Signout all other sessions",
		SignOutAllSuccessfullyTips: "Sign out all successfully",
		ChangePassword:             "Change password",
		LoginSessions:              "Login Sessions",
		LoginSessionsTips:          "Places where you are connected to the administrator.",
		Expired:                    "Expired",
		Active:                     "Active",
		CurrentSession:             "Current Session",
		Time:                       "Time",
		Device:                     "Device",
		IPAddress:                  "IP Address",
		HideIPTips:                 "Hide IPTips",
		Status:                     "Status",
	}

	Messages_pt_BR = &Messages{
		ChangePassword:             "Alterar Senha",
		LoginSessions:              "Sessões de Login",
		LoginSessionsTips:          "Locais onde você está conectado ao administrador.",
		SignOutAllOtherSessions:    "Sair de todas as outras sessões",
		Expired:                    "Expirado",
		Active:                     "Ativo",
		CurrentSession:             "Sessão Atual",
		Time:                       "Horário",
		Device:                     "Dispositivo",
		IPAddress:                  "Endereço IP",
		HideIPTips:                 "Invisível devido a questões de segurança",
		SignOutAllSuccessfullyTips: "Todas as outras sessões foram desconectadas com sucesso.",
		Status:                     "Situação",
	}
)
