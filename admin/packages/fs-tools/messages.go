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
	TerminalOff                   string    `i18n:"label='Terminal off', hint='Why the terminal is not available.', fields=(;'%s'='the reason')"`
	KeyPermissions                string    `i18n:"label='Permissions of the access key', hint='Title of what an access key needs to reach the WebDAV.'"`
	KeyPermissionsHint            string    `i18n:"label='How the key is limited', hint='Says that a key only denies, on top of what its user may, and what it must not deny.'"`
	KeyPermRead                   string    `i18n:"label='Read by the key', hint='What the action @get allows through the WebDAV.'"`
	KeyPermWrite                  string    `i18n:"label='Write by the key', hint='What the action !write allows through the WebDAV.'"`
	KeyPermExample                string    `i18n:"label='A key only to read', hint='Introduces the line that makes a key only read the files.'"`
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
		KeyPermissions:    "Permissions of the access key",
		KeyPermissionsHint: "A key only denies, on top of what you may: a key with no line of permissions reaches what you reach. " +
			"For the WebDAV, its lines must not deny this page's resource:",
		KeyPermRead:    "read and download the files",
		KeyPermWrite:   "upload, change, move and delete",
		KeyPermExample: "A key that only reads the files:",
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
		KeyPermissions:    "Permissões da chave de acesso",
		KeyPermissionsHint: "Uma chave só nega, sobre o que você pode: uma chave sem nenhuma linha de permissão alcança o que você alcança. " +
			"Para o WebDAV, as linhas dela não podem negar o recurso desta página:",
		KeyPermRead:    "ler e baixar os arquivos",
		KeyPermWrite:   "enviar, mudar, mover e apagar",
		KeyPermExample: "Uma chave que só lê os arquivos:",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
