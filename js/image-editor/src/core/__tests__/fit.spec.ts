import {describe, expect, it} from 'vitest'
import {fitSize, limitFor, scaleRect} from '../fit'

const fullHD = {width: 1920, height: 1080}

describe('fit (the limit of an upload)', () => {
  it('a smaller image untouched', () => {
    expect(fitSize(1280, 720, fullHD)).toEqual({width: 1280, height: 720, ratio: 1})
    expect(fitSize(1920, 1080, fullHD).ratio).toBe(1)
  })
  it('a larger one made smaller, its proportions kept', () => {
    expect(fitSize(4000, 3000, fullHD)).toMatchObject({width: 1440, height: 1080})
    expect(fitSize(4000, 1200, fullHD)).toMatchObject({width: 1920, height: 576})
  })
  it('a portrait one fits the limit on its side', () => {
    expect(limitFor(3000, 4000, fullHD)).toEqual({width: 1080, height: 1920})
    expect(fitSize(3000, 4000, fullHD)).toMatchObject({width: 1080, height: 1440})
    // a square image, the limit as it is
    expect(fitSize(3000, 3000, fullHD)).toMatchObject({width: 1080, height: 1080})
  })
  it('no limit, or a side of 0: none on it', () => {
    expect(fitSize(5000, 5000, null).ratio).toBe(1)
    expect(fitSize(5000, 100, {width: 2000, height: 0})).toMatchObject({width: 2000, height: 40})
  })
  it('a crop at the new scale', () => {
    expect(scaleRect({x: 100, y: 50, width: 1000, height: 500}, 0.5)).toEqual({x: 50, y: 25, width: 500, height: 250})
  })
})
