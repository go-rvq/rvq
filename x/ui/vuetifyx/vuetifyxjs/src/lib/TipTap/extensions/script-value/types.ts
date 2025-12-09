import { Node } from '@tiptap/core'
import { VueNodeViewRenderer } from '@tiptap/vue-3'
import View from './View.vue'
import { makeInnerTextAsValueAttribute, makeMessagesAttribute } from '../script-block/utils'

export interface ScriptValueOptions {
  /**
   * HTML attributes to add to the italic element.
   * @default {}
   * @example { class: 'foo' }
   */
  HTMLAttributes: Record<string, any>
}

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    scriptValue: {
      addScriptValue: () => ReturnType
    }
  }
}


/**
 * This extension allows you to create italic text.
 * @see https://www.tiptap.dev/api/marks/italic
 */
export const ScriptValue = Node.create<ScriptValueOptions>({
  name: 'scriptValue',

  inline: true,
  group: 'inline',

  content: 'text*',


  addAttributes() {
    return {
      ...makeInnerTextAsValueAttribute(),
      ...makeMessagesAttribute()
    }
  },

  parseHTML() {
    return [{
      tag: `script[type="text/scriptValue"]`
    }]
  },

  renderHTML({ node }) {
    const attrs = { ...node.attrs }
    delete attrs.messages
    return ['script', { 'type': 'text/scriptValue' }, attrs.value]
  },

  addCommands() {
    return {
      addScriptValue:
        () =>
          ({ commands, editor }) => {
            const node = editor.schema.nodes.scriptValue.create({ value: '' })
            return commands.insertContent(node)
          }
    }
  },

  addNodeView() {
    return VueNodeViewRenderer(View as any)
  }
})
