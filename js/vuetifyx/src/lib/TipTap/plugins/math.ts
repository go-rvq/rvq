// The plugin "math" of vx-tiptap-editor: formulas (LaTeX, drawn by KaTeX —
// tiptap's mathematics: span[data-type=inline-math], div[data-type=block-math],
// the HTML x/markup makes of $…$, $$…$$) — its button opens the formula's
// dialog (MathAction: the LaTeX, in the line or on its own, its preview); a
// formula clicked, the same on it.
import {Extension} from '@tiptap/core'
import Mathematics from '@tiptap/extension-mathematics'
import MathAction from '../components/MathAction.vue'
import type {TiptapPlugin} from './index'

// MathButton is the button of the dialog, and the clicks on a formula that
// open it (editor.storage.mdkitMath.open: the button's)
export const MathButton = Extension.create({
  name: 'mdkitMath',
  addStorage() {
    return {open: null as null | ((latex: string, pos: number | null, block: boolean) => void)}
  },
  addOptions() {
    return {
      divider: false,
      spacer: false,
      button: ({editor}: any) => ({component: MathAction, componentProps: {editor}})
    } as any
  }
})

const math: TiptapPlugin = {
  extensions: () => {
    let storage: any = null
    const open = (block: boolean) => (node: any, pos: number) => storage?.open?.(node.attrs.latex, pos, block)
    const button = MathButton.extend({
      onCreate() {
        storage = this.storage
      }
    })
    return [
      Mathematics.configure({
        inlineOptions: {onClick: open(false)},
        blockOptions: {onClick: open(true)},
        katexOptions: {throwOnError: false}
      }),
      button
    ]
  }
}

export default math
