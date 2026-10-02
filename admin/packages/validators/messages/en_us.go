package messages

// Messages_en_US is the American English translation.
var Messages_en_US = &Messages{
	ModuleDescription: "The validators: the rules a value is checked against.",
	Validators:        "Validators",
	Validator:         "Validator",

	Name:        "Name",
	Description: "Description",
	DefaultLang: "Default language",
	Value:       "Algorithm",
	Doc:         "Documentation",
	Messages:    "Messages",

	ValueHint:     "GAD algorithm; receives `param value` and returns True/False. Empty uses the registered initial value.",
	DocHint:       "Markdown documentation explaining the algorithm. Empty uses the initial documentation.",
	MessagesHint:  "Per-language error messages, in the shape {language: {key: text}}. Empty restores the initial messages.",
	MessagesLang:  "Language",
	MessagesKey:   "Key",
	MessagesValue: "Message",
	Algorithm:     "Validation algorithm",
	EmptyResets:   "Leave empty to restore the registered initial value.",

	ErrMessagesFormat: "Messages must be in the shape {language: {key: text}} with string values.",
}
