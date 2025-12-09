import type { GeneralOptions } from 'vuetify-pro-tiptap'
import { Extension } from '@tiptap/core'

import Button from './Button.vue'
import {ToolsItem} from './types'

export interface ToolsOptions extends GeneralOptions<ToolsOptions> {
  component?: Function
  items?:ToolsItem[]
}

export const ToolsButton = Extension.create<ToolsOptions>({
  name: 'toolsButton',
  addOptions() {
    return {
      divider: false,
      spacer: false,
      title: 'Tools',
      button: ({ editor, extension }) => ({
        component: Button,
        componentProps: {
          editor,
          component: extension.options.component,
          items: extension.options.items
        }
      })
    } as ToolsOptions
  }
})
