package profile

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "admin/packages/profile"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	ModuleDescription  string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	ChangePassword     string `i18n:"hint='Button of the user\\'s profile that opens the change of their password.'"`
	ProfileAccountName string `i18n:"hint='Label of the account name field of the profile.'"`
	ProfileAvatar      string `i18n:"hint='Label of the picture of the user in their profile (shown in the top bar).'"`
	MyProfile          string `i18n:"hint='Item of the menu of the user: their profile.'"`
	Login              string `i18n:"hint='Button of the top bar when no one is logged in.'"`
	OtherAccounts      string `i18n:"hint='Title of the section of the menu of the user with the other accounts logged in the browser.'"`
	SessionExpired     string `i18n:"hint='Under an account of the menu of the user whose session expired: its login asks it again.'"`
	AddAccount         string `i18n:"hint='Item of the menu of the user that logs in another account, the others staying logged in.'"`
	AccountsLimit      string `i18n:"hint='Under the item that logs in another account, when the browser holds as many as it can; %d is how many.'"`
	Logout             string `i18n:"hint='Item of the menu of the user: logs out.'"`
	LogoutThisAccount  string `i18n:"hint='Item of the menu of the user, with several accounts: logs out only the current one.'"`
	LogoutAll          string `i18n:"hint='Item of the menu of the user, with several accounts: logs out every one.'"`
}

var (
	Messages_en_US = &Messages{
		ModuleDescription:  "The profile of the signed-in user.",
		ChangePassword:     "Change Password",
		ProfileAccountName: "Account name",
		ProfileAvatar:      "Picture",
		MyProfile:          "My profile",
		Login:              "Log in",
		OtherAccounts:      "Other accounts",
		SessionExpired:     "session expired · log in again",
		AddAccount:         "Log in with another account",
		AccountsLimit:      "Limit of %d accounts: log out of one to add another",
		Logout:             "Log out",
		LogoutThisAccount:  "Log out of this account",
		LogoutAll:          "Log out of all accounts",
	}

	Messages_pt_BR = &Messages{
		ModuleDescription:  "O perfil do usuário conectado.",
		ChangePassword:     "Alterar Senha",
		ProfileAccountName: "Nome da conta",
		ProfileAvatar:      "Foto",
		MyProfile:          "Meu perfil",
		Login:              "Entrar",
		OtherAccounts:      "Outras contas",
		SessionExpired:     "sessão expirada · entrar de novo",
		AddAccount:         "Entrar com outra conta",
		AccountsLimit:      "Limite de %d contas: saia de uma para entrar com outra",
		Logout:             "Sair",
		LogoutThisAccount:  "Sair desta conta",
		LogoutAll:          "Sair de todas as contas",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
