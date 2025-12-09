<script lang="ts" setup>
import {ActionButton, extActionButtonProps} from 'vuetify-pro-tiptap'
import { mdiCodeBraces } from '@mdi/js'
import { useI18n } from '@/lib/locale'

defineProps(extActionButtonProps)
const { t } = useI18n()
</script>

<template>
  <ActionButton
    :editor="editor"
    :action="() => editor.chain().focus().addScript().run()"
    :is-active="() => editor.isActive('script')"
    :disabled="editor.isActive('scriptBlock')"
    :tooltip="t('script')"
    :raw-icon="`svg:${mdiCodeBraces}`"
  />
</template>

<!-- <style lang="scss" scoped></style> -->
