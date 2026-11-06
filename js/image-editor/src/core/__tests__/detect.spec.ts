import {describe, expect, it} from 'vitest'
import {detectCorners} from '../detect'
import {pointInPolygon} from '../geometry'

describe('"Detectar bordas"', () => {
  it('finds the corners of a light document on a dark background', () => {
    const w = 200, h = 150
    const doc = [
      {x: 40, y: 30},
      {x: 170, y: 20},
      {x: 180, y: 120},
      {x: 30, y: 130},
    ]
    const data = new Uint8ClampedArray(w * h * 4)
    for (let y = 0; y < h; y++)
      for (let x = 0; x < w; x++) {
        const o = (y * w + x) * 4
        const v = pointInPolygon({x: x + 0.5, y: y + 0.5}, doc) ? 230 : 40
        data[o] = data[o + 1] = data[o + 2] = v
        data[o + 3] = 255
      }
    const q = detectCorners({width: w, height: h, data})!
    expect(q).not.toBeNull()
    q.forEach((p, i) => {
      expect(Math.hypot(p.x - doc[i].x, p.y - doc[i].y)).toBeLessThan(4)
    })
  })

  it('nothing on a flat image', () => {
    const data = new Uint8ClampedArray(50 * 50 * 4).fill(128)
    expect(detectCorners({width: 50, height: 50, data})).toBeNull()
  })
})
