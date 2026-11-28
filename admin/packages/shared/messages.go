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
	MyInvites         string `i18n:"label='My invites', hint='Title of the shares the user was invited to.'"`
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
	RequestSent       string `i18n:"hint='Shown when a sharing waits for the owner.'"`
	ParentRead        string `i18n:"hint='Switch of the sharing of a record under another: the subject reads its parent and the parent\\'s supporting records too (read only).'"`
	Granted           string `i18n:"hint='Shown when a sharing was done at once (the role None: no owner to ask).'"`
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

	EventConfirmationRequested string `i18n:"hint='An event of the history: the owner was asked to confirm by e-mail.'"`
	EventOwnerConfirmed        string `i18n:"hint='An event of the history: the owner confirmed.'"`
	EventOwnerRejected         string `i18n:"hint='An event of the history: the owner declined.'"`
	EventConfirmationExpired   string `i18n:"hint='An event of the history: the confirmation link expired.'"`
	EventInvited               string `i18n:"hint='An event of the history: the user was invited.'"`
	EventAccepted              string `i18n:"hint='An event of the history: the user accepted the invite.'"`
	EventRejected              string `i18n:"hint='An event of the history: the user declined the invite.'"`
	EventGranted               string `i18n:"hint='An event of the history: the access was given.'"`
	EventRoleChangeRequested   string `i18n:"hint='An event of the history: a new role was asked.'"`
	EventRoleChanged           string `i18n:"hint='An event of the history: the role changed.'"`
	EventRevoked               string `i18n:"hint='An event of the history: the access was taken away.'"`
	EventOwnerTransferRequest  string `i18n:"hint='An event of the history: a new owner was asked.'"`
	EventOwnerTransferred      string `i18n:"hint='An event of the history: the record has a new owner.'"`
	EventMemberAdded           string `i18n:"hint='An event of the history: a user joined a group with access.'"`
	EventMemberRemoved         string `i18n:"hint='An event of the history: a user left a group with access.'"`
	KindShare                  string `i18n:"hint='The kind of a request: a new sharing.'"`
	KindRoleChange             string `i18n:"hint='The kind of a request: a new role.'"`
	KindOwnerTransfer          string `i18n:"hint='The kind of a request: a new owner.'"`
	StatePending               string `i18n:"hint='The state of a request or invite: waiting.'"`
	StateConfirmed             string `i18n:"hint='The state of a request: confirmed by the owner.'"`
	StateAccepted              string `i18n:"hint='The state of an invite: accepted.'"`
	StateRejected              string `i18n:"hint='The state of a request or invite: declined.'"`
	StateExpired               string `i18n:"hint='The state of a request: its link expired.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "The sharing of records with users and groups.",
	ShareTitle:        "Share",
	MyInvites:         "Share invites",
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
	RequestSent:       "Sent to the owner: the sharing waits for their confirmation by e-mail.",
	ParentRead:        "See its parent and the parent's supporting records (read only)",
	Granted:           "Done: the record is out of what they reach.",
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

	EventConfirmationRequested: "confirmation asked to the owner",
	EventOwnerConfirmed:        "confirmed by the owner",
	EventOwnerRejected:         "declined by the owner",
	EventConfirmationExpired:   "confirmation expired",
	EventInvited:               "invited",
	EventAccepted:              "invite accepted",
	EventRejected:              "invite declined",
	EventGranted:               "access granted",
	EventRoleChangeRequested:   "role change asked",
	EventRoleChanged:           "role changed",
	EventRevoked:               "access revoked",
	EventOwnerTransferRequest:  "ownership transfer asked",
	EventOwnerTransferred:      "ownership transferred",
	EventMemberAdded:           "joined the group",
	EventMemberRemoved:         "left the group",
	KindShare:                  "Sharing",
	KindRoleChange:             "Role change",
	KindOwnerTransfer:          "Ownership transfer",
	StatePending:               "pending",
	StateConfirmed:             "confirmed",
	StateAccepted:              "accepted",
	StateRejected:              "declined",
	StateExpired:               "expired",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "O compartilhamento de registros com usuários e grupos.",
	ShareTitle:        "Compartilhar",
	MyInvites:         "Convites de compartilhamento",
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
	RequestSent:       "Enviado ao proprietário: o compartilhamento espera a confirmação dele por e-mail.",
	ParentRead:        "Ver o pai e os registros de apoio dele (somente leitura)",
	Granted:           "Feito: o registro saiu do que ele alcança.",
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

	EventConfirmationRequested: "confirmação pedida ao proprietário",
	EventOwnerConfirmed:        "confirmado pelo proprietário",
	EventOwnerRejected:         "recusado pelo proprietário",
	EventConfirmationExpired:   "confirmação expirada",
	EventInvited:               "convidado",
	EventAccepted:              "convite aceito",
	EventRejected:              "convite recusado",
	EventGranted:               "acesso concedido",
	EventRoleChangeRequested:   "troca de papel pedida",
	EventRoleChanged:           "papel trocado",
	EventRevoked:               "acesso revogado",
	EventOwnerTransferRequest:  "transferência de propriedade pedida",
	EventOwnerTransferred:      "propriedade transferida",
	EventMemberAdded:           "entrou no grupo",
	EventMemberRemoved:         "saiu do grupo",
	KindShare:                  "Compartilhamento",
	KindRoleChange:             "Troca de papel",
	KindOwnerTransfer:          "Transferência de propriedade",
	StatePending:               "pendente",
	StateConfirmed:             "confirmado",
	StateAccepted:              "aceito",
	StateRejected:              "recusado",
	StateExpired:               "expirado",
}

func registerMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, I18nSharedKey, Messages_en_US)
	b.RegisterForModules(language.BrazilianPortuguese, I18nSharedKey, Messages_pt_BR)
}

func msgs(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nSharedKey, Messages_en_US).(*Messages)
}

// EventLabel is the name of the event e of the history in the language of
// ctx (e itself when unknown).
func EventLabel(ctx context.Context, e string) string { return msgs(ctx).EventName(e) }

// EventName is the name of the event e of the history (e itself when
// unknown).
func (m *Messages) EventName(e string) string {
	switch e {
	case EventConfirmationRequested:
		return m.EventConfirmationRequested
	case EventOwnerConfirmed:
		return m.EventOwnerConfirmed
	case EventOwnerRejected:
		return m.EventOwnerRejected
	case EventConfirmationExpired:
		return m.EventConfirmationExpired
	case EventInvited:
		return m.EventInvited
	case EventAccepted:
		return m.EventAccepted
	case EventRejected:
		return m.EventRejected
	case EventGranted:
		return m.EventGranted
	case EventRoleChangeRequested:
		return m.EventRoleChangeRequested
	case EventRoleChanged:
		return m.EventRoleChanged
	case EventRevoked:
		return m.EventRevoked
	case EventOwnerTransferRequest:
		return m.EventOwnerTransferRequest
	case EventOwnerTransferred:
		return m.EventOwnerTransferred
	case EventMemberAdded:
		return m.EventMemberAdded
	case EventMemberRemoved:
		return m.EventMemberRemoved
	}
	return e
}

// KindName is the name of the kind of a request.
func (m *Messages) KindName(k string) string {
	switch k {
	case RequestShare:
		return m.KindShare
	case RequestRoleChange:
		return m.KindRoleChange
	case RequestOwnerTransfer:
		return m.KindOwnerTransfer
	}
	return k
}

// StateName is the name of the state of a request or of an invite.
func (m *Messages) StateName(st string) string {
	switch st {
	case StatePending:
		return m.StatePending
	case StateConfirmed:
		return m.StateConfirmed
	case StateAccepted:
		return m.StateAccepted
	case StateRejected:
		return m.StateRejected
	case StateExpired:
		return m.StateExpired
	}
	return st
}
