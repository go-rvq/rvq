import {describe, expect, it} from 'vitest'
import {registerTiptapPlugin, tiptapPlugins} from './index'

// The plugins of vx-tiptap-editor: rvq's registered (footnotes, tables), an
// application's too, by name and in order; the footnotes' output, the HTML
// kept (goldmark's: the labels in the ids, numbered).
describe('tiptap plugins', () => {
  it('finds the plugins by name, in order; none unknown', () => {
    registerTiptapPlugin('mine', {output: (h) => h + '!'})
    const ps = tiptapPlugins(['tables', 'nope', 'footnotes', 'mine'])
    expect(ps.length).toBe(3)
    expect(ps[0].extensions?.({t: (k) => k}).map((e) => e.name)).toEqual(['mdkitTables'])
    expect(ps[1].extensions?.({t: (k) => k}).map((e) => e.name)).toEqual(['footnoteRef', 'footnotes', 'mdkitFootnotes'])
    expect(ps[2].output?.('x')).toBe('x!')
  })
  it('the footnotes: the HTML kept', () => {
    const [fn] = tiptapPlugins(['footnotes'])
    const html = '<p>a<sup data-fn="n1" class="mdkit-fn-ref"></sup></p><div data-fn-notes="[{&quot;label&quot;:&quot;n1&quot;,&quot;html&quot;:&quot;<p>nota</p>&quot;}]"></div>'
    expect(fn.output?.(html)).toBe('<p>a<sup id="fnref:n1"><a href="#fn:n1" class="footnote-ref">1</a></sup></p>' +
      '<div class="footnotes"><hr><ol><li id="fn:n1"><p>nota&#160;<a href="#fnref:n1" class="footnote-backref">&#x21a9;&#xfe0e;</a></p></li></ol></div>')
  })
})
