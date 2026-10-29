import { config } from './config.js'

// footnoteBalloons shows the note of a reference of a footnote (a.footnote-ref,
// its note the element of its #id) in a balloon by it: the pointer over it,
// or the focus on it; a touch shows it — a second one follows the link —; a
// touch elsewhere, Escape, the pointer out, hide it.
export function footnoteBalloons() {
  if (window.__mdkitFootnotes) return
  window.__mdkitFootnotes = true
  let balloon = null
  let shownFor = null
  const hide = () => {
    if (balloon) balloon.remove()
    balloon = null
    shownFor = null
  }
  const show = (a) => {
    const id = decodeURIComponent((a.getAttribute('href') || '').replace(/^#/, ''))
    const note = id && document.getElementById(id)
    if (!note) return false
    hide()
    balloon = document.createElement('div')
    balloon.className = 'mdkit-fn-balloon ' + config.textClass
    balloon.setAttribute('role', 'tooltip')
    balloon.innerHTML = note.innerHTML
    document.body.append(balloon)
    const r = a.getBoundingClientRect()
    const w = balloon.offsetWidth
    const left = Math.max(8, Math.min(window.scrollX + r.left - w / 2 + r.width / 2, window.scrollX + document.documentElement.clientWidth - w - 8))
    balloon.style.left = left + 'px'
    balloon.style.top = (window.scrollY + r.bottom + 6) + 'px'
    shownFor = a
    return true
  }
  const refOf = (e) => e.target && e.target.closest && e.target.closest('a.footnote-ref')
  document.addEventListener('mouseover', (e) => { const a = refOf(e); if (a && a !== shownFor) show(a) })
  document.addEventListener('mouseout', (e) => { const a = refOf(e); if (a && a === shownFor) hide() })
  document.addEventListener('focusin', (e) => { const a = refOf(e); if (a) show(a) })
  document.addEventListener('focusout', (e) => { if (refOf(e)) hide() })
  document.addEventListener('keydown', (e) => { if (e.key === 'Escape') hide() })
  document.addEventListener('touchstart', (e) => {
    const a = refOf(e)
    if (!a) return hide()
    if (a !== shownFor && show(a)) e.preventDefault()
  }, { passive: false })
}
