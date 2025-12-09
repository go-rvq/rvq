import { ToolsItem } from '../../types'
import Item from './Item.vue'

export function CopyAndPasteContentItem() {
  return {
    component: {
      view: Item,
    }
  } as ToolsItem
}
