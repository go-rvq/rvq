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

type Messages struct {
	History         string
	HistoryEmpty    string
	Hash            string
	Published       string
	Tag             string
	Accesses        string
	InitialRevision string
	Unchanged       string
	Old             string
	New             string

	Revision          string
	Author            string
	When              string
	Status            string
	Compare           string
	CompareCurrent    string
	SelectTwoHint     string
	PickToCompare     string
	Merged            string
	Revert            string
	Reverted          string
	Field             string
	RevertSelected    string
	PartialRevert     string
	PartialRevertHint string
	Current           string
	NoHunks           string
	NoChanges         string
	Invert            string
}

var (
	Messages_en_US = &Messages{
		History:           "History",
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
		Author:            "Author",
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
		RevertSelected:    "Revert selected",
		PartialRevert:     "Partial revert",
		PartialRevertHint: "Click the highlighted changes to select what to revert; the selected regions are restored to this revision.",
		Current:           "Current",
		NoHunks:           "This field is identical to the current value.",
		NoChanges:         "No differences between the selected revisions.",
		Invert:            "Invert comparison",
	}

	Messages_pt_BR = &Messages{
		History:           "Histórico",
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
		Author:            "Autor",
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
		RevertSelected:    "Reverter selecionados",
		PartialRevert:     "Reversão parcial",
		PartialRevertHint: "Clique nas alterações destacadas para escolher o que reverter; as regiões marcadas voltam ao valor desta revisão.",
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
