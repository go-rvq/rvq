package role

import (
	"context"

	"github.com/qor5/x/v3/i18n"
)

const I18nRoleKey = "I18nRoleKey"

type Messages struct {
	Roles string
	Role  string
}

var Messages_en_US = &Messages{
	Roles: "Roles",
	Role:  "Role",
}

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nRoleKey, Messages_en_US).(*Messages)
}
