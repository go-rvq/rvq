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
// back to the address an overlay was opened from closes it — and what was
// opened above it: a jump of several steps back (with the history writes
// suppressed: the browser already moved) —, as the user's close: a form with
// changes not saved holds it, asks, and keeps its address. Forward to the address of an
// overlay Back closed shows nothing there — that overlay is gone —, so the
// browser is sent back to where it was: forward is a no-op, never a close.

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
// the addresses of the overlays Back closed: what Forward would go to
const forwardUrls: string[] = []
// set while we send the browser back from such a Forward (its popstate ours)
let bouncing = false

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
    if (bouncing) {
      bouncing = false
      return
    }
    const here = currentHref()
    // Back to where an overlay was opened from: it closes, and what is above
    // it — from the top down, each as a close the user asked for (requestClose):
    // one holding (a form with changes not saved) stays, its confirmation
    // shown, and its address comes back (the browser moved already). The
    // close callbacks run while suppressed, so they unwind the stack without
    // writing history again.
    for (let i = stack.length - 1; i >= 0; i--) {
      if (stack[i].prev !== here) {
        continue
      }
      for (let j = stack.length - 1; j >= i; j--) {
        const e = stack[j]
        if (!e) {
          continue // gone with one below it
        }
        let closed = true
        suppressed = true
        try {
          if (typeof e.closer.requestClose === 'function') {
            closed = e.closer.requestClose() !== false
          } else {
            e.closer.show = false
          }
        } finally {
          suppressed = false
        }
        if (!closed) {
          window.history.pushState(null, '', e.url)
          return
        }
        if (!forwardUrls.includes(e.url)) {
          forwardUrls.push(e.url)
        }
      }
      return
    }
    // Forward to an overlay Back closed: back to where the browser was
    if (forwardUrls.includes(here)) {
      bouncing = true
      window.history.back()
    }
  })
}

/**
 * Says the popstate under way is the overlays' — Back to where an overlay was
 * opened from, Forward to one Back closed, or the browser sent back from it —,
 * not a navigation of the page: plaid must not reload the page for it (its
 * popstate handler runs first, before ours unwinds the stack), which would
 * take away every overlay, not only the ones the Back closes.
 */
export function closerUrlOwnsPopstate(): boolean {
  if (bouncing) {
    return true
  }
  const here = currentHref()
  return stack.some((e) => e.prev === here) || forwardUrls.includes(here)
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
    // a new entry drops what was ahead: nothing to go forward to
    forwardUrls.splice(0)
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
  forwardUrls.splice(0)
  bouncing = false
}
