package messages

// Messages holds the org module UI labels. pt-BR and en-US must keep all fields
// filled.
type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Organizacoes      string `i18n:"label='Organizations', hint='Name of the organizations model in the plural (menu, listing title).'"`

	Organizacao          string `i18n:"label='Organization', hint='Name of the organizations model in the singular (detail and form titles).'"`
	OrganizacaoNome      string `i18n:"label='Name', hint='Label of an organization\\'s name field.'"`
	OrganizacaoDescricao string `i18n:"label='Description', hint='Label of an organization\\'s description field.'"`
	OrganizacaoOwner     string `i18n:"label='Owner', hint='Label of the user who owns an organization.'"`
	OrganizacaoMyDefault string `i18n:"label='My default', hint='Column of the organizations: the one that is the user\\'s default (opened in the menu when none was chosen).'"`

	SharedBy string `i18n:"hint='Beside the name of an organization of another, shared with the user: before the name of its owner.'"`

	// the default organization of the user
	SetMyDefault   string `i18n:"label='Make it my default', hint='Action that makes the organization the default of the user (opened in the menu when none was chosen).'"`
	UnsetMyDefault string `i18n:"label='No longer my default', hint='Action that leaves the user with no default organization.'"`
	MyDefaultSet   string `i18n:"hint='Shown when the organization became the default of the user.'"`
	MyDefaultUnset string `i18n:"hint='Shown when the user was left with no default organization.'"`

	// selector
	SelecioneOrganizacao string `i18n:"label='Select the organization', hint='Shown when an organization must be chosen and none was.'"`

	// the panel of an organization and its widgets
	Painel          string `i18n:"label='Dashboard', hint='Title of the page of an organization with its widgets (its dashboard).'"`
	WidgetsGroup    string `i18n:"label='Organization', hint='Title of the group of the widgets of an organization (its members, what waits, its activity).'"`
	Membros         string `i18n:"label='Members', hint='Title of the widget of who an organization is shared with, and their roles.'"`
	MembrosCaption  string `i18n:"label='members', hint='Under the number of the members of an organization.'"`
	SemMembros      string `i18n:"hint='Shown when an organization is shared with nobody.'"`
	Pendentes       string `i18n:"label='Pending requests and invites', hint='Title of the widget of the sharings of an organization waiting for its owner or for whom was invited.'"`
	PendentesNenhum string `i18n:"hint='Shown when nothing of an organization waits.'"`
	PedidoDono      string `i18n:"hint='Beside a request: it waits for the owner to confirm it by e-mail.'"`
	ConviteAceite   string `i18n:"hint='Beside an invite: it waits for whom was invited to accept it.'"`
	Atividade       string `i18n:"label='Recent activity', hint='Title of the widget of the latest events of the sharing of an organization.'"`
	SemAtividade    string `i18n:"hint='Shown when an organization has no events of sharing.'"`
	VerTudo         string `i18n:"label='See all', hint='Link of a widget to the page with everything it summarizes.'"`
}

// Messages_pt_BR is the Brazilian Portuguese translation (module default).
var Messages_pt_BR = &Messages{
	ModuleDescription: "As organizações.",
	Organizacoes:      "Organizações",

	Organizacao:          "Organização",
	OrganizacaoNome:      "Nome",
	OrganizacaoDescricao: "Descrição",
	OrganizacaoOwner:     "Proprietário",
	OrganizacaoMyDefault: "Minha padrão",

	SharedBy: "compartilhada por",

	SetMyDefault:   "Tornar minha padrão",
	UnsetMyDefault: "Deixar de ser minha padrão",
	MyDefaultSet:   "Agora é a sua organização padrão.",
	MyDefaultUnset: "Você não tem mais organização padrão.",

	SelecioneOrganizacao: "Selecione a organização",

	Painel:          "Painel",
	WidgetsGroup:    "Organização",
	Membros:         "Membros",
	MembrosCaption:  "membros",
	SemMembros:      "A organização não está compartilhada com ninguém.",
	Pendentes:       "Pedidos e convites pendentes",
	PendentesNenhum: "Nada aguardando.",
	PedidoDono:      "aguarda o proprietário",
	ConviteAceite:   "aguarda o aceite",
	Atividade:       "Atividade recente",
	SemAtividade:    "Nenhuma atividade de compartilhamento.",
	VerTudo:         "Ver tudo",
}

// Messages_en_US is the American English translation.
var Messages_en_US = &Messages{
	ModuleDescription: "The organizations.",
	Organizacoes:      "Organizations",

	Organizacao:          "Organization",
	OrganizacaoNome:      "Name",
	OrganizacaoDescricao: "Description",
	OrganizacaoOwner:     "Owner",
	OrganizacaoMyDefault: "My default",

	SharedBy: "shared by",

	SetMyDefault:   "Make it my default",
	UnsetMyDefault: "No longer my default",
	MyDefaultSet:   "It is your default organization now.",
	MyDefaultUnset: "You have no default organization any more.",

	SelecioneOrganizacao: "Select the organization",

	Painel:          "Dashboard",
	WidgetsGroup:    "Organization",
	Membros:         "Members",
	MembrosCaption:  "members",
	SemMembros:      "The organization is shared with nobody.",
	Pendentes:       "Pending requests and invites",
	PendentesNenhum: "Nothing waiting.",
	PedidoDono:      "waits for the owner",
	ConviteAceite:   "waits to be accepted",
	Atividade:       "Recent activity",
	SemAtividade:    "No sharing activity.",
	VerTudo:         "See all",
}
