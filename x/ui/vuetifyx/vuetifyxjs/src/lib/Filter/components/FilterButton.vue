<script setup lang="ts">
import { FilterItem } from '@/lib/Filter/Model'
import FilterButtonBody from '@/lib/Filter/components/FilterButtonBody.vue'

const props = defineProps<{
  op: FilterItem
  isFoldedItem: boolean
  slotProps: any
  compTranslations: any
}>()

const emit = defineEmits(['clear'])
const clear = (e: any) => {
  emit('clear', e)
}
</script>

<template>
  <v-list-item v-if="isFoldedItem" variant="outlined" class="my-1 px-2" v-bind="slotProps">
    <filter-button-body :comp-translations="compTranslations" :op="op" @clear="clear" />
  </v-list-item>
  <v-chip
    v-else
    variant="outlined"
    class="mr-2 my-1"
    :class="op.selected ? '' : 'text-grey-darken-1'"
    :style="{ borderStyle: op.selected ? 'solid' : 'dashed' }"
    v-bind="slotProps"
  >
    <filter-button-body :comp-translations="compTranslations" :op="op" @clear="clear" />
  </v-chip>
</template>

<style scoped></style>
