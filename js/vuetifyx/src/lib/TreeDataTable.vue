<template>
  <v-data-table-virtual
    :headers="currentHeaders"
    :items="rows"
    v-bind="tableProps"
    class="vx-tree-data-table"
    :class="{ 'vx-tree-data-table--guides': guides }"
  >
    <template v-slot:item.$control="{ item, value }">
      <!-- the row's own height (a title of lines): not forced to the cell's -->
      <div class="vx-tree-control d-flex align-center">
        <slot
          v-if="getControl(item).path.length"
          name="item.$control.prefix"
          :item="item"
          :control="getControl(item)"
        >
          <template v-if="guides">
            <span
              class="vx-tree-guide-space"
              :style="`width:${(getControl(item).lasts.length - guideFrom()) * guideStep}px`"
            ></span>
            <span v-for="(g, i) in guideLines(item)" :key="i" :class="g.class" :style="g.style"></span>
          </template>
          <div v-else :style="`width:${(getControl(item).path.length - 1) * 30}px`"></div>
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
        <div v-else-if="guides" class="vx-tree-guide-leaf"></div>
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
          <v-btn icon="mdi-cog" density="compact" flat v-bind="props" class="noprint"></v-btn>
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
import {isReactive, nextTick, onMounted, reactive, ref, watch} from 'vue'

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
  },
  // search: the rows whose searchKeys hold it (any case), with their
  // ancestors, opened; empty, the tree as it is
  search: {
    type: String,
    default: ''
  },
  searchKeys: {
    type: Array<string>,
    default: () => ['title']
  },
  // openAll: every row opened (true) or closed (false) as it changes; null,
  // each as it is
  openAll: {
    type: Boolean,
    default: null
  },
  // controlTitle: the header of the column of the tree (the titles)
  controlTitle: {
    type: String,
    default: ''
  },
  // settings: the column of the settings (which columns are shown)
  settings: {
    type: Boolean,
    default: true
  },
  // guides: who is whose child marked by lines, as the command tree marks
  // it (├── └── │), the whole height of the rows
  guides: {
    type: Boolean,
    default: false
  },
  // rootGuides: the rows at the root marked too, as children of a root of no
  // name (├── └── before them); off, the guides begin at their children
  rootGuides: {
    type: Boolean,
    default: false
  }
})

const rows = ref<any[]>([])

export interface Control {
  path: number[]
  // lasts: of each level of the path, the row there the last of its
  // siblings shown (the guides: └── and no │ under it)
  lasts: boolean[]
  // shownChildren: how many of its children are shown under it
  shownChildren: number
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
    lasts: [],
    shownChildren: 0,
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
    reactive({ title: props.controlTitle, key: '$control', sortable: false, noSelectable: true, visible: true })
  ]
  value.forEach((item: any) => {
    if (item.visible === undefined) {
      item.visible = true
    }
    h.push(item)
  })
  h.push(
    reactive({
      title: '',
      key: '$columns-control',
      visible: props.settings,
      sortable: false,
      width: '0%',
      noSelectable: true,
      headerProps: {
        class: 'noprint'
      },
      cellProps: {
        class: 'noprint'
      }
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

// matches says the item holds the search q in one of its searchKeys
const matches = (item: any, q: string): boolean =>
  (props.searchKeys as string[]).some((k) =>
    String(item?.[k] ?? '')
      .toLowerCase()
      .includes(q)
  )

const updateRows = () => {
  // a search shows the rows that hold it and their ancestors, opened
  const q = (props.search || '').trim().toLowerCase(),
    shown = (item: any): boolean =>
      !q || matches(item, q) || getControl(item).children.some((child: any) => shown(child))
  const items: any[] = [],
    doItem = (parent: any, item: any, path: any[], lasts: boolean[]) => {
      const c = getControl(item)
      c.path = path
      c.lasts = lasts
      c.parent = parent

      c.shownChildren = 0
      if (c.expanded || q) {
        const children = c.children.filter((child: any) => shown(child))
        c.shownChildren = children.length
        children.forEach((child: any, i: any) => {
          items[items.length] = child
          doItem(item, child, [...path, i], [...lasts, i === children.length - 1])
        })
      }
    }

  const roots = initialItems.filter((item: any) => shown(item))
  roots.forEach((item, i) => {
    items.push(item)
    doItem(null, item, [i], [i === roots.length - 1])
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

// guideStep is the width of a level of the guides; guideStart, where the
// first begins in the cell (its padding)
const guideStep = 20,
  guideStart = 16

// guideLines are the lines before a row, as the command tree draws them, the
// whole height of its cell: of each level above it, │ while that one has
// siblings after it; of its own, ├── (└── the last of its siblings) — a
// row at the root, only with rootGuides —; and, its children shown under it,
// │ from it down to them
// guideFrom is the first level the guides mark: the root's (rootGuides), else
// its children's
const guideFrom = () => (props.rootGuides ? 0 : 1)

const guideLines = (item: any): { class: string; style: string }[] => {
  const c = getControl(item),
    lasts = c?.lasts || [],
    from = guideFrom(),
    out: { class: string; style: string }[] = [],
    x = (level: number) => `left:${guideStart + (level - from) * guideStep + guideStep / 2}px`
  for (let i = from; i < lasts.length - 1; i++) {
    if (!lasts[i]) {
      out.push({ class: 'vx-tree-guide-v', style: x(i) })
    }
  }
  const own = lasts.length - 1
  if (own >= from) {
    out.push({ class: 'vx-tree-guide-v' + (lasts[own] ? ' vx-tree-guide-v--top' : ''), style: x(own) })
    // to its button, or — a leaf — almost to its title
    const leaf = !(props.expandable || c.children.length)
    out.push({
      class: 'vx-tree-guide-h',
      style: `${x(own)};width:${leaf ? guideStep / 2 + 24 : guideStep / 2}px`
    })
  }
  if (c?.shownChildren) {
    out.push({ class: 'vx-tree-guide-v vx-tree-guide-v--bottom', style: x(own + 1) })
  }
  return out
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

watch(
  () => props.search,
  () => updateRows()
)

// every row opened or closed (openAll), as it changes
watch(
  () => props.openAll,
  (open) => {
    if (open === null || open === undefined) {
      return
    }
    controllers.forEach((c) => {
      c.expanded = open
    })
    updateRows()
  }
)
</script>

<style>
.vx-tree-data-table--guides tbody td {
  position: relative;
}

/* a cell of lines — a long title, a value — apart from the row's borders */
.vx-tree-data-table tbody td {
  padding-top: 4px !important;
  padding-bottom: 4px !important;
}

.vx-tree-guide-space {
  flex-shrink: 0;
}

.vx-tree-guide-v,
.vx-tree-guide-h {
  position: absolute;
  pointer-events: none;
  border-color: rgba(var(--v-theme-on-surface), 0.3);
  border-style: solid;
  border-width: 0;
}

.vx-tree-guide-v {
  top: 0;
  bottom: 0;
  border-left-width: 1px;
}

.vx-tree-guide-v--top {
  bottom: 50%;
}

.vx-tree-guide-v--bottom {
  top: 50%;
}

.vx-tree-guide-h {
  top: 50%;
  border-top-width: 1px;
}

.vx-tree-data-table--guides .vx-tree-guide-leaf {
  width: 28px;
  flex-shrink: 0;
}
</style>
