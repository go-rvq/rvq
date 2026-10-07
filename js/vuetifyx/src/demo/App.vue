<script setup lang="ts">
import { ref } from 'vue'
import DiffBrowserExample from '@/demo/components/DiffBrowserExample.vue'
import TipTapEditorExample from '@/demo/components/TipTapEditorExample.vue'

// the examples, a tab each
const tab = ref('diff-browser')
</script>
<template>
  <VApp>
    <VMain>
      <VContainer>
        <VTabs v-model="tab" class="mb-4" color="primary">
          <VTab value="diff-browser">Diff browser</VTab>
          <VTab value="tiptap">TipTap editor</VTab>
        </VTabs>
        <DiffBrowserExample v-if="tab === 'diff-browser'" />
        <TipTapEditorExample v-else />
      </VContainer>
    </VMain>
  </VApp>
</template>
