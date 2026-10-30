// The blocks of a text the page draws (rvq mdkit blocks): its formulas by
// KaTeX — inline, a block of its own —, its diagrams by Mermaid (the code
// made its picture, its source kept; one that does not draw stays code);
// what is in an editor is the editor's.
import {afterEach, describe, expect, it} from 'vitest'
import {drawBlocks} from '@rvq/mdkit'

const w = window as any
afterEach(() => {
  delete w.katex
  delete w.mermaid
})

describe('blocks of a text', () => {
  it('draws the formulas and the diagrams, not the editor\'s', async () => {
    w.katex = {render: (latex: string, el: HTMLElement, o: any) => { el.textContent = 'K[' + latex + (o.displayMode ? ']!' : ']') }}
    w.mermaid = {
      initialize() {},
      render: async (_id: string, src: string) => {
        if (src.includes('bad')) throw new Error('bad')
        return {svg: '<svg data-src="' + src.trim() + '"></svg>'}
      }
    }
    document.body.innerHTML =
      '<p>a <span data-type="inline-math" data-latex="x^2">x^2</span></p>' +
      '<div data-type="block-math" data-latex="\\frac{1}{2}">\\frac{1}{2}</div>' +
      '<pre><code class="language-mermaid">graph TD\nA--&gt;B</code></pre>' +
      '<pre><code class="language-mermaid">bad</code></pre>' +
      '<div contenteditable="true"><span data-type="inline-math" data-latex="y">y</span></div>'
    await (drawBlocks as any)(document.body, {katex: {src: 'x'}, mermaid: {src: 'x'}})
    expect(document.querySelector('span[data-latex="x^2"]')!.textContent).toBe('K[x^2]')
    expect(document.querySelector('div[data-type="block-math"]')!.textContent).toBe('K[\\frac{1}{2}]!')
    const d = document.querySelector('.mdkit-diagram')!
    expect(d.getAttribute('data-source')).toBe('graph TD\nA-->B')
    expect(d.querySelector('svg')).not.toBeNull()
    expect(document.querySelector('pre.mdkit-diagram-error code')!.textContent).toBe('bad')
    expect(document.querySelector('[contenteditable] span')!.textContent).toBe('y')
  })
})
