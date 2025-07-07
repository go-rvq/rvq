package l10n

import (
	"context"
	"reflect"

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
