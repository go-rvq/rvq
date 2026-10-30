// The blocks of a text the page draws (x/markup): its formulas —
// span[data-type=inline-math], div[data-type=block-math], their LaTeX in
// data-latex — by KaTeX, its diagrams — pre > code.language-mermaid — by
// Mermaid, its code — pre > code.language-* — by Prism (code.js), at once
// and as more come (a preview, a comment drawn); KaTeX and
// Mermaid loaded only when there is one to draw, from where the page's meta
// says — <meta name="mdkit-vendor" content='{"katex": {"src", "local",
// "integrity"}, "katex-css": …, "mermaid": …}'>: the CDN, else the app's copy
// (rvq js/api/vendor) —. What is in an editor (contenteditable) is its own.

import { drawCode, codeSel } from './code.js'

const mathSel = '[data-type="inline-math"][data-latex],[data-type="block-math"][data-latex]'
const diagramSel = 'pre > code.language-mermaid'

// vendorOf is the assets of the page's meta (mdkit-vendor); {} none.
export function vendorOf(doc = (typeof document !== 'undefined' ? document : null)) {
  const meta = doc && doc.querySelector('meta[name="mdkit-vendor"]')
  try {
    return (meta && JSON.parse(meta.getAttribute('content') || '{}')) || {}
  } catch (e) {
    return {}
  }
}

const loading = {}

// load is the asset name loaded once — a script (its global then there) or a
// stylesheet —: its CDN, else its copy.
export function load(name, global, assets = vendorOf()) {
  if (global && typeof window !== 'undefined' && window[global]) return Promise.resolve(window[global])
  if (loading[name]) return loading[name]
  const a = assets[name]
  if (!a) return Promise.reject(new Error('mdkit: no asset ' + name))
  loading[name] = new Promise((resolve, reject) => {
    const style = !global
    const add = (src, integrity, onerror) => {
      const el = document.createElement(style ? 'link' : 'script')
      if (style) {
        el.rel = 'stylesheet'
        el.href = src
      } else {
        el.src = src
      }
      if (integrity) {
        el.integrity = integrity
        el.crossOrigin = 'anonymous'
      }
      el.onload = () => resolve(global ? window[global] : true)
      el.onerror = onerror
      document.head.append(el)
    }
    add(a.src, a.integrity, () => (a.local ? add(a.local, '', () => reject(new Error('mdkit: ' + name + ' not loaded'))) : reject(new Error('mdkit: ' + name))))
  })
  return loading[name]
}

// own says el is what a page draws: not in an editor (one read only too:
// tiptap's, ProseMirror)
export const own = (el) => !el.closest('[contenteditable="true"],.ProseMirror')

// drawMath draws the formulas under root not drawn yet.
export async function drawMath(root, assets = vendorOf()) {
  const els = [...root.querySelectorAll(mathSel)].filter((el) => !el.hasAttribute('data-drawn') && own(el))
  if (!els.length) return
  els.forEach((el) => el.setAttribute('data-drawn', 'asked'))
  let katex
  try {
    load('katex-css', '', assets).catch(() => {})
    katex = await load('katex', 'katex', assets)
  } catch (e) {
    els.forEach((el) => el.setAttribute('data-drawn', 'none'))
    return
  }
  for (const el of els) {
    try {
      katex.render(el.getAttribute('data-latex') || '', el, {
        displayMode: el.getAttribute('data-type') === 'block-math', throwOnError: false,
      })
      el.setAttribute('data-drawn', 'drawn')
    } catch (e) {
      el.setAttribute('data-drawn', 'error')
    }
  }
}

let diagrams = 0

// drawDiagrams draws the diagrams under root not drawn yet: each code made
// its picture (an svg in div.mdkit-diagram, its source kept); one that does
// not draw stays code.
export async function drawDiagrams(root, assets = vendorOf()) {
  const codes = [...root.querySelectorAll(diagramSel)].filter((c) => !c.hasAttribute('data-drawn') && own(c))
  if (!codes.length) return
  codes.forEach((c) => c.setAttribute('data-drawn', 'asked'))
  let mermaid
  try {
    mermaid = await load('mermaid', 'mermaid', assets)
    if (!mermaid.__mdkit) {
      const dark = typeof matchMedia !== 'undefined' && matchMedia('(prefers-color-scheme: dark)').matches
      mermaid.initialize({ startOnLoad: false, securityLevel: 'strict', theme: dark ? 'dark' : 'default' })
      mermaid.__mdkit = true
    }
  } catch (e) {
    codes.forEach((c) => c.setAttribute('data-drawn', 'none'))
    return
  }
  for (const code of codes) {
    const source = code.textContent || ''
    try {
      const { svg } = await mermaid.render('mdkit-diagram-' + (++diagrams), source)
      const div = document.createElement('div')
      div.className = 'mdkit-diagram'
      div.setAttribute('data-source', source)
      div.innerHTML = svg
      code.parentElement.replaceWith(div)
    } catch (e) {
      code.setAttribute('data-drawn', 'error')
      code.parentElement.classList.add('mdkit-diagram-error')
    }
  }
}

// drawBlocks draws the formulas, the diagrams and the code under root.
export function drawBlocks(root, assets = vendorOf()) {
  return Promise.all([drawMath(root, assets), drawDiagrams(root, assets), drawCode(root, assets)])
}

// watchBlocks draws the blocks under root now and as they come (once per
// page).
export function watchBlocks(root = (typeof document !== 'undefined' ? document.body : null), assets = vendorOf()) {
  if (!root || typeof window === 'undefined' || window.__mdkitBlocks || !Object.keys(assets).length) return
  window.__mdkitBlocks = true
  let timer = null
  const run = () => {
    clearTimeout(timer)
    timer = setTimeout(() => drawBlocks(root, assets), 60)
  }
  run()
  new MutationObserver((ms) => {
    if (ms.some((m) => [...m.addedNodes].some((n) => n.nodeType === 1 && (n.matches?.(mathSel + ',' + diagramSel + ',pre') ||
      n.querySelector?.(mathSel + ',' + diagramSel + ',' + codeSel))))) run()
  }).observe(root, { childList: true, subtree: true })
}
