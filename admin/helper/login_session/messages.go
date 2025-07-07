package login_session

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "admin/helper/login_session"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}
func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_pt_BR).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

type Messages struct {
	SignOutAllOtherSessions    string
	SignOutAllSuccessfullyTips string
	ChangePassword             string
	LoginSessions              string
	LoginSessionsTips          string
	Expired                    string
	Active                     string
	CurrentSession             string
	Time                       string
	Device                     string
	IPAddress                  string
	HideIPTips                 string
	Status                     string
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
