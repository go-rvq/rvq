import { Editor } from '@tiptap/vue-3'

export interface ToolsItemComponent {
  props?: Object
  view: any
}

export interface ToolsItem {
  component ?: ToolsItemComponent
  title ?:()=>string
  action: (args:CommandArgs) => void
  divider?:boolean
  icon?:string
}

export interface Message {
  text: string
  color: string
}

export type Notify = (message:string, color?:string) => void

export interface CommandArgs {
  editor: Editor
  notify: Notify
}
