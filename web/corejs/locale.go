package corejs

var VueI18nLocaleAlias = map[string]string{
	"pt_BR": "pt",
}

func GetVueI18nLocale(locale string) (v string) {
	v = VueI18nLocaleAlias[locale]
	if len(v) == 0 {
		v = locale
	}
	return
}
