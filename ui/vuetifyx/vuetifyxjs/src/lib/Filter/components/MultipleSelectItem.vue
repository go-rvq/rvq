<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  translations: any
}>()

const t = props.translations

class Item {
  text: string = ''
  value: string = ''
}

class Model {
  options: Item[] = []
  valuesAre: Item[] = []
  modifier: String = 'in'
}

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
    <v-select
      chips
      v-model="(model as Model).valuesAre"
      :items="(model as Model).options"
      item-title="text"
      item-value="value"
      multiple
      density="comfortable"
    ></v-select>
  </div>
</template>
