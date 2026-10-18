package gadcode

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// MessagesKey is the module of the words of the editor of gad code
// (vx-gad-code): its buttons and panels — the same keys as its labels prop.
const MessagesKey i18n.ModuleKey = "rvq-admin/gadcode"

// GetMessages are the words of the editor in ctx's language.
func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

// ConfigureMessages registers the words of the editor.
func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModule(language.English, MessagesKey, Messages_en_US).
		RegisterForModule(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

// Labels are the words of the editor in ctx's language, as its labels prop
// (VXGadCodeBuilder.Labels).
func Labels(ctx context.Context) map[string]string {
	m := GetMessages(ctx)
	return map[string]string{
		"format": m.Format, "expand": m.Expand, "collapse": m.Collapse, "structure": m.Structure,
		"doc": m.Doc, "refresh": m.Refresh, "noStructure": m.NoStructure, "noDoc": m.NoDoc,
		"notFormatted": m.NotFormatted, "undo": m.Undo, "redo": m.Redo, "compare": m.Compare,
		"edit": m.Edit, "prevChange": m.PrevChange, "nextChange": m.NextChange, "revert": m.Revert,
	}
}

// Messages are the words of the editor of gad code.
type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Format            string `i18n:"hint='Button that formats the code (gad fmt).'"`
	Expand            string `i18n:"hint='Button that expands the editor over the page.'"`
	Collapse          string `i18n:"hint='Button that brings the editor expanded back.'"`
	Structure         string `i18n:"hint='Title of the panel with the declarations of the code, as a tree.'"`
	Doc               string `i18n:"hint='Title of the panel with the documentation of the code (gad doc).'"`
	Refresh           string `i18n:"hint='Button that reads the structure or the documentation again.'"`
	NoStructure       string `i18n:"hint='Shown when the code declares nothing.'"`
	NoDoc             string `i18n:"hint='Shown when the code has no documentation.'"`
	NotFormatted      string `i18n:"hint='Shown when the code could not be formatted (it does not parse).'"`
	Undo              string `i18n:"hint='Button that undoes the last edit.'"`
	Redo              string `i18n:"hint='Button that redoes the edit undone.'"`
	Compare           string `i18n:"hint='Button that compares the code with its value when the page opened.'"`
	Edit              string `i18n:"hint='Button that leaves the compare, back to the editor.'"`
	PrevChange        string `i18n:"hint='Button that goes to the previous change of the compare.'"`
	NextChange        string `i18n:"hint='Button that goes to the next change of the compare.'"`
	Revert            string `i18n:"hint='Button of a change of the compare: takes the old part back.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "The editor of gad code of the admin (a form structure, a banner body, the texts of gadx and gadt).",
	Format:            "Format",
	Expand:            "Expand",
	Collapse:          "Collapse",
	Structure:         "Structure",
	Doc:               "Documentation",
	Refresh:           "Refresh",
	NoStructure:       "No declarations.",
	NoDoc:             "No documentation.",
	NotFormatted:      "The code could not be formatted: fix what is wrong in it.",
	Undo:              "Undo",
	Redo:              "Redo",
	Compare:           "Compare with the initial value",
	Edit:              "Back to the editor",
	PrevChange:        "Previous change",
	NextChange:        "Next change",
	Revert:            "Revert",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "O editor de código gad do admin (a estrutura de um formulário, o corpo de um banner, os textos gadx e gadt).",
	Format:            "Formatar",
	Expand:            "Expandir",
	Collapse:          "Recolher",
	Structure:         "Estrutura",
	Doc:               "Documentação",
	Refresh:           "Atualizar",
	NoStructure:       "Nenhuma declaração.",
	NoDoc:             "Nenhuma documentação.",
	NotFormatted:      "O código não pôde ser formatado: corrija o que está errado nele.",
	Undo:              "Desfazer",
	Redo:              "Refazer",
	Compare:           "Comparar com o valor inicial",
	Edit:              "Voltar ao editor",
	PrevChange:        "Mudança anterior",
	NextChange:        "Próxima mudança",
	Revert:            "Reverter",
}
