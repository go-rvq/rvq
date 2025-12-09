<script setup lang="ts">
import { ref } from 'vue'
import VXGadLangCodeEditor from '@/lib/GadLangCodeEditor'
const code = ref('const a = 1')
</script>

<template>
  <VXGadLangCodeEditor single-line v-model="code" line-numberes></VXGadLangCodeEditor>
  <VTextarea v-model="code" />
</template>

<style scoped>

</style>
