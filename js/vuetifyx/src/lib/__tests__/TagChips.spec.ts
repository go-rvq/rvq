import { describe, expect, it } from 'vitest'
import { defineComponent, h, nextTick, ref } from 'vue'
import TagChips from '../TagChips'
import { mountTemplate } from '@/lib/__tests__/testutils'

const items = [
  { value: '1', title: 'Urgente', color: 'red', hint: 'por Ana', group: 'a' },
  { value: '2', title: 'Retornar', group: 'a' },
  { value: '3', title: 'Outro tipo', group: 'b' }
]

function page(model: ReturnType<typeof ref<string[]>>, extra: Record<string, any> = {}) {
  const added: string[] = []
  const removed: string[] = []
  const Page = defineComponent({
    setup() {
      return () =>
        h(TagChips, {
          items,
          modelValue: model.value,
          'onUpdate:modelValue': (v: string[]) => (model.value = v),
          onAdd: (v: string) => added.push(v),
          onRemove: (v: string) => removed.push(v),
          ...extra
        })
    }
  })
  return { w: mountTemplate(Page, {}), added, removed }
}

describe('the tag chips', () => {
  it('a click on a tag not chosen chooses it; the trash of one chosen takes it off', async () => {
    const model = ref<string[]>(['1'])
    const { w, added, removed } = page(model, { group: 'a' })
    await nextTick()
    // the tags of the group alone
    expect(w.findAll('[data-tag]').map((c) => c.attributes('data-tag'))).toEqual(['1', '2'])
    expect(w.find('[data-tag="1"]').attributes('data-tag-on')).toBe('true')

    await w.find('[data-tag="2"]').trigger('click')
    expect(model.value).toEqual(['1', '2'])
    expect(added).toEqual(['2'])

    // the chip itself chosen: nothing (only its trash takes it off)
    await w.find('[data-tag="1"]').trigger('click')
    expect(model.value).toEqual(['1', '2'])

    await nextTick()
    await w.find('[data-tag="1"] .v-chip__close').trigger('click')
    expect(model.value).toEqual(['2'])
    expect(removed).toEqual(['1'])
    w.unmount()
  })

  it('read-only: no change', async () => {
    const model = ref<string[]>([])
    const { w, added } = page(model, { readonly: true })
    await nextTick()
    await w.find('[data-tag="2"]').trigger('click')
    expect(model.value).toEqual([])
    expect(added).toEqual([])
    w.unmount()
  })
})
