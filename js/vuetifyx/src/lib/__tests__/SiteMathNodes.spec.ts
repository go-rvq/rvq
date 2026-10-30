// The formulas of the rich text editor of a site (rvq mdkit: KaTeX loaded by
// the page, not bundled): read from and written as tiptap's mathematics'
// HTML (x/markup's); put and changed by their commands — another kind in its
// place —; drawn by the page's KaTeX.
import {afterEach, describe, expect, it} from 'vitest'
import {Editor} from '@tiptap/core'
import Document from '@tiptap/extension-document'
import Paragraph from '@tiptap/extension-paragraph'
import Text from '@tiptap/extension-text'
import {BlockMath, InlineMath, MathCommands} from '@rvq/mdkit'

let editor: Editor | null = null
afterEach(() => {
  editor?.destroy()
  delete (window as any).katex
})

describe('site math nodes', () => {
  it('reads, writes, puts and changes formulas', async () => {
    ;(window as any).katex = {render: (l: string, el: HTMLElement, o: any) => { el.textContent = 'K[' + l + (o.displayMode ? ']!' : ']') }}
    editor = new Editor({
      element: document.createElement('div'),
      extensions: [Document, Paragraph, Text, InlineMath, BlockMath, MathCommands],
      content: '<p>a <span data-type="inline-math" data-latex="x^2"></span></p><div data-type="block-math" data-latex="y"></div>'
    })
    const e = editor
    const names: string[] = []
    e.state.doc.forEach((n) => names.push(n.type.name))
    expect(names).toEqual(['paragraph', 'blockMath'])
    expect(e.getHTML()).toMatch(/<span[^>]*data-type="inline-math"[^>]*>|<span[^>]*data-latex="x\^2"/)
    await new Promise((r) => setTimeout(r, 10))
    expect(e.view.dom.querySelector('[data-type="block-math"]')!.textContent).toBe('K[y]!')
    ;(e.commands as any).insertMath('z', false)
    expect(e.getHTML()).toContain('data-latex="z"')
    let pos = -1
    e.state.doc.forEach((n, p) => { if (n.type.name === 'blockMath') pos = p })
    ;(e.commands as any).updateMath(pos, 'w', false)
    const after: string[] = []
    e.state.doc.descendants((n) => { if (n.type.name.endsWith('Math')) after.push(n.type.name + ':' + n.attrs.latex) })
    expect(after).toContain('inlineMath:w')
    expect(after).not.toContain('blockMath:y')
  })
})
