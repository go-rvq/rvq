// The link card of the rich text editor (rvq mdkit LinkCard): read from and
// written as x/markup's HTML (div.mdkit-link-card[data-href]); put by its
// command — a path of the site, a record path; never another site —; an
// address pasted alone made a card — before the link's paste —, Ctrl+Z the
// text pasted.
import {afterEach, describe, expect, it} from 'vitest'
import {Editor} from '@tiptap/core'
import Document from '@tiptap/extension-document'
import Paragraph from '@tiptap/extension-paragraph'
import Text from '@tiptap/extension-text'
import {UndoRedo} from '@tiptap/extensions'
import Link from '@tiptap/extension-link'
import {LinkCard} from '@rvq/mdkit'

// jsdom has no ClipboardEvent (ProseMirror's pasteText makes one)
if (typeof (globalThis as any).ClipboardEvent === 'undefined') {
  (globalThis as any).ClipboardEvent = class extends Event {
    clipboardData = null
  }
}

let editor: Editor | null = null
afterEach(() => editor?.destroy())

const make = (content = '') => {
  editor = new Editor({element: document.createElement('div'), extensions: [Document, Paragraph, Text, UndoRedo, Link, LinkCard], content})
  return editor
}

describe('link card node', () => {
  it('is no first block: an empty text is a paragraph, not a card', () => {
    const e = make('')
    expect(e.state.doc.firstChild!.type.name).toBe('paragraph')
    expect(e.getHTML()).toBe('<p></p>')
  })

  it('reads and writes the card of x/markup', () => {
    const e = make('<p>a</p><div class="mdkit-link-card" data-href="/en-us/about"><a href="/en-us/about">/en-us/about</a></div>' +
      '<div class="mdkit-link-card" data-href="//evil.com/x">x</div>')
    const html = e.getHTML()
    expect(html).toContain('<div class="mdkit-link-card" data-href="/en-us/about"><a href="/en-us/about">/en-us/about</a></div>')
    expect(html).not.toContain('evil.com')
    expect(e.view.dom.querySelector('[data-link-card-edit] .mdkit-link-card-bin')).not.toBeNull()
  })

  it('puts a card by its command, a target a card may link only', () => {
    const e = make('<p>a</p>')
    expect((e.commands as any).insertLinkCard('admin:posts:<1>:')).toBe(true)
    expect((e.commands as any).insertLinkCard('https://evil.com/')).toBe(false)
    // a card chosen (the one put): the next after it, not over it
    expect((e.commands as any).insertLinkCard('/en-us/about')).toBe(true)
    expect(e.getHTML().match(/mdkit-link-card/g)!.length).toBe(2)
    // (its < > unescaped in the attribute: valid HTML; the server's policy
    // escapes them)
    expect(e.getHTML()).toContain('<div class="mdkit-link-card" data-href="admin:posts:<1>:">admin:posts:&lt;1&gt;:</div>')
  })

  it('makes an address pasted alone a card; Ctrl+Z, the text', () => {
    const e = make('<p></p>')
    e.commands.focus()
    e.view.pasteText('/en-us/about')
    expect(e.getHTML()).toContain('data-href="/en-us/about"')
    // the empty paragraph it was pasted in: the card in its place
    expect(e.getHTML()).not.toContain('<p></p><div')
    e.commands.undo()
    expect(e.getHTML()).not.toContain('mdkit-link-card')
    expect(e.getText()).toContain('/en-us/about')
    // more than an address: text
    e.commands.setContent('<p></p>')
    e.view.pasteText('see /en-us/about')
    expect(e.getHTML()).not.toContain('mdkit-link-card')
  })
})
