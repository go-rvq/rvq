<template>
  <div class="v-input vx-input-editorjs" style="display: block">
    <!---->
    <div class="v-input__control" style="display: block">
      <label v-if="label" class="v-label">{{ label }}</label>
      <div
        class="editorjs v-card--variant-outlined v-card--variant-outlined"
        ref="htmlelement"
      ></div>
    </div>
    <div class="v-input__details" v-if="!hideDetails && errorMessages.length">
      <div class="v-messages" role="alert" aria-live="polite">
        <div class="v-messages__message" v-for="m in errorMessages">{{ m }}</div>
      </div>
    </div>
  </div>
</template>
<style scoped>
.v-input.vx-input-editorjs {
  display: block !important;
}

.v-input__control {
  display: block;
}

.vx-input-editorjs :deep(.codex-editor--narrow .ce-toolbar__actions) {
  right: auto !important;
  left: -55px;
}
.vx-input-editorjs :deep(.codex-editor__redactor) {
  margin-right: 0;
}
</style>
<script setup>
import EditorJS from '@editorjs/editorjs'
import EmbedTool from '@editorjs/embed'
import ListTool from '@editorjs/list'
import ImageTool from '@editorjs/image'
import HeaderTool from '@editorjs/header'
import Checklist from '@editorjs/checklist'

import Quote from '@editorjs/quote'
import VideoTool from './video.js'

import Marker from '@editorjs/marker'
import InlineCode from '@editorjs/inline-code'
import Underline from '@editorjs/underline'
import ChangeCase from 'editorjs-change-case'
import Strikethrough from '@sotaproject/strikethrough'
import Undo from 'editorjs-undo'
import NestedCheckList from '@calumk/editorjs-nested-checklist'
import Table from '@editorjs/table'
import Warning from '@editorjs/warning'
import CodeTool from '@editorjs/code'
import RawTool from '@editorjs/raw'
import Delimiter from '@editorjs/delimiter'

// TODO: import LinkAutocomplete from '@editorjs/link-autocomplete';
// TODO: import Personality from '@editorjs/personality'
// import Paragraph   from "@editorjs/paragraph"; // no change on view
// import TextVariantTune from '@editorjs/text-variant-tune'; // bug view
import { onMounted, onUnmounted, ref, watch } from 'vue'

const htmlelement = ref(null),
  props = defineProps({
    modelValue: { type: String },
    placeholder: { type: String },
    label: { type: String },
    hideDetails: { type: Boolean, default: false },
    errorMessages: {
      type: Array,
      default: []
    }
  }),
  emit = defineEmits(['update:modelValue'])

let editor,
  ready = false,
  updatingModel = false

// model -> view
function modelToView() {
  if (!editor || !ready) {
    return
  }
  if (!props.modelValue) {
    return
  }
  if (typeof props.modelValue === 'string') {
    const s = props.modelValue.toString()
    if (s[0] === '{' && s[s.length - 1] === '}') {
      const v = new Function('return ' + s).call()
      editor.render(v)
    } else {
      editor.blocks.renderFromHTML(props.modelValue)
    }
  }
  new Undo({ editor })
}

// view -> model
function viewToModel(api, event) {
  updatingModel = true
  editor
    .save()
    .then((outputData) => {
      emit('update:modelValue', JSON.stringify(outputData))
    })
    .catch((error) => {
      console.log(event, 'Saving failed: ', error)
    })
    .finally(() => {
      updatingModel = false
    })
}

onMounted(() => {
  editor = new EditorJS({
    holder: htmlelement.value,
    placeholder: props.placeholder,
    tools: {
      //paragraph: Paragraph,
      embed: EmbedTool,
      list: ListTool,
      image: ImageTool,
      video: VideoTool,
      header: {
        class: HeaderTool,
        config: {
          placeholder: 'Enter a header',
          levels: [1, 2, 3, 4, 5, 6],
          defaultLevel: 2
        }
      },
      quote: Quote,
      checklist: {
        class: Checklist,
        inlineToolbar: true
      },
      Marker: {
        class: Marker,
        shortcut: 'CMD+SHIFT+M'
      },
      inlineCode: InlineCode,
      underline: Underline,
      changeCase: {
        class: ChangeCase //,
        //config: {
        //  showLocaleOption: true, // enable locale case options
        // locale: 'tr' // or ['tr', 'TR', 'tr-TR']
        //}
      },
      strikethrough: Strikethrough,
      nestedCheckList: NestedCheckList,
      table: Table,
      warning: Warning,
      code: CodeTool,
      raw: RawTool,
      delimiter: Delimiter
    },
    minHeight: 'auto',
    data: props.modelValue,
    onChange: viewToModel
  })

  editor.isReady.then(() => {
    ready = true
    modelToView()
  })
})
watch(
  () => props.modelValue,
  () => {
    if (!updatingModel) {
      modelToView()
    }
  }
)
onUnmounted(() => {
  editor.destroy()
})
</script>
