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
	BackToAdmin       string `i18n:"hint='Button of the header of the editor (opened in a tab of its own) that goes back to the admin.'"`
	SignOut           string `i18n:"hint='Button of the header of the editor that signs the user out.'"`
	ThemeLight        string `i18n:"hint='Title of the button that turns the editor light.'"`
	ThemeDark         string `i18n:"hint='Title of the button that turns the editor dark.'"`

	CommitAction      string `i18n:"label='Commit', hint='Action that records the changes of the draft, with a message.'"`
	UpdateAction      string `i18n:"label='Update', hint='Action that brings into the draft what others published.'"`
	PublishAction     string `i18n:"label='Publish', hint='Action that puts the commits of the draft on the site.'"`
	DiscardAction     string `i18n:"label='Discard', hint='Action that drops the changes of a file not committed.'"`
	CopyDiff          string `i18n:"label='Copy', hint='Button of the diff of a file that copies it.'"`
	DiffCopied        string `i18n:"label='Copied!', hint='Said by the button of a diff once it copied it.'"`
	CopyDiffError     string `i18n:"label='Press Ctrl+C to copy', hint='Said by the button of a diff when it could not copy it: how to copy it by hand.'"`
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
	GitAction_Desc     string `i18n:"hint='Description of the permission of cloning and pushing the files by git.'"`

	Committed         string `i18n:"hint='Shown once the changes were committed.', fields=(;'%s'='the short hash')"`
	Updated           string `i18n:"hint='Shown once the draft was updated.'"`
	Published         string `i18n:"hint='Shown once the draft was published: the site changed.'"`
	Discarded         string `i18n:"hint='Shown once the changes of a file were dropped.', fields=(;'%s'='the file')"`
	Reset             string `i18n:"hint='Shown once the draft was dropped.'"`
	ErrUncommitted    string `i18n:"hint='Error: there are changes not committed, commit or discard them before publishing.'"`
	ErrBehind         string `i18n:"hint='Error: others published since; update the draft before publishing.'"`
	ErrSiteChanged    string `i18n:"hint='Error: files of the site were changed out of the editor; nothing was overwritten.'"`
	ErrNothing        string `i18n:"hint='Error: nothing to commit.'"`
	ErrInvalid        string `i18n:"hint='Error: the files of the draft do not work; the commit or the publishing was not made.', fields=(;'%s'='what is wrong')"`
	ShareAction       string `i18n:"hint='Action that shares the draft of the user with other users.'"`
	ShareAction_Desc  string `i18n:"hint='What sharing the draft allows.'"`
	DraftsAction      string `i18n:"hint='Permission to see and reach every draft.'"`
	DraftsAction_Desc string `i18n:"hint='What seeing every draft allows.'"`
	RevokeAction      string `i18n:"hint='Action that revokes a sharing of a draft.'"`

	PublicPreview            string `i18n:"hint='Title of the public preview of the draft (seen with no login).'"`
	PublicPreviewAction      string `i18n:"hint='Action (and permission) that makes the preview of the draft public.'"`
	PublicPreviewAction_Desc string `i18n:"hint='What making the preview public allows.'"`
	PublicPreviewOffAction   string `i18n:"hint='Action that ends the public preview.'"`
	PublicPreviewNone        string `i18n:"hint='Shown when the preview of the draft is not public.'"`
	PublicPreviewOn          string `i18n:"hint='Shown when the preview of the draft is public.'"`
	PublicPreviewLink        string `i18n:"hint='Label of the link of the public preview.'"`
	PublicPreviewVisit       string `i18n:"hint='Button that opens the public preview.'"`
	PublicPreviewWarn        string `i18n:"hint='Warning in the form that makes the preview public.'"`
	PublicPreviewExpiresHint string `i18n:"hint='Hint of the date the public preview ends.'"`
	PublicPreviewOpened      string `i18n:"hint='Shown once the preview was made public.'"`
	PublicPreviewEnded       string `i18n:"hint='Shown once the public preview was ended.'"`
	Shared                   string `i18n:"hint='Shown once a draft was shared.'"`
	Revoked                  string `i18n:"hint='Shown once a sharing was revoked.'"`
	SharesTab                string `i18n:"hint='Tab of the sharings of the draft.'"`
	DraftsTab                string `i18n:"hint='Tab of the drafts the user reaches.'"`
	SharesHint               string `i18n:"hint='Explains the sharing of a draft.'"`
	SharesNone               string `i18n:"hint='Said when the draft is shared with no one.'"`
	SharesHistory            string `i18n:"hint='Title of the past sharings (revoked, expired).'"`
	ShareUser                string `i18n:"hint='Header of the column of whom a draft is shared with.'"`
	ShareSince               string `i18n:"hint='Header of the column of when and by whom it was shared.'"`
	ShareState               string `i18n:"hint='Header of the column of how a sharing stands.'"`
	ShareForever             string `i18n:"hint='A sharing with no end.'"`
	ShareUntil               string `i18n:"hint='A sharing until a date.'"`
	ShareExpiredOn           string `i18n:"hint='A sharing that expired.'"`
	ShareRevokedOn           string `i18n:"hint='A sharing revoked.'"`
	ShareExpiresHint         string `i18n:"hint='Hint of the date a sharing ends.'"`
	ShareFormUser            string `i18n:"hint='Label of whom to share the draft with.'"`
	ShareFormExpiresAt       string `i18n:"hint='Label of the date a sharing ends.'"`
	ErrShareUser             string `i18n:"hint='Error of a sharing with no user chosen.'"`
	ErrShareSelf             string `i18n:"hint='Error of a draft shared with its owner.'"`
	ErrShareExpired          string `i18n:"hint='Error of a sharing ending in the past.'"`
	SharedWithMe             string `i18n:"hint='Title of the drafts shared with the user.'"`
	SharedWithMeNone         string `i18n:"hint='Said when no draft is shared with the user.'"`
	AllDrafts                string `i18n:"hint='Title of every draft (with the permission).'"`
	NoDrafts                 string `i18n:"hint='Said when there is no draft.'"`
	DraftOwner               string `i18n:"hint='Header of the column of the owner of a draft.'"`
	LastCommit               string `i18n:"hint='Header of the column of the last commit of a draft.'"`
	SharedWithCol            string `i18n:"hint='Header of the column of whom a draft is shared with.'"`
	ChangesCount             string `i18n:"hint='How many files of a draft changed, not committed.'"`
	DraftOf                  string `i18n:"hint='Title of the page of the draft of another user.'"`
	ViaShare                 string `i18n:"hint='Said of a draft reached by a sharing.'"`
	ViaDrafts                string `i18n:"hint='Said of a draft reached by the permission of every draft.'"`
	GitHookHint              string `i18n:"hint='Explains the hook that signs the commits of a clone.'"`
	GitHook                  string `i18n:"hint='Label of the command that installs the hook.'"`
	DraftOther               string `i18n:"hint='Tab of the draft of another user (shared, or every draft).'"`
	DraftHintOther           string `i18n:"hint='Explains the draft of another user: whoever reaches it edits it; the site changes only when published.'"`
	OpenDraft                string `i18n:"hint='Link to the page of a draft, in the list of the drafts.'"`
	GitTab                   string `i18n:"hint='Tab of the repositories by git.'"`
	GitDraftTitle            string `i18n:"hint='Title of the section of the draft of the user by git.'"`
	GitOtherDraftTitle       string `i18n:"hint='Title of the section of the draft of another user by git.'"`
	GitSiteTitle             string `i18n:"hint='Title of the section of the site by git.'"`
	GitDraftWhat             string `i18n:"hint='What a push to the draft does.'"`
	GitSiteWhat              string `i18n:"hint='What a push to the site does.'"`
	GitURLLabel              string `i18n:"hint='Label of the URL of a repository by git.'"`
	GitCloneLabel            string `i18n:"hint='Label of the command that clones a repository.'"`
	GitPermsTitle            string `i18n:"hint='Title of the permissions a repository by git asks.'"`
	GitPermsMissing          string `i18n:"hint='Said when the user lacks a permission a repository asks.'"`
	GitPermGet               string `i18n:"hint='A permission of a repository by git: seeing the files.'"`
	GitPermPush              string `i18n:"hint='A permission of a repository by git: pushing.'"`
	GitKeyHint               string `i18n:"hint='Explains the access key for git.'"`
	GitKeyHelp               string `i18n:"hint='Link to the documentation of git with an access key.'"`
	GitPermsAll              string `i18n:"hint='Label of all the permissions a repository asks, as one (a policy).'"`
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
	BackToAdmin:        "Admin panel",
	SignOut:            "Sign out",
	ThemeLight:         "Light theme",
	ThemeDark:          "Dark theme",
	CommitAction:       "Commit",
	UpdateAction:       "Update",
	PublishAction:      "Publish",
	DiscardAction:      "Discard",
	CopyDiff:           "Copy",
	DiffCopied:         "Copied!",
	CopyDiffError:      "Press Ctrl+C to copy",
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
	ShareAction:        "Share",
	ShareAction_Desc:   "Share your draft with other users: they reach it whole — the editor, git, the preview — and commit in it with their own name; publishing still asks its permission.",
	DraftsAction:       "Every draft",
	DraftsAction_Desc:  "See and reach every user's draft — its editor, git, preview — and revoke any sharing.",
	RevokeAction:       "Revoke",

	PublicPreview:            "Public preview",
	PublicPreviewAction:      "Make the preview public",
	PublicPreviewAction_Desc: "Make the preview of your draft seen by anyone with its link, with no login — until the date you choose or until you end it.",
	PublicPreviewOffAction:   "End the public preview",
	PublicPreviewNone:        "Only who logs in sees the preview of your draft. Make it public to show it to someone without an account.",
	PublicPreviewOn:          "Anyone with the link sees the site as your draft makes it, with no login",
	PublicPreviewLink:        "Link",
	PublicPreviewVisit:       "Open",
	PublicPreviewWarn:        "Anyone who has the link sees the site as your draft makes it — its uncommitted changes too —, with no login. The link is not indexed by search engines. Ending it, or making it again, stops the link.",
	PublicPreviewExpiresHint: "Empty: until you end it.",
	PublicPreviewOpened:      "The preview of your draft is public.",
	PublicPreviewEnded:       "The public preview was ended: its link no longer opens.",
	Shared:                   "Shared with %s.",
	Revoked:                  "Sharing with %s revoked.",
	SharesTab:                "Sharing",
	DraftsTab:                "Drafts",
	SharesHint:               "Whom you share your draft with reaches it whole — the editor, git, the preview — and commits in it with their own name, until the date you choose or until you revoke it. Publishing still asks the permission to publish.",
	SharesNone:               "Not shared with anyone.",
	SharesHistory:            "Past sharings",
	ShareUser:                "User",
	ShareSince:               "Shared on, by",
	ShareState:               "State",
	ShareForever:             "until revoked",
	ShareUntil:               "until %s",
	ShareExpiredOn:           "expired on %s",
	ShareRevokedOn:           "revoked on %s by %s",
	ShareExpiresHint:         "Empty: until revoked.",
	ShareFormUser:            "User",
	ShareFormExpiresAt:       "Until",
	ErrShareUser:             "Choose a user.",
	ErrShareSelf:             "A draft is not shared with its owner.",
	ErrShareExpired:          "The date is past.",
	SharedWithMe:             "Shared with you",
	SharedWithMeNone:         "No draft is shared with you.",
	AllDrafts:                "Every draft",
	NoDrafts:                 "No draft yet.",
	DraftOwner:               "Draft of",
	LastCommit:               "Last commit",
	SharedWithCol:            "Shared with",
	ChangesCount:             "%d not committed",
	DraftOf:                  "Draft of %s",
	ViaShare:                 "shared with you",
	ViaDrafts:                "you reach every draft",
	GitHookHint:              "Each commit pushed says who made it on this site (Site-User). In each clone, install the hook that adds it:",
	GitHook:                  "The hook commit-msg",
	DraftOther:               "Draft",
	DraftHintOther:           "You edit the draft of another user, with them: your commits are yours; the site changes only when someone publishes.",
	OpenDraft:                "Open",
	GitTab:                   "Git",
	GitDraftTitle:            "Your draft by git",
	GitOtherDraftTitle:       "The draft of %s by git",
	GitSiteTitle:             "The site by git",
	GitDraftWhat:             "A push changes the draft: the editor and the preview show it at once; publish from the page.",
	GitSiteWhat:              "A push publishes at once — checking what Publish does: the templates compile, a fast-forward, the files of the site untouched.",
	GitURLLabel:              "Address",
	GitCloneLabel:            "Clone",
	GitPermsTitle:            "The permissions it asks",
	GitPermsMissing:          "Without the ones marked ✗, git refuses what asks them.",
	GitPermGet:               "Clone and fetch (each file)",
	GitPermPush:              "Push",
	GitKeyHint:               "Use an access key as the password: it needs these same permissions, and your role too.",
	GitKeyHelp:               "How to",
	GitPermsAll:              "All of them, one permission",
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
	BackToAdmin:        "Painel admin",
	SignOut:            "Sair",
	ThemeLight:         "Tema claro",
	ThemeDark:          "Tema escuro",
	CommitAction:       "Commit",
	UpdateAction:       "Atualizar",
	PublishAction:      "Publicar",
	DiscardAction:      "Descartar",
	CopyDiff:           "Copiar",
	DiffCopied:         "Copiado!",
	CopyDiffError:      "Use Ctrl+C para copiar",
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
	ShareAction:        "Compartilhar",
	ShareAction_Desc:   "Compartilhar o seu rascunho com outros usuários: eles o acessam por inteiro — o editor, o git, o preview — e fazem commits nele com o próprio nome; publicar continua pedindo a permissão de publicar.",
	DraftsAction:       "Todos os rascunhos",
	DraftsAction_Desc:  "Ver e acessar o rascunho de qualquer usuário — o editor, o git, o preview — e revogar qualquer compartilhamento.",
	RevokeAction:       "Revogar",

	PublicPreview:            "Preview público",
	PublicPreviewAction:      "Tornar o preview público",
	PublicPreviewAction_Desc: "Deixar o preview do seu rascunho visível a quem tiver o link, sem login — até a data escolhida ou até você encerrar.",
	PublicPreviewOffAction:   "Encerrar o preview público",
	PublicPreviewNone:        "Só quem entra no admin vê o preview do seu rascunho. Torne-o público para mostrá-lo a quem não tem conta.",
	PublicPreviewOn:          "Quem tiver o link vê o site como o seu rascunho o faz, sem login",
	PublicPreviewLink:        "Link",
	PublicPreviewVisit:       "Abrir",
	PublicPreviewWarn:        "Quem tiver o link vê o site como o seu rascunho o faz — também as mudanças não commitadas —, sem login. O link não é indexado pelos buscadores. Encerrar, ou tornar público de novo, invalida o link.",
	PublicPreviewExpiresHint: "Vazio: até você encerrar.",
	PublicPreviewOpened:      "O preview do seu rascunho está público.",
	PublicPreviewEnded:       "O preview público foi encerrado: o link não abre mais.",
	Shared:                   "Compartilhado com %s.",
	Revoked:                  "Compartilhamento com %s revogado.",
	SharesTab:                "Compartilhamentos",
	DraftsTab:                "Rascunhos",
	SharesHint:               "Com quem você compartilha o seu rascunho o acessa por inteiro — o editor, o git, o preview — e faz commits nele com o próprio nome, até a data que você escolher ou até você revogar. Publicar continua pedindo a permissão de publicar.",
	SharesNone:               "Não está compartilhado com ninguém.",
	SharesHistory:            "Compartilhamentos anteriores",
	ShareUser:                "Usuário",
	ShareSince:               "Compartilhado em, por",
	ShareState:               "Situação",
	ShareForever:             "até ser revogado",
	ShareUntil:               "até %s",
	ShareExpiredOn:           "expirou em %s",
	ShareRevokedOn:           "revogado em %s por %s",
	ShareExpiresHint:         "Vazio: até ser revogado.",
	ShareFormUser:            "Usuário",
	ShareFormExpiresAt:       "Até",
	ErrShareUser:             "Escolha um usuário.",
	ErrShareSelf:             "Um rascunho não é compartilhado com o próprio dono.",
	ErrShareExpired:          "A data já passou.",
	SharedWithMe:             "Compartilhados com você",
	SharedWithMeNone:         "Nenhum rascunho está compartilhado com você.",
	AllDrafts:                "Todos os rascunhos",
	NoDrafts:                 "Nenhum rascunho ainda.",
	DraftOwner:               "Rascunho de",
	LastCommit:               "Último commit",
	SharedWithCol:            "Compartilhado com",
	ChangesCount:             "%d sem commit",
	DraftOf:                  "Rascunho de %s",
	ViaShare:                 "compartilhado com você",
	ViaDrafts:                "você acessa todos os rascunhos",
	GitHookHint:              "Cada commit enviado diz quem o fez neste site (Site-User). Em cada clone, instale o hook que o acrescenta:",
	GitHook:                  "O hook commit-msg",
	DraftOther:               "Rascunho",
	DraftHintOther:           "Você edita o rascunho de outro usuário, junto com ele: os seus commits são seus; o site muda só quando alguém publica.",
	OpenDraft:                "Abrir",
	GitTab:                   "Git",
	GitDraftTitle:            "O seu rascunho pelo git",
	GitOtherDraftTitle:       "O rascunho de %s pelo git",
	GitSiteTitle:             "O site pelo git",
	GitDraftWhat:             "Um push muda o rascunho: o editor e o preview mostram na hora; publique pela página.",
	GitSiteWhat:              "Um push publica direto — conferindo o mesmo que Publicar: os templates compilam, um fast-forward, os arquivos do site intactos.",
	GitURLLabel:              "Endereço",
	GitCloneLabel:            "Clonar",
	GitPermsTitle:            "As permissões que pede",
	GitPermsMissing:          "Sem as marcadas com ✗, o git recusa o que as pede.",
	GitPermGet:               "Clonar e buscar (cada arquivo)",
	GitPermPush:              "Enviar (push)",
	GitKeyHint:               "Use uma chave de acesso como senha: ela precisa destas mesmas permissões, e o seu papel também.",
	GitKeyHelp:               "Como",
	GitPermsAll:              "Todas, numa permissão",
}
