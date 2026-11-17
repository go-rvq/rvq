package fs_tools

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// FilesMessagesKey is the module of the words of the components (FilesMessages).
const FilesMessagesKey i18n.ModuleKey = "rvq-admin/fs-tools/files"

// FilesMessages are the words of <vx-files> and <vx-terminal>, in the
// user's language: the page's vue-i18n has them in the namespace vxFiles
// (UseAssets), where the components read them; {n}, {path}… are
// vue-i18n's placeholders, kept as they are.
type FilesMessages struct {
	ModuleDescription string `json:"-" i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Home              string `json:"home" i18n:"hint='The root of the explorer in its path, when it has no title.'"`
	Back              string `json:"back" i18n:"hint='Button of the explorer’s bar: the folder before.'"`
	Forward           string `json:"forward" i18n:"hint='Button of the explorer’s bar: the folder after.'"`
	Up                string `json:"up" i18n:"hint='Button of the explorer’s bar: the parent folder.'"`
	TypePath          string `json:"typePath" i18n:"hint='Button beside the path: edit it as text.'"`
	Filter            string `json:"filter" i18n:"hint='Placeholder of the filter of the files.'"`
	Refresh           string `json:"refresh" i18n:"hint='Button of the explorer’s bar: read the folder again.'"`
	New               string `json:"new" i18n:"hint='Button of the explorer’s bar: the menu of what to create.'"`
	Folder            string `json:"folder" i18n:"hint='Item of the New menu: a folder.'"`
	TextFile          string `json:"textFile" i18n:"hint='Item of the New menu: a text file.'"`
	UploadFiles       string `json:"uploadFiles" i18n:"hint='Item of a menu: choose files to upload.'"`
	UploadFolder      string `json:"uploadFolder" i18n:"hint='Item of a menu: choose a folder to upload.'"`
	Fetch             string `json:"fetch" i18n:"hint='Item of a menu (and title of its dialog): a file by its URL.'"`
	View              string `json:"view" i18n:"hint='Accessible name of the group of the view buttons.'"`
	ViewList          string `json:"viewList" i18n:"hint='View of the files: a list.'"`
	ViewTree          string `json:"viewTree" i18n:"hint='View of the files: a tree.'"`
	ViewIcons         string `json:"viewIcons" i18n:"hint='View of the files: icons.'"`
	TermClose         string `json:"termClose" i18n:"hint='Terminal button, the terminal open.'"`
	TermOpen          string `json:"termOpen" i18n:"hint='Terminal button, the terminal closed.'"`
	ThumbsHide        string `json:"thumbsHide" i18n:"hint='Images button, the images shown.'"`
	ThumbsShow        string `json:"thumbsShow" i18n:"hint='Images button, the images not shown.'"`
	HiddenHide        string `json:"hiddenHide" i18n:"hint='Hidden files button, them shown.'"`
	HiddenShow        string `json:"hiddenShow" i18n:"hint='Hidden files button, them not shown.'"`
	GitOpen           string `json:"gitOpen" i18n:"hint='Button of a folder that is a git repository.'"`
	CloneCopy         string `json:"cloneCopy" i18n:"hint='Button and menu item: copy the URL to clone the repository.'"`
	IdeOpen           string `json:"ideOpen" i18n:"hint='Button of the explorer’s bar: the IDE in this folder.'"`
	Uploading         string `json:"uploading" i18n:"hint='Progress of an upload: files done, their total, bytes sent of the total.', fields=(;'{done}'='done', '{total}'='total', '{sent}'='sent', '{bytes}'='bytes')"`
	DropHere          string `json:"dropHere" i18n:"hint='Shown while files are dragged over the explorer: the folder they go to.', fields=(;'{path}'='path')"`
	Name              string `json:"name" i18n:"hint='Column of the list, and field of a dialog: the name.'"`
	Size              string `json:"size" i18n:"hint='Column of the list: the size.'"`
	Dimensions        string `json:"dimensions" i18n:"hint='Column of the list: the dimensions of an image.'"`
	Modified          string `json:"modified" i18n:"hint='Column of the list: when the file changed.'"`
	Permissions       string `json:"permissions" i18n:"hint='Column of the list: the permissions of the file (rwx).'"`
	Columns           string `json:"columns" i18n:"hint='Button and title of the menu of the columns.'"`
	ImageDims         string `json:"imageDims" i18n:"hint='Item of the menu of the columns: show the dimensions.'"`
	GitRepo           string `json:"gitRepo" i18n:"hint='Title of the icon of a folder that is a git repository.'"`
	NoMatch           string `json:"noMatch" i18n:"hint='Shown when the filter hides every file.'"`
	EmptyFolder       string `json:"emptyFolder" i18n:"hint='Shown for an empty folder.'"`
	ItemsOne          string `json:"itemsOne" i18n:"hint='Status bar: one file in the folder.'"`
	ItemsMany         string `json:"itemsMany" i18n:"hint='Status bar (and the delete dialog): how many files.', fields=(;'{n}'='n')"`
	Selected          string `json:"selected" i18n:"hint='Status bar: how many files are selected.', fields=(;'{n}'='n')"`
	Hint              string `json:"hint" i18n:"hint='Status bar: how to use the explorer.'"`
	TermHeight        string `json:"termHeight" i18n:"hint='Accessible name of the bar between the files and the terminal.'"`
	TermSash          string `json:"termSash" i18n:"hint='Title of the bar between the files and the terminal.'"`
	Open              string `json:"open" i18n:"hint='Item of the menu of a file: open it.'"`
	ExtractDots       string `json:"extractDots" i18n:"hint='Item of the menu of an archive: extract it (a dialog follows).'"`
	Edit              string `json:"edit" i18n:"hint='Item of the menu of a file, and button of the viewer: edit it.'"`
	OpenIde           string `json:"openIde" i18n:"hint='Item of the menu of a folder: the IDE in it.'"`
	OpenGit           string `json:"openGit" i18n:"hint='Item of the menu of a repository: its Git page.'"`
	Download          string `json:"download" i18n:"hint='Item of the menu of a file, button of the image viewer and of the fetch dialog.'"`
	DownloadN         string `json:"downloadN" i18n:"hint='Title of the button of the download of several items selected (it opens a menu).', fields=(;'{n}'='n')"`
	DownloadNTgz      string `json:"downloadNTgz" i18n:"hint='Item of the menu of several items: download them in one tar.gz.', fields=(;'{n}'='n')"`
	DownloadNEach     string `json:"downloadNEach" i18n:"hint='Item of the menu of several items: download each by itself.', fields=(;'{n}'='n')"`
	DownloadTgz       string `json:"downloadTgz" i18n:"hint='Item of the menu of a folder: download it as tar.gz.'"`
	RenameF2          string `json:"renameF2" i18n:"hint='Item of the menu of a file: rename it.'"`
	CopyTo            string `json:"copyTo" i18n:"hint='Item of the menu of the files: copy them (a dialog follows).'"`
	MoveTo            string `json:"moveTo" i18n:"hint='Item of the menu of the files: move them (a dialog follows).'"`
	DeleteKey         string `json:"deleteKey" i18n:"hint='Item of the menu of the files: delete them.'"`
	NewFolder         string `json:"newFolder" i18n:"hint='Item of the menu of the folder, and title of its dialog.'"`
	NewTextFile       string `json:"newTextFile" i18n:"hint='Item of the menu of the folder: a text file.'"`
	DownloadFolder    string `json:"downloadFolder" i18n:"hint='Item of the menu of the folder open (and title of the download button with nothing selected): download it as tar.gz.'"`
	SelectAll         string `json:"selectAll" i18n:"hint='Item of the menu of the folder: select every file.'"`
	NewFile           string `json:"newFile" i18n:"hint='Title of the dialog of a new text file.'"`
	Rename            string `json:"rename" i18n:"hint='Title and button of the dialog that renames.'"`
	CopyToTitle       string `json:"copyToTitle" i18n:"hint='Title of the dialog that copies.'"`
	MoveToTitle       string `json:"moveToTitle" i18n:"hint='Title of the dialog that moves.'"`
	Delete            string `json:"delete" i18n:"hint='Title and button of the dialog that deletes.'"`
	Extract           string `json:"extract" i18n:"hint='Title and button of the dialog that extracts an archive.'"`
	Copy              string `json:"copy" i18n:"hint='Button of the dialog that copies.'"`
	Move              string `json:"move" i18n:"hint='Button of the dialog that moves.'"`
	Create            string `json:"create" i18n:"hint='Button of the dialogs that create a folder or a file.'"`
	ToFolder          string `json:"toFolder" i18n:"hint='Field of the copy, move and extract dialogs: where to.'"`
	ToFolderHint      string `json:"toFolderHint" i18n:"hint='Hint of the field To the folder.'"`
	Address           string `json:"address" i18n:"hint='Field of the fetch dialog: the URL.'"`
	SaveAs            string `json:"saveAs" i18n:"hint='Field of the fetch dialog: the name of the file.'"`
	InPath            string `json:"inPath" i18n:"hint='Hint of the field Save as: the folder.', fields=(;'{path}'='path')"`
	DeleteConfirm     string `json:"deleteConfirm" i18n:"hint='Text of the delete dialog: what (a name, or how many).', fields=(;'{what}'='what')"`
	Cancel            string `json:"cancel" i18n:"hint='Button that cancels a dialog or an edit.'"`
	ExistTitle        string `json:"existTitle" i18n:"hint='Title of the dialog of the uploaded files that exist.'"`
	ExistText         string `json:"existText" i18n:"hint='Text of the dialog of the uploaded files that exist.', fields=(;'{n}'='n', '{dir}'='dir')"`
	ColFile           string `json:"colFile" i18n:"hint='Column of the dialog of the files that exist: the file.'"`
	ColThere          string `json:"colThere" i18n:"hint='Column of the dialog of the files that exist: the one in the folder.'"`
	ColSent           string `json:"colSent" i18n:"hint='Column of the dialog of the files that exist: the one uploaded.'"`
	ColDiff           string `json:"colDiff" i18n:"hint='Column of the dialog of the files that exist: how they differ.'"`
	KeepTheirs        string `json:"keepTheirs" i18n:"hint='Button of the dialog of the files that exist: keep the ones there.'"`
	Replace           string `json:"replace" i18n:"hint='Button of the dialog of the files that exist: replace them.'"`
	CloseEsc          string `json:"closeEsc" i18n:"hint='Button that closes the viewer of a file or an image.'"`
	Changed           string `json:"changed" i18n:"hint='Title of the mark of a file edited and not saved.'"`
	Text              string `json:"text" i18n:"hint='The language of a file of plain text, in the viewer.'"`
	Truncated         string `json:"truncated" i18n:"hint='Chip of the viewer: the file is shown in part.'"`
	NoWrite           string `json:"noWrite" i18n:"hint='Chip of the viewer: the file is not written.'"`
	SaveCtrlS         string `json:"saveCtrlS" i18n:"hint='Button of the editor: save the file.'"`
	EditCtrlE         string `json:"editCtrlE" i18n:"hint='Title of the button of the viewer that edits.'"`
	IdeFolder         string `json:"ideFolder" i18n:"hint='Title of the IDE button of the viewer.'"`
	FitScreen         string `json:"fitScreen" i18n:"hint='Button of the image viewer: fit the image.'"`
	ActualSize        string `json:"actualSize" i18n:"hint='Button of the image viewer: the image in its size.'"`
	Prev              string `json:"prev" i18n:"hint='Button of the image viewer: the image before.'"`
	Next              string `json:"next" i18n:"hint='Button of the image viewer: the image after.'"`
	DiscardTitle      string `json:"discardTitle" i18n:"hint='Title of the dialog closing an edited file.'"`
	DiscardText       string `json:"discardText" i18n:"hint='Text of the dialog closing an edited file.', fields=(;'{path}'='path')"`
	KeepEditing       string `json:"keepEditing" i18n:"hint='Button of the dialog closing an edited file: go back.'"`
	Discard           string `json:"discard" i18n:"hint='Button of the dialog closing an edited file: lose the changes.'"`
	Expired           string `json:"expired" i18n:"hint='Error of a request refused (the session ended).'"`
	DoneMkdir         string `json:"doneMkdir" i18n:"hint='Message after a folder is created.'"`
	DoneNewfile       string `json:"doneNewfile" i18n:"hint='Message after a file is created.'"`
	DoneRename        string `json:"doneRename" i18n:"hint='Message after a rename.'"`
	DoneCopy          string `json:"doneCopy" i18n:"hint='Message after a copy.'"`
	DoneMove          string `json:"doneMove" i18n:"hint='Message after a move.'"`
	DoneDelete        string `json:"doneDelete" i18n:"hint='Message after a delete.'"`
	DoneExtract       string `json:"doneExtract" i18n:"hint='Message after an archive is extracted.'"`
	DoneExtractN      string `json:"doneExtractN" i18n:"hint='Message after an archive is extracted: how many files.', fields=(;'{n}'='n')"`
	DoneFetch         string `json:"doneFetch" i18n:"hint='Message after a file is downloaded from the internet.'"`
	DoneFetchPath     string `json:"doneFetchPath" i18n:"hint='Message after a file is downloaded from the internet: where.', fields=(;'{path}'='path')"`
	CloneCopied       string `json:"cloneCopied" i18n:"hint='Message after the clone URL is copied.', fields=(;'{url}'='url')"`
	Saved             string `json:"saved" i18n:"hint='Message after a file is saved.'"`
	NotSent           string `json:"notSent" i18n:"hint='Message of the files that failed to upload.', fields=(;'{list}'='list')"`
	SentOne           string `json:"sentOne" i18n:"hint='Message after one file is uploaded.'"`
	SentMany          string `json:"sentMany" i18n:"hint='Message after several files are uploaded.', fields=(;'{n}'='n')"`
	SameTime          string `json:"sameTime" i18n:"hint='Difference of time of a file uploaded and the one there: none.'"`
	Newer             string `json:"newer" i18n:"hint='Difference of time: the one uploaded is newer (by 3 h).', fields=(;'{n}'='n')"`
	Older             string `json:"older" i18n:"hint='Difference of time: the one uploaded is older (by 3 h).', fields=(;'{n}'='n')"`
	Day               string `json:"day" i18n:"hint='A day, in a difference of time (1 day).'"`
	Days              string `json:"days" i18n:"hint='Days, in a difference of time (3 days).'"`
	SameSize          string `json:"sameSize" i18n:"hint='Difference of size of a file uploaded and the one there: none.'"`
	MovedTo           string `json:"movedTo" i18n:"hint='Message after files are dragged into a folder.', fields=(;'{name}'='name')"`
	ArchiveName       string `json:"archiveName" i18n:"hint='Name of the tar.gz of the root folder.'"`
	Connected         string `json:"connected" i18n:"hint='Chip of the terminal: connected.'"`
	Closed            string `json:"closed" i18n:"hint='Chip of the terminal: the session ended.'"`
	Connecting        string `json:"connecting" i18n:"hint='Chip of the terminal: connecting.'"`
	Tmux              string `json:"tmux" i18n:"hint='Chip of a terminal kept by tmux.'"`
	Recorded          string `json:"recorded" i18n:"hint='Chip of a terminal whose session is recorded.'"`
	Search            string `json:"search" i18n:"hint='Button of the terminal: search its text.'"`
	SearchPrompt      string `json:"searchPrompt" i18n:"hint='Question of the search of the terminal.'"`
	FontDown          string `json:"fontDown" i18n:"hint='Button of the terminal: a smaller font.'"`
	FontUp            string `json:"fontUp" i18n:"hint='Button of the terminal: a larger font.'"`
	Reopen            string `json:"reopen" i18n:"hint='Button of the terminal: a new session.'"`
	Restore           string `json:"restore" i18n:"hint='Button of the terminal expanded: back.'"`
	Expand            string `json:"expand" i18n:"hint='Button of the terminal: over the whole window.'"`
	ConnClosed        string `json:"connClosed" i18n:"hint='Why the session ended: the connection closed.'"`
	SessionEnded      string `json:"sessionEnded" i18n:"hint='Shown when the session ended: how.', fields=(;'{exit}'='exit')"`
}

var FilesMessages_en_US = &FilesMessages{
	ModuleDescription: "The words of the explorer of files and its terminal.",
	Home:              "Home",
	Back:              "Back",
	Forward:           "Forward",
	Up:                "Up (Backspace)",
	TypePath:          "Type the path (Ctrl+L)",
	Filter:            "Filter",
	Refresh:           "Refresh",
	New:               "New",
	Folder:            "Folder",
	TextFile:          "Text file",
	UploadFiles:       "Upload files",
	UploadFolder:      "Upload a folder",
	Fetch:             "Download from the internet",
	View:              "View",
	ViewList:          "List",
	ViewTree:          "Tree",
	ViewIcons:         "Icons",
	TermClose:         "Close the terminal (F4)",
	TermOpen:          "Terminal in this folder (F4)",
	ThumbsHide:        "Don't show the images",
	ThumbsShow:        "Show the images",
	HiddenHide:        "Hide the hidden files",
	HiddenShow:        "Show the hidden files",
	GitOpen:           "Open this repository in Git",
	CloneCopy:         "Copy the clone URL (git)",
	IdeOpen:           "Open this folder in the IDE",
	Uploading:         "Uploading {done}/{total} · {sent} of {bytes}",
	DropHere:          "Drop to upload to {path}",
	Name:              "Name",
	Size:              "Size",
	Dimensions:        "Dimensions",
	Modified:          "Modified",
	Permissions:       "Permissions",
	Columns:           "Columns",
	ImageDims:         "Dimensions of the images",
	GitRepo:           "Git repository",
	NoMatch:           "Nothing with this filter.",
	EmptyFolder:       "Empty folder. Drag files or folders here.",
	ItemsOne:          "1 item",
	ItemsMany:         "{n} items",
	Selected:          "{n} selected",
	Hint:              "Drag files or folders to upload · double click opens · right button: more",
	TermHeight:        "Terminal height",
	TermSash:          "Drag to change the terminal's height (double click: the default height; ↑ ↓ by the keyboard)",
	Open:              "Open",
	ExtractDots:       "Extract…",
	Edit:              "Edit",
	OpenIde:           "Open in the IDE",
	OpenGit:           "Open in Git",
	Download:          "Download",
	DownloadN:         "Download ({n})",
	DownloadNTgz:      "Download {n} items in one .tar.gz",
	DownloadNEach:     "Download {n} items one by one",
	DownloadTgz:       "Download (.tar.gz)",
	RenameF2:          "Rename (F2)",
	CopyTo:            "Copy to…",
	MoveTo:            "Move to…",
	DeleteKey:         "Delete (Delete)",
	NewFolder:         "New folder",
	NewTextFile:       "New text file",
	DownloadFolder:    "Download this folder (.tar.gz)",
	SelectAll:         "Select all (Ctrl+A)",
	NewFile:           "New file",
	Rename:            "Rename",
	CopyToTitle:       "Copy to",
	MoveToTitle:       "Move to",
	Delete:            "Delete",
	Extract:           "Extract",
	Copy:              "Copy",
	Move:              "Move",
	Create:            "Create",
	ToFolder:          "To the folder",
	ToFolderHint:      "A path of this root, like /a/b",
	Address:           "Address (http, https)",
	SaveAs:            "Save as (optional)",
	InPath:            "In {path}",
	DeleteConfirm:     "Delete {what}? Folders go with all they have. It can't be undone.",
	Cancel:            "Cancel",
	ExistTitle:        "Files that already exist",
	ExistText:         "{n} file(s) already exist in {dir}. Replace?",
	ColFile:           "File",
	ColThere:          "There",
	ColSent:           "Sent",
	ColDiff:           "Difference",
	KeepTheirs:        "Keep theirs",
	Replace:           "Replace",
	CloseEsc:          "Close (Esc)",
	Changed:           "Changed",
	Text:              "text",
	Truncated:         "only the first 2 MB",
	NoWrite:           "no permission to write",
	SaveCtrlS:         "Save (Ctrl+S)",
	EditCtrlE:         "Edit (Ctrl+E)",
	IdeFolder:         "Open the folder in the IDE",
	FitScreen:         "Fit to the screen",
	ActualSize:        "Actual size",
	Prev:              "Previous (←)",
	Next:              "Next (→)",
	DiscardTitle:      "Discard the changes?",
	DiscardText:       "{path} has changes not saved.",
	KeepEditing:       "Keep editing",
	Discard:           "Discard",
	Expired:           "The page expired: reload it.",
	DoneMkdir:         "Folder created.",
	DoneNewfile:       "File created.",
	DoneRename:        "Renamed.",
	DoneCopy:          "Copied.",
	DoneMove:          "Moved.",
	DoneDelete:        "Deleted.",
	DoneExtract:       "Extracted.",
	DoneExtractN:      "Extracted ({n} files).",
	DoneFetch:         "Downloaded.",
	DoneFetchPath:     "Downloaded: {path}",
	CloneCopied:       "Clone URL copied: {url}",
	Saved:             "Saved.",
	NotSent:           "Not sent: {list}",
	SentOne:           "File uploaded.",
	SentMany:          "{n} files uploaded.",
	SameTime:          "same time",
	Newer:             "newer by {n}",
	Older:             "older by {n}",
	Day:               "day",
	Days:              "days",
	SameSize:          "same size",
	MovedTo:           "Moved to {name}.",
	ArchiveName:       "folder",
	Connected:         "connected",
	Closed:            "closed",
	Connecting:        "connecting…",
	Tmux:              "tmux (the session goes on when closed)",
	Recorded:          "recorded",
	Search:            "Search",
	SearchPrompt:      "Search in the terminal:",
	FontDown:          "Smaller font",
	FontUp:            "Larger font",
	Reopen:            "Reopen",
	Restore:           "Back to the size before",
	Expand:            "Expand the terminal (the whole window)",
	ConnClosed:        "connection closed",
	SessionEnded:      "Session ended: {exit}",
}

var FilesMessages_pt_BR = &FilesMessages{
	ModuleDescription: "As palavras do explorador de arquivos e do seu terminal.",
	Home:              "Início",
	Back:              "Voltar",
	Forward:           "Avançar",
	Up:                "Subir (Backspace)",
	TypePath:          "Digitar o caminho (Ctrl+L)",
	Filter:            "Filtrar",
	Refresh:           "Atualizar",
	New:               "Novo",
	Folder:            "Pasta",
	TextFile:          "Arquivo de texto",
	UploadFiles:       "Enviar arquivos",
	UploadFolder:      "Enviar uma pasta",
	Fetch:             "Baixar da internet",
	View:              "Visão",
	ViewList:          "Lista",
	ViewTree:          "Árvore",
	ViewIcons:         "Ícones",
	TermClose:         "Fechar o terminal (F4)",
	TermOpen:          "Terminal nesta pasta (F4)",
	ThumbsHide:        "Não mostrar as imagens",
	ThumbsShow:        "Mostrar as imagens",
	HiddenHide:        "Esconder os ocultos",
	HiddenShow:        "Mostrar os ocultos",
	GitOpen:           "Abrir este repositório no Git",
	CloneCopy:         "Copiar a URL de clone (git)",
	IdeOpen:           "Abrir esta pasta na IDE",
	Uploading:         "Enviando {done}/{total} · {sent} de {bytes}",
	DropHere:          "Solte para enviar a {path}",
	Name:              "Nome",
	Size:              "Tamanho",
	Dimensions:        "Dimensões",
	Modified:          "Modificado",
	Permissions:       "Permissões",
	Columns:           "Colunas",
	ImageDims:         "Dimensões das imagens",
	GitRepo:           "Repositório git",
	NoMatch:           "Nada com esse filtro.",
	EmptyFolder:       "Pasta vazia. Arraste arquivos ou pastas para cá.",
	ItemsOne:          "1 item",
	ItemsMany:         "{n} itens",
	Selected:          "{n} selecionados",
	Hint:              "Arraste arquivos ou pastas para enviar · duplo clique abre · botão direito: mais",
	TermHeight:        "Altura do terminal",
	TermSash:          "Arraste para mudar a altura do terminal (duplo clique: a altura padrão; ↑ ↓ pelo teclado)",
	Open:              "Abrir",
	ExtractDots:       "Extrair…",
	Edit:              "Editar",
	OpenIde:           "Abrir na IDE",
	OpenGit:           "Abrir no Git",
	Download:          "Baixar",
	DownloadN:         "Baixar ({n})",
	DownloadNTgz:      "Baixar {n} itens num .tar.gz",
	DownloadNEach:     "Baixar {n} itens um a um",
	DownloadTgz:       "Baixar (.tar.gz)",
	RenameF2:          "Renomear (F2)",
	CopyTo:            "Copiar para…",
	MoveTo:            "Mover para…",
	DeleteKey:         "Excluir (Delete)",
	NewFolder:         "Nova pasta",
	NewTextFile:       "Novo arquivo de texto",
	DownloadFolder:    "Baixar esta pasta (.tar.gz)",
	SelectAll:         "Selecionar tudo (Ctrl+A)",
	NewFile:           "Novo arquivo",
	Rename:            "Renomear",
	CopyToTitle:       "Copiar para",
	MoveToTitle:       "Mover para",
	Delete:            "Excluir",
	Extract:           "Extrair",
	Copy:              "Copiar",
	Move:              "Mover",
	Create:            "Criar",
	ToFolder:          "Para a pasta",
	ToFolderHint:      "Um caminho desta raiz, como /a/b",
	Address:           "Endereço (http, https)",
	SaveAs:            "Salvar como (opcional)",
	InPath:            "Em {path}",
	DeleteConfirm:     "Excluir {what}? As pastas vão com tudo o que têm. Não dá para desfazer.",
	Cancel:            "Cancelar",
	ExistTitle:        "Arquivos que já existem",
	ExistText:         "{n} arquivo(s) já existem em {dir}. Substituir?",
	ColFile:           "Arquivo",
	ColThere:          "Lá",
	ColSent:           "Enviado",
	ColDiff:           "Diferença",
	KeepTheirs:        "Manter os de lá",
	Replace:           "Substituir",
	CloseEsc:          "Fechar (Esc)",
	Changed:           "Alterado",
	Text:              "texto",
	Truncated:         "só os primeiros 2 MB",
	NoWrite:           "sem permissão de escrita",
	SaveCtrlS:         "Salvar (Ctrl+S)",
	EditCtrlE:         "Editar (Ctrl+E)",
	IdeFolder:         "Abrir a pasta na IDE",
	FitScreen:         "Ajustar à tela",
	ActualSize:        "Tamanho real",
	Prev:              "Anterior (←)",
	Next:              "Próxima (→)",
	DiscardTitle:      "Descartar as alterações?",
	DiscardText:       "{path} tem alterações que não foram salvas.",
	KeepEditing:       "Continuar editando",
	Discard:           "Descartar",
	Expired:           "A página expirou: recarregue-a.",
	DoneMkdir:         "Pasta criada.",
	DoneNewfile:       "Arquivo criado.",
	DoneRename:        "Renomeado.",
	DoneCopy:          "Copiado.",
	DoneMove:          "Movido.",
	DoneDelete:        "Excluído.",
	DoneExtract:       "Extraído.",
	DoneExtractN:      "Extraído ({n} arquivos).",
	DoneFetch:         "Baixado.",
	DoneFetchPath:     "Baixado: {path}",
	CloneCopied:       "URL de clone copiada: {url}",
	Saved:             "Salvo.",
	NotSent:           "Não enviados: {list}",
	SentOne:           "Arquivo enviado.",
	SentMany:          "{n} arquivos enviados.",
	SameTime:          "mesmo horário",
	Newer:             "mais novo {n}",
	Older:             "mais antigo {n}",
	Day:               "dia",
	Days:              "dias",
	SameSize:          "mesmo tamanho",
	MovedTo:           "Movido para {name}.",
	ArchiveName:       "pasta",
	Connected:         "conectado",
	Closed:            "encerrado",
	Connecting:        "conectando…",
	Tmux:              "tmux (a sessão continua ao fechar)",
	Recorded:          "gravado",
	Search:            "Buscar",
	SearchPrompt:      "Buscar no terminal:",
	FontDown:          "Diminuir a fonte",
	FontUp:            "Aumentar a fonte",
	Reopen:            "Reabrir",
	Restore:           "Voltar ao tamanho de antes",
	Expand:            "Expandir o terminal (a janela toda)",
	ConnClosed:        "conexão encerrada",
	SessionEnded:      "Sessão encerrada: {exit}",
}

// ConfigureFilesMessages registers FilesMessages.
func ConfigureFilesMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, FilesMessagesKey, FilesMessages_en_US).
		RegisterForModules(language.BrazilianPortuguese, FilesMessagesKey, FilesMessages_pt_BR)
}

// GetFilesMessages are the words of the components in ctx's language.
func GetFilesMessages(ctx context.Context) *FilesMessages {
	return i18n.MustGetModuleMessages(ctx, FilesMessagesKey, FilesMessages_en_US).(*FilesMessages)
}
