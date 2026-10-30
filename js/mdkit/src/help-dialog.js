// mdkit-help-dialog: the guide of an editor in a dialog — a Vue component
// (open(src)) —, the documentation's own: src is the page of the guide (rvq
// userdocs DocHandler: the same files and edits as the documentation page,
// so a change of the guide is a change here), loaded as the content of a
// rvq-portal (?portal=1; js/portal: the app's, the admin's; else the page's
// RvqPortal, a site's) in the language chosen before, else the page's. Its
// links to that page — the other guides, the languages it can be read in
// (one chosen kept) — load in the portal too; its code, formulas and
// diagrams drawn by the page's blocks (blocks.js). Its button opens in a
// new tab href (the admin's documentation), else the page of the guide.
import { wordsOf } from './config.js'

// the language chosen, kept in the browser
const chosenKey = 'mdkit.help.locale'

const WORDS = {
  en: { title: 'Guide', open: 'Open in a new tab', close: 'Close', failed: 'The guide could not be loaded.', back: 'Back' },
  pt: { title: 'Guia', open: 'Abrir em nova aba', close: 'Fechar', failed: 'Não foi possível carregar o guia.', back: 'Voltar' },
  es: { title: 'Guía', open: 'Abrir en una pestaña nueva', close: 'Cerrar', failed: 'No se pudo cargar la guía.', back: 'Volver' },
}

// pageLang is the language of the page (its html lang).
export function pageLang() {
  return (typeof document !== 'undefined' && document.documentElement.lang) || ''
}

// helpOf is where the page's meta (mdkit-help: {rich, markdown}) says the
// guide of an editor is asked; '' none.
export function helpOf(kind) {
  const meta = typeof document !== 'undefined' && document.querySelector('meta[name="mdkit-help"]')
  try {
    return ((meta && JSON.parse(meta.getAttribute('content') || '{}')) || {})[kind] || ''
  } catch (e) {
    return ''
  }
}

// withParams is src with the params of q set (none, those not there only).
export function withParams(src, q, keep = false) {
  const u = new URL(src, typeof location !== 'undefined' ? location.href : 'http://x/')
  for (const [k, v] of Object.entries(q)) if (v && !(keep && u.searchParams.has(k))) u.searchParams.set(k, v)
  return u.pathname + u.search
}

// titleOf is the title of the content of a guide (its first heading).
export function titleOf(html) {
  const m = /<h1[^>]*>([\s\S]*?)<\/h1>/i.exec(html || '')
  return m ? m[1].replace(/<[^>]+>/g, '').replace(/&amp;/g, '&').replace(/&lt;/g, '<').replace(/&gt;/g, '>').replace(/&quot;/g, '"').replace(/&#39;/g, "'").trim() : ''
}

const read = () => {
  try { return localStorage.getItem(chosenKey) || '' } catch (e) { return '' }
}
const keep = (l) => {
  try { localStorage.setItem(chosenKey, l) } catch (e) { /* none kept */ }
}

export default {
  name: 'MdkitHelpDialog',
  props: {
    words: { type: Object, default: () => wordsOf('help') },
    // href: the guide's page in the admin's documentation, opened in a new
    // tab (none: the page of the guide — a site's)
    href: { type: String, default: '' },
  },
  data: () => ({ show: false, src: '', content: '', loading: false, failed: false, history: [] }),
  computed: {
    w() {
      const l = pageLang().slice(0, 2).toLowerCase()
      return { ...(WORDS[l] || WORDS.en), ...this.words }
    },
    // the rvq-portal: the app's, else the page's runtime's (a site's)
    portal() {
      const c = this.$.appContext.components
      return c['rvq-portal'] || c.RvqPortal || (typeof window !== 'undefined' && window.RvqPortal && window.RvqPortal.Portal) || null
    },
    title() { return titleOf(this.content) || this.w.title },
    // what the button opens in a new tab: the documentation's page given
    // (the first guide), else the page of the one shown
    tabHref() {
      if (!this.src) return ''
      return this.href && !this.history.length ? this.href : this.src
    },
  },
  template: `
<v-dialog v-model="show" max-width="920" scrollable data-help-dialog>
  <v-card>
    <v-card-title class="d-flex align-center ga-2">
      <v-btn v-if="history.length" icon="mdi-arrow-left" size="small" variant="text" :title="w.back" @click="back" data-help-back></v-btn>
      <v-icon icon="mdi-help-circle-outline" color="primary"></v-icon>
      <span class="text-truncate" data-help-title>{{ title }}</span>
      <v-spacer></v-spacer>
      <v-btn v-if="tabHref" icon="mdi-open-in-new" size="small" variant="text" :href="tabHref" target="_blank" rel="noopener"
        :title="w.open" data-help-open></v-btn>
      <v-btn icon="mdi-close" size="small" variant="text" :title="w.close" @click="show = false"></v-btn>
    </v-card-title>
    <v-progress-linear v-if="loading" indeterminate color="primary" height="2"></v-progress-linear>
    <v-divider></v-divider>
    <v-card-text style="max-height: 75vh">
      <v-alert v-if="failed" type="warning" variant="tonal" density="compact" data-help-failed>{{ w.failed }}</v-alert>
      <div v-if="content && portal" @click="follow" data-help-doc>
        <component :is="portal" tag="div" :model-value="content"></component>
      </div>
    </v-card-text>
  </v-card>
</v-dialog>`,
  methods: {
    // open shows the guide whose page is src
    open(src) {
      this.history = []
      this.show = true
      this.load(withParams(src, { locale: read() || pageLang() }, true))
    },
    // load puts the guide whose page is src in the portal
    async load(src) {
      this.src = src
      this.loading = true
      this.failed = false
      try {
        const r = await fetch(withParams(src, { portal: '1' }), { credentials: 'same-origin' })
        if (!r.ok) throw new Error(String(r.status))
        this.content = await r.text()
      } catch (e) {
        this.failed = true
      } finally {
        this.loading = false
      }
    },
    back() {
      const src = this.history.pop()
      if (src) this.load(src)
    },
    // follow: a link to the page of a guide (another, a language) loads
    // it here — a language chosen kept for the next —; another, in a new tab
    follow(e) {
      const a = e.target && e.target.closest && e.target.closest('a[href]')
      if (!a || e.ctrlKey || e.metaKey || e.shiftKey || e.button) return
      const u = new URL(a.getAttribute('href'), location.href)
      const page = new URL(this.src, location.href)
      if (u.origin === location.origin && u.pathname === page.pathname) {
        e.preventDefault()
        if (a.hasAttribute('data-locale')) keep(a.getAttribute('data-locale'))
        else this.history.push(this.src)
        this.load(u.pathname + u.search)
        return
      }
      if (a.getAttribute('href').startsWith('#')) return
      e.preventDefault()
      window.open(u.href, '_blank', 'noopener')
    },
  },
}
