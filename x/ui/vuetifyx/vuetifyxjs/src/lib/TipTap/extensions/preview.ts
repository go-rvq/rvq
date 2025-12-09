import type { GeneralOptions } from 'vuetify-pro-tiptap'
import { Extension } from '@tiptap/core'

import PreviewActionButton from '../components/PreviewActionButton.vue'

export interface PreviewOptions extends GeneralOptions<PreviewOptions> {
  title?: string
  component?: any
}

export default Extension.create<PreviewOptions>({
  name: 'preview',
  addOptions() {
    return {
      divider: false,
      spacer: false,
      title: 'Preview',
      button: ({ editor, extension }) => ({
        component: PreviewActionButton,
        componentProps: {
          editor,
          component: extension.options.component
        }
      })
    }
  }
})
