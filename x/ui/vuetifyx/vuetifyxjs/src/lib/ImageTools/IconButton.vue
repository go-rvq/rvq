<script setup lang="ts">
const props = defineProps({
  title: {
    type: String
  }
})
</script>

<template>
  <v-tooltip :text="title">
    <template v-slot:activator="{ props }">
      <v-btn v-bind="{ ...props, ...$attrs }">
        <template v-for="(_, name) in $slots" v-slot:[name]="slotData"
          ><slot :name="name" v-bind="slotData"
        /></template>
      </v-btn>
    </template>
  </v-tooltip>
</template>

<style scoped></style>
