<script lang="ts" setup>
import {ActionButton, extActionButtonProps} from 'vuetify-pro-tiptap'
import { mdiEqual } from '@mdi/js'
import { useI18n } from '@/lib/locale'

const { t } = useI18n()
const { editor } = defineProps(extActionButtonProps)

function action() {
  return editor.chain().focus().addScriptValue().run()
}
</script>

<template>
  <ActionButton
    :editor="editor"
    :action="action"
    :is-active="() => editor.isActive('scriptValue') ||  false"
    :disabled="editor.isActive('scriptBlock')"
    :tooltip="t('scriptResult')"
    :raw-icon="`svg:${mdiEqual}`"
  />
</template>

<!-- <style lang="scss" scoped></style> -->
