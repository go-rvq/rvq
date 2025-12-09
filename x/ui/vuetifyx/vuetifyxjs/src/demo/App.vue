<script lang="ts">
import { defineComponent } from 'vue'
import TipTapEditorExample from '@/demo/components/TipTapEditorExample.vue'

export default defineComponent({
  components: { TipTapEditorExample }
})
</script>
<template>
  <VApp>
    <VMain>
      <VContainer>
        <TipTapEditorExample />
      </VContainer>
    </VMain>
  </VApp>
</template>
