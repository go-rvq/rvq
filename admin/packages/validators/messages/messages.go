package messages

// Messages holds the validators admin UI labels. The pt-BR and en-US variants
// must keep all fields filled (see TestMessagesCompleteness).
type Messages struct {
	// Validators is the menu group / resource title.
	Validators string
	Validator  string

	Name        string
	Description string
	DefaultLang string
	Value       string
	Doc         string
	Messages    string

	// hints / detail
	ValueHint     string
	DocHint       string
	MessagesHint  string
	MessagesLang  string
	MessagesKey   string
	MessagesValue string
	Algorithm     string
	EmptyResets   string

	// validation errors
	ErrMessagesFormat string
}
