import {describe, expect, it} from 'vitest'
import {History, MAX_STEPS} from '../history'

describe('History', () => {
  it('undoes and redoes the steps, in order', () => {
    const h = new History({n: 0})
    h.push({n: 1})
    h.push({n: 2})
    expect(h.undoCount).toBe(2)
    expect(h.undo()).toEqual({n: 1})
    expect(h.undo()).toEqual({n: 0})
    expect(h.undo()).toBeNull()
    expect(h.redo()).toEqual({n: 1})
    expect(h.redoCount).toBe(1)
  })

  it('a new step forgets the steps undone', () => {
    const h = new History({n: 0})
    h.push({n: 1})
    h.undo()
    h.push({n: 5})
    expect(h.canRedo).toBe(false)
    expect(h.state).toEqual({n: 5})
  })

  it('the same state is no step', () => {
    const h = new History({n: 0})
    expect(h.push({n: 0})).toBe(false)
    expect(h.undoCount).toBe(0)
  })

  it('a drag is one step, at its end', () => {
    const h = new History({p: [0, 0]})
    h.begin()
    for (let i = 1; i <= 10; i++) h.set({p: [i, i]})
    expect(h.undoCount).toBe(0)
    expect(h.end()).toBe(true)
    expect(h.undoCount).toBe(1)
    expect(h.undo()).toEqual({p: [0, 0]})
    // a drag that ends where it began: nothing
    h.begin()
    h.set({p: [3, 3]})
    h.set({p: [0, 0]})
    expect(h.end()).toBe(false)
  })

  it('keeps at most MAX_STEPS', () => {
    const h = new History(0)
    for (let i = 1; i <= MAX_STEPS + 20; i++) h.push(i)
    expect(h.undoCount).toBe(MAX_STEPS)
  })

  it('never shares the objects of a state', () => {
    const s = {a: [1]}
    const h = new History(s)
    s.a.push(2)
    expect(h.state).toEqual({a: [1]})
  })
})
