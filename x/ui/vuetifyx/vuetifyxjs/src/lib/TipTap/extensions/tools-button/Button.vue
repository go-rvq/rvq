<script lang="ts" setup>
import {  extActionButtonProps } from 'vuetify-pro-tiptap'
import { useI18n } from '@/lib/locale'
import { PropType, ref } from 'vue'
import { Editor } from '@tiptap/vue-3'
import { CommandArgs, ToolsItem } from '@/lib/TipTap/extensions/tools-button/types'
import { mdiTools } from '@mdi/js'
import {Message} from './types'

const props = defineProps({
  items:{
    type: Object as PropType<ToolsItem[]>,
    default: () => []
  },
  ... extActionButtonProps
})
const editor: Editor = props.editor

const { t } = useI18n()
const messages = ref<Message[]>([])

function notify(msg: string, color: string = 'success') {
  messages.value.push({ text: msg, color: color } as Message)
}

const commandArgs = {editor, notify} as CommandArgs
</script>

<template>
  <VSnackbarQueue
    v-model="messages"
    location="top"
    timeout="3000"
    closable
  />
  <VMenu v-if="items?.length">
    <template #activator="{props}">
      <VBtn v-bind="props" icon density="comfortable" size="small" class="rounded me-1 ms-0">
        <VIcon :icon="`svg:`+mdiTools" />
        <VTooltip :eager="false" activator="parent" location="top" :text="t('tools')" />
      </VBtn>
    </template>
    <VList>
      <template v-for="(item, index) in items"
                 :key="index"
                 :value="index">
        <component v-if="item.component" :is="item.component?.view" v-bind="item.component.props" :t="t" :editor="editor" :notify="notify" />
        <VListItem v-else :prepend-icon="item.icon" :title="item.title?.()" @click="item.action?.(commandArgs)" />
      </template>
    </VList>
  </VMenu>
</template>
