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
	PermissionsHelp          string `i18n:"hint='Title of the help button of the permissions of a role: it opens the documentation of policies and permissions.'"`
	RoleSystemKey            string `i18n:"label='System', hint='Label of the mark of a role of the system.'"`
	SystemRole               string `i18n:"hint='Said of a role of the system: made by the application, not deleted.'"`
	ResetPermissions         string `i18n:"hint='Action that resets the permissions of a role of the system to its originals.'"`
	PermissionsReset         string `i18n:"hint='Shown once the permissions of a role were reset.'"`
	ErrSystemRoleDelete      string `i18n:"hint='Error: a role of the system is not deleted.'"`
	ErrSystemRoleRename      string `i18n:"hint='Error: a role of the system is not renamed.'"`
	ErrSystemRoleFixed       string `i18n:"hint='Error: the permissions of a fixed role of the system are not changed.'"`
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
		PermissionsHelp:          "Help: policies and permissions",
		RoleSystemKey:            "System",
		SystemRole:               "System role",
		ResetPermissions:         "Reset the permissions",
		PermissionsReset:         "The permissions are the originals again.",
		ErrSystemRoleDelete:      "A role of the system is not deleted.",
		ErrSystemRoleRename:      "A role of the system keeps its name.",
		ErrSystemRoleFixed:       "This role may everything, by the application: its permissions are not changed.",
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
		PermissionsHelp:          "Ajuda: políticas e permissões",
		RoleSystemKey:            "Sistema",
		SystemRole:               "Papel do sistema",
		ResetPermissions:         "Restaurar as permissões",
		PermissionsReset:         "As permissões voltaram às originais.",
		ErrSystemRoleDelete:      "Um papel do sistema não é excluído.",
		ErrSystemRoleRename:      "Um papel do sistema mantém o nome.",
		ErrSystemRoleFixed:       "Este papel pode tudo, pela aplicação: as permissões dele não mudam.",
	}
)

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
