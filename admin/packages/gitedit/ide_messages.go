package gitedit

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// IdeMessagesKey is the module of the words of the gad IDE's Changes and Git
// panels (and of the diff browser in them), as the IDE is given them
// (GadIde's messages: IdeMessages of @gad-lang/ide-vuetify — the same keys).
const IdeMessagesKey i18n.ModuleKey = "rvq-admin/gitedit-ide"

// GetIdeMessages are the words of the IDE's panels in ctx's language.
func GetIdeMessages(ctx context.Context) *IdeMessages {
	return i18n.MustGetModuleMessages(ctx, IdeMessagesKey, IdeMessages_en_US).(*IdeMessages)
}

// ConfigureIdeMessages registers the words of the IDE's panels.
func ConfigureIdeMessages(b *i18n.Builder) {
	b.RegisterForModule(language.English, IdeMessagesKey, IdeMessages_en_US).
		RegisterForModule(language.BrazilianPortuguese, IdeMessagesKey, IdeMessages_pt_BR)
}

// IdeMessages are the texts of the IDE's Changes and Git panels; their JSON
// is the IDE's messages prop.
type IdeMessages struct {
	ModuleDescription   string `json:"-" i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Changes             string `json:"changes" i18n:"hint='Title of the panel of the IDE with the files changed and not committed (its tab and its button).'"`
	Git                 string `json:"git" i18n:"hint='Title of the panel of the IDE with the branches and their commits.'"`
	Refresh             string `json:"refresh" i18n:"hint='Title of the button that reads the changes, or the branches, again.'"`
	Files               string `json:"files" i18n:"hint='Title of the tree of the files of a diff.'"`
	Branches            string `json:"branches" i18n:"hint='Title of the list of the local branches.'"`
	RemoteBranches      string `json:"remoteBranches" i18n:"hint='Title of the list of the branches of the remotes.'"`
	More                string `json:"more" i18n:"hint='Button at the end of the log that loads more commits.'"`
	ChooseCommit        string `json:"chooseCommit" i18n:"hint='Shown in the Git panel with no commit chosen.'"`
	Patch               string `json:"patch" i18n:"hint='Button that downloads the patch of the commit chosen.'"`
	DownloadCommitPatch string `json:"downloadCommitPatch" i18n:"hint='Title of the button that downloads the patch of the commit chosen.'"`
	DownloadFile        string `json:"downloadFile" i18n:"hint='Title of the button of the diff of a file that downloads it as the commit has it.'"`
	DownloadPatch       string `json:"downloadPatch" i18n:"hint='Title of the button of the diff of a file that downloads its patch from the commit before.'"`
	NoParent            string `json:"noParent" i18n:"hint='Title of the old side of a file of the first commit (none before it).'"`
	Expand              string `json:"expand" i18n:"hint='Title of the button that expands the panel over the others.'"`
	Restore             string `json:"restore" i18n:"hint='Title of the button that brings the panel expanded back.'"`
	Old                 string `json:"old" i18n:"hint='Title of the side of a file as it was.'"`
	Current             string `json:"current" i18n:"hint='Title of the side of a file as it is.'"`
	Binary              string `json:"binary" i18n:"hint='Shown for a binary file changed: it is not compared.'"`
	NoChanges           string `json:"noChanges" i18n:"hint='Shown when no file was changed.'"`
	RenamedTo           string `json:"renamedTo" i18n:"hint='Precedes the new path of a file renamed or moved.'"`
	Unchanged           string `json:"unchanged" i18n:"hint='Said of a file renamed or moved whose content did not change.'"`
	Loading             string `json:"loading" i18n:"hint='Shown in the tab of a file while its diff is loaded.'"`
	LoadError           string `json:"loadError" i18n:"hint='Shown in the tab of a file whose diff could not be loaded.'"`
	Undo                string `json:"undo" i18n:"hint='Button of the diff of a file: undoes the last edit of its current side.'"`
	Redo                string `json:"redo" i18n:"hint='Button of the diff of a file: redoes the edit undone.'"`
	Save                string `json:"save" i18n:"hint='Button of the diff of a file: saves its current side as edited.'"`
	Saving              string `json:"saving" i18n:"hint='Said while the file edited in the diff is saved.'"`
	Saved               string `json:"saved" i18n:"hint='Said once the file edited in the diff was saved.'"`
	Unsaved             string `json:"unsaved" i18n:"hint='Said of a file edited in the diff and not saved.'"`
	Revert              string `json:"revert" i18n:"hint='Title of the button of a change: takes the old part into the current.'"`
	Prev                string `json:"prev" i18n:"hint='Title of the button that goes to the previous change of the diff (Shift+F7).'"`
	Next                string `json:"next" i18n:"hint='Title of the button that goes to the next change of the diff (F7).'"`
}

var IdeMessages_en_US = &IdeMessages{
	ModuleDescription:   "The Changes and Git panels of the editor of the files (the IDE).",
	Changes:             "Changes",
	Git:                 "Git",
	Refresh:             "Refresh",
	Files:               "Files",
	Branches:            "Branches",
	RemoteBranches:      "Remote",
	More:                "More",
	ChooseCommit:        "Choose a commit.",
	Patch:               "Patch",
	DownloadCommitPatch: "Download the patch of the commit",
	DownloadFile:        "Download the file (this commit's)",
	DownloadPatch:       "Download the patch (from the commit before)",
	NoParent:            "(none)",
	Expand:              "Expand",
	Restore:             "Restore",
	Old:                 "Old",
	Current:             "Current",
	Binary:              "A binary file: it is not compared.",
	NoChanges:           "No changes.",
	RenamedTo:           "Renamed to",
	Unchanged:           "Its content did not change.",
	Loading:             "Loading…",
	LoadError:           "The diff could not be loaded.",
	Undo:                "Undo",
	Redo:                "Redo",
	Save:                "Save",
	Saving:              "Saving…",
	Saved:               "Saved",
	Unsaved:             "Not saved",
	Revert:              "Revert: the old part into the current",
	Prev:                "Previous change (Shift+F7)",
	Next:                "Next change (F7)",
}

var IdeMessages_pt_BR = &IdeMessages{
	ModuleDescription:   "Os painéis Alterações e Git do editor dos arquivos (a IDE).",
	Changes:             "Alterações",
	Git:                 "Git",
	Refresh:             "Atualizar",
	Files:               "Arquivos",
	Branches:            "Branches",
	RemoteBranches:      "Remotas",
	More:                "Mais",
	ChooseCommit:        "Escolha um commit.",
	Patch:               "Patch",
	DownloadCommitPatch: "Baixar o patch do commit",
	DownloadFile:        "Baixar o arquivo (deste commit)",
	DownloadPatch:       "Baixar o patch (desde o commit anterior)",
	NoParent:            "(nenhum)",
	Expand:              "Expandir",
	Restore:             "Restaurar",
	Old:                 "Antigo",
	Current:             "Atual",
	Binary:              "Um arquivo binário: não é comparado.",
	NoChanges:           "Nenhuma alteração.",
	RenamedTo:           "Renomeado para",
	Unchanged:           "O conteúdo não mudou.",
	Loading:             "Carregando…",
	LoadError:           "Não foi possível carregar o diff.",
	Undo:                "Desfazer",
	Redo:                "Refazer",
	Save:                "Salvar",
	Saving:              "Salvando…",
	Saved:               "Salvo",
	Unsaved:             "Não salvo",
	Revert:              "Desfazer: a parte antiga no atual",
	Prev:                "Alteração anterior (Shift+F7)",
	Next:                "Próxima alteração (F7)",
}
