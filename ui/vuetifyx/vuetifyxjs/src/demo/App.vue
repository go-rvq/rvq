<template>
  <v-app>
    <v-container>
      <EditorJS v-model="body" placeholder="Insert body now"> </EditorJS>
      <TreeDataTable
        :items="items"
        :title="(v: any) => v.text"
        :headers="headers"
        @apply-settings="() => console.log('settings changed')"
      ></TreeDataTable>
    </v-container>
  </v-app>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { default as TreeDataTable } from '@/lib/TreeDataTable.vue'
import EditorJS from '@/lib/EditorJS/EditorJS.vue'

const body = 'hello'

const headers = [
  {
    title: 'V1',
    key: 'value'
  },
  {
    title: 'V2',
    key: 'value2'
  }
]

const items = ref([
  {
    id: 1,
    text: 'hello',
    value: 2,
    value2: 5
  }
])
</script>
