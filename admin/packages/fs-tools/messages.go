package fs_tools

import (
	"context"
	"fmt"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

const MessagesKey i18n.ModuleKey = "rqv-admin/fs-tools"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

type Messages struct {
	FileSystem                    string
	WebDavAccessTemplate          h.RawHTML
	WebDavProtocolTitle           string
	WebDavProtocolSoftwareExample h.RawHTML
}

func (m *Messages) WebDavAccess(url string) h.RawHTML {
	return h.RawHTML(fmt.Sprintf(string(m.WebDavAccessTemplate), url))
}

var (
	Messages_en_US = &Messages{
		FileSystem:          "File System",
		WebDavProtocolTitle: "WEBDAV Protocol",
	}

	Messages_pt_BR = &Messages{
		FileSystem: "Sistema de Arquivos",
		WebDavAccessTemplate: "Acesse os arquivos através do protocolo WEBDAV, pela URL: <code class='text-primary'>%s" +
			"</code>, usando seu login e senha de usuário.",
		WebDavProtocolTitle: "Protocolo WEBDAV",
		WebDavProtocolSoftwareExample: "<div class='mt-2'>Alguns programas para acesso de Protocolo WEBDAV: " +
			"<a href='https://winscp.net/eng/index.php' target='_blank'>WinSCP</a> (Windows); " +
			"Dolphin and Nautilus (Linux). </div>",
	}
)

func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModules(language.English, MessagesKey, Messages_pt_BR).
		RegisterForModules(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}
