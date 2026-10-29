// The plugin "tables" of vx-tiptap-editor: rvq mdkit's table builder (a
// sheet) on the table under the cursor, or a new one (TableAction).
import {Extension} from '@tiptap/core'
import {Plugin, PluginKey} from '@tiptap/pm/state'
import {Decoration, DecorationSet} from '@tiptap/pm/view'
import TableAction from '../components/TableAction.vue'
import type {TiptapPlugin} from './index'

export const TableBuilderButton = Extension.create({
  name: 'mdkitTables',
  // the layout of a table (the builder's): fluid, or not and placed — its
  // class mdkit-fit mdkit-<position> kept by the table node
  addGlobalAttributes() {
    return [{
      types: ['table'],
      attributes: {
        mdkitLayout: {
          default: null,
          parseHTML: (el: HTMLElement) => (el.getAttribute('class') || '').split(/\s+/).filter((c) => /^mdkit-(fit|left|center|right)$/.test(c)).join(' ') || null,
          renderHTML: (attrs: any) => (attrs.mdkitLayout ? {class: attrs.mdkitLayout} : {})
        }
      }
    }]
  },
  addOptions() {
    return {
      divider: false,
      spacer: false,
      button: ({editor}: any) => ({component: TableAction, componentProps: {editor}})
    } as any
  },
  // the layout seen in the editor: tiptap's view of a table puts no class on
  // it — its class on the view (a decoration), the table inside placed by it
  addProseMirrorPlugins() {
    return [new Plugin({
      key: new PluginKey('mdkitTableLayout'),
      props: {
        // (any: the types of ProseMirror the tsc sees, not the alias of vite)
        decorations(state): any {
          const decos: Decoration[] = []
          state.doc.descendants((node, pos) => {
            if (node.type.name === 'table' && node.attrs.mdkitLayout) {
              decos.push(Decoration.node(pos, pos + node.nodeSize, {class: node.attrs.mdkitLayout}))
            }
          })
          return DecorationSet.create(state.doc, decos)
        }
      }
    })]
  }
})

const tablesPlugin: TiptapPlugin = {
  extensions: () => [TableBuilderButton]
}

export default tablesPlugin
