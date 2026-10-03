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

	CommitAction      string `i18n:"label='Commit', hint='Action that records the changes of the draft, with a message.'"`
	UpdateAction      string `i18n:"label='Update', hint='Action that brings into the draft what others published.'"`
	PublishAction     string `i18n:"label='Publish', hint='Action that puts the commits of the draft on the site.'"`
	DiscardAction     string `i18n:"label='Discard', hint='Action that drops the changes of a file not committed.'"`
	ResetAction       string `i18n:"label='Start over', hint='Action that drops the draft: a new one is made from the site.'"`
	CommitFormMessage string `i18n:"label='Message', hint='Label of the message of a commit.'"`

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
	ModuleDescription: "The editor of the files of the site: drafts, commits and publishing.",
	Title:             "Site files",
	Draft:             "Your draft",
	DraftHint:         "You edit a draft of your own: the site changes only when you publish.",
	Ahead:             "%d commit(s) to publish",
	Behind:            "%d commit(s) published by others: update",
	UpToDate:          "Up to date with the site",
	Changes:           "Changes not committed",
	NoChanges:         "Nothing changed since the last commit.",
	History:           "History",
	Files:             "Files",
	Preview:           "Preview the site",
	Help:              "Help",
	CommitAction:      "Commit",
	UpdateAction:      "Update",
	PublishAction:     "Publish",
	DiscardAction:     "Discard",
	ResetAction:       "Start over",
	CommitFormMessage: "Message",
	Committed:         "Committed: %s",
	Updated:           "The draft has what was published.",
	Published:         "Published: the site has your commits.",
	Discarded:         "Changes of %s discarded.",
	Reset:             "The draft was dropped.",
	ErrUncommitted:    "Commit or discard the changes before publishing.",
	ErrBehind:         "Others published since: update the draft before publishing.",
	ErrSiteChanged:    "Files of the site were changed out of the editor: nothing was overwritten.",
	ErrNothing:        "Nothing to commit.",
	ErrInvalid:        "The files of the draft do not work: %s",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "O editor dos arquivos do site: rascunhos, commits e publicação.",
	Title:             "Arquivos do site",
	Draft:             "Seu rascunho",
	DraftHint:         "Você edita um rascunho só seu: o site muda só quando você publica.",
	Ahead:             "%d commit(s) para publicar",
	Behind:            "%d commit(s) publicado(s) por outros: atualize",
	UpToDate:          "Igual ao site",
	Changes:           "Alterações sem commit",
	NoChanges:         "Nada mudou desde o último commit.",
	History:           "Histórico",
	Files:             "Arquivos",
	Preview:           "Ver o site do rascunho",
	Help:              "Ajuda",
	CommitAction:      "Commit",
	UpdateAction:      "Atualizar",
	PublishAction:     "Publicar",
	DiscardAction:     "Descartar",
	ResetAction:       "Recomeçar",
	CommitFormMessage: "Mensagem",
	Committed:         "Commit feito: %s",
	Updated:           "O rascunho tem o que foi publicado.",
	Published:         "Publicado: o site tem os seus commits.",
	Discarded:         "Alterações de %s descartadas.",
	Reset:             "O rascunho foi descartado.",
	ErrUncommitted:    "Faça o commit ou descarte as alterações antes de publicar.",
	ErrBehind:         "Outros publicaram depois: atualize o rascunho antes de publicar.",
	ErrSiteChanged:    "Arquivos do site foram alterados fora do editor: nada foi sobrescrito.",
	ErrNothing:        "Nada para o commit.",
	ErrInvalid:        "Os arquivos do rascunho não funcionam: %s",
}
