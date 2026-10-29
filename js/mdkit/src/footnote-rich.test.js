// bun test (site/components): the footnotes of the rich text editor's HTML.
import { expect, test } from 'bun:test'
import { finalize, noteHTML } from './footnote-rich.js'

test('the editor\'s HTML made goldmark\'s: numbered, ids the labels, the notes in order, ways back', () => {
  const notes = JSON.stringify([{ label: 'b', html: '<p>Nota b</p>' }, { label: 'a', html: '<p>Nota <strong>a</strong></p>' }, { label: 'x', html: '<p>Sem ref</p>' }])
  const html = '<p>Um<sup data-fn="a" class="hermon-fn-ref"></sup> dois<sup data-fn="b" class="hermon-fn-ref"></sup> e<sup data-fn="a" class="hermon-fn-ref"></sup></p>' +
    `<div data-fn-notes="${notes.replace(/&/g, '&amp;').replace(/"/g, '&quot;')}"></div>`
  const out = finalize(html)
  expect(out).toBe('<p>Um<sup id="fnref:a"><a href="#fn:a" class="footnote-ref">1</a></sup> dois<sup id="fnref:b"><a href="#fn:b" class="footnote-ref">2</a></sup> e<sup id="fnref1:a"><a href="#fn:a" class="footnote-ref">1</a></sup></p>' +
    '<div class="footnotes"><hr><ol>' +
    '<li id="fn:a"><p>Nota <strong>a</strong>&#160;<a href="#fnref:a" class="footnote-backref">&#x21a9;&#xfe0e;</a>&#160;<a href="#fnref1:a" class="footnote-backref">&#x21a9;&#xfe0e;</a></p></li>' +
    '<li id="fn:b"><p>Nota b&#160;<a href="#fnref:b" class="footnote-backref">&#x21a9;&#xfe0e;</a></p></li>' +
    '<li id="fn:x"><p>Sem ref</p></li></ol></div>')
  // none: as it is
  expect(finalize('<p>x</p>')).toBe('<p>x</p>')
  // the empty paragraphs before the notes: none
  expect(finalize('<p>a<sup data-fn="1"></sup></p><p></p><p></p><div data-fn-notes="[{&quot;label&quot;:&quot;1&quot;,&quot;html&quot;:&quot;<p>n</p>&quot;}]"></div>'))
    .toBe('<p>a<sup id="fnref:1"><a href="#fn:1" class="footnote-ref">1</a></sup></p><div class="footnotes"><hr><ol><li id="fn:1"><p>n&#160;<a href="#fnref:1" class="footnote-backref">&#x21a9;&#xfe0e;</a></p></li></ol></div>')
})

test('a note read back: without its ways back', () => {
  expect(noteHTML('<p>Nota <strong>a</strong>&#160;<a href="#fnref:a" class="footnote-backref">↩</a>&#160;<a href="#fnref1:a" class="footnote-backref">↩</a></p>'))
    .toBe('<p>Nota <strong>a</strong></p>')
  expect(noteHTML('<p>Nota <a href="#fnref:a" class="footnote-backref" role="doc-backlink">↩︎</a></p>\n')).toBe('<p>Nota</p>')
})
