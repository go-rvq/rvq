package schemaform

import (
	"context"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/x/i18n"
)

// MessagesKey is the i18n module of the schemaform's messages: an application
// registers the languages it offers for it (i18n.Builder.RegisterForModule —
// Messages_pt_BR, Messages_en_US), English when none is.
const MessagesKey i18n.ModuleKey = "admin/presets/fields/schemaform"

// Messages are what the schemaform tells the one filling a form, or reading
// one: the errors of a value posted, the buttons, and — when the schema asks
// for what cannot be drawn — why a field is not there.
type Messages struct {
	// NoSchema: a field drawn without a schema.
	NoSchema string
	// FormWithoutSchema: a form field with no schema; %q is the field.
	FormWithoutSchema string
	// ItemsUnavailable: the values a field may hold could not be fetched; %q
	// the field, %v why.
	ItemsUnavailable string
	// NoItems: a field that is a closed list has nothing to choose; %q the
	// field.
	NoItems string
	// ListTypeUnknown: a list's items are of a type with no component; %q the
	// type, %v the ones there are.
	ListTypeUnknown string
	// FieldTypeUnknown: a field of a type with no component; %q the field,
	// %q the type, %v the ones there are.
	FieldTypeUnknown string

	// ChooseValue: a required choice was left empty.
	ChooseValue string
	// NotAmongItems: a value the list does not offer; %q the value, %s the
	// ones it does.
	NotAmongItems string
	// ItemsUnavailableOnSave: the values could not be fetched to check the one
	// posted (why follows it).
	ItemsUnavailableOnSave string

	// Add is the button that adds an item to a list.
	Add string
	// ShowWhole is the hint of a listing cell cut with "…".
	ShowWhole string
	// OrderAsc, OrderDesc and OrderNone are the directions of a field of an
	// order: ascending, descending, not ordered.
	OrderAsc  string
	OrderDesc string
	OrderNone string
}

var Messages_en_US = &Messages{
	NoSchema:               "schemaform: the field has no schema",
	FormWithoutSchema:      "schemaform: the field %q is a form with no schema",
	ItemsUnavailable:       "the values of %q could not be fetched: %v",
	NoItems:                "the field %q has no values to choose from",
	ListTypeUnknown:        "schemaform: the list asks for the type %q, which has no component and is no enum (there are: %v)",
	FieldTypeUnknown:       "schemaform: the field %q asks for the type %q, which has no component and is no enum (there are: %v)",
	ChooseValue:            "choose a value",
	NotAmongItems:          "%q is not one of the values it may hold (%s)",
	ItemsUnavailableOnSave: "the values it may hold could not be fetched",
	Add:                    "Add",
	ShowWhole:              "click to see it whole",
	OrderAsc:               "ASC",
	OrderDesc:              "DESC",
	OrderNone:              "—",
}

var Messages_pt_BR = &Messages{
	NoSchema:               "schemaform: o campo não tem schema",
	FormWithoutSchema:      "schemaform: o campo %q é um form sem schema",
	ItemsUnavailable:       "os valores de %q não puderam ser obtidos: %v",
	NoItems:                "o campo %q não tem valores para escolher",
	ListTypeUnknown:        "schemaform: a lista pede o type %q, que não tem componente registrado nem é enum (há: %v)",
	FieldTypeUnknown:       "schemaform: o campo %q pede o type %q, que não tem componente registrado nem é enum (há: %v)",
	ChooseValue:            "escolha um valor",
	NotAmongItems:          "%q não está entre os valores possíveis (%s)",
	ItemsUnavailableOnSave: "os valores possíveis não puderam ser obtidos",
	Add:                    "Adicionar",
	ShowWhole:              "clique para ver inteiro",
	OrderAsc:               "ASC",
	OrderDesc:              "DESC",
	OrderNone:              "—",
}

// GetMessages are the schemaform's messages in the language of ctx.
func GetMessages(ctx context.Context) *Messages {
	if ctx == nil {
		ctx = context.Background()
	}
	return i18n.MustGetModuleMessages(ctx, MessagesKey, Messages_en_US).(*Messages)
}

// messagesOf are the messages in the language of the request being answered.
func messagesOf(ev *web.EventContext) *Messages {
	if ev == nil || ev.R == nil {
		// no request: the default language's
		return GetMessages(context.Background())
	}
	return GetMessages(ev.R.Context())
}

// Messages are the schemaform's messages in the language of the request.
func (c *Context) Messages() *Messages { return messagesOf(c.Event) }
