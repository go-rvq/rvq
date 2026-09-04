package l10n

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
)

const I18nLocalizeKey i18n.ModuleKey = "I18nLocalizeKey"

type Messages struct {
	Localize                         string
	LocalizeFrom                     string
	LocalizeTo                       string
	CurrentLocalizations             string
	Actions                          string
	Localizations                    string
	SuccessfullyLocalized            string
	Location                         string
	Colon                            string
	International                    string
	China                            string
	Japan                            string
	ErrDeleteInternationalizedRecord string
	ChangeLocale                     string
	SuccessfullyChangedLocale        string
	ErrChangeLocaleEmpty             string
	ErrChangeLocaleUnavailable       string
}

var Messages_en_US = &Messages{
	Localize:                         "Localize",
	LocalizeFrom:                     "From",
	LocalizeTo:                       "To",
	CurrentLocalizations:             "Current Localizations",
	Actions:                          "Actions",
	Localizations:                    "Localizations",
	SuccessfullyLocalized:            "Successfully Localized",
	Location:                         "Location",
	Colon:                            ":",
	International:                    "International",
	China:                            "China",
	Japan:                            "Japan",
	ErrDeleteInternationalizedRecord: "It is not possible to delete the standard language record when it has is internationalized.",
	ChangeLocale:                     "Change Location",
	SuccessfullyChangedLocale:        "Location changed",
	ErrChangeLocaleEmpty:             "Choose the location to move this record to.",
	ErrChangeLocaleUnavailable:       "This location is no longer available: the record already exists there.",
}

var Messages_zh_CN = &Messages{
	Localize:                         "本地化",
	LocalizeFrom:                     "从",
	LocalizeTo:                       "到",
	CurrentLocalizations:             Messages_en_US.CurrentLocalizations,
	Actions:                          Messages_en_US.Actions,
	Localizations:                    Messages_en_US.Localizations,
	SuccessfullyLocalized:            "本地化成功",
	Location:                         "地区",
	Colon:                            "：",
	International:                    "全球",
	China:                            "中国",
	Japan:                            "日本",
	ErrDeleteInternationalizedRecord: Messages_en_US.ErrDeleteInternationalizedRecord,
	ChangeLocale:                     Messages_en_US.ChangeLocale,
	SuccessfullyChangedLocale:        Messages_en_US.SuccessfullyChangedLocale,
	ErrChangeLocaleEmpty:             Messages_en_US.ErrChangeLocaleEmpty,
	ErrChangeLocaleUnavailable:       Messages_en_US.ErrChangeLocaleUnavailable,
}

var Messages_ja_JP = &Messages{
	Localize:                         "ローカライズ",
	LocalizeFrom:                     "から",
	LocalizeTo:                       "に",
	CurrentLocalizations:             Messages_en_US.CurrentLocalizations,
	Actions:                          Messages_en_US.Actions,
	Localizations:                    Messages_en_US.Localizations,
	SuccessfullyLocalized:            "ローカライズに成功しました",
	Location:                         "場所",
	Colon:                            ":",
	International:                    "インターナショナル",
	China:                            "中国",
	Japan:                            "日本",
	ErrDeleteInternationalizedRecord: Messages_en_US.ErrDeleteInternationalizedRecord,
	ChangeLocale:                     Messages_en_US.ChangeLocale,
	SuccessfullyChangedLocale:        Messages_en_US.SuccessfullyChangedLocale,
	ErrChangeLocaleEmpty:             Messages_en_US.ErrChangeLocaleEmpty,
	ErrChangeLocaleUnavailable:       Messages_en_US.ErrChangeLocaleUnavailable,
}

func MustGetTranslation(ctx context.Context, key string) string {
	return i18n.T(ctx, I18nLocalizeKey, key)
}

func MustGetMessages(ctx context.Context) *Messages {
	return i18n.MustGetModuleMessages(ctx, I18nLocalizeKey, Messages_en_US).(*Messages)
}

var Messages_pt_BR = &Messages{
	Localize:                         "Localizar",
	LocalizeFrom:                     "De",
	LocalizeTo:                       "Para",
	CurrentLocalizations:             "Localizações atuais",
	Actions:                          "Ações",
	Localizations:                    "Localizações",
	SuccessfullyLocalized:            "Localizado com sucesso",
	Location:                         "Local",
	Colon:                            ":",
	International:                    "Internacional",
	China:                            "China",
	Japan:                            "Japão",
	ErrDeleteInternationalizedRecord: "Não é possível excluir o registro do idioma padrão enquanto ele estiver internacionalizado.",
	ChangeLocale:                     "Alterar local",
	SuccessfullyChangedLocale:        "Local alterado",
	ErrChangeLocaleEmpty:             "Escolha o local para onde mover este registro.",
	ErrChangeLocaleUnavailable:       "Este local não está mais disponível: o registro já existe nele.",
}
