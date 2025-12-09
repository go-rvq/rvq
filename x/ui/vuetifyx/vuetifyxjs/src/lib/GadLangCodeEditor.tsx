import { computed, defineComponent, ExtractPublicPropTypes, PropType, useModel, nextTick } from 'vue'

// import Prism Editor
import { PrismEditor, EditorProps } from 'vue-prism-editor'
import 'vue-prism-editor/dist/prismeditor.min.css'

import Prism from 'prismjs'
const { highlight, languages } = Prism

import 'prismjs/components/prism-clike';
import 'prismjs/components/prism-gad';
import 'prismjs/themes/prism-coy.css'; // import syntax highlighting styles
import "./styles/GadLangCodeEditor.scss"

import { useI18n } from '@/lib/locale'

const propsOptions = {
  modelValue: {
    type: String,
    required: true
  },

  singleLine: Boolean,
  noPlaceholder: Boolean,

  editorProps: {
    type: Object as PropType<EditorProps>,
    default: () => ({
      lineNumbers: true,
    }) as EditorProps
  }
} as const

export type Props = ExtractPublicPropTypes<typeof propsOptions>

const closeChars:Record<string, string> = {
  '(': ')',
  '[': ']',
  '{': '}',
  "'": "'",
  '"': '"'
}

export default defineComponent({
  name: 'VXGadLangCodeEditor',

  props: propsOptions,

  emits: {
    'update:modelValue': (value:string) => true
  },

  setup(props:Props) {
    const
      { t } = useI18n(),
      code = useModel(props, 'modelValue'),
      highlighter = (code:string) => highlight(code, languages.gad, 'gad'),
      editorProps = { ...props.editorProps},

      handleInput = async (event:any) => {
      const { target } = event;
      const selectionStart = target.selectionStart as number, value = target.value as string;
      const char = value.slice(selectionStart - 1, selectionStart)

      if (closeChars[char]) {
        const newValue = value.slice(0, selectionStart) + closeChars[char] + value.slice(selectionStart);
        code.value = newValue; // Update the v-model for PrismEditor
        // Use nextTick to ensure DOM is updated before setting selectionEnd
        return nextTick().then(() => {
          target.selectionEnd = selectionStart;
        });
      }
      return new Promise((resolve)=>resolve(''))
    }

    if (props.singleLine) {
      editorProps.lineNumbers = false
    }

    if (!props.noPlaceholder) {
      editorProps.placeholder = t('pleaseEnterTheCode')
    }

    return () => <div class={['gadlang-editor language-gad', {
      'gadlang-editor--single-line': props.singleLine || false
    }]}>
      <PrismEditor
        v-model={code.value}
      {...{
        highlight: highlighter,
        onInput: handleInput,
        ...editorProps
      }}
    ></PrismEditor></div>
  }
})
