<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3'
import { mdiEye } from '@mdi/js'
import { useTheme } from 'vuetify'
import Dialog from '@/lib/Dialog.vue'
import { useI18n } from '@/lib/locale'

import { ActionButton } from 'vuetify-pro-tiptap'

withDefaults(defineProps<Props>(), {
  disabled: false
})

const { global: { name: originalTheme } } = useTheme()

interface Props {
  editor: Editor
  disabled?: boolean
  component?:any
  componentProps?:Object
}

const { t } = useI18n()
const title = t('preview')
</script>

<template>
  <Dialog expandable closable :title="title" density="compact" v-if="component">
    <template #activator="{props: activatorProps}">
      <ActionButton :editor="editor" :tooltip="title" :disabled="disabled" v-bind="activatorProps">
        <VIcon>{{ `svg:${mdiEye}` }}</VIcon>
      </ActionButton>
    </template>
    <template #body>
      <component :is="component" :editor="editor" :t="t" v-bind="componentProps || {}" />
    </template>
  </Dialog>
</template>
