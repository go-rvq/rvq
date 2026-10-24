import { afterEach, describe, expect, it } from 'vitest'
import { defineComponent, h, nextTick, ref } from 'vue'
import { VLayout } from 'vuetify/components'

import NavigationDrawer from '../NavigationDrawer.vue'
import { escCloseStack, pushEscClose, removeEscClose } from '../esc-close'
import { mountTemplate } from '@/lib/__tests__/testutils'

const esc = () => {
  const e = new KeyboardEvent('keydown', { key: 'Escape', bubbles: true, cancelable: true })
  document.body.dispatchEvent(e)
  return e
}

const ticks = async () => {
  await nextTick()
  await nextTick()
}

describe('Esc and the floating drawers', () => {
  afterEach(() => {
    document.querySelectorAll('.fake-overlay').forEach((o) => o.remove())
  })

  it('closes the top one only, the last opened', () => {
    const closed: string[] = []
    const a = {}
    const b = {}
    pushEscClose(a, () => null, () => (closed.push('a'), removeEscClose(a)))
    pushEscClose(b, () => null, () => (closed.push('b'), removeEscClose(b)))
    expect(esc().defaultPrevented).toBe(true)
    expect(closed).toEqual(['b'])
    esc()
    expect(closed).toEqual(['b', 'a'])
    expect(esc().defaultPrevented).toBe(false)
    expect(escCloseStack()).toBe(0)
  })

  it('a Vuetify overlay above the drawer (a dialog opened from it) takes the key', () => {
    const closed: string[] = []
    const a = {}
    const nav = document.createElement('nav')
    document.body.appendChild(nav)
    pushEscClose(a, () => nav, () => closed.push('a'))
    const dialog = document.createElement('div')
    dialog.className = 'v-overlay v-overlay--active fake-overlay'
    document.body.appendChild(dialog)
    expect(esc().defaultPrevented).toBe(false)
    expect(closed).toEqual([])

    // the drawer inside the dialog (the dialog's menu): the drawer's
    dialog.appendChild(nav)
    expect(esc().defaultPrevented).toBe(true)
    expect(closed).toEqual(['a'])

    // a tooltip does not count
    document.body.appendChild(nav)
    dialog.className = 'v-overlay v-overlay--active v-tooltip fake-overlay'
    esc()
    expect(closed).toEqual(['a', 'a'])
    removeEscClose(a)
    nav.remove()
  })

  it('closes the menu of a drawer, then the drawer', async () => {
    const open = ref(true)
    const Page = defineComponent({
      setup() {
        return () =>
          h(VLayout, () =>
            h(
              NavigationDrawer,
              {
                temporary: true,
                modelValue: open.value,
                'onUpdate:modelValue': (v: boolean) => (open.value = v)
              },
              { default: () => h('div', 'x'), mainMenu: () => h('div', 'menu') }
            )
          )
      }
    })
    const w = mountTemplate(Page, {})
    await ticks()
    const drawer = w.findComponent(NavigationDrawer).vm as any
    expect(drawer.drawerEl()?.tagName).toBe('NAV')
    drawer.isMainMenuOpen = true
    await ticks()
    expect(escCloseStack()).toBe(2)
    esc()
    await ticks()
    expect(drawer.isMainMenuOpen).toBe(false)
    expect(open.value).toBe(true)
    esc()
    await ticks()
    expect(open.value).toBe(false)
    expect(escCloseStack()).toBe(0)
    w.unmount()
  })
})
