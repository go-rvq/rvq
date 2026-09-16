import type {GeneralOptions} from 'vuetify-pro-tiptap'
import {Extension} from '@tiptap/core'

import SourceCodeActionButton from '../components/SourceCodeActionButton.vue'

export interface SourceCodeOptions extends GeneralOptions<SourceCodeOptions> {}

// sourcecode adds a toolbar button that opens an expandable dialog with a
// CodeMirror editor for the content's HTML source (replacing the inline "code"
// mark button). Applying the dialog sets the editor content from the HTML.
export default Extension.create<SourceCodeOptions>({
  name: 'sourceCode',
  addOptions() {
    return {
      divider: false,
      spacer: false,
      button: ({editor}) => {
        return {
          component: SourceCodeActionButton,
          componentProps: {
            editor
          }
        }
      }
    }
  }
})
