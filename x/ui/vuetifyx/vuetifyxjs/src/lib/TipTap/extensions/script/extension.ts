import { GeneralOptions } from 'vuetify-pro-tiptap'
import Button from './Button.vue'
import { InlineTemplateCode as ExtInlineTemplateCode, InlineTemplateCodeOptions as ExtInlineTemplateCodeOptions } from './types'

export interface InlineTemplateCodeOptions extends ExtInlineTemplateCodeOptions, GeneralOptions<InlineTemplateCodeOptions> {
}

export const InlineTemplateCode = /* @__PURE__*/ ExtInlineTemplateCode.extend<InlineTemplateCodeOptions>({
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
