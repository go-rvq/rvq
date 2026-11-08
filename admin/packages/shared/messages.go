package shared

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// PermShare is the perm verb (and detail-action name) guarding the share action.
const PermShare = "!share"

// I18nSharedKey is the i18n module key of the shared package.
const I18nSharedKey i18n.ModuleKey = "SharedManager"

// Messages are the sharing UI labels.
type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	ShareTitle        string `i18n:"label='Share: title', hint='Title of the dialog that shares a record.'"`
	ShareWith         string `i18n:"hint='Label of the users a record is shared with.'"`
	ShareWithHint     string `i18n:"label='Share with: hint', hint='Hint of the users field: one per line, or separated by commas or semicolons.'"`
	Permissions       string `i18n:"hint='Label of what the users a record is shared with may do.'"`
	CanView           string `i18n:"hint='Permission to see a shared record.'"`
	CanEdit           string `i18n:"hint='Permission to change a shared record.'"`
	CanDelete         string `i18n:"hint='Permission to delete a shared record.'"`
	SharedWith        string `i18n:"hint='Title of the users a record is shared with.'"`
	NoShares          string `i18n:"hint='Shown when a record is shared with no one.'"`
	Add               string `i18n:"label='Share', hint='Button that shares the record with the users chosen.'"`
	Remove            string `i18n:"hint='Button that stops sharing the record with a user.'"`
	MyInvites         string `i18n:"label='My invites', hint='Title of the shares the user was invited to.'"`
	NoInvites         string `i18n:"hint='Shown when the user has no invite waiting.'"`
	AcceptInvite      string `i18n:"label='Accept invite', hint='Button that accepts an invite to a shared record.'"`
	RejectInvite      string `i18n:"label='Decline invite', hint='Button that declines an invite to a shared record.'"`
	AccessesTitle     string `i18n:"hint='Title of the dialog of who reaches a record.'"`
	HistoryTitle      string `i18n:"hint='Title of the dialog of the history of the sharing of a record.'"`
	Subject           string `i18n:"hint='Label of whom a record is shared with.'"`
	SubjectUser       string `i18n:"hint='Kind of a subject: a user.'"`
	SubjectGroup      string `i18n:"hint='Kind of a subject: a group.'"`
	Role              string `i18n:"hint='Label of the role of an access.'"`
	Owner             string `i18n:"hint='The owner of a record.'"`
	Inherited         string `i18n:"hint='Says an access comes from a record above.'"`
	Revoke            string `i18n:"hint='Button that takes an access away.'"`
	ChangeRole        string `i18n:"hint='Button that changes the role of an access.'"`
	RequestSent       string `i18n:"hint='Shown when a sharing waits for the owner.'"`
	RoleChanged       string `i18n:"hint='Shown when a role changed at once.'"`
	Revoked           string `i18n:"hint='Shown when an access was taken away.'"`
	ConfirmTitle      string `i18n:"hint='Title of the page where the owner confirms a request.'"`
	Confirm           string `i18n:"hint='Button that confirms a request.'"`
	Decline           string `i18n:"hint='Button that declines a request.'"`
	Confirmed         string `i18n:"hint='Shown when a request was confirmed.'"`
	Declined          string `i18n:"hint='Shown when a request was declined.'"`
	RequestedBy       string `i18n:"hint='Label of who asked.'"`
	Record            string `i18n:"hint='Label of the record of a request.'"`
	Event             string `i18n:"hint='Label of an event of the history.'"`
	When              string `i18n:"hint='Label of when an event happened.'"`
	Who               string `i18n:"hint='Label of who did an event.'"`
	NoHistory         string `i18n:"hint='Shown when a history is empty.'"`
	TransferOwner     string `i18n:"hint='Button that asks a new owner of a record.'"`
	KeepRole          string `i18n:"hint='Label of the role the former owner keeps.'"`
	InviteFrom        string `i18n:"label='Invited by', hint='Label of who invited the user to a shared record.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "Words shared by the packages.",
	ShareTitle:        "Share",
	ShareWith:         "Share with",
	Permissions:       "Permission",
	ShareWithHint:     "Users to share this record with (one per line, or comma/semicolon separated).",
	CanView:           "Can view",
	CanEdit:           "Can edit",
	CanDelete:         "Can delete",
	SharedWith:        "Shared with",
	NoShares:          "Not shared with anyone yet.",
	Add:               "Share",
	Remove:            "Remove",
	MyInvites:         "Share invites",
	NoInvites:         "No pending invites.",
	AcceptInvite:      "Accept",
	RejectInvite:      "Decline",
	InviteFrom:        "Invited by",
	AccessesTitle:     "Accesses",
	HistoryTitle:      "Access history",
	Subject:           "User or group",
	SubjectUser:       "User",
	SubjectGroup:      "Group",
	Role:              "Role",
	Owner:             "Owner",
	Inherited:         "inherited",
	Revoke:            "Revoke",
	ChangeRole:        "Change role",
	RequestSent:       "Sent to the owner: the sharing waits for their confirmation by e-mail.",
	RoleChanged:       "Role changed.",
	Revoked:           "Access revoked.",
	ConfirmTitle:      "Confirm the request",
	Confirm:           "Confirm",
	Decline:           "Decline",
	Confirmed:         "Confirmed.",
	Declined:          "Declined.",
	RequestedBy:       "Asked by",
	Record:            "Record",
	Event:             "Event",
	When:              "When",
	Who:               "By",
	NoHistory:         "Nothing yet.",
	TransferOwner:     "Transfer ownership",
	KeepRole:          "Role you keep",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "Palavras compartilhadas pelos pacotes.",
	ShareTitle:        "Compartilhar",
	ShareWith:         "Compartilhar com",
	Permissions:       "Permissão",
	ShareWithHint:     "Usuários com quem compartilhar este registro (um por linha, ou separados por vírgula/ponto e vírgula).",
	CanView:           "Pode ver",
	CanEdit:           "Pode editar",
	CanDelete:         "Pode excluir",
	SharedWith:        "Compartilhado com",
	NoShares:          "Ainda não compartilhado com ninguém.",
	Add:               "Compartilhar",
	Remove:            "Remover",
	MyInvites:         "Convites de compartilhamento",
	NoInvites:         "Nenhum convite pendente.",
	AcceptInvite:      "Aceitar",
	RejectInvite:      "Recusar",
	InviteFrom:        "Convidado por",
	AccessesTitle:     "Acessos",
	HistoryTitle:      "Histórico de acessos",
	Subject:           "Usuário ou grupo",
	SubjectUser:       "Usuário",
	SubjectGroup:      "Grupo",
	Role:              "Papel",
	Owner:             "Proprietário",
	Inherited:         "herdado",
	Revoke:            "Revogar",
	ChangeRole:        "Trocar papel",
	RequestSent:       "Enviado ao proprietário: o compartilhamento espera a confirmação dele por e-mail.",
	RoleChanged:       "Papel trocado.",
	Revoked:           "Acesso revogado.",
	ConfirmTitle:      "Confirme o pedido",
	Confirm:           "Confirmar",
	Decline:           "Recusar",
	Confirmed:         "Confirmado.",
	Declined:          "Recusado.",
	RequestedBy:       "Pedido por",
	Record:            "Registro",
	Event:             "Evento",
	When:              "Quando",
	Who:               "Por",
	NoHistory:         "Nada ainda.",
	TransferOwner:     "Transferir propriedade",
	KeepRole:          "Papel que você mantém",
}

func registerMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, I18nSharedKey, Messages_en_US)
	b.RegisterForModules(language.BrazilianPortuguese, I18nSharedKey, Messages_pt_BR)
}

func msgs(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nSharedKey, Messages_en_US).(*Messages)
}
