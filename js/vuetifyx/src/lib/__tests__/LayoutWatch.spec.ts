import {expect, it, vi} from 'vitest'
import {mount} from '@vue/test-utils'
import LayoutWatch from '../LayoutWatch.vue'

// a window of a width (its media queries)
function width(px: number) {
  const listeners: ((e: any) => void)[] = []
  window.matchMedia = vi.fn((q: string) => {
    const max = Number(/max-width:\s*([\d.]+)px/.exec(q)?.[1] ?? Infinity)
    return {matches: px <= max, media: q, addEventListener: (_: string, f: any) => listeners.push(f), removeEventListener: () => {}} as any
  }) as any
  return listeners
}

// the page all mounted (two frames): then it asks
const frames = () => new Promise((r) => setTimeout(r, 60))

it("a phone, the wide screen drawn: asks the phone's — once the page is mounted", async () => {
  width(400)
  const w = mount(LayoutWatch, {props: {mobile: false}})
  // not at its own mount: the inputs after it have not put their values yet
  expect(w.emitted('change')).toBeUndefined()
  await frames()
  expect(w.emitted('change')).toEqual([[true]])
  w.unmount()
})

it("the layout drawn is the width's: asks nothing", async () => {
  width(400)
  const w = mount(LayoutWatch, {props: {mobile: true}})
  await frames()
  expect(w.emitted('change')).toBeUndefined()
  w.unmount()
  width(1400)
  const x = mount(LayoutWatch, {props: {mobile: false}})
  await frames()
  expect(x.emitted('change')).toBeUndefined()
  x.unmount()
})

it('the width crosses: asks the other', () => {
  const listeners = width(1400)
  const w = mount(LayoutWatch, {props: {mobile: false}})
  listeners.forEach((f) => f({matches: true}))
  expect(w.emitted('change')).toEqual([[true]])
  w.unmount()
})
