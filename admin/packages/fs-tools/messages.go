package fs_tools

import (
	"context"
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "rvq-admin/fs-tools"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	ModuleDescription             string    `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	FileSystem                    string    `i18n:"hint='Title of the file system tools page.'"`
	WebDavAccessTemplate          h.RawHTML `i18n:"type=html, label='WebDAV access', hint='How to reach the files through WebDAV (HTML).', fields=(;'%s'='the WebDAV URL')"`
	WebDavProtocolTitle           string    `i18n:"label='WebDAV protocol title', hint='Title of the section about the WebDAV protocol.'"`
	WebDavProtocolSoftwareExample h.RawHTML `i18n:"type=html, label='WebDAV programs', hint='Programs that access files through WebDAV, with their links (HTML).'"`
	WebDavCredentials             string    `i18n:"label='WebDAV credentials', hint='Says what the user and the password of the WebDAV are.'"`
	Folders                       string    `i18n:"hint='Title of the folders of the file system (in the WebDAV panel).'"`
	ReadOnly                      string    `i18n:"label='Read only', hint='Said of a folder that is only read.'"`
	ReadWrite                     string    `i18n:"label='Read and write', hint='Said of a folder that is read and written.'"`
	ServiceFsTitle                string    `i18n:"hint='Title of the service of the file system (its stable paths).'"`
	ServiceFsDesc                 string    `i18n:"hint='What the service of the file system serves.'"`
	RouteFsList                   string    `i18n:"hint='What the route list of the explorer gives.'"`
	RouteFsDownload               string    `i18n:"hint='What the route download of the explorer gives.'"`
	RouteFsOp                     string    `i18n:"hint='What the routes of the explorer that change the files do.'"`
	RouteFsTerminal               string    `i18n:"hint='What the route terminal gives.'"`
	RouteFsVendor                 string    `i18n:"hint='What the route vendor gives.'"`
	ServiceDavTitle               string    `i18n:"hint='Title of the service of the WebDAV.'"`
	ServiceDavDesc                string    `i18n:"hint='What the service of the WebDAV serves.'"`
	RouteDav                      string    `i18n:"hint='What the route of the WebDAV does.'"`
	ServiceGitTitle               string    `i18n:"hint='Title of the service of the repositories of git.'"`
	ServiceGitDesc                string    `i18n:"hint='What the service of the repositories serves.'"`
	RouteGit                      string    `i18n:"hint='What the route of a repository does.'"`
	TerminalOff                   string    `i18n:"label='Terminal off', hint='Why the terminal is not available.', fields=(;'%s'='the reason')"`
	KeyPermissions                string    `i18n:"label='Permissions of the access key', hint='Title of what an access key needs to reach the WebDAV.'"`
	KeyPermissionsHint            string    `i18n:"label='How the key is limited', hint='Says that a key only denies, on top of what its user may, and what it must not deny.'"`
	KeyPermRead                   string    `i18n:"label='Read by the key', hint='What the action @get allows through the WebDAV.'"`
	KeyPermWrite                  string    `i18n:"label='Write by the key', hint='What the action !write allows through the WebDAV.'"`
	KeyPermExample                string    `i18n:"label='A key only to read', hint='Introduces the line that makes a key only read the files.'"`
	TermIsolated                  string    `i18n:"label='Isolated', hint='The chip of the terminal kept apart by a sandbox.', fields=(;'%s'='the sandbox: Landlock, bubblewrap')"`
	TermNotIsolated               string    `i18n:"label='Not isolated', hint='The chip of a terminal with no sandbox.'"`
	TermIsolation                 string    `i18n:"label='Isolation', hint='Label of the line saying how the terminal is kept apart.'"`
	TermLandlock                  string    `i18n:"label='Landlock coverage', hint='What Landlock keeps apart.', fields=(;'%d'='its version (ABI)', '%s'='what it covers: files, TCP, signals')"`
	TermCovFiles                  string    `i18n:"label='Files (covered)', hint='Landlock keeps the files apart.'"`
	TermCovSignals                string    `i18n:"label='Signals (covered)', hint='Landlock keeps the signals apart.'"`
	TermBwrap                     string    `i18n:"label='bubblewrap coverage', hint='What bubblewrap keeps apart.'"`
	TermSignalsWarn               string    `i18n:"label='Signals not isolated', hint='Warns that the shell may signal the server.'"`
	TermOffWarn                   string    `i18n:"label='No isolation', hint='Warns that the shell reaches what the server does.'"`
	TermFolders                   string    `i18n:"label='Folders of the terminal', hint='Label of the folders the terminal reaches.'"`
	TermNetwork                   string    `i18n:"label='Network', hint='Label of whether the terminal reaches the network.'"`
	TermNetOn                     string    `i18n:"label='Network on', hint='The terminal reaches the network.'"`
	TermNetOff                    string    `i18n:"label='Network off', hint='The terminal does not reach the network; the server administrator turns it on.'"`
	TermLimits                    string    `i18n:"label='Limits', hint='Label of the limits of the terminal.'"`
	TermMemory                    string    `i18n:"label='Memory (limit)', hint='The limit of memory, in the limits line.'"`
	TermFileSize                  string    `i18n:"label='File (limit)', hint='The limit of the size of a file, in the limits line.'"`
	TermProcesses                 string    `i18n:"label='Processes (limit)', hint='The limit of processes, in the limits line.'"`
	TermNone                      string    `i18n:"label='No limit', hint='Said when there is no limit.'"`
	TermIdle                      string    `i18n:"label='Idle', hint='Label of when the idle terminal closes.'"`
	TermIdleCloses                string    `i18n:"label='Closes when idle', hint='When the idle terminal closes.', fields=(;'%s'='the time: 30 min')"`
	TermSession                   string    `i18n:"label='Session', hint='Label of what the session keeps.'"`
	TermSessionTemp               string    `i18n:"label='Temporary session', hint='Says HOME and /tmp are removed at the end.'"`
}

func (m *Messages) WebDavAccess(url string) h.RawHTML {
	return h.RawHTML(fmt.Sprintf(string(m.WebDavAccessTemplate), url))
}

var (
	Messages_en_US = &Messages{
		ModuleDescription: "The file system tools: browsing the files and reaching them through WebDAV.",
		FileSystem:        "File System",
		WebDavAccessTemplate: "Access the files through the WEBDAV protocol at the URL: <code class='text-primary'>%s" +
			"</code>, using your user login and password.",
		WebDavProtocolTitle: "WEBDAV Protocol",
		WebDavProtocolSoftwareExample: "<div class='mt-2'>Some programs for WEBDAV protocol access: " +
			"<a href='https://winscp.net/eng/index.php' target='_blank'>WinSCP</a> (Windows); " +
			"Dolphin and Nautilus (Linux). </div>",
		WebDavCredentials: "User: your e-mail of the admin. Password: an access key of yours (your profile › Access keys) — better than your password: it may be revoked, and limited to the files.",
		Folders:           "Folders",
		ReadOnly:          "read only",
		ReadWrite:         "read and write",
		TerminalOff:       "The terminal is off: %s",
		ServiceFsTitle:    "File system",
		ServiceFsDesc:     "The API of the explorer of the page File system, its terminal and the scripts of its components.",
		RouteFsList:       "the files of a folder (JSON)",
		RouteFsDownload:   "a file; a folder as zip",
		RouteFsOp:         "write, mkdir, rename, copy, delete, upload, upload-part, upload-done, upload-abort, extract",
		RouteFsTerminal:   "the terminal (a websocket), in the folder dir",
		RouteFsVendor:     "xterm, the terminal's library (static)",
		ServiceDavTitle:   "WebDAV",
		ServiceDavDesc:    "The files of the file system by WebDAV: a network drive.",
		RouteDav:          "WebDAV (PROPFIND, GET, PUT, MKCOL, MOVE, COPY, DELETE…): the account and its password, or an access key",
		ServiceGitTitle:   "Repositories (git)",
		ServiceGitDesc:    "The repositories of the admin (the site's files, the admin's) by git over HTTP: clone, pull, push. Their paths are of git's protocol, `<admin>/<name>.git`: the suffix tells them apart from a page's.",
		RouteGit:          "git over HTTP (clone, fetch, push): the account and its password, or an access key",
		KeyPermissions:    "Permissions of the access key",
		KeyPermissionsHint: "A key only denies, on top of what you may: a key with no line of permissions reaches what you reach. " +
			"For the WebDAV, its lines must not deny this page's resource:",
		KeyPermRead:     "read and download the files",
		KeyPermWrite:    "upload, change, move and delete",
		KeyPermExample:  "A key that only reads the files:",
		TermIsolated:    "Isolated · %s",
		TermNotIsolated: "Not isolated",
		TermIsolation:   "Isolation",
		TermLandlock:    "Landlock (ABI %d): %s",
		TermCovFiles:    "files",
		TermCovSignals:  "signals",
		TermBwrap:       "bubblewrap: mounts, processes and network of its own",
		TermSignalsWarn: "signals not isolated (Linux < 6.12): it may signal the server's processes",
		TermOffWarn:     "no sandbox: the shell reaches what the server reaches",
		TermFolders:     "Folders",
		TermNetwork:     "Network",
		TermNetOn:       "on",
		TermNetOff:      "off — to use it, ask the server's administrator",
		TermLimits:      "Limits",
		TermMemory:      "memory",
		TermFileSize:    "file",
		TermProcesses:   "processes",
		TermNone:        "none",
		TermIdle:        "Idle",
		TermIdleCloses:  "closes after %s",
		TermSession:     "Session",
		TermSessionTemp: "HOME and /tmp temporary, removed when it closes",
	}

	Messages_pt_BR = &Messages{
		ModuleDescription: "As ferramentas do sistema de arquivos: navegar pelos arquivos e acessá-los por WebDAV.",
		FileSystem:        "Sistema de Arquivos",
		WebDavAccessTemplate: "Acesse os arquivos através do protocolo WEBDAV, pela URL: <code class='text-primary'>%s" +
			"</code>, usando seu login e senha de usuário.",
		WebDavProtocolTitle: "Protocolo WEBDAV",
		WebDavProtocolSoftwareExample: "<div class='mt-2'>Alguns programas para acesso de Protocolo WEBDAV: " +
			"<a href='https://winscp.net/eng/index.php' target='_blank'>WinSCP</a> (Windows); " +
			"Dolphin and Nautilus (Linux). </div>",
		WebDavCredentials: "Usuário: o seu e-mail do admin. Senha: uma chave de acesso sua (o seu perfil › Chaves de acesso) — melhor que a sua senha: pode ser revogada, e limitada aos arquivos.",
		Folders:           "Pastas",
		ReadOnly:          "só leitura",
		ReadWrite:         "leitura e escrita",
		TerminalOff:       "O terminal está desligado: %s",
		ServiceFsTitle:    "Sistema de arquivos",
		ServiceFsDesc:     "A API do explorador da página Sistema de arquivos, o seu terminal e os scripts dos seus componentes.",
		RouteFsList:       "os arquivos de uma pasta (JSON)",
		RouteFsDownload:   "um arquivo; uma pasta em zip",
		RouteFsOp:         "write, mkdir, rename, copy, delete, upload, upload-part, upload-done, upload-abort, extract",
		RouteFsTerminal:   "o terminal (um websocket), na pasta dir",
		RouteFsVendor:     "o xterm, a biblioteca do terminal (estático)",
		ServiceDavTitle:   "WebDAV",
		ServiceDavDesc:    "Os arquivos do sistema de arquivos por WebDAV: uma unidade de rede.",
		RouteDav:          "WebDAV (PROPFIND, GET, PUT, MKCOL, MOVE, COPY, DELETE…): a conta e a sua senha, ou uma chave de acesso",
		ServiceGitTitle:   "Repositórios (git)",
		ServiceGitDesc:    "Os repositórios do admin (os arquivos do site, os do admin) por git sobre HTTP: clone, pull, push. Os seus caminhos são do protocolo do git, `<admin>/<nome>.git`: o sufixo os separa dos de uma página.",
		RouteGit:          "git sobre HTTP (clone, fetch, push): a conta e a sua senha, ou uma chave de acesso",
		KeyPermissions:    "Permissões da chave de acesso",
		KeyPermissionsHint: "Uma chave só nega, sobre o que você pode: uma chave sem nenhuma linha de permissão alcança o que você alcança. " +
			"Para o WebDAV, as linhas dela não podem negar o recurso desta página:",
		KeyPermRead:     "ler e baixar os arquivos",
		KeyPermWrite:    "enviar, mudar, mover e apagar",
		KeyPermExample:  "Uma chave que só lê os arquivos:",
		TermIsolated:    "Isolado · %s",
		TermNotIsolated: "Sem isolamento",
		TermIsolation:   "Isolamento",
		TermLandlock:    "Landlock (ABI %d): %s",
		TermCovFiles:    "arquivos",
		TermCovSignals:  "sinais",
		TermBwrap:       "bubblewrap: montagens, processos e rede próprios",
		TermSignalsWarn: "sinais não isolados (Linux < 6.12): pode sinalizar os processos do servidor",
		TermOffWarn:     "sem sandbox: o shell alcança o que o servidor alcança",
		TermFolders:     "Pastas",
		TermNetwork:     "Rede",
		TermNetOn:       "ativada",
		TermNetOff:      "desativada — para usá-la, consulte o administrador do servidor",
		TermLimits:      "Limites",
		TermMemory:      "memória",
		TermFileSize:    "arquivo",
		TermProcesses:   "processos",
		TermNone:        "nenhum",
		TermIdle:        "Inativo",
		TermIdleCloses:  "fecha após %s",
		TermSession:     "Sessão",
		TermSessionTemp: "HOME e /tmp temporários, apagados ao fechar",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
