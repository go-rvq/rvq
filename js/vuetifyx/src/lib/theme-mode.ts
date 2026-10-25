// The theme of the page as the user chose it: light, dark, or auto — the
// system's, followed as it changes (prefers-color-scheme). Kept in the
// browser under a key; vx-theme-toggle and vx-theme-select share it.

export type ThemeMode = 'light' | 'dark' | 'auto'

export const readMode = (key: string): ThemeMode | null => {
  try {
    const v = localStorage.getItem(key)
    return v === 'light' || v === 'dark' || v === 'auto' ? v : null
  } catch {
    return null
  }
}

export const writeMode = (key: string, v: ThemeMode) => {
  try {
    localStorage.setItem(key, v)
  } catch {
    /* private window */
  }
}

const media = () =>
  typeof window !== 'undefined' && window.matchMedia ? window.matchMedia('(prefers-color-scheme: dark)') : null

// systemTheme is the system's theme now.
export const systemTheme = (): 'light' | 'dark' => (media()?.matches ? 'dark' : 'light')

let follow: ((e: MediaQueryListEvent) => void) | null = null

// applyMode sets the theme of mode (set: the theme by its name); auto
// follows the system from now on, any other mode stops following it.
export function applyMode(mode: ThemeMode, set: (name: string) => void) {
  const m = media()
  if (follow && m) {
    m.removeEventListener('change', follow)
    follow = null
  }
  if (mode === 'auto') {
    set(systemTheme())
    if (m) {
      follow = (e) => set(e.matches ? 'dark' : 'light')
      m.addEventListener('change', follow)
    }
    return
  }
  set(mode)
}
