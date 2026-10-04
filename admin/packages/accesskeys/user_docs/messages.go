package user_docs

import (
	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/x/i18n"
)

// MessagesKey is the module of the words of the documentation of the package.
const MessagesKey i18n.ModuleKey = "rvq-admin/packages/accesskeys/user_docs"

// Messages are the words of the documentation of the package.
type Messages struct {
	userdocs.PackageMessages
}

var Messages_en_US = &Messages{
	PackageMessages: userdocs.PackageMessagesEn("Access keys",
		"Codes for automations — git, WebDAV, scripts — that act in a user's name, only within what the key allows."),
}

var Messages_pt_BR = &Messages{
	PackageMessages: userdocs.PackageMessagesPtBR("Chaves de acesso",
		"Códigos para automações — git, WebDAV, scripts — que agem em nome de um usuário, só dentro do que a chave permite."),
}
