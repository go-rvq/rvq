package gitedit

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// MessagesKey is the module of the words of the editor of the files.
const MessagesKey i18n.ModuleKey = "rvq-admin/gitedit"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

// ConfigureMessages registers the words of the editor of the files.
func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModule(language.English, MessagesKey, Messages_en_US).
		RegisterForModule(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Title             string `i18n:"hint='Title of the page (and its menu item) that edits the files of the site.'"`
	Draft             string `i18n:"hint='Title of the box that says where the draft of the user stands.'"`
	DraftHint         string `i18n:"hint='Explains the draft: the changes are the user own, the site changes only when published.'"`
	Ahead             string `i18n:"hint='How many commits of the draft are not published.', fields=(;'%d'='the count')"`
	Behind            string `i18n:"hint='How many commits published by others the draft does not have.', fields=(;'%d'='the count')"`
	UpToDate          string `i18n:"hint='Said of a draft with nothing to publish nor to update.'"`
	Changes           string `i18n:"hint='Title of the list of the files changed and not committed.'"`
	NoChanges         string `i18n:"hint='Shown when no file was changed since the last commit.'"`
	History           string `i18n:"hint='Title of the list of the last commits.'"`
	Files             string `i18n:"hint='Title of the editor of the files (the IDE).'"`
	Preview           string `i18n:"hint='Button that opens the site as the draft makes it.'"`
	Help              string `i18n:"hint='Button that opens the documentation of the editor of the files.'"`
	OpenIDE           string `i18n:"hint='Button that opens the editor of the files (the IDE) in a tab of its own.'"`

	CommitAction      string `i18n:"label='Commit', hint='Action that records the changes of the draft, with a message.'"`
	UpdateAction      string `i18n:"label='Update', hint='Action that brings into the draft what others published.'"`
	PublishAction     string `i18n:"label='Publish', hint='Action that puts the commits of the draft on the site.'"`
	DiscardAction     string `i18n:"label='Discard', hint='Action that drops the changes of a file not committed.'"`
	ResetAction       string `i18n:"label='Start over', hint='Action that drops the draft: a new one is made from the site.'"`
	CommitFormMessage string `i18n:"label='Message', hint='Label of the message of a commit.'"`

	// the permissions of the page's actions: titles and descriptions (the
	// tree of the roles)
	CommitAction_Desc  string `i18n:"hint='Description of the permission of committing.'"`
	UpdateAction_Desc  string `i18n:"hint='Description of the permission of updating the draft.'"`
	PublishAction_Desc string `i18n:"hint='Description of the permission of publishing.'"`
	DiscardAction_Desc string `i18n:"hint='Description of the permission of discarding changes.'"`
	ResetAction_Desc   string `i18n:"hint='Description of the permission of starting over.'"`
	PreviewAction      string `i18n:"label='Preview', hint='Title of the permission of seeing the site of the draft.'"`
	PreviewAction_Desc string `i18n:"hint='Description of the permission of seeing the site of the draft.'"`
	CreateAction       string `i18n:"label='Create files', hint='Title of the permission of creating files and folders.'"`
	CreateAction_Desc  string `i18n:"hint='Description of the permission of creating files and folders.'"`
	EditAction         string `i18n:"label='Edit files', hint='Title of the permission of changing files that are there.'"`
	EditAction_Desc    string `i18n:"hint='Description of the permission of changing files that are there.'"`
	RenameAction       string `i18n:"label='Rename files', hint='Title of the permission of renaming in the folder.'"`
	RenameAction_Desc  string `i18n:"hint='Description of the permission of renaming in the folder.'"`
	MoveAction         string `i18n:"label='Move files', hint='Title of the permission of moving to another folder.'"`
	MoveAction_Desc    string `i18n:"hint='Description of the permission of moving to another folder.'"`
	DeleteAction       string `i18n:"label='Delete files', hint='Title of the permission of deleting files and folders.'"`
	DeleteAction_Desc  string `i18n:"hint='Description of the permission of deleting files and folders.'"`
	ImportAction       string `i18n:"label='Import files', hint='Title of the permission of uploading and downloading from the internet.'"`
	ImportAction_Desc  string `i18n:"hint='Description of the permission of uploading and downloading from the internet.'"`
	GitAction          string `i18n:"label='Access by git', hint='Title of the permission of cloning and pushing the files by git.'"`
	ByGit              string `i18n:"hint='Title of the box of the URLs of the files by git.'"`
	ByGitHint          string `i18n:"hint='Explains the URLs by git: cloned and pushed with the login and password.'"`
	GitDraftURL        string `i18n:"hint='Label of the URL by git of the draft of the user.'"`
	GitSiteURL         string `i18n:"hint='Label of the URL by git of the site: a push to it publishes.'"`
	GitAction_Desc     string `i18n:"hint='Description of the permission of cloning and pushing the files by git.'"`

	Committed      string `i18n:"hint='Shown once the changes were committed.', fields=(;'%s'='the short hash')"`
	Updated        string `i18n:"hint='Shown once the draft was updated.'"`
	Published      string `i18n:"hint='Shown once the draft was published: the site changed.'"`
	Discarded      string `i18n:"hint='Shown once the changes of a file were dropped.', fields=(;'%s'='the file')"`
	Reset          string `i18n:"hint='Shown once the draft was dropped.'"`
	ErrUncommitted string `i18n:"hint='Error: there are changes not committed, commit or discard them before publishing.'"`
	ErrBehind      string `i18n:"hint='Error: others published since; update the draft before publishing.'"`
	ErrSiteChanged string `i18n:"hint='Error: files of the site were changed out of the editor; nothing was overwritten.'"`
	ErrNothing     string `i18n:"hint='Error: nothing to commit.'"`
	ErrInvalid     string `i18n:"hint='Error: the files of the draft do not work; the commit or the publishing was not made.', fields=(;'%s'='what is wrong')"`
}

var Messages_en_US = &Messages{
	ModuleDescription:  "The editor of the files of the site: drafts, commits and publishing.",
	Title:              "Site files",
	Draft:              "Your draft",
	DraftHint:          "You edit a draft of your own: the site changes only when you publish.",
	Ahead:              "%d commit(s) to publish",
	Behind:             "%d commit(s) published by others: update",
	UpToDate:           "Up to date with the site",
	Changes:            "Changes not committed",
	NoChanges:          "Nothing changed since the last commit.",
	History:            "History",
	Files:              "Files",
	Preview:            "Preview the site",
	Help:               "Help",
	OpenIDE:            "Open the editor",
	CommitAction:       "Commit",
	UpdateAction:       "Update",
	PublishAction:      "Publish",
	DiscardAction:      "Discard",
	ResetAction:        "Start over",
	CommitFormMessage:  "Message",
	CommitAction_Desc:  "Records the changes of the draft, with a message.",
	UpdateAction_Desc:  "Brings into the draft what others published.",
	PublishAction_Desc: "Puts the commits of the draft on the site.",
	DiscardAction_Desc: "Drops the changes of a file not committed.",
	ResetAction_Desc:   "Drops the whole draft: a new one is made from the site.",
	PreviewAction:      "Preview",
	PreviewAction_Desc: "Sees the site as the draft makes it.",
	CreateAction:       "Create files",
	CreateAction_Desc:  "Creates files and folders that are not there — in the editor, by uploading or importing, by the WebDAV.",
	EditAction:         "Edit files",
	EditAction_Desc:    "Changes files that are there; without it, the editor opens them read only.",
	RenameAction:       "Rename files",
	RenameAction_Desc:  "Renames a file or a folder in its folder.",
	MoveAction:         "Move files",
	MoveAction_Desc:    "Moves a file or a folder to another folder.",
	DeleteAction:       "Delete files",
	DeleteAction_Desc:  "Deletes files and folders.",
	ImportAction:       "Import files",
	ImportAction_Desc:  "Uploads files from the computer and downloads them from the internet (with Create files or Edit files for what it writes).",
	GitAction:          "Access by git",
	ByGit:              "By git",
	ByGitHint:          "Clone and push with git, with your login and password.",
	GitDraftURL:        "Your draft (a push changes it; publish here)",
	GitSiteURL:         "The site (a push publishes)",
	GitAction_Desc:     "Clones and pushes the files by git — the draft, and the site with Publish —, each file changed asking what the editor asks.",
	Committed:          "Committed: %s",
	Updated:            "The draft has what was published.",
	Published:          "Published: the site has your commits.",
	Discarded:          "Changes of %s discarded.",
	Reset:              "The draft was dropped.",
	ErrUncommitted:     "Commit or discard the changes before publishing.",
	ErrBehind:          "Others published since: update the draft before publishing.",
	ErrSiteChanged:     "Files of the site were changed out of the editor: nothing was overwritten.",
	ErrNothing:         "Nothing to commit.",
	ErrInvalid:         "The files of the draft do not work: %s",
}

var Messages_pt_BR = &Messages{
	ModuleDescription:  "O editor dos arquivos do site: rascunhos, commits e publicação.",
	Title:              "Arquivos do site",
	Draft:              "Seu rascunho",
	DraftHint:          "Você edita um rascunho só seu: o site muda só quando você publica.",
	Ahead:              "%d commit(s) para publicar",
	Behind:             "%d commit(s) publicado(s) por outros: atualize",
	UpToDate:           "Igual ao site",
	Changes:            "Alterações sem commit",
	NoChanges:          "Nada mudou desde o último commit.",
	History:            "Histórico",
	Files:              "Arquivos",
	Preview:            "Ver o site do rascunho",
	Help:               "Ajuda",
	OpenIDE:            "Abrir o editor",
	CommitAction:       "Commit",
	UpdateAction:       "Atualizar",
	PublishAction:      "Publicar",
	DiscardAction:      "Descartar",
	ResetAction:        "Recomeçar",
	CommitFormMessage:  "Mensagem",
	CommitAction_Desc:  "Registra as alterações do rascunho, com uma mensagem.",
	UpdateAction_Desc:  "Traz para o rascunho o que outros publicaram.",
	PublishAction_Desc: "Põe os commits do rascunho no site.",
	DiscardAction_Desc: "Descarta as alterações sem commit de um arquivo.",
	ResetAction_Desc:   "Descarta o rascunho inteiro: um novo é feito a partir do site.",
	PreviewAction:      "Pré-visualizar",
	PreviewAction_Desc: "Vê o site como o rascunho o faz.",
	CreateAction:       "Criar arquivos",
	CreateAction_Desc:  "Cria arquivos e pastas que não existem — no editor, por upload ou importação, pelo WebDAV.",
	EditAction:         "Editar arquivos",
	EditAction_Desc:    "Muda arquivos que existem; sem ela, o editor os abre só para leitura.",
	RenameAction:       "Renomear arquivos",
	RenameAction_Desc:  "Renomeia um arquivo ou uma pasta na sua pasta.",
	MoveAction:         "Mover arquivos",
	MoveAction_Desc:    "Move um arquivo ou uma pasta para outra pasta.",
	DeleteAction:       "Excluir arquivos",
	DeleteAction_Desc:  "Exclui arquivos e pastas.",
	ImportAction:       "Importar arquivos",
	ImportAction_Desc:  "Faz upload de arquivos do computador e os baixa da internet (com Criar arquivos ou Editar arquivos para o que grava).",
	GitAction:          "Acessar por git",
	ByGit:              "Pelo git",
	ByGitHint:          "Clone e faça push com o git, com o seu login e senha.",
	GitDraftURL:        "Seu rascunho (um push o muda; publique aqui)",
	GitSiteURL:         "O site (um push publica)",
	GitAction_Desc:     "Clona e envia (push) os arquivos por git — o rascunho, e o site com Publicar —, cada arquivo alterado pedindo o que o editor pede.",
	Committed:          "Commit feito: %s",
	Updated:            "O rascunho tem o que foi publicado.",
	Published:          "Publicado: o site tem os seus commits.",
	Discarded:          "Alterações de %s descartadas.",
	Reset:              "O rascunho foi descartado.",
	ErrUncommitted:     "Faça o commit ou descarte as alterações antes de publicar.",
	ErrBehind:          "Outros publicaram depois: atualize o rascunho antes de publicar.",
	ErrSiteChanged:     "Arquivos do site foram alterados fora do editor: nada foi sobrescrito.",
	ErrNothing:         "Nada para o commit.",
	ErrInvalid:         "Os arquivos do rascunho não funcionam: %s",
}
