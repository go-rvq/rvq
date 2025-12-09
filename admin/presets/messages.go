package presets

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-rvq/rvq/x/i18n"
)

func MustGetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, CoreI18nModuleKey, DefaultMessages).(*Messages)
}

type TimeFormatMessages struct {
	Date     string
	Time     string
	DateTime string
}

type StrMap map[string]string

func (m StrMap) Get(key string) string {
	return m[key]
}

func (m StrMap) Set(pair ...string) {
	if len(pair)%2 != 0 {
		panic("pairs must have pairs")
	}
	for i := 0; i < len(pair); i += 2 {
		m[pair[i]] = pair[i+1]
	}
}

func (m StrMap) Update(pair ...string) {
	if len(pair)%2 != 0 {
		panic("pairs must have pairs")
	}

	for i := 0; i < len(pair); i += 2 {
		m[pair[i]] = pair[i+1]
	}
}

func (m StrMap) Merge(other StrMap) {
	for k, v := range other {
		m[k] = v
	}
}

type PrinterOptionsMessages struct {
	Title          string
	Print          string
	WithHeaders    string
	WithoutHeaders string
	Preview        string
}

type Messages struct {
	SuccessfullyUpdated        string
	SuccessfullyCreated        string
	SuccessfullyDeleted        string
	SuccessfullyExecutedAction string
	Search                     string
	TheFemaleTitle             string
	TheMaleTitle               string

	YouAreHere                                 string
	New                                        string
	Update                                     string
	Execute                                    string
	Delete                                     string
	Edit                                       string
	FormTitle                                  string
	OK                                         string
	Cancel                                     string
	Clear                                      string
	Create                                     string
	DeleteConfirmationTextTemplate             string
	CreatingFemaleObjectTitleTemplate          string
	EditingTitleTemplate                       string
	CreatingObjectTitleTemplate                string
	EditingObjectTitleTemplate                 string
	ListingObjectTitleTemplate                 string
	DetailingObjectTitleTemplate               string
	FiltersClear                               string
	FiltersAdd                                 string
	FilterApply                                string
	FilterByTemplate                           string
	FiltersDateInTheLast                       string
	FiltersDateEquals                          string
	FiltersDateBetween                         string
	FiltersDateIsAfter                         string
	FiltersDateIsAfterOrOn                     string
	FiltersDateIsBefore                        string
	FiltersDateIsBeforeOrOn                    string
	FiltersDateDays                            string
	FiltersDateMonths                          string
	FiltersDateAnd                             string
	FiltersTo                                  string
	FiltersNumberEquals                        string
	FiltersNumberBetween                       string
	FiltersNumberGreaterThan                   string
	FiltersNumberLessThan                      string
	FiltersNumberAnd                           string
	FiltersStringEquals                        string
	FiltersStringContains                      string
	FiltersMultipleSelectIn                    string
	FiltersMultipleSelectNotIn                 string
	Month                                      string
	MonthNames                                 [time.December + 1]string
	Year                                       string
	PaginationRowsPerPage                      string
	PaginationPageInfo                         string
	PaginationPage                             string
	PaginationOfPage                           string
	ListingNoRecordToShow                      string
	ListingSelectedCountNotice                 string
	ListingClearSelection                      string
	BulkActionNoAvailableRecords               string
	BulkActionSelectedIdsProcessNoticeTemplate string
	ConfirmDialogPromptTitle                   string
	ConfirmDialogPromptText                    string
	Language                                   string
	Colon                                      string
	NotFoundPageNotice                         string
	AddRow                                     string
	PleaseSelectRecord                         string
	PrinterOptions                             PrinterOptionsMessages
	BulkActionConfirmationTextTemplate         string

	TimeFormats TimeFormatMessages

	Common StrMap

	Error               string
	ErrEmptyParamID     i18n.ErrorString
	ErrPermissionDenied i18n.ErrorString

	CopiedToClipboard string
}

func (msgr *Messages) TheTitle(female bool, title string, args ...string) string {
	if female {
		return fmt.Sprintf(msgr.TheFemaleTitle, title)
	}
	return fmt.Sprintf(msgr.TheMaleTitle, title)
}

func (msgr *Messages) DeleteConfirmationText(model, theModelTitle, title string) string {
	return strings.NewReplacer("{model}", model, "{the_model}", theModelTitle, "{title}", title).
		Replace(msgr.DeleteConfirmationTextTemplate)
}

func (msgr *Messages) DeleteConfirmationHtml(model, theModelTitle, title string) string {
	return strings.NewReplacer("{model}", model, "{the_model}", theModelTitle, "{title}", "<b>"+title+"</b>").
		Replace(msgr.DeleteConfirmationTextTemplate)
}

func (msgr *Messages) EditingTitle(label string) string {
	return strings.NewReplacer("{modelName}", label).
		Replace(msgr.EditingTitleTemplate)
}

func (msgr *Messages) CreatingObjectTitle(modelName string, female bool) string {
	tmpl := msgr.CreatingObjectTitleTemplate
	if female && msgr.CreatingFemaleObjectTitleTemplate != "" {
		tmpl = msgr.CreatingFemaleObjectTitleTemplate
	}
	return strings.NewReplacer("{modelName}", modelName).
		Replace(tmpl)
}

func (msgr *Messages) EditingObjectTitle(label string, name string) string {
	return strings.NewReplacer("{id}", name, "{modelName}", label).
		Replace(msgr.EditingObjectTitleTemplate)
}

func (msgr *Messages) ListingObjectTitle(label string) string {
	return strings.NewReplacer("{modelName}", label).
		Replace(msgr.ListingObjectTitleTemplate)
}

func (msgr *Messages) DetailingObjectTitle(label string, name string) string {
	return strings.NewReplacer("{id}", name, "{modelName}", label).
		Replace(msgr.DetailingObjectTitleTemplate)
}

func (msgr *Messages) BulkActionSelectedIdsProcessNotice(ids string) string {
	return strings.NewReplacer("{ids}", ids).
		Replace(msgr.BulkActionSelectedIdsProcessNoticeTemplate)
}

func (msgr *Messages) FilterBy(filter string) string {
	return strings.NewReplacer("{filter}", filter).
		Replace(msgr.FilterByTemplate)
}

func (msgr *Messages) BulkActionConfirmationText(action string, count string) string {
	return strings.NewReplacer("{Action}", action, "{count}", count).
		Replace(msgr.BulkActionConfirmationTextTemplate)
}

func (msgr *Messages) ListingSelectedCountNoticeText(count int) string {
	return strings.NewReplacer("{count}", fmt.Sprint(count)).
		Replace(msgr.ListingSelectedCountNotice)
}

var Messages_en_US = &Messages{
	YouAreHere:                        "You Are Here",
	SuccessfullyUpdated:               "Successfully Updated",
	SuccessfullyCreated:               "Successfully Created",
	SuccessfullyDeleted:               "Successfully Deleted",
	Search:                            "Search",
	New:                               "New",
	Update:                            "Update",
	Execute:                           "Execute",
	Delete:                            "Delete",
	Edit:                              "Edit",
	FormTitle:                         "Form",
	OK:                                "OK",
	Cancel:                            "Cancel",
	Clear:                             "Clear",
	Create:                            "Create",
	DeleteConfirmationTextTemplate:    "Are you sure you want to delete {the_model}: {title}?",
	CreatingObjectTitleTemplate:       "New {modelName}",
	CreatingFemaleObjectTitleTemplate: "New {modelName}",
	EditingTitleTemplate:              "Editing {modelName}",
	EditingObjectTitleTemplate:        "Editing {modelName} {id}",
	ListingObjectTitleTemplate:        "Listing {modelName}",
	DetailingObjectTitleTemplate:      "{modelName} {id}",
	FiltersClear:                      "Clear Filters",
	FiltersAdd:                        "Add Filters",
	FilterApply:                       "Apply",
	FilterByTemplate:                  "Filter by {filter}",
	FiltersDateInTheLast:              "is in the last",
	FiltersDateEquals:                 "is equal to",
	FiltersDateBetween:                "is between",
	FiltersDateIsAfter:                "is after",
	FiltersDateIsAfterOrOn:            "is on or after",
	FiltersDateIsBefore:               "is before",
	FiltersDateIsBeforeOrOn:           "is before or on",
	FiltersDateDays:                   "days",
	FiltersDateMonths:                 "months",
	FiltersDateAnd:                    "and",
	FiltersTo:                         "to",
	FiltersNumberEquals:               "is equal to",
	FiltersNumberBetween:              "between",
	FiltersNumberGreaterThan:          "is greater than",
	FiltersNumberLessThan:             "is less than",
	FiltersNumberAnd:                  "and",
	FiltersStringEquals:               "is equal to",
	FiltersStringContains:             "contains",
	FiltersMultipleSelectIn:           "in",
	FiltersMultipleSelectNotIn:        "not in",
	Month:                             "Month",
	MonthNames: [time.December + 1]string{
		"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	},
	Year:                                       "Year",
	PaginationRowsPerPage:                      "Rows per page: ",
	ListingNoRecordToShow:                      "No records to show",
	ListingSelectedCountNotice:                 "{count} records are selected. ",
	ListingClearSelection:                      "clear selection",
	BulkActionNoAvailableRecords:               "None of the selected records can be executed with this action.",
	BulkActionSelectedIdsProcessNoticeTemplate: "Partially selected records cannot be executed with this action: {ids}.",
	ConfirmDialogPromptText:                    "Are you sure?",
	Language:                                   "Language",
	Colon:                                      ":",
	NotFoundPageNotice:                         "Sorry, the requested page cannot be found. Please check the URL.",
	PleaseSelectRecord:                         "Please select a record",
	AddRow:                                     "Add Row",

	BulkActionConfirmationTextTemplate: "Are you sure you want to <b>{Action}</b> then {count} records?",

	Error:             "ERROR",
	ErrEmptyParamID:   "Empty param ID",
	CopiedToClipboard: "Copied to clipboard",
}

var DefaultMessages = Messages_en_US

var Messages_pt_BR = &Messages{
	YouAreHere:                 "Você está aqui",
	CopiedToClipboard:          "Copiado para a área de transferência!",
	TheFemaleTitle:             "A %s",
	TheMaleTitle:               "O %s",
	SuccessfullyUpdated:        "Atualizado com Sucesso",
	SuccessfullyCreated:        "Cadastrado com Sucesso",
	SuccessfullyDeleted:        "Excluído com Sucesso",
	SuccessfullyExecutedAction: "Ação executada com Sucesso",
	Search:                     "Pesquisa",
	New:                        "Novo",
	Update:                     "Atualizar",
	Execute:                    "Executar",
	Delete:                     "Excluir",
	Edit:                       "Editar",
	FormTitle:                  "Formulário",
	OK:                         "OK",
	Cancel:                     "Cancelar",
	Clear:                      "Limpar",

	Create:                                     "Cadastrar",
	DeleteConfirmationTextTemplate:             "Tem certeza de que deseja excluir {the_model}: {title}?",
	CreatingFemaleObjectTitleTemplate:          "Nova ‹{modelName}›",
	CreatingObjectTitleTemplate:                "Novo ‹{modelName}›",
	EditingTitleTemplate:                       "Editando ‹{modelName}›",
	EditingObjectTitleTemplate:                 "Editando ‹{modelName}› {id}",
	ListingObjectTitleTemplate:                 "Listando ‹{modelName}›",
	DetailingObjectTitleTemplate:               "‹{modelName}› {id}",
	FiltersClear:                               "Limpar Filtros",
	FiltersAdd:                                 "Adicionar Filtro",
	FilterApply:                                "Aplicar",
	FilterByTemplate:                           "Filtrar por {filter}",
	FiltersDateInTheLast:                       "está no último",
	FiltersDateEquals:                          "é igual a",
	FiltersDateBetween:                         "é entre",
	FiltersDateIsAfter:                         "é depois",
	FiltersDateIsAfterOrOn:                     "está depois ou em",
	FiltersDateIsBefore:                        "é antes",
	FiltersDateIsBeforeOrOn:                    "está antes ou em",
	FiltersDateDays:                            "dias",
	FiltersDateMonths:                          "meses",
	FiltersDateAnd:                             "e",
	FiltersTo:                                  "até",
	FiltersNumberEquals:                        "é igual a",
	FiltersNumberBetween:                       "entre",
	FiltersNumberGreaterThan:                   "é maior que",
	FiltersNumberLessThan:                      "é menor que",
	FiltersNumberAnd:                           "e",
	FiltersStringEquals:                        "é igual a",
	FiltersStringContains:                      "contém",
	FiltersMultipleSelectIn:                    "em",
	FiltersMultipleSelectNotIn:                 "fora de de",
	PaginationRowsPerPage:                      "Registros por página: ",
	PaginationPageInfo:                         "{currPageStart}-{currPageEnd} de {total}",
	PaginationPage:                             "Página:",
	PaginationOfPage:                           "de {total}",
	ListingNoRecordToShow:                      "Nenhum registro a ser mostrado",
	ListingSelectedCountNotice:                 "{count} registros selecionados. ",
	ListingClearSelection:                      "limpar seleção",
	BulkActionNoAvailableRecords:               "Nenhum dos registros selecionados podem ser executados com esta ação.",
	BulkActionSelectedIdsProcessNoticeTemplate: "Registros parcialmente selecionados não podem ser executados com esta ação Pesquisar isso no Goo: {ids}.",
	BulkActionConfirmationTextTemplate:         "Tem certeza que deseja executar a ação <b>{Action}</b> nos {count} registros?",
	ConfirmDialogPromptText:                    "Tem certeza?",
	ConfirmDialogPromptTitle:                   "Confirmação",
	Language:                                   "Idioma",
	Colon:                                      ":",
	NotFoundPageNotice:                         "Desculpe, a página solicitada não pode ser encontrada. Verifique o URL.",
	PleaseSelectRecord:                         "Selecione pelo menos um registro.",
	AddRow:                                     "Adicionar",
	Error:                                      "Erro",
	Month:                                      "Mês",
	Year:                                       "Ano",
	ErrEmptyParamID:                            "Parâmetro ID não informado",
	ErrPermissionDenied:                        "Permissão negada",

	PrinterOptions: PrinterOptionsMessages{
		Title:          "Opções de Impressão",
		WithHeaders:    "Com Cabeçalhos",
		WithoutHeaders: "Sem Cabeçalhos",
		Preview:        "Previsualização de Impressão",
		Print:          "Imprimir",
	},

	TimeFormats: TimeFormatMessages{
		Time:     "15:04:05Z07:00",
		Date:     "02/01/2006",
		DateTime: "02/01/2006 15:04:05 Z07:00",
	},

	Common: map[string]string{
		"CreatedAt":           "Cadastro",
		"UpdatedAt":           "Atualização",
		"DeletedAt":           "Exclusão",
		"Title":               "Título",
		"Status":              "Situação",
		"Body":                "Corpo",
		"Cover":               "Destaque",
		"Type":                "Tipo",
		"Live":                "Site",
		"Name":                "Nome",
		"Summary":             "Sumário",
		"Page":                "Página",
		"Multiple Statuses":   "Múltiplas Situações",
		"Path":                "Caminho",
		"Enabled":             "Habilitado",
		"Translate":           "Traduzir",
		"Publication":         "Publicação",
		"Description":         "Descrição",
		"Action":              "Ação",
		"Actions":             "Ações",
		"YouAreHere":          "Voçê está aqui",
		"Size":                "Tamanho",
		"Position":            "Posição",
		"Link":                "URL",
		"LinkQuery":           "Parâmetros da URL",
		"ID":                  "ID",
		"Layout":              "Layout",
		"Galleries":           "Galerias de Imagens",
		"LayoutConfig":        "Configuração do Layout",
		"Config":              "Configuração",
		"TitleWithSlug":       "Endereço",
		"PageOptions":         "Opções da Página",
		"Value":               "Valor",
		"File":                "Arquivo",
		"LocaleCode":          "Idioma",
		"L10nTitle":           "Título em Outros Idiomas",
		"L10nDescription":     "Descrição em Outros Idiomas",
		"L10nLink":            "Endereço em Outros Idiomas",
		"Roles":               "Papéis",
		"PostListing":         "Listagem",
		"PostListingEnabled":  "Listado",
		"PostListingDisabled": "Não Listado",
		"Ext":                 "Extensão",
		"FileName":            "Nome do arquivo",
		"Profile":             "Perfil",
		"OldPassword":         "Senha Atual",
		"NewPassword":         "Nova Senha",
		"ConfirmPassword":     "Repita a Nova Senha",
	},

	MonthNames: [time.December + 1]string{
		"",
		"Janeiro",
		"Fevereiro",
		"Março",
		"Abril",
		"Maio",
		"Junho",
		"Julho",
		"Agosto",
		"Setembro",
		"Outubro",
		"Novembro",
		"Dezembro",
	},
}
