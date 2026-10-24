import { afterEach, describe, expect, it } from 'vitest'
import { defineComponent, h, nextTick, ref } from 'vue'
import { VLayout } from 'vuetify/components'

import NavigationDrawer from '../NavigationDrawer.vue'
import { lockScroll, scrollLocks, unlockScroll } from '../scroll-lock'
import { mountTemplate } from '@/lib/__tests__/testutils'

const html = () => document.documentElement.classList

describe('the lock of the page scroll', () => {
  afterEach(() => {
    html().remove('overflow-hidden', 'overflow-y-hidden')
  })

  it('a dialog opened from a dialog: closed in any order, the scroll comes back', () => {
    const drawer = {}
    const texts = {}
    const text = {}
    lockScroll(drawer, 'overflow-y-hidden')
    lockScroll(texts, 'overflow-hidden')
    lockScroll(text, 'overflow-hidden')
    unlockScroll(text)
    expect(html().contains('overflow-hidden')).toBe(true)
    unlockScroll(texts)
    expect(html().contains('overflow-hidden')).toBe(false)
    expect(html().contains('overflow-y-hidden')).toBe(true)
    unlockScroll(drawer)
    expect(html().contains('overflow-y-hidden')).toBe(false)
    expect(scrollLocks()).toBe(0)

    // the parent closed first
    lockScroll(texts, 'overflow-hidden')
    lockScroll(text, 'overflow-hidden')
    unlockScroll(texts)
    expect(html().contains('overflow-hidden')).toBe(true)
    unlockScroll(text)
    expect(html().contains('overflow-hidden')).toBe(false)
  })

  it('a class <html> had before is kept', () => {
    html().add('overflow-hidden')
    const d = {}
    lockScroll(d, 'overflow-hidden')
    unlockScroll(d)
    expect(html().contains('overflow-hidden')).toBe(true)
  })

  it('locked twice, given back once; an unknown owner changes nothing', () => {
    const d = {}
    lockScroll(d, 'overflow-hidden')
    lockScroll(d, 'overflow-hidden')
    unlockScroll({})
    expect(html().contains('overflow-hidden')).toBe(true)
    unlockScroll(d)
    expect(html().contains('overflow-hidden')).toBe(false)
  })
})

describe('a floating drawer and the page scroll', () => {
  afterEach(() => {
    html().remove('overflow-hidden', 'overflow-y-hidden')
  })

  function page(open: ReturnType<typeof ref<boolean>>, shown: ReturnType<typeof ref<boolean>>) {
    const Page = defineComponent({
      setup() {
        return () =>
          h(VLayout, () =>
            shown.value
              ? h(
                  NavigationDrawer,
                  {
                    temporary: true,
                    modelValue: open.value,
                    'onUpdate:modelValue': (v: boolean) => (open.value = v)
                  },
                  () => h('div', 'x')
                )
              : null
          )
      }
    })
    return mountTemplate(Page, {})
  }

  it('locks it while open', async () => {
    const open = ref(false)
    const shown = ref(true)
    const w = page(open, shown)
    open.value = true
    await nextTick()
    expect(html().contains('overflow-y-hidden')).toBe(true)
    open.value = false
    await nextTick()
    expect(html().contains('overflow-y-hidden')).toBe(false)
    w.unmount()
  })

  it('closed with its menu open, the menu closes and gives it back', async () => {
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
    const drawer = w.findComponent(NavigationDrawer).vm as any
    drawer.isMainMenuOpen = true
    await nextTick()
    await nextTick()
    expect(scrollLocks()).toBe(2)
    open.value = false
    await nextTick()
    await nextTick()
    expect(drawer.isMainMenuOpen).toBe(false)
    expect(scrollLocks()).toBe(0)
    expect(html().contains('overflow-y-hidden')).toBe(false)
    w.unmount()
  })

  it('its menu closed by its own close (the arrow), the drawer stays open', async () => {
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
    const drawer = w.findComponent(NavigationDrawer).vm as any
    drawer.isMainMenuOpen = true
    await nextTick()
    await nextTick()
    const menu = w.findAllComponents(NavigationDrawer)[1].vm as any
    menu.cIsOpen = false
    await nextTick()
    await nextTick()
    expect(drawer.isMainMenuOpen).toBe(false)
    expect(open.value).toBe(true)
    expect(scrollLocks()).toBe(1)
    w.unmount()
  })

  it('unmounted while open, gives it back', async () => {
    const open = ref(true)
    const shown = ref(true)
    const w = page(open, shown)
    await nextTick()
    expect(html().contains('overflow-y-hidden')).toBe(true)
    shown.value = false
    await nextTick()
    expect(html().contains('overflow-y-hidden')).toBe(false)
    w.unmount()
  })
})
