package profile

import (
	"context"

	"github.com/qor5/x/v3/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "admin/helper/profile"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	ChangePassword string
}

var (
	Messages_en_US = &Messages{}

	Messages_pt_BR = &Messages{
		ChangePassword: "Alterar Senha",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_pt_BR).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
