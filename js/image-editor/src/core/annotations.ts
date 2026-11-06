// The annotations (the Anotar tab, as Spectacle's): vectors in the pixels of
// the image as cut (after the crop), drawn over it — the colors never touch
// them. Each is content of the image: it grows with the zoom.

import type {Point} from './view'
import type {Rect} from './geometry'
import {bounds as pointsBounds, distToSegment} from './geometry'

export type Tool =
  | 'select'
  | 'pen'
  | 'highlighter'
  | 'line'
  | 'arrow'
  | 'rect'
  | 'ellipse'
  | 'text'
  | 'number'
  | 'pixelate'
  | 'blur'

export const TOOLS: Tool[] = ['select', 'pen', 'highlighter', 'line', 'arrow', 'rect', 'ellipse', 'text', 'number', 'pixelate', 'blur']

export interface Annotation {
  id: string
  type: Exclude<Tool, 'select'>
  // pen, highlighter: the path; line, arrow: from, to; rect, ellipse,
  // pixelate, blur: two corners; text, number: where
  points: Point[]
  stroke: string
  fill: string // '' none
  width: number
  opacity: number // 0..1
  shadow: boolean
  font: string
  fontSize: number
  text: string
  n: number // number: its count
}

export interface Style {
  stroke: string
  fill: string
  width: number
  opacity: number
  shadow: boolean
  font: string
  fontSize: number
}

export function defaultStyle(): Style {
  return {stroke: '#e53935', fill: '', width: 4, opacity: 1, shadow: false, font: 'sans-serif', fontSize: 32}
}

let seq = 0
export function newId(): string {
  seq++
  return `a${Date.now().toString(36)}${seq}`
}

export function create(type: Annotation['type'], at: Point, style: Style, n = 0): Annotation {
  const a: Annotation = {
    id: newId(),
    type,
    points: [at, at],
    ...style,
    text: '',
    n,
  }
  if (type === 'pen' || type === 'highlighter') a.points = [at]
  if (type === 'text' || type === 'number') a.points = [at]
  if (type === 'highlighter') {
    a.stroke = style.stroke === '#e53935' ? '#ffeb3b' : style.stroke
    a.width = Math.max(style.width, 18)
    a.opacity = Math.min(style.opacity, 0.45)
  }
  if (type === 'number') a.fill = a.fill || a.stroke
  return a
}

// nextNumber: the count of a new number — the highest + 1.
export function nextNumber(list: Annotation[]): number {
  return list.filter((a) => a.type === 'number').reduce((m, a) => Math.max(m, a.n), 0) + 1
}

function textSize(a: Annotation): {width: number; height: number} {
  const lines = (a.text || ' ').split('\n')
  const w = Math.max(...lines.map((l) => l.length)) * a.fontSize * 0.6
  return {width: Math.max(a.fontSize, w), height: lines.length * a.fontSize * 1.25}
}

export function numberRadius(a: Annotation): number {
  return a.fontSize * 0.75
}

// bounds: the box of an annotation, in the image's pixels.
export function bounds(a: Annotation): Rect {
  switch (a.type) {
    case 'text': {
      const s = textSize(a)
      return {x: a.points[0].x, y: a.points[0].y, ...s}
    }
    case 'number': {
      const r = numberRadius(a)
      return {x: a.points[0].x - r, y: a.points[0].y - r, width: 2 * r, height: 2 * r}
    }
    default: {
      const b = pointsBounds(a.points)
      const pad = a.type === 'pen' || a.type === 'highlighter' || a.type === 'line' || a.type === 'arrow' ? a.width / 2 : 0
      return {x: b.x - pad, y: b.y - pad, width: b.width + 2 * pad, height: b.height + 2 * pad}
    }
  }
}

function inRect(p: Point, r: Rect, tol: number): boolean {
  return p.x >= r.x - tol && p.x <= r.x + r.width + tol && p.y >= r.y - tol && p.y <= r.y + r.height + tol
}

// hit: whether p (the image's pixels) touches a, within tol.
export function hit(a: Annotation, p: Point, tol: number): boolean {
  switch (a.type) {
    case 'pen':
    case 'highlighter':
    case 'line':
    case 'arrow':
      for (let i = 0; i < a.points.length - 1; i++) {
        if (distToSegment(p, a.points[i], a.points[i + 1]) <= tol + a.width / 2) return true
      }
      return a.points.length === 1 && Math.hypot(p.x - a.points[0].x, p.y - a.points[0].y) <= tol + a.width
    case 'rect':
    case 'ellipse':
      if (a.fill) return inRect(p, bounds(a), tol)
      {
        const b = bounds(a)
        const inner = {x: b.x + tol + a.width, y: b.y + tol + a.width, width: b.width - 2 * (tol + a.width), height: b.height - 2 * (tol + a.width)}
        return inRect(p, b, tol + a.width / 2) && !(inner.width > 0 && inner.height > 0 && inRect(p, inner, 0))
      }
    default:
      return inRect(p, bounds(a), tol)
  }
}

// topmost: the last (drawn on top) annotation that p touches.
export function topmost(list: Annotation[], p: Point, tol: number): Annotation | null {
  for (let i = list.length - 1; i >= 0; i--) if (hit(list[i], p, tol)) return list[i]
  return null
}

export function move(a: Annotation, dx: number, dy: number): Annotation {
  return {...a, points: a.points.map((p) => ({x: p.x + dx, y: p.y + dy}))}
}

// resize: the annotation's box from `from` to `to` (its points scaled in
// it); a text or a number grows its font.
export function resize(a: Annotation, from: Rect, to: Rect): Annotation {
  const sx = from.width ? to.width / from.width : 1
  const sy = from.height ? to.height / from.height : 1
  if (a.type === 'text' || a.type === 'number') {
    const k = Math.max(0.1, a.type === 'number' ? Math.min(sx, sy) : sy)
    const fontSize = Math.max(6, Math.round(a.fontSize * k))
    const p = a.type === 'number' ? {x: to.x + to.width / 2, y: to.y + to.height / 2} : {x: to.x, y: to.y}
    return {...a, fontSize, points: [p]}
  }
  return {
    ...a,
    points: a.points.map((p) => ({x: to.x + (p.x - from.x) * sx, y: to.y + (p.y - from.y) * sy})),
  }
}

export function duplicate(a: Annotation, list: Annotation[], offset = 16): Annotation {
  const d = move({...a, id: newId()}, offset, offset)
  if (d.type === 'number') d.n = nextNumber(list)
  return d
}

// --- drawing (a 2D canvas) ---

function arrowHead(ctx: CanvasRenderingContext2D, from: Point, to: Point, w: number) {
  const ang = Math.atan2(to.y - from.y, to.x - from.x)
  const len = Math.max(12, w * 4)
  ctx.beginPath()
  ctx.moveTo(to.x, to.y)
  ctx.lineTo(to.x - len * Math.cos(ang - Math.PI / 7), to.y - len * Math.sin(ang - Math.PI / 7))
  ctx.lineTo(to.x - len * Math.cos(ang + Math.PI / 7), to.y - len * Math.sin(ang + Math.PI / 7))
  ctx.closePath()
  ctx.fillStyle = ctx.strokeStyle
  ctx.fill()
}

function norm(r: Rect): Rect {
  return {
    x: Math.min(r.x, r.x + r.width),
    y: Math.min(r.y, r.y + r.height),
    width: Math.abs(r.width),
    height: Math.abs(r.height),
  }
}

export function cornersRect(a: Annotation): Rect {
  const [p, q] = [a.points[0], a.points[a.points.length - 1]]
  return norm({x: p.x, y: p.y, width: q.x - p.x, height: q.y - p.y})
}

// draw paints a over ctx (in the image's pixels) — base is the image below
// (what pixelate and blur read), k its pixels per pixel of the image (a
// preview made smaller: k < 1).
export function draw(ctx: CanvasRenderingContext2D, a: Annotation, base?: CanvasImageSource, k = 1) {
  ctx.save()
  ctx.globalAlpha = a.opacity
  ctx.strokeStyle = a.stroke
  ctx.fillStyle = a.fill || a.stroke
  ctx.lineWidth = a.width
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  if (a.shadow) {
    ctx.shadowColor = 'rgba(0,0,0,0.5)'
    ctx.shadowBlur = Math.max(4, a.width * 2)
    ctx.shadowOffsetX = ctx.shadowOffsetY = Math.max(2, a.width / 2)
  }
  const [p0] = a.points
  const pl = a.points[a.points.length - 1]
  switch (a.type) {
    case 'pen':
    case 'highlighter':
      ctx.beginPath()
      ctx.moveTo(p0.x, p0.y)
      if (a.points.length === 1) ctx.lineTo(p0.x + 0.01, p0.y)
      for (const p of a.points.slice(1)) ctx.lineTo(p.x, p.y)
      if (a.type === 'highlighter') ctx.globalCompositeOperation = 'multiply'
      ctx.stroke()
      break
    case 'line':
    case 'arrow':
      ctx.beginPath()
      ctx.moveTo(p0.x, p0.y)
      ctx.lineTo(pl.x, pl.y)
      ctx.stroke()
      if (a.type === 'arrow') arrowHead(ctx, p0, pl, a.width)
      break
    case 'rect': {
      const r = cornersRect(a)
      if (a.fill) {
        ctx.fillStyle = a.fill
        ctx.fillRect(r.x, r.y, r.width, r.height)
      }
      ctx.strokeRect(r.x, r.y, r.width, r.height)
      break
    }
    case 'ellipse': {
      const r = cornersRect(a)
      ctx.beginPath()
      ctx.ellipse(r.x + r.width / 2, r.y + r.height / 2, r.width / 2, r.height / 2, 0, 0, Math.PI * 2)
      if (a.fill) {
        ctx.fillStyle = a.fill
        ctx.fill()
      }
      ctx.stroke()
      break
    }
    case 'text': {
      ctx.font = `${a.fontSize}px ${a.font}`
      ctx.textBaseline = 'top'
      ctx.fillStyle = a.stroke
      const lines = (a.text || '').split('\n')
      if (a.fill) {
        const b = bounds(a)
        ctx.save()
        ctx.fillStyle = a.fill
        ctx.fillRect(b.x - 4, b.y - 2, b.width + 8, b.height + 4)
        ctx.restore()
      }
      lines.forEach((l, i) => ctx.fillText(l, p0.x, p0.y + i * a.fontSize * 1.25))
      break
    }
    case 'number': {
      const r = numberRadius(a)
      ctx.beginPath()
      ctx.arc(p0.x, p0.y, r, 0, Math.PI * 2)
      ctx.fillStyle = a.fill || a.stroke
      ctx.fill()
      ctx.shadowColor = 'transparent'
      ctx.fillStyle = '#fff'
      ctx.font = `bold ${Math.round(r * 1.1)}px ${a.font}`
      ctx.textAlign = 'center'
      ctx.textBaseline = 'middle'
      ctx.fillText(String(a.n), p0.x, p0.y + r * 0.05)
      break
    }
    case 'pixelate':
    case 'blur': {
      if (!base) break
      const r = cornersRect(a)
      if (r.width < 1 || r.height < 1) break
      ctx.shadowColor = 'transparent'
      ctx.beginPath()
      ctx.rect(r.x, r.y, r.width, r.height)
      ctx.clip()
      const bw = (base as any).width as number
      const bh = (base as any).height as number
      if (a.type === 'blur') {
        ctx.filter = `blur(${Math.max(4, a.width * 2)}px)`
        ctx.drawImage(base, 0, 0, bw, bh, 0, 0, bw / k, bh / k)
      } else {
        // blocks: the region drawn small, then large without smoothing
        const block = Math.max(6, a.width * 3)
        const w = Math.max(1, Math.round(r.width / block))
        const h = Math.max(1, Math.round(r.height / block))
        const tmp = document.createElement('canvas')
        tmp.width = w
        tmp.height = h
        tmp.getContext('2d')!.drawImage(base, r.x * k, r.y * k, r.width * k, r.height * k, 0, 0, w, h)
        ctx.imageSmoothingEnabled = false
        ctx.drawImage(tmp, 0, 0, w, h, r.x, r.y, r.width, r.height)
      }
      break
    }
  }
  ctx.restore()
}

export function drawAll(ctx: CanvasRenderingContext2D, list: Annotation[], base?: CanvasImageSource, k = 1) {
  for (const a of list) draw(ctx, a, base, k)
}
