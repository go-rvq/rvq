package l10n

import (
	"context"
	"reflect"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/admin/model"

	"github.com/go-rvq/rvq/admin/reflect_utils"
	"github.com/go-rvq/rvq/web"
)

const FieldLocalizedEntries = "Localizations"

type LocaleInterface interface {
	EmbedLocale() *Locale
}

// Locale embed this struct into GROM-backend models to enable localization feature for your model
type Locale struct {
	LocaleCode string `sql:"size:20" gorm:"primaryKey;default:''"`
}

type ModelLocalizeOptions struct {
	LocalizeCallback func(ctx *web.EventContext, from, to any) (post func() (err error), err error)

	// Name is the name of the action that localizes a record — its items of
	// the menus —, "Localize" when empty; Title its title, the translation of
	// "Localize" when nil. SetAction renames it.
	Name  string
	Title func(ctx context.Context) string
	// Switches are what a localization may take along besides the records
	// chosen: each a switch of its dialogs, whose state Expand is told.
	Switches []*LocalizeSwitch
	// Expand is what a localization of mids to `to` takes, in the order it is
	// done, with the switches on in on: mids alone when nil. to may be empty:
	// the locales not chosen yet.
	Expand func(ctx *web.EventContext, mids model.IDSlice, to []string, on map[string]bool) (*Expansion, error)
	// Do localizes mids, in that order, to `to`: DoLocalizeTo each when nil.
	Do func(ctx *web.EventContext, mids model.IDSlice, to []string) error
	// Notice, when not nil, is put at the top of the dialog.
	Notice func(ctx *web.EventContext) h.HTMLComponent
	// Done is the message of n records localized: SuccessfullyLocalized when
	// nil.
	Done func(ctx context.Context, n int) string
}

// LocalizeSwitch is a switch of the dialogs of a localization (Switches).
type LocalizeSwitch struct {
	// Name is its key in Expand's on, and its field ("LocalizeSwitches.<Name>").
	Name        string
	Label, Hint func(ctx context.Context) string
	Default     bool
	// Count says the n records it takes along ("+2 ancestors"), in the
	// dialog, before they are localized.
	Count func(ctx context.Context, n int) string
}

// Expansion is what a localization takes (Expand): the records, in the
// order they are localized, and, by the name of a switch, how many it added.
type Expansion struct {
	IDs   model.IDSlice
	Added map[string]int
}

// GetLocale get model's locale
func (l *Locale) EmbedLocale() *Locale {
	return l
}

func EmbedLocale(v any) *Locale {
	return v.(LocaleInterface).EmbedLocale()
}

func IsLocalizable(obj interface{}) (isLocalizable bool) {
	_, isLocalizable = reflect_utils.GetStruct(reflect.TypeOf(obj)).(LocaleInterface)
	return
}

func IsLocalizableFromContext(ctx context.Context) (localeCode string, isLocalizable bool) {
	locale := ctx.Value(LocaleCode)
	if locale != nil {
		localeCode = locale.(string)
		isLocalizable = true
	}
	return
}
