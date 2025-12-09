package web

import (
	"fmt"
	"strings"

	h "github.com/go-rvq/htmlgo"
	"github.com/go-rvq/rvq/web/corejs"
	"github.com/go-rvq/rvq/web/vue"
)

type VueLayoutLocaleSetter struct {
	scriptBuilders []func(lang string) string
	lang           string
}

func (v *VueLayoutLocaleSetter) SetLang(lang string) {
	v.lang = lang
}

func (v *VueLayoutLocaleSetter) Write(w *h.Context) (err error) {
	lang := strings.ReplaceAll(v.lang, "-", "_")
	s := make([]string, len(v.scriptBuilders)+1)
	s[0] = fmt.Sprintf(`window.VueI18n.useI18n().locale.value = %q`, corejs.GetVueI18nLocale(lang))

	for i, script := range v.scriptBuilders {
		s[i+1] = script(lang)
	}

	return vue.UserComponent().Setup(fmt.Sprintf(`({window}) => {%s}`, strings.Join(s, ";"))).Write(w)
}

func (v *VueLayoutLocaleSetter) Append(scriptBuilder ...func(lang string) string) {
	v.scriptBuilders = append(v.scriptBuilders, scriptBuilder...)
}
