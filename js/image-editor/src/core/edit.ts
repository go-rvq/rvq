// The edit of an image: what the dialog does to it, in the order the image is
// made — the turn (90°, mirror), the crop (geometry), the colors, the
// annotations on top. A plain value (JSON): the history's states, and what is
// saved beside the original to open the edit again.

import type {Point} from './view'
import type {Rect} from './geometry'
import {rotatePoint, center, usefulArea, turnedCorners, type TurnedRect} from './geometry'
import type {Quad, Pixels} from './perspective'
import {isWholeImage, outputSize, warp} from './perspective'
import type {Colors} from './colors'
import {applyColors, isNeutral, neutral} from './colors'
import type {Annotation} from './annotations'
import {drawAll} from './annotations'

export type CropMode = 'fixo' | 'retangulo' | 'perspectiva' | 'livre'
export type Tab = 'recortar' | 'cores' | 'anotar'

export interface Edit {
  version: 1
  // the image turned (clockwise) and mirrored, before anything
  rotate: 0 | 90 | 180 | 270
  flipH: boolean
  mode: CropMode
  // the fixed crop: a rectangle (of the size's ratio, when it has one)
  fixo: Rect | null
  // its ratio (width / height) when it is no size's: 0, free; -1, the
  // image's own
  fixoRatio?: number
  // the free rectangle, turned: the result is it, straight
  retangulo?: TurnedRect | null
  perspectiva: Quad | null
  livre: {points: Point[]; rotation: number} | null
  colors: Colors
  annotations: Annotation[]
}

export function emptyEdit(mode: CropMode = 'fixo'): Edit {
  return {version: 1, rotate: 0, flipH: false, mode, fixo: null, perspectiva: null, livre: null, colors: neutral(), annotations: []}
}

// parseEdit: an edit saved (JSON), completed with what it lacks.
export function parseEdit(v: unknown, mode: CropMode = 'fixo'): Edit {
  const e = emptyEdit(mode)
  if (!v || typeof v !== 'object') return e
  const o = v as Partial<Edit>
  return {
    ...e,
    ...o,
    version: 1,
    colors: {...e.colors, ...(o.colors || {})},
    annotations: Array.isArray(o.annotations) ? o.annotations : [],
  }
}

// geometryOf: whether the crop changes the image — the mode's shape is set
// (a fixed crop of a size is not: it is the size's, CropOptions).
export function hasGeometry(e: Edit, fixoIsSize: boolean): boolean {
  if (e.rotate || e.flipH) return true
  switch (e.mode) {
    case 'fixo':
      return !fixoIsSize && !!e.fixo
    case 'retangulo':
      return !!e.retangulo
    case 'perspectiva':
      return !!e.perspectiva
    case 'livre':
      return !!e.livre && e.livre.points.length >= 3
  }
}

// changes: whether the edit makes a new image (else only the fixed crop of a
// size, or nothing).
export function changes(e: Edit, fixoIsSize: boolean): boolean {
  return hasGeometry(e, fixoIsSize) || !isNeutral(e.colors) || e.annotations.length > 0
}

// transparent: the result has transparent pixels (a PNG): the free crop.
export function transparent(e: Edit): boolean {
  return e.mode === 'livre' && !!e.livre && e.livre.points.length >= 3
}

// normalized: the edit with what does not cut taken out — a perspective on
// the image's own corners (w × h, the image turned), a free path of fewer
// than 3 points.
export function normalized(e: Edit, w: number, h: number): Edit {
  const out = {...e}
  if (out.fixo && Math.abs(out.fixo.x) <= 0.5 && Math.abs(out.fixo.y) <= 0.5 && Math.abs(out.fixo.width - w) <= 0.5 && Math.abs(out.fixo.height - h) <= 0.5) out.fixo = null
  if (out.perspectiva && isWholeImage(out.perspectiva, w, h)) out.perspectiva = null
  const r = out.retangulo
  if (r && Math.abs(r.angle) < 0.05 && Math.abs(r.cx - w / 2) <= 0.5 && Math.abs(r.cy - h / 2) <= 0.5 && Math.abs(r.width - w) <= 0.5 && Math.abs(r.height - h) <= 0.5) out.retangulo = null
  if (out.livre && out.livre.points.length < 3) out.livre = null
  return out
}

// --- the canvas (the browser) ---

export type Source = CanvasImageSource & {width: number; height: number}

export function canvas(w: number, h: number): HTMLCanvasElement {
  const c = document.createElement('canvas')
  c.width = Math.max(1, Math.round(w))
  c.height = Math.max(1, Math.round(h))
  return c
}

function sizeOf(s: CanvasImageSource): {width: number; height: number} {
  const any = s as any
  return {width: any.naturalWidth || any.videoWidth || any.width, height: any.naturalHeight || any.videoHeight || any.height}
}

// oriented: the image turned and mirrored.
export function oriented(src: CanvasImageSource, rotate: number, flipH: boolean): HTMLCanvasElement {
  const {width: w, height: h} = sizeOf(src)
  const swap = rotate === 90 || rotate === 270
  const c = canvas(swap ? h : w, swap ? w : h)
  const ctx = c.getContext('2d')!
  ctx.translate(c.width / 2, c.height / 2)
  ctx.rotate((rotate * Math.PI) / 180)
  if (flipH) ctx.scale(-1, 1)
  ctx.drawImage(src, -w / 2, -h / 2)
  return c
}

function pixelsOf(c: HTMLCanvasElement): Pixels {
  const d = c.getContext('2d')!.getImageData(0, 0, c.width, c.height)
  return {width: d.width, height: d.height, data: d.data}
}

function canvasOf(p: Pixels): HTMLCanvasElement {
  const c = canvas(p.width, p.height)
  c.getContext('2d')!.putImageData(new ImageData(p.data as any, p.width, p.height), 0, 0)
  return c
}

// cropped: the image (oriented) cut by the crop of the edit — k scales the
// edit's coordinates (a preview made smaller).
export function cropped(src: HTMLCanvasElement, e: Edit, fixoIsSize: boolean, k = 1): HTMLCanvasElement {
  switch (e.mode) {
    case 'fixo': {
      if (fixoIsSize || !e.fixo) return src
      const r = e.fixo
      const c = canvas(r.width * k, r.height * k)
      c.getContext('2d')!.drawImage(src, r.x * k, r.y * k, r.width * k, r.height * k, 0, 0, c.width, c.height)
      return c
    }
    case 'perspectiva': {
      if (!e.perspectiva) return src
      const q = e.perspectiva.map((p) => ({x: p.x * k, y: p.y * k})) as Quad
      return canvasOf(warp(pixelsOf(src), q, outputSize(q)))
    }
    case 'retangulo':
    case 'livre': {
      // the rectangle: a path of its corners, turned back straight
      const shape = e.mode === 'retangulo' ? (e.retangulo ? {points: turnedCorners(e.retangulo), rotation: -e.retangulo.angle} : null) : e.livre
      if (!shape || shape.points.length < 3) return src
      const pts = shape.points.map((p) => ({x: p.x * k, y: p.y * k}))
      const deg = shape.rotation || 0
      const area = usefulArea(pts, deg)
      const c = canvas(area.width, area.height)
      const ctx = c.getContext('2d')!
      const ctr = center(pts)
      ctx.translate(-area.x, -area.y)
      // turned around the center of the cut
      ctx.translate(ctr.x, ctr.y)
      ctx.rotate((deg * Math.PI) / 180)
      ctx.translate(-ctr.x, -ctr.y)
      ctx.beginPath()
      pts.forEach((p, i) => (i ? ctx.lineTo(p.x, p.y) : ctx.moveTo(p.x, p.y)))
      ctx.closePath()
      ctx.clip()
      ctx.drawImage(src, 0, 0)
      return c
    }
  }
}

// croppedSize: the size of the image as cut — the pixels the annotations
// are in —, without making it (the oriented image of w × h).
export function croppedSize(e: Edit, w: number, h: number, fixoIsSize: boolean): {width: number; height: number} {
  switch (e.mode) {
    case 'fixo':
      return !fixoIsSize && e.fixo ? {width: Math.round(e.fixo.width), height: Math.round(e.fixo.height)} : {width: w, height: h}
    case 'retangulo':
      return e.retangulo ? {width: Math.max(1, Math.round(e.retangulo.width)), height: Math.max(1, Math.round(e.retangulo.height))} : {width: w, height: h}
    case 'perspectiva':
      return e.perspectiva ? outputSize(e.perspectiva) : {width: w, height: h}
    case 'livre':
      if (!e.livre || e.livre.points.length < 3) return {width: w, height: h}
      {
        const a = usefulArea(e.livre.points, e.livre.rotation || 0)
        return {width: a.width, height: a.height}
      }
  }
}

// orientedSizeOf: the size of the image of w × h turned by rotate.
export function turnedSize(w: number, h: number, rotate: number): {width: number; height: number} {
  return rotate === 90 || rotate === 270 ? {width: h, height: w} : {width: w, height: h}
}

// colored: the colors of the edit on the image.
export function colored(src: HTMLCanvasElement, e: Edit): HTMLCanvasElement {
  if (isNeutral(e.colors)) return src
  return canvasOf(applyColors(pixelsOf(src), e.colors))
}

// annotated: the annotations over the image (scaled by k).
export function annotated(src: HTMLCanvasElement, e: Edit, k = 1): HTMLCanvasElement {
  if (!e.annotations.length) return src
  const c = canvas(src.width, src.height)
  const ctx = c.getContext('2d')!
  ctx.drawImage(src, 0, 0)
  if (k !== 1) ctx.scale(k, k)
  drawAll(ctx, e.annotations, src, k)
  return c
}

// render: the image made by the edit, at full resolution (k = 1) or a
// preview (k < 1, the source given made smaller by k).
export function render(src: CanvasImageSource, e: Edit, fixoIsSize: boolean, k = 1): HTMLCanvasElement {
  const o = oriented(src, e.rotate, e.flipH)
  return annotated(colored(cropped(o, e, fixoIsSize, k), e), e, k)
}

// rotatedPoint: where a point of the cut lands in the free crop's result.
export function livrePoint(e: Edit, p: Point): Point {
  if (!e.livre) return p
  const pts = e.livre.points
  const area = usefulArea(pts, e.livre.rotation)
  const q = rotatePoint(p, e.livre.rotation, center(pts))
  return {x: q.x - area.x, y: q.y - area.y}
}

// toBlob: the result as a file — PNG when it has transparency (the free
// crop), else the type of the original (JPEG, WebP at quality).
export function toBlob(c: HTMLCanvasElement, type: string, quality = 0.92): Promise<Blob> {
  return new Promise((resolve, reject) =>
    c.toBlob((b) => (b ? resolve(b) : reject(new Error('image-editor: no image'))), type, quality),
  )
}

export function resultType(e: Edit, original: string): string {
  if (transparent(e)) return 'image/png'
  if (original === 'image/png' || original === 'image/webp' || original === 'image/jpeg') return original
  return 'image/jpeg'
}

export function resultName(name: string, type: string): string {
  const ext = type === 'image/png' ? '.png' : type === 'image/webp' ? '.webp' : '.jpg'
  const base = name.replace(/\.[^./]+$/, '') || 'image'
  return base + ext
}
