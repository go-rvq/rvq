package role

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "admin/roles"

type Messages struct {
	ModuleDescription        string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Roles                    string `i18n:"hint='Name of the roles model in the plural (menu, listing title).'"`
	Role                     string `i18n:"hint='Name of the roles model in the singular (detail and form titles).'"`
	RolePermissions          string `i18n:"label='Permissions', hint='Label of a role\\'s permissions field.'"`
	DefaultDbpolicyEffect    string `i18n:"label='Effect', hint='Label of whether a permission of a role allows or denies.'"`
	DefaultDbpolicyResources string `i18n:"label='Resources', hint='Label of the resources a permission of a role applies to.'"`
	Allowed                  string `i18n:"hint='Effect of a permission that allows what it names.'"`
	Denied                   string `i18n:"hint='Effect of a permission that denies what it names.'"`
	DefaultDbpolicy          string `i18n:"hint='Name of a permission of a role.'"`
	DefaultDbpolicyActions   string `i18n:"hint='Label of the actions a permission of a role applies to.'"`
}

var (
	Messages_en_US = &Messages{
		ModuleDescription:        "The roles of the users.",
		Roles:                    "Roles",
		Role:                     "Role",
		RolePermissions:          "Permissions",
		DefaultDbpolicyEffect:    "Effect",
		DefaultDbpolicyResources: "Resources",
		Allowed:                  "Allowed",
		Denied:                   "Denied",
		DefaultDbpolicy:          "Permission",
		DefaultDbpolicyActions:   "Actions",
	}
	Messages_pt_BR = &Messages{
		ModuleDescription:        "Os papéis dos usuários.",
		Roles:                    "Papéis de Usuário",
		Role:                     "Papél de Usuário",
		RolePermissions:          "Permissões",
		DefaultDbpolicyEffect:    "Efeito",
		DefaultDbpolicyResources: "Recursos",
		Allowed:                  "Permitir",
		Denied:                   "Negar",
		DefaultDbpolicy:          "Permissão",
		DefaultDbpolicyActions:   "Ações",
	}
)

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
