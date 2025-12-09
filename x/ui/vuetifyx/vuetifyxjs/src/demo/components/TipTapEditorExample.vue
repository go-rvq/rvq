<script setup lang="ts">

import { computed, onMounted, ref, watch } from 'vue'
import { default as VXTipTapEditor } from '@/lib/TipTap/VXTipTapEditor'
// import { highlight, languages } from 'prismjs/components/prism-core'

// import { languages } from 'prismjs/components'
// import { highlight } from 'prismjs'

import Prism from 'prismjs'
const { highlight, languages } = Prism

import { PrismEditor } from 'vue-prism-editor'
import 'prismjs/components/prism-markup'
import 'prismjs/components/prism-json'

import { VuetifyViewer } from 'vuetify-pro-tiptap'
declare const window: any

const readonly = ref(false)
const withLabel = ref(true)
const withTemplate = ref(true)
const withHint = ref(true)
const withErrorMessages = ref(false)
const withError = ref(false)
const initialContent = ref<any>('')
const checkContent = (t: boolean) => {
  initialContent.value = t ? `
<pre><code>a<`+`script type="text/scriptValue">123456789abcdefghijklmno</`+`script></code></pre>
<p><`+`script type="text/script">ssd</`+`script></p>
`
    : `
<p>Example</p>
`
}


const
  editor = ref<InstanceType<typeof VXTipTapEditor> | null>(null),
  htmlContent = ref(''),
  jsonContent = ref(''),
  renderContents = () => {
    const e = (editor.value as any).vuetifyTiptap.editor
    htmlContent.value = e.getHTML()
    jsonContent.value = JSON.stringify(e.getJSON(), null, 4)
  },
  contentModel = computed({
    get() {
      return initialContent.value
    },

    set(v) {
      initialContent.value = v
      renderContents()
    }
  })

watch(withTemplate, (v) => {
  checkContent(v)
})


onMounted(() => {
  renderContents()
})

checkContent(withTemplate.value)

</script>

<template>
  <VRow>
    <VCol>
      <VCheckbox v-model="withTemplate" label="Template Script" />
    </VCol>
    <VCol>
      <VCheckbox v-model="readonly" label="ReadOnly" />
    </VCol>
    <VCol>
      <VCheckbox v-model="withLabel" label="Show Label" />
    </VCol>
    <VCol>
      <VCheckbox v-model="withHint" label="Show Hint" />
    </VCol>
    <VCol>
      <VCheckbox v-model="withErrorMessages" label="Show Error Messages" />
    </VCol>
    <VCol>
      <VCheckbox v-model="withError" label="Force Error State" />
    </VCol>
  </VRow>
  <VXTipTapEditor
    :template="withTemplate"
    :label="withLabel ? 'Page Content': undefined"
    :hint="withHint ? 'The content of page': undefined"
    :error-messages="withErrorMessages ? ['Bad content value', 'Another error reason']: undefined"
    :error="withError"
    class="mb-3"
    :readonly="readonly"
    output="html"
    v-model="contentModel"
    ref="editor">
    <template #help>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message</p>
      <p>Help Message w</p>
    </template>
  </VXTipTapEditor>
  <VRow>
    <VCol>
      <VCard title="HTML">
        <VCardText  style="overflow-y: auto; max-height: 440px">
          <PrismEditor readonly line-numbers v-model="htmlContent"
                       :highlight="(code:string) => highlight(code, languages.markup, 'markup')" />
        </VCardText>
      </VCard>
    </VCol>
    <VCol>
      <VCard title="JSON">
        <VCardText style="overflow-y: auto; max-height: 440px">
          <PrismEditor readonly line-numbers v-model="jsonContent"
                       :highlight="(code:string) => highlight(code, languages.json, 'json')" />
        </VCardText>
      </VCard>
    </VCol>
  </VRow>
</template>
