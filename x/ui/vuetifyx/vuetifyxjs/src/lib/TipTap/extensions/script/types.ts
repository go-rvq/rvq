import { Node } from '@tiptap/core'
import { VueNodeViewRenderer } from '@tiptap/vue-3'
import View from './View.vue'
import { makeInnerTextAsValueAttribute, makeMessagesAttribute } from '../script-block/utils'

export interface InlineTemplateCodeOptions {
  value ?: string
}

declare module '@tiptap/core' {
  interface Commands<ReturnType> {
    script: {
      addScript: () => ReturnType
    }
  }
}


/**
 * This extension allows you to create italic text.
 * @see https://www.tiptap.dev/api/marks/italic
 */
export const InlineTemplateCode = Node.create<InlineTemplateCodeOptions>({
  name: 'script',

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
      tag: `script[type="text/script"]`,
    }]
  },

  renderHTML({ node }) {
    return ['script', { 'type': "text/script" }, node.attrs.value]
  },

  addCommands() {
    return {
      addScript:
        () =>
          ({ commands, editor }) => {
            const node = editor.schema.nodes.script.create({value:''});
            return commands.insertContent(node);
          },
    }
  },

  addNodeView() {
    return VueNodeViewRenderer(View as any)
  },
})
