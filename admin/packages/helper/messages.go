package helper

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "rqv-admin/helper"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	ErrFieldRequired i18n.ErrorString
}

func (m *Messages) FromError(err error) error {
	switch err {
	case ErrFieldRequired:
		return m.ErrFieldRequired
	default:
		return err
	}
}

var (
	Messages_en_US = &Messages{
		ErrFieldRequired: i18n.ErrorString(ErrFieldRequired.Error()),
	}

	Messages_pt_BR = &Messages{
		ErrFieldRequired: "Este campo não pode ser vazio",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_pt_BR).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
