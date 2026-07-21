package messages

// Messages_pt_BR is the Brazilian Portuguese translation (module default).
var Messages_pt_BR = &Messages{
	Validators: "Validadores",
	Validator:  "Validador",

	Name:        "Nome",
	Description: "Descrição",
	DefaultLang: "Idioma padrão",
	Value:       "Algoritmo",
	Doc:         "Documentação",
	Messages:    "Mensagens",

	ValueHint:     "Algoritmo em GAD; recebe `param value` e retorna True/False. Vazio usa o valor inicial registrado.",
	DocHint:       "Documentação em Markdown explicando o algoritmo. Vazio usa a documentação inicial.",
	MessagesHint:  "Mensagens de erro por idioma, no formato {idioma: {chave: texto}}. Vazio restaura as mensagens iniciais.",
	MessagesLang:  "Idioma",
	MessagesKey:   "Chave",
	MessagesValue: "Mensagem",
	Algorithm:     "Algoritmo de validação",
	EmptyResets:   "Deixe em branco para restaurar o valor inicial registrado.",

	ErrMessagesFormat: "As mensagens devem estar no formato {idioma: {chave: texto}} com valores de texto.",
}
