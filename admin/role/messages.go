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
	// %s: the line
	ErrLine                string `i18n:"hint='Error: a line of the permissions of a role whose prefix (deny, allow, reset) has no resource after it.'"`
	RoleSystemKey          string `i18n:"label='System', hint='Label of the mark of a role of the system.'"`
	SystemRole             string `i18n:"hint='Said of a role of the system: made by the application, not deleted.'"`
	ResetPermissions       string `i18n:"hint='Action that resets the permissions of a role of the system to its originals.'"`
	PermissionsReset       string `i18n:"hint='Shown once the permissions of a role were reset.'"`
	ErrSystemRoleDelete    string `i18n:"hint='Error: a role of the system is not deleted.'"`
	ErrSystemRoleRename    string `i18n:"hint='Error: a role of the system is not renamed.'"`
	ErrSystemRoleFixed     string `i18n:"hint='Error: the permissions of a fixed role of the system are not changed.'"`
	PermSets               string `i18n:"hint='Name of the sets of permissions model in the plural (menu, listing title).'"`
	PermSet                string `i18n:"hint='Name of the sets of permissions model in the singular.'"`
	PermSetName            string `i18n:"hint='Label of the name of a set of permissions, as a line writes it (%name).'"`
	PermSetName_Hint       string `i18n:"hint='Hint of the name of a set.'"`
	PermSetDefinition      string `i18n:"hint='Label of what a set gives, by kind of resource.'"`
	PermSetDefinition_Hint string `i18n:"hint='Hint of the definition of a set.'"`
	PermSetTitle           string `i18n:"hint='Label of the title of a set, in the language of its localization.'"`
	PermSetDescription     string `i18n:"hint='Label of the description of a set, in the language of its localization.'"`
	PermSetSystem          string `i18n:"hint='Label of the mark of a set of the system: of the code, its name and definition read only, its title and description translatable.'"`
	PermSetCreatedAt       string `i18n:"hint='Label of when a set was made.'"`
	PermSetCreatedBy       string `i18n:"hint='Label of who made a set.'"`
	PermSetUpdatedAt       string `i18n:"hint='Label of when a set was changed last.'"`
	PermSetUpdatedBy       string `i18n:"hint='Label of who changed a set last.'"`
	PermSetLocale          string `i18n:"hint='Label of the language of a set: the language of its title and description.'"`
	ErrPermSetSystem       string `i18n:"hint='Error: a set of the system is not deleted.'"`
	ErrPermSetName         string `i18n:"hint='Error: a name a set may not have.'"`
	ErrPermSetTaken        string `i18n:"hint='Error: the name of a set taken, by another or by one of the system.'"`
	ErrPermSetKind         string `i18n:"hint='Error: a kind of resource a definition names that the admin has not; %s the kind.'"`
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
		ErrLine:                  "%q: a prefix (deny, allow, reset) is its word and a space, before a resource (\"deny admin:users:*\").",
		RoleSystemKey:            "System",
		SystemRole:               "System role",
		ResetPermissions:         "Reset the permissions",
		PermissionsReset:         "The permissions are the originals again.",
		ErrSystemRoleDelete:      "A role of the system is not deleted.",
		ErrSystemRoleRename:      "A role of the system keeps its name.",
		ErrSystemRoleFixed:       "This role may everything, by the application: its permissions are not changed.",
		PermSets:                 "Permission sets",
		PermSet:                  "Permission set",
		PermSetName:              "Name",
		PermSetName_Hint:         "as a permission writes it, after %: lower case letters, digits and _",
		PermSetDefinition:        "Definition",
		PermSetDefinition_Hint:   "by kind of resource, what it gives: {model: {@list, <*>:@get}, singleton: @get}",
		PermSetTitle:             "Title",
		PermSetDescription:       "Description",
		PermSetSystem:            "Of the system",
		PermSetCreatedAt:         "Created at",
		PermSetCreatedBy:         "Created by",
		PermSetUpdatedAt:         "Updated at",
		PermSetUpdatedBy:         "Updated by",
		PermSetLocale:            "Language",
		ErrPermSetSystem:         "A set of the system is not deleted: its title and description may be translated.",
		ErrPermSetName:           "A name of lower case letters, digits and _, a letter first.",
		ErrPermSetTaken:          "There is a set of this name already.",
		ErrPermSetKind:           "Unknown kind of resource: %s.",
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
		ErrLine:                  "%q: um prefixo (deny, allow, reset) é a palavra e um espaço, antes de um recurso (\"deny admin:users:*\").",
		RoleSystemKey:            "Sistema",
		SystemRole:               "Papel do sistema",
		ResetPermissions:         "Restaurar as permissões",
		PermissionsReset:         "As permissões voltaram às originais.",
		ErrSystemRoleDelete:      "Um papel do sistema não é excluído.",
		ErrSystemRoleRename:      "Um papel do sistema mantém o nome.",
		ErrSystemRoleFixed:       "Este papel pode tudo, pela aplicação: as permissões dele não mudam.",
		PermSets:                 "Conjuntos de permissões",
		PermSet:                  "Conjunto de permissões",
		PermSetName:              "Nome",
		PermSetName_Hint:         "como uma permissão o escreve, depois de %: letras minúsculas, números e _",
		PermSetDefinition:        "Definição",
		PermSetDefinition_Hint:   "por tipo de recurso, o que ele dá: {model: {@list, <*>:@get}, singleton: @get}",
		PermSetTitle:             "Título",
		PermSetDescription:       "Descrição",
		PermSetSystem:            "Do sistema",
		PermSetCreatedAt:         "Criado em",
		PermSetCreatedBy:         "Criado por",
		PermSetUpdatedAt:         "Atualizado em",
		PermSetUpdatedBy:         "Atualizado por",
		PermSetLocale:            "Idioma",
		ErrPermSetSystem:         "Um conjunto do sistema não é excluído: o título e a descrição dele podem ser traduzidos.",
		ErrPermSetName:           "Um nome de letras minúsculas, números e _, começando por uma letra.",
		ErrPermSetTaken:          "Já existe um conjunto com este nome.",
		ErrPermSetKind:           "Tipo de recurso desconhecido: %s.",
	}
)

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
