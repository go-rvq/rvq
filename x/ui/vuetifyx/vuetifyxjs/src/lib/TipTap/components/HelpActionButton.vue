<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3'
import { mdiHelpCircle } from '@mdi/js'

import { ActionButton } from 'vuetify-pro-tiptap'
import { useI18n } from '@/lib/locale'

const { t } = useI18n()

const props = withDefaults(defineProps<Props>(), {
  modelValue: false
})

interface Props {
  editor: Editor
  tooltip?: string
  disabled?: boolean
  title?: string
  component: any
  modelValue?:boolean
}

const visible = defineModel<boolean>()

function action(v:any) {
  const newValue = !visible.value
  visible.value = newValue
}

function isActive() {
  return visible.value as boolean
}
</script>

<template>
  <ActionButton :rawIcon="`svg:${mdiHelpCircle}`" :editor="editor" :tooltip="t('help')" :disabled="disabled" :action="action" :is-active="isActive">
  </ActionButton>
</template>
