package role

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "admin/roles"

type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Roles             string `i18n:"hint='Name of the roles model in the plural (menu, listing title).'"`
	Role              string `i18n:"hint='Name of the roles model in the singular (detail and form titles).'"`
	RolePermissions   string `i18n:"label='Permissions', hint='Label of a role\\'s permissions field.'"`
	RoleEffect        string `i18n:"label='Effect', hint='Label of whether a permission of a role allows or denies.'"`
	RoleResources     string `i18n:"label='Resources', hint='Label of the resources a permission of a role applies to.'"`
	Allowed           string `i18n:"hint='Effect of a permission that allows what it names.'"`
	Denied            string `i18n:"hint='Effect of a permission that denies what it names.'"`
}

var (
	Messages_en_US = &Messages{
		ModuleDescription: "The roles of the users.",
		Roles:             "Roles",
		Role:              "Role",
		RolePermissions:   "Permissions",
		RoleEffect:        "Effect",
		RoleResources:     "Resources",
		Allowed:           "Allowed",
		Denied:            "Denied",
	}
	Messages_pt_BR = &Messages{
		ModuleDescription: "Os papéis dos usuários.",
		Roles:             "Papéis de Usuário",
		Role:              "Papél de Usuário",
		RolePermissions:   "Permissões",
		RoleEffect:        "Efeito",
		RoleResources:     "Recursos",
		Allowed:           "Permitir",
		Denied:            "Negar",
	}
)

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
