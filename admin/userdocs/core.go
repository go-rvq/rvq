package userdocs

import (
	"embed"
	"io/fs"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// coreFS is the documentation's own documents: the forms of every model
// (forms/new.md, forms/edit.md, forms/detail.md), which a model's own
// documents may stand in for; the policies of the admin
// (guides/policies/…), which an application puts in its tree (Custom).
//
//go:embed user_docs/en user_docs/pt-BR
var coreFS embed.FS

// CorePackage is the package of the documentation's own documents.
const CorePackage = "github.com/go-rvq/rvq/admin/userdocs"

// coreMessagesKey is the module of the words of the documentation's own
// documents.
const coreMessagesKey i18n.ModuleKey = "rvq-admin/userdocs/user_docs"

type coreMessages struct{ PackageMessages }

// coreSource is the source of the documentation's own documents.
func coreSource() *Source {
	sub, err := fs.Sub(coreFS, "user_docs")
	if err != nil {
		panic(err)
	}
	return &Source{
		Package:     CorePackage,
		FS:          sub,
		MessagesKey: coreMessagesKey,
		Messages: map[language.Tag]i18n.Messages{
			language.English: &coreMessages{PackageMessagesEn("Forms",
				"The forms of every model: of a new record, of an edit, and the detail with its menu.")},
			language.BrazilianPortuguese: &coreMessages{PackageMessagesPtBR("Formulários",
				"Os formulários de todo model: de um novo registro, de uma edição, e o detalhe com o menu dele.")},
		},
	}
}
