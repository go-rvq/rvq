import { GeneralOptions } from 'vuetify-pro-tiptap'
import Button from './Button.vue'
import {
  CodeBlockScriptValue as ExtCodeBlockScriptValue,
  Options as ExtOptions
} from './types'

export interface Options extends ExtOptions, GeneralOptions<Options> {
}

export const CodeBlockScriptValue = /* @__PURE__*/ ExtCodeBlockScriptValue.extend<Options>({
  addOptions() {
    return {
      HTMLAttributes:{},
      button: ({ editor, extension, t }) => {
        return {
          component: Button,
          componentProps: {
            ...extension,
            editor,
            t
          }
        }
      }
    }
  }
})
