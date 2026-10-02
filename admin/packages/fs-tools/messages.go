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
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_en_US).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
