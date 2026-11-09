package user_docs

import (
	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/x/i18n"
)

// MessagesKey is the module of the words of the documentation of the package.
const MessagesKey i18n.ModuleKey = "rvq-admin/packages/notifications/user_docs"

// Messages are the words of the documentation of the package.
type Messages struct {
	userdocs.PackageMessages
}

var Messages_en_US = &Messages{
	PackageMessages: userdocs.PackageMessagesEn("Notifications",
		"What the system and the other users tell you: the bell, your notifications, the widget of your profile."),
}

var Messages_pt_BR = &Messages{
	PackageMessages: userdocs.PackageMessagesPtBR("Notificações",
		"O que o sistema e os outros usuários lhe dizem: o sino, as suas notificações, o widget do seu perfil."),
}
