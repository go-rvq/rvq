package accesskeys

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
	"golang.org/x/text/language"
)

// MessagesKey is the module of the words of the access keys.
const MessagesKey i18n.ModuleKey = "admin/access_keys"

func GetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

// ConfigureMessages registers the words of the access keys.
func ConfigureMessages(b *i18n.Builder) {
	b.RegisterForModule(language.English, MessagesKey, Messages_en_US).
		RegisterForModule(language.BrazilianPortuguese, MessagesKey, Messages_pt_BR)
}

type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`

	AccessKey        string `i18n:"hint='Name of an access key (singular).'"`
	AccessKeys       string `i18n:"hint='The access keys of a user (plural: listing, menu).'"`
	AccessKey_Desc   string `i18n:"hint='What an access key is.'"`
	MyAccessKey      string `i18n:"hint='One of the keys of the user of the request (singular).'"`
	MyAccessKeys     string `i18n:"hint='The keys of the user of the request: the menu item My keys.'"`
	MyAccessKey_Desc string `i18n:"hint='What My keys is.'"`
	AccessKeyUse     string `i18n:"hint='A request made by an access key (singular).'"`
	AccessKeyUses    string `i18n:"hint='The history of an access key: its requests.'"`

	AccessKeyName               string `i18n:"label='Name', hint='Label of the field Name of a key.'"`
	AccessKeyName_Desc          string `i18n:"hint='Description of the field Name of a key.'"`
	AccessKeyName_Hint          string `i18n:"hint='What the form of a key adds to the description of the field Name.'"`
	AccessKeyDescription        string `i18n:"label='Description', hint='Label of the field Description of a key.'"`
	AccessKeyDescription_Desc   string `i18n:"hint='Description of the field Description of a key.'"`
	AccessKeyDescription_Hint   string `i18n:"hint='What the form of a key adds to the description of the field Description.'"`
	AccessKeyPrefix             string `i18n:"label='Prefix', hint='Label of the field Prefix of a key.'"`
	AccessKeyPrefix_Desc        string `i18n:"hint='Description of the field Prefix of a key.'"`
	AccessKeyEnabled            string `i18n:"label='Enabled', hint='Label of the field Enabled of a key.'"`
	AccessKeyEnabled_Desc       string `i18n:"hint='Description of the field Enabled of a key.'"`
	AccessKeyEnabled_Hint       string `i18n:"hint='What the form of a key adds to the description of the field Enabled.'"`
	AccessKeyExpiresAt          string `i18n:"label='Expires at', hint='Label of the field ExpiresAt of a key.'"`
	AccessKeyExpiresAt_Desc     string `i18n:"hint='Description of the field ExpiresAt of a key.'"`
	AccessKeyExpiresAt_Hint     string `i18n:"hint='What the form of a key adds to the description of the field ExpiresAt.'"`
	AccessKeyLastUsedAt         string `i18n:"label='Last used', hint='Label of the field LastUsedAt of a key.'"`
	AccessKeyLastUsedIP         string `i18n:"label='Last address', hint='Label of the field LastUsedIP of a key.'"`
	AccessKeyCreatedAt          string `i18n:"label='Created at', hint='Label of the field CreatedAt of a key.'"`
	AccessKeyPermissions        string `i18n:"label='Permissions', hint='Label of the field Permissions of a key.'"`
	AccessKeyPermissions_Desc   string `i18n:"hint='Description of the field Permissions of a key.'"`
	MyAccessKeyName             string `i18n:"label='Name', hint='Label of the field Name of one of my keys.'"`
	MyAccessKeyName_Desc        string `i18n:"hint='Description of the field Name of one of my keys.'"`
	MyAccessKeyName_Hint        string `i18n:"hint='What the form of one of my keys adds to the description of the field Name.'"`
	MyAccessKeyDescription      string `i18n:"label='Description', hint='Label of the field Description of one of my keys.'"`
	MyAccessKeyDescription_Desc string `i18n:"hint='Description of the field Description of one of my keys.'"`
	MyAccessKeyDescription_Hint string `i18n:"hint='What the form of one of my keys adds to the description of the field Description.'"`
	MyAccessKeyPrefix           string `i18n:"label='Prefix', hint='Label of the field Prefix of one of my keys.'"`
	MyAccessKeyPrefix_Desc      string `i18n:"hint='Description of the field Prefix of one of my keys.'"`
	MyAccessKeyEnabled          string `i18n:"label='Enabled', hint='Label of the field Enabled of one of my keys.'"`
	MyAccessKeyEnabled_Desc     string `i18n:"hint='Description of the field Enabled of one of my keys.'"`
	MyAccessKeyEnabled_Hint     string `i18n:"hint='What the form of one of my keys adds to the description of the field Enabled.'"`
	MyAccessKeyExpiresAt        string `i18n:"label='Expires at', hint='Label of the field ExpiresAt of one of my keys.'"`
	MyAccessKeyExpiresAt_Desc   string `i18n:"hint='Description of the field ExpiresAt of one of my keys.'"`
	MyAccessKeyExpiresAt_Hint   string `i18n:"hint='What the form of one of my keys adds to the description of the field ExpiresAt.'"`
	MyAccessKeyLastUsedAt       string `i18n:"label='Last used', hint='Label of the field LastUsedAt of one of my keys.'"`
	MyAccessKeyLastUsedIP       string `i18n:"label='Last address', hint='Label of the field LastUsedIP of one of my keys.'"`
	MyAccessKeyCreatedAt        string `i18n:"label='Created at', hint='Label of the field CreatedAt of one of my keys.'"`
	MyAccessKeyPermissions      string `i18n:"label='Permissions', hint='Label of the field Permissions of one of my keys.'"`
	MyAccessKeyPermissions_Desc string `i18n:"hint='Description of the field Permissions of one of my keys.'"`

	AccessKeyUseCreatedAt string `i18n:"label='When', hint='Label of when a request by a key was made.'"`
	AccessKeyUseKind      string `i18n:"label='Kind', hint='Label of what a request by a key reached (git, webdav, api).'"`
	AccessKeyUseMethod    string `i18n:"label='Method', hint='Label of the HTTP method of a request by a key.'"`
	AccessKeyUsePath      string `i18n:"label='Path', hint='Label of the path of a request by a key.'"`
	AccessKeyUseStatus    string `i18n:"label='Status', hint='Label of the HTTP status a request by a key was answered.'"`
	AccessKeyUseIP        string `i18n:"label='Address', hint='Label of the address of a request by a key.'"`
	AccessKeyUsePlace     string `i18n:"label='Place', hint='Label of where the address of a request by a key is.'"`
	AccessKeyUseUserAgent string `i18n:"label='Program', hint='Label of the program (user agent) of a request by a key.'"`

	DefaultDbpolicyActions   string `i18n:"label='Actions', hint='Label of the actions a permission of a key applies to.'"`
	DefaultDbpolicyResources string `i18n:"label='Resources', hint='Label of the resources a permission of a key applies to.'"`

	PermissionsHelp string `i18n:"hint='Title of the help button of the permissions of a key.'"`
	CodeCreated     string `i18n:"hint='Shown once a key was made: its code, to be copied now.', fields=(;'%s'='the code')"`
	CodeOnce        string `i18n:"hint='Says the code is shown only now.'"`
	Code            string `i18n:"hint='Label of the code of a key just made.'"`
	ErrNameRequired string `i18n:"hint='Error: the key has no name.'"`
	NoPermissions   string `i18n:"hint='Said of a key with no permission: it does nothing.'"`
	ErrTooLong      string `i18n:"hint='Error: the key expires too far ahead.', fields=(;'%d'='the most days')"`
}

var Messages_en_US = &Messages{
	ModuleDescription: "The access keys of the users: codes for automations (git, WebDAV, scripts).",

	AccessKey:        "Access key",
	AccessKeys:       "Access keys",
	AccessKey_Desc:   "A code an automation — git, WebDAV, a script — authenticates with in the user's name. Its permissions only restrict: it does what the user may and it allows.",
	MyAccessKey:      "My access key",
	MyAccessKeys:     "My access keys",
	MyAccessKey_Desc: "Your access keys: codes for your automations — git, WebDAV, scripts —, each doing only what you may and it allows.",
	AccessKeyUse:     "Access",
	AccessKeyUses:    "History",

	AccessKeyName:               "Name",
	AccessKeyName_Desc:          "What the key is for.",
	AccessKeyName_Hint:          "",
	AccessKeyDescription:        "Description",
	AccessKeyDescription_Desc:   "More about the key: where it is used, by whom.",
	AccessKeyDescription_Hint:   "",
	AccessKeyPrefix:             "Prefix",
	AccessKeyPrefix_Desc:        "The public part of the code, to tell the key.",
	AccessKeyEnabled:            "Enabled",
	AccessKeyEnabled_Desc:       "Whether the key works: untick to stop it at once.",
	AccessKeyEnabled_Hint:       "",
	AccessKeyExpiresAt:          "Expires at",
	AccessKeyExpiresAt_Desc:     "When the key stops working.",
	AccessKeyExpiresAt_Hint:     "Empty: 90 days from now; 1 year at most.",
	AccessKeyLastUsedAt:         "Last used",
	AccessKeyLastUsedIP:         "Last address",
	AccessKeyCreatedAt:          "Created at",
	AccessKeyPermissions:        "Permissions",
	AccessKeyPermissions_Desc:   "What the key may do — only of what the user may. With none, it does nothing.",
	MyAccessKeyName:             "Name",
	MyAccessKeyName_Desc:        "What the key is for.",
	MyAccessKeyName_Hint:        "",
	MyAccessKeyDescription:      "Description",
	MyAccessKeyDescription_Desc: "More about the key: where it is used, by whom.",
	MyAccessKeyDescription_Hint: "",
	MyAccessKeyPrefix:           "Prefix",
	MyAccessKeyPrefix_Desc:      "The public part of the code, to tell the key.",
	MyAccessKeyEnabled:          "Enabled",
	MyAccessKeyEnabled_Desc:     "Whether the key works: untick to stop it at once.",
	MyAccessKeyEnabled_Hint:     "",
	MyAccessKeyExpiresAt:        "Expires at",
	MyAccessKeyExpiresAt_Desc:   "When the key stops working.",
	MyAccessKeyExpiresAt_Hint:   "Empty: 90 days from now; 1 year at most.",
	MyAccessKeyLastUsedAt:       "Last used",
	MyAccessKeyLastUsedIP:       "Last address",
	MyAccessKeyCreatedAt:        "Created at",
	MyAccessKeyPermissions:      "Permissions",
	MyAccessKeyPermissions_Desc: "What the key may do — only of what the user may. With none, it does nothing.",

	AccessKeyUseCreatedAt: "When",
	AccessKeyUseKind:      "Kind",
	AccessKeyUseMethod:    "Method",
	AccessKeyUsePath:      "Path",
	AccessKeyUseStatus:    "Status",
	AccessKeyUseIP:        "Address",
	AccessKeyUsePlace:     "Place",
	AccessKeyUseUserAgent: "Program",

	DefaultDbpolicyActions:   "Actions",
	DefaultDbpolicyResources: "Resources",

	PermissionsHelp: "Help: policies and permissions",
	CodeCreated:     "Access key created. Copy its code now — it is shown only this time: %s",
	CodeOnce:        "This is the only time the code is shown. Copy it and keep it safe: it acts in your name.",
	Code:            "Code",
	ErrNameRequired: "Give the key a name.",
	NoPermissions:   "None: the key does nothing.",
	ErrTooLong:      "A key lasts at most %d days.",
}

var Messages_pt_BR = &Messages{
	ModuleDescription: "As chaves de acesso dos usuários: códigos para automações (git, WebDAV, scripts).",

	AccessKey:        "Chave de acesso",
	AccessKeys:       "Chaves de acesso",
	AccessKey_Desc:   "Um código com que uma automação — git, WebDAV, um script — se autentica em nome do usuário. As permissões dela só restringem: ela faz o que o usuário pode e ela permite.",
	MyAccessKey:      "Minha chave de acesso",
	MyAccessKeys:     "Minhas chaves",
	MyAccessKey_Desc: "As suas chaves de acesso: códigos para as suas automações — git, WebDAV, scripts —, cada uma fazendo só o que você pode e ela permite.",
	AccessKeyUse:     "Acesso",
	AccessKeyUses:    "Histórico",

	AccessKeyName:               "Nome",
	AccessKeyName_Desc:          "Para que a chave serve.",
	AccessKeyName_Hint:          "",
	AccessKeyDescription:        "Descrição",
	AccessKeyDescription_Desc:   "Mais sobre a chave: onde é usada, por quem.",
	AccessKeyDescription_Hint:   "",
	AccessKeyPrefix:             "Prefixo",
	AccessKeyPrefix_Desc:        "A parte pública do código, para reconhecer a chave.",
	AccessKeyEnabled:            "Ativada",
	AccessKeyEnabled_Desc:       "Se a chave funciona: desmarque para pará-la na hora.",
	AccessKeyEnabled_Hint:       "",
	AccessKeyExpiresAt:          "Expira em",
	AccessKeyExpiresAt_Desc:     "Quando a chave deixa de funcionar.",
	AccessKeyExpiresAt_Hint:     "Vazio: 90 dias a partir de agora; no máximo 1 ano.",
	AccessKeyLastUsedAt:         "Último uso",
	AccessKeyLastUsedIP:         "Último endereço",
	AccessKeyCreatedAt:          "Criada em",
	AccessKeyPermissions:        "Permissões",
	AccessKeyPermissions_Desc:   "O que a chave pode fazer — só do que o usuário pode. Sem nenhuma, ela não faz nada.",
	MyAccessKeyName:             "Nome",
	MyAccessKeyName_Desc:        "Para que a chave serve.",
	MyAccessKeyName_Hint:        "",
	MyAccessKeyDescription:      "Descrição",
	MyAccessKeyDescription_Desc: "Mais sobre a chave: onde é usada, por quem.",
	MyAccessKeyDescription_Hint: "",
	MyAccessKeyPrefix:           "Prefixo",
	MyAccessKeyPrefix_Desc:      "A parte pública do código, para reconhecer a chave.",
	MyAccessKeyEnabled:          "Ativada",
	MyAccessKeyEnabled_Desc:     "Se a chave funciona: desmarque para pará-la na hora.",
	MyAccessKeyEnabled_Hint:     "",
	MyAccessKeyExpiresAt:        "Expira em",
	MyAccessKeyExpiresAt_Desc:   "Quando a chave deixa de funcionar.",
	MyAccessKeyExpiresAt_Hint:   "Vazio: 90 dias a partir de agora; no máximo 1 ano.",
	MyAccessKeyLastUsedAt:       "Último uso",
	MyAccessKeyLastUsedIP:       "Último endereço",
	MyAccessKeyCreatedAt:        "Criada em",
	MyAccessKeyPermissions:      "Permissões",
	MyAccessKeyPermissions_Desc: "O que a chave pode fazer — só do que o usuário pode. Sem nenhuma, ela não faz nada.",

	AccessKeyUseCreatedAt: "Quando",
	AccessKeyUseKind:      "Tipo",
	AccessKeyUseMethod:    "Método",
	AccessKeyUsePath:      "Caminho",
	AccessKeyUseStatus:    "Resultado",
	AccessKeyUseIP:        "Endereço",
	AccessKeyUsePlace:     "Lugar",
	AccessKeyUseUserAgent: "Programa",

	DefaultDbpolicyActions:   "Ações",
	DefaultDbpolicyResources: "Recursos",

	PermissionsHelp: "Ajuda: políticas e permissões",
	CodeCreated:     "Chave de acesso criada. Copie o código agora — ele só aparece desta vez: %s",
	CodeOnce:        "É a única vez que o código aparece. Copie e guarde num lugar seguro: ele age em seu nome.",
	Code:            "Código",
	ErrNameRequired: "Dê um nome à chave.",
	NoPermissions:   "Nenhuma: a chave não faz nada.",
	ErrTooLong:      "Uma chave dura no máximo %d dias.",
}
