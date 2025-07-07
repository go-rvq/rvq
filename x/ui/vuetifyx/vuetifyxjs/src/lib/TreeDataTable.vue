<template>
  <v-data-table-virtual :headers="currentHeaders" :items="rows" v-bind="tableProps">
    <template v-slot:item.$control="{ item, value }">
      <div class="d-inline-flex" style="height: 100%; vertical-align: middle">
        <slot
          v-if="getControl(item).path.length"
          name="item.$control.prefix"
          :item="item"
          :control="getControl(item)"
        >
          <div :style="`width:${(getControl(item).path.length - 1) * 30}px`"></div>
        </slot>
        <slot
          v-if="expandable || getControl(item).children.length"
          name="item.$control.buttom"
          :item="item"
          :control="getControl(item)"
          :tooggle="tooggleExpansion"
        >
          <div style="display: flex; flex-direction: row; align-items: center">
            <v-btn
              density="compact"
              variant="text"
              :icon="getControl(item).expanded ? 'mdi-menu-down' : 'mdi-menu-right'"
              size="small"
              @click="tooggleExpansion(item)"
            ></v-btn>
          </div>
        </slot>
        <div style="display: flex; flex-direction: row; align-items: center">
          <slot name="item.$control.title" :item="item" :control="getControl(item)">
            {{ props.title(item) }}
          </slot>
        </div>
      </div>
    </template>
    <template v-for="(_, name, index) in $slots as {}" #[name]="scope" :key="index">
      <slot :name v-bind="scope" />
    </template>
    <template v-slot:header.$columns-control="{ column }">
      <v-dialog max-width="450" scrollable v-bind="settingsProps">
        <template v-slot:activator="{ props }">
          <v-btn icon="mdi-cog" density="compact" flat v-bind="props"></v-btn>
        </template>
        <template v-slot:default="{ isActive }">
          <v-card density="compact">
            <v-toolbar density="compact">
              <v-toolbar-title>{{ settingsTitle }}</v-toolbar-title>
              <v-spacer />
              <v-btn
                density="compact"
                class="ml-auto"
                icon="mdi-close"
                flat
                @click="isActive.value = false"
              ></v-btn>
            </v-toolbar>
            <v-divider />
            <v-card-text>
              <slot name="prepend-settings"></slot>
              <v-list density="compact" select-strategy="leaf">
                <v-list-subheader>{{ settingsColumnsTitle }}</v-list-subheader>
                <v-list-item
                  v-for="item in selectableHeaders"
                  density="compact"
                  :key="item.key"
                  :title="item.title"
                  :value="item.key"
                >
                  <template #prepend>
                    <v-list-item-action start>
                      <v-checkbox-btn v-model="item.visible"></v-checkbox-btn>
                    </v-list-item-action>
                  </template>
                </v-list-item>
              </v-list>
              <slot name="append-settings"></slot>
            </v-card-text>
            <v-divider />
            <v-card-actions>
              <v-btn
                class="ml-auto"
                prepend-icon="mdi-check"
                color="primary"
                variant="flat"
                @click="
                  () => {
                    applySettings()
                    isActive.value = false
                  }
                "
                >{{ settingsOkText }}
              </v-btn>
            </v-card-actions>
          </v-card>
        </template>
      </v-dialog>
    </template>
  </v-data-table-virtual>
</template>
<script setup lang="ts">
import { isReactive, nextTick, onMounted, reactive, ref } from 'vue'

const props = defineProps({
  settingsTitle: {
    type: String,
    default: 'Settings'
  },
  settingsColumnsTitle: {
    type: String,
    default: 'Colums'
  },
  settingsOkText: {
    type: String,
    default: 'OK'
  },
  tableProps: {
    type: Object,
    default: {}
  },
  settingsProps: {
    type: Object,
    default: {}
  },
  headers: {
    type: Array<any>,
    default: null
  },
  items: {
    type: Array<any>,
    default: null
  },
  title: {
    type: Function,
    default: (item: any) => item.title
  },
  children: {
    type: Function,
    default: (item: any) => item.children
  },
  expanded: {
    type: Boolean,
    default: false
  },
  expandable: {
    type: Boolean,
    default: false
  }
})

const rows = ref<any[]>([])

export interface Control {
  path: number[]
  expanded: boolean
  item: any
  children: any[]
  parent: any
}

const controllers = new Map<any, Control>(),
  getControl = (item: any): Control => {
    return controllers.get(item) as Control
  }

function initItem(item: any) {
  const c: Control = {
    path: [],
    expanded: props.expanded,
    item: item,
    children: reactive(((props.children(item) || []) as any[]).map(reactive)),
    parent: null
  }

  controllers.set(item, c)
  c.children.forEach((child: any) => initItem(child))
}

let initialItems = reactive<any[]>([])

const allHeaders = ref<any[]>([{}]),
  currentHeaders = ref<any[]>([{}]),
  selectableHeaders = ref<any[]>([{}])

const updateHeaders = (value: any[]) => {
  const h: any[] = [
    reactive({ title: '', key: '$control', sortable: false, noSelectable: true, visible: true })
  ]
  value.forEach((item: any) => {
    console.log(item, item.visible)
    if (item.visible === undefined) {
      item.visible = true
    }
    h.push(item)
  })
  h.push(
    reactive({
      title: '',
      key: '$columns-control',
      sortable: false,
      width: '0%',
      noSelectable: true,
      visible: true
    })
  )

  allHeaders.value = h
  currentHeaders.value = h.filter((v) => v.visible)
  selectableHeaders.value = h.filter((v) => !v.noSelectable)
}

const selectionHeadersChange = () => {
  currentHeaders.value = allHeaders.value.filter((v: any) => v.visible)
}

const init = () => {
  if (props.items) {
    if (isReactive(props.items)) {
      initialItems = props.items
    } else {
      initialItems = reactive((props.items as any[]).map(reactive))
    }
  }

  initialItems.forEach((item: any) => {
    initItem(item)
  })

  updateRows()

  if (props.headers) {
    let headers: any[] = []

    if (Array.isArray(props.headers)) {
      headers = props.headers as any[]
    } else {
      headers = (props.headers as any).value as any[]
    }

    updateHeaders(headers)
  }
}

const updateRows = () => {
  const items: any[] = [],
    doItem = (parent: any, item: any, path: any[]) => {
      const c = getControl(item)
      c.path = path
      c.parent = parent

      if (c.expanded) {
        c.children.forEach((child: any, i: any) => {
          items[items.length] = child
          doItem(item, child, [...path, i])
        })
      }
    }

  initialItems.forEach((item, i) => {
    items.push(item)
    doItem(null, item, [i])
  })

  rows.value = items

  emit('postLoad', {
    items: items,
    controllers: controllers,
    toggle: tooggleExpansion,
    count: count,
    walk: walk,
    update: updateRows,
    setChildren: setChildren,
    remove: remove
  })
}

const internalWalk = (items: any[], handler: (data: ItemControl) => number) => {
  for (let i = 0; i < items.length; i++) {
    const item = items[i]
    const control = controllers.get(item) as Control
    const ret = handler({ item: item, control: control })
    if (!ret) {
      if (internalWalk(control?.children, handler) == 2) {
        // skip tree
        return 2
      }
    } else if (ret == 1) {
      // skip siblings
      break
    } else {
      // skip tree
      return 2
    }
  }
}

const walk = (handler: (c: ItemControl) => number) => {
  internalWalk(rows.value, handler)
}

const setChildren = (item: any, children: any[]) => {
  const c = getControl(item)
  internalWalk(c.children, (c: ItemControl) => {
    controllers.delete(c.item)
    return 0
  })
  children.forEach(initItem)
  c.children = children
  updateRows()
}

const remove = (item: any) => {
  const c = getControl(item)
  if (c.parent) {
    const p = getControl(c.parent)
    p.children = p.children.filter((child: any) => child !== item)
  } else {
    initialItems = initialItems.filter((child: any) => child !== item)
  }

  updateRows()
}

export interface ItemControl {
  item: any
  control: Control
}

export interface PostLoadEvent {
  items: any[]
  controllers: Map<any, Control>
  toggle: (item: any) => void
  count: (accept?: (item: any) => boolean) => number
  walk: (handler: (data: ItemControl) => number) => void
  update: () => void
  remove: (item: any) => void
  setChildren: (item: any, items: any[]) => void
}

export interface ExpandEvent {
  c: Control
  update: () => void
}

const count = (accept?: (item: any) => boolean): number => {
  let val = 0
  const doAccept: (item: any) => boolean = accept || ((item: any) => true),
    doItem = (item: any, path: any[]) => {
      if (doAccept(item)) {
        val++
        getControl(item).children.forEach((child: any, i: any) => {
          doItem(child, [...path, i])
        })
      }
    }

  initialItems.forEach((item, i) => {
    doItem(item, [i])
  })

  return val
}

const tooggleExpansion = (item: any) => {
  const c = getControl(item)
  c.expanded = !c.expanded

  if (c.expanded) {
    emit('expand', {
      c,
      update: updateRows
    })
  }

  updateRows()
}

const applySettings = () => {
  selectionHeadersChange()
  emit('applySettings')
}

const emit = defineEmits<{
  (e: 'postLoad', props: PostLoadEvent): void
  (e: 'expand', props: ExpandEvent): void
  (e: 'applySettings'): void
}>()

onMounted(() => {
  nextTick(() => {
    init()
  })
})
</script>
