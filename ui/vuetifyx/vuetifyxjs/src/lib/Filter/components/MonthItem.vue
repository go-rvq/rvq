<script setup lang="ts">
import { ref } from 'vue'
import Monthpicker from '@/lib/Monthpicker.vue'

const props = defineProps<{
  modelValue: any
  translations: any
  min: string
  max: string
}>()

const monthPickerVisible = ref(false)

const emit = defineEmits(['update:modelValue'])
</script>

<template>
  <div style="width: 150px">
    <monthpicker
      v-model="modelValue.valueIs"
      :visible="monthPickerVisible"
      :hide-details="true"
      :year-text="translations['year']"
      :month-text="translations['month']"
      :ok-text="translations['apply']"
      :clear-text="translations['clear']"
      :month-names="translations['monthNames']"
      :min="min"
      :max="max"
    />
  </div>
</template>

<style scoped></style>
