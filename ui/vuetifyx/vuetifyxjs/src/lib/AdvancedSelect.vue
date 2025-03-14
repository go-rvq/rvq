<template>
  <div class="vx-advanced-select" :class="'vx-advanced-select__' + (many ? 'many' : 'one')">
    <label
      v-if="label"
      class="v-label theme--light"
      v-html="label"
      :style="errorMessages.length ? 'color: rgb(var(--v-theme-error));opacity: inherit;' : ''"
    ></label>
    <v-card
      :color="errorMessages.length ? 'error' : ''"
      v-if="internalSelectedItems.length > 0"
      variant="flat"
      class="mb-2 vx-advanced-select"
    >
      <v-chip-group v-if="chips">
        <v-chip v-for="item in internalSelectedItems">
          {{ item[itemText] }}
          <template v-slot:append>
            <v-btn
              density="compact"
              variant="text"
              icon="mdi-delete"
              @click="removeItem(item[itemValue])"
            ></v-btn>
          </template>
        </v-chip>
      </v-chip-group>

      <v-list v-else density="compact" style="padding: 0">
        <draggable
          v-model="internalSelectedItems"
          :item-key="itemValue"
          @dragstart="startDragItems = internalSelectedItems"
          @change="changeOrder"
          handle=".handle"
        >
          <template #item="{ element }">
            <v-list-item
              :prepend-avatar="element[itemImage]"
              :title="element[itemText]"
              animation="300"
            >
              <template v-slot:prepend>
                <v-icon
                  v-if="many"
                  density="compact"
                  icon="mdi-drag"
                  class="handle mx-2 cursor-grab"
                ></v-icon>
                <v-btn
                  @click="removeItem(element[itemValue])"
                  variant="text"
                  icon="mdi-delete"
                  density="compact"
                ></v-btn>
              </template>
            </v-list-item>
          </template>
        </draggable>
      </v-list>
    </v-card>

    <v-autocomplete
      v-if="selecting"
      :item-value="itemValue"
      :item-title="itemText"
      :items="internalItems"
      :label="addItemLabel"
      v-model="autocompleteValue"
      @update:modelValue="addItem"
      @update:search="setSearchTerm"
      @update:focused="focused"
      @keyup.enter="entered"
      :loading="isLoading"
      :no-filter="noFilter"
      :error-messages="errorMessages"
      return-object
      variant="underlined"
      density="compact"
      :style="autoComplete?.style"
      :class="autoComplete?.class"
    >
      <template v-slot:item="{ props, item }">
        <v-list-item
          v-bind="props"
          :prepend-avatar="item.raw[itemImage]"
          :title="item.raw[itemText]"
        ></v-list-item>
      </template>
    </v-autocomplete>

    <div
      class="v-messages mb-3"
      v-if="!selecting && errorMessages.length"
      aria-live="polite"
      style="opacity: 1; color: rgb(var(--v-theme-error))"
    >
      <div style="height: 5px; border-top: 1px solid rgb(var(--v-theme-error))"></div>
      <div role="alert" v-for="(item, i) in errorMessages" class="v-messages__message">
        {{ item }}
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import draggable from 'vuedraggable'

import {isReactive, onMounted, ref, toRaw, watch} from 'vue'

const props = defineProps({
  items: {
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
  modelValue: {},
  searchItemsFunc: {
    type: Function,
    default: null
  },
  many: {
    type: Boolean,
    default: false
  },
  errorMessages: {
    type: Array<string>,
    default: []
  },
  chips: {
    type: Boolean,
    default: false
  },
  autoComplete: {
    type: Object,
    default: {}
  }
})

const modelValue = defineModel()

const startDragItems = ref<any[]>([])
const internalSelectedItems = ref<any[]>([])
const internalItems = ref<any[]>([])
const autocompleteValue = ref<any[]>([])
const searchTerm = ref<string>('')
const isLoading = ref(false)
const noFilter = ref(false)
const selecting = ref(false)
const updating = { v: false }

const emit = defineEmits()

const modelValueArray = (modelValue: any): any[] => {
  if (!modelValue) {
    return []
  }

  const rawValue = modelValue
  const rawModelValue = isReactive(rawValue) ? toRaw(rawValue) : rawValue
  let rawModelValueArray: any[] = Array.isArray(rawModelValue)
    ? (rawModelValue as any[])
    : [rawModelValue as any]
  return rawModelValueArray
}

const initValue = (modelValue: any) => {
  const oldItems = internalSelectedItems.value as any[]
  let items: any[] = []

  if (modelValue) {
    const rawModelValueArray = modelValueArray(modelValue)

    if (internalItems.value.length) {
      items = rawModelValueArray.map((id: any) => {
        return internalItems.value.find(
          (item: any) => item[props.itemValue].toString() === id.toString()
        )
      })
    } else {
      items = rawModelValueArray.map((id: any) => {
        return oldItems.find((item: any) => item[props.itemValue].toString() === id.toString())
      })
    }
  }

  selecting.value = many || items.length == 0
  internalSelectedItems.value = items
}

onMounted(() => {
  if (props.items && props.items.length) {
    internalItems.value = props.items
  }

  const value = modelValueArray(props.modelValue)
  if (value.length && internalItems.value.length) {
    internalSelectedItems.value = value.map((id: any) => {
      return internalItems.value.find(
        (item: any) => item[props.itemValue].toString() === id.toString()
      )
    })
  }

  selecting.value = many || value.length == 0

  setTimeout(() => {
    watch(modelValue, (newValue) => {
      if (!updating.v) {
        initValue(newValue)
      }
    })
  }, 100)
})

const setSearchTerm = (event: any) => {
  searchTerm.value = event
}

const many = props.many

// methods
const addItem = (event: any) => {
  if (!event) return

  const oldItems = internalSelectedItems.value
  autocompleteValue.value = []

  // check if previously selected
  if (
    internalSelectedItems.value.find(
      (element) => element[props.itemValue] == event[props.itemValue]
    )
  ) {
    return
  }

  const item = internalItems.value.find((element) => {
    return element[props.itemValue] == event[props.itemValue]
  })

  if (many) {
    internalSelectedItems.value.push(item)
  } else {
    internalSelectedItems.value = [item]
  }

  setValue(oldItems)
}

const changeOrder = (event: any) => {
  const old = startDragItems.value
  startDragItems.value = []
  setValue(old)
}

const removeItem = (id: string) => {
  const oldItems = internalSelectedItems.value
  internalSelectedItems.value = internalSelectedItems.value.filter(
    (element) => element[props.itemValue] != id
  )
  setValue(oldItems)
}

const setValue = (oldItems: any[]) => {
  const vals = internalSelectedItems.value.map((item) => item[props.itemValue]),
    val = many ? vals : vals.length ? vals[0] : '',
    old = many ? oldItems : oldItems.length ? oldItems[0] : ''

  selecting.value = many || internalSelectedItems.value.length == 0
  updating.v = true
  emit('update:modelValue', val)
}

let search: any,
  focused: any,
  doSearch: any,
  entered = () => {}

if (props.searchItemsFunc) {
  const searchDebounce = (func: Function, delay: number) => {
    let timeoutId: number

    return function (val: String) {
      clearTimeout(timeoutId)
      timeoutId = setTimeout(func, delay, val)
    }
  }

  interface ListResponse {
    Records: any[]
  }

  doSearch = (val: String) => {
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
    } else {
      internalItems.value = []
    }
  }

  entered = () => {
    internalItems.value = []
    doSearch(searchTerm.value)
  }
}
</script>
