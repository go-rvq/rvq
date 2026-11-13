// notify tells something in a small note at the bottom of the page, gone in
// a moment — discreet: it takes no focus, asks nothing (an error of a card
// not read, say). Its text as text (never HTML).
export function notify(text, { ms = 4500, kind = 'warning', doc = (typeof document !== 'undefined' ? document : null) } = {}) {
  if (!doc || !doc.body || !text) return null
  const el = doc.createElement('div')
  el.setAttribute('role', 'status')
  el.setAttribute('data-mdkit-notice', kind)
  el.textContent = text
  el.style.cssText =
    'position:fixed;left:16px;bottom:16px;z-index:3000;max-width:min(32rem,calc(100vw - 32px));' +
    'padding:6px 12px;border-radius:6px;font:13px/1.4 sans-serif;color:#fff;opacity:0;transition:opacity .2s;' +
    'box-shadow:0 2px 8px rgba(0,0,0,.25);pointer-events:none;background:' + (kind === 'error' ? '#b3261e' : '#5f4b00')
  doc.body.append(el)
  const show = () => (el.style.opacity = '0.95')
  if (typeof requestAnimationFrame !== 'undefined') requestAnimationFrame(show)
  else show()
  setTimeout(() => {
    el.style.opacity = '0'
    setTimeout(() => el.remove(), 250)
  }, ms)
  return el
}
