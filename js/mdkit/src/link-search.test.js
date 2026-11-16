// bun test: the facets of the search of the links — their values, the ones
// no longer there, their items in groups — and the answers read.
import { expect, test } from 'bun:test'
import { countText, facetValues, foundOf, itemsOf, pruneValues, resultOf, tabOf } from './link-search.js'

const facets = [
  { key: 'kind', kind: 'chips', default: ['pages', 'posts'], items: [{ value: 'pages', title: 'Pages' }, { value: 'posts', title: 'Posts' }] },
  { key: 'postTag', kind: 'autocomplete', items: [
    { value: 't1', title: 'Big', group: 'News' }, { value: 't2', title: 'Small', group: 'News' }, { value: 't3', title: 'Wet', group: 'Services' },
  ] },
]

test('resultOf: the groups alone, or with the facets', () => {
  expect(resultOf([{ items: [] }])).toEqual({ groups: [{ items: [] }], facets: null })
  expect(resultOf({ groups: [], facets })).toEqual({ groups: [], facets })
  expect(resultOf(null)).toEqual({ groups: [], facets: null })
})

test('facetValues: the chosen, else the defaults', () => {
  expect(facetValues(facets, {})).toEqual({ kind: ['pages', 'posts'], postTag: [] })
  expect(facetValues(facets, { kind: [] })).toEqual({ kind: [], postTag: [] })
})

test('pruneValues: a value no longer an item, out', () => {
  expect(pruneValues(facets, { kind: ['posts'], postTag: ['t1'] })).toEqual({ values: { kind: ['posts'], postTag: ['t1'] }, pruned: false })
  expect(pruneValues(facets, { postTag: ['t1', 'gone'] })).toEqual({ values: { postTag: ['t1'] }, pruned: true })
  // a facet no longer there
  expect(pruneValues([facets[0]], { kind: ['pages'], postTag: ['t1'] })).toEqual({ values: { kind: ['pages'] }, pruned: true })
})

test('itemsOf: a header before each group', () => {
  expect(itemsOf(facets[1]).map((it) => (it.header ? '#' + it.title : it.value))).toEqual(['#News', 't1', 't2', '#Services', 't3'])
  expect(itemsOf(facets[0]).map((it) => it.value)).toEqual(['pages', 'posts'])
})

test('flatten and isTree: the options of a tree by their targets; a group a tree', async () => {
  const { flatten, isTree } = await import('./link-search.js')
  const items = [{ target: '/a', title: 'A', context: true, children: [{ target: '/a/b', title: 'B', children: [{ target: '/a/b/c', title: 'C' }] }] }]
  expect(Object.keys(flatten(items))).toEqual(['/a', '/a/b', '/a/b/c'])
  expect(isTree({ items })).toBe(true)
  expect(isTree({ kind: 'tree', items: [{ target: '/x' }] })).toBe(true)
  expect(isTree({ kind: 'posts', items: [{ target: '/x' }] })).toBe(false)
})

test('foundOf and countText: the count, else the options found (not those holding others); + of more', () => {
  expect(foundOf({ count: 5, items: [{}] })).toBe(5)
  const tree = { items: [{ context: true, children: [{ title: 'a' }, { title: 'b', children: [{ title: 'c' }] }] }] }
  expect(foundOf(tree)).toBe(3)
  expect(countText({ count: 8, more: true, items: [] })).toBe('8+')
  expect(countText({ items: [{}, {}] })).toBe('2')
  expect(tabOf({ kind: 'files', title: 'Arquivos' })).toBe('files')
  expect(tabOf({ title: 'Users' })).toBe('Users')
})
