package messages

// Messages holds the validators admin UI labels. The pt-BR and en-US variants
// must keep all fields filled (see TestMessagesCompleteness).
type Messages struct {
	ModuleDescription string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	// Validators is the menu group / resource title.
	Validators string `i18n:"hint='Name of the validators model in the plural (menu, listing title).'"`
	Validator  string `i18n:"hint='Name of the validators model in the singular (detail and form titles).'"`

	Name        string `i18n:"hint='Label of a validator\\'s name.'"`
	Description string `i18n:"hint='Label of a validator\\'s description.'"`
	DefaultLang string `i18n:"label='Default language', hint='Label of the language of a validator\\'s messages used when the user\\'s has none.'"`
	Value       string `i18n:"label='Algorithm', hint='Label of the field with a validator\\'s algorithm.'"`
	Doc         string `i18n:"label='Documentation', hint='Label of the field with a validator\\'s documentation.'"`
	Messages    string `i18n:"hint='Label of the field with a validator\\'s error messages, per language.'"`

	// hints / detail
	ValueHint     string `i18n:"label='Algorithm: hint', hint='Hint of the algorithm field: what the GAD script receives and returns.'"`
	DocHint       string `i18n:"label='Documentation: hint', hint='Hint of the documentation field (Markdown).'"`
	MessagesHint  string `i18n:"label='Messages: hint', hint='Hint of the messages field: their shape, {language: {key: text}}.'"`
	MessagesLang  string `i18n:"label='Messages: language', hint='Column with the language of a validator\\'s message.'"`
	MessagesKey   string `i18n:"label='Messages: key', hint='Column with the key of a validator\\'s message.'"`
	MessagesValue string `i18n:"label='Messages: text', hint='Column with the text of a validator\\'s message.'"`
	Algorithm     string `i18n:"label='Algorithm (detail)', hint='Title of the algorithm in a validator\\'s detail.'"`
	EmptyResets   string `i18n:"label='Empty resets', hint='Hint that leaving a field empty restores its registered initial value.'"`

	// validation errors
	ErrMessagesFormat string `i18n:"hint='Error when the messages are not in the shape {language: {key: text}}.'"`
}
