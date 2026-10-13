import { afterEach, beforeAll, describe, expect, it } from 'vitest'
import { defineComponent, h, nextTick } from 'vue'
import { VLayout } from 'vuetify/components'

import NavigationDrawer from '../NavigationDrawer.vue'
import { dragWidth, isEndSide, savedWidth, saveWidth } from '../navigation-drawer-resize'
import { mountTemplate } from '@/lib/__tests__/testutils'

// a localStorage of the test's (the page of jsdom here has no origin: none)
beforeAll(() => {
  if (window.localStorage) return
  const data = new Map<string, string>()
  Object.defineProperty(window, 'localStorage', {
    configurable: true,
    value: {
      getItem: (k: string) => (data.has(k) ? data.get(k)! : null),
      setItem: (k: string, v: string) => void data.set(k, String(v)),
      removeItem: (k: string) => void data.delete(k),
      clear: () => data.clear()
    }
  })
})

describe('the width of a resizable drawer', () => {
  afterEach(() => window.localStorage.clear())

  it('follows the drag, kept between its min and max', () => {
    expect(dragWidth(320, 80, 200, 640, false)).toBe(400)
    expect(dragWidth(320, -500, 200, 640, false)).toBe(200)
    expect(dragWidth(320, 900, 200, 640, false)).toBe(640)
    // on the right: a drag to the left widens it
    expect(isEndSide('right')).toBe(true)
    expect(dragWidth(320, -80, 200, 640, true)).toBe(400)
  })

  it('is saved under its key, and read kept in its limits', () => {
    expect(savedWidth('k', 200, 640)).toBeNull()
    saveWidth('k', 410)
    expect(savedWidth('k', 200, 640)).toBe(410)
    saveWidth('k', 900)
    expect(savedWidth('k', 200, 640)).toBe(640)
    saveWidth('k', null)
    expect(savedWidth('k', 200, 640)).toBeNull()
    expect(savedWidth('', 200, 640)).toBeNull()
  })
})

// The drawer in a layout, as a page has it.
function mountDrawer(props: Record<string, any>) {
  const Page = defineComponent({
    setup() {
      return () => h(VLayout, () => h(NavigationDrawer, props, () => h('div', 'menu')))
    }
  })
  return mountTemplate(Page, {})
}

describe('NavigationDrawer resizable', () => {
  afterEach(() => window.localStorage.clear())

  it('has no handle unless resizable', async () => {
    const w = mountDrawer({ modelValue: true, width: 320 })
    await nextTick()
    expect(w.find('[data-drawer-resizer]').exists()).toBe(false)
  })

  it('drags its width by the handle, saves it, and a double click resets it', async () => {
    const w = mountDrawer({ modelValue: true, width: 320, resizable: 'menu.width' })
    await nextTick()
    const handle = w.find('[data-drawer-resizer]')
    expect(handle.exists()).toBe(true)
    expect(handle.classes()).toContain('vx-navigation-drawer__resizer--end')

    await handle.trigger('pointerdown', { clientX: 100 })
    window.dispatchEvent(new MouseEvent('pointermove', { clientX: 180 }))
    window.dispatchEvent(new MouseEvent('pointerup', { clientX: 180 }))
    await nextTick()
    expect(window.localStorage.getItem('menu.width')).toBe('400')
    expect(w.find('.vx-navigation-drawer').attributes('style')).toContain('width: 400px')

    await handle.trigger('dblclick')
    await nextTick()
    expect(window.localStorage.getItem('menu.width')).toBeNull()
    expect(w.find('.vx-navigation-drawer').attributes('style')).toContain('width: 320px')
  })

  it('starts from the width saved', async () => {
    window.localStorage.setItem('menu.width', '450')
    const w = mountDrawer({ modelValue: true, width: 320, resizable: 'menu.width' })
    await nextTick()
    await nextTick()
    expect(w.find('.vx-navigation-drawer').attributes('style')).toContain('width: 450px')
  })

  it('on the right, its handle on its left edge', async () => {
    const w = mountDrawer({ modelValue: true, width: 320, resizable: true, location: 'right' })
    await nextTick()
    expect(w.find('[data-drawer-resizer]').classes()).toContain('vx-navigation-drawer__resizer--start')
  })
})
