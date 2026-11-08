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

	// the default organization of the user
	SetMyDefault   string `i18n:"label='Make it my default', hint='Action that makes the organization the default of the user (opened in the menu when none was chosen).'"`
	UnsetMyDefault string `i18n:"label='No longer my default', hint='Action that leaves the user with no default organization.'"`
	MyDefaultSet   string `i18n:"hint='Shown when the organization became the default of the user.'"`
	MyDefaultUnset string `i18n:"hint='Shown when the user was left with no default organization.'"`

	// selector
	SelecioneOrganizacao string `i18n:"label='Select the organization', hint='Shown when an organization must be chosen and none was.'"`
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

	SetMyDefault:   "Tornar minha padrão",
	UnsetMyDefault: "Deixar de ser minha padrão",
	MyDefaultSet:   "Agora é a sua organização padrão.",
	MyDefaultUnset: "Você não tem mais organização padrão.",

	SelecioneOrganizacao: "Selecione a organização",
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

	SetMyDefault:   "Make it my default",
	UnsetMyDefault: "No longer my default",
	MyDefaultSet:   "It is your default organization now.",
	MyDefaultUnset: "You have no default organization any more.",

	SelecioneOrganizacao: "Select the organization",
}
