<script setup lang="ts">
// RecordPickerItem filters by records of a model chosen in the dialog of the
// AdvancedSelect: the search is the records event of the model's listing
// (config.url). The value goes as a list (key.in=a,b), one record or many;
// the options are the records chosen, their labels.
import AdvancedSelect from '@/lib/AdvancedSelect.vue'
import {computed} from 'vue'

interface Option {
  text: string
  value: string
}

interface Texts {
  select?: string
  change?: string
  search?: string
  noItems?: string
  close?: string
}

interface Config {
  // the listing of the model: its records event answers the search
  url: string
  event?: string
  // more of the query of the search (the encoder of the titles, say)
  query?: Record<string, string>
  itemValue?: string
  itemText?: string
  many?: boolean
  texts?: Texts
}

const props = defineProps<{
  translations: any
  config: Config
}>()

const model = defineModel<any>()
const config: Config = props.config || ({ url: '' } as Config)
const itemValue = config.itemValue || 'ID'
const itemText = config.itemText || 'Text'
const texts: Texts = config.texts || {}

if (!model.value.valuesAre) {
  model.value.valuesAre = []
}
if (!model.value.options) {
  model.value.options = []
}
model.value.modifier = 'in'

// the records known: the options, as the AdvancedSelect reads them
const items = computed(() =>
  (model.value.options || []).map((o: Option) => ({ [itemValue]: o.value, [itemText]: o.text }))
)

const value = computed({
  get: () => (config.many ? model.value.valuesAre : model.value.valuesAre?.[0] ?? ''),
  set: (v: any) => {
    model.value.valuesAre = Array.isArray(v) ? v.map(String) : v ? [String(v)] : []
  }
})

const selected = (records: any[]) => {
  model.value.options = (records || [])
    .filter((r: any) => r)
    .map((r: any) => ({ value: String(r[itemValue]), text: String(r[itemText] ?? r[itemValue]) }))
}

const searchPage = (keyword: string, page: number) => {
  const u = new URL(config.url, window.location.href)
  u.searchParams.set('__execute_event__', config.event || 'presets_ListData')
  u.searchParams.set('keyword', keyword || '')
  u.searchParams.set('page', String(page || 1))
  for (const k in config.query || {}) {
    u.searchParams.set(k, config.query![k])
  }
  return fetch(u.toString(), {
    credentials: 'same-origin',
    // the page asking, not the browser navigating (web.PlaidRequestHeader)
    headers: { 'X-Plaid-Request': '1' }
  })
    .then((r) => r.json())
    .then((r) => r?.data || {})
}
</script>

<template>
  <div style="min-width: 260px">
    <advanced-select
      v-model="value"
      :items="items"
      :item-value="itemValue"
      :item-text="itemText"
      :many="!!config.many"
      :search-page-func="searchPage"
      dialog
      chips
      :dialog-title="model.label"
      dialog-on-top
      hide-details
      :select-text="texts.select"
      :change-text="texts.change"
      :search-text="texts.search"
      :no-items-text="texts.noItems"
      :close-text="texts.close"
      @update:selected-items="selected"
    />
  </div>
</template>
