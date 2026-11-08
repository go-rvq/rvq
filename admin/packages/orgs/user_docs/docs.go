// Package user_docs is the documentation of github.com/go-rvq/rvq/admin/packages/orgs for whoever uses the
// admin (userdocs), in en/ and pt-BR/: the organizations and their sharing.
//
// The model of the organizations is out of the menu (its item is the
// organization chosen, at the top): the application documents it by
// userdocs.Builder.Models(admin.OrganizacaoModel(p)).
package user_docs

import (
	"embed"

	"github.com/go-rvq/rvq/admin/userdocs"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// FS is the documents, in the directory of their language.
//
//go:embed en pt-BR
var FS embed.FS

// Package is the package documented.
const Package = "github.com/go-rvq/rvq/admin/packages/orgs"

// Source is the documentation of the package.
func Source() *userdocs.Source {
	return &userdocs.Source{
		Package:     Package,
		FS:          FS,
		MessagesKey: MessagesKey,
		Messages: map[language.Tag]i18n.Messages{
			language.English:             Messages_en_US,
			language.BrazilianPortuguese: Messages_pt_BR,
		},
	}
}
