// The plugin "details" of vx-tiptap-editor: a section that folds (tiptap's
// details: <details><summary>…</summary>…</details>, as GitHub's) — its button
// makes the blocks chosen one, or the one under the cursor its blocks again.
import {Extension} from '@tiptap/core'
import {Details, DetailsContent, DetailsSummary} from '@tiptap/extension-details'
import DetailsAction from '../components/DetailsAction.vue'
import type {TiptapPlugin} from './index'

export const DetailsButton = Extension.create({
  name: 'mdkitDetailsButton',
  addOptions() {
    return {
      divider: false,
      spacer: false,
      button: ({editor}: any) => ({component: DetailsAction, componentProps: {editor}})
    } as any
  }
})

const details: TiptapPlugin = {
  extensions: () => [Details.configure({persist: true, HTMLAttributes: {class: 'mdkit-details'}}), DetailsSummary, DetailsContent, DetailsButton]
}

export default details
