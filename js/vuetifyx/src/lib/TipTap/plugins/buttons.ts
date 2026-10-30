// The plugin "buttons" of vx-tiptap-editor: rvq mdkit's buttons — links as
// buttons, each its title (else its page's), description, color and icons
// —, made and changed in their dialog (ButtonsAction); buttons clicked, the
// same on them.
import {Extension} from '@tiptap/core'
import {Buttons} from '@rvq/mdkit'
import ButtonsAction from '../components/ButtonsAction.vue'
import type {TiptapPlugin} from './index'

export const ButtonsButton = Extension.create({
  name: 'mdkitButtonsButton',
  addOptions() {
    return {
      divider: false,
      spacer: false,
      button: ({editor}: any) => ({component: ButtonsAction, componentProps: {editor}})
    } as any
  }
})

const buttons: TiptapPlugin = {
  extensions: () => [Buttons, ButtonsButton]
}

export default buttons
