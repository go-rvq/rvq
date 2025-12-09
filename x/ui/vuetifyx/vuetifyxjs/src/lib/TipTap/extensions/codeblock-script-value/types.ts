import { Node } from '@tiptap/core'
import { VueNodeViewRenderer } from '@tiptap/vue-3'
import View from './View.vue'
import { makeInnerTextAsValueAttribute, makeMessagesAttribute } from '../script-block/utils'

export interface Options {
  value ?: string
}

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    codeBlockScriptValue: {
      addCodeBlockScriptValue: () => ReturnType
    }
  }
}


/**
 * This extension allows you to create italic text.
 * @see https://www.tiptap.dev/api/marks/italic
 */
export const CodeBlockScriptValue = Node.create<Options>({
  name: 'codeBlockScriptValue',

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
      tag: `script[type="text/codeBlockScriptValue"]`,
    }]
  },

  renderHTML({ node }) {
    return ['script', { 'type': "text/codeBlockScriptValue" }, node.attrs.value]
  },

  addCommands() {
    return {
      addCodeBlockScriptValue:
        () =>
          ({ commands, editor }) => {
            const node = editor.schema.nodes.codeBlockScriptValue.create({value:''});
            return commands.insertContent(node);
          },
    }
  },

  addNodeView() {
    return VueNodeViewRenderer(View as any)
  },
})
