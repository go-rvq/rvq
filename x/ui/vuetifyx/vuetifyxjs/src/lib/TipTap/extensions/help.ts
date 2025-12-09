import type { GeneralOptions } from 'vuetify-pro-tiptap'
import { Extension } from '@tiptap/core'

import HelpActionButton from '../components/HelpActionButton.vue'
import { Ref } from 'vue'

export interface HelpOptions extends GeneralOptions<HelpOptions> {
  modelValue?:Ref<boolean>
}

export default Extension.create<HelpOptions>({
  name: 'help',
  addOptions() {
    return {
      divider: false,
      spacer: false,
      button: ({ editor, extension }) => {
        return {
          component: HelpActionButton,
          componentProps: {
            editor,
            modelValue: extension.options.modelValue?.value,
            'onUpdate:modelValue': (v:boolean) => {
              if (extension.options.modelValue) {
                extension.options.modelValue.value = v
              }
            }
          }
        }
      }
    }
  }
})
