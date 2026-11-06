import {describe, expect, it} from 'vitest'
import {apply, outputSize, rectToQuad, warp, type Quad} from '../perspective'

describe('the perspective', () => {
  const q: Quad = [
    {x: 10, y: 20},
    {x: 210, y: 0},
    {x: 230, y: 160},
    {x: 0, y: 140},
  ]

  it('the size: the mean of the opposite sides', () => {
    const s = outputSize(q)
    expect(s.width).toBe(Math.round((Math.hypot(200, 20) + Math.hypot(230, 20)) / 2))
  })

  it('the rectangle\'s corners go to the quadrilateral\'s', () => {
    const m = rectToQuad(q, 100, 50)
    const c = [apply(m, {x: 0, y: 0}), apply(m, {x: 100, y: 0}), apply(m, {x: 100, y: 50}), apply(m, {x: 0, y: 50})]
    c.forEach((p, i) => {
      expect(p.x).toBeCloseTo(q[i].x)
      expect(p.y).toBeCloseTo(q[i].y)
    })
  })

  it('rectifies: a tilted bright square comes out filling the rectangle', () => {
    const w = 64, h = 64
    const data = new Uint8ClampedArray(w * h * 4)
    // a white quadrilateral on black
    const quad: Quad = [
      {x: 12, y: 10},
      {x: 52, y: 14},
      {x: 50, y: 54},
      {x: 10, y: 50},
    ]
    for (let y = 0; y < h; y++)
      for (let x = 0; x < w; x++) {
        const inside = x >= 13 && x <= 49 && y >= 15 && y <= 49
        const o = (y * w + x) * 4
        data[o] = data[o + 1] = data[o + 2] = inside ? 255 : 0
        data[o + 3] = 255
      }
    const out = warp({width: w, height: h, data}, quad, {width: 20, height: 20})
    // the middle of the result is white
    const mid = (10 * 20 + 10) * 4
    expect(out.data[mid]).toBe(255)
    expect(out.width).toBe(20)
  })
})
