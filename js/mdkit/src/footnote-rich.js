// The footnotes of the rich text editor's HTML — text only (no DOM) —: the
// tiptap nodes' HTML (a reference `<sup data-fn="label">`, the notes in
// `<div data-fn-notes="[{label, html}]">`) made the HTML kept (finalize),
// goldmark's — the references numbered by their order, their ids the labels
// (`<sup id="fnref:label"><a href="#fn:label" class="footnote-ref">1</a></sup>`,
// the second to a note fnref1:label), the notes at the end in that order,
// each with its ways back (`<div class="footnotes"><hr><ol><li id="fn:label">…`)
// —; the HTML of a note read back without its ways back (noteHTML).

const unescapeAttr = (s) => String(s).replace(/&quot;/g, '"').replace(/&#39;/g, "'").replace(/&lt;/g, '<').replace(/&gt;/g, '>').replace(/&amp;/g, '&')
const escapeAttr = (s) => String(s).replace(/&/g, '&amp;').replace(/"/g, '&quot;')

// backHTML is goldmark's way back to a reference.
const backHTML = '&#x21a9;&#xfe0e;'

// noteHTML is the HTML of a note (its li's) without its ways back and the
// space before them.
export function noteHTML(html) {
  return String(html || '')
    .replace(/(?:&#160;|&nbsp;| )?<a[^>]*class="footnote-backref"[^>]*>[\s\S]*?<\/a>/g, '')
    .replace(/(?:&#160;|&nbsp;| |\s)+(<\/p>)?\s*$/, '$1')
    .trim()
}

// finalize is the editor's HTML with its footnotes as kept (see above).
export function finalize(html) {
  let notes = []
  let out = String(html || '').replace(/<div[^>]*\bdata-fn-notes="([^"]*)"[^>]*>\s*<\/div>/g, (m, json) => {
    try { notes = JSON.parse(unescapeAttr(json)) || [] } catch (e) { notes = [] }
    return ''
  })
  const number = new Map()
  const count = new Map()
  out = out.replace(/<sup[^>]*\bdata-fn="([^"]*)"[^>]*>\s*<\/sup>/g, (m, raw) => {
    const label = unescapeAttr(raw)
    if (!number.has(label)) number.set(label, number.size + 1)
    const k = count.get(label) || 0
    count.set(label, k + 1)
    const l = escapeAttr(label)
    return `<sup id="fnref${k || ''}:${l}"><a href="#fn:${l}" class="footnote-ref">${number.get(label)}</a></sup>`
  })
  if (!notes.length) return out
  // the empty paragraphs before the notes (the editor's, after the text): none
  out = out.replace(/(?:<p>\s*(?:<br[^>]*>)?\s*<\/p>\s*)+$/, '')
  const ordered = notes.map((n, i) => ({ n, i })).sort((a, b) =>
    ((number.get(a.n.label) || Infinity) - (number.get(b.n.label) || Infinity)) || a.i - b.i).map((x) => x.n)
  const items = ordered.map((n) => {
    const l = escapeAttr(n.label)
    const back = Array.from({ length: count.get(n.label) || 0 }, (_, k) =>
      `&#160;<a href="#fnref${k || ''}:${l}" class="footnote-backref">${backHTML}</a>`).join('')
    let h = String(n.html || '').trim() || '<p></p>'
    h = /<\/p>\s*$/.test(h) ? h.replace(/<\/p>\s*$/, back + '</p>') : h + back
    return `<li id="fn:${l}">${h}</li>`
  })
  return out + `<div class="footnotes"><hr><ol>${items.join('')}</ol></div>`
}
