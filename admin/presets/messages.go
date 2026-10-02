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
	Date     string `i18n:"hint='How a date is written, as a Go layout: the reference date 2006-01-02 written the way wanted.'"`
	Time     string `i18n:"hint='How a time is written, as a Go layout: the reference time 15:04:05 Z07:00 written the way wanted.'"`
	DateTime string `i18n:"label='Date and time', hint='How a date with its time is written, as a Go layout (reference 2006-01-02 15:04:05 Z07:00).'"`
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
	Title          string `i18n:"hint='Title of the print options.'"`
	Print          string `i18n:"hint='Button that prints.'"`
	WithHeaders    string `i18n:"hint='Option that prints the headers.'"`
	WithoutHeaders string `i18n:"hint='Option that prints without the headers.'"`
	Preview        string `i18n:"hint='Title of the preview of a print.'"`
}

type Messages struct {
	SuccessfullyUpdated        string `i18n:"hint='Shown once a record was saved.'"`
	SuccessfullyCreated        string `i18n:"hint='Shown once a record was created.'"`
	SuccessfullyDeleted        string `i18n:"hint='Shown once a record was deleted.'"`
	SuccessfullyExecutedAction string `i18n:"hint='Shown once an action ran.'"`
	Search                     string `i18n:"hint='Placeholder of the search of a listing.'"`
	TheFemaleTitle             string `i18n:"label='The (feminine)', hint='A title with its article, for a feminine noun.', fields=(;'%s'='the title')"`
	TheMaleTitle               string `i18n:"label='The (masculine)', hint='A title with its article, for a masculine noun.', fields=(;'%s'='the title')"`

	YouAreHere                                 string                 `i18n:"hint='Label of the breadcrumb (where the user is).'"`
	New                                        string                 `i18n:"hint='Button that creates a record.'"`
	Update                                     string                 `i18n:"hint='Button that saves a record.'"`
	Execute                                    string                 `i18n:"hint='Button that runs an action.'"`
	Delete                                     string                 `i18n:"hint='Button that deletes a record.'"`
	DeleteRelated                              string                 `i18n:"hint='Option that deletes the related records too.'"`
	DeleteRelatedHint                          string                 `i18n:"label='Delete related: hint', hint='Hint of the option that deletes the related records too.'"`
	ShowRelatedItemsTitle                      string                 `i18n:"hint='Button that shows the records related to the one being deleted.'"`
	RelatedItemsForDeletionActionTitle         string                 `i18n:"label='Related items for deletion', hint='Title of the records deleted along with a record.'"`
	Edit                                       string                 `i18n:"hint='Button that opens the form of a record.'"`
	FormTitle                                  string                 `i18n:"hint='Title of a form with no other.'"`
	OK                                         string                 `i18n:"hint='Button that confirms.'"`
	Cancel                                     string                 `i18n:"hint='Button that closes without doing anything.'"`
	Clear                                      string                 `i18n:"hint='Button that empties a field.'"`
	Create                                     string                 `i18n:"hint='Button that saves a new record.'"`
	DeleteConfirmationTextTemplate             string                 `i18n:"label='Delete confirmation', hint='Asks to confirm a deletion.', fields=(;'{the_model}'='the model, with its article', '{title}'='the record\\'s title')"`
	CreatingFemaleObjectTitleTemplate          string                 `i18n:"label='New record title (feminine)', hint='Title of the creation of a record of a feminine model.', fields=(;'{modelName}'='the model')"`
	EditingTitleTemplate                       string                 `i18n:"label='Editing title', hint='Title of the form of a record.', fields=(;'{modelName}'='the model')"`
	CreatingObjectTitleTemplate                string                 `i18n:"label='New record title', hint='Title of the creation of a record.', fields=(;'{modelName}'='the model')"`
	EditingObjectTitleTemplate                 string                 `i18n:"label='Editing record title', hint='Title of the form of a record.', fields=(;'{modelName}'='the model', '{id}'='the record')"`
	ListingObjectTitleTemplate                 string                 `i18n:"label='Listing title', hint='Title of a listing.', fields=(;'{modelName}'='the model')"`
	DetailingObjectTitleTemplate               string                 `i18n:"label='Detail title', hint='Title of the detail of a record.', fields=(;'{modelName}'='the model', '{id}'='the record')"`
	FiltersClear                               string                 `i18n:"hint='Button that removes every filter.'"`
	FiltersAdd                                 string                 `i18n:"hint='Button that opens the filters.'"`
	FilterApply                                string                 `i18n:"hint='Button that applies a filter.'"`
	FilterByTemplate                           string                 `i18n:"label='Filter by', hint='Title of a filter.', fields=(;'{filter}'='its name')"`
	FiltersDateInTheLast                       string                 `i18n:"label='Date: in the last', hint='Date filter operator: within the last days or months.'"`
	FiltersDateEquals                          string                 `i18n:"label='Date: equals', hint='Date filter operator: on the day.'"`
	FiltersDateBetween                         string                 `i18n:"label='Date: between', hint='Date filter operator: between two days.'"`
	FiltersDateIsAfter                         string                 `i18n:"label='Date: after', hint='Date filter operator: after the day.'"`
	FiltersDateIsAfterOrOn                     string                 `i18n:"label='Date: on or after', hint='Date filter operator: on or after the day.'"`
	FiltersDateIsBefore                        string                 `i18n:"label='Date: before', hint='Date filter operator: before the day.'"`
	FiltersDateIsBeforeOrOn                    string                 `i18n:"label='Date: on or before', hint='Date filter operator: on or before the day.'"`
	FiltersDateDays                            string                 `i18n:"label='Date: days', hint='Unit of the date filter: days.'"`
	FiltersDateMonths                          string                 `i18n:"label='Date: months', hint='Unit of the date filter: months.'"`
	FiltersDateAnd                             string                 `i18n:"label='Date: and', hint='Word between the two days of a range.'"`
	FiltersTo                                  string                 `i18n:"label='To', hint='Word between the two ends of a range.'"`
	FiltersNumberEquals                        string                 `i18n:"label='Number: equals', hint='Number filter operator: equal to.'"`
	FiltersNumberBetween                       string                 `i18n:"label='Number: between', hint='Number filter operator: between two numbers.'"`
	FiltersNumberGreaterThan                   string                 `i18n:"label='Number: greater than', hint='Number filter operator: greater than.'"`
	FiltersNumberLessThan                      string                 `i18n:"label='Number: less than', hint='Number filter operator: less than.'"`
	FiltersNumberAnd                           string                 `i18n:"label='Number: and', hint='Word between the two numbers of a range.'"`
	FiltersStringEquals                        string                 `i18n:"label='Text: equals', hint='Text filter operator: equal to.'"`
	FiltersStringContains                      string                 `i18n:"label='Text: contains', hint='Text filter operator: contains.'"`
	FiltersMultipleSelectIn                    string                 `i18n:"label='Choice: in', hint='Choice filter operator: one of those selected.'"`
	FiltersMultipleSelectNotIn                 string                 `i18n:"label='Choice: not in', hint='Choice filter operator: none of those selected.'"`
	Month                                      string                 `i18n:"hint='Label of a month.'"`
	MonthNames                                 [13]string             `i18n:"hint='The names of the months, January first (the item 0 is unused).'"`
	Year                                       string                 `i18n:"hint='Label of a year.'"`
	PaginationRowsPerPage                      string                 `i18n:"label='Rows per page', hint='Label of how many records a page of a listing shows.'"`
	PaginationPageInfo                         string                 `i18n:"label='Page info', hint='Which records the page shows.', fields=(;'{currPageStart}'='the first record', '{currPageEnd}'='the last record', '{total}'='how many there are')"`
	PaginationPage                             string                 `i18n:"label='Page', hint='Label before the page number.'"`
	PaginationOfPage                           string                 `i18n:"label='Of pages', hint='After the page number.', fields=(;'{total}'='how many pages there are')"`
	ListingNoRecordToShow                      string                 `i18n:"label='No records', hint='Shown when a listing is empty.'"`
	ListingSelectedCountNotice                 string                 `i18n:"label='Selected count', hint='How many records are selected.', fields=(;'{count}'='the number')"`
	ListingClearSelection                      string                 `i18n:"label='Clear selection', hint='Link that unselects every record.'"`
	BulkActionNoAvailableRecords               string                 `i18n:"label='Bulk action: none possible', hint='Shown when no selected record can take the action.'"`
	BulkActionSelectedIdsProcessNoticeTemplate string                 `i18n:"label='Bulk action: some impossible', hint='Shown when some selected records cannot take the action.', fields=(;'{ids}'='those records')"`
	ConfirmDialogPromptTitle                   string                 `i18n:"label='Confirmation title', hint='Title of a confirmation dialog.'"`
	ConfirmDialogPromptText                    string                 `i18n:"label='Confirmation text', hint='Question of a confirmation dialog.'"`
	Language                                   string                 `i18n:"hint='Label of a language.'"`
	Colon                                      string                 `i18n:"hint='The colon put after a label, as the language writes it.'"`
	NotFoundPageNotice                         string                 `i18n:"label='Page not found', hint='Shown when the page asked for does not exist.'"`
	AddRow                                     string                 `i18n:"hint='Button that adds a row to a list.'"`
	ListEditorDeletedItem                      string                 `i18n:"label='List: removed item', hint='Mark of an item removed from a list, before saving.'"`
	ListEditorRevertDeletion                   string                 `i18n:"label='List: undo', hint='Button that brings back an item removed from a list.'"`
	ListEditorRemoveItem                       string                 `i18n:"label='List: remove', hint='Button that removes an item from a list.'"`
	ListEditorActions                          string                 `i18n:"label='List: actions', hint='Title of the actions column of a list.'"`
	PleaseSelectRecord                         string                 `i18n:"hint='Shown when an action needs a record and none was selected.'"`
	PrinterOptions                             PrinterOptionsMessages `i18n:"hint='The words of the print options.'"`
	BulkActionConfirmationTextTemplate         string                 `i18n:"type=html, label='Bulk action: confirmation', hint='Asks to confirm an action on the selected records (HTML).', fields=(;'{Action}'='the action', '{count}'='how many records')"`

	TimeFormats TimeFormatMessages `i18n:"hint='How times and dates are written.'"`

	Common StrMap `i18n:"hint='Labels of the fields every model may have, by the field\\'s name (CreatedAt, Title…): what a field is called when its model says nothing.'"`

	Error               string           `i18n:"hint='Title of an error.'"`
	ErrEmptyParamID     i18n.ErrorString `i18n:"label='Empty ID', hint='Error when a record was asked for without its id.'"`
	ErrPermissionDenied i18n.ErrorString `i18n:"hint='Error when the user may not do what was asked.'"`
	ErrFieldRequired    i18n.ErrorString `i18n:"hint='Error of a required field left empty.'"`

	// Optimistic locking of the edit form (see record_stamp.go).
	// ErrRecordChangedBy takes the author's name and e-mail.
	// ErrRecordChangedUnknownWhen is for a model with no UpdatedAt: the state
	// hash says the record moved, and nothing says when or by whom.
	ErrRecordChanged            string           `i18n:"hint='Error when someone else saved the record after the form was opened.', fields=(;'%s'='when it was saved')"`
	ErrRecordChangedBy          string           `i18n:"hint='Error when someone else saved the record after the form was opened.', fields=(;'%[1]s'='who saved it', '%[2]s'='their account', '%[3]s'='when it was saved')"`
	ErrRecordChangedUnknownWhen i18n.ErrorString `i18n:"hint='Error when someone else saved the record after the form was opened, when not known.'"`
	ErrRecordStampMissing       i18n.ErrorString `i18n:"hint='Error of a form too old to be saved.'"`

	CopiedToClipboard string `i18n:"hint='Shown once something was copied.'"`
}

// FormatDateTime writes an instant the way this language does. Falls back to a
// readable layout when the messages carry none — Format("") returns "".
func (msgr *Messages) FormatDateTime(t time.Time) string {
	layout := msgr.TimeFormats.DateTime
	if layout == "" {
		layout = "2006-01-02 15:04:05 Z07:00"
	}
	return t.Format(layout)
}

// RecordChangedMessage is the stale-record message: WHEN the record was changed
// and, when the application could load it, by WHOM.
func (msgr *Messages) RecordChangedMessage(name, email string, at time.Time) string {
	when := msgr.FormatDateTime(at)
	if name == "" && email == "" {
		return fmt.Sprintf(msgr.ErrRecordChanged, when)
	}
	return fmt.Sprintf(msgr.ErrRecordChangedBy, name, email, when)
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
	YouAreHere:                         "You Are Here",
	SuccessfullyUpdated:                "Successfully Updated",
	SuccessfullyCreated:                "Successfully Created",
	SuccessfullyDeleted:                "Successfully Deleted",
	Search:                             "Search",
	New:                                "New",
	Update:                             "Update",
	Execute:                            "Execute",
	Delete:                             "Delete",
	DeleteRelated:                      "Delete Related",
	DeleteRelatedHint:                  "Delete all related records",
	ShowRelatedItemsTitle:              "Show related items",
	RelatedItemsForDeletionActionTitle: "Related Items For Deletion",
	Edit:                               "Edit",
	FormTitle:                          "Form",
	OK:                                 "OK",
	Cancel:                             "Cancel",
	Clear:                              "Clear",
	Create:                             "Create",
	DeleteConfirmationTextTemplate:     "Are you sure you want to delete {the_model}: {title}?",
	CreatingObjectTitleTemplate:        "New {modelName}",
	CreatingFemaleObjectTitleTemplate:  "New {modelName}",
	EditingTitleTemplate:               "Editing {modelName}",
	EditingObjectTitleTemplate:         "Editing {modelName} {id}",
	ListingObjectTitleTemplate:         "Listing {modelName}",
	DetailingObjectTitleTemplate:       "{modelName} {id}",
	FiltersClear:                       "Clear Filters",
	FiltersAdd:                         "Add Filters",
	FilterApply:                        "Apply",
	FilterByTemplate:                   "Filter by {filter}",
	FiltersDateInTheLast:               "is in the last",
	FiltersDateEquals:                  "is equal to",
	FiltersDateBetween:                 "is between",
	FiltersDateIsAfter:                 "is after",
	FiltersDateIsAfterOrOn:             "is on or after",
	FiltersDateIsBefore:                "is before",
	FiltersDateIsBeforeOrOn:            "is before or on",
	FiltersDateDays:                    "days",
	FiltersDateMonths:                  "months",
	FiltersDateAnd:                     "and",
	FiltersTo:                          "to",
	FiltersNumberEquals:                "is equal to",
	FiltersNumberBetween:               "between",
	FiltersNumberGreaterThan:           "is greater than",
	FiltersNumberLessThan:              "is less than",
	FiltersNumberAnd:                   "and",
	FiltersStringEquals:                "is equal to",
	FiltersStringContains:              "contains",
	FiltersMultipleSelectIn:            "in",
	FiltersMultipleSelectNotIn:         "not in",
	Month:                              "Month",
	MonthNames: [13]string{
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
	ListEditorDeletedItem:                      "Removed item",
	ListEditorRevertDeletion:                   "Undo",
	ListEditorRemoveItem:                       "Remove",
	ListEditorActions:                          "Actions",

	BulkActionConfirmationTextTemplate: "Are you sure you want to <b>{Action}</b> then {count} records?",

	TimeFormats: TimeFormatMessages{
		Time:     "15:04:05 Z07:00",
		Date:     "2006-01-02",
		DateTime: "2006-01-02 15:04:05 Z07:00",
	},

	Error:           "ERROR",
	ErrEmptyParamID: "Empty param ID",
	ErrRecordChanged: "This record was changed by someone else at %s, after you opened this form. " +
		"Reload it and make your changes again, so that nothing that has been saved in the meantime is lost.",
	ErrRecordChangedBy: "This record was changed by %s (%s) at %s, after you opened this form. " +
		"Reload it and make your changes again, so that nothing that has been saved in the meantime is lost.",
	ErrRecordChangedUnknownWhen: "This record was changed by someone else after you opened this form. " +
		"Reload it and make your changes again, so that nothing saved in the meantime is lost.",
	ErrRecordStampMissing: "This form is out of date and cannot be saved. Reload it and make your changes again.",
	ErrFieldRequired:      i18n.ErrorString(ErrFieldRequired.Error()),
	CopiedToClipboard:     "Copied to clipboard",
}

var DefaultMessages = Messages_en_US

var Messages_pt_BR = &Messages{
	YouAreHere:                         "Você está aqui",
	CopiedToClipboard:                  "Copiado para a área de transferência!",
	TheFemaleTitle:                     "A %s",
	TheMaleTitle:                       "O %s",
	SuccessfullyUpdated:                "Atualizado com Sucesso",
	SuccessfullyCreated:                "Cadastrado com Sucesso",
	SuccessfullyDeleted:                "Excluído com Sucesso",
	SuccessfullyExecutedAction:         "Ação executada com Sucesso",
	Search:                             "Pesquisa",
	New:                                "Novo",
	Update:                             "Atualizar",
	Execute:                            "Executar",
	Delete:                             "Excluir",
	DeleteRelated:                      "Excluir items relacionados",
	DeleteRelatedHint:                  "Exclui automaticamente todas entidades relacionados",
	RelatedItemsForDeletionActionTitle: "Itens a serem removidos durante a exclusão",
	ShowRelatedItemsTitle:              "Ver itens relacionados",
	Edit:                               "Editar",
	FormTitle:                          "Formulário",
	OK:                                 "OK",
	Cancel:                             "Cancelar",
	Clear:                              "Limpar",

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
	ListEditorDeletedItem:                      "Item removido",
	ListEditorRevertDeletion:                   "Desfazer",
	ListEditorRemoveItem:                       "Remover",
	ListEditorActions:                          "Ações",
	Error:                                      "Erro",
	Month:                                      "Mês",
	Year:                                       "Ano",
	ErrEmptyParamID:                            "Parâmetro ID não informado",
	ErrPermissionDenied:                        "Permissão negada",
	ErrRecordChanged: "Este registro foi alterado por outra pessoa em %s, depois que você abriu este formulário. " +
		"Recarregue-o e refaça suas alterações, para que nada do que foi salvo nesse meio tempo se perca.",
	ErrRecordChangedBy: "Este registro foi alterado por %s (%s) em %s, depois que você abriu este formulário. " +
		"Recarregue-o e refaça suas alterações, para que nada do que foi salvo nesse meio tempo se perca.",
	ErrRecordChangedUnknownWhen: "Este registro foi alterado por outra pessoa depois que você abriu este formulário. " +
		"Recarregue-o e refaça suas alterações, para que nada do que foi salvo nesse meio tempo se perca.",
	ErrRecordStampMissing: "Este formulário está desatualizado e não pode ser salvo. Recarregue-o e refaça suas alterações.",
	ErrFieldRequired:      "Este campo não pode ser vazio",

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
		"RecordTimes":         "Horas do Registro",
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

	MonthNames: [13]string{
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
