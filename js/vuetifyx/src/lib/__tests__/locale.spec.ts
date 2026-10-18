import {expect, it} from 'vitest'
import {mount} from '@vue/test-utils'
import {defineComponent, h, nextTick} from 'vue'
import {createI18n} from 'vue-i18n'
import {fmt, intlMonths, useVxMessages} from '@/lib/locale'
import type {VxMessages} from '@/lib/locale'

// probe mounts a component that shows a word of useVxMessages (given over it)
function probe(key: keyof VxMessages, plugins: any[] = [], given?: Partial<VxMessages>) {
  let words: any
  const C = defineComponent({
    setup() {
      words = useVxMessages(() => given)
      return () => h('span', String(words.value[key]))
    },
  })
  const w = mount(C, {global: {plugins}})
  return {w, words: () => words.value as VxMessages}
}

it('the page vue-i18n words win, its locale chosen', async () => {
  const i18n = createI18n({legacy: false, locale: 'pt', fallbackLocale: 'en'})
  // the admin merges its words (the layout), the months its i18n default module's
  i18n.global.mergeLocaleMessage('pt', {vx: {filterClearAll: 'Limpar tudo', monthNames: ['Jan.', 'Fev.', 'Mar.', 'Abr.', 'Mai.', 'Jun.', 'Jul.', 'Ago.', 'Set.', 'Out.', 'Nov.', 'Dez.']}})
  const {w, words} = probe('filterClearAll', [i18n])
  expect(w.text()).toBe('Limpar tudo')
  // the rest: the defaults (pt.ts), merged under
  expect(words().clear).toBe('Limpar')
  expect(words().monthNames[1]).toBe('Fev.')
  // the locale changed: its words
  i18n.global.locale.value = 'en'
  await nextTick()
  expect(w.text()).toBe('Clear Filters')
})

it('a locale of the language: pt_BR reads pt', () => {
  const i18n = createI18n({legacy: false, locale: 'pt_BR', fallbackLocale: 'en'})
  const {w} = probe('apply', [i18n])
  expect(w.text()).toBe('Aplicar')
})

it('the months out of the page: Intl in the vue-i18n locale, not the browser', () => {
  const i18n = createI18n({legacy: false, locale: 'pt', fallbackLocale: 'en'})
  const {words} = probe('monthNames', [i18n])
  expect(words().monthNames).toEqual(intlMonths('pt'))
  expect(words().monthNames).toHaveLength(12)
})

it('the given words win', () => {
  const i18n = createI18n({legacy: false, locale: 'en'})
  const {w} = probe('apply', [i18n], {apply: 'Go'})
  expect(w.text()).toBe('Go')
})

it('with no vue-i18n: the defaults', () => {
  const {words} = probe('apply')
  expect(['Apply', 'Aplicar']).toContain(words().apply)
})

it('fmt fills {name}', () => {
  expect(fmt('Filter by {filter}', {filter: 'Name'})).toBe('Filter by Name')
  expect(fmt('{a} {b}', {a: 1})).toBe('1 {b}')
})
