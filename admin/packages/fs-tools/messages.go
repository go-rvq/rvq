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
	Close                         string    `i18n:"hint='Button that closes a dialog.'"`
	TerminalOff                   string    `i18n:"label='Terminal off', hint='Why the terminal is not available.', fields=(;'%s'='the reason')"`
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
		Close:             "Close",
		TerminalOff:       "The terminal is off: %s",
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
		Close:             "Fechar",
		TerminalOff:       "O terminal está desligado: %s",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
