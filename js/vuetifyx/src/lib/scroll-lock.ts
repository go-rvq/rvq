// The page's scroll, locked while a dialog or a floating (temporary) drawer is
// open: each one holds a lock of its own, a class on <html>
// (overflow-hidden, overflow-y-hidden), and the class goes away when no lock
// asks for it any more — whatever the order they close in, and when one is
// unmounted while open. A class <html> already had before the first lock is
// never taken away.

const locks = new Map<object, string>()
const kept = new Set<string>()

function html(): DOMTokenList | undefined {
  return typeof document === 'undefined' ? undefined : document.documentElement?.classList
}

// lockScroll is owner's lock of the page's scroll, by the class cls; a lock it
// already holds is changed to cls.
export function lockScroll(owner: object, cls: string) {
  const cl = html()
  if (!cl) return
  const had = locks.get(owner)
  if (had === cls) return
  if (![...locks.values()].includes(cls) && cl.contains(cls)) {
    kept.add(cls)
  }
  locks.set(owner, cls)
  cl.add(cls)
  if (had) release(cl, had)
}

// unlockScroll gives owner's lock back: its class leaves <html> when no other
// lock holds it.
export function unlockScroll(owner: object) {
  const cl = html()
  const cls = locks.get(owner)
  if (!cl || cls === undefined) return
  locks.delete(owner)
  release(cl, cls)
}

function release(cl: DOMTokenList, cls: string) {
  if ([...locks.values()].includes(cls)) return
  if (kept.has(cls)) {
    kept.delete(cls)
    return
  }
  cl.remove(cls)
}

// scrollLocks is how many locks are held (the tests').
export function scrollLocks(): number {
  return locks.size
}
