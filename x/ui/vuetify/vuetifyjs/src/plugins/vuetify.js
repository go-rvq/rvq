/**
 * plugins/vuetify.js
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'
import * as components from 'vuetify/components'
import * as labComponents from 'vuetify/labs/components'
import * as directives from 'vuetify/directives'
import '@/scss/index.scss'

// Composables
import { createVuetify } from 'vuetify'
import { themes } from './theme'
import * as lang from 'vuetify/locale'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: "light",
    themes,
  },
  locale: {
    locale: 'en',
    fallback: 'en',
    messages: lang
  },
  autoImport: { labs: true },
  components,
  labComponents,
  directives,
});
