// The blocks of the rich text editor: a formula (tiptap's mathematics), a
// diagram (rvq mdkit's: a ```mermaid fence's HTML, before the code block), a
// section that folds (tiptap's details) — read from the HTML x/markup keeps
// and written as it; their commands.
import {afterEach, describe, expect, it} from 'vitest'
import {Editor} from '@tiptap/core'
import Document from '@tiptap/extension-document'
import Paragraph from '@tiptap/extension-paragraph'
import Text from '@tiptap/extension-text'
import {CodeBlock} from '@tiptap/extension-code-block'
import Mathematics from '@tiptap/extension-mathematics'
import {Details, DetailsContent, DetailsSummary} from '@tiptap/extension-details'
import {Diagram} from '@rvq/mdkit'

let editor: Editor | null = null
afterEach(() => editor?.destroy())

const make = (content: string) => {
  editor = new Editor({
    element: document.createElement('div'),
    extensions: [Document, Paragraph, Text, CodeBlock, Mathematics, Details, DetailsSummary, DetailsContent, Diagram],
    content
  })
  return editor
}

describe('block nodes', () => {
  it('reads and writes formulas, diagrams and sections that fold', () => {
    const e = make('<p>a <span data-type="inline-math" data-latex="x^2"></span></p>' +
      '<div data-type="block-math" data-latex="\\frac{1}{2}"></div>' +
      '<pre><code class="language-mermaid">graph TD\nA--&gt;B</code></pre>' +
      '<pre><code class="language-go">x := 1</code></pre>' +
      '<details><summary>More</summary><div data-type="detailsContent"><p>inside</p></div></details>')
    const names: string[] = []
    e.state.doc.forEach((n) => names.push(n.type.name))
    expect(names).toEqual(['paragraph', 'blockMath', 'mdkitDiagram', 'codeBlock', 'details'])
    const html = e.getHTML()
    expect(html).toMatch(/<span (data-latex="x\^2" data-type="inline-math"|data-type="inline-math" data-latex="x\^2")>/)
    expect(html).toMatch(/<div (data-latex="\\frac\{1\}\{2\}" data-type="block-math"|data-type="block-math" data-latex="\\frac\{1\}\{2\}")>/)
    expect(html).toContain('<pre><code class="language-mermaid">graph TD\nA--&gt;B</code></pre>')
    expect(html).toContain('<summary>More</summary>')
    expect(html).toContain('data-type="detailsContent"')
  })

  it('puts and changes a diagram', () => {
    const e = make('<p>a</p>')
    ;(e.commands as any).insertDiagram('graph LR\nX-->Y')
    let pos = -1
    e.state.doc.forEach((n, p) => { if (n.type.name === 'mdkitDiagram') pos = p })
    expect(pos).toBeGreaterThan(-1)
    ;(e.commands as any).updateDiagram(pos, 'graph LR\nX-->Z')
    expect(e.getHTML()).toContain('X--&gt;Z')
  })
})
