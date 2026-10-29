// The link cards of a text (x/markup: `[[/en-us/about]]` on a line of its
// own, a record of the admin `[[admin:posts:<id>:]]`): each link marked —
// .mdkit-link-card, its target data-href — made a card, at once and as more
// come (a preview, a comment drawn): its cover, its title, its summary, in
// the language of the page, asked of the server (the targets of a moment
// together); one the server does not answer (not there, not theirs to see)
// stays a link.
//
// Where to ask: the page's <meta name="mdkit-link-cards" content="<address>"
// data-locale="<locale>"> — the address answers a POST of target (many) and
// locale with the cards by target ({data: …}, {cards: …}, or the map) —, or
// a resolver of the targets given to watchLinkCards.

export const cardClass = 'mdkit-link-card'

// safeHref is href when a card may link it: a path of the site (not
// //host), or an address of this site; else ''.
export function safeHref(href, origin = (typeof location !== 'undefined' ? location.origin : '')) {
  const h = String(href || '').trim()
  if (/^\/(?![/\\])/.test(h)) return h
  if (origin && (h === origin || h.startsWith(origin + '/'))) return h
  return ''
}

// targetRe is what a card may link: a path of the site (not //host), or a
// record path of the admin (admin:…:, an action after) — x/markup's rule.
const targetRe = /^(\/|\/[^/\\\s"'<>][^\s"'<>]*|admin:[^\s"'\\]*:(![A-Za-z0-9_]+)?)$/

// cardTarget is the target of a card of text — typed or pasted: a path of
// the site, an address of this site (its path), a record path —; '', none.
export function cardTarget(text, origin = (typeof location !== 'undefined' ? location.origin : '')) {
  let t = String(text || '').trim()
  if (origin && (t === origin || t.startsWith(origin + '/'))) t = t.slice(origin.length) || '/'
  return targetRe.test(t) ? t : ''
}

// cardsOf is the map of the cards of what the server answered.
export function cardsOf(out) {
  if (!out || typeof out !== 'object') return {}
  return out.data || out.cards || out
}

// style is the look of a card, once for the page.
function style() {
  if (typeof document === 'undefined' || document.getElementById('mdkit-link-card-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-link-card-style'
  st.textContent =
    '.' + cardClass + '{margin:.75em 0}' +
    '.mdkit-card{display:flex;gap:.75rem;align-items:stretch;max-width:36rem;border:1px solid rgba(128,128,128,.35);border-radius:8px;' +
    'overflow:hidden;text-decoration:none!important;color:inherit!important;background:rgba(128,128,128,.04)}' +
    '.mdkit-card:hover{border-color:rgba(37,99,235,.6);background:rgba(37,99,235,.05)}' +
    '.mdkit-card-cover{flex:0 0 7rem;width:7rem;min-height:5rem;object-fit:cover;display:block;margin:0!important;border-radius:0!important}' +
    '.mdkit-card-body{padding:.6rem .75rem .6rem 0;min-width:0;display:flex;flex-direction:column;justify-content:center}' +
    '.mdkit-card-body:first-child{padding-left:.75rem}' +
    '.mdkit-card-title{font-weight:600;line-height:1.3}' +
    '.mdkit-card-summary{font-size:.875em;opacity:.75;margin-top:.2rem;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}' +
    '.mdkit-card-href{font-size:.75em;opacity:.55;margin-top:.2rem;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}'
  document.head.append(st)
}

// drawCard makes el the card of card: its texts as text (never HTML).
export function drawCard(el, card) {
  const href = safeHref(card && card.href)
  if (!href || !card.title) return false
  const doc = el.ownerDocument
  const a = doc.createElement('a')
  a.className = 'mdkit-card'
  a.href = href
  a.setAttribute('data-card-href', href)
  const cover = safeHref(card.cover) || (/^https:\/\//.test(card.cover || '') ? card.cover : '')
  if (cover) {
    const img = doc.createElement('img')
    img.className = 'mdkit-card-cover'
    img.src = cover
    img.alt = ''
    img.loading = 'lazy'
    a.append(img)
  }
  const body = doc.createElement('span')
  body.className = 'mdkit-card-body'
  const title = doc.createElement('span')
  title.className = 'mdkit-card-title'
  title.textContent = card.title
  body.append(title)
  if (card.summary) {
    const s = doc.createElement('span')
    s.className = 'mdkit-card-summary'
    s.textContent = card.summary
    body.append(s)
  }
  const h = doc.createElement('span')
  h.className = 'mdkit-card-href'
  h.textContent = href
  body.append(h)
  a.append(body)
  el.replaceChildren(a)
  el.setAttribute('data-card', 'drawn')
  return true
}

// metaResolver is the resolver of the page's meta (mdkit-link-cards); null,
// none.
export function metaResolver(doc = (typeof document !== 'undefined' ? document : null)) {
  const meta = doc && doc.querySelector('meta[name="mdkit-link-cards"]')
  const url = meta && meta.getAttribute('content')
  if (!url) return null
  const locale = meta.getAttribute('data-locale') || (doc.documentElement.lang || '')
  return async (targets) => {
    const body = new URLSearchParams()
    targets.forEach((t) => body.append('target', t))
    if (locale) body.append('locale', locale)
    const res = await fetch(url, { method: 'POST', credentials: 'same-origin', body })
    if (!res.ok) return {}
    return cardsOf(await res.json().catch(() => ({})))
  }
}

// known are the cards asked of a resolver, by target (null: none): a text
// drawn again (a preview) asks nothing again.
const known = new WeakMap()

// linkCards draws the cards under root not drawn yet, asking resolve for
// the targets it was not asked for yet (all of them at once).
export async function linkCards(root, resolve) {
  if (!root || !resolve) return
  style()
  const els = [...root.querySelectorAll('.' + cardClass + '[data-href]:not([data-card])')]
  if (!els.length) return
  els.forEach((el) => el.setAttribute('data-card', 'asked'))
  if (!known.has(resolve)) known.set(resolve, new Map())
  const cache = known.get(resolve)
  const targets = [...new Set(els.map((el) => el.getAttribute('data-href')))].filter((t) => !cache.has(t))
  if (targets.length) {
    let cards = {}
    try {
      cards = (await resolve(targets)) || {}
    } catch (e) {
      cards = {}
    }
    targets.forEach((t) => cache.set(t, cards[t] || null))
  }
  for (const el of els) {
    if (!drawCard(el, cache.get(el.getAttribute('data-href')))) el.setAttribute('data-card', 'none')
  }
}

// watchLinkCards draws the cards under root now and as they come (once per
// page: the admin's components and the site's may both start it); resolve,
// by default the page's meta's.
export function watchLinkCards(root = (typeof document !== 'undefined' ? document.body : null), resolve = metaResolver()) {
  if (!root || !resolve || typeof window === 'undefined' || window.__mdkitLinkCards) return
  window.__mdkitLinkCards = true
  let timer = null
  const run = () => {
    clearTimeout(timer)
    timer = setTimeout(() => linkCards(root, resolve), 60)
  }
  run()
  new MutationObserver((ms) => {
    if (ms.some((m) => [...m.addedNodes].some((n) => n.nodeType === 1 && (n.matches?.('.' + cardClass) || n.querySelector?.('.' + cardClass))))) run()
  }).observe(root, { childList: true, subtree: true })
}
