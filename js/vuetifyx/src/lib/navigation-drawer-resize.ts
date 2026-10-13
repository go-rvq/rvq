// The width of a resizable drawer (NavigationDrawer's `resizable`): dragged by
// the handle on its inner edge, kept between its min and max, saved under its
// key in localStorage — what a reload, or another page, starts from.

// isEndSide says the drawer sits on the right (its handle on its left edge):
// a drag to the left widens it.
export function isEndSide(location?: string): boolean {
  return location === 'right' || location === 'end'
}

// dragWidth is the width a drag of dx pixels from a start width makes, kept
// in [min, max].
export function dragWidth(
  start: number,
  dx: number,
  min: number,
  max: number,
  endSide: boolean
): number {
  const w = start + (endSide ? -dx : dx)
  return Math.round(Math.min(Math.max(w, min), max))
}

// savedWidth is the width saved under key, kept in [min, max]; null when there
// is none (or the storage cannot be read: a private window).
export function savedWidth(key: string, min: number, max: number): number | null {
  if (!key) return null
  try {
    const v = Number(window.localStorage.getItem(key))
    if (!v || !Number.isFinite(v)) return null
    return Math.min(Math.max(v, min), max)
  } catch {
    return null
  }
}

// saveWidth saves the width under key; null takes it out (the default again).
export function saveWidth(key: string, width: number | null) {
  if (!key) return
  try {
    if (width == null) window.localStorage.removeItem(key)
    else window.localStorage.setItem(key, String(width))
  } catch {
    // no storage: the width holds until the page goes
  }
}
