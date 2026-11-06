// The limit of an image uploaded (MediaBoxConfig.OriginalLimit): the largest
// it is sent in, its proportions kept — the limit turned with the image (a
// portrait one fits it on its side), as the server checks it.

export interface Limit {
  width: number
  height: number
}

// limitFor: the limit for an image w × h — turned when the image is turned
// the other way.
export function limitFor(w: number, h: number, max: Limit): Limit {
  const {width, height} = max
  if (width > 0 && height > 0 && w !== h && width !== height && w > h !== width > height) {
    return {width: height, height: width}
  }
  return {width, height}
}

// fitSize: w × h made smaller to fit max (turned as limitFor); ratio, the
// scale (1: it fits, untouched). A side of the limit 0 or less: no limit on
// it.
export function fitSize(w: number, h: number, max: Limit | null | undefined): {width: number; height: number; ratio: number} {
  if (!max || !(w > 0) || !(h > 0)) return {width: w, height: h, ratio: 1}
  const l = limitFor(w, h, max)
  let ratio = 1
  if (l.width > 0 && w > l.width) ratio = l.width / w
  if (l.height > 0 && h > l.height) ratio = Math.min(ratio, l.height / h)
  if (ratio === 1) return {width: w, height: h, ratio}
  return {width: Math.max(1, Math.round(w * ratio)), height: Math.max(1, Math.round(h * ratio)), ratio}
}

// scaleRect: a rectangle of the image (a fixed crop) at the image's new
// scale, in whole pixels.
export function scaleRect<T extends {x: number; y: number; width: number; height: number}>(r: T, ratio: number): T {
  if (ratio === 1) return r
  return {
    ...r,
    x: Math.round(r.x * ratio),
    y: Math.round(r.y * ratio),
    width: Math.max(1, Math.round(r.width * ratio)),
    height: Math.max(1, Math.round(r.height * ratio)),
  }
}
