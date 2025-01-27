<script setup lang="ts">
const props = defineProps<{
  modelValue: any
  translations: any
}>()

props.modelValue.modifier = props.modelValue.modifier || 'equals'

const emit = defineEmits(['update:modelValue'])
const modifier = props.modelValue.modifier
const t = props.translations
</script>

<template>
  <div style="width: 300px">
    <div v-if="modifier == 'between'">
      <v-container>
        <v-row>
          <v-icon
            class="py-8 mr-4 float-md-start"
            icon="mdi-subdirectory-arrow-right"
            size="large"
          ></v-icon>
          <v-text-field
            class="d-inline-block"
            style="width: 80px"
            type="number"
            variant="underlined"
            v-model="props.modelValue.valueFrom"
          />
          <div class="pt-6 px-3">{{ t.and }}</div>
          <v-text-field
            class="d-inline-block"
            style="width: 80px"
            type="number"
            variant="underlined"
            v-model="props.modelValue.valueTo"
          />
        </v-row>
      </v-container>
    </div>
    <div v-else>
      <v-icon
        class="py-8 mr-4 float-md-start"
        icon="mdi-subdirectory-arrow-right"
        size="large"
      ></v-icon>
      <v-text-field
        class="d-inline-block"
        style="width: 120px"
        type="number"
        variant="underlined"
        v-model="props.modelValue.valueIs"
      />
    </div>
  </div>
</template>

<style scoped></style>
