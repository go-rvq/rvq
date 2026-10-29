// The script of the link cards of a page (dist/link-cards.js, small: the
// pages of a site that show a text with cards): the cards drawn, asked where
// the page's meta says (mdkit-link-cards).
import { watchLinkCards } from './link-cards.js'

const start = () => watchLinkCards()
if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', start)
else start()
