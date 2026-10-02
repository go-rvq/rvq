package messages

// Messages holds every people UI label, action name and error text. The pt-BR
// and en-US variants must keep all fields filled (see TestMessagesCompleteness).
type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	// People is the menu group title of the package.
	People string `i18n:"hint='Title of the people section in the menu.'"`

	// listing / trash
	TabAll  string `i18n:"label='Tab: all', hint='Tab of the people listing that shows everyone.'"`
	Trash   string `i18n:"hint='Tab of the people listing that shows the deleted records.'"`
	Restore string `i18n:"hint='Action that restores deleted people.'"`

	// Person
	Person             string `i18n:"hint='Name of the people model in the singular (detail and form titles).'"`
	Persons            string `i18n:"label='People', hint='Name of the people model in the plural (listing title).'"`
	PersonName         string `i18n:"label='Name', hint='Label of a person\\'s name.'"`
	PersonAddress      string `i18n:"label='Address', hint='Label of a person\\'s address.'"`
	PersonDocumentType string `i18n:"label='Document type', hint='Label of the kind of a person\\'s document.'"`
	PersonDocument     string `i18n:"label='Document', hint='Label of a person\\'s document number.'"`
	PersonNotes        string `i18n:"label='Notes', hint='Label of the notes about a person.'"`
	PersonActive       string `i18n:"label='Active', hint='Label of whether a person is active.'"`
	DocumentTypeCPF    string `i18n:"label='CPF', hint='Document type: the Brazilian individual taxpayer number.'"`
	DocumentTypeCNPJ   string `i18n:"label='CNPJ', hint='Document type: the Brazilian company taxpayer number.'"`
	DocumentTypeOther  string `i18n:"hint='Document type: any other document.'"`

	// validation error messages
	ErrInvalidCPF          string `i18n:"label='Invalid CPF', hint='Error of a CPF that is not valid.'"`
	ErrInvalidCNPJ         string `i18n:"label='Invalid CNPJ', hint='Error of a CNPJ that is not valid.'"`
	ErrInvalidDocumentType string `i18n:"hint='Error of a document type that is not one of the known ones.'"`
}
