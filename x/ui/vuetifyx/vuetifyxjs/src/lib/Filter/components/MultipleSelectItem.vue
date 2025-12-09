<script setup lang="ts">
import * as Vue from 'vue'
import { ref } from 'vue'

declare var window: any

class Item {
  text: string = ''
  value: string = ''
}

class Model {
  options: Item[] = []
  valuesAre: Item[] = []
  modifier: String = 'in'
}

interface CompArgs {
  window: any
  Vue: any
}

interface Config {
  component?: ((scope: CompArgs) => void)
}

const props = defineProps<{
  translations: any
  config: Config
}>()

const t = props.translations
const config = props.config || {}

const model = defineModel<Model>()

if (!model.value) {
  model.value = new Model()
}

if (!model.value.options) {
  model.value.options = []
}

if (!model.value.modifier) {
  model.value.modifier = 'in'
}

const items = ref([
  { text: t.in, value: 'in' },
  { text: t.notIn, value: 'notIn' }
])

const itemSelectorComp = config.component?.({
  window,
  Vue
})

</script>

<template>
  <div>
    <div>
      <v-select
        class="d-inline-block"
        style="width: 200px"
        v-model="(model as Model).modifier"
        :items="items"
        item-title="text"
        item-value="value"
        variant="underlined"
        hide-details
      ></v-select>
    </div>
    <itemSelectorComp
      v-if="itemSelectorComp"
      v-model="(model as Model).valuesAre"
      :items="(model as Model).options"
      item-title="text"
      item-value="value" />
    <v-select
      v-else
      chips
      v-model="(model as Model).valuesAre"
      :items="(model as Model).options"
      item-title="text"
      item-value="value"
      multiple
      density="comfortable"
    />
  </div>
</template>
