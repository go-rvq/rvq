// The plugin "cards" of vx-tiptap-editor: rvq mdkit's link cards — a page of
// the site or a record of the admin, on a line of its own, shown as its card
// (its cover, title and summary), a bin to take it out —: its button opens
// the picker (the pages and records found as one types, or an address
// typed); an address of the site pasted alone is a card (Ctrl+Z, the text).
import {Extension} from '@tiptap/core'
import {LinkCard} from '@rvq/mdkit'
import LinkCardAction from '../components/LinkCardAction.vue'
import {mdkitWords} from '@/lib/locale/mdkit'
import type {TiptapPlugin} from './index'

export const LinkCardButton = Extension.create({
  name: 'mdkitLinkCards',
  addOptions() {
    return {
      divider: false,
      spacer: false,
      button: ({editor}: any) => ({component: LinkCardAction, componentProps: {editor}})
    } as any
  }
})

const cards: TiptapPlugin = {
  extensions: ({t}) => [LinkCard.configure({removeTitle: mdkitWords(t).card.remove}), LinkCardButton]
}

export default cards
