// The buttons and the tabs of the rich text editor (tiptap): blocks whole
// (atoms), their value an attribute — read from their HTML and written as it
// (widgets-model.js) —, shown as the page shows them (widgets.js), a click
// opening the dialog of one (editor.storage.mdkitWidgets.open: the page's,
// given its kind, value and place).
import { Node } from '@tiptap/core'
import { buttonsHTML, buttonsOf, tabsHTML, tabsOf } from './widgets-model.js'
import { drawButtons, drawTabs } from './widgets.js'
import { vendorOf } from './blocks.js'

// element is the element of html (one).
function element(html) {
  const t = document.createElement('template')
  t.innerHTML = html
  return t.content.firstElementChild
}

const widget = (name, kind, html, of) => Node.create({
  name,
  group: 'block',
  atom: true,
  draggable: true,
  selectable: true,

  addStorage() {
    return { open: null }
  },

  addAttributes() {
    return { value: { default: { align: 'start', items: [] }, rendered: false } }
  },

  parseHTML() {
    return [{ tag: 'div[data-type="' + kind + '"]', priority: 60, getAttrs: (el) => ({ value: of(el) }) }]
  },

  renderHTML({ node }) {
    return element(html(node.attrs.value))
  },

  addNodeView() {
    return ({ node, editor, getPos }) => {
      const dom = document.createElement('div')
      dom.className = 'mdkit-widget-edit'
      dom.setAttribute('data-widget-edit', kind)
      dom.contentEditable = 'false'
      const draw = (n) => {
        dom.replaceChildren(element(html(n.attrs.value)))
        if (kind === 'tabs') drawTabs(dom, vendorOf(), true)
        else drawButtons(dom, vendorOf(), undefined, true)
      }
      draw(node)
      dom.addEventListener('click', (e) => {
        // a tab chosen: only that
        if (e.target.closest && e.target.closest('.mdkit-tabs-bar')) return
        e.preventDefault()
        const open = editor.storage[name] && editor.storage[name].open
        const pos = getPos()
        if (open && editor.isEditable && typeof pos === 'number') open(current.attrs.value, pos)
      })
      let current = node
      return {
        dom,
        ignoreMutation: () => true,
        stopEvent: (e) => e.type === 'click' || e.type === 'mousedown',
        update(n) {
          if (n.type !== current.type) return false
          if (n.attrs.value !== current.attrs.value) draw(n)
          current = n
          return true
        },
      }
    }
  },

  addCommands() {
    const cap = name[0].toUpperCase() + name.slice(1)
    return {
      // insertMdkitButtons, insertMdkitTabs: at the cursor — after a block
      // selected (a widget just put: not in its place) —, a paragraph after
      // it at the end of the text (to write on)
      ['insert' + cap]: (value) => ({ state, commands }) => {
        const { selection, doc } = state
        const content = [{ type: name, attrs: { value } }]
        const at = selection.node ? selection.to : null
        const end = at ?? (selection.$to.depth >= 1 ? selection.$to.after(1) : selection.to)
        if (end >= doc.content.size) content.push({ type: 'paragraph' })
        return at != null ? commands.insertContentAt(at, content) : commands.insertContent(content)
      },
      // updateMdkitButtons, updateMdkitTabs: the one at pos, null taken out
      ['update' + cap]: (pos, value) => ({ tr, dispatch }) => {
        const n = tr.doc.nodeAt(pos)
        if (!n || n.type.name !== name) return false
        if (dispatch) {
          if (value) tr.setNodeMarkup(pos, undefined, { ...n.attrs, value })
          else tr.delete(pos, pos + n.nodeSize)
        }
        return true
      },
    }
  },
})

// style of a widget in an editor: hovered, selected
function editStyle() {
  if (typeof document === 'undefined' || document.getElementById('mdkit-widget-edit-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-widget-edit-style'
  st.textContent = '.mdkit-widget-edit{margin:.5em 0;padding:2px 6px;border:1px dashed rgba(127,127,127,.4);border-radius:6px;cursor:pointer}' +
    '.mdkit-widget-edit:hover{border-color:#3b82f6}.mdkit-widget-edit.ProseMirror-selectednode{outline:2px solid #3b82f6}' +
    '.mdkit-widget-edit a[data-type="button"]{pointer-events:none}'
  document.head.append(st)
}
if (typeof document !== 'undefined') editStyle()

export const Buttons = widget('mdkitButtons', 'buttons', buttonsHTML, buttonsOf)
export const Tabs = widget('mdkitTabs', 'tabs', tabsHTML, tabsOf)
