package presets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-rvq/rvq/web"
	"github.com/go-rvq/rvq/web/corejs"
	"github.com/go-rvq/rvq/x/ui/vuetify"
	vx "github.com/go-rvq/rvq/x/ui/vuetifyx"
)

// registerVxMessages registers the words of the vuetifyx components (their
// module) and makes them the vue-i18n's namespace vx.
func (b *Builder) registerVxMessages() {
	vx.ConfigureMessages(b.i18nBuilder)
	b.VueI18nMessages(vx.MessagesNamespace, vx.VueMessages)
}

// vueI18nMessages is a namespace of the page's vue-i18n and its words in a
// request's language.
type vueI18nMessages struct {
	ns    string
	words func(ctx context.Context) any
}

// VueI18nMessages makes the words of the namespace ns — words(ctx), in the
// user's language — the page's vue-i18n's: the layout merges them once a page
// (with its locale), and the components read them there (vuetifyx's `vx`, the
// IDE's `gadIde`/`gadGit`), not given by props at each render. A namespace
// registered again is replaced.
func (b *Builder) VueI18nMessages(ns string, words func(ctx context.Context) any) *Builder {
	for i, m := range b.vueI18nMessages {
		if m.ns == ns {
			b.vueI18nMessages[i].words = words
			return b
		}
	}
	b.vueI18nMessages = append(b.vueI18nMessages, vueI18nMessages{ns, words})
	return b
}

// vueI18nScript is the script that merges the registered namespaces, in ctx's
// language, into the page's vue-i18n under the locale of lang (the layout's,
// "pt_BR"); "" with none.
func (b *Builder) vueI18nScript(ctx context.Context) func(lang string) string {
	if len(b.vueI18nMessages) == 0 {
		return nil
	}
	all := map[string]any{}
	for _, m := range b.vueI18nMessages {
		all[m.ns] = m.words(ctx)
	}
	data, err := json.Marshal(all)
	if err != nil {
		return nil
	}
	return func(lang string) string {
		return fmt.Sprintf(`window.VueI18n.useI18n().mergeLocaleMessage(%q, %s)`, corejs.GetVueI18nLocale(lang), data)
	}
}

// appendLocaleScripts adds to the page's locale setter (the layout's) the
// locale of Vuetify and the components' words (VueI18nMessages) merged into
// the vue-i18n, once a page.
func (b *Builder) appendLocaleScripts(ctx *web.EventContext) {
	ls := web.VueLocaleSetterFromContext(ctx)
	if ls == nil {
		return
	}
	ls.Append(func(lang string) string {
		return fmt.Sprintf(`window.Vuetify.useLocale().current.value = %q`, vuetify.GetLocale(lang))
	})
	if script := b.vueI18nScript(ctx.Context()); script != nil {
		ls.Append(script)
	}
}
