import { GeneralOptions } from 'vuetify-pro-tiptap'
import Button from './Button.vue'
import { ScriptValue as ExtTemplateCode, ScriptValueOptions as ExtTemplateCodeOptions } from './types'

export interface TemplateCodeValueOptions extends ExtTemplateCodeOptions, GeneralOptions<TemplateCodeValueOptions> {
}

export const TemplateCodeValue = /* @__PURE__*/ ExtTemplateCode.extend<TemplateCodeValueOptions>({
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
