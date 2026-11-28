package helper

import (
	"context"
	"sync"

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

	SelectorSelect  string `i18n:"hint='The chip that opens the dialog to choose a record.'"`
	SelectorChange  string `i18n:"hint='The chip that opens the dialog to change the record chosen.'"`
	SelectorSearch  string `i18n:"hint='The search field of the dialog of a selector.'"`
	SelectorNoItems string `i18n:"hint='The list of the dialog of a selector when nothing is found.'"`
	SelectorClose   string `i18n:"hint='The button that closes the dialog of a selector.'"`
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
		SelectorSelect:    "Select",
		SelectorChange:    "Change",
		SelectorSearch:    "Search",
		SelectorNoItems:   "Nothing found",
		SelectorClose:     "Close",
	}

	Messages_pt_BR = &Messages{
		ModuleDescription: "Palavras compartilhadas pelos componentes auxiliares.",
		ErrFieldRequired:  "Este campo não pode ser vazio",
		SelectorSelect:    "Selecionar",
		SelectorChange:    "Trocar",
		SelectorSearch:    "Pesquisar",
		SelectorNoItems:   "Nada encontrado",
		SelectorClose:     "Fechar",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

var (
	configuredMu sync.Mutex
	configured   = map[*i18n.Builder]bool{}
)

// ensureMessages registers the messages of the package in b once.
// EnsureMessages registers the messages of the selectors in b once (a
// RecordSelect's: before the admin serves).
func EnsureMessages(b *i18n.Builder) { ensureMessages(b) }

func ensureMessages(b *i18n.Builder) {
	if b == nil {
		return
	}
	configuredMu.Lock()
	defer configuredMu.Unlock()
	if configured[b] {
		return
	}
	configured[b] = true
	ConfigureMessages(b)
}
