<script setup lang="ts">
import { FilterItem } from '@/lib/Filter/Model'
import FilterButton from '@/lib/Filter/components/FilterButton.vue'
import { ref } from 'vue'

interface WithBeforeDoneHandler {
  beforeDoneHandler: (cb:(empty:boolean)=>void)=>void
}

const props = defineProps<{
  modelValue: FilterItem
  isFoldedItem: boolean
  itemComp: any
  translations: any
  compTranslations: any
  internalValue: any
  index: number
  config: any
}>()


let beforeDoneHandler:(cb:(empty:boolean)=>void)=>void

const value = ref({ ...props.modelValue })
const menu = ref(false)
const emit = defineEmits(['update:modelValue', 'change', 'clear'])

const save = (empty:boolean) => {
  if (empty) {
    clear(new Event("save"))
    return
  }
  if (
    !value.value.valueIs &&
    (!value.value.valuesAre || value.value.valuesAre.length == 0) &&
    !value.value.valueFrom &&
    !value.value.valueTo
  ) {
    return
  }
  value.value.selected = true
  Object.assign(props.modelValue, value.value)
  emit('update:modelValue', props.modelValue)
  emit('change', null)
}

const clickDone = () => {
  menu.value = false

  if (beforeDoneHandler) {
    beforeDoneHandler(save)
  } else {
    save(false)
  }
}

const clear = (e: any) => {
  if (!value.value.selected) {
    return
  }
  value.value.valueIs = ''
  value.value.valuesAre = []
  value.value.valueFrom = ''
  value.value.valueTo = ''
  value.value.selected = false
  Object.assign(props.modelValue, value.value)
  emit('update:modelValue', props.modelValue)
  emit('clear', e)
}
</script>

<template>
  <v-menu :close-on-content-click="false" class="rounded-lg" v-model="menu">
    <template v-slot:activator="{ props }">
      <filter-button
        :comp-translations="compTranslations"
        :op="value"
        :is-folded-item="isFoldedItem"
        :slotProps="props"
        @clear="clear"
      />
    </template>
    <v-card class="pa-3 bg-white">
      <div>{{ modelValue.translations?.filterBy }}</div>
      <component ref="itemEl"
        v-model="value"
        :is="itemComp"
        :translations="compTranslations"
        :config="config"
        @postSetup="(o:WithBeforeDoneHandler) => beforeDoneHandler = o.beforeDoneHandler"
      ></component>
      <div>
        <v-btn class="mt-5 float-right" color="primary" rounded @click="clickDone">{{
          translations.apply
        }}</v-btn>
      </div>
    </v-card>
  </v-menu>
</template>

<style scoped></style>
