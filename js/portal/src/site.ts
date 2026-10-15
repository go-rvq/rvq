// The runtime of the portals of a site (the bundle rvq-portal.umd.js): each
// <rvq-portal source="TOKEN"> of the page — the page names no address, only
// the opaque token — is loaded when it comes into view: a progress while the
// server answers, then the content (a Vue template and its data) set as the
// value of its v-model, which the portal compiles and shows.
//
// The content and what it does (an action: send a form…) are asked at one
// fixed path (path), the token and the action in headers (the Go side:
// js/api/portal, Server). Vue — with its compiler: the templates are compiled
// here — is the page's global (window.Vue).
import { type App, createApp, h, markRaw, reactive } from 'vue'

import Portal from './Portal'
import type { ContentType } from './types'

// path is where the server of the portals answers: the Go server writes its
// own (Server.Path) here when it serves this bundle.
export const path = '/_rvq/portal'
export const headerToken = 'X-Rvq-Portal'
export const headerAction = 'X-Rvq-Portal-Action'

// how far before the viewport a portal starts loading
const rootMargin = '200px 0px'

// the words of the runtime, the page's when its element says them
// (data-loading-text, data-error-text, data-expired-text, data-retry-text)
const words = {
  loading: 'Loading…',
  error: 'This content could not be loaded.',
  expired: 'This page is out of date.',
  retry: 'Try again'
}

const styleId = 'rvq-portal-style'
const style = `
.rvq-portal-progress { position: relative; height: 4px; overflow: hidden; border-radius: 2px;
  background: rgba(0, 0, 0, .08); margin: 1rem 0; }
.rvq-portal-progress::after { content: ''; position: absolute; inset: 0; width: 40%;
  background: currentColor; opacity: .6; animation: rvq-portal-progress 1.1s ease-in-out infinite; }
@keyframes rvq-portal-progress { 0% { transform: translateX(-100%) } 100% { transform: translateX(250%) } }
.rvq-portal-error { padding: 1rem; text-align: center; }
.rvq-portal-error button { margin-left: .5rem; text-decoration: underline; cursor: pointer;
  background: none; border: 0; color: inherit; font: inherit; }
`

function injectStyle(doc: Document) {
  if (doc.getElementById(styleId)) return
  const el = doc.createElement('style')
  el.id = styleId
  el.textContent = style
  doc.head.appendChild(el)
}

// Client is how a content talks to the server: its action (name) with a
// body — a FormData (files), a string, or an object as JSON —; the answer is
// the fetch's.
export interface Client {
  action(name: string, body?: BodyInit | Record<string, any> | null, init?: RequestInit): Promise<Response>
}

export function client(token: string, base = path): Client {
  return {
    action(name, body, init = {}) {
      const headers = new Headers(init.headers)
      headers.set(headerToken, token)
      headers.set(headerAction, name)
      let b: BodyInit | null | undefined = body as any
      if (body && typeof body === 'object' && !(body instanceof FormData) && !(body instanceof Blob)
        && !(body instanceof URLSearchParams) && !(body instanceof ArrayBuffer)) {
        b = JSON.stringify(body)
        headers.set('Content-Type', 'application/json')
      }
      return fetch(base, { ...init, method: 'POST', credentials: 'same-origin', headers, body: b })
    }
  }
}

// Loaded is what the server answers for a portal: its template and data;
// the scripts it needs (loaded once for the page, before it is shown); the
// kind of content that makes it (kinds), when not the template as it is.
export interface Loaded {
  template: string
  data?: any
  scripts?: string[]
  kind?: string
}

// Kind makes the content of a portal of its kind from what the server
// answered and the portal's client: its template (as it is, or made anew),
// its data, its behavior (options.mixins) — registered by the scripts of the
// content (RvqPortal.kinds[name] = …).
export type Kind = (loaded: Loaded, client: Client | null) => { template: string; options?: Record<string, any> }

export const kinds: Record<string, Kind> = {}

// scripts loaded for the page, by their address: once each
const scripts: Record<string, Promise<void>> = {}

export function loadScript(src: string, doc: Document = document): Promise<void> {
  return (scripts[src] ||= new Promise<void>((resolve, reject) => {
    const el = doc.createElement('script')
    el.src = src
    el.async = false
    el.onload = () => resolve()
    el.onerror = () => {
      delete scripts[src]
      reject(new Error('script ' + src))
    }
    doc.head.appendChild(el)
  }))
}

// contentOf is the value of the v-model of a portal of what the server
// answered: the template and its data (reactive), the client of the content
// (portal.action) — or, of a kind, what the kind makes.
export function contentOf(c: Loaded, api: Client | null): ContentType {
  if (c.kind) {
    const kind = kinds[c.kind]
    if (!kind) throw new Error('kind ' + c.kind)
    const made = kind(c, api)
    const options = made.options || {}
    return { template: made.template, options: { ...options, data: reactive({ ...(options.data || {}), portal: api && markRaw(api) }) } }
  }
  return { template: c.template, options: { data: reactive({ ...(c.data || {}), portal: api && markRaw(api) }) } }
}

// load asks the content of token: {template, data, scripts, kind}, or an
// error ('expired' for a page older than its token).
export async function load(token: string, base = path): Promise<Loaded> {
  const res = await fetch(base, {
    method: 'POST',
    credentials: 'same-origin',
    headers: { [headerToken]: token, Accept: 'application/json' }
  })
  if (res.status === 410) throw new Error('expired')
  if (!res.ok) throw new Error(String(res.status))
  return res.json()
}

// mount makes the element el — <rvq-portal source="…"> of the page — a
// portal: it waits until it is in view, shows a progress, loads its content
// and shows it; an error says so, with a button to try again (or, for a
// page out of date, to reload it).
export function mount(el: Element, opts: { base?: string; observe?: boolean } = {}): App | null {
  const token = el.getAttribute('source')
  if (!token || (el as any).__rvqPortal) return null
  const doc = el.ownerDocument
  injectStyle(doc)
  const t = {
    loading: el.getAttribute('data-loading-text') || words.loading,
    error: el.getAttribute('data-error-text') || words.error,
    expired: el.getAttribute('data-expired-text') || words.expired,
    retry: el.getAttribute('data-retry-text') || words.retry
  }
  const base = opts.base || path
  const state = reactive({
    content: '' as ContentType,
    status: 'waiting' as 'waiting' | 'loading' | 'shown' | 'error' | 'expired'
  })
  const api = client(token, base)

  const fetchContent = async () => {
    state.status = 'loading'
    try {
      const c = await load(token, base)
      for (const src of c.scripts || []) await loadScript(src, doc)
      // the value of the v-model: the template, its data reactive, and the
      // client of the content (portal.action)
      state.content = contentOf(c, api)
      state.status = 'shown'
    } catch (e: any) {
      state.status = e?.message === 'expired' ? 'expired' : 'error'
    }
  }

  // the element of the page becomes the place of the app
  const host = doc.createElement('div')
  host.className = el.className
  host.setAttribute('data-rvq-portal', '')
  el.replaceWith(host)

  const app = createApp({
    setup() {
      return () => {
        switch (state.status) {
          case 'waiting':
          case 'loading':
            return h('div', { class: 'rvq-portal-progress', role: 'progressbar', 'aria-label': t.loading,
              'aria-busy': 'true' })
          case 'error':
          case 'expired': {
            const expired = state.status === 'expired'
            return h('div', { class: 'rvq-portal-error', role: 'alert' }, [
              expired ? t.expired : t.error,
              h('button', { type: 'button', onClick: () => (expired ? location.reload() : fetchContent()) }, t.retry)
            ])
          }
          default:
            return h(Portal, {
              tag: 'div',
              modelValue: state.content,
              'onUpdate:modelValue': (v: ContentType) => (state.content = v)
            })
        }
      }
    }
  })
  app.mount(host)
  ;(host as any).__rvqPortal = app

  const start = () => {
    if (state.status === 'waiting') fetchContent()
  }
  const IO = (doc.defaultView as any)?.IntersectionObserver
  if (opts.observe === false || !IO) {
    start()
  } else {
    const io = new IO((entries: IntersectionObserverEntry[]) => {
      if (entries.some((e) => e.isIntersecting)) {
        io.disconnect()
        start()
      }
    }, { rootMargin })
    io.observe(host)
  }
  return app
}

// mountContent shows a content given — not asked to the server: a preview,
// with no token — in el, as a portal does: its kind's (loaded), else the
// template as it is; it acts nowhere (no client).
export function mountContent(el: Element, c: Loaded): App {
  const content = contentOf(c, null)
  const app = createApp({ setup: () => () => h(Portal, { tag: 'div', modelValue: content }) })
  app.mount(el)
  return app
}

// mountAll mounts the portals of root — the page when it loads.
export function mountAll(root: ParentNode = document, opts: { base?: string } = {}): App[] {
  const apps: App[] = []
  root.querySelectorAll('rvq-portal[source]').forEach((el) => {
    const app = mount(el, opts)
    if (app) apps.push(app)
  })
  return apps
}

export { Portal }

if (typeof document !== 'undefined' && !(globalThis as any).__rvqPortalNoAuto) {
  if (document.readyState === 'loading') document.addEventListener('DOMContentLoaded', () => mountAll())
  else mountAll()
}
