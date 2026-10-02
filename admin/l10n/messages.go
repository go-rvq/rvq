package l10n

import (
	"context"

	"github.com/go-rvq/rvq/x/i18n"
)

const I18nLocalizeKey i18n.ModuleKey = "I18nLocalizeKey"

type Messages struct {
	ModuleDescription                string `i18n:"hint='Says what the module is, for whoever translates its texts.'"`
	Localize                         string `i18n:"hint='Action that copies a record to other languages.'"`
	LocalizeFrom                     string `i18n:"label='Localize: from', hint='Label of the language a record is localized from.'"`
	LocalizeTo                       string `i18n:"label='Localize: to', hint='Label of the languages a record is localized to.'"`
	CurrentLocalizations             string `i18n:"hint='Title of the languages a record already exists in.'"`
	Actions                          string `i18n:"hint='Title of the actions column of the localizations.'"`
	Localizations                    string `i18n:"hint='Title of a record\\'s versions in other languages.'"`
	SuccessfullyLocalized            string `i18n:"hint='Shown once a record was localized.'"`
	Location                         string `i18n:"label='Language', hint='Label of the language (location) of a record.'"`
	Colon                            string `i18n:"hint='The colon put after a label, as the language writes it.'"`
	International                    string `i18n:"hint='Name of the international (default) language option.'"`
	China                            string `i18n:"hint='Name of the China location option.'"`
	Japan                            string `i18n:"hint='Name of the Japan location option.'"`
	ErrDeleteInternationalizedRecord string `i18n:"label='Delete: localized record', hint='Error when deleting the default-language record of a record that has localizations.'"`
	ChangeLocale                     string `i18n:"label='Change language', hint='Action that moves a record to another language.'"`
	SuccessfullyChangedLocale        string `i18n:"label='Language changed', hint='Shown once a record was moved to another language.'"`
	ErrChangeLocaleEmpty             string `i18n:"label='Change language: none chosen', hint='Error when no language was chosen to move a record to.'"`
	ErrChangeLocaleUnavailable       string `i18n:"label='Change language: taken', hint='Error when the record already exists in the chosen language.'"`
}

var Messages_en_US = &Messages{
	ModuleDescription:                "The localization of records: their versions in each language.",
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
	ModuleDescription:                "记录的本地化：每种语言的版本。",
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
	ModuleDescription:                "レコードのローカライズ：言語ごとのバージョン。",
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
	ModuleDescription:                "A localização dos registros: suas versões em cada idioma.",
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
