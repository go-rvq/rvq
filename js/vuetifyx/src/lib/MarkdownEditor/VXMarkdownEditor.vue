<script setup lang="ts">
// vx-markdown-editor: a text in Markdown in the admin — rvq mdkit's editor
// (tables, footnotes, files included, pictures, preview, help) with its
// label, hint and errors as a field's, its words the page's vue-i18n's (vx
// Mdkit…). preview and upload: rvq admin/presets/fields/markdown's events (a
// function), or addresses.
import {computed} from 'vue'
import {MarkdownEditor, configure, textStyle} from '@rvq/mdkit'
import {useI18n} from '@/lib/locale'
import {mdkitWords} from '@/lib/locale/mdkit'

const props = withDefaults(defineProps<{
  modelValue?: string
  label?: string
  hint?: string
  errorMessages?: string[] | string
  rows?: number | string
  disabled?: boolean
  placeholder?: string
  preview?: string | ((text: string) => Promise<string>)
  upload?: string | ((body: FormData) => Promise<any>)
}>(), {modelValue: '', label: '', hint: '', errorMessages: () => [], rows: 8, disabled: false, placeholder: '', preview: '', upload: ''})

const emit = defineEmits<{ (e: 'update:modelValue', v: string): void }>()
const {t} = useI18n()
const words = mdkitWords(t)
configure({words})
textStyle()

const errors = computed(() => (Array.isArray(props.errorMessages) ? props.errorMessages : props.errorMessages ? [props.errorMessages] : []))
</script>

<template>
  <div class="vx-markdown-editor mb-4" :class="{'vx-markdown-editor--error': errors.length}" data-vx-markdown-editor>
    <div v-if="label" class="text-caption text-medium-emphasis mb-1">{{ label }}</div>
    <MarkdownEditor :model-value="modelValue" @update:model-value="(v: string) => emit('update:modelValue', v)"
                    :words="words.markdown" :bar-words="words.rich" :rows="rows" :disabled="disabled"
                    :placeholder="placeholder" :invalid="errors.length > 0" :preview="preview" :upload="upload"/>
    <div v-for="e in errors" :key="e" class="text-caption text-error mt-1">{{ e }}</div>
    <div v-if="hint && !errors.length" class="text-caption text-medium-emphasis mt-1">{{ hint }}</div>
  </div>
</template>
