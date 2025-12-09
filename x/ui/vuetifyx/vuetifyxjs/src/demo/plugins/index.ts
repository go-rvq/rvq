/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

// Plugins
import vuetify from './vuetify'
import { createI18n } from 'vue-i18n'
// Types
import type { App } from 'vue'

export function registerPlugins(app: App) {
  app.use(createI18n({
    locale: 'pt'
  }))
  app.use(vuetify)

  // fix warning injected property "decorationClasses" is a ref and will be auto-unwrapped
  // https://github.com/ueberdosis/tiptap/issues/1719
}
