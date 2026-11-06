import {describe, expect, it} from 'vitest'
import {bounds, create, defaultStyle, duplicate, hit, move, nextNumber, resize, topmost} from '../annotations'

describe('the annotations', () => {
  const st = defaultStyle()

  it('a number counts 1, 2, 3…', () => {
    const list = [create('number', {x: 10, y: 10}, st, 1), create('number', {x: 50, y: 10}, st, 2)]
    expect(nextNumber(list)).toBe(3)
    const d = duplicate(list[0], list)
    expect(d.n).toBe(3)
    expect(d.id).not.toBe(list[0].id)
    expect(d.points[0]).toEqual({x: 26, y: 26})
  })

  it('a line is hit along it, not away from it', () => {
    const a = create('line', {x: 0, y: 0}, st)
    a.points = [{x: 0, y: 0}, {x: 100, y: 0}]
    expect(hit(a, {x: 50, y: 3}, 4)).toBe(true)
    expect(hit(a, {x: 50, y: 30}, 4)).toBe(false)
  })

  it('an empty rectangle is hit on its border, not inside', () => {
    const a = create('rect', {x: 0, y: 0}, st)
    a.points = [{x: 0, y: 0}, {x: 100, y: 100}]
    expect(hit(a, {x: 0, y: 50}, 4)).toBe(true)
    expect(hit(a, {x: 50, y: 50}, 4)).toBe(false)
    a.fill = '#fff'
    expect(hit(a, {x: 50, y: 50}, 4)).toBe(true)
  })

  it('the topmost is the last drawn', () => {
    const a = create('rect', {x: 0, y: 0}, {...st, fill: '#000'})
    a.points = [{x: 0, y: 0}, {x: 100, y: 100}]
    const b = move({...a, id: 'b'}, 10, 10)
    expect(topmost([a, b], {x: 50, y: 50}, 2)!.id).toBe('b')
  })

  it('resizes its box', () => {
    const a = create('ellipse', {x: 0, y: 0}, st)
    a.points = [{x: 0, y: 0}, {x: 100, y: 50}]
    const r = resize(a, bounds(a), {x: 10, y: 10, width: 200, height: 100})
    expect(r.points).toEqual([{x: 10, y: 10}, {x: 210, y: 110}])
  })
})
