package user_docs

import (
	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/x/i18n"
)

// MessagesKey is the module of the words of the documentation of the package.
const MessagesKey i18n.ModuleKey = "rvq-admin/packages/orgs/user_docs"

// Messages are the words of the documentation of the package.
type Messages struct {
	userdocs.PackageMessages
}

var Messages_en_US = &Messages{
	PackageMessages: userdocs.PackageMessagesEn("Organizations",
		"The organizations — what is under each one — and their sharing with users and groups, confirmed by the owner."),
}

var Messages_pt_BR = &Messages{
	PackageMessages: userdocs.PackageMessagesPtBR("Organizações",
		"As organizações — o que fica sob cada uma — e o compartilhamento delas com usuários e grupos, confirmado pelo proprietário."),
}
