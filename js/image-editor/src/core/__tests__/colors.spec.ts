import {describe, expect, it} from 'vitest'
import {analyze, applyColors, fixLighting, grayWorld, neutral} from '../colors'

// a synthetic photo: a gradient between lo and hi, tinted by (r, g, b)
function photo(lo: number, hi: number, tint: [number, number, number] = [1, 1, 1]) {
  const w = 64, h = 64
  const data = new Uint8ClampedArray(w * h * 4)
  for (let y = 0; y < h; y++)
    for (let x = 0; x < w; x++) {
      const v = lo + ((hi - lo) * (x + y)) / (w + h - 2)
      const o = (y * w + x) * 4
      data[o] = Math.min(255, v * tint[0])
      data[o + 1] = Math.min(255, v * tint[1])
      data[o + 2] = Math.min(255, v * tint[2])
      data[o + 3] = 255
    }
  return {width: w, height: h, data}
}

describe('"Corrigir iluminação"', () => {
  it('a dark photo: its mean to the middle, nothing blown', () => {
    const dark = photo(5, 60)
    const fixed = applyColors(dark, fixLighting(analyze(dark)))
    const s = analyze(fixed)
    expect(s.mean).toBeGreaterThan(100)
    expect(s.mean).toBeLessThan(156)
    expect(s.clippedHigh).toBeLessThan(0.03)
  })

  it('a bright photo: its mean to the middle, nothing blown', () => {
    const bright = photo(190, 254)
    const fixed = applyColors(bright, fixLighting(analyze(bright)))
    const s = analyze(fixed)
    expect(s.mean).toBeGreaterThan(100)
    expect(s.mean).toBeLessThan(156)
    expect(s.clippedHigh).toBeLessThan(0.03)
  })

  it('sets the controls (it does not paint)', () => {
    const c = fixLighting(analyze(photo(5, 60)))
    expect(c.black).toBeLessThanOrEqual(10)
    expect(c.white).toBeGreaterThanOrEqual(55)
    // a dark photo: its shadows opened
    expect(c.shadows).toBeGreaterThan(0)
  })

  it('the gray world takes out a yellow cast', () => {
    const yellow = photo(40, 200, [1.15, 1.05, 0.7])
    const gains = grayWorld(analyze(yellow).avg)
    const s = analyze(applyColors(yellow, {...neutral(), wb: gains}))
    expect(Math.abs(s.avg[0] - s.avg[2])).toBeLessThan(8)
  })

  it('the neutral controls change nothing', () => {
    const p = photo(10, 200)
    expect(Array.from(applyColors(p, neutral()).data)).toEqual(Array.from(p.data))
  })
})

describe('"Corrigir iluminação" on a document', () => {
  it('a light sheet with dark text keeps its paper light (no gray)', () => {
    const w = 100, h = 100
    const data = new Uint8ClampedArray(w * h * 4)
    for (let i = 0; i < w * h; i++) {
      const text = i % 50 === 0
      const v = text ? 40 : 210
      data.set([v, v, v * 0.85, 255], i * 4)
    }
    const px = {width: w, height: h, data}
    const c = fixLighting(analyze(px))
    expect(c.gamma).toBeGreaterThanOrEqual(0.6)
    const s = analyze(applyColors(px, c))
    expect(s.mean).toBeGreaterThan(180)
  })
})
