<script setup lang="ts">
import {ref} from 'vue'

// template builds the text inserted at the cursor from a clicked tag. It
// defaults to the `{{tag}}` form; a caller using another syntax (e.g. the
// single-brace `{expr}` SEO templates) passes its own builder.
const props = defineProps<{
  template?: (tag: string) => string
}>()

const vnode = ref()

const tagInputsFocus = (v: any) => {
  vnode.value = v
}
const addTags = (tag: any) => {
  if (!vnode.value) {
    return
  }
  let lazyValue = vnode.value.modelValue
  let selectionStart = vnode.value.selectionStart
  let selectionEnd = vnode.value.selectionEnd
  const input = vnode.value.$el.querySelector('input, textarea')
  if (input) {
    selectionStart = input.selectionStart
    selectionEnd = input.selectionEnd
  }
  let startString = lazyValue.substring(0, selectionStart)
  let endString = lazyValue.substring(selectionEnd, lazyValue.length)

  const text = props.template ? props.template(tag) : '{{' + tag + '}}'
  vnode.value.$emit('update:modelValue', startString + text + endString)
  vnode.value.focus()
}
defineExpose({
  tagInputsFocus,
  addTags
})
</script>

<template>
  <div>
    <slot></slot>
  </div>
</template>

<style scoped></style>
