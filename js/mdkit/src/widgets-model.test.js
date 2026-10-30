// bun test: the buttons and the tabs, their HTML and Markdown, read back.
import { expect, test } from 'bun:test'
import { Window } from 'happy-dom'
import { blockAt, buttonsHTML, buttonsOf, tabsHTML, tabsMarkdown, tabsOf } from './widgets-model.js'

const buttons = {
  align: 'center',
  items: [
    { href: '/en-us/about', target: '/en-us/about', title: '', description: 'Who we are', color: 'primary', iconLeft: 'mdi-information', iconRight: '', newTab: false },
    { href: 'https://x.test/?a=1&b="2"', target: '', title: 'Site <x>', description: '', color: 'nope', iconLeft: 'bad icon', iconRight: 'mdi-open-in-new', newTab: true },
  ],
}

const doc = () => new Window().document

test('buttons: their HTML (escaped, only known colors and icons) and back', () => {
  const html = buttonsHTML(buttons)
  expect(html).toBe('<div data-type="buttons" data-align="center">' +
    '<a data-type="button" href="/en-us/about" data-href="/en-us/about" title="Who we are" data-color="primary" data-icon-left="mdi-information"></a>' +
    '<a data-type="button" href="https://x.test/?a=1&amp;b=&quot;2&quot;" data-icon-right="mdi-open-in-new" target="_blank" rel="noopener">Site &lt;x&gt;</a></div>')
  const d = doc()
  d.body.innerHTML = html
  const v = buttonsOf(d.body.firstChild)
  expect(v.align).toBe('center')
  expect(v.items[0]).toEqual({ ...buttons.items[0] })
  expect(v.items[1]).toEqual({ ...buttons.items[1], color: '', iconLeft: '' })
  // in Markdown: the same line, read back by blockAt
  const md = 'Before\n\n' + html + '\n\nAfter'
  const b = blockAt(md, md.indexOf('Site'))
  expect(b.kind).toBe('buttons')
  expect(md.slice(b.from, b.to)).toBe(html)
  expect(b.value).toEqual(v)
  expect(blockAt(md, 2)).toBe(null)
})

const tabs = {
  align: 'end',
  items: [
    { title: 'One "1"', color: 'success', iconLeft: '', iconRight: 'mdi-star', body: '**bold**\n\n<div class="mdkit-link-card" data-href="/x">/x</div>' },
    { title: 'Two', color: '', iconLeft: 'mdi-home', iconRight: '', body: '- a\n- b' },
  ],
}

test('tabs: their Markdown (bodies between blank lines) and back', () => {
  const md = tabsMarkdown(tabs)
  expect(md).toBe('<div data-type="tabs" data-align="end">\n' +
    '<div data-type="tab" data-title="One &quot;1&quot;" data-color="success" data-icon-right="mdi-star">\n\n**bold**\n\n<div class="mdkit-link-card" data-href="/x">/x</div>\n\n</div>\n' +
    '<div data-type="tab" data-title="Two" data-icon-left="mdi-home">\n\n- a\n- b\n\n</div>\n</div>')
  const text = '# Doc\n\n' + md + '\n\nEnd'
  const b = blockAt(text, text.indexOf('- b'))
  expect(b.kind).toBe('tabs')
  expect(text.slice(b.from, b.to)).toBe(md)
  expect(b.value).toEqual(tabs)
})

test('tabs: their HTML and back', () => {
  const v = { ...tabs, items: tabs.items.map((it) => ({ ...it, body: '<p>' + it.title + '</p>' })) }
  const d = doc()
  d.body.innerHTML = tabsHTML(v)
  const got = tabsOf(d.body.firstChild)
  expect(got.align).toBe('end')
  expect(got.items.map((it) => [it.title, it.color, it.iconLeft, it.iconRight])).toEqual([['One "1"', 'success', '', 'mdi-star'], ['Two', '', 'mdi-home', '']])
  expect(got.items[1].body).toBe('<p>Two</p>')
})
