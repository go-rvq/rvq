// The orientation of an image: turned clockwise (rotate) and mirrored
// horizontally (flipH) — a vertical mirror is a turn of 180° and a horizontal
// one. Of the editor and of the camera (vuetifyx, by its source).

export type Rotation = 0 | 90 | 180 | 270

export interface Orientation {
  rotate: Rotation
  flipH: boolean
}

// flipped: the turn of an image (rotate, flipH) mirrored as it is seen —
// horizontally (h: left and right swapped) or vertically (v: top and
// bottom) —, whatever it was turned: the mirror is applied after the turn.
export function flipped(rotate: Rotation, flipH: boolean, axis: 'h' | 'v'): Orientation {
  // H·R(r)·F = R(-r)·H·F; V = R(180)·H
  const r = axis === 'h' ? 360 - rotate : 180 - rotate
  return {rotate: (((r % 360) + 360) % 360) as Rotation, flipH: !flipH}
}

export function canvas(w: number, h: number): HTMLCanvasElement {
  const c = document.createElement('canvas')
  c.width = Math.max(1, Math.round(w))
  c.height = Math.max(1, Math.round(h))
  return c
}

// sizeOf: the size of an image, a video's frame, a canvas.
export function sizeOf(s: CanvasImageSource): {width: number; height: number} {
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
