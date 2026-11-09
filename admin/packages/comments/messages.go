package comments

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// MessagesKey is the module of the words of the comments.
const MessagesKey i18n.ModuleKey = "admin/packages/comments"

// Messages are the words of the comments.
type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	None              string `i18n:"hint='The comments of a record: none yet.'"`
	Placeholder       string `i18n:"hint='Placeholder of the box that writes a comment.'"`
	ReplyPlaceholder  string `i18n:"hint='Placeholder of the box that answers a comment.'"`
	Send              string `i18n:"hint='Button that sends a comment.'"`
	Reply             string `i18n:"hint='Button that answers a comment.'"`
	Edit              string `i18n:"hint='Button that edits ones own comment.'"`
	Cancel            string `i18n:"hint='Button that cancels writing or editing a comment.'"`
	Delete            string `i18n:"hint='Button that deletes ones own comment.'"`
	DeleteConfirm     string `i18n:"hint='Asks whether to delete a comment.'"`
	DeleteYes         string `i18n:"hint='Confirms deleting a comment.'"`
	RemovedUser       string `i18n:"hint='The author of a comment whose user is gone.'"`
	Empty             string `i18n:"hint='Shown when a comment is sent with no text and no file.'"`
	Attachments       string `i18n:"hint='Field of the files of a comment.'"`
	EditFile          string `i18n:"hint='Button of an image chosen for a comment: edits it before it is sent.'"`
	RemoveFile        string `i18n:"hint='Button of a file chosen for a comment: leaves it out.'"`
	FileTooLarge      string `i18n:"hint='A file of a comment larger than its limit: %s its name, %s its size, %s the limit.'"`
	MessageTooLarge   string `i18n:"hint='The files of a comment larger together than the limit: %s their size, %s the limit.'"`
	DropFiles         string `i18n:"hint='The place the files of a comment are dropped on, or clicked to choose them.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "The comments of the records: a conversation, its answers.",
	None:              "No comments yet.",
	Placeholder:       "Write a comment…",
	ReplyPlaceholder:  "Write an answer…",
	Send:              "Comment",
	Reply:             "Reply",
	Edit:              "Edit",
	Cancel:            "Cancel",
	Delete:            "Delete",
	DeleteConfirm:     "Delete this comment?",
	DeleteYes:         "Yes, delete",
	RemovedUser:       "Removed user",
	Empty:             "Write something to comment.",
	Attachments:       "Attachments",
	EditFile:          "Edit the image",
	RemoveFile:        "Remove",
	FileTooLarge:      "%s has %s, more than the %s a file may have.",
	MessageTooLarge:   "The files have %s together, more than the %s a message may have.",
	DropFiles:         "Drop files here, or click to choose them",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "Os comentários dos registros: uma conversa, as suas respostas.",
	None:              "Nenhum comentário ainda.",
	Placeholder:       "Escreva um comentário…",
	ReplyPlaceholder:  "Escreva uma resposta…",
	Send:              "Comentar",
	Reply:             "Responder",
	Edit:              "Editar",
	Cancel:            "Cancelar",
	Delete:            "Excluir",
	DeleteConfirm:     "Excluir este comentário?",
	DeleteYes:         "Sim, excluir",
	RemovedUser:       "Usuário removido",
	Empty:             "Escreva algo para comentar.",
	Attachments:       "Anexos",
	EditFile:          "Editar a imagem",
	RemoveFile:        "Remover",
	FileTooLarge:      "%s tem %s, mais que os %s que um arquivo pode ter.",
	MessageTooLarge:   "Os arquivos somam %s, mais que os %s que uma mensagem pode ter.",
	DropFiles:         "Arraste os arquivos para cá, ou clique para escolhê-los",
}

// ConfigureMessages registers the words.
func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

// GetMessages are the words in the language of ctx.
func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}
