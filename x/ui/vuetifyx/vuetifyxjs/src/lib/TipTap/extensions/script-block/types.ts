import { Node } from '@tiptap/core'
import { VueNodeViewRenderer } from '@tiptap/vue-3'
import { type GeneralOptions, type DisabledCheckOptions } from 'vuetify-pro-tiptap'
import View from './View.vue'
import { makeAttributes } from './utils'

export interface ScriptBlockOptions extends GeneralOptions<ScriptBlockOptions>, DisabledCheckOptions {
  value ?: string
}

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    scriptBlock: {
      setScriptBlock: () => ReturnType
      toggleTemplateCodeBlock: () => ReturnType
    }
  }
}

/**
 * This extension allows you to create paragraphs.
 * @see https://www.tiptap.dev/api/nodes/paragraph
 */
export const ScriptBlock = Node.create<ScriptBlockOptions>({
  name: 'scriptBlock',

  priority: 1000,

  addAttributes() {
    return makeAttributes()
  },

  marks: '',

  group: 'block',

  content: 'inline*',

  draggable: true,

  code: true,

  whitespace: 'pre',

  parseHTML() {
    return [{
      tag: `script[type="text/scriptBlock"]`,
    }]
  },

  renderHTML({node }) {
    return ['script', { 'type': "text/scriptBlock" }, node.attrs.value]
  },

  addNodeView() {
    return VueNodeViewRenderer(View as any)
  },

  addCommands() {
    return {
      setScriptBlock: () => ({ commands }) => commands.setNode(this.name),
      toggleTemplateCodeBlock: () => ({ commands }) => commands.toggleNode?.(this.name, 'paragraph', {})
    }
  }
})
