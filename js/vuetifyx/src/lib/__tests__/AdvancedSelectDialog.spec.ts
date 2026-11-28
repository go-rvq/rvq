import { expect, it } from 'vitest'

import AdvancedSelect from '../AdvancedSelect.vue'
import { filterData } from '../Filter/FilterData'
import { mountTemplate } from '@/lib/__tests__/testutils'
import { nextTick } from 'vue'

// jsdom has no visualViewport, which the VDialog reads
;(globalThis as any).visualViewport ??= {
  width: 1024,
  height: 768,
  offsetLeft: 0,
  offsetTop: 0,
  scale: 1,
  addEventListener() {},
  removeEventListener() {}
}

const flush = () => new Promise((r) => setTimeout(r, 0))

// a search of 25 records, 10 a page
const searchPage = (keyword: string, page: number) => {
  const all = Array.from({ length: 25 }, (_, i) => ({
    ID: String(i + 1),
    Text: `rec ${i + 1}`
  })).filter((r) => !keyword || r.Text.includes(keyword))
  return Promise.resolve({
    Records: all.slice((page - 1) * 10, page * 10),
    TotalCount: all.length,
    PerPage: 10
  })
}

it('the dialog: a page of the search, picking one closes it', async () => {
  const wrapper = mountTemplate(AdvancedSelect, {
    dialog: true,
    itemValue: 'ID',
    itemText: 'Text',
    selectText: 'Pick',
    searchPageFunc: searchPage
  })
  await nextTick()
  expect(wrapper.text()).toContain('Pick')
  // no autocomplete in the dialog mode
  expect(wrapper.findComponent({ name: 'VAutocomplete' }).exists()).toBe(false)

  const vm: any = wrapper.vm
  vm.openDialog()
  await flush()
  await nextTick()
  const dialog = wrapper.findComponent({ name: 'VDialog' })
  expect(dialog.props('modelValue')).toBe(true)

  vm.pick({ ID: '3', Text: 'rec 3' })
  await nextTick()
  expect(wrapper.emitted('update:modelValue')!.at(-1)).toEqual(['3'])
  expect((wrapper.emitted('update:selectedItems')!.at(-1) as any)[0]).toEqual([
    { ID: '3', Text: 'rec 3' }
  ])
  expect(dialog.props('modelValue')).toBe(false)
  expect(wrapper.text()).toContain('rec 3')
  // chosen: no "select" chip, the change is the icon on the right of it
  expect(wrapper.text()).not.toContain('Pick')
  expect(wrapper.find('.vx-advanced-select__change').exists()).toBe(true)
})

it('the dialog many: picks mark and unmark, the dialog stays open', async () => {
  const wrapper = mountTemplate(AdvancedSelect, {
    dialog: true,
    many: true,
    itemValue: 'ID',
    itemText: 'Text',
    searchPageFunc: searchPage
  })
  const vm: any = wrapper.vm
  vm.openDialog()
  await flush()
  vm.pick({ ID: '1', Text: 'rec 1' })
  vm.pick({ ID: '2', Text: 'rec 2' })
  vm.pick({ ID: '1', Text: 'rec 1' })
  await nextTick()
  expect(wrapper.emitted('update:modelValue')!.at(-1)).toEqual([['2']])
  expect(wrapper.findComponent({ name: 'VDialog' }).props('modelValue')).toBe(true)
})

it('a RecordPickerItem goes as key.in, one record or many', () => {
  expect(
    filterData([
      { key: 'moeda', itemType: 'RecordPickerItem', selected: true, valuesAre: ['a', 'b'] }
    ])
  ).toEqual([['moeda.in', 'a,b']])
})

it('a value no item has (of another key): not drawn, no error', async () => {
  const wrapper = mountTemplate(AdvancedSelect, {
    dialog: true,
    itemValue: 'ID',
    itemText: 'Text',
    label: 'Parent',
    selectText: 'Pick',
    items: [{ ID: 'a_en-US', Text: 'About' }],
    modelValue: 'a',
    searchPageFunc: searchPage
  })
  await nextTick()
  expect(wrapper.text()).toContain('Pick')
  expect(wrapper.text()).toContain('Parent')
  // the label on a line of its own
  expect(wrapper.find('label.d-block').exists()).toBe(true)
})

it('itemsTree: each under its parent when there, else at the root by its full text', async () => {
  const { itemsTree } = await import('../advancedSelectTree')
  const items = [
    { ID: 'a', Name: 'About', Title: 'About', Parent: '' },
    { ID: 't', Name: 'Team', Title: 'About › Team', Parent: 'a' },
    { ID: 'l', Name: 'Leadership', Title: 'About › Team › Leadership', Parent: 't' },
    { ID: 's', Name: 'Small', Title: 'Services › Small', Parent: 'x' }
  ]
  const tree = itemsTree(items, 'ID', 'Parent', 'Name', 'Title')
  expect(tree.map((n) => n.title)).toEqual(['About', 'Services › Small'])
  expect(tree[0].children!.map((n) => n.title)).toEqual(['Team'])
  expect(tree[0].children![0].children!.map((n) => n.title)).toEqual(['Leadership'])
})

// the one chosen: its change beside its removal, both on its right — of its
// line in the list, of its chip (chips)
it('the change and the removal of the one chosen: on its right, together', async () => {
  for (const chips of [false, true]) {
    const wrapper = mountTemplate(AdvancedSelect, {
      dialog: true,
      chips,
      itemValue: 'ID',
      itemText: 'Text',
      items: [{ ID: '3', Text: 'rec 3' }],
      modelValue: '3',
      searchPageFunc: searchPage
    })
    await flush()
    await nextTick()
    const change = wrapper.find('.vx-advanced-select__change')
    const remove = wrapper.find('.vx-advanced-select__remove')
    expect(change.exists() && remove.exists(), `chips ${chips}`).toBe(true)
    // the same parent, the change first
    expect(change.element.parentElement, `chips ${chips}`).toBe(remove.element.parentElement)
    expect(change.element.nextElementSibling, `chips ${chips}`).toBe(remove.element)
    const side = chips ? '.v-chip__append' : '.v-list-item__append'
    expect(change.element.closest(side), `chips ${chips}: on the right`).not.toBeNull()
    if (chips) {
      const chip = wrapper.find('.v-chip.vx-advanced-select__chip')
      expect(chip.text()).toContain('rec 3')
      // a color of its own (chipColor: primary by default), tonal
      expect(chip.classes()).toContain('text-primary')
      expect(chip.classes()).toContain('v-chip--variant-tonal')
      // its text apart from its first icon
      expect(chip.find('.vx-advanced-select__item-text').classes()).toContain('me-2')
    }
  }
})

// itemHref: the one chosen opened by its icon on its right, in place of its
// change — its href for the middle button (a new tab), clicked the event
// openItem (its detail in a dialog) —; its text, clicked, the choice again
// (its dialog) — of its chip and of its line in the list
it('itemHref: its open in place of its change, its text the choice again', async () => {
  for (const chips of [false, true]) {
    const wrapper = mountTemplate(AdvancedSelect, {
      dialog: true,
      chips,
      itemValue: 'ID',
      itemText: 'Text',
      items: [{ ID: 'a b', Text: 'rec 3' }],
      modelValue: 'a b',
      itemHref: '/admin/recs/{id}',
      searchPageFunc: searchPage
    })
    await flush()
    await nextTick()
    const at = `chips ${chips}`
    const open = wrapper.find('a.vx-advanced-select__item-link')
    expect(open.exists(), at).toBe(true)
    expect(open.attributes('href'), at).toBe('/admin/recs/a%20b')
    expect(wrapper.find('.vx-advanced-select__change').exists(), at).toBe(false)
    // beside the removal, on the right
    expect(open.element.nextElementSibling, at).toBe(
      wrapper.find('.vx-advanced-select__remove').element
    )
    await open.trigger('click')
    expect((wrapper.emitted('openItem') as any)?.at(-1)?.[0], at).toEqual({
      ID: 'a b',
      Text: 'rec 3'
    })
    expect((wrapper.vm as any).dialogOpen, at).toBeFalsy()
    // its text: the choice again
    const text = wrapper.find('.vx-advanced-select__item-text')
    expect(text.text(), at).toBe('rec 3')
    await text.trigger('click')
    await flush()
    expect((wrapper.vm as any).dialogOpen, at).toBe(true)
  }
})

// readonly: the one chosen alone — no change (its text changes nothing), no
// removal, no choice — its open still (itemHref)
it('readonly: no change, no removal, no choice', async () => {
  for (const chips of [false, true]) {
    for (const modelValue of ['3', '']) {
      const wrapper = mountTemplate(AdvancedSelect, {
        dialog: true,
        chips,
        readonly: true,
        itemValue: 'ID',
        itemText: 'Text',
        items: [{ ID: '3', Text: 'rec 3' }],
        modelValue,
        selectText: 'Pick',
        itemHref: '/admin/recs/{id}',
        searchPageFunc: searchPage
      })
      await flush()
      await nextTick()
      const at = `chips ${chips}, value "${modelValue}"`
      expect(wrapper.find('.vx-advanced-select__change').exists(), at).toBe(false)
      expect(wrapper.find('.vx-advanced-select__remove').exists(), at).toBe(false)
      expect(wrapper.find('.vx-advanced-select__open').exists(), at).toBe(false)
      expect(wrapper.text(), at).not.toContain('Pick')
      if (modelValue) {
        // its open still; its text changes nothing
        expect(wrapper.find('a.vx-advanced-select__item-link').exists(), at).toBe(true)
        const text = wrapper.find('.vx-advanced-select__item-text')
        expect(text.text(), at).toBe('rec 3')
        await text.trigger('click')
        await flush()
        expect((wrapper.vm as any).dialogOpen, at).toBeFalsy()
      }
    }
  }
})

// the icons of the one chosen: their tooltips (and labels) the texts given —
// the language's
it('the icons of the one chosen: their texts the ones given', async () => {
  for (const [itemHref, icons] of [
    ['/admin/recs/{id}', { 'item-link': 'Abrir', remove: 'Remover' }],
    ['', { change: 'Trocar', remove: 'Remover' }]
  ] as const) {
    const wrapper = mountTemplate(AdvancedSelect, {
      dialog: true,
      chips: true,
      itemValue: 'ID',
      itemText: 'Text',
      items: [{ ID: '3', Text: 'rec 3' }],
      modelValue: '3',
      itemHref,
      openText: 'Abrir',
      changeText: 'Trocar',
      removeText: 'Remover',
      searchPageFunc: searchPage
    })
    await flush()
    await nextTick()
    for (const [cls, text] of Object.entries(icons)) {
      const btn = wrapper.find('.vx-advanced-select__' + cls)
      expect(btn.attributes('aria-label'), `${cls} (${itemHref})`).toBe(text)
      expect(
        btn.findComponent({ name: 'VTooltip' }).exists(),
        `${cls} (${itemHref}): tooltip`
      ).toBe(true)
      // its icon drawn (a tooltip in the button's slot takes the place of
      // its icon prop: the icon its own)
      const icon = {
        'item-link': 'mdi-open-in-new',
        change: 'mdi-swap-horizontal',
        remove: 'mdi-close'
      }[cls]!
      expect(btn.find('.v-icon.' + icon).exists(), `${cls} (${itemHref}): its icon`).toBe(true)
      if (cls === 'remove') {
        // the X of the removal: the color of an error
        expect(btn.find('.v-icon').classes(), `remove (${itemHref}): its color`).toContain(
          'text-error'
        )
      }
    }
  }
})
