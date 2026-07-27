// The `closer` — the control object of every overlay (dialog, drawer, form
// host). It is created in three places (the app root, `go-plaid-scope`'s
// `:closer` prop and a `user-component` scope variable, see presets.FormHost)
// and ADOPTED in several more (a portal seeds `closer` in the content scope, a
// dialog binds `v-model="closer.show"`), so its shape has to come from ONE
// factory — this file.
//
// `show` is an accessor over `_show`: turning it on fires `openCallbacks`,
// turning it off fires `closeCallbacks`, and assigning the value it already has
// fires nothing. The trigger lives on the OBJECT (not in a component `watch`)
// precisely because the same closer is adopted by several components — a
// watcher per component would fire the callbacks once per adopter.

import { reactive } from 'vue'
import { installCloserUrlSync } from './closer-url'

export type CloserCallback = (closer: any) => void

/** A closer's URL: either the address itself or a function of the closer. */
export type CloserUrl = string | ((closer: any) => string)

export interface Closer {
  show: boolean
  openCallbacks: CloserCallback[]
  closeCallbacks: CloserCallback[]
  /** Registers an open callback and returns the function that removes it. */
  onOpen(fn: CloserCallback): () => void
  onClose(fn: CloserCallback): () => void
  /** Optional clean page URL this overlay stands for, see closer-url.ts. */
  url?: CloserUrl
  $parent?: any

  [k: string]: any
}

export function isCloser(v: any): boolean {
  return (
    !!v &&
    typeof v === 'object' &&
    Array.isArray(v.openCallbacks) &&
    Array.isArray(v.closeCallbacks)
  )
}

function fire(callbacks: CloserCallback[], closer: any) {
  // a copy: a callback may add or remove callbacks (onOpen returns a remover)
  for (const fn of [...callbacks]) {
    try {
      fn(closer)
    } catch (e) {
      console.error('closer callback failed:', e)
    }
  }
}

/**
 * Normalizes an object INTO a closer, in place, and returns it. Idempotent: an
 * object that already is one is returned untouched, which is what makes
 * adoption safe (a scope receiving an existing closer must not re-wrap it).
 *
 * Must run on the raw object, before `reactive()`.
 */
export function asCloser(target: any): any {
  if (!target || typeof target !== 'object') {
    return target
  }
  if (isCloser(target)) {
    return target
  }

  // whatever `show` the initializer carried becomes the initial state, WITHOUT
  // firing anything: there is nobody subscribed at creation time.
  const initial = !!target.show
  delete target.show
  target._show = initial

  if (!Array.isArray(target.openCallbacks)) {
    target.openCallbacks = []
  }
  if (!Array.isArray(target.closeCallbacks)) {
    target.closeCallbacks = []
  }

  Object.defineProperty(target, 'show', {
    enumerable: true,
    configurable: true,
    get(): boolean {
      return this._show
    },
    set(v: boolean) {
      const next = !!v
      if (next === this._show) {
        return
      }
      this._show = next
      fire(next ? this.openCallbacks : this.closeCallbacks, this)
    }
  })

  const register = (list: string) => (fn: CloserCallback) => {
    target[list].push(fn)
    return () => {
      const i = target[list].indexOf(fn)
      if (i >= 0) {
        target[list].splice(i, 1)
      }
    }
  }

  Object.defineProperty(target, 'onOpen', {
    enumerable: false,
    configurable: true,
    value: register('openCallbacks')
  })
  Object.defineProperty(target, 'onClose', {
    enumerable: false,
    configurable: true,
    value: register('closeCallbacks')
  })

  return target
}

/**
 * Creates a reactive closer from an initializer (`{show, url, …}` plus whatever
 * state its owner keeps on it — see presets.FormHost, which stores the record
 * id and the reload/onSave hooks there).
 */
export function createCloser(init?: any): Closer {
  const raw = asCloser({ ...(init || {}) })
  // registered for every closer, not only the ones that carry a url today: the
  // owner may set `closer.url` later, and the hooks are a no-op without one.
  installCloserUrlSync(raw)
  return reactive(raw) as Closer
}
