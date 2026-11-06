import {expect, it} from 'vitest'
import {mount, flushPromises} from '@vue/test-utils'
import {createVuetify} from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import Dialog from '../Dialog.vue'

global.ResizeObserver = require('resize-observer-polyfill')
;(globalThis as any).visualViewport ??= Object.assign(new EventTarget(), {width: 1024, height: 768, scale: 1, offsetLeft: 0, offsetTop: 0})

// a vuetify of a screen width (its display: mobile below md)
function open(width: number, props: Record<string, any>) {
  // the window's width, read by Vuetify's display when it is made
  Object.defineProperty(window, 'innerWidth', {value: width, configurable: true, writable: true})
  const vuetify = createVuetify({components, directives})
  return mount(Dialog, {props: {modelValue: true, expandable: true, title: 'T', ...props}, global: {plugins: [vuetify]}, attachTo: document.body})
}

it('a phone: the whole screen, no expand button (fullscreenOnMobile, the default)', async () => {
  const w = open(400, {})
  await flushPromises()
  const vm = w.vm as any
  expect(vm.isFullscreen).toBe(true)
  expect(vm.showExpand).toBe(false)
  w.unmount()
})

it('a phone, the option off: as on a wide screen', async () => {
  const w = open(400, {fullscreenOnMobile: false})
  await flushPromises()
  const vm = w.vm as any
  expect(vm.isFullscreen).toBe(false)
  expect(vm.showExpand).toBe(true)
  w.unmount()
})

it('a wide screen: the expand button, expanded only when asked', async () => {
  const w = open(1400, {})
  await flushPromises()
  const vm = w.vm as any
  expect(vm.isFullscreen).toBe(false)
  expect(vm.showExpand).toBe(true)
  w.unmount()
  const e = open(1400, {expanded: true})
  await flushPromises()
  expect((e.vm as any).isFullscreen).toBe(true)
  e.unmount()
})
