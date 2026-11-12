import {expect, it} from 'vitest'

import AdvancedSelect from '../AdvancedSelect.vue'
import {filterData} from '../Filter/FilterData'
import {mountTemplate} from '@/lib/__tests__/testutils'
import {nextTick} from 'vue'

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
  const all = Array.from({ length: 25 }, (_, i) => ({ ID: String(i + 1), Text: `rec ${i + 1}` })).filter(
    (r) => !keyword || r.Text.includes(keyword)
  )
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
  expect((wrapper.emitted('update:selectedItems')!.at(-1) as any)[0]).toEqual([{ ID: '3', Text: 'rec 3' }])
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
    filterData([{ key: 'moeda', itemType: 'RecordPickerItem', selected: true, valuesAre: ['a', 'b'] }])
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
  const {itemsTree} = await import('../advancedSelectTree')
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
