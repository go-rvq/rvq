// The formulas of the rich text editor of a site (tiptap): inline and on a
// line of their own — their HTML tiptap's mathematics' (and x/markup's):
// span[data-type=inline-math], div[data-type=block-math], the LaTeX in
// data-latex —, drawn by KaTeX loaded from where the page's meta
// mdkit-vendor says (not in the bundle); a click opening one
// (editor.storage.mdkitMath.open: the page's dialog, given its LaTeX, place
// and kind).
import { Extension, Node } from '@tiptap/core'
import { load } from './blocks.js'

// drawMath draws latex in el (KaTeX, its styles too); not drawn, its text.
export async function drawMath(el, latex, block) {
  el.textContent = latex
  try {
    load('katex-css', '').catch(() => {})
    const katex = await load('katex', 'katex')
    katex.render(latex, el, { displayMode: block, throwOnError: false })
    return true
  } catch (e) {
    return false
  }
}

const mathNode = (name, block) => Node.create({
  name,
  group: block ? 'block' : 'inline',
  inline: !block,
  atom: true,
  selectable: true,

  addAttributes() {
    return {
      latex: {
        default: '',
        parseHTML: (el) => el.getAttribute('data-latex') || '',
        renderHTML: (attrs) => ({ 'data-latex': attrs.latex }),
      },
    }
  },

  parseHTML() {
    return [{ tag: (block ? 'div' : 'span') + '[data-type="' + (block ? 'block-math' : 'inline-math') + '"]' }]
  },

  renderHTML({ HTMLAttributes }) {
    return [block ? 'div' : 'span', { ...HTMLAttributes, 'data-type': block ? 'block-math' : 'inline-math' }]
  },

  addNodeView() {
    return ({ node, editor, getPos }) => {
      const dom = document.createElement(block ? 'div' : 'span')
      dom.setAttribute('data-type', block ? 'block-math' : 'inline-math')
      dom.setAttribute('data-latex', node.attrs.latex)
      dom.style.cursor = 'pointer'
      if (block) dom.style.textAlign = 'center'
      drawMath(dom, node.attrs.latex, block)
      dom.addEventListener('click', () => {
        const pos = getPos()
        const open = editor.storage.mdkitMath && editor.storage.mdkitMath.open
        if (open && editor.isEditable && typeof pos === 'number') open(node.attrs.latex, pos, block)
      })
      let current = node
      return {
        dom,
        ignoreMutation: () => true,
        update(n) {
          if (n.type !== current.type) return false
          if (n.attrs.latex !== current.attrs.latex) drawMath(dom, n.attrs.latex, block)
          current = n
          return true
        },
      }
    }
  },
})

export const InlineMath = mathNode('inlineMath', false)
export const BlockMath = mathNode('blockMath', true)

// MathCommands are the commands of the formulas — insertMath (at the cursor)
// and updateMath (the one at pos: another kind, in its place) — and the
// storage of the dialog that opens one (mdkitMath).
export const MathCommands = Extension.create({
  name: 'mdkitMath',

  addStorage() {
    return { open: null }
  },

  addCommands() {
    const type = (schema, block) => schema.nodes[block ? 'blockMath' : 'inlineMath']
    return {
      insertMath: (latex, block) => ({ state, commands }) =>
        !!type(state.schema, block) && commands.insertContent({ type: block ? 'blockMath' : 'inlineMath', attrs: { latex } }),
      updateMath: (pos, latex, block) => ({ tr, state, dispatch }) => {
        const node = tr.doc.nodeAt(pos)
        if (!node || (node.type.name !== 'inlineMath' && node.type.name !== 'blockMath')) return false
        if (dispatch) tr.replaceWith(pos, pos + node.nodeSize, type(state.schema, block).create({ latex }))
        return true
      },
    }
  },
})
