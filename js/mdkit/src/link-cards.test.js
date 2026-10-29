// bun test: the link cards — what a card may link, what the server answers.
import { expect, test } from 'bun:test'
import { cardTarget, cardsOf, safeHref } from './link-cards.js'

test('a card links a path of the site or an address of it; never another site or a script', () => {
  const o = 'https://site.test'
  expect(safeHref('/en-us/about', o)).toBe('/en-us/about')
  expect(safeHref('/admin/posts/1', o)).toBe('/admin/posts/1')
  expect(safeHref('https://site.test/x', o)).toBe('https://site.test/x')
  for (const bad of ['//evil.com/x', '/\\evil.com', 'https://evil.com/x', 'https://site.test.evil.com/', 'javascript:alert(1)', 'data:text/html,x', '', null]) {
    expect(safeHref(bad, o)).toBe('')
  }
})

test('the cards of what the server answers: an event\'s data, cards, or the map', () => {
  const m = { '/a': { href: '/a', title: 'A' } }
  expect(cardsOf({ data: m })).toEqual(m)
  expect(cardsOf({ cards: m })).toEqual(m)
  expect(cardsOf(m)).toEqual(m)
  expect(cardsOf(null)).toEqual({})
})

test('the target of a card typed or pasted: a path, an address of the site, a record path', () => {
  const o = 'https://site.test'
  expect(cardTarget('/en-us/about', o)).toBe('/en-us/about')
  expect(cardTarget('  https://site.test/en-us/about?x=1  ', o)).toBe('/en-us/about?x=1')
  expect(cardTarget('https://site.test', o)).toBe('/')
  expect(cardTarget('admin:posts:<01a0_en-US>:', o)).toBe('admin:posts:<01a0_en-US>:')
  for (const bad of ['//evil.com', 'https://evil.com/x', 'about', 'javascript:alert(1)', '/a b', 'admin:posts', '']) {
    expect(cardTarget(bad, o)).toBe('')
  }
})

test('the address of a link chosen: a path as it is; a record, its address in the admin', async () => {
  const { hrefOf } = await import('./link-picker.js')
  const resolve = async (targets) => Object.fromEntries(targets.map((t) => [t, { href: '/admin/posts/1', title: 'A post' }]))
  expect(await hrefOf('/en-us/about', resolve)).toBe('/en-us/about')
  expect(await hrefOf('admin:posts:<1>:', resolve)).toBe('/admin/posts/1')
  expect(await hrefOf('admin:posts:<2>:', async () => ({}))).toBe('')
  expect(await hrefOf('admin:posts:<2>:', null)).toBe('')
})
