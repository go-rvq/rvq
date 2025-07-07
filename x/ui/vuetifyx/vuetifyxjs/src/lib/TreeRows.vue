<template>
  <template v-for="item in rows">
    <slot name="row" :item="item">
      <tr>
        <slot name="item-start" :item="item" />
        <slot name="item-controller" :item="item">
          <td>
            <div class="d-inline-flex">
              <slot
                v-if="getControl(item).path.length"
                name="depth"
                :item="item"
                :control="getControl(item)"
              >
                <div :style="`width:${getControl(item).path.length * 20}px`"></div>
              </slot>
              <slot
                v-if="item[childrenKey] && item[childrenKey].length"
                name="item-controller-buttom"
                :item="item"
                :control="getControl(item)"
                :tooggle="tooggleExpansion"
              >
                <v-btn
                  density="comfortable"
                  variant="text"
                  :icon="getControl(item).expanded ? 'mdi-chevron-up' : 'mdi-chevron-down'"
                  size="small"
                  @click="tooggleExpansion(item)"
                ></v-btn>
              </slot>
              <slot name="title" :item="item" :control="getControl(item)">{{
                item[titleKey]
              }}</slot>
            </div>
          </td>
        </slot>
        <slot name="item-container" :item="item" :control="getControl(item)">
          <slot name="pre-columns" :item="item" :control="getControl(item)" />
          <slot
            v-for="column in columns"
            :item="item"
            :control="getControl(item)"
            :column="column"
            :name="`column[${column}]`.toLowerCase()"
          >
            <td>
              <slot
                :name="`column-content[${column}]`.toLowerCase()"
                :item="item"
                :control="getControl(item)"
                :column="column"
              >
                {{ item[column] }}
              </slot>
            </td>
          </slot>
          <slot name="post-columns" :item="item" :control="getControl(item)" />
        </slot>
        <slot name="post-item" :item="item" :control="getControl(item)" />
      </tr>
    </slot>
  </template>
</template>
<script setup lang="ts">
import { isReactive, nextTick, onMounted, reactive, ref } from 'vue'

const props = defineProps({
  items: {
    type: Array<any>,
    default: null
  },
  columns: {
    type: Array<string>,
    default: []
  },
  titleKey: {
    type: String,
    default: 'title'
  },
  actionsKey: {
    type: String,
    default: 'actions'
  },
  childrenKey: {
    type: String,
    default: 'children'
  },
  expanded: {
    type: Boolean,
    default: false
  }
})

const rows = ref<any[]>([])

class Control {
  path: number[]
  expanded: boolean

  public constructor() {
    this.path = []
    this.expanded = props.expanded
  }
}

const controller = new Map<any, Control>(),
  getControl = (item: any): Control => controller.get(item) as Control

function initItem(item: any) {
  controller.set(item, new Control())
  item[props.childrenKey]?.forEach((child: any) => initItem(child))
}

let initialItems = reactive<any[]>([])

const init = () => {
  if (props.items) {
    if (isReactive(props.items)) {
      initialItems = props.items
    } else {
      initialItems = (props.items as any[]).map((item: any) =>
        isReactive(item) ? item : reactive(item)
      )
    }
  }

  console.log(initialItems)

  initialItems.forEach((item: any) => {
    initItem(item)
  })

  updateRows()
}

const updateRows = () => {
  const items: any[] = [],
    doItem = (item: any, path: any[]) => {
      const c = getControl(item)
      c.path = path

      if (c.expanded) {
        item[props.childrenKey]?.forEach((child: any, i: any) => {
          items[items.length] = child
          doItem(child, [...path, i])
        })
      }
    }

  initialItems.forEach((item, i) => {
    items.push(item)
    doItem(item, [i])
  })

  rows.value = items
}

const tooggleExpansion = (item: any) => {
  const c = getControl(item)
  c.expanded = !c.expanded
  updateRows()
}

onMounted(() => {
  nextTick(() => {
    init()
  })
})
</script>
