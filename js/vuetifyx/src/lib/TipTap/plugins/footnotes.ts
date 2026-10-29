// The plugin "footnotes" of vx-tiptap-editor: rvq mdkit's footnotes — a
// reference (its number) in the text, the notes at its end —, made, edited,
// renamed, deleted by mdkit's footnotes builder: its button on the bar, a
// reference clicked (its note, the reference to point at another), a note
// clicked; the HTML kept goldmark's (finalizeFootnotes: the labels in the ids).
import {Extension} from '@tiptap/core'
import {Plugin, PluginKey} from '@tiptap/pm/state'
import {FootnoteRef, Footnotes, finalizeFootnotes} from '@rvq/mdkit'
import FootnoteAction from '../components/FootnoteAction.vue'
import type {TiptapPlugin} from './index'

// FootnoteControl is the button of the builder, and the clicks that open it
// (editor.storage.mdkitFootnotes.open: the button's)
export const FootnoteControl = Extension.create({
  name: 'mdkitFootnotes',
  addStorage() {
    return {open: null as null | ((label?: string, pos?: number | null) => void)}
  },
  addOptions() {
    return {
      divider: false,
      spacer: false,
      button: ({editor}: any) => ({component: FootnoteAction, componentProps: {editor}})
    } as any
  },
  addProseMirrorPlugins() {
    const storage = this.storage
    return [new Plugin({
      key: new PluginKey('mdkitFootnoteClicks'),
      props: {
        handleClickOn: (_view, _pos, node, nodePos, event) => {
          if (!storage.open) return false
          if (node.type.name === 'footnoteRef') {
            storage.open(node.attrs.label, nodePos)
            return true
          }
          const li = node.type.name === 'footnotes' && (event.target as HTMLElement)?.closest?.('[data-fn-note]') as HTMLElement | null
          if (li) {
            storage.open(li.dataset.fnNote, null)
            return true
          }
          return false
        }
      }
    })]
  }
})

const footnotes: TiptapPlugin = {
  extensions: () => [FootnoteRef, Footnotes, FootnoteControl],
  output: finalizeFootnotes
}

export default footnotes
