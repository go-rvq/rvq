package user_docs

import (
	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/x/i18n"
)

// MessagesKey is the module of the words of the documentation of the package.
const MessagesKey i18n.ModuleKey = "rvq-admin/packages/people/user_docs"

// Messages are the words of the documentation of the package.
type Messages struct {
	userdocs.PackageMessages
}

var Messages_en_US = &Messages{
	PackageMessages: userdocs.PackageMessagesEn("People",
		"The people of each organization: its registry of individuals and companies."),
}

var Messages_pt_BR = &Messages{
	PackageMessages: userdocs.PackageMessagesPtBR("Pessoas",
		"As pessoas de cada organização: o cadastro dela de pessoas físicas e jurídicas."),
}
