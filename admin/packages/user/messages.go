package user

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "admin/packages/user"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	User                      string `i18n:"hint='Name of the users model in the singular (detail and form titles).'"`
	Users                     string `i18n:"hint='Name of the users model in the plural (menu, listing title).'"`
	UserCreatedAt             string `i18n:"label='Created', hint='Label of when a user was created.'"`
	UserName                  string `i18n:"label='Name', hint='Label of a user\\'s name.'"`
	UserStatus                string `i18n:"label='Status', hint='Label of whether a user is active.'"`
	UserRegistrationDate      string `i18n:"label='Registration date', hint='Label of the date a user registered.'"`
	UserRegistrationDateRange string `i18n:"label='Registration date range', hint='Filter of the users listing by registration date.'"`
	MailSentSuccessfully      string `i18n:"label='Email sent', hint='Shown once an email was sent to a user.'"`

	AllSessionLogsExpiredSuccessfully string `i18n:"label='Sessions expired', hint='Shown once all sessions of a user were ended.'"`
	UserUnlockedSuccessfully          string `i18n:"label='User unlocked', hint='Shown once a locked user was unlocked.'"`

	ErrorAccountRequired string `i18n:"label='Account required', hint='Error when a user is saved without an account or email.'"`
	Active               string `i18n:"hint='Status of a user who can log in.'"`
	Inactive             string `i18n:"hint='Status of a user who cannot log in.'"`
	Actives              string `i18n:"hint='Tab of the users listing that shows the active users.'"`
	Inactives            string `i18n:"hint='Tab of the users listing that shows the inactive users.'"`

	// AnonymousUserName is the display name of the static anonymous user,
	// resolved per request from its language.
	AnonymousUserName string `i18n:"label='Anonymous user name', hint='Name shown for the anonymous user (what is done with no one logged in).'"`

	SendResetPasswordEmail string `i18n:"label='Send password reset email', hint='Action that emails a user a link to reset their password.'"`
	Unlock                 string `i18n:"hint='Action that unlocks a user locked out by failed logins.'"`
	RevokeTOTP             string `i18n:"label='Revoke TOTP', hint='Action that removes a user\\'s two-factor (TOTP) authentication.'"`
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
		AnonymousUserName:                 "Anonymous",
		SendResetPasswordEmail:            "Send password reset email",
		Unlock:                            "Unlock",
		RevokeTOTP:                        "Revoke TOTP",
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
		AnonymousUserName:                 "Anônimo",
		SendResetPasswordEmail:            "Enviar email para alterar a senha",
		Unlock:                            "Desbloquear",
		RevokeTOTP:                        "Revogar TOTP",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
