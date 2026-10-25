import { afterEach, beforeAll, describe, expect, it, vi } from 'vitest'
import { applyMode, readMode, writeMode } from '../theme-mode'

// a matchMedia of the test's: the system's theme, changed at will
let dark = false
const listeners: ((e: any) => void)[] = []
beforeAll(() => {
  if (!window.localStorage) {
    const data = new Map<string, string>()
    Object.defineProperty(window, 'localStorage', {
      configurable: true,
      value: {
        getItem: (k: string) => (data.has(k) ? data.get(k)! : null),
        setItem: (k: string, v: string) => void data.set(k, String(v)),
        removeItem: (k: string) => void data.delete(k),
        clear: () => data.clear()
      }
    })
  }
  ;(window as any).matchMedia = () => ({
    get matches() {
      return dark
    },
    addEventListener: (_: string, f: any) => listeners.push(f),
    removeEventListener: (_: string, f: any) => {
      const i = listeners.indexOf(f)
      if (i >= 0) listeners.splice(i, 1)
    }
  })
})

const systemChanges = (d: boolean) => {
  dark = d
  listeners.slice().forEach((f) => f({ matches: d }))
}

describe('the theme mode', () => {
  afterEach(() => {
    window.localStorage.clear()
    applyMode('light', () => {})
  })

  it('is kept: light, dark, auto; anything else is none', () => {
    writeMode('k', 'auto')
    expect(readMode('k')).toBe('auto')
    window.localStorage.setItem('k', 'purple')
    expect(readMode('k')).toBe(null)
  })

  it('auto: the system now, and as it changes; another mode stops following it', () => {
    const set = vi.fn()
    dark = true
    applyMode('auto', set)
    expect(set).toHaveBeenLastCalledWith('dark')
    systemChanges(false)
    expect(set).toHaveBeenLastCalledWith('light')
    applyMode('dark', set)
    expect(set).toHaveBeenLastCalledWith('dark')
    systemChanges(false)
    expect(set).toHaveBeenLastCalledWith('dark')
    expect(listeners.length).toBe(0)
  })
})
