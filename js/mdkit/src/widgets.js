// The buttons and the tabs of a text as the page shows them (blocks.js
// drawBlocks; the HTML: widgets-model.js): a button its icons and the color
// of the theme — of Vuetify's when the page has it —, its title, when it has
// none, its target's (the link cards'); the tabs a bar of their titles over
// the content of the one chosen — printed, each a section of its own (a
// card, open). The icons are Material Design's (the vendor's "mdi-css",
// loaded when one is there).
import { load, own } from './blocks.js'
import { metaResolver } from './link-cards.js'

export const buttonsSel = 'div[data-type="buttons"]'
export const tabsSel = 'div[data-type="tabs"]'

// the colors of the theme: Vuetify's, else these
const RGB = {
  primary: '25,118,210', secondary: '72,169,166', success: '76,175,80', info: '33,150,243',
  warning: '251,140,0', error: '176,0,32',
}

// widgetStyle is the look of the buttons and the tabs, once for the page.
export function widgetStyle() {
  if (typeof document === 'undefined' || document.getElementById('mdkit-widgets-style')) return
  const st = document.createElement('style')
  st.id = 'mdkit-widgets-style'
  let colors = ''
  for (const [c, rgb] of Object.entries(RGB)) {
    colors += `[data-color="${c}"]{--mdkit-c:var(--v-theme-${c},${rgb});--mdkit-on:var(--v-theme-on-${c},255,255,255)}`
  }
  st.textContent = colors +
    '[data-type="buttons"]{display:flex;flex-wrap:wrap;gap:.5rem;margin:.75em 0;justify-content:flex-start}' +
    '[data-type="buttons"][data-align="center"]{justify-content:center}' +
    '[data-type="buttons"][data-align="end"]{justify-content:flex-end}' +
    '[data-type="buttons"][data-align="stretch"]>a{flex:1 1 0}' +
    'a[data-type="button"]{display:inline-flex;align-items:center;justify-content:center;gap:.4em;min-height:2.25rem;padding:.35rem 1rem;' +
    'border-radius:6px;border:1px solid rgba(var(--mdkit-c,127,127,127),.6);color:rgb(var(--mdkit-c,60,60,60));background:transparent;' +
    'text-decoration:none!important;font-weight:500;line-height:1.2;transition:filter .15s,background .15s}' +
    'a[data-type="button"][data-color]{background:rgb(var(--mdkit-c));border-color:rgb(var(--mdkit-c));color:rgb(var(--mdkit-on))}' +
    'a[data-type="button"]:hover{filter:brightness(.92);background-color:rgba(var(--mdkit-c,127,127,127),.08)}' +
    'a[data-type="button"][data-color]:hover{background-color:rgb(var(--mdkit-c))}' +
    'a[data-type="button"]:empty::before{content:attr(href);opacity:.7}' +
    '[data-type="tabs"]{margin:.75em 0}' +
    '.mdkit-tabs-bar{display:flex;flex-wrap:wrap;gap:.15rem;border-bottom:1px solid rgba(127,127,127,.3);justify-content:flex-start}' +
    '[data-type="tabs"][data-align="center"]>.mdkit-tabs-bar{justify-content:center}' +
    '[data-type="tabs"][data-align="end"]>.mdkit-tabs-bar{justify-content:flex-end}' +
    '[data-type="tabs"][data-align="stretch"]>.mdkit-tabs-bar>button{flex:1 1 0}' +
    '.mdkit-tabs-bar>button{display:inline-flex;align-items:center;justify-content:center;gap:.4em;padding:.5rem 1rem;border:0;' +
    'border-bottom:2px solid transparent;margin-bottom:-1px;background:none;color:inherit;font:inherit;font-weight:500;cursor:pointer;opacity:.75}' +
    '.mdkit-tabs-bar>button[data-color]{color:rgb(var(--mdkit-c))}' +
    '.mdkit-tabs-bar>button[aria-selected="true"]{opacity:1;border-bottom-color:rgb(var(--mdkit-c,var(--v-theme-primary,25,118,210)))}' +
    '.mdkit-tabs-bar>button:hover{background:rgba(127,127,127,.08)}' +
    '[data-type="tab"]{padding:.75rem .25rem}' +
    '[data-type="tabs"][data-drawn] [data-type="tab"]:not(.mdkit-tab-shown){display:none}' +
    '.mdkit-tab-title{display:none}' +
    // printed: the browsers leave out the backgrounds — a button's text,
    // white on its color, would be white on the paper —: dark text, its
    // color in its border
    // printed, a section that folds open (beforeprint opens them: this, where
    // the browser has ::details-content, its content even so)
    '@media print{.mdkit-tabs-bar{display:none}details::details-content{display:block;content-visibility:visible}' +
    'a[data-type="button"],a[data-type="button"][data-color]{background:none!important;background-color:transparent!important;color:#111!important;' +
    'border:1.5px solid rgb(var(--mdkit-c,90,90,90))!important;filter:none!important}' +
    'a[data-type="button"] i.mdi{color:rgb(var(--mdkit-c,90,90,90))}' +
    '[data-type="tabs"][data-drawn] [data-type="tab"]:not(.mdkit-tab-shown){display:block}' +
    '[data-type="tab"]{display:block;border:1px solid rgba(127,127,127,.4);border-radius:6px;padding:0 .75rem .75rem;margin:.75rem 0;break-inside:avoid}' +
    '.mdkit-tab-title{display:flex;align-items:center;gap:.4em;font-weight:600;padding:.5rem 0;margin-bottom:.5rem;border-bottom:1px solid rgba(127,127,127,.3)}' +
    '.mdkit-tab-title[data-color]{color:rgb(var(--mdkit-c))}}'
  document.head.append(st)
}

// iconsOf puts in el the icons of item (data-icon-left, data-icon-right of
// from) around its text; true when it has one.
export function iconsOf(el, from = el) {
  const l = from.getAttribute('data-icon-left')
  const r = from.getAttribute('data-icon-right')
  const i = (name) => {
    const e = document.createElement('i')
    e.className = 'mdi ' + name
    e.setAttribute('aria-hidden', 'true')
    return e
  }
  if (l) el.prepend(i(l))
  if (r) el.append(i(r))
  return !!(l || r)
}

// icons loads the icons' styles when root has an icon.
function icons(root, assets) {
  if (root.querySelector('[data-icon-left],[data-icon-right]')) load('mdi-css', '', assets).catch(() => {})
}

// drawButtons draws the buttons under root not drawn yet: their icons, the
// titles they have not (their targets', else their addresses).
// (all: those in an editor too — its preview's)
export async function drawButtons(root, assets, resolve = metaResolver(), all = false) {
  const groups = [...root.querySelectorAll(buttonsSel)].filter((g) => !g.hasAttribute('data-drawn') && (all || own(g)))
  if (!groups.length) return
  widgetStyle()
  icons(root, assets)
  const untitled = []
  for (const g of groups) {
    g.setAttribute('data-drawn', '')
    for (const a of g.querySelectorAll('a[data-type="button"]')) {
      const empty = !(a.textContent || '').trim()
      if (empty) untitled.push(a)
      else wrap(a)
      iconsOf(a)
    }
  }
  if (!untitled.length) return
  const targets = [...new Set(untitled.map((a) => a.getAttribute('data-href')).filter(Boolean))]
  let cards = {}
  if (targets.length && resolve) cards = (await resolve(targets).catch(() => ({}))) || {}
  for (const a of untitled) {
    const card = cards[a.getAttribute('data-href')]
    const span = document.createElement('span')
    span.textContent = (card && card.title) || a.getAttribute('href')
    // between its icons
    const right = a.querySelector('i.mdi:last-child')
    if (right && a.getAttribute('data-icon-right')) a.insertBefore(span, right)
    else a.append(span)
  }
}

// wrap puts the text of a in a span (its icons around it).
function wrap(a) {
  const span = document.createElement('span')
  span.append(...a.childNodes)
  a.append(span)
}

let tabSeq = 0

// drawTabs draws the tabs under root not drawn yet: a bar of their titles
// (their icons, their colors), the first shown; printed, each a section of
// its own, its title over it (all: those in an editor too).
export function drawTabs(root, assets, all = false) {
  const sets = [...root.querySelectorAll(tabsSel)].filter((s) => !s.hasAttribute('data-drawn') && (all || own(s)))
  if (!sets.length) return
  widgetStyle()
  icons(root, assets)
  for (const set of sets) {
    set.setAttribute('data-drawn', '')
    const tabs = [...set.children].filter((c) => c.getAttribute('data-type') === 'tab')
    const bar = document.createElement('div')
    bar.className = 'mdkit-tabs-bar'
    bar.setAttribute('role', 'tablist')
    const n = ++tabSeq
    const buttons = tabs.map((tab, i) => {
      const b = document.createElement('button')
      b.type = 'button'
      b.id = 'mdkit-tab-' + n + '-' + i
      b.setAttribute('role', 'tab')
      b.setAttribute('aria-controls', 'mdkit-tabpanel-' + n + '-' + i)
      if (tab.getAttribute('data-color')) b.setAttribute('data-color', tab.getAttribute('data-color'))
      const span = document.createElement('span')
      span.textContent = tab.getAttribute('data-title') || String(i + 1)
      b.append(span)
      iconsOf(b, tab)
      tab.id = 'mdkit-tabpanel-' + n + '-' + i
      tab.setAttribute('role', 'tabpanel')
      tab.setAttribute('aria-labelledby', b.id)
      // its title, printed
      const title = b.cloneNode(true)
      title.removeAttribute('id')
      title.removeAttribute('role')
      title.removeAttribute('aria-controls')
      const t = document.createElement('div')
      t.className = 'mdkit-tab-title'
      if (tab.getAttribute('data-color')) t.setAttribute('data-color', tab.getAttribute('data-color'))
      t.append(...title.childNodes)
      tab.prepend(t)
      b.addEventListener('click', () => show(i))
      b.addEventListener('keydown', (e) => {
        const d = e.key === 'ArrowRight' ? 1 : e.key === 'ArrowLeft' ? -1 : 0
        if (!d) return
        e.preventDefault()
        const j = (i + d + tabs.length) % tabs.length
        show(j)
        buttons[j].focus()
      })
      return b
    })
    const show = (i) => tabs.forEach((tab, j) => {
      tab.classList.toggle('mdkit-tab-shown', i === j)
      buttons[j].setAttribute('aria-selected', String(i === j))
      buttons[j].tabIndex = i === j ? 0 : -1
    })
    bar.append(...buttons)
    set.prepend(bar)
    show(0)
  }
}

// printDetails opens, while the page prints, the sections that fold closed
// under root (not an editor's) — printed, each its card, as the tabs —, and
// closes them again after; once a page.
export function printDetails(root) {
  if (typeof window === 'undefined' || window.__mdkitPrintDetails) return
  window.__mdkitPrintDetails = true
  let opened = []
  window.addEventListener('beforeprint', () => {
    opened = [...root.querySelectorAll('details:not([open])')].filter((d) => own(d))
    opened.forEach((d) => { d.open = true })
  })
  window.addEventListener('afterprint', () => {
    opened.forEach((d) => { d.open = false })
    opened = []
  })
}
