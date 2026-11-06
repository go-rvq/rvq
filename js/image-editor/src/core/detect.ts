// "Detectar bordas" (the Perspectiva mode): the four corners of a document
// (a check, a sheet) on a background — a heuristic of contours: the image
// made small and gray, split in two by Otsu's threshold, the side that is not
// the background (the one touching the image's border the least), its
// largest connected region, and the extremes of it (x+y, x−y).

import type {Pixels, Quad} from './perspective'
import {luma} from './colors'

function otsu(hist: number[], total: number): number {
  let sum = 0
  for (let i = 0; i < 256; i++) sum += i * hist[i]
  let sumB = 0, wB = 0, best = 0, max = 0
  for (let t = 0; t < 256; t++) {
    wB += hist[t]
    if (wB === 0) continue
    const wF = total - wB
    if (wF === 0) break
    sumB += t * hist[t]
    const mB = sumB / wB
    const mF = (sum - sumB) / wF
    const between = wB * wF * (mB - mF) * (mB - mF)
    if (between > max) {
      max = between
      best = t
    }
  }
  return best
}

// detectCorners: the corners in the pixels of px, or null when nothing
// stands out.
export function detectCorners(px: Pixels, maxSide = 256): Quad | null {
  const k = Math.min(1, maxSide / Math.max(px.width, px.height))
  const w = Math.max(2, Math.round(px.width * k))
  const h = Math.max(2, Math.round(px.height * k))
  // the gray of the image made small (nearest sample)
  const g = new Uint8Array(w * h)
  const hist = new Array(256).fill(0)
  for (let y = 0; y < h; y++) {
    for (let x = 0; x < w; x++) {
      const sx = Math.min(px.width - 1, Math.floor(x / k))
      const sy = Math.min(px.height - 1, Math.floor(y / k))
      const o = (sy * px.width + sx) * 4
      const v = Math.round(luma(px.data[o], px.data[o + 1], px.data[o + 2]))
      g[y * w + x] = v
      hist[v]++
    }
  }
  const t = otsu(hist, w * h)
  // the side touching the border the least is the document
  let borderLight = 0, border = 0
  for (let x = 0; x < w; x++) {
    for (const y of [0, h - 1]) {
      border++
      if (g[y * w + x] > t) borderLight++
    }
  }
  for (let y = 0; y < h; y++) {
    for (const x of [0, w - 1]) {
      border++
      if (g[y * w + x] > t) borderLight++
    }
  }
  const light = borderLight / border < 0.5
  const fg = new Uint8Array(w * h)
  for (let i = 0; i < g.length; i++) fg[i] = (light ? g[i] > t : g[i] <= t) ? 1 : 0
  // the largest connected region (4-neighbours)
  const label = new Int32Array(w * h).fill(-1)
  let bestLabel = -1, bestSize = 0, n = 0
  const stack: number[] = []
  for (let i = 0; i < fg.length; i++) {
    if (!fg[i] || label[i] >= 0) continue
    let size = 0
    stack.push(i)
    label[i] = n
    while (stack.length) {
      const j = stack.pop()!
      size++
      const x = j % w, y = (j - x) / w
      const nb = [x > 0 ? j - 1 : -1, x < w - 1 ? j + 1 : -1, y > 0 ? j - w : -1, y < h - 1 ? j + w : -1]
      for (const q of nb) {
        if (q >= 0 && fg[q] && label[q] < 0) {
          label[q] = n
          stack.push(q)
        }
      }
    }
    if (size > bestSize) {
      bestSize = size
      bestLabel = n
    }
    n++
  }
  // too small (or the whole image): nothing stands out
  if (bestLabel < 0 || bestSize < w * h * 0.05 || bestSize > w * h * 0.98) return null
  let tl = {s: Infinity, x: 0, y: 0}, br = {s: -Infinity, x: 0, y: 0}
  let tr = {s: -Infinity, x: 0, y: 0}, bl = {s: Infinity, x: 0, y: 0}
  for (let i = 0; i < label.length; i++) {
    if (label[i] !== bestLabel) continue
    const x = i % w, y = (i - x) / w
    if (x + y < tl.s) tl = {s: x + y, x, y}
    if (x + y > br.s) br = {s: x + y, x, y}
    if (x - y > tr.s) tr = {s: x - y, x, y}
    if (x - y < bl.s) bl = {s: x - y, x, y}
  }
  const back = (p: {x: number; y: number}, dx: number, dy: number) => ({
    x: Math.min(px.width, Math.max(0, (p.x + dx) / k)),
    y: Math.min(px.height, Math.max(0, (p.y + dy) / k)),
  })
  return [back(tl, 0, 0), back(tr, 1, 0), back(br, 1, 1), back(bl, 0, 1)]
}
