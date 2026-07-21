package messages

// Messages holds every people UI label, action name and error text. The pt-BR
// and en-US variants must keep all fields filled (see TestMessagesCompleteness).
type Messages struct {
	// People is the menu group title of the package.
	People string

	// listing / trash
	TabAll  string
	Trash   string
	Restore string

	// Person
	Person             string
	Persons            string
	PersonName         string
	PersonAddress      string
	PersonDocumentType string
	PersonDocument     string
	PersonNotes        string
	PersonActive       string
	DocumentTypeCPF    string
	DocumentTypeCNPJ   string
	DocumentTypeOther  string

	// validation error messages
	ErrInvalidCPF          string
	ErrInvalidCNPJ         string
	ErrInvalidDocumentType string
}
