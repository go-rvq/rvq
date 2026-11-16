import {expect, it, vi} from 'vitest'
import {flushPromises} from '@vue/test-utils'
// @ts-ignore its source, JavaScript
import {LinkFacets} from '../../../../mdkit/src/link-search.js'
import {mountTemplate} from '@/lib/__tests__/testutils'

// jsdom has no visualViewport, which the VDialog reads
;(globalThis as any).visualViewport ??= {
  width: 1024, height: 768, offsetLeft: 0, offsetTop: 0, scale: 1,
  addEventListener() {}, removeEventListener() {},
}

// A filter of a list: its values chosen as chips, a button opening its dialog
// — its items found as one types (a moment after), checked, applied.
it('a filter of a list opens in a dialog: found, checked, applied', async () => {
  vi.useFakeTimers({toFake: ['setTimeout']})
  const facets = [{key: 'postType', kind: 'autocomplete', label: 'Types', items: [
    {value: 'a', title: 'Alpha', group: 'G1'}, {value: 'b', title: 'Beta', group: 'G1'}, {value: 'c', title: 'Gamma', group: 'G2'}]}]
  let value: any = {postType: ['a']}
  const w = mountTemplate(LinkFacets as any, {facets, modelValue: value, words: {facetChoose: 'Escolher'},
    'onUpdate:modelValue': (v: any) => {
      value = v
      w.setProps({modelValue: v})
    }})
  // the chosen, a chip; the button
  expect(w.find('[data-link-facet-chosen="a"]').text()).toBe('Alpha')
  await w.find('[data-link-facet-open="postType"]').trigger('click')
  await flushPromises()
  const body = () => document.body
  expect(body().querySelectorAll('[data-link-facet-items] [data-link-facet-item]').length).toBe(3)
  // found as one types: a moment after
  const q = body().querySelector('[data-link-facet-q] input') as HTMLInputElement
  q.value = 'gam'
  q.dispatchEvent(new Event('input'))
  await flushPromises()
  expect(body().querySelectorAll('[data-link-facet-items] [data-link-facet-item]').length).toBe(3)
  vi.advanceTimersByTime(300)
  await flushPromises()
  const items = body().querySelectorAll('[data-link-facet-items] [data-link-facet-item]')
  expect([...items].map((e) => e.getAttribute('data-link-facet-item'))).toEqual(['c'])
  ;(items[0] as HTMLElement).click()
  await flushPromises()
  ;(body().querySelector('[data-link-facet-apply]') as HTMLElement).click()
  await flushPromises()
  expect(value).toEqual({postType: ['a', 'c']})
  vi.useRealTimers()
  w.unmount()
})
