import { describe, expect, it } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent, getCurrentInstance, nextTick, ref } from 'vue'

import Portal from '../Portal'
import type { UpdateEvent } from '../types'

// A node told by what it is: an element by its tag and class, a text by its
// data.
const desc = (n: Node) =>
  n.nodeType === Node.ELEMENT_NODE
    ? (n as Element).tagName.toLowerCase() + ((n as Element).className ? '.' + (n as Element).className : '')
    : n.nodeType === Node.TEXT_NODE
      ? 'text:' + (n as Text).data
      : '#' + n.nodeName

// A component of several roots, put at the root of a content.
const Multi = defineComponent({ template: '<i class="m1">1</i><i class="m2">2</i>' })

interface Logged {
  name: string
  elements: string[]
  added?: string[]
  removed?: string[]
}

// The page of the portal: siblings before and after it (the elements must
// not reach them), its content (v-model), variables of the content (data:
// refs, reactive in it), every event logged with the portal's elements then.
const page = (tag: string | undefined, slot = '') =>
  defineComponent({
    components: { Portal },
    setup() {
      const content = ref<any>('')
      const log: Logged[] = []
      const data = { show: ref(true), items: ref([1, 2]), Multi }
      const on = (name: string) => (e: UpdateEvent) =>
        log.push({
          name,
          elements: (e.portal.elements || []).map(desc),
          ...(e.added ? { added: e.added.map(desc) } : {}),
          ...(e.removed ? { removed: e.removed.map(desc) } : {})
        })
      const portal = ref<any>()
      return { content, log, data, on, portal, inst: getCurrentInstance() }
    },
    template: `<div class="parent"><span class="before"></span><Portal ref="portal" name="p" ${tag ? `tag="${tag}" class="box"` : ''}
      v-model="content" :data="data"
      @setup="on('setup')($event)" @mounted="on('mounted')($event)" @unmounted="on('unmounted')($event)"
      @pre-content-update="on('pre')($event)" @post-content-update="on('post')($event)"
      @change="on('change')($event)" @elements="on('elements')($event)">${slot}</Portal><span class="after"></span></div>`
  })

const settle = async () => {
  for (let i = 0; i < 4; i++) {
    await flushPromises()
    await nextTick()
  }
  await new Promise((r) => setTimeout(r, 5))
  await flushPromises()
  await nextTick()
}

const names = (log: Logged[]) => log.map((l) => l.name)
const last = (log: Logged[], name: string) => [...log].reverse().find((l) => l.name === name)

for (const tag of [undefined, 'div']) {
  describe(`vx-portal ${tag ? 'with a tag' : 'without a tag'}`, () => {
    const setup = async (slot = '') => {
      const w = mount(page(tag, slot))
      await settle()
      const vm = w.vm as any
      const set = async (v: any) => {
        vm.log.length = 0
        vm.content = v
        await settle()
        return vm.log as Logged[]
      }
      const act = async (f: () => void) => {
        vm.log.length = 0
        f()
        await settle()
        return vm.log as Logged[]
      }
      const elements = () => (vm.portal.elements as Node[]).map(desc)
      return { w, vm, set, act, elements }
    }

    it('mounted with nothing: no content, no element', async () => {
      const { vm, elements } = await setup()
      expect(names(vm.log)).toEqual(['pre', 'post', 'change'])
      expect(elements()).toEqual([])
    })

    it('one root: its events in order, its element', async () => {
      const { set, elements, vm } = await setup()
      const log = await set('<p class="a">A</p>')
      expect(names(log)).toEqual(['pre', 'setup', 'mounted', 'post', 'change'])
      expect(last(log, 'mounted')!.elements).toEqual(['p.a'])
      expect(last(log, 'change')!.elements).toEqual(['p.a'])
      expect(elements()).toEqual(['p.a'])
      // the element the content is in: the tag's, or none
      if (tag) {
        expect((vm.portal.containerElement as Element).className).toContain('box')
      } else {
        expect(vm.portal.containerElement).toBeNull()
      }
    })

    it('several roots: all of them, and none of the siblings of the portal', async () => {
      const { set, elements } = await setup()
      await set('<p class="a">A</p>')
      const log = await set('<p class="a">A</p><p class="b">B</p><p class="c">C</p>')
      // the old content goes (with its elements) before the new is mounted
      expect(names(log)).toEqual(['pre', 'setup', 'unmounted', 'mounted', 'post', 'change'])
      expect(last(log, 'unmounted')!.elements).toEqual(['p.a'])
      expect(last(log, 'mounted')!.elements).toEqual(['p.a', 'p.b', 'p.c'])
      expect(elements()).toEqual(['p.a', 'p.b', 'p.c'])
    })

    it('texts among the roots are roots; empty ones and comments are not', async () => {
      const { set, elements } = await setup()
      await set('x<p class="a">A</p>y<!-- note -->')
      expect(elements()).toEqual(['text:x', 'p.a', 'text:y'])
      await set('hello')
      expect(elements()).toEqual(['text:hello'])
    })

    it('a v-if at the root: the place it leaves is no element; toggled, the elements event', async () => {
      const { set, act, elements, vm } = await setup()
      await set('<p v-if="show" class="x">X</p><p class="y">Y</p>')
      expect(elements()).toEqual(['p.x', 'p.y'])

      let log = await act(() => (vm.data.show.value = false))
      expect(names(log)).toEqual(['elements'])
      expect(log[0]).toMatchObject({ elements: ['p.y'], added: [], removed: ['p.x'] })
      expect(elements()).toEqual(['p.y'])

      log = await act(() => (vm.data.show.value = true))
      expect(log).toEqual([{ name: 'elements', elements: ['p.x', 'p.y'], added: ['p.x'], removed: [] }])
      expect(elements()).toEqual(['p.x', 'p.y'])
    })

    it('a v-for at the root: its items; added and removed by the elements event', async () => {
      const { set, act, elements, vm } = await setup()
      await set('<p v-for="i in items" :key="i" :class="\'i\' + i">{{ i }}</p>')
      expect(elements()).toEqual(['p.i1', 'p.i2'])

      let log = await act(() => vm.data.items.value.push(3))
      expect(log).toEqual([{ name: 'elements', elements: ['p.i1', 'p.i2', 'p.i3'], added: ['p.i3'], removed: [] }])

      log = await act(() => vm.data.items.value.splice(0, 1))
      expect(log).toEqual([{ name: 'elements', elements: ['p.i2', 'p.i3'], added: [], removed: ['p.i1'] }])

      // none left: no element — and not the sibling after the portal
      log = await act(() => (vm.data.items.value = []))
      expect(log).toEqual([{ name: 'elements', elements: [], added: [], removed: ['p.i2', 'p.i3'] }])
      expect(elements()).toEqual([])
    })

    it('a v-for with no item from the start: no element, no sibling', async () => {
      const { set, elements, vm } = await setup()
      vm.data.items.value = []
      await set('<p v-for="i in items" :key="i">{{ i }}</p>')
      expect(elements()).toEqual([])
    })

    it('a component of several roots at the root: its elements', async () => {
      const { set, elements } = await setup()
      await set('<component :is="Multi" />')
      expect(elements()).toEqual(['i.m1', 'i.m2'])
      await set('<b class="b">b</b><component :is="Multi" />')
      expect(elements()).toEqual(['b.b', 'i.m1', 'i.m2'])
    })

    it('a change inside a root is no change of the elements', async () => {
      const { set, act, vm } = await setup()
      await set('<div class="w"><p v-if="show" class="x">X</p></div>')
      const log = await act(() => (vm.data.show.value = false))
      expect(log).toEqual([])
    })

    it('emptied: the content goes first, then post and change with no element', async () => {
      const { set, elements } = await setup()
      await set('<p class="a">A</p><p class="b">B</p>')
      const log = await set('')
      expect(names(log)).toEqual(['pre', 'unmounted', 'post', 'change'])
      expect(last(log, 'unmounted')!.elements).toEqual(['p.a', 'p.b'])
      expect(last(log, 'post')!.elements).toEqual([])
      expect(elements()).toEqual([])
    })

    it('a content with options: its data, its callbacks with the elements', async () => {
      const { set, elements } = await setup()
      const seen: string[][] = []
      await set({
        template: '<p class="v">{{ v }}</p><p class="w">{{ w }}</p>',
        options: { data: { v: 'V', w: 'W' }, mounted: [({ locals }: any) => void seen.push(locals.nodes.map(desc))] }
      })
      expect(elements()).toEqual(['p.v', 'p.w'])
      expect(seen).toEqual([['p.v', 'p.w']])
    })

    it('registered by its name; its content set through it', async () => {
      const { vm, w, elements } = await setup()
      const portals = (vm.inst.appContext.config.globalProperties as any).vx$portals
      expect(portals.p.name).toBe("p")
      expect(portals.p.elements).toEqual([])
      portals.p.content.value = '<p class="n">N</p>'
      await settle()
      expect(elements()).toEqual(['p.n'])
      w.unmount()
      expect(portals.p).toBeUndefined()
    })
  })

  describe(`vx-portal ${tag ? 'with a tag' : 'without a tag'}, its slot`, () => {
    it('null: the slot again, its several roots its elements', async () => {
      const w = mount(page(tag, '<template #default="{show}"><p class="s1">S</p><p class="s2">{{ show }}</p></template>'),
        {})
      await settle()
      const vm = w.vm as any
      const elements = () => (vm.portal.elements as Node[]).map(desc)
      expect(elements()).toEqual(['p.s1', 'p.s2'])
      vm.content = '<p class="a">A</p>'
      await settle()
      expect(elements()).toEqual(['p.a'])
      vm.content = null
      await settle()
      expect(elements()).toEqual(['p.s1', 'p.s2'])
    })

    it('v-pre: the slot as a template of the portal, its roots its elements', async () => {
      const w = mount(page(tag, '<template v-pre><p class="t1">{{ typeof show }}</p><p class="t2">T</p></template>'),
        {})
      await settle()
      const vm = w.vm as any
      expect((vm.portal.elements as Node[]).map(desc)).toEqual(['p.t1', 'p.t2'])
      // the portal's variables: data's (show, a ref unwrapped)
      expect(w.find('.t1').text()).toBe('boolean')
    })
  })
}
