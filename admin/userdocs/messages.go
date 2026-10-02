package userdocs

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// MessagesKey is the module of the words of the documentation page.
const MessagesKey i18n.ModuleKey = "rvq-admin/userdocs"

// GetMessages are the words of the documentation page of ctx's language.
func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

// Messages are the words of the documentation page.
type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Contents          string `i18n:"hint='Title of the tree of the documentation.'"`
	NothingWritten    string `i18n:"hint='Shown for a part of the admin nothing was written about yet.'"`
	SeeAlso           string `i18n:"hint='Title of the list of the parts inside the one shown.'"`
	EditDocuments     string `i18n:"hint='Link to the form of the documents of the package of the page shown.'"`
	RenderError       string `i18n:"hint='Shown when a document could not be shown.', fields=(;'%s'='what went wrong')"`
	ResetDocuments    string `i18n:"hint='Button that brings the documents of a package back to the initial ones.'"`
	Actions           string `i18n:"hint='Node of the tree with the actions of a part of the admin.'"`
	Pages             string `i18n:"hint='Node of the tree with the pages of a part of the admin.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "The documentation of the admin: the page that explains each part of it.",
	Contents:          "Contents",
	NothingWritten:    "Nothing was written about this yet.",
	SeeAlso:           "In this part",
	EditDocuments:     "Edit the documents",
	RenderError:       "This document could not be shown: %s",
	ResetDocuments:    "Reset to the initial documents",
	Actions:           "Actions",
	Pages:             "Pages",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "A documentação do admin: a página que explica cada parte dele.",
	Contents:          "Conteúdo",
	NothingWritten:    "Ainda não há nada escrito sobre isto.",
	SeeAlso:           "Nesta parte",
	EditDocuments:     "Editar os documentos",
	RenderError:       "Este documento não pôde ser mostrado: %s",
	ResetDocuments:    "Restaurar os documentos iniciais",
	Actions:           "Ações",
	Pages:             "Páginas",
}

// PackageMessages are the words of the documentation of a package: what its
// user_docs/messages.go declares, embedding them —
//
//	type Messages struct{ userdocs.PackageMessages }
//
// — with the values of its language. They are the labels and the hints of the
// form of its documents.
type PackageMessages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Title             string `i18n:"hint='The title of the documentation of the package.'"`
	Path              string `i18n:"hint='Label of the path of a document of the package.'"`
	PathHint          string `i18n:"hint='Hint of the path of a document of the package.'"`
	Content           string `i18n:"hint='Label of the content of a document of the package.'"`
	ContentHint       string `i18n:"hint='Hint of the content of a document of the package.'"`
}

// PackageMessagesEn are the words of a package's documentation in English,
// its title and description given.
func PackageMessagesEn(title, description string) PackageMessages {
	return PackageMessages{
		ModuleDescription: description,
		Title:             title,
		Path:              "Document",
		PathHint:          "Where the document is, in the documentation of the package.",
		Content:           "Content",
		ContentHint:       "Markdown; {%= admin.model(\"id\").link %} links to a part of the admin.",
	}
}

// packageMessages are the words of src's documentation in ctx's language.
func packageMessages(ctx context.Context, src *Source) PackageMessages {
	var def i18n.Messages
	for _, m := range src.Messages {
		def = m
		if pm, ok := asPackageMessages(m); ok && pm.Title != "" {
			break
		}
	}
	if m, ok := asPackageMessages(i18n.MustGetModuleMessages(ctx, src.MessagesKey, def)); ok {
		return m
	}
	return PackageMessagesEn(src.Package, "")
}

// packageMessagesOf is the PackageMessages a package's messages embed.
type packageMessagesOf interface {
	packageMessages() PackageMessages
}

func (m *PackageMessages) packageMessages() PackageMessages { return *m }

func asPackageMessages(m any) (PackageMessages, bool) {
	if p, ok := m.(packageMessagesOf); ok {
		return p.packageMessages(), true
	}
	return PackageMessages{}, false
}

// ConfigureMessages registers the words of the documentation page.
func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModule(language.English, MessagesKey, Messages_en_US).
		RegisterForModule(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
