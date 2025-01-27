<script setup lang="ts">
import { ref } from 'vue'
import * as constants from '../Constants'
import Monthpicker from '@/lib/Monthpicker.vue'

const props = defineProps<{
  modelValue: any
  translations: any
  config: any
}>()

props.modelValue.modifier = props.modelValue.modifier || constants.ModifierBetween

const datePickerVisible = ref(false)

const emit = defineEmits(['update:modelValue'])
const modifier = props.modelValue.modifier
</script>

<template>
  <div style="width: 200px">
    <monthpicker
      v-model="modelValue.valueFrom"
      :key="modifier + 'form'"
      :visible="datePickerVisible"
      :hide-details="true"
      :year-text="translations['year']"
      :month-text="translations['month']"
      :ok-text="translations['apply']"
      :clear-text="translations['clear']"
      :month-names="translations['monthNames']"
      :min="config?.min"
      :max="modelValue.valueTo"
    />
    <div style="height: 34px" class="pl-2 pt-4">
      <span>{{ translations['to'] }}</span>
    </div>
    <monthpicker
      v-model="modelValue.valueTo"
      :key="modifier + 'to'"
      :hide-details="true"
      :year-text="translations['year']"
      :month-text="translations['month']"
      :ok-text="translations['apply']"
      :clear-text="translations['clear']"
      :month-names="translations['monthNames']"
      :min="modelValue.valueFrom"
      :max="config?.max"
    />
  </div>
</template>

<style scoped></style>
