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
	ModuleDescription     string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	NothingWritten        string `i18n:"hint='Shown for a part of the admin nothing was written about yet.'"`
	SeeAlso               string `i18n:"hint='Title of the list of the parts inside the one shown.'"`
	EditDocuments         string `i18n:"hint='Link to the form of the documents of the package of the page shown.'"`
	RenderError           string `i18n:"hint='Shown when a document could not be shown.', fields=(;'%s'='what went wrong')"`
	ResetDocuments        string `i18n:"hint='Button that brings the documents of a package back to the initial ones.'"`
	Actions               string `i18n:"hint='Node of the tree with the actions of a part of the admin.'"`
	Pages                 string `i18n:"hint='Node of the tree with the pages of a part of the admin.'"`
	FormNew               string `i18n:"hint='Node of the tree with the form of a new record of a model.'"`
	FormEdit              string `i18n:"hint='Node of the tree with the form of an edit of a record of a model.'"`
	FormDetail            string `i18n:"hint='Node of the tree with the detail of a record of a model.'"`
	Field                 string `i18n:"hint='Column of the table of the fields of a form: the field.'"`
	Description           string `i18n:"hint='Column of the table of the fields of a form: what the field is.'"`
	Conditional           string `i18n:"hint='Said of a field shown depending on the record.'"`
	NoFields              string `i18n:"hint='Shown for a form with no field.'"`
	Children              string `i18n:"hint='Title of the models nested in a record, in the menu of its detail.'"`
	NothingInMenu         string `i18n:"hint='Shown for a detail whose menu has nothing.'"`
	RuleHeading           string `i18n:"hint='Heading of the section of the document of an action that says when it is available: the text under it is shown in the menu of the detail. The documents use the same.'"`
	RuleUndocumented      string `i18n:"hint='Shown for an action available by a rule its document does not say.'"`
	Permissions           string `i18n:"hint='Node of the tree with the permissions of a part of the admin.'"`
	PermNone              string `i18n:"hint='Shown for a part with no permission of its own.'"`
	PermResources         string `i18n:"hint='Title of the resources of a part.'"`
	PermBy                string `i18n:"hint='Column of the table of the resources: how the resource names the part.'"`
	PermResource          string `i18n:"hint='Column with a resource of a permission.'"`
	PermUnique            string `i18n:"hint='Says a resource names the part by its unique name: it decides first.'"`
	PermByGroups          string `i18n:"hint='Says a resource names the part through the groups of the menu.'"`
	PermRecords           string `i18n:"hint='Title of the permissions of the records of a model.'"`
	PermWhat              string `i18n:"hint='Column with what a permission allows.'"`
	PermVerb              string `i18n:"hint='Column with the verb (action) of a permission.'"`
	PermList              string `i18n:"hint='Permission to list the records.'"`
	PermGet               string `i18n:"hint='Permission to see a record.'"`
	PermCreate            string `i18n:"hint='Permission to create a record.'"`
	PermUpdate            string `i18n:"hint='Permission to edit a record.'"`
	PermDelete            string `i18n:"hint='Permission to delete a record.'"`
	PermDeleteWithRelated string `i18n:"hint='Permission to delete a record with the records related to it.'"`
	PermFields            string `i18n:"hint='Title of the permissions of the fields of a model.'"`
	PermFieldsHint        string `i18n:"hint='Explains the permissions of the fields.'"`
	PermSections          string `i18n:"hint='Title of the permissions of the sections of a detail.'"`
	PermSection           string `i18n:"hint='Column with a section of a detail.'"`
	PermSectionsHint      string `i18n:"hint='Explains the permissions of the sections.'"`
	PermParts             string `i18n:"hint='Title of the parts of a group, with their unique names.'"`
	PermScopeAdminDesc    string `i18n:"hint='Description of the scope of the admin in the tree of the permissions.'"`
	PermScopeAdmin        string `i18n:"hint='Name of the scope of the admin (its resources begin with admin:) in the tree of the permissions.'"`
	PermDescGroup         string `i18n:"hint='Description of a group of the menu in the tree of the permissions.', fields=(;'%s'='the group')"`
	PermDescModel         string `i18n:"hint='Description of a model (its listing) in the tree of the permissions.', fields=(;'%s'='the model')"`
	PermDescNested        string `i18n:"hint='Description of a model nested in a record of another, in the tree of the permissions.', fields=(;'%s'='the model; the second %s, the parent')"`
	PermDescSingleton     string `i18n:"hint='Description of a model with a single record in the tree of the permissions.', fields=(;'%s'='the model')"`
	PermDescRecord        string `i18n:"hint='Description of the records of a model in the tree of the permissions.', fields=(;'%s'='the model')"`
	PermDescField         string `i18n:"hint='Description of a field in the tree of the permissions.', fields=(;'%s'='the field; the second %s, its record or field')"`
	PermDescSection       string `i18n:"hint='Description of a section of a detail in the tree of the permissions.', fields=(;'%s'='the section; the second %s, the record')"`
	PermDescPage          string `i18n:"hint='Description of a page of the admin in the tree of the permissions.', fields=(;'%s'='the page')"`
	PermDescPageOf        string `i18n:"hint='Description of a page of a model or record in the tree of the permissions.', fields=(;'%s'='the page; the second %s, its model or record')"`
	PermDescCheck         string `i18n:"hint='Description of a permission of its own (a verifier) in the tree of the permissions.', fields=(;'%s'='what it allows')"`
	PermUntitled          string `i18n:"hint='Label of a part of the tree of the permissions with no title of its own.', fields=(;'%s'='its name')"`
	PermPartsHint         string `i18n:"hint='Explains the parts of a group.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription:     "The documentation of the admin: the page that explains each part of it.",
	NothingWritten:        "Nothing was written about this yet.",
	SeeAlso:               "In this part",
	EditDocuments:         "Edit the documents",
	RenderError:           "This document could not be shown: %s",
	ResetDocuments:        "Reset to the initial documents",
	Actions:               "Actions",
	Pages:                 "Pages",
	FormNew:               "New record",
	FormEdit:              "Edit",
	FormDetail:            "Detail",
	Field:                 "Field",
	Description:           "Description",
	Conditional:           "Shown depending on the record.",
	NoFields:              "No fields.",
	Children:              "Nested records",
	NothingInMenu:         "The menu has nothing.",
	RuleHeading:           "When it is available",
	RuleUndocumented:      "Available depending on the record (the rule is not documented).",
	Permissions:           "Permissions",
	PermNone:              "This part has no permission of its own.",
	PermResources:         "Resources",
	PermBy:                "By",
	PermResource:          "Resource",
	PermUnique:            "Its unique name (decides first)",
	PermByGroups:          "Its groups",
	PermRecords:           "Records",
	PermWhat:              "What",
	PermVerb:              "Verb",
	PermList:              "List",
	PermGet:               "See a record",
	PermCreate:            "Create",
	PermUpdate:            "Edit",
	PermDelete:            "Delete",
	PermDeleteWithRelated: "Delete with the related",
	PermFields:            "Fields",
	PermFieldsHint:        "A field allowed to see but not to edit is shown read only; one not allowed to see is not shown.",
	PermSections:          "Sections",
	PermSection:           "Section",
	PermSectionsHint:      "A section is seen and edited in place by its own permission; the fields it writes still ask theirs.",
	PermParts:             "Inside the group",
	PermScopeAdminDesc:    "Everything of the admin: the groups of the menu, its models with their records, fields, sections and actions, and its pages",
	PermDescGroup:         "Everything in the group %s of the menu",
	PermDescModel:         "The listing of %s: listing and creating records, and its actions",
	PermDescNested:        "The records of %s inside a record of %s",
	PermDescSingleton:     "The single record of %s: seeing and editing it, and its actions",
	PermDescRecord:        "Any record of %s (<*>), or one by its id (<7>): seeing, editing, deleting it, and its actions",
	PermDescField:         "The field %s of %s, in the forms that hold it",
	PermDescSection:       "The section %s of the detail of %s, seen and edited in place",
	PermDescPage:          "The page %s",
	PermDescPageOf:        "The page %s of %s",
	PermDescCheck:         "A permission of its own: %s",
	PermUntitled:          "Permission %s",
	PermScopeAdmin:        "Admin",
	PermPartsHint:         "What the group's resource allows covers all of them; a permission by the unique name of one of them decides before it — and a deny wins.",
}

var Messages_pt_BR = &Messages{
	ModuleDescription:     "A documentação do admin: a página que explica cada parte dele.",
	NothingWritten:        "Ainda não há nada escrito sobre isto.",
	SeeAlso:               "Nesta parte",
	EditDocuments:         "Editar os documentos",
	RenderError:           "Este documento não pôde ser mostrado: %s",
	ResetDocuments:        "Restaurar os documentos iniciais",
	Actions:               "Ações",
	Pages:                 "Páginas",
	FormNew:               "Cadastro",
	FormEdit:              "Edição",
	FormDetail:            "Detalhe",
	Field:                 "Campo",
	Description:           "Descrição",
	Conditional:           "Aparece conforme o registro.",
	NoFields:              "Nenhum campo.",
	Children:              "Registros filhos",
	NothingInMenu:         "O menu não tem nada.",
	RuleHeading:           "Quando está disponível",
	RuleUndocumented:      "Disponível conforme o registro (a regra não está documentada).",
	Permissions:           "Permissões",
	PermNone:              "Esta parte não tem permissão própria.",
	PermResources:         "Recursos",
	PermBy:                "Por",
	PermResource:          "Recurso",
	PermUnique:            "Nome único (decide primeiro)",
	PermByGroups:          "Grupos",
	PermRecords:           "Registros",
	PermWhat:              "O quê",
	PermVerb:              "Verbo",
	PermList:              "Listar",
	PermGet:               "Ver um registro",
	PermCreate:            "Criar",
	PermUpdate:            "Editar",
	PermDelete:            "Excluir",
	PermDeleteWithRelated: "Excluir com os relacionados",
	PermFields:            "Campos",
	PermFieldsHint:        "Um campo permitido para ver mas não para editar aparece somente leitura; um não permitido para ver não aparece.",
	PermSections:          "Seções",
	PermSection:           "Seção",
	PermSectionsHint:      "Uma seção é vista e editada no lugar pela própria permissão; os campos que ela grava ainda pedem as deles.",
	PermParts:             "Dentro do grupo",
	PermScopeAdminDesc:    "Tudo do admin: os grupos do menu, os modelos com seus registros, campos, seções e ações, e as páginas",
	PermDescGroup:         "Tudo o que está no grupo %s do menu",
	PermDescModel:         "A listagem de %s: listar e criar registros, e as ações dela",
	PermDescNested:        "Os registros de %s dentro de um registro de %s",
	PermDescSingleton:     "O registro único de %s: ver e editar, e as ações dele",
	PermDescRecord:        "Qualquer registro de %s (<*>), ou um pelo id (<7>): ver, editar, excluir, e as ações dele",
	PermDescField:         "O campo %s de %s, nos formulários que o contêm",
	PermDescSection:       "A seção %s do detalhe de %s, vista e editada no lugar",
	PermDescPage:          "A página %s",
	PermDescPageOf:        "A página %s de %s",
	PermDescCheck:         "Uma permissão própria: %s",
	PermUntitled:          "Permissão %s",
	PermScopeAdmin:        "Administração",
	PermPartsHint:         "O que o recurso do grupo permite vale para todos eles; uma permissão pelo nome único de um deles decide antes — e negar prevalece.",
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

// PackageMessagesPtBR are the words of a package's documentation in
// Portuguese, its title and description given.
func PackageMessagesPtBR(title, description string) PackageMessages {
	return PackageMessages{
		ModuleDescription: description,
		Title:             title,
		Path:              "Documento",
		PathHint:          "Onde o documento fica, na documentação do pacote.",
		Content:           "Conteúdo",
		ContentHint:       "Markdown; {%= admin.model(\"id\").link %} cria um link para uma parte do admin.",
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
