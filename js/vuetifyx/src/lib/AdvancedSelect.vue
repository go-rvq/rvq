<template>
  <div class="vx-advanced-select mb-2" :class="'vx-advanced-select__' + (many ? 'many' : 'one')">
    <!-- the label on a line of its own: the chosen, the button below it -->
    <label
      v-if="label"
      class="v-label theme--light d-block mb-1"
      v-html="label"
      :style="errorMessages.length ? 'color: rgb(var(--v-theme-error));opacity: inherit;' : ''"
    ></label>
    <p v-if="hint" v-html="hint"></p>
    <v-card
      :color="errorMessages.length ? 'error' : ''"
      v-if="internalSelectedItems.length > 0"
      variant="flat"
      class="mb-2 vx-advanced-select"
    >
      <template v-if="chips">
        <v-chip-group column>
          <v-chip v-for="item in internalSelectedItems">
            {{ item[itemText] }}
            <template v-slot:append>
              <v-btn
                v-if="dialog && !many"
                density="compact"
                variant="text"
                icon="mdi-swap-horizontal"
                :title="changeText"
                @click="openDialog"
              ></v-btn>
              <v-btn
                density="compact"
                variant="text"
                icon="mdi-delete"
                @click="removeItem(item[itemValue])"
              ></v-btn>
            </template>
          </v-chip>
        </v-chip-group>
      </template>

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
              <template v-if="dialog && !many" v-slot:append>
                <!-- the change of the one chosen: on its right -->
                <v-btn
                  variant="text"
                  icon="mdi-swap-horizontal"
                  density="compact"
                  :title="changeText"
                  class="vx-advanced-select__change"
                  @click="openDialog"
                ></v-btn>
              </template>
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

    <template v-if="dialog">
      <!-- the selector in a dialog: the search, the list, the pages -->
      <v-chip
        v-if="many || !internalSelectedItems.length"
        class="vx-advanced-select__open mb-2"
        prepend-icon="mdi-magnify"
        variant="outlined"
        size="small"
        :color="errorMessages.length ? 'error' : undefined"
        @click="openDialog"
        >{{ selectText }}</v-chip
      >
      <div v-if="errorMessages.length" class="v-messages" style="opacity: 1; color: rgb(var(--v-theme-error))">
        <div role="alert" v-for="(item, i) in errorMessages" class="v-messages__message">{{ item }}</div>
      </div>
      <vx-dialog
        v-model="dialogOpen"
        :title="dialogTitle || label"
        :container-props="{ maxWidth: 600 }"
        density="compact"
        expandable
        closable
        :class="['vx-advanced-select__dialog', { 'vx-advanced-select__dialog--top': dialogOnTop }]"
      >
        <template v-slot:body>
          <v-text-field
            v-model="dialogKeyword"
            prepend-inner-icon="mdi-magnify"
            :label="searchText"
            variant="outlined"
            density="compact"
            hide-details
            clearable
            autofocus
            autocomplete="off"
            class="mb-2"
            @update:model-value="dialogSearch"
          />
          <v-progress-linear v-if="isLoading" indeterminate />
          <!-- a tree (itemParent): each under its parent, when it is there -->
          <v-treeview
            v-if="itemParent && dialogItems.length"
            :items="dialogTree"
            item-title="title"
            item-value="value"
            item-children="children"
            open-all
            activatable
            active-strategy="single-independent"
            :activated="selectedValues"
            density="compact"
            color="primary"
            class="vx-advanced-select__tree"
            @update:activated="pickValue"
          />
          <v-list v-else density="compact">
            <v-list-item
              v-for="item in dialogItems"
              :key="item[itemValue]"
              :title="item[itemText]"
              :prepend-avatar="item[itemImage]"
              :active="isSelected(item)"
              color="primary"
              @click="pick(item)"
            >
              <template v-if="many" v-slot:prepend>
                <v-checkbox-btn :model-value="isSelected(item)" density="compact" />
              </template>
            </v-list-item>
            <v-list-item v-if="!isLoading && !dialogItems.length" :title="noItemsText" disabled />
          </v-list>
          <v-pagination
            v-if="dialogPages > 1"
            v-model="dialogPage"
            :length="dialogPages"
            :total-visible="7"
            density="compact"
            @update:model-value="loadDialog"
          />
        </template>
      </vx-dialog>
    </template>

    <v-autocomplete
      v-else-if="selecting"
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
      autocomplete="off"
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

    <div v-if="!hideDetails" class="v-input__details" role="alert" aria-live="polite">
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
  </div>
</template>
<script setup lang="ts">
import draggable from 'vuedraggable'
import VxDialog from '@/lib/Dialog.vue'

import {computed, isReactive, onMounted, ref, toRaw, watch} from 'vue'
import {itemsTree} from './advancedSelectTree'

const props = defineProps({
  hint: {
    type: String,
    default: ''
  },
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
  // itemParent: the key of an item's parent (its value) — the dialog a tree
  // (advancedSelectTree.ts) —; itemNodeText, the text of an item under its
  // parent (none: itemText)
  itemParent: {
    type: String,
    default: ''
  },
  itemNodeText: {
    type: String,
    default: ''
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
  hideDetails: {
    type: Boolean,
    default: false
  },
  autoComplete: {
    type: Object,
    default: {}
  },
  // dialog: the selector opens in a dialog — the search, the list, the pages
  // — in place of the autocomplete; the items, the value and the cache of the
  // chosen ones are the same
  dialog: {
    type: Boolean,
    default: false
  },
  // searchPageFunc(keyword, page) is a page of the search: {Records,
  // TotalCount, PerPage}; none, searchItemsFunc(keyword), a page alone
  searchPageFunc: {
    type: Function,
    default: null
  },
  dialogTitle: { type: String, default: '' },
  // dialogOnTop: the dialog above the overlay it opens from (the menu of a
  // filter, whose z-index is high)
  dialogOnTop: { type: Boolean, default: false },
  selectText: { type: String, default: 'Select' },
  changeText: { type: String, default: 'Change' },
  searchText: { type: String, default: 'Search' },
  noItemsText: { type: String, default: 'No items' },
  closeText: { type: String, default: 'Close' }
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

    // a value with no item (not given, or of another key) is not drawn:
    // nothing of it to show
    if (internalItems.value.length) {
      items = rawModelValueArray
        .map((id: any) => {
          return internalItems.value.find(
            (item: any) => item[props.itemValue].toString() === id.toString()
          )
        })
        .filter(Boolean)
    } else {
      items = rawModelValueArray
        .map((id: any) => {
          return oldItems.find((item: any) => item[props.itemValue].toString() === id.toString())
        })
        .filter(Boolean)
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
    internalSelectedItems.value = value
      .map((id: any) => {
        return internalItems.value.find(
          (item: any) => item[props.itemValue].toString() === id.toString()
        )
      })
      .filter(Boolean)
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
  // the records chosen, whole: their labels for who shows them (a filter)
  emit('update:selectedItems', internalSelectedItems.value)
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

// the dialog of the selector
const dialogOpen = ref(false)
const dialogKeyword = ref<string>('')
const dialogItems = ref<any[]>([])
const dialogPage = ref(1)
const dialogPages = ref(0)

const isSelected = (item: any) =>
  internalSelectedItems.value.some(
    (s: any) => s && item && String(s[props.itemValue]) === String(item[props.itemValue])
  )

const loadDialog = () => {
  const keyword = dialogKeyword.value || ''
  let p: Promise<any>
  if (props.searchPageFunc) {
    p = props.searchPageFunc(keyword, dialogPage.value)
  } else if (props.searchItemsFunc) {
    p = props.searchItemsFunc(keyword).then((records: any[]) => ({ Records: records }))
  } else {
    // no search: the items given, filtered here
    const k = keyword.toLowerCase()
    p = Promise.resolve({
      Records: (props.items as any[]).filter(
        (i: any) => !k || String(i[props.itemText] ?? '').toLowerCase().includes(k)
      )
    })
  }
  isLoading.value = true
  p.then((r: any) => {
    const records = r?.Records || []
    dialogItems.value = records
    // the items known: the labels of the chosen ones
    internalItems.value = records
    const perPage = r?.PerPage || records.length
    dialogPages.value = perPage && r?.TotalCount ? Math.ceil(r.TotalCount / perPage) : 0
  }).finally(() => {
    isLoading.value = false
  })
}

let dialogSearchTimeout: any
const dialogSearch = () => {
  clearTimeout(dialogSearchTimeout)
  dialogSearchTimeout = setTimeout(() => {
    dialogPage.value = 1
    loadDialog()
  }, 400)
}

const openDialog = () => {
  dialogOpen.value = true
  dialogPage.value = 1
  loadDialog()
}

// pick chooses item: one, the dialog closes; many, it is marked or unmarked
// the dialog's items as a tree (itemParent), and the values chosen
const dialogTree = computed(() =>
  itemsTree(dialogItems.value, props.itemValue, props.itemParent, props.itemNodeText, props.itemText)
)
const selectedValues = computed(() => internalSelectedItems.value.map((it: any) => String(it[props.itemValue])))
// pickValue picks the item of the value activated in the tree
const pickValue = (values: any[]) => {
  const v = values && values.length ? String(values[values.length - 1]) : ''
  const item = dialogItems.value.find((it: any) => String(it[props.itemValue]) === v)
  if (item && !isSelected(item)) pick(item)
}

const pick = (item: any) => {
  const oldItems = internalSelectedItems.value
  if (isSelected(item)) {
    if (!many) {
      dialogOpen.value = false
      return
    }
    internalSelectedItems.value = oldItems.filter(
      (s: any) => String(s[props.itemValue]) !== String(item[props.itemValue])
    )
  } else if (many) {
    internalSelectedItems.value = [...oldItems, item]
  } else {
    internalSelectedItems.value = [item]
  }
  setValue(oldItems)
  if (!many) {
    dialogOpen.value = false
  }
}

defineExpose({ openDialog, pick })
</script>

<style>
.vx-advanced-select__dialog--top {
  z-index: 100010 !important;
}
</style>
