package backup

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// MessagesKey is the module of the words of the backup.
const MessagesKey i18n.ModuleKey = "admin/backup"

// GetMessages are the words of the backup in ctx's language.
func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

// ConfigureMessages registers the words of the backup.
func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModule(language.English, MessagesKey, Messages_en_US).
		RegisterForModule(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

// Messages are the words of the backup: the model's labels
// (BackupState<Field>), its statuses, its texts. {name} are placeholders.
type Messages struct {
	ModuleDescription           string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	BackupState                 string `i18n:"hint='Name of a state of the backup (singular).'"`
	BackupStates                string `i18n:"hint='The states of the backup (plural: listing, menu).'"`
	BackupState_Desc            string `i18n:"hint='What a state of the backup is.'"`
	BackupStateID               string `i18n:"label='ID', hint='Label of the field ID of a state.'"`
	BackupStateCreatedAt        string `i18n:"label='Created at', hint='Label of the field CreatedAt of a state.'"`
	BackupStateDescription      string `i18n:"label='Description', hint='Label of the field Description of a state.'"`
	BackupStateDescription_Hint string `i18n:"hint='What the form adds to the field Description.'"`
	BackupStateSources          string `i18n:"label='Sources', hint='Label of the field Sources of a state.'"`
	BackupStateSources_Hint     string `i18n:"hint='What the form adds to the field Sources.'"`
	BackupStateKeep             string `i18n:"label='Keep', hint='Label of the field Keep of a state.'"`
	BackupStateKeep_Hint        string `i18n:"hint='What the form adds to the field Keep.'"`
	BackupStateStatus           string `i18n:"label='Status', hint='Label of the field Status of a state.'"`
	BackupStateOrigin           string `i18n:"label='Origin', hint='Label of the field Origin of a state.'"`
	BackupStateCreatedByName    string `i18n:"label='Created by', hint='Label of the field CreatedByName of a state.'"`
	BackupStateStartedAt        string `i18n:"label='Started at', hint='Label of the field StartedAt of a state.'"`
	BackupStateFinishedAt       string `i18n:"label='Finished at', hint='Label of the field FinishedAt of a state.'"`
	BackupStateError            string `i18n:"label='Error', hint='Label of the field Error of a state.'"`
	BackupStateFiles            string `i18n:"label='Files', hint='Label of the field Files of a state.'"`
	BackupStateBytes            string `i18n:"label='Size', hint='Label of the field Bytes of a state.'"`
	BackupStatePartial          string `i18n:"label='Partial', hint='Label of the field Partial of a state.'"`
	BackupStateManifestSum      string `i18n:"label='Manifest', hint='Label of the field ManifestSum of a state.'"`
	BackupStateResults          string `i18n:"label='By source', hint='Label of the field Results of a state.'"`
	BackupStateDownload         string `i18n:"label='Download', hint='Label of the section Download of a state.'"`
	BackupState_Action_cancel   string `i18n:"hint='Action of a state queued or being made: cancel it.'"`
	BackupSession               string `i18n:"hint='A download session of a state (singular).'"`
	BackupSessions              string `i18n:"hint='The download sessions of a state (plural).'"`
	BackupSessionCreatedAt      string `i18n:"label='Opened', hint='Label of the field CreatedAt of a session.'"`
	BackupSessionUserName       string `i18n:"label='User', hint='Label of the field UserName of a session.'"`
	BackupSessionIP             string `i18n:"label='Address', hint='Label of the field IP of a session.'"`
	BackupSessionUserAgent      string `i18n:"label='Program', hint='Label of the field UserAgent of a session.'"`
	BackupSessionExpiresAt      string `i18n:"label='Expires', hint='Label of the field ExpiresAt of a session.'"`
	BackupSessionFiles          string `i18n:"label='Files', hint='Label of the field Files of a session.'"`
	BackupSessionBytes          string `i18n:"label='Bytes', hint='Label of the field Bytes of a session.'"`
	StatusQueued                string `i18n:"hint='Status of a state: waiting to be made.'"`
	StatusBuilding              string `i18n:"hint='Status of a state: being made.'"`
	StatusReady                 string `i18n:"hint='Status of a state: made, downloadable.'"`
	StatusFailed                string `i18n:"hint='Status of a state: its making failed.'"`
	StatusCanceled              string `i18n:"hint='Status of a state: canceled.'"`
	OriginManual                string `i18n:"hint='Origin of a state: asked by a user.'"`
	OriginScheduled             string `i18n:"hint='Origin of a state: made by the schedule.'"`
	PartialChip                 string `i18n:"hint='Chip of a partial state (some sources left out).'"`
	KeptChip                    string `i18n:"hint='Title of the mark of a kept state.'"`
	ProgressText                string `i18n:"hint='The progress of a state being made: the source, its files and size so far. {source}, {files}, {size}.'"`
	PermDownloadTitle           string `i18n:"hint='Permission of the backup: download the states (the manifest, the files).'"`
	PermDownloadDesc            string `i18n:"hint='Description of the permission of downloading the states.'"`
	JobBuildState               string `i18n:"hint='Title of the job that makes a state.'"`
	JobScheduled                string `i18n:"hint='Title of the job that makes the scheduled state.'"`
	NotReady                    string `i18n:"hint='The download of a state not ready yet. {status}.'"`
	NotDownloadable             string `i18n:"hint='The download of a state that failed or was canceled.'"`
	ManifestLink                string `i18n:"hint='Link to download the manifest of a state.'"`
	ClientURL                   string `i18n:"hint='Label of the routes the client downloads by.'"`
	ClientRoutes                string `i18n:"hint='Says the client’s routes are not the admin’s pages.'"`
	RouteLatest                 string `i18n:"hint='What the route !latest gives.'"`
	RouteSession                string `i18n:"hint='What the route !session gives.'"`
	RouteManifest               string `i18n:"hint='What the route !manifest gives.'"`
	RouteFiles                  string `i18n:"hint='What the route !files gives.'"`
	StateIDLabel                string `i18n:"hint='Label of the id of the state, for the client.'"`
	KeyNeeded                   string `i18n:"hint='Says what an access key needs to download.'"`
	PartialWarning              string `i18n:"hint='Warning of a partial state.'"`
	ColSource                   string `i18n:"hint='Column of the table by source: the source.'"`
	ColFiles                    string `i18n:"hint='Column of the table by source: the files.'"`
	ColSize                     string `i18n:"hint='Column of the table by source: the size.'"`
	ColDuration                 string `i18n:"hint='Column of the table by source: how long.'"`
	ColWarnings                 string `i18n:"hint='Column of the table by source: its warnings.'"`
	ErrInProgress               string `i18n:"hint='Error: a state being made cannot be removed.'"`
	ErrNoSource                 string `i18n:"hint='Error: no source chosen.'"`
	Interrupted                 string `i18n:"hint='Error of a state whose making the restart of the application interrupted.'"`
	Gone                        string `i18n:"hint='Error of a state ready whose files are gone (its folder removed).'"`
	Canceled                    string `i18n:"hint='Flash after a state is canceled.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription:           "The backup: states of the whole application, downloaded by a client.",
	BackupState:                 "Backup",
	BackupStates:                "Backups",
	BackupState_Desc:            "A copy of the whole application — database, site files, uploads, binary, configuration — that a client downloads into a directory where the application starts.",
	BackupStateID:               "ID",
	BackupStateCreatedAt:        "Created",
	BackupStateDescription:      "Description",
	BackupStateDescription_Hint: "Why this state (\"before the migration to 2.0\").",
	BackupStateSources:          "Sources",
	BackupStateSources_Hint:     "A state without some is partial: it does not start the application alone.",
	BackupStateKeep:             "Keep",
	BackupStateKeep_Hint:        "The retention never removes it.",
	BackupStateStatus:           "Status",
	BackupStateOrigin:           "Origin",
	BackupStateCreatedByName:    "Created by",
	BackupStateStartedAt:        "Started",
	BackupStateFinishedAt:       "Finished",
	BackupStateError:            "Error",
	BackupStateFiles:            "Files",
	BackupStateBytes:            "Size",
	BackupStatePartial:          "Partial",
	BackupStateManifestSum:      "Manifest (sha256)",
	BackupStateResults:          "By source",
	BackupStateDownload:         "Download",
	BackupState_Action_cancel:   "Cancel",
	BackupSession:               "Download session",
	BackupSessions:              "Downloads",
	BackupSessionCreatedAt:      "Opened",
	BackupSessionUserName:       "User",
	BackupSessionIP:             "Address",
	BackupSessionUserAgent:      "Program",
	BackupSessionExpiresAt:      "Expires",
	BackupSessionFiles:          "Files",
	BackupSessionBytes:          "Served",
	StatusQueued:                "queued",
	StatusBuilding:              "building",
	StatusReady:                 "ready",
	StatusFailed:                "failed",
	StatusCanceled:              "canceled",
	OriginManual:                "manual",
	OriginScheduled:             "scheduled",
	PartialChip:                 "partial",
	KeptChip:                    "kept",
	ProgressText:                "{source}: {files} files, {size}",
	PermDownloadTitle:           "Download the states",
	PermDownloadDesc:            "The manifest and the files of a state: all the application's data and secrets",
	JobBuildState:               "Make a backup state",
	JobScheduled:                "Scheduled backup",
	NotReady:                    "Not ready to download: {status}.",
	NotDownloadable:             "This state cannot be downloaded.",
	ManifestLink:                "Download the manifest",
	ClientURL:                   "Routes of the client",
	ClientRoutes:                "Out of the admin's pages and of its login: an access key opens a session (its token for the others).",
	RouteLatest:                 "the id of the newest state ready and whole",
	RouteSession:                "a session: its token (the access key)",
	RouteManifest:               "the manifest (the token)",
	RouteFiles:                  "a file of the manifest (the token; Range, If-Match)",
	StateIDLabel:                "State",
	KeyNeeded:                   "The client opens a session with an access key (your profile › Access keys) that does not deny this page's Download permission.",
	PartialWarning:              "Partial: some sources were left out — it does not start the application alone.",
	ColSource:                   "Source",
	ColFiles:                    "Files",
	ColSize:                     "Size",
	ColDuration:                 "Time",
	ColWarnings:                 "Warnings",
	ErrInProgress:               "A state being made cannot be removed: cancel it first.",
	ErrNoSource:                 "Choose at least one source.",
	Interrupted:                 "interrupted by a restart of the application",
	Gone:                        "its files are gone (its folder was removed)",
	Canceled:                    "Canceled.",
}

var Messages_pt_BR = &Messages{
	ModuleDescription:           "O backup: states da aplicação inteira, baixados por um cliente.",
	BackupState:                 "Backup",
	BackupStates:                "Backups",
	BackupState_Desc:            "Uma cópia da aplicação inteira — banco, arquivos do site, uploads, binário, configuração — que um cliente baixa num diretório onde a aplicação sobe.",
	BackupStateID:               "ID",
	BackupStateCreatedAt:        "Criado",
	BackupStateDescription:      "Descrição",
	BackupStateDescription_Hint: "Por que este state (\"antes da migração 2.0\").",
	BackupStateSources:          "Fontes",
	BackupStateSources_Hint:     "Um state sem alguma é parcial: não sobe a aplicação sozinho.",
	BackupStateKeep:             "Manter",
	BackupStateKeep_Hint:        "A retenção nunca o apaga.",
	BackupStateStatus:           "Estado",
	BackupStateOrigin:           "Origem",
	BackupStateCreatedByName:    "Criado por",
	BackupStateStartedAt:        "Início",
	BackupStateFinishedAt:       "Fim",
	BackupStateError:            "Erro",
	BackupStateFiles:            "Arquivos",
	BackupStateBytes:            "Tamanho",
	BackupStatePartial:          "Parcial",
	BackupStateManifestSum:      "Manifesto (sha256)",
	BackupStateResults:          "Por fonte",
	BackupStateDownload:         "Download",
	BackupState_Action_cancel:   "Cancelar",
	BackupSession:               "Sessão de download",
	BackupSessions:              "Downloads",
	BackupSessionCreatedAt:      "Aberta",
	BackupSessionUserName:       "Usuário",
	BackupSessionIP:             "Endereço",
	BackupSessionUserAgent:      "Programa",
	BackupSessionExpiresAt:      "Expira",
	BackupSessionFiles:          "Arquivos",
	BackupSessionBytes:          "Servido",
	StatusQueued:                "na fila",
	StatusBuilding:              "gerando",
	StatusReady:                 "pronto",
	StatusFailed:                "falhou",
	StatusCanceled:              "cancelado",
	OriginManual:                "manual",
	OriginScheduled:             "agendado",
	PartialChip:                 "parcial",
	KeptChip:                    "fixado",
	ProgressText:                "{source}: {files} arquivos, {size}",
	PermDownloadTitle:           "Baixar os states",
	PermDownloadDesc:            "O manifesto e os arquivos de um state: todos os dados e segredos da aplicação",
	JobBuildState:               "Gerar um state de backup",
	JobScheduled:                "Backup agendado",
	NotReady:                    "Ainda não pode ser baixado: {status}.",
	NotDownloadable:             "Este state não pode ser baixado.",
	ManifestLink:                "Baixar o manifesto",
	ClientURL:                   "Rotas do cliente",
	ClientRoutes:                "Fora das páginas do admin e do seu login: uma chave de acesso abre uma sessão (o token dela vale para as outras).",
	RouteLatest:                 "o id do state pronto e completo mais novo",
	RouteSession:                "uma sessão: o seu token (pela chave de acesso)",
	RouteManifest:               "o manifesto (pelo token)",
	RouteFiles:                  "um arquivo do manifesto (pelo token; Range, If-Match)",
	StateIDLabel:                "State",
	KeyNeeded:                   "O cliente abre uma sessão com uma chave de acesso (o seu perfil › Chaves de acesso) que não negue a permissão Baixar desta página.",
	PartialWarning:              "Parcial: algumas fontes ficaram de fora — não sobe a aplicação sozinho.",
	ColSource:                   "Fonte",
	ColFiles:                    "Arquivos",
	ColSize:                     "Tamanho",
	ColDuration:                 "Tempo",
	ColWarnings:                 "Avisos",
	ErrInProgress:               "Um state sendo gerado não pode ser excluído: cancele-o antes.",
	ErrNoSource:                 "Escolha ao menos uma fonte.",
	Interrupted:                 "interrompido por um reinício da aplicação",
	Gone:                        "os seus arquivos sumiram (a pasta foi removida)",
	Canceled:                    "Cancelado.",
}
