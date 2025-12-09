package vuetify

var LocaleAlias = map[string]string{
	"pt_BR": "pt",
}

func GetLocale(locale string) (v string) {
	v = LocaleAlias[locale]
	if len(v) == 0 {
		v = locale
	}
	return
}
