<template>
  <draggable
    :list="items"
    :group="group"
    :item-key="itemKey"
    :animation="150"
    :force-fallback="true"
    :fallback-on-body="true"
    :swap-threshold="0.65"
    :disabled="disabled"
    handle=".vx-tree-draggable__handle"
    ghost-class="vx-tree-draggable__ghost"
    tag="div"
    :class="['vx-tree-draggable__list', {'vx-tree-draggable__list--root': depth === 0}]"
    @start="onStart"
    @end="onEnd"
  >
    <template #item="{element}">
      <div class="vx-tree-draggable__node" :data-id="element[itemKey]">
        <div class="vx-tree-draggable__row" :style="{paddingInlineStart: depth * indent + 'px'}">
          <!-- the handle, at the start of the node: the only place that drags
          (a touch there drags, not scrolls — on a phone too) -->
          <span v-if="!disabled" class="vx-tree-draggable__handle">
            <v-icon size="small" icon="mdi-drag" />
          </span>
          <v-btn
            :icon="isOpen(element) ? 'mdi-chevron-down' : 'mdi-chevron-right'"
            :style="{visibility: childrenOf(element).length ? 'visible' : 'hidden'}"
            density="compact"
            variant="text"
            size="small"
            @click="toggle(element)"
          />
          <v-icon
            class="me-1"
            size="small"
            :icon="childrenOf(element).length ? (isOpen(element) ? 'mdi-folder-open-outline' : 'mdi-folder-outline') : 'mdi-file-outline'"
          />
          <div class="vx-tree-draggable__content">
            <slot name="item" :item="element" :depth="depth">{{ element[titleKey] }}</slot>
          </div>
        </div>
        <!-- always there — empty, it takes a node dropped as a child —, hidden
        when closed and not dragging -->
        <vx-tree-draggable
          v-show="isOpen(element) || dragging"
          :items="childrenOf(element)"
          :item-key="itemKey"
          :title-key="titleKey"
          :children-key="childrenKey"
          :group="group"
          :indent="indent"
          :disabled="disabled"
          :depth="depth + 1"
        >
          <template #item="s: {item: any; depth: number}">
            <slot name="item" :item="s.item" :depth="s.depth" />
          </template>
        </vx-tree-draggable>
      </div>
    </template>
  </draggable>
</template>

<script setup lang="ts">
// A tree whose nodes are dragged and dropped — in another place of their
// level, or into another node (a child), as the folders of a file system:
// items is changed in place (each node's children in childrenKey, made when
// missing), and the root tells "change" with it after each drop.
import draggable from 'vuedraggable'
import {inject, provide, reactive, ref, type Ref} from 'vue'

defineOptions({name: 'VxTreeDraggable'})

const props = withDefaults(
  defineProps<{
    items: any[]
    itemKey?: string
    titleKey?: string
    childrenKey?: string
    group?: string
    indent?: number
    disabled?: boolean
    // the nodes open at first: all of them (true), or their keys
    open?: boolean | any[]
    depth?: number
  }>(),
  {
    itemKey: 'id',
    titleKey: 'title',
    childrenKey: 'children',
    group: 'vx-tree-draggable',
    indent: 24,
    disabled: false,
    open: true,
    depth: 0
  }
)

const emit = defineEmits<{(e: 'change', items: any[]): void}>()

type Shared = {opened: Record<string, boolean>; dragging: Ref<boolean>; changed: () => void}

const shared: Shared =
  props.depth === 0
    ? {
        opened: reactive({} as Record<string, boolean>),
        dragging: ref(false),
        changed: () => emit('change', props.items)
      }
    : (inject('vxTreeDraggable') as Shared)
if (props.depth === 0) {
  provide('vxTreeDraggable', shared)
}
const dragging = shared.dragging

const childrenOf = (node: any) => {
  if (!Array.isArray(node[props.childrenKey])) {
    node[props.childrenKey] = []
  }
  return node[props.childrenKey]
}

const keyOf = (node: any) => String(node[props.itemKey])

const isOpen = (node: any) => {
  const k = keyOf(node)
  if (k in shared.opened) {
    return shared.opened[k]
  }
  return props.open === true || (Array.isArray(props.open) && props.open.map(String).includes(k))
}

const toggle = (node: any) => {
  shared.opened[keyOf(node)] = !isOpen(node)
}

const onStart = () => {
  dragging.value = true
}

const onEnd = () => {
  dragging.value = false
  shared.changed()
}
</script>

<style scoped>
.vx-tree-draggable__list {
  min-height: 6px;
}
.vx-tree-draggable__row {
  display: flex;
  align-items: center;
  gap: 2px;
  min-height: 36px;
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}
.vx-tree-draggable__row:hover {
  background: rgba(var(--v-theme-on-surface), 0.04);
}
.vx-tree-draggable__handle {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  min-height: 36px;
  cursor: grab;
  opacity: 0.6;
  touch-action: none;
  user-select: none;
  -webkit-user-select: none;
}
.vx-tree-draggable__content {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 4px;
}
.vx-tree-draggable__ghost > .vx-tree-draggable__row {
  background: rgba(var(--v-theme-primary), 0.12);
}
</style>
