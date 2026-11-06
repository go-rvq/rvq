// The perspective: four corners of a quadrilateral of the image (top left,
// top right, bottom right, bottom left) made a rectangle — a projective
// transformation, the homography of the rectangle's corners to them.

import type {Point} from './view'
import {dist} from './geometry'

export type Quad = [Point, Point, Point, Point]

// Matrix3: row major, [a b c; d e f; g h 1].
export type Matrix3 = number[]

// outputSize: the rectangle's size — the mean of the opposite sides.
export function outputSize(q: Quad): {width: number; height: number} {
  const [tl, tr, br, bl] = q
  return {
    width: Math.max(1, Math.round((dist(tl, tr) + dist(bl, br)) / 2)),
    height: Math.max(1, Math.round((dist(tl, bl) + dist(tr, br)) / 2)),
  }
}

// solve: Ax = b by Gauss with partial pivoting (n×n).
function solve(A: number[][], b: number[]): number[] {
  const n = b.length
  const M = A.map((row, i) => [...row, b[i]])
  for (let c = 0; c < n; c++) {
    let p = c
    for (let r = c + 1; r < n; r++) if (Math.abs(M[r][c]) > Math.abs(M[p][c])) p = r
    if (Math.abs(M[p][c]) < 1e-12) throw new Error('perspective: degenerate quadrilateral')
    ;[M[c], M[p]] = [M[p], M[c]]
    for (let r = 0; r < n; r++) {
      if (r === c) continue
      const f = M[r][c] / M[c][c]
      for (let k = c; k <= n; k++) M[r][k] -= f * M[c][k]
    }
  }
  return M.map((row, i) => row[n] / row[i])
}

// homography: the matrix taking each from[i] to to[i].
export function homography(from: Point[], to: Point[]): Matrix3 {
  const A: number[][] = []
  const b: number[] = []
  for (let i = 0; i < 4; i++) {
    const {x, y} = from[i]
    const {x: u, y: v} = to[i]
    A.push([x, y, 1, 0, 0, 0, -u * x, -u * y])
    b.push(u)
    A.push([0, 0, 0, x, y, 1, -v * x, -v * y])
    b.push(v)
  }
  return [...solve(A, b), 1]
}

export function apply(m: Matrix3, p: Point): Point {
  const w = m[6] * p.x + m[7] * p.y + m[8]
  return {x: (m[0] * p.x + m[1] * p.y + m[2]) / w, y: (m[3] * p.x + m[4] * p.y + m[5]) / w}
}

// rectToQuad: the homography of the output rectangle (width × height) to the
// quadrilateral — each pixel of the result read where it lands in the image.
export function rectToQuad(q: Quad, width: number, height: number): Matrix3 {
  return homography(
    [
      {x: 0, y: 0},
      {x: width, y: 0},
      {x: width, y: height},
      {x: 0, y: height},
    ],
    q,
  )
}

export interface Pixels {
  width: number
  height: number
  data: Uint8ClampedArray
}

// warp: the quadrilateral q of src made a rectangle (outputSize, or size):
// each pixel sampled bilinearly where the homography puts it; outside the
// image, transparent.
export function warp(src: Pixels, q: Quad, size = outputSize(q)): Pixels {
  const {width, height} = size
  const m = rectToQuad(q, width, height)
  const out = new Uint8ClampedArray(width * height * 4)
  const sw = src.width
  const sh = src.height
  const s = src.data
  for (let y = 0; y < height; y++) {
    for (let x = 0; x < width; x++) {
      const p = apply(m, {x: x + 0.5, y: y + 0.5})
      const fx = p.x - 0.5
      const fy = p.y - 0.5
      const o = (y * width + x) * 4
      if (fx < -0.5 || fy < -0.5 || fx > sw - 0.5 || fy > sh - 0.5) continue
      const x0 = Math.max(0, Math.min(sw - 1, Math.floor(fx)))
      const y0 = Math.max(0, Math.min(sh - 1, Math.floor(fy)))
      const x1 = Math.min(sw - 1, x0 + 1)
      const y1 = Math.min(sh - 1, y0 + 1)
      const ax = Math.max(0, Math.min(1, fx - x0))
      const ay = Math.max(0, Math.min(1, fy - y0))
      for (let c = 0; c < 4; c++) {
        const v00 = s[(y0 * sw + x0) * 4 + c]
        const v10 = s[(y0 * sw + x1) * 4 + c]
        const v01 = s[(y1 * sw + x0) * 4 + c]
        const v11 = s[(y1 * sw + x1) * 4 + c]
        out[o + c] = (v00 * (1 - ax) + v10 * ax) * (1 - ay) + (v01 * (1 - ax) + v11 * ax) * ay
      }
    }
  }
  return {width, height, data: out}
}

// defaultQuad: the image's corners (inset: a share in from them).
export function defaultQuad(w: number, h: number, inset = 0): Quad {
  const dx = w * inset
  const dy = h * inset
  return [
    {x: dx, y: dy},
    {x: w - dx, y: dy},
    {x: w - dx, y: h - dy},
    {x: dx, y: h - dy},
  ]
}

// isWholeImage: the quadrilateral is the image's corners (within tol px) —
// no crop.
export function isWholeImage(q: Quad, w: number, h: number, tol = 0.5): boolean {
  const c = defaultQuad(w, h)
  return q.every((p, i) => Math.abs(p.x - c[i].x) <= tol && Math.abs(p.y - c[i].y) <= tol)
}
