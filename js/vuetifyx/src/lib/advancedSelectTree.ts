// The items of the dialog of vx-advanced-select as a tree (its itemParent):
// each under its parent when its parent is among them — named by its own text
// (itemNodeText) —, else at the root, named by its full text (itemText: a
// breadcrumb, its ancestors and it) when it has a parent not there.

export interface ItemsTreeNode {
  value: string
  title: string
  raw: any
  children?: ItemsTreeNode[]
}

export function itemsTree(
  items: any[],
  valueKey: string,
  parentKey: string,
  nodeTextKey: string,
  textKey: string
): ItemsTreeNode[] {
  const nodes = new Map<string, ItemsTreeNode>()
  for (const it of items || []) {
    const value = String(it[valueKey])
    nodes.set(value, { value, title: '', raw: it })
  }
  const roots: ItemsTreeNode[] = []
  for (const node of nodes.values()) {
    const parent = node.raw[parentKey] ? String(node.raw[parentKey]) : ''
    const own = node.raw[nodeTextKey || textKey] ?? node.raw[textKey]
    const holder = parent && parent !== node.value ? nodes.get(parent) : undefined
    if (holder) {
      node.title = String(own ?? '')
      ;(holder.children ||= []).push(node)
    } else {
      // at the root: its own name when it has no parent, else where it is
      node.title = String((parent ? node.raw[textKey] : own) ?? '')
      roots.push(node)
    }
  }
  return roots
}
