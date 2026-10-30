// The diagram of the rich text editor (tiptap): a block, whole (atom), its
// Mermaid source; its HTML the one x/markup keeps of a ```mermaid fence —
// pre > code.language-mermaid —, read back from it (before the code block);
// shown drawn (Mermaid, loaded from where the page's meta mdkit-vendor
// says), a click opening it (editor.storage.mdkitDiagram.open: the page's
// dialog, given its source and place).
import { Node } from '@tiptap/core'
import { load } from './blocks.js'

let drawn = 0

// drawInto draws source in el (an svg); one that does not draw, its source
// as code (and false).
export async function drawInto(el, source) {
  try {
    const mermaid = await load('mermaid', 'mermaid')
    if (!mermaid.__mdkit) {
      mermaid.initialize({ startOnLoad: false, securityLevel: 'strict' })
      mermaid.__mdkit = true
    }
    const { svg } = await mermaid.render('mdkit-diagram-edit-' + (++drawn), source)
    el.innerHTML = svg
    return true
  } catch (e) {
    const pre = document.createElement('pre')
    pre.className = 'mdkit-diagram-error'
    pre.textContent = source
    el.replaceChildren(pre)
    return false
  }
}

// style is the look of a diagram in an editor, once for the page.
function style() {
  if (typeof document === 'undefined' || document.getElementById('mdkit-diagram-edit-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-diagram-edit-style'
  st.textContent = '.mdkit-diagram-edit{margin:.75em 0;padding:.5rem;border:1px dashed rgba(128,128,128,.4);border-radius:6px;cursor:pointer;text-align:center}' +
    '.mdkit-diagram-edit:hover{border-color:#3b82f6}' +
    '.mdkit-diagram-edit.ProseMirror-selectednode{outline:2px solid #3b82f6}' +
    '.mdkit-diagram-edit svg{max-width:100%;height:auto}' +
    '.mdkit-diagram-error{text-align:left;color:#b91c1c;white-space:pre-wrap}'
  document.head.append(st)
}

export const Diagram = Node.create({
  name: 'mdkitDiagram',
  group: 'block',
  atom: true,
  draggable: true,
  selectable: true,

  addStorage() {
    return { open: null }
  },

  addAttributes() {
    return { source: { default: '' } }
  },

  parseHTML() {
    return [{
      tag: 'pre',
      // before the code block's
      priority: 60,
      preserveWhitespace: 'full',
      getAttrs: (el) => {
        const code = el.querySelector('code.language-mermaid')
        return code ? { source: code.textContent || '' } : false
      },
    }, {
      // a diagram drawn (the page's): its source
      tag: 'div.mdkit-diagram[data-source]',
      getAttrs: (el) => ({ source: el.getAttribute('data-source') || '' }),
    }]
  },

  renderHTML({ node }) {
    return ['pre', {}, ['code', { class: 'language-mermaid' }, node.attrs.source]]
  },

  addCommands() {
    return {
      // insertDiagram puts a diagram of source at the cursor
      insertDiagram: (source) => ({ commands }) => commands.insertContent({ type: this.name, attrs: { source } }),
      // updateDiagram makes the diagram at pos of source
      updateDiagram: (pos, source) => ({ tr, dispatch }) => {
        const node = tr.doc.nodeAt(pos)
        if (!node || node.type.name !== this.name) return false
        if (dispatch) tr.setNodeMarkup(pos, undefined, { ...node.attrs, source })
        return true
      },
    }
  },

  addNodeView() {
    return ({ node, editor, getPos }) => {
      style()
      const dom = document.createElement('div')
      dom.className = 'mdkit-diagram-edit'
      dom.setAttribute('data-diagram-edit', '')
      drawInto(dom, node.attrs.source)
      dom.addEventListener('click', () => {
        const open = editor.storage.mdkitDiagram && editor.storage.mdkitDiagram.open
        const pos = getPos()
        if (open && editor.isEditable && typeof pos === 'number') open(node.attrs.source, pos)
      })
      let current = node
      return {
        dom,
        ignoreMutation: () => true,
        update(n) {
          if (n.type !== current.type) return false
          if (n.attrs.source !== current.attrs.source) drawInto(dom, n.attrs.source)
          current = n
          return true
        },
      }
    }
  },
})
