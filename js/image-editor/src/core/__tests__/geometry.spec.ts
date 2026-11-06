import {describe, expect, it} from 'vitest'
import {alignEdge, edgeAngle, insertPoint, normDeg, removePoint, resizeRect, rotatePoints, snapRotation, stepRotation, usefulArea} from '../geometry'

const square = [
  {x: 0, y: 0},
  {x: 100, y: 0},
  {x: 100, y: 100},
  {x: 0, y: 100},
]

describe('the free crop', () => {
  it('a click inserts the point on the nearest edge (the path closed, not crossed)', () => {
    const {points, index} = insertPoint(square, {x: 50, y: 103})
    expect(index).toBe(3)
    expect(points[3]).toEqual({x: 50, y: 103})
    expect(points).toHaveLength(5)
  })

  it('removes a point, never below 3', () => {
    const tri = removePoint(square, 0)!
    expect(tri).toHaveLength(3)
    expect(removePoint(tri, 0)).toBeNull()
  })

  it('the magnet: an edge within ~3° of horizontal sticks', () => {
    // the square turned 2°: its edges near 0° and 90°
    expect(snapRotation(square, 2)).toBeCloseTo(0)
    expect(snapRotation(square, -2.5)).toBeCloseTo(0)
    // 10°: far from all of them, kept
    expect(snapRotation(square, 10)).toBeCloseTo(10)
  })

  it('Shift: 15° by 15°', () => {
    expect(stepRotation(22)).toBe(15)
    expect(stepRotation(23)).toBe(30)
  })

  it('"Alinhar borda" leaves the edge horizontal, the side nearest the angle now', () => {
    // a crooked check: its bottom edge (2 → 3) at ~7°
    const check = [
      {x: 100, y: 100},
      {x: 500, y: 150},
      {x: 480, y: 310},
      {x: 80, y: 260},
    ]
    const r = alignEdge(check, 2, 0)
    const turned = rotatePoints(check, r)
    expect(Math.abs(normDeg(edgeAngle(turned[2], turned[3]) - 180)) % 180).toBeCloseTo(0, 5)
    expect(Math.abs(r)).toBeLessThan(90)
  })

  it('the useful area: the bounds of the turned points', () => {
    const pts = [
      {x: 180, y: 48},
      {x: 660, y: 100},
      {x: 600, y: 262},
      {x: 200, y: 200},
    ]
    expect(usefulArea(pts, 0)).toEqual({x: 180, y: 48, width: 480, height: 214})
    // turned 90°: the sides swap
    const a = usefulArea(pts, 90)
    expect(a.width).toBe(214)
    expect(a.height).toBe(480)
  })
})

describe('the fixed crop', () => {
  it('a corner keeps the ratio, the opposite corner fixed', () => {
    const r = resizeRect({x: 0, y: 0, width: 200, height: 100}, 'se', {x: 300, y: 300}, 2)
    expect(r.x).toBe(0)
    expect(r.y).toBe(0)
    expect(r.width / r.height).toBeCloseTo(2)
  })
})

describe('the free rectangle (turned)', () => {
  it('its corners turned around its center', async () => {
    const {turnedCorners} = await import('../geometry')
    const c = turnedCorners({cx: 50, cy: 50, width: 40, height: 20, angle: 90})
    expect(c[0].x).toBeCloseTo(60)
    expect(c[0].y).toBeCloseTo(30)
  })

  it('a side dragged: the opposite side stays, also turned', async () => {
    const {resizeTurned, turnedCorners} = await import('../geometry')
    const r = {cx: 50, cy: 50, width: 40, height: 20, angle: 30}
    const before = turnedCorners(r)
    // the right side (e) to 10 px further along the rectangle's x
    const rad = (30 * Math.PI) / 180
    const e = {x: 50 + 30 * Math.cos(rad), y: 50 + 30 * Math.sin(rad)}
    const n = resizeTurned(r, 'e', e)
    expect(n.width).toBeCloseTo(50)
    const after = turnedCorners(n)
    // the left side (corners 0 and 3) where it was
    expect(after[0].x).toBeCloseTo(before[0].x)
    expect(after[3].y).toBeCloseTo(before[3].y)
  })

  it('the magnet at 0°, 90°…; free between', async () => {
    const {snapAngle} = await import('../geometry')
    expect(snapAngle(2)).toBeCloseTo(0)
    expect(snapAngle(88.5)).toBeCloseTo(90)
    expect(snapAngle(10)).toBeCloseTo(10)
  })
})
