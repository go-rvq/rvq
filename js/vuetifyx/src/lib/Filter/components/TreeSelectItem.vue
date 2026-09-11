<script setup lang="ts">
import { computed } from 'vue'

// A filter item that presents its options as a checkbox tree. The options are
// flat paths ("PageOptions", "PageOptions.Layout", "Tags[0].Name"); this builds
// the tree from them, and the selected values are the chosen paths.

class Item {
  text: string = ''
  value: string = ''
}

class Model {
  options: Item[] = []
  valuesAre: string[] = []
  modifier: String = 'in'
}

const props = defineProps<{
  translations: any
}>()

const model = defineModel<Model>()

if (!model.value) {
  model.value = new Model()
}
if (!model.value.options) {
  model.value.options = []
}
if (!model.value.valuesAre) {
  model.value.valuesAre = []
}
if (!model.value.modifier) {
  model.value.modifier = 'in'
}

interface TreeNode {
  title: string
  value: string
  children?: TreeNode[]
}

// splitPath breaks "a.b[0].c" into ["a","b","b[0]","b[0].c"] cumulative segments,
// so each level of the tree is a valid path to the value at that level.
function segments(path: string): { title: string; full: string }[] {
  const out: { title: string; full: string }[] = []
  // tokens split on '.', but a '[n]' stays attached to its field.
  const parts = path.split('.')
  let acc = ''
  for (const p of parts) {
    // p may be "field" or "field[0]" or "field[0][1]"
    const m = p.match(/^([^[]*)((\[\d+\])*)$/)
    const name = m ? m[1] : p
    const idx = m ? m[2] : ''
    if (name) {
      acc = acc ? acc + '.' + name : name
      out.push({ title: name, full: acc })
    }
    if (idx) {
      // one node per index segment
      for (const im of idx.matchAll(/\[(\d+)\]/g)) {
        acc = acc + '[' + im[1] + ']'
        out.push({ title: '[' + im[1] + ']', full: acc })
      }
    }
  }
  return out
}

const tree = computed<TreeNode[]>(() => {
  const roots: TreeNode[] = []
  const byPath = new Map<string, TreeNode>()
  for (const opt of model.value!.options) {
    const segs = segments(opt.value)
    let parent: TreeNode[] = roots
    for (const s of segs) {
      let node = byPath.get(s.full)
      if (!node) {
        node = { title: s.title, value: s.full }
        byPath.set(s.full, node)
        parent.push(node)
      }
      if (!node.children) node.children = []
      parent = node.children
    }
  }
  // drop empty children arrays so leaves are leaves
  byPath.forEach((n) => {
    if (n.children && n.children.length === 0) delete n.children
  })
  return roots
})

const selected = computed<string[]>({
  get: () => model.value!.valuesAre || [],
  set: (v) => {
    model.value!.valuesAre = v
  }
})
</script>

<template>
  <div>
    <v-treeview
      :items="tree"
      item-value="value"
      item-title="title"
      select-strategy="independent"
      selectable
      open-all
      density="compact"
      v-model:selected="selected"
    />
  </div>
</template>
