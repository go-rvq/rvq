// The plugin "diagrams" of vx-tiptap-editor: rvq mdkit's diagram (Mermaid,
// drawn in the editor; its HTML a ```mermaid fence's) — its button opens the
// diagram's dialog (DiagramAction: its source, its preview); a diagram
// clicked, the same on it.
import {Extension} from '@tiptap/core'
import {Diagram} from '@rvq/mdkit'
import DiagramAction from '../components/DiagramAction.vue'
import type {TiptapPlugin} from './index'

export const DiagramButton = Extension.create({
  name: 'mdkitDiagramButton',
  addOptions() {
    return {
      divider: false,
      spacer: false,
      button: ({editor}: any) => ({component: DiagramAction, componentProps: {editor}})
    } as any
  }
})

const diagrams: TiptapPlugin = {
  extensions: () => [Diagram, DiagramButton]
}

export default diagrams
