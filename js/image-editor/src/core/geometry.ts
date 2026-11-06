// The geometry of the crops, in the image's pixels: the free polygon (its
// points, the rotation of the part cut out), the perspective's corners.

import type {Point} from './view'

export interface Rect {
  x: number
  y: number
  width: number
  height: number
}

export const MIN_POINTS = 3

export function dist(a: Point, b: Point): number {
  return Math.hypot(a.x - b.x, a.y - b.y)
}

// distToSegment: the distance from p to the segment a–b.
export function distToSegment(p: Point, a: Point, b: Point): number {
  const dx = b.x - a.x
  const dy = b.y - a.y
  const len2 = dx * dx + dy * dy
  if (len2 === 0) return dist(p, a)
  let t = ((p.x - a.x) * dx + (p.y - a.y) * dy) / len2
  t = Math.max(0, Math.min(1, t))
  return dist(p, {x: a.x + t * dx, y: a.y + t * dy})
}

// nearestPoint: the index of the point within radius of p (the nearest), -1
// when none.
export function nearestPoint(points: Point[], p: Point, radius: number): number {
  let best = -1
  let bestD = radius
  points.forEach((q, i) => {
    const d = dist(p, q)
    if (d <= bestD) {
      best = i
      bestD = d
    }
  })
  return best
}

// nearestEdge: the index i of the edge points[i]–points[i+1] (the path closed:
// the last goes to the first) nearest p — within radius when one is given.
export function nearestEdge(points: Point[], p: Point, radius = Infinity): number {
  let best = -1
  let bestD = radius
  for (let i = 0; i < points.length; i++) {
    const d = distToSegment(p, points[i], points[(i + 1) % points.length])
    if (d <= bestD) {
      best = i
      bestD = d
    }
  }
  return best
}

// insertPoint adds p to the closed path on the edge nearest it — so the path
// stays closed and does not cross itself —; with fewer than 2 points, at the
// end. It answers the new points and the index of p.
export function insertPoint(points: Point[], p: Point): {points: Point[]; index: number} {
  if (points.length < 2) return {points: [...points, p], index: points.length}
  const e = nearestEdge(points, p)
  const next = [...points]
  next.splice(e + 1, 0, p)
  return {points: next, index: e + 1}
}

// removePoint takes the point i out — never below MIN_POINTS: null then.
export function removePoint(points: Point[], i: number): Point[] | null {
  if (points.length <= MIN_POINTS || i < 0 || i >= points.length) return null
  return points.filter((_, j) => j !== i)
}

export function bounds(points: Point[]): Rect {
  if (!points.length) return {x: 0, y: 0, width: 0, height: 0}
  let x0 = Infinity, y0 = Infinity, x1 = -Infinity, y1 = -Infinity
  for (const p of points) {
    x0 = Math.min(x0, p.x)
    y0 = Math.min(y0, p.y)
    x1 = Math.max(x1, p.x)
    y1 = Math.max(y1, p.y)
  }
  return {x: x0, y: y0, width: x1 - x0, height: y1 - y0}
}

// center: the middle of the points' bounds — what the cut is turned around.
export function center(points: Point[]): Point {
  const b = bounds(points)
  return {x: b.x + b.width / 2, y: b.y + b.height / 2}
}

export function rotatePoint(p: Point, deg: number, c: Point): Point {
  const r = (deg * Math.PI) / 180
  const cos = Math.cos(r)
  const sin = Math.sin(r)
  const dx = p.x - c.x
  const dy = p.y - c.y
  return {x: c.x + dx * cos - dy * sin, y: c.y + dx * sin + dy * cos}
}

export function rotatePoints(points: Point[], deg: number, c = center(points)): Point[] {
  return points.map((p) => rotatePoint(p, deg, c))
}

// normDeg puts an angle in (-180, 180].
export function normDeg(a: number): number {
  let r = a % 360
  if (r <= -180) r += 360
  if (r > 180) r -= 360
  return r
}

// edgeAngle: the angle of a–b, in degrees (0: horizontal, to the right).
export function edgeAngle(a: Point, b: Point): number {
  return (Math.atan2(b.y - a.y, b.x - a.x) * 180) / Math.PI
}

// The magnet: an edge within SNAP_DEG of horizontal or vertical sticks there.
export const SNAP_DEG = 3

// snapRotation: the rotation deg of the cut, stuck to the nearest angle that
// makes one of its edges horizontal or vertical when one is within tol.
export function snapRotation(points: Point[], deg: number, tol = SNAP_DEG): number {
  let best = deg
  let bestD = tol + 1e-9
  for (let i = 0; i < points.length; i++) {
    const a = points[i]
    const b = points[(i + 1) % points.length]
    if (dist(a, b) === 0) continue
    const e = edgeAngle(a, b) + deg
    // to the nearest multiple of 90°
    const off = normDeg(e - Math.round(e / 90) * 90)
    if (Math.abs(off) < bestD) {
      bestD = Math.abs(off)
      best = deg - off
    }
  }
  return normDeg(best)
}

// stepRotation: with Shift, 15° by 15°.
export function stepRotation(deg: number, step = 15): number {
  return normDeg(Math.round(deg / step) * step)
}

// alignEdge: the rotation that makes edge i horizontal — the one of the two
// (the edge on top or at the bottom) nearest the rotation now.
export function alignEdge(points: Point[], i: number, current: number): number {
  const a = points[i]
  const b = points[(i + 1) % points.length]
  const e = edgeAngle(a, b)
  // e + r ≡ 0 (mod 180)
  let best = -e
  let bestD = Infinity
  for (let k = -2; k <= 2; k++) {
    const r = -e + k * 180
    const d = Math.abs(normDeg(r - current))
    if (d < bestD) {
      bestD = d
      best = r
    }
  }
  return normDeg(best)
}

// usefulArea: the rectangle the result is cut to — the bounds of the points
// turned by deg around their center —, in whole pixels.
export function usefulArea(points: Point[], deg: number): Rect {
  const b = bounds(rotatePoints(points, deg))
  const x = Math.floor(b.x + 1e-6)
  const y = Math.floor(b.y + 1e-6)
  return {
    x,
    y,
    width: Math.max(1, Math.ceil(b.x + b.width - 1e-6) - x),
    height: Math.max(1, Math.ceil(b.y + b.height - 1e-6) - y),
  }
}

// A rectangle's points, clockwise from the top left.
export function rectPoints(r: Rect): Point[] {
  return [
    {x: r.x, y: r.y},
    {x: r.x + r.width, y: r.y},
    {x: r.x + r.width, y: r.y + r.height},
    {x: r.x, y: r.y + r.height},
  ]
}

// clampRect keeps r in the image (w × h), its size kept when it fits.
export function clampRect(r: Rect, w: number, h: number): Rect {
  const width = Math.min(r.width, w)
  const height = Math.min(r.height, h)
  return {
    x: Math.min(Math.max(0, r.x), w - width),
    y: Math.min(Math.max(0, r.y), h - height),
    width,
    height,
  }
}

// fitAspect: the largest rectangle of ratio (width / height) centered in the
// image — the fixed crop's start.
export function fitAspect(w: number, h: number, ratio: number): Rect {
  if (!(ratio > 0)) return {x: 0, y: 0, width: w, height: h}
  let rw = w
  let rh = w / ratio
  if (rh > h) {
    rh = h
    rw = h * ratio
  }
  return {x: (w - rw) / 2, y: (h - rh) / 2, width: rw, height: rh}
}

// resizeRect moves the corner (or side) handle of r to p: handle is one of
// "nw", "n", "ne", "e", "se", "s", "sw", "w"; ratio, when given, is kept
// (the opposite corner fixed). min: the smallest side.
export function resizeRect(r: Rect, handle: string, p: Point, ratio = 0, min = 4): Rect {
  let x0 = r.x, y0 = r.y, x1 = r.x + r.width, y1 = r.y + r.height
  if (handle.includes('w')) x0 = Math.min(p.x, x1 - min)
  if (handle.includes('e')) x1 = Math.max(p.x, x0 + min)
  if (handle.includes('n')) y0 = Math.min(p.y, y1 - min)
  if (handle.includes('s')) y1 = Math.max(p.y, y0 + min)
  let out = {x: x0, y: y0, width: x1 - x0, height: y1 - y0}
  if (ratio > 0) {
    const horizontal = handle === 'e' || handle === 'w'
    const vertical = handle === 'n' || handle === 's'
    let w = out.width
    let h = out.height
    if (vertical) w = h * ratio
    else if (horizontal) h = w / ratio
    else if (w / h > ratio) h = w / ratio
    else w = h * ratio
    // the opposite side stays
    const x = handle.includes('w') ? x1 - w : horizontal || handle.includes('e') ? x0 : r.x + (r.width - w) / 2
    const y = handle.includes('n') ? y1 - h : vertical || handle.includes('s') ? y0 : r.y + (r.height - h) / 2
    out = {x, y, width: w, height: h}
  }
  return out
}

export function pointInPolygon(p: Point, poly: Point[]): boolean {
  let inside = false
  for (let i = 0, j = poly.length - 1; i < poly.length; j = i++) {
    const a = poly[i]
    const b = poly[j]
    if (a.y > p.y !== b.y > p.y && p.x < ((b.x - a.x) * (p.y - a.y)) / (b.y - a.y) + a.x) inside = !inside
  }
  return inside
}

// A rectangle turned: its center, its size, its angle (degrees, clockwise) —
// the crop "Retângulo": free, turned, always a rectangle.
export interface TurnedRect {
  cx: number
  cy: number
  width: number
  height: number
  angle: number
}

// turnedCorners: its corners in the image, clockwise from the top left.
export function turnedCorners(r: TurnedRect): Point[] {
  const c = {x: r.cx, y: r.cy}
  return rectPoints({x: r.cx - r.width / 2, y: r.cy - r.height / 2, width: r.width, height: r.height}).map((p) =>
    rotatePoint(p, r.angle, c),
  )
}

// toLocal: a point of the image in the rectangle's frame (its center the
// origin, its sides along the axes).
export function toLocal(r: TurnedRect, p: Point): Point {
  const q = rotatePoint(p, -r.angle, {x: r.cx, y: r.cy})
  return {x: q.x - r.cx, y: q.y - r.cy}
}

// resizeTurned moves the handle of r to p (the image's): the opposite side
// stays where it is; ratio, when given, is kept.
export function resizeTurned(r: TurnedRect, handle: string, p: Point, ratio = 0, min = 4): TurnedRect {
  const local = {x: -r.width / 2, y: -r.height / 2, width: r.width, height: r.height}
  const n = resizeRect(local, handle, toLocal(r, p), ratio, min)
  const lc = {x: n.x + n.width / 2, y: n.y + n.height / 2}
  const c = rotatePoint(lc, r.angle, {x: 0, y: 0})
  return {cx: r.cx + c.x, cy: r.cy + c.y, width: n.width, height: n.height, angle: r.angle}
}

// snapAngle: the angle stuck to 0°, 90°… within tol.
export function snapAngle(deg: number, tol = SNAP_DEG): number {
  const off = normDeg(deg - Math.round(deg / 90) * 90)
  return Math.abs(off) <= tol ? normDeg(deg - off) : normDeg(deg)
}
