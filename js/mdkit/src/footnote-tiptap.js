// The footnotes of the rich text editor (tiptap): two nodes of ours —
//   - footnoteRef: in a line, whole (atom), its label; shown as its number
//     (the order of the first reference to its note: its node view);
//   - footnotes: a block, whole, the notes of the text [{label, html}];
//     shown at its end as a list in the order of the references.
// Their HTML in the editor (renderHTML) is made the HTML kept by
// footnote-rich finalize; both read back (parseHTML) the HTML kept —
// goldmark's: sup > a.footnote-ref, div.footnotes > ol > li#fn:label — and
// their own.
import { Node, mergeAttributes } from '@tiptap/core'
import { noteHTML } from './footnote-rich.js'

const labelOf = (s) => (/fn(?:ref\d*)?:(.+)$/.exec(s || '') || [])[1] || ''

// style is the look of the footnotes in an editor, once for the page: a
// reference its number; the notes at the end, to click.
function style() {
  if (typeof document === 'undefined' || document.getElementById('mdkit-footnote-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-footnote-style'
  st.textContent = 'sup.mdkit-fn-ref{cursor:pointer;padding:0 .1em}' +
    'sup.mdkit-fn-ref::after{content:attr(data-number);color:#2563eb;font-weight:600}' +
    'sup.mdkit-fn-ref.ProseMirror-selectednode{outline:1px solid #3b82f6;border-radius:2px}' +
    '.mdkit-fn-notes{font-size:.875em;color:#4b5563;margin-top:1em}' +
    '.mdkit-fn-notes hr{border:0;border-top:1px solid #d1d5db;margin:0 0 .5em;width:30%}' +
    '.mdkit-fn-notes ol{list-style:decimal;padding-left:1.5em}' +
    '.mdkit-fn-notes li{cursor:pointer}' +
    '.mdkit-fn-notes li:hover{background:#f3f4f6}' +
    '.mdkit-fn-notes li p{margin:0}'
  document.head.append(st)
}

// numbers are the numbers of the labels of doc: the order of their first
// references.
export function numbers(doc) {
  const out = new Map()
  doc.descendants((node) => {
    if (node.type.name === 'footnoteRef' && !out.has(node.attrs.label)) out.set(node.attrs.label, out.size + 1)
  })
  return out
}

export const FootnoteRef = Node.create({
  name: 'footnoteRef',
  group: 'inline',
  inline: true,
  atom: true,
  selectable: true,
  addAttributes() {
    // its label: in data-fn (renderHTML), not an attribute of its own
    return { label: { default: '', renderHTML: () => ({}) } }
  },
  parseHTML() {
    return [
      { tag: 'sup[data-fn]', getAttrs: (el) => ({ label: el.getAttribute('data-fn') }) },
      {
        tag: 'sup',
        priority: 60,
        getAttrs: (el) => {
          const a = el.querySelector('a.footnote-ref')
          const label = a && labelOf(a.getAttribute('href'))
          return label ? { label } : false
        },
      },
    ]
  },
  renderHTML({ node, HTMLAttributes }) {
    return ['sup', mergeAttributes(HTMLAttributes, { 'data-fn': node.attrs.label, class: 'mdkit-fn-ref' })]
  },
  // its number (its place among the others) drawn by a view of its own,
  // again at each change of the text — no class of the editor's view: any
  // editor's, its copy of ProseMirror whichever —
  addNodeView() {
    return ({ editor, node }) => {
      style()
      let current = node
      const dom = document.createElement('sup')
      dom.className = 'mdkit-fn-ref'
      dom.contentEditable = 'false'
      const draw = () => {
        const n = numbers(editor.state.doc).get(current.attrs.label)
        dom.dataset.fn = current.attrs.label
        dom.dataset.number = n ? String(n) : ''
        dom.title = '[^' + current.attrs.label + ']'
      }
      draw()
      editor.on('transaction', draw)
      return {
        dom,
        update(updated) {
          if (updated.type !== current.type) return false
          current = updated
          draw()
          return true
        },
        destroy() { editor.off('transaction', draw) },
        ignoreMutation: () => true,
      }
    }
  },
})

export const Footnotes = Node.create({
  name: 'footnotes',
  group: 'block',
  atom: true,
  selectable: false,
  draggable: false,
  addAttributes() {
    return { notes: { default: [] } }
  },
  parseHTML() {
    return [
      {
        tag: 'div[data-fn-notes]',
        getAttrs: (el) => {
          try { return { notes: JSON.parse(el.getAttribute('data-fn-notes') || '[]') } } catch (e) { return { notes: [] } }
        },
      },
      {
        tag: 'div.footnotes',
        priority: 60,
        getAttrs: (el) => ({
          notes: [...el.querySelectorAll('li')].map((li) => ({ label: labelOf(li.id), html: noteHTML(li.innerHTML) }))
            .filter((n) => n.label),
        }),
      },
    ]
  },
  renderHTML({ node }) {
    return ['div', { 'data-fn-notes': JSON.stringify(node.attrs.notes || []) }]
  },
  // the list of the notes, in the order of the references, drawn again at
  // each change of the text
  addNodeView() {
    return ({ editor, node }) => {
      style()
      let current = node
      const dom = document.createElement('div')
      dom.className = 'mdkit-fn-notes footnotes'
      dom.contentEditable = 'false'
      const draw = () => {
        const nums = numbers(editor.state.doc)
        const notes = [...(current.attrs.notes || [])].map((n, i) => ({ n, i }))
          .sort((a, b) => ((nums.get(a.n.label) || Infinity) - (nums.get(b.n.label) || Infinity)) || a.i - b.i)
        dom.innerHTML = notes.length ? '<hr><ol></ol>' : ''
        const ol = dom.querySelector('ol')
        for (const { n } of notes) {
          const li = document.createElement('li')
          li.dataset.fnNote = n.label
          li.title = '[^' + n.label + ']'
          li.innerHTML = n.html || ''
          ol.append(li)
        }
      }
      draw()
      editor.on('transaction', draw)
      return {
        dom,
        update(updated) {
          if (updated.type !== current.type) return false
          current = updated
          draw()
          return true
        },
        destroy() { editor.off('transaction', draw) },
        ignoreMutation: () => true,
      }
    }
  },
})
