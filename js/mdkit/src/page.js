// The script of a page that shows a text (dist/page.js, small): its link
// cards drawn (asked where the meta mdkit-link-cards says), its formulas and
// diagrams (KaTeX, Mermaid: loaded from where the meta mdkit-vendor says,
// only when there is one).
import { watchLinkCards } from './link-cards.js'
import { watchBlocks } from './blocks.js'

const start = () => {
  watchLinkCards()
  watchBlocks()
}
if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', start)
else start()
