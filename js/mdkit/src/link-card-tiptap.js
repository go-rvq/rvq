// The link card of the rich text editor (tiptap): a block, whole (atom), its
// target (href: a path of the site, a record path of the admin); its HTML the
// one x/markup makes of `[[target]]` — div.mdkit-link-card[data-href] —,
// read back from it; shown as its card (link-cards draws it), with a bin
// that takes it out. Pasting an address of the site alone (a path, an
// address of this site) puts a card in its place (Ctrl+Z: the text pasted).
import { Extension, Node } from '@tiptap/core'
import { Plugin, PluginKey } from '@tiptap/pm/state'
import { closeHistory } from '@tiptap/pm/history'
import { cardClass, cardTarget } from './link-cards.js'

// style is the look of a card in an editor, once for the page.
function style() {
  if (typeof document === 'undefined' || document.getElementById('mdkit-link-card-edit-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-link-card-edit-style'
  st.textContent = '.mdkit-link-card-edit{position:relative;margin:.75em 0}' +
    '.mdkit-link-card-edit>.' + cardClass + '{margin:0;pointer-events:none}' +
    '.mdkit-link-card-edit.ProseMirror-selectednode>.' + cardClass + '{outline:2px solid #3b82f6;border-radius:8px}' +
    '.mdkit-link-card-bin{position:absolute;top:.35rem;left:min(36rem,100%);transform:translateX(-2.1rem);width:1.7rem;height:1.7rem;border:0;' +
    'border-radius:50%;background:rgba(255,255,255,.9);color:#b91c1c;cursor:pointer;font-size:1rem;line-height:1.7rem;box-shadow:0 1px 3px rgba(0,0,0,.25)}' +
    '.mdkit-link-card-bin:hover{background:#fee2e2}'
  document.head.append(st)
}

export const LinkCard = Node.create({
  name: 'linkCard',
  group: 'block',
  atom: true,
  draggable: true,
  selectable: true,

  addOptions() {
    return {
      // removeTitle: the bin's title (the words of the page)
      removeTitle: 'Remove the card',
      // paste: an address pasted alone made a card
      paste: true,
    }
  },

  addAttributes() {
    return { href: { default: '' } }
  },

  parseHTML() {
    return [{
      tag: 'div.' + cardClass + '[data-href]',
      getAttrs: (el) => (cardTarget(el.getAttribute('data-href')) ? { href: el.getAttribute('data-href') } : false),
    }]
  },

  renderHTML({ node }) {
    const href = node.attrs.href
    return ['div', { class: cardClass, 'data-href': href }, href.startsWith('/') ? ['a', { href }, href] : href]
  },

  addCommands() {
    return {
      // insertLinkCard puts a card of href at the cursor (in place of the
      // text chosen; after a block chosen whole — a card —, not over it)
      insertLinkCard: (href) => ({ state, commands }) => {
        const target = cardTarget(href)
        if (!target) return false
        const card = { type: this.name, attrs: { href: target } }
        const { selection } = state
        return selection.node ? commands.insertContentAt(selection.to, card) : commands.insertContent(card)
      },
    }
  },

  addNodeView() {
    return ({ node, editor, getPos }) => {
      style()
      const dom = document.createElement('div')
      dom.className = 'mdkit-link-card-edit'
      dom.setAttribute('data-link-card-edit', node.attrs.href)
      const card = document.createElement('div')
      card.className = cardClass
      card.setAttribute('data-href', node.attrs.href)
      card.textContent = node.attrs.href
      dom.append(card)
      if (editor.isEditable) {
        const bin = document.createElement('button')
        bin.type = 'button'
        bin.className = 'mdkit-link-card-bin'
        bin.title = this.options.removeTitle
        bin.setAttribute('data-link-card-remove', '')
        bin.textContent = '🗑'
        bin.addEventListener('mousedown', (e) => e.preventDefault())
        bin.addEventListener('click', () => {
          const pos = getPos()
          if (typeof pos === 'number') editor.chain().focus().deleteRange({ from: pos, to: pos + node.nodeSize }).run()
        })
        dom.append(bin)
      }
      return {
        dom,
        // the card drawn in it (link-cards) is not the editor's
        ignoreMutation: () => true,
        stopEvent: (e) => !!(e.target && e.target.closest && e.target.closest('.mdkit-link-card-bin')),
      }
    }
  },

  // its paste: an extension of its own, before the link's
  addExtensions() {
    return this.options.paste ? [LinkCardPaste] : []
  },
})

// LinkCardPaste makes an address of the site pasted alone a card: before the
// link's paste (1000: an address made a link) — a priority of its own, not
// the card's (the first block of the schema fills an empty text).
export const LinkCardPaste = Extension.create({
  name: 'linkCardPaste',
  priority: 1100,

  addProseMirrorPlugins() {
    const type = this.editor.schema.nodes.linkCard
    if (!type) return []
    return [new Plugin({
      key: new PluginKey('mdkitLinkCardPaste'),
      props: {
        // an address of the site alone, pasted: its card
        handlePaste(view, event, slice) {
          const text = (event && event.clipboardData && event.clipboardData.getData('text/plain')) ||
            (slice && slice.content.textBetween(0, slice.content.size, '\n'))
          if (!text || /\n/.test(text.trim())) return false
          const href = cardTarget(text)
          if (!href || !href.startsWith('/') && !href.startsWith('admin:')) return false
          // the text pasted, then its card — a step of its own: Ctrl+Z, the text
          const t = text.trim()
          const from = view.state.selection.from
          view.dispatch(view.state.tr.insertText(t))
          // (the paragraph whole, when it was the text pasted alone)
          const tr = view.state.tr.replaceRangeWith(from, from + t.length, type.create({ href }))
          view.dispatch(closeHistory(tr).scrollIntoView())
          return true
        },
      },
    })]
  },
})
