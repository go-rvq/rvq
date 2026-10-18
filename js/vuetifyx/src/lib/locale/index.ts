// The words of the vuetifyx components: the namespace `vx` of the page's
// vue-i18n — the admin's layout merges them once a page, in the user's
// language (rvq x/ui/vuetifyx Messages, translatable) —, in its locale (the
// locale's language too: pt_BR, pt), over the defaults (en.ts, pt.ts), which
// are merged under what the page has. With no vue-i18n installed, the
// defaults of the browser's language. A component's props, when given, are
// over all of it.
import {computed, getCurrentInstance, type ComputedRef} from 'vue'
import {useI18n as vueUseI18n, type Composer} from 'vue-i18n'
import en from './en'
import pt from './pt'
import type {VxMessages} from './types'

export type {VxMessages}

/** VX_NAMESPACE is the namespace of the words in the vue-i18n. */
export const VX_NAMESPACE = 'vx'

type Words = Omit<VxMessages, 'monthNames' | 'abbrMonthNames'>

/** VX_DEFAULTS are the defaults by locale (the months aside: the page's, or
 * Intl's). */
export const VX_DEFAULTS: Record<string, Words> = {en, pt}

/** intlMonths are the names of the months in loc by Intl (long or short),
 * January first. */
export function intlMonths(loc: string, month: 'long' | 'short' = 'long'): string[] {
  let f: Intl.DateTimeFormat
  try {
    f = new Intl.DateTimeFormat(loc.replace('_', '-'), {month, timeZone: 'UTC'})
  } catch {
    f = new Intl.DateTimeFormat('en', {month, timeZone: 'UTC'})
  }
  return Array.from({length: 12}, (_, i) => {
    const s = f.format(new Date(Date.UTC(2000, i, 15)))
    return s.charAt(0).toUpperCase() + s.slice(1)
  })
}

/** globalI18n is the page's vue-i18n (its global composer), if installed. */
function globalI18n(): Composer | undefined {
  if (!getCurrentInstance()) return undefined
  try {
    return vueUseI18n({useScope: 'global'}) as unknown as Composer
  } catch {
    return undefined
  }
}

// the composers the defaults were merged into
const merged = new WeakSet<object>()

/** installDefaults merges the defaults under what the vue-i18n has of vx. */
function installDefaults(i18n: Composer) {
  if (merged.has(i18n)) return
  merged.add(i18n)
  for (const [loc, words] of Object.entries(VX_DEFAULTS)) {
    const have = ((i18n.getLocaleMessage(loc) as Record<string, unknown>)[VX_NAMESPACE] ?? {}) as Partial<VxMessages>
    i18n.mergeLocaleMessage(loc, {[VX_NAMESPACE]: {...words, ...have}})
  }
}

/** langOf is the language of a locale: "pt" of "pt_BR", "pt-BR". */
const langOf = (loc: string) => loc.split(/[-_]/)[0].toLowerCase()

/** pick are the words of loc in byLocale: its own, its language's, a locale
 * of its language's. */
function pick(byLocale: Record<string, Partial<VxMessages> | undefined>, loc: string): Partial<VxMessages> {
  const lang = langOf(loc)
  return byLocale[loc] ?? byLocale[lang] ?? byLocale[Object.keys(byLocale).find((k) => langOf(k) === lang) ?? ''] ?? {}
}

/** useVxMessages are the words — the vue-i18n's vx, in its locale, over the
 * defaults —, the given ones over them. In setup. */
export function useVxMessages(given?: () => Partial<VxMessages> | undefined): ComputedRef<VxMessages> {
  const i18n = globalI18n()
  if (i18n) installDefaults(i18n)
  return computed<VxMessages>(() => {
    let words: Partial<VxMessages>
    const loc = String(i18n?.locale.value || (typeof navigator !== 'undefined' && navigator.language) || 'en')
    if (i18n) {
      const all = i18n.messages.value as Record<string, Record<string, Partial<VxMessages>> | undefined>
      const byLocale: Record<string, Partial<VxMessages> | undefined> = {}
      for (const [loc, m] of Object.entries(all)) byLocale[loc] = m?.[VX_NAMESPACE]
      words = pick(byLocale, loc)
    } else {
      words = pick(VX_DEFAULTS, loc)
    }
    // the months: the page's (the admin's i18n), else Intl's in the locale
    const out = {...en, monthNames: intlMonths(loc), abbrMonthNames: intlMonths(loc, 'short')} as Record<string, unknown>
    for (const src of [words, given?.() ?? {}]) {
      for (const [k, v] of Object.entries(src)) {
        if (v === undefined || v === null || v === '') continue
        if (Array.isArray(v) && !v.some((x) => x)) continue
        out[k] = v
      }
    }
    return out as unknown as VxMessages
  })
}

/** capitalized is s with its first letter capital (a month's name: the
 * picker's, as a title). */
export const capitalized = (s: string) => (s ? s.charAt(0).toUpperCase() + s.slice(1) : s)

/** fmt fills a text's {name} with vars. */
export const fmt = (text: string, vars: Record<string, string | number>) =>
  text.replace(/\{(\w+)\}/g, (m, k: string) => (k in vars ? String(vars[k]) : m))

/** useI18n is the components' t(key) — the word of key (useVxMessages) — and
 * the locale (the vue-i18n's, else the browser's language). */
export function useI18n() {
  const i18n = globalI18n()
  const words = useVxMessages()
  const locale = computed(() => String(i18n?.locale.value || (typeof navigator !== 'undefined' && navigator.language) || 'en'))
  return {
    t: (key: string) => {
      const v = (words.value as unknown as Record<string, unknown>)[key]
      return typeof v === 'string' ? v : key
    },
    locale,
  }
}
