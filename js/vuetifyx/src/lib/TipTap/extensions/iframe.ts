import {Node, mergeAttributes} from '@tiptap/core'

// Iframe is a block, atomic node so TipTap keeps `<iframe>` embeds (Google Maps,
// generic embeds) in the content instead of stripping them — anything not in the
// schema is dropped when content is parsed. Trust is enforced on SAVE by the
// server-side sanitizer (only allowlisted hosts survive), so keeping the node in
// the editor is safe.
export interface IframeOptions {
  HTMLAttributes: Record<string, any>
}

export const Iframe = Node.create<IframeOptions>({
  name: 'iframe',
  group: 'block',
  atom: true,
  selectable: true,
  draggable: false,

  addOptions() {
    return {HTMLAttributes: {class: 'vx-tiptap-iframe'}}
  },

  addAttributes() {
    return {
      src: {default: null},
      width: {default: null},
      height: {default: null},
      frameborder: {default: '0'},
      allow: {default: null},
      allowfullscreen: {default: null},
      loading: {default: null},
      referrerpolicy: {default: null},
      style: {default: null},
    }
  },

  parseHTML() {
    return [{tag: 'iframe'}]
  },

  renderHTML({HTMLAttributes}) {
    return ['iframe', mergeAttributes(this.options.HTMLAttributes, HTMLAttributes)]
  },
})

export default Iframe
