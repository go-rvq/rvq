import {describe, expect, it} from 'vitest'
import {clampScale, fit, imageToScreen, nextStep, onScreen, screenToImage, zoomAt} from '../view'

describe('view (zoom)', () => {
  it('screen ↔ image', () => {
    const v = {scale: 2, x: 10, y: 20}
    expect(screenToImage(v, {x: 30, y: 60})).toEqual({x: 10, y: 20})
    expect(imageToScreen(v, {x: 10, y: 20})).toEqual({x: 30, y: 60})
  })

  it('zooms around the cursor: the point under it stays', () => {
    const v = {scale: 1, x: 0, y: 0}
    const at = {x: 120, y: 80}
    const before = screenToImage(v, at)
    for (const s of [0.25, 1, 8]) {
      const z = zoomAt(v, s, at)
      const after = screenToImage(z, at)
      expect(after.x).toBeCloseTo(before.x)
      expect(after.y).toBeCloseTo(before.y)
    }
  })

  it('10% to 1600%', () => {
    expect(clampScale(0.01)).toBe(0.1)
    expect(clampScale(100)).toBe(16)
    expect(nextStep(1, 1)).toBe(1.25)
    expect(nextStep(1, -1)).toBe(0.75)
  })

  it('the marks the same size on the screen at 25%, 100% and 800%', () => {
    for (const s of [0.25, 1, 8]) {
      const v = {scale: s, x: 0, y: 0}
      // a handle of 6 px: in the image, 6/s; back on the screen, 6
      expect(onScreen(v, 6) * s).toBeCloseTo(6)
    }
  })

  it('the points kept the same in any zoom', () => {
    const p = {x: 123.5, y: 77}
    for (const s of [0.25, 1, 8]) {
      const v = zoomAt({scale: 1, x: 5, y: 7}, s, {x: 40, y: 40})
      const back = screenToImage(v, imageToScreen(v, p))
      expect(back.x).toBeCloseTo(p.x)
      expect(back.y).toBeCloseTo(p.y)
    }
  })

  it('fits the image in the stage, centered', () => {
    const v = fit(800, 600, 1600, 600, 0)
    expect(v.scale).toBeCloseTo(0.5)
    expect(v.x).toBeCloseTo(0)
    expect(v.y).toBeCloseTo(150)
  })
})
