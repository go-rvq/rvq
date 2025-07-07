package user

import (
	"context"

	"github.com/qor5/x/v3/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "admin/helper/user"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	User                      string
	Users                     string
	UserCreatedAt             string
	UserName                  string
	UserStatus                string
	UserRegistrationDate      string
	UserRegistrationDateRange string
	MailSentSuccessfully      string

	AllSessionLogsExpiredSuccessfully string
	UserUnlockedSuccessfully          string

	ErrorAccountRequired string
	Active               string
	Inactive             string
	Actives              string
	Inactives            string
}

var (
	Messages_en_US = &Messages{
		User:                              "User",
		Users:                             "Users",
		UserCreatedAt:                     "Created",
		UserName:                          "Name",
		UserStatus:                        "Status",
		UserRegistrationDate:              "Registration date",
		UserRegistrationDateRange:         "Registration date Range",
		AllSessionLogsExpiredSuccessfully: "All session logs expired successfully",
		UserUnlockedSuccessfully:          "User Unlocked Successfully",
		MailSentSuccessfully:              "Email sent successfully",
		ErrorAccountRequired:              "Account/Email required",
		Active:                            "Active",
		Inactive:                          "Inactive",
		Actives:                           "Actives",
		Inactives:                         "Inactives",
	}

	Messages_pt_BR = &Messages{
		User:                              "Usuário",
		Users:                             "Usuários",
		UserCreatedAt:                     "Cadastro",
		UserName:                          "Nome",
		UserStatus:                        "Situação",
		UserRegistrationDate:              "Registro",
		UserRegistrationDateRange:         "Período de Registro",
		AllSessionLogsExpiredSuccessfully: "Você foi desconectado de todas as sessões ativas",
		UserUnlockedSuccessfully:          "Usuário desbloqueado com sucesso",
		ErrorAccountRequired:              "Nome da conta/Email é obrigatório",
		MailSentSuccessfully:              "Email enviado com sucesso",
		Active:                            "Ativo",
		Inactive:                          "Inativo",
		Actives:                           "Ativos",
		Inactives:                         "Inativos",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
