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
	ManageTitle    string `i18n:"hint='Title of the dialog that manages who may do what with a record.'"`
	Subject        string `i18n:"hint='Label of the roles or users a permission is given to.'"`
	SubjectHint    string `i18n:"label='Subject: hint', hint='Hint of the roles or users field: one per line, or separated by commas.'"`
	Permissions    string `i18n:"hint='Title of the permissions given over a record.'"`
	Actions        string `i18n:"hint='Title of the actions a permission covers.'"`
	Pages          string `i18n:"hint='Title of the pages a permission covers.'"`
	NoGrants       string `i18n:"hint='Shown when no permission was given over a record.'"`
	Help           string `i18n:"hint='Explanation of the permissions dialog: how to give and how to revoke.'"`
	View           string `i18n:"hint='Permission to see a record.'"`
	Edit           string `i18n:"hint='Permission to change a record.'"`
	Delete         string `i18n:"hint='Permission to delete a record.'"`
	FieldsView     string `i18n:"label='Fields: view', hint='Permission to see only some fields of a record.'"`
	FieldsEdit     string `i18n:"label='Fields: edit', hint='Permission to change only some fields of a record.'"`
	FieldsHelp     string `i18n:"label='Fields: hint', hint='Hint of the fields a permission may be restricted to.'"`
	RevokeHint     string `i18n:"label='Revoke: hint', hint='Hint that saving with nothing checked revokes the permissions of the roles or users given.'"`
	TabAll         string `i18n:"label='Tab: all', hint='Tab of a listing that shows the records not deleted.'"`
	TabTrash       string `i18n:"label='Tab: trash', hint='Tab of a listing that shows the deleted records.'"`
	Restore        string `i18n:"hint='Action that restores deleted records.'"`
	TrashEmptyHint string `i18n:"label='Trash empty', hint='Shown when there is no deleted record.'"`
	// The trash's columns — when, by whom and from where a record was
	// deleted —, the action that shows where from, and its words.
	// DeletedNotice is the warning of the detail of a deleted record.
	DeletedNotice     string `i18n:"hint='Warning at the top of the detail of a deleted record.'"`
	DeletedAt         string `i18n:"label='Deleted on', hint='Label of when a record was deleted.'"`
	DeletedBy         string `i18n:"hint='Label of who deleted a record.'"`
	DeletedOrigin     string `i18n:"label='Deleted from', hint='Label of where a record was deleted from, with the action that shows it on a map.'"`
	OriginIP          string `i18n:"label='Origin: address', hint='Label of the address a record was deleted from.'"`
	OriginBrowser     string `i18n:"label='Origin: browser', hint='Label of the browser a record was deleted from.'"`
	OriginPlace       string `i18n:"label='Origin: place', hint='Label of the place (city, region, country) a record was deleted from.'"`
	OriginCoordinates string `i18n:"label='Origin: coordinates', hint='Label of the coordinates of the place a record was deleted from.'"`
	OriginNoMap       string `i18n:"label='Origin: no map', hint='Shown when the place of the address is not known, so there is no map.'"`
	Shared            string `i18n:"label='Share', hint='Mark of a permission given by sharing the record.'"`
	SharedNote        string `i18n:"label='Shared: note', hint='Hint that a permission given by sharing is managed in Sharing, not here.'"`
}

var Messages_en_US = &Messages{
	ManageTitle:       "Manage Permissions",
	Subject:           "Roles or users",
	SubjectHint:       "One per line (or comma-separated). Applies to all of them.",
	Permissions:       "Permissions",
	Actions:           "Actions",
	Pages:             "Pages",
	NoGrants:          "No permissions granted for this record.",
	Help:              "Grant a role/user permission over this record. Clear everything and save to revoke.",
	View:              "View",
	Edit:              "Edit",
	Delete:            "Delete",
	FieldsView:        "Fields — view",
	FieldsEdit:        "Fields — edit",
	FieldsHelp:        "Optional: restrict the grant to specific fields (nested included).",
	RevokeHint:        "Saving with no option checked revokes the given roles/users' permissions.",
	TabAll:            "All",
	TabTrash:          "Trash",
	Restore:           "Restore",
	TrashEmptyHint:    "No deleted records.",
	DeletedNotice:     "This record was deleted.",
	DeletedAt:         "Deleted on",
	DeletedBy:         "Deleted by",
	DeletedOrigin:     "Deleted from",
	OriginIP:          "Address",
	OriginBrowser:     "Browser",
	OriginPlace:       "Place",
	OriginCoordinates: "Coordinates",
	OriginNoMap:       "Where the address is is not known: no map.",
	Shared:            "Share",
	SharedNote:        "Managed in Sharing — cannot be removed here.",
}

var Messages_pt_BR = &Messages{
	ManageTitle:       "Gerenciar Permissões",
	Subject:           "Papéis ou usuários",
	SubjectHint:       "Um por linha (ou separados por vírgula). Aplica a todos os informados.",
	Permissions:       "Permissões",
	Actions:           "Ações",
	Pages:             "Páginas",
	NoGrants:          "Nenhuma permissão concedida para este registro.",
	Help:              "Conceda a um papel/usuário permissão sobre este registro. Desmarque tudo e salve para revogar.",
	View:              "Visualizar",
	Edit:              "Editar",
	Delete:            "Excluir",
	FieldsView:        "Campos — visualizar",
	FieldsEdit:        "Campos — editar",
	FieldsHelp:        "Opcional: restringe a concessão a campos específicos (nested incluídos).",
	RevokeHint:        "Salvar sem nenhuma opção marcada revoga as permissões dos papéis/usuários informados.",
	TabAll:            "Tudo",
	TabTrash:          "Lixeira",
	Restore:           "Restaurar",
	TrashEmptyHint:    "Nenhum registro excluído.",
	DeletedNotice:     "Este registro foi excluído.",
	DeletedAt:         "Excluído em",
	DeletedBy:         "Excluído por",
	DeletedOrigin:     "Excluído de",
	OriginIP:          "Endereço",
	OriginBrowser:     "Navegador",
	OriginPlace:       "Lugar",
	OriginCoordinates: "Coordenadas",
	OriginNoMap:       "Não se sabe onde fica o endereço: sem mapa.",
	Shared:            "Compartilhamento",
	SharedNote:        "Gerenciado em Compartilhamento — não pode ser removido aqui.",
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
