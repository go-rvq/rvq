package messages

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// Key is the i18n module key of the orgs module.
const Key i18n.ModuleKey = "orgs"

// Register registers the org messages for every supported language.
func Register(b *i18n.Builder) {
	b.RegisterForModules(language.BrazilianPortuguese, Key, Messages_pt_BR)
	b.RegisterForModules(language.English, Key, Messages_en_US)
}

// Get returns the org messages for the request context, falling back to pt-BR.
func Get(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, Key, Messages_pt_BR).(*Messages)
}
