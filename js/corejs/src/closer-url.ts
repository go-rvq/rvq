// Address bar for overlays: a closer that carries a `url` puts that URL in the
// address bar while it is open, and restores the previous one when it closes —
// so a dialog/drawer showing a listing, a record's detail or an edit form reads
// like the PAGE that renders the same thing (a clean URL, not an action), and
// the address bar never keeps the URL of something that is no longer on screen.
//
// Overlays stack, so the addresses stack the same way: LIFO. Opening pushes an
// entry with the URL that was current at the time; closing an entry restores
// ITS previous URL and drops everything opened above it (those overlays live
// inside the one being closed — they go away with it).
//
// Going Back is the same gesture as closing the top overlay, so a popstate
// closes it (with the history writes suppressed: the browser already moved).

import type { CloserUrl } from './closer'

interface Entry {
  closer: any
  url: string
  prev: string
}

const stack: Entry[] = []

// set while the browser (not us) is driving the change
let suppressed = false
let popstateInstalled = false

function resolveUrl(closer: any): string {
  const url: CloserUrl | undefined = closer?.url
  if (!url) {
    return ''
  }
  try {
    return typeof url === 'function' ? url(closer) || '' : String(url)
  } catch (e) {
    console.error('closer url failed:', e)
    return ''
  }
}

function currentHref(): string {
  return window.location.pathname + window.location.search + window.location.hash
}

function ensurePopstate() {
  if (popstateInstalled || typeof window === 'undefined') {
    return
  }
  popstateInstalled = true
  window.addEventListener('popstate', () => {
    // Back with overlays open closes the top one. Its close callback runs while
    // suppressed, so it unwinds the stack without writing history again.
    const top = stack[stack.length - 1]
    if (!top) {
      return
    }
    suppressed = true
    try {
      top.closer.show = false
    } finally {
      suppressed = false
    }
  })
}

/** Entries currently stacked — the address bar's LIFO. Exposed for tests. */
export function closerUrlStack(): { url: string; prev: string }[] {
  return stack.map((e) => ({ url: e.url, prev: e.prev }))
}

export function closerUrlOpen(closer: any) {
  const url = resolveUrl(closer)
  if (!url) {
    return
  }
  ensurePopstate()
  stack.push({ closer, url, prev: currentHref() })
  if (!suppressed) {
    // state stays null on purpose: plaid's popstate handler only reacts to a
    // state it pushed itself, so navigating here must not reload the page.
    window.history.pushState(null, '', url)
  }
}

export function closerUrlClose(closer: any) {
  const i = stack.findIndex((e) => e.closer === closer)
  if (i < 0) {
    return
  }
  const { prev } = stack[i]
  stack.splice(i)
  if (!suppressed) {
    // replace, not push: the overlay's address is not a place to come back to.
    window.history.replaceState(null, '', prev)
  }
}

/** Registers the address-bar hooks on a closer (no-op while it has no url). */
export function installCloserUrlSync(closer: any) {
  closer.openCallbacks.push(closerUrlOpen)
  closer.closeCallbacks.push(closerUrlClose)
}

/** Test helper: forgets every stacked address. */
export function resetCloserUrlStack() {
  stack.splice(0)
}
