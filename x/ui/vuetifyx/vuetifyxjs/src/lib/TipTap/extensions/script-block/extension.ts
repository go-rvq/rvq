import { Editor } from '@tiptap/vue-3'
import { mdiCodeBlockBraces } from '@mdi/js'
import { type ActionMenuButtonItem, type GeneralOptions, type DisabledCheckOptions } from 'vuetify-pro-tiptap'

export const ScriptBlockMenuItemBuilder = (editor: Editor, t:(path:string)=>string) => ({
  action: (() => editor.chain().focus().setScriptBlock().run()),
  isActive: () => editor.isActive('scriptBlock') || false,
  disabled: !editor.can().toggleTemplateCodeBlock(),
  title: t('script'),
  rawIcon: `svg:${mdiCodeBlockBraces}`
} as ActionMenuButtonItem)

