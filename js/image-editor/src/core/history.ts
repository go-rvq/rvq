// The undo / redo of the editor's dialog: one history of states for all of it
// (the crop, the colors, the annotations). A state is a plain value, cloned
// on the way in and out, so a step never shares objects with the next.
//
// A drag is one step: begin() takes the state as it was, the moves change it
// freely (set, not recorded), and end() records it — on release — when it
// changed.

export const MAX_STEPS = 100

export function clone<T>(v: T): T {
  return v === undefined ? v : JSON.parse(JSON.stringify(v))
}

function same(a: unknown, b: unknown): boolean {
  return JSON.stringify(a) === JSON.stringify(b)
}

export class History<T> {
  private past: T[] = []
  private future: T[] = []
  private current: T
  // the state at begin(), while a drag goes on
  private pending: T | null = null

  constructor(initial: T, private max = MAX_STEPS) {
    this.current = clone(initial)
  }

  // the state now (a copy)
  get state(): T {
    return clone(this.current)
  }

  get undoCount(): number {
    return this.past.length
  }

  get redoCount(): number {
    return this.future.length
  }

  get canUndo(): boolean {
    return this.past.length > 0
  }

  get canRedo(): boolean {
    return this.future.length > 0
  }

  // push records next as a step (nothing when it is the same state); the
  // steps undone are gone.
  push(next: T): boolean {
    if (this.pending !== null) {
      this.current = clone(next)
      return false
    }
    if (same(next, this.current)) return false
    this.past.push(this.current)
    if (this.past.length > this.max) this.past.shift()
    this.current = clone(next)
    this.future = []
    return true
  }

  // begin starts a drag: the moves (set) are one step at end().
  begin(): void {
    if (this.pending === null) this.pending = clone(this.current)
  }

  // set changes the state with no step: the moves of a drag.
  set(next: T): void {
    this.current = clone(next)
  }

  // end records the drag as one step, when it changed anything.
  end(): boolean {
    if (this.pending === null) return false
    const before = this.pending
    this.pending = null
    if (same(before, this.current)) return false
    this.past.push(before)
    if (this.past.length > this.max) this.past.shift()
    this.future = []
    return true
  }

  get dragging(): boolean {
    return this.pending !== null
  }

  undo(): T | null {
    if (this.pending !== null) this.end()
    const prev = this.past.pop()
    if (prev === undefined) return null
    this.future.push(this.current)
    this.current = prev
    return this.state
  }

  redo(): T | null {
    const next = this.future.pop()
    if (next === undefined) return null
    this.past.push(this.current)
    this.current = next
    return this.state
  }

  // reset forgets the steps: initial is the state of a new history.
  reset(initial: T): void {
    this.past = []
    this.future = []
    this.pending = null
    this.current = clone(initial)
  }
}
