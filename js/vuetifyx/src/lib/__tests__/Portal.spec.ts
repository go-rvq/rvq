import { describe, expect, it } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import { defineComponent, nextTick, ref } from 'vue'
import { createVuetify } from 'vuetify'

import Portal from '../Portal/Portal.vue'

// the page of the portal: a variable of its scope (page), one passed on as
// the portal's (data), and the content (v-model)
const Page = defineComponent({
  components: { Portal },
  setup() {
    const content = ref<any>('')
    const page = ref('of the page')
    return { content, page, data: { given: 'passed on' } }
  },
  template: `
    <Portal tag="div" v-model="content" :data="data">
      <template #default="{given}">
        <span class="slot">{{ page }} / {{ given }}</span>
      </template>
    </Portal>`
})

const settle = async () => {
  await flushPromises()
  await nextTick()
  await nextTick()
}

describe('vx-portal', () => {
  it('shows its content, its variables (data), its slot where it says <slot>', async () => {
    const w = mount(Page, { global: { plugins: [createVuetify()] } })
    await settle()
    // mounted: the initial content — the slot, the model being no template
    expect(w.find('.slot').text()).toBe('of the page / passed on')

    // a template: its variables only (data); the slot of the page, with the
    // page's scope, where it says
    ;(w.vm as any).content = '<p class="c">{{ given }} — {{ typeof page }}</p><slot></slot>'
    await settle()
    expect(w.find('.c').text()).toBe('passed on — undefined')
    expect(w.find('.slot').text()).toBe('of the page / passed on')

    // "": nothing
    ;(w.vm as any).content = ''
    await settle()
    expect(w.find('.c').exists()).toBe(false)
    expect(w.find('.slot').exists()).toBe(false)

    // null: the initial content again (the slot)
    ;(w.vm as any).content = null
    await settle()
    expect(w.find('.slot').text()).toBe('of the page / passed on')
  })
})
