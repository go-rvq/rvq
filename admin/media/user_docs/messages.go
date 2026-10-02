package user_docs

import (
	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/x/i18n"
)

// MessagesKey is the module of the words of the documentation of the package.
const MessagesKey i18n.ModuleKey = "rvq-admin/media/user_docs"

// Messages are the words of the documentation of the package.
type Messages struct {
	userdocs.PackageMessages
}

var Messages_en_US = &Messages{
	PackageMessages: userdocs.PackageMessagesEn("Media library",
		"The images and files of the site: uploading, choosing, cropping."),
}
