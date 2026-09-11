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
	Published       string
	Tag             string
	Accesses        string
	InitialRevision string
	Unchanged       string
	Old             string
	New             string

	Revision       string
	Author         string
	When           string
	Status         string
	Compare        string
	CompareCurrent string
	SelectTwoHint  string
	PickToCompare  string
	Merged         string
	Revert         string
	Reverted       string
	Field          string
	RevertSelected string
	PartialRevert  string
	NoHunks        string
	NoChanges      string
}

var (
	Messages_en_US = &Messages{
		History:         "History",
		HistoryEmpty:    "No revisions yet.",
		Published:       "Published",
		Tag:             "Tag",
		Accesses:        "Accesses",
		InitialRevision: "Initial revision",
		Unchanged:       "unchanged",
		Old:             "Before",
		New:             "After",
		Revision:        "Revision",
		Author:          "Author",
		When:            "When",
		Status:          "Status",
		Compare:         "Compare selected",
		CompareCurrent:  "Compare with current",
		SelectTwoHint:   "Select two revisions to compare, or use “Compare with current”.",
		PickToCompare:   "Pick revisions above to see the differences here.",
		Merged:          "Merged",
		Revert:          "Revert to this revision",
		Reverted:        "Reverted to the selected revision",
		Field:           "Field",
		RevertSelected:  "Revert selected hunks",
		PartialRevert:   "Partial revert",
		NoHunks:         "This field is identical to the current value.",
		NoChanges:       "No differences between the selected revisions.",
	}

	Messages_pt_BR = &Messages{
		History:         "Histórico",
		HistoryEmpty:    "Ainda não há revisões.",
		Published:       "Publicada",
		Tag:             "Tag",
		Accesses:        "Acessos",
		InitialRevision: "Revisão inicial",
		Unchanged:       "sem alteração",
		Old:             "Antes",
		New:             "Depois",
		Revision:        "Revisão",
		Author:          "Autor",
		When:            "Quando",
		Status:          "Situação",
		Compare:         "Comparar selecionadas",
		CompareCurrent:  "Comparar com a atual",
		SelectTwoHint:   "Selecione duas revisões para comparar, ou use “Comparar com a atual”.",
		PickToCompare:   "Escolha revisões acima para ver as diferenças aqui.",
		Merged:          "Mesclado",
		Revert:          "Reverter para esta revisão",
		Reverted:        "Revertido para a revisão selecionada",
		Field:           "Campo",
		RevertSelected:  "Reverter trechos selecionados",
		PartialRevert:   "Reversão parcial",
		NoHunks:         "Este campo é idêntico ao valor atual.",
		NoChanges:       "Sem diferenças entre as revisões selecionadas.",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
