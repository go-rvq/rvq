// The zoom of the stage: the image drawn at scale, moved by (x, y) — its
// origin on the screen. The marks (points, handles, lines) are kept in the
// image's pixels and drawn the same size on the screen whatever the zoom:
// their sizes are divided by the scale.

export interface Point {
  x: number
  y: number
}

export interface View {
  scale: number
  x: number
  y: number
}

export const MIN_SCALE = 0.1
export const MAX_SCALE = 16

export function clampScale(s: number): number {
  if (!isFinite(s) || s <= 0) return MIN_SCALE
  return Math.min(MAX_SCALE, Math.max(MIN_SCALE, s))
}

// screenToImage: a point of the stage (its pixels) in the image's pixels.
export function screenToImage(v: View, p: Point): Point {
  return {x: (p.x - v.x) / v.scale, y: (p.y - v.y) / v.scale}
}

export function imageToScreen(v: View, p: Point): Point {
  return {x: p.x * v.scale + v.x, y: p.y * v.scale + v.y}
}

// zoomAt sets the scale keeping the point of the image under at (the
// cursor, the middle of a pinch) where it is.
export function zoomAt(v: View, scale: number, at: Point): View {
  const s = clampScale(scale)
  const img = screenToImage(v, at)
  return {scale: s, x: at.x - img.x * s, y: at.y - img.y * s}
}

// fit: the whole image in the stage, centered (never above 100% for a small
// image... unless it does not fill: then up to it).
export function fit(stageW: number, stageH: number, imgW: number, imgH: number, pad = 16): View {
  if (imgW <= 0 || imgH <= 0) return {scale: 1, x: 0, y: 0}
  const s = clampScale(Math.min((stageW - 2 * pad) / imgW, (stageH - 2 * pad) / imgH))
  return {scale: s, x: (stageW - imgW * s) / 2, y: (stageH - imgH * s) / 2}
}

// actual: 100% (1 pixel of the image, 1 of the screen), centered.
export function actual(stageW: number, stageH: number, imgW: number, imgH: number): View {
  return {scale: 1, x: (stageW - imgW) / 2, y: (stageH - imgH) / 2}
}

// The steps of the − / + buttons.
const STEPS = [0.1, 0.125, 0.167, 0.25, 0.333, 0.5, 0.667, 0.75, 1, 1.25, 1.5, 2, 3, 4, 6, 8, 12, 16]

export function nextStep(scale: number, dir: 1 | -1): number {
  if (dir > 0) return STEPS.find((s) => s > scale * 1.001) ?? MAX_SCALE
  return [...STEPS].reverse().find((s) => s < scale / 1.001) ?? MIN_SCALE
}

// onScreen: a size of the screen (a handle's radius, a line's width) in the
// image's pixels — what the marks are drawn with, so they keep their size.
export function onScreen(v: View, px: number): number {
  return px / v.scale
}

// The radius a point or a handle is grabbed in, on the screen.
export const HIT_RADIUS = 12
export const HIT_RADIUS_TOUCH = 20
