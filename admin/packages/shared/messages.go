package shared

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// PermShare is the perm verb (and detail-action name) guarding the share action.
const PermShare = "share"

// I18nSharedKey is the i18n module key of the shared package.
const I18nSharedKey i18n.ModuleKey = "SharedManager"

// Messages are the sharing UI labels.
type Messages struct {
	ShareTitle    string
	ShareWith     string
	ShareWithHint string
	Permissions   string
	CanView       string
	CanEdit       string
	CanDelete     string
	SharedWith    string
	NoShares      string
	Add           string
	Remove        string
	MyInvites     string
	NoInvites     string
	AcceptInvite  string
	RejectInvite  string
	InviteFrom    string
}

var Messages_en_US = &Messages{
	ShareTitle:    "Share",
	ShareWith:     "Share with",
	Permissions:   "Permission",
	ShareWithHint: "Users to share this record with (one per line, or comma/semicolon separated).",
	CanView:       "Can view",
	CanEdit:       "Can edit",
	CanDelete:     "Can delete",
	SharedWith:    "Shared with",
	NoShares:      "Not shared with anyone yet.",
	Add:           "Share",
	Remove:        "Remove",
	MyInvites:     "Share invites",
	NoInvites:     "No pending invites.",
	AcceptInvite:  "Accept",
	RejectInvite:  "Decline",
	InviteFrom:    "Invited by",
}

var Messages_pt_BR = &Messages{
	ShareTitle:    "Compartilhar",
	ShareWith:     "Compartilhar com",
	Permissions:   "Permissão",
	ShareWithHint: "Usuários com quem compartilhar este registro (um por linha, ou separados por vírgula/ponto e vírgula).",
	CanView:       "Pode ver",
	CanEdit:       "Pode editar",
	CanDelete:     "Pode excluir",
	SharedWith:    "Compartilhado com",
	NoShares:      "Ainda não compartilhado com ninguém.",
	Add:           "Compartilhar",
	Remove:        "Remover",
	MyInvites:     "Convites de compartilhamento",
	NoInvites:     "Nenhum convite pendente.",
	AcceptInvite:  "Aceitar",
	RejectInvite:  "Recusar",
	InviteFrom:    "Convidado por",
}

func registerMessages(b *i18n.Builder) {
	b.RegisterForModules(language.AmericanEnglish, I18nSharedKey, Messages_en_US)
	b.RegisterForModules(language.BrazilianPortuguese, I18nSharedKey, Messages_pt_BR)
}

func msgs(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nSharedKey, Messages_en_US).(*Messages)
}
