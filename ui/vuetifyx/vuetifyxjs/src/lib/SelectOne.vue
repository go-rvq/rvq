<template>
  <label class="v-label theme--light" v-html="label"></label>
  <v-card v-if="internalSelectedItems.length > 0" variant="flat" class="mb-2">
    <v-list>
      <v-list-item
        v-for="element in internalSelectedItems"
        :prepend-avatar="element[itemImage]"
        :title="element[itemText]"
        animation="300"
      >
        <template v-slot:append>
          <v-btn @click="removeItem(element[itemValue])" variant="text" icon="mdi-delete"> </v-btn>
        </template>
      </v-list-item>
    </v-list>
  </v-card>

  <v-autocomplete
    v-if="internalSelectedItems.length == 0"
    :item-value="itemValue"
    :item-title="itemText"
    :items="internalItems"
    :label="addItemLabel"
    v-model="autocompleteValue"
    auto-select-first
    @update:modelValue="addItem"
    @update:search="search"
    @update:focused="focused"
    :loading="isLoading"
    :no-filter="noFilter"
    return-object
    variant="underlined"
  >
    <template v-slot:item="{ props, item }">
      <v-list-item
        v-bind="props"
        :prepend-avatar="item.raw[itemImage]"
        :title="item.raw[itemText]"
      ></v-list-item>
    </template>
  </v-autocomplete>
</template>
<script setup lang="ts">
import draggable from 'vuedraggable'

import { onMounted, Ref, ref } from 'vue'
import { aN } from 'vitest/dist/reporters-P7C2ytIv'

const props = defineProps({
  items: {
    type: Array<object>,
    default: () => []
  },
  selectedItems: {
    type: Array<object>,
    default: () => []
  },
  itemValue: {
    type: String,
    default: 'id'
  },
  itemText: {
    type: String,
    default: 'text'
  },
  itemImage: {
    type: String,
    default: 'image'
  },
  label: {
    type: String,
    default: ''
  },
  addItemLabel: {
    type: String,
    default: ''
  },
  modelValue: {
    type: Array<any>,
    default: []
  },
  searchItemsFunc: {
    type: Function,
    default: null
  }
})
const internalSelectedItems: Ref<any[]> = ref([])
const internalItems: Ref<any[]> = ref([])
const autocompleteValue: Ref<any> = ref([])
const isLoading = ref(false)
const noFilter = ref(false)

const emit = defineEmits(['update:modelValue'])

onMounted(() => {
  internalSelectedItems.value = props.selectedItems
  internalItems.value = props.items
  if (!internalSelectedItems.value || internalSelectedItems.value.length === 0) {
    internalSelectedItems.value = props.modelValue.map((id: any) => {
      return props.items.find((item: any) => item[props.itemValue] === id)
    })
  }
})

// methods
const addItem = (event: any) => {
  autocompleteValue.value = []

  if (
    internalSelectedItems.value.find(
      (element) => element[props.itemValue] == event[props.itemValue]
    )
  ) {
    return
  }

  let newValues: any[] = []
  newValues.push(
    internalItems.value.find((element) => element[props.itemValue] == event[props.itemValue])
  )
  internalSelectedItems.value = newValues
  setValue()
}

const removeItem = (id: string) => {
  internalSelectedItems.value = []
  setValue()
}

const setValue = () => {
  let val = null
  if (internalSelectedItems.value.length > 0) {
    val = internalSelectedItems.value[0][props.itemValue]
  }
  emit('update:modelValue', val)
}

let search: any, focused: any

if (props.searchItemsFunc) {
  const searchDebounce = (func: Function, delay: number) => {
    let timeoutId: number

    return function (val: String) {
      clearTimeout(timeoutId)
      timeoutId = setTimeout(func, delay, val)
    }
  }

  const doSearch = (val: String) => {
    isLoading.value = true
    props
      .searchItemsFunc(val)
      .then((res: any[]) => {
        internalItems.value = res
      })
      .finally(() => {
        isLoading.value = false
      })
  }

  search = searchDebounce(doSearch, 800)
  focused = (val: boolean) => {
    if (val) {
      doSearch('')
    }
  }
}
</script>

<style scoped></style>
