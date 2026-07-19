package perms

import (
	"context"

	"github.com/go-rvq/rvq/admin/presets"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// I18nKey is the i18n module key of the perms package.
const I18nKey i18n.ModuleKey = "PermsI18n"

// Messages holds the labels of the permission manager and the trash.
type Messages struct {
	ManageTitle    string
	Subject        string
	SubjectHint    string
	Permissions    string
	Actions        string
	NoGrants       string
	Help           string
	View           string
	Edit           string
	Delete         string
	FieldsView     string
	FieldsEdit     string
	FieldsHelp     string
	RevokeHint     string
	TabAll         string
	TabTrash       string
	Restore        string
	TrashEmptyHint string
}

var Messages_en_US = &Messages{
	ManageTitle:    "Manage Permissions",
	Subject:        "Roles or users",
	SubjectHint:    "One per line (or comma-separated). Applies to all of them.",
	Permissions:    "Permissions",
	Actions:        "Actions",
	NoGrants:       "No permissions granted for this record.",
	Help:           "Grant a role/user permission over this record. Clear everything and save to revoke.",
	View:           "View",
	Edit:           "Edit",
	Delete:         "Delete",
	FieldsView:     "Fields — view",
	FieldsEdit:     "Fields — edit",
	FieldsHelp:     "Optional: restrict the grant to specific fields (nested included).",
	RevokeHint:     "Saving with no option checked revokes the given roles/users' permissions.",
	TabAll:         "All",
	TabTrash:       "Trash",
	Restore:        "Restore",
	TrashEmptyHint: "No deleted records.",
}

var Messages_pt_BR = &Messages{
	ManageTitle:    "Gerenciar Permissões",
	Subject:        "Papéis ou usuários",
	SubjectHint:    "Um por linha (ou separados por vírgula). Aplica a todos os informados.",
	Permissions:    "Permissões",
	Actions:        "Ações",
	NoGrants:       "Nenhuma permissão concedida para este registro.",
	Help:           "Conceda a um papel/usuário permissão sobre este registro. Desmarque tudo e salve para revogar.",
	View:           "Visualizar",
	Edit:           "Editar",
	Delete:         "Excluir",
	FieldsView:     "Campos — visualizar",
	FieldsEdit:     "Campos — editar",
	FieldsHelp:     "Opcional: restringe a concessão a campos específicos (nested incluídos).",
	RevokeHint:     "Salvar sem nenhuma opção marcada revoga as permissões dos papéis/usuários informados.",
	TabAll:         "Tudo",
	TabTrash:       "Lixeira",
	Restore:        "Restaurar",
	TrashEmptyHint: "Nenhum registro excluído.",
}

func msgs(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nKey, Messages_en_US).(*Messages)
}

// registerMessages registers the perms messages on the builder (idempotent).
func registerMessages(b *presets.Builder) {
	b.I18n().
		RegisterForModule(language.AmericanEnglish, I18nKey, Messages_en_US).
		RegisterForModule(language.BrazilianPortuguese, I18nKey, Messages_pt_BR)
}
