// The width of a resizable drawer (NavigationDrawer's `resizable`): dragged by
// the handle on its inner edge, kept between its min and max, saved under its
// key in localStorage — what a reload, or another page, starts from.

// isEndSide says the drawer sits on the right (its handle on its left edge):
// a drag to the left widens it.
export function isEndSide(location?: string): boolean {
  return location === 'right' || location === 'end'
}

// pixelWidth is a width given in pixels: a percentage ("92%") of the
// viewport's width, a number or "320px" as it is; undefined when it is none of
// them (Vuetify's drawer takes only pixels: a "92%" made its width NaN, and a
// closed drawer stayed on the page — translateX(NaN) is no transform).
export function pixelWidth(width: unknown, viewport: number): number | undefined {
  if (typeof width === 'number') return Number.isFinite(width) ? width : undefined
  if (typeof width !== 'string' || width.trim() === '') return undefined
  const s = width.trim()
  if (s.endsWith('%')) {
    const p = Number(s.slice(0, -1))
    return Number.isFinite(p) ? Math.round((viewport * p) / 100) : undefined
  }
  const n = Number(s.endsWith('px') ? s.slice(0, -2) : s)
  return Number.isFinite(n) ? n : undefined
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
