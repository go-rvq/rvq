package messages

// Messages holds the org module UI labels. pt-BR and en-US must keep all fields
// filled.
type Messages struct {
	Organizacoes string

	Organizacao          string
	OrganizacaoNome      string
	OrganizacaoDescricao string
	OrganizacaoOwner     string

	// sharing
	Compartilhar        string
	Compartilhamento    string
	CompartilharComDica string

	// selector
	SelecioneOrganizacao string
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
