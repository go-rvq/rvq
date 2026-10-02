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
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	ChangePassword    string `i18n:"hint='Button of the user\\'s profile that opens the change of their password.'"`
}

var (
	Messages_en_US = &Messages{
		ModuleDescription: "The profile of the signed-in user.",
		ChangePassword:    "Change Password",
	}

	Messages_pt_BR = &Messages{
		ModuleDescription: "O perfil do usuário conectado.",
		ChangePassword:    "Alterar Senha",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
