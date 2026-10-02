package admin

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "admin/packages/history"

func getMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

// GetMessages exposes this package's localized messages so an application can
// reuse the same labels (e.g. link to a model's revisions from elsewhere).
func GetMessages(ctx context.Context) *Messages { return getMessages(ctx) }

type Messages struct {
	History         string `i18n:"hint='Title of a record\\'s history (its revisions).'"`
	Revisions       string `i18n:"hint='Name of the revisions model in the plural (tab, listing title).'"`
	HistoryEmpty    string `i18n:"hint='Shown when a record has no revision yet.'"`
	Hash            string `i18n:"hint='Column with the short identifier of a revision.'"`
	Published       string `i18n:"hint='Column with whether a revision was published.'"`
	Tag             string `i18n:"hint='Column with the name given to a published revision.'"`
	Accesses        string `i18n:"hint='Column with how many times a published revision was accessed.'"`
	InitialRevision string `i18n:"hint='Shown for the first revision of a record.'"`
	Unchanged       string `i18n:"hint='Shown for a field the revision did not change.'"`
	Old             string `i18n:"label='Before', hint='Column with a field\\'s value before a change.'"`
	New             string `i18n:"label='After', hint='Column with a field\\'s value after a change.'"`

	Revision string `i18n:"hint='Name of the revisions model in the singular (detail title).'"`
	// RevisionSnapshot is the detail's field of the record as of the revision.
	RevisionSnapshot string `i18n:"label='Snapshot', hint='Label of the record as of a revision.'"`
	Author           string `i18n:"hint='Column with who made a revision.'"`
	Origin           string `i18n:"hint='Column with where a revision was made from, with the action that shows it on a map.'"`
	// EventDeleted marks the revision of a record's deletion.
	EventDeleted      string `i18n:"label='Event: deleted', hint='Mark of the revision of a record\\'s deletion.'"`
	OriginIP          string `i18n:"label='Origin: address', hint='Label of the address a revision was made from.'"`
	OriginBrowser     string `i18n:"label='Origin: browser', hint='Label of the browser a revision was made from.'"`
	OriginPlace       string `i18n:"label='Origin: place', hint='Label of the place a revision was made from.'"`
	OriginCoordinates string `i18n:"label='Origin: coordinates', hint='Label of the coordinates of that place.'"`
	OriginNoMap       string `i18n:"label='Origin: no map', hint='Shown when the place of the address is not known, so there is no map.'"`
	When              string `i18n:"hint='Column with when a revision was made.'"`
	Status            string `i18n:"hint='Column with the state of a revision.'"`
	Compare           string `i18n:"label='Compare selected', hint='Button that compares the two revisions selected.'"`
	CompareCurrent    string `i18n:"hint='Button that compares a revision with the record as it is.'"`
	SelectTwoHint     string `i18n:"label='Select two: hint', hint='Hint of how to compare revisions.'"`
	PickToCompare     string `i18n:"hint='Shown before revisions were picked to compare.'"`
	Merged            string `i18n:"hint='Mark of a revision that merged others.'"`
	Revert            string `i18n:"label='Revert to this revision', hint='Action that brings the record back to a revision.'"`
	Reverted          string `i18n:"hint='Shown once the record was reverted.'"`
	Field             string `i18n:"hint='Column with the field of a change.'"`
	Changes           string `i18n:"hint='Column with what a revision changed.'"`
	RevertSelected    string `i18n:"hint='Button that reverts only the changes selected.'"`
	PartialRevert     string `i18n:"hint='Title of the revert of only some changes of a field.'"`
	PartialRevertHint string `i18n:"label='Partial revert: hint', hint='How to choose the changes of a partial revert.'"`
	RevertConfirmHint string `i18n:"label='Revert: confirmation', hint='Title of the list of changes a revert will undo.'"`
	ToggleAll         string `i18n:"label='Toggle all', hint='Button that selects or clears every change.'"`
	Current           string `i18n:"hint='Mark of the record as it is now, in a comparison.'"`
	NoHunks           string `i18n:"label='No hunks', hint='Shown when a field is the same as now.'"`
	NoChanges         string `i18n:"hint='Shown when the revisions compared do not differ.'"`
	Invert            string `i18n:"label='Invert comparison', hint='Button that swaps the two sides of a comparison.'"`
}

var (
	Messages_en_US = &Messages{
		History:           "History",
		Revisions:         "Revisions",
		HistoryEmpty:      "No revisions yet.",
		Hash:              "Hash",
		Published:         "Published",
		Tag:               "Tag",
		Accesses:          "Accesses",
		InitialRevision:   "Initial revision",
		Unchanged:         "unchanged",
		Old:               "Before",
		New:               "After",
		Revision:          "Revision",
		RevisionSnapshot:  "Snapshot",
		Author:            "Author",
		Origin:            "Origin",
		EventDeleted:      "Deleted",
		OriginIP:          "IP address",
		OriginBrowser:     "Browser",
		OriginPlace:       "Place",
		OriginCoordinates: "Coordinates",
		OriginNoMap:       "Where the address is is not known: no map.",
		When:              "When",
		Status:            "Status",
		Compare:           "Compare selected",
		CompareCurrent:    "Compare with current",
		SelectTwoHint:     "Select two revisions to compare, or use “Compare with current”.",
		PickToCompare:     "Pick revisions above to see the differences here.",
		Merged:            "Merged",
		Revert:            "Revert to this revision",
		Reverted:          "Reverted to the selected revision",
		Field:             "Field",
		Changes:           "Changes",
		RevertSelected:    "Revert selected",
		PartialRevert:     "Partial revert",
		PartialRevertHint: "Click the highlighted changes to select what to revert; the selected regions are restored to this revision.",
		RevertConfirmHint: "These changes will be reverted:",
		ToggleAll:         "Select / clear all changes",
		Current:           "Current",
		NoHunks:           "This field is identical to the current value.",
		NoChanges:         "No differences between the selected revisions.",
		Invert:            "Invert comparison",
	}

	Messages_pt_BR = &Messages{
		History:           "Histórico",
		Revisions:         "Revisões",
		HistoryEmpty:      "Ainda não há revisões.",
		Hash:              "Hash",
		Published:         "Publicada",
		Tag:               "Tag",
		Accesses:          "Acessos",
		InitialRevision:   "Revisão inicial",
		Unchanged:         "sem alteração",
		Old:               "Antes",
		New:               "Depois",
		Revision:          "Revisão",
		RevisionSnapshot:  "Registro nesta revisão",
		Author:            "Autor",
		Origin:            "Origem",
		EventDeleted:      "Exclusão",
		OriginIP:          "Endereço IP",
		OriginBrowser:     "Navegador",
		OriginPlace:       "Local",
		OriginCoordinates: "Coordenadas",
		OriginNoMap:       "Não se sabe onde fica o endereço: sem mapa.",
		When:              "Quando",
		Status:            "Situação",
		Compare:           "Comparar selecionadas",
		CompareCurrent:    "Comparar com a atual",
		SelectTwoHint:     "Selecione duas revisões para comparar, ou use “Comparar com a atual”.",
		PickToCompare:     "Escolha revisões acima para ver as diferenças aqui.",
		Merged:            "Mesclado",
		Revert:            "Reverter para esta revisão",
		Reverted:          "Revertido para a revisão selecionada",
		Field:             "Campo",
		Changes:           "Alterações",
		RevertSelected:    "Reverter selecionados",
		PartialRevert:     "Reversão parcial",
		PartialRevertHint: "Clique nas alterações destacadas para escolher o que reverter; as regiões marcadas voltam ao valor desta revisão.",
		RevertConfirmHint: "Estas alterações serão revertidas:",
		ToggleAll:         "Marcar / desmarcar todas as alterações",
		Current:           "Atual",
		NoHunks:           "Este campo é idêntico ao valor atual.",
		NoChanges:         "Sem diferenças entre as revisões selecionadas.",
		Invert:            "Inverter comparação",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
