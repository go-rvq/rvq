// The plugin "tables" of vx-tiptap-editor: rvq mdkit's table builder (a
// sheet) on the table under the cursor, or a new one (TableAction).
import {Extension} from '@tiptap/core'
import TableAction from '../components/TableAction.vue'
import type {TiptapPlugin} from './index'

export const TableBuilderButton = Extension.create({
  name: 'mdkitTables',
  addOptions() {
    return {
      divider: false,
      spacer: false,
      button: ({editor}: any) => ({component: TableAction, componentProps: {editor}})
    } as any
  }
})

const tablesPlugin: TiptapPlugin = {
  extensions: () => [TableBuilderButton]
}

export default tablesPlugin
