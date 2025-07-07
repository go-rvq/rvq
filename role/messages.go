package role

import (
	"context"

	"github.com/qor5/x/v3/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "admin/roles"

type Messages struct {
	Roles           string
	Role            string
	RolePermissions string
	RoleEffect      string
	RoleResources   string
	Allowed         string
	Denied          string
}

var (
	Messages_en_US = &Messages{
		Roles:           "Roles",
		Role:            "Role",
		RolePermissions: "Permissions",
		RoleEffect:      "Effect",
		RoleResources:   "Resources",
		Allowed:         "Allowed",
		Denied:          "Denied",
	}
	Messages_pt_BR = &Messages{
		Roles:           "Papéis de Usuário",
		Role:            "Papél de Usuário",
		RolePermissions: "Permissões",
		RoleEffect:      "Efeito",
		RoleResources:   "Recursos",
		Allowed:         "Permitir",
		Denied:          "Negar",
	}
)

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
