// The plugin "tabs" of vx-tiptap-editor: rvq mdkit's tabs — each its title,
// color, icons and content, the content written in this editor (its
// plugins, but tabs) —, made and changed in their dialog (TabsAction); tabs
// clicked, the same on them.
import {Extension} from '@tiptap/core'
import {Tabs} from '@rvq/mdkit'
import TabsAction from '../components/TabsAction.vue'
import type {TiptapPlugin} from './index'

const tabs: TiptapPlugin = {
  extensions: ({plugins}) => [Tabs, Extension.create({
    name: 'mdkitTabsButton',
    addOptions() {
      return {
        divider: false,
        spacer: false,
        button: ({editor}: any) => ({
          component: TabsAction,
          componentProps: {editor, plugins: (plugins || []).filter((p) => p !== 'tabs')}
        })
      } as any
    }
  })]
}

export default tabs
