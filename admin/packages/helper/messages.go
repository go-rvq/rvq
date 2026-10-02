package helper

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "rvq-admin/helper"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	ModuleDescription string           `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	ErrFieldRequired  i18n.ErrorString `i18n:"hint='Error of a required field left empty.'"`
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
		ModuleDescription: "Words shared by the helper components.",
		ErrFieldRequired:  i18n.ErrorString(ErrFieldRequired.Error()),
	}

	Messages_pt_BR = &Messages{
		ModuleDescription: "Palavras compartilhadas pelos componentes auxiliares.",
		ErrFieldRequired:  "Este campo não pode ser vazio",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
