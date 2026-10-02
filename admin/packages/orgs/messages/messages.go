package messages

// Messages holds the org module UI labels. pt-BR and en-US must keep all fields
// filled.
type Messages struct {
	Organizacoes string `i18n:"label='Organizations', hint='Name of the organizations model in the plural (menu, listing title).'"`

	Organizacao          string `i18n:"label='Organization', hint='Name of the organizations model in the singular (detail and form titles).'"`
	OrganizacaoNome      string `i18n:"label='Name', hint='Label of an organization\\'s name field.'"`
	OrganizacaoDescricao string `i18n:"label='Description', hint='Label of an organization\\'s description field.'"`
	OrganizacaoOwner     string `i18n:"label='Owner', hint='Label of the user who owns an organization.'"`

	// sharing
	Compartilhar        string `i18n:"label='Share', hint='Action that shares an organization with other users.'"`
	Compartilhamento    string `i18n:"label='Sharing', hint='Title of the users an organization is shared with.'"`
	CompartilharComDica string `i18n:"label='Share: hint', hint='Hint of the field that chooses the users to share an organization with.'"`

	// selector
	SelecioneOrganizacao string `i18n:"label='Select the organization', hint='Shown when an organization must be chosen and none was.'"`
}

// Messages_pt_BR is the Brazilian Portuguese translation (module default).
var Messages_pt_BR = &Messages{
	Organizacoes: "Organizações",

	Organizacao:          "Organização",
	OrganizacaoNome:      "Nome",
	OrganizacaoDescricao: "Descrição",
	OrganizacaoOwner:     "Proprietário",

	Compartilhar:        "Compartilhar",
	Compartilhamento:    "Compartilhamento",
	CompartilharComDica: "Selecione os usuários com quem compartilhar esta organização.",

	SelecioneOrganizacao: "Selecione a organização",
}

// Messages_en_US is the American English translation.
var Messages_en_US = &Messages{
	Organizacoes: "Organizations",

	Organizacao:          "Organization",
	OrganizacaoNome:      "Name",
	OrganizacaoDescricao: "Description",
	OrganizacaoOwner:     "Owner",

	Compartilhar:        "Share",
	Compartilhamento:    "Sharing",
	CompartilharComDica: "Select the users to share this organization with.",

	SelecioneOrganizacao: "Select the organization",
}
