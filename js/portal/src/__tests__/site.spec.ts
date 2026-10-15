import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises } from '@vue/test-utils'
import { nextTick } from 'vue'

;(globalThis as any).__rvqPortalNoAuto = true
const site = await import('../site')

// an IntersectionObserver of the test: entering is said by the test (enter)
let observers: { cb: (e: any[]) => void; els: Element[] }[] = []
class FakeIO {
  cb: (e: any[]) => void
  els: Element[] = []
  constructor(cb: (e: any[]) => void) {
    this.cb = cb
    observers.push(this)
  }
  observe(el: Element) { this.els.push(el) }
  disconnect() { observers = observers.filter((o) => o !== this) }
}
const enter = () => [...observers].forEach((o) => o.cb(o.els.map((target) => ({ isIntersecting: true, target }))))

// a fetch of the test: what each token answers, and the requests seen
let answers: Record<string, () => Response>
let seen: { url: string; headers: Headers; body: any }[]
const fakeFetch = vi.fn(async (url: string, init: RequestInit) => {
  const headers = new Headers(init.headers)
  seen.push({ url, headers, body: init.body })
  const a = answers[headers.get(site.headerToken) + '|' + (headers.get(site.headerAction) || '')]
  return a ? a() : new Response(JSON.stringify({ error: 'token' }), { status: 403 })
})
const json = (v: any, status = 200) => () => new Response(JSON.stringify(v), { status, headers: { 'content-type': 'application/json' } })

const settle = async () => {
  for (let i = 0; i < 4; i++) {
    await flushPromises()
    await nextTick()
  }
}

beforeEach(() => {
  observers = []
  seen = []
  answers = {}
  ;(window as any).IntersectionObserver = FakeIO
  vi.stubGlobal('fetch', fakeFetch)
  document.body.innerHTML = ''
})
afterEach(() => vi.unstubAllGlobals())

describe('the portals of a site', () => {
  it('load when in view: a progress, then the content with its data', async () => {
    answers['T1|'] = json({ template: '<p class="c">{{ name }}</p>', data: { name: 'Ana' } })
    document.body.innerHTML = '<p>before</p><rvq-portal source="T1" class="box" data-loading-text="Carregando"></rvq-portal>'
    site.mountAll()
    await settle()
    // not in view: nothing asked, the progress shown
    expect(seen.length).toBe(0)
    const bar = document.querySelector('.rvq-portal-progress')!
    expect(bar.getAttribute('role')).toBe('progressbar')
    expect(bar.getAttribute('aria-label')).toBe('Carregando')
    // the page's element replaced, its class kept, no token left in the DOM
    expect(document.querySelector('rvq-portal')).toBeNull()
    expect(document.querySelector('[data-rvq-portal]')!.className).toBe('box')
    expect(document.body.innerHTML).not.toContain('T1')

    enter()
    await settle()
    expect(seen.length).toBe(1)
    expect(seen[0].url).toBe(site.path)
    expect(seen[0].headers.get(site.headerToken)).toBe('T1')
    expect(document.querySelector('.c')!.textContent).toBe('Ana')
    expect(document.querySelector('.rvq-portal-progress')).toBeNull()
  })

  it('the data of a content is reactive, and it acts through its client', async () => {
    answers['T2|'] = json({
      template: '<button class="b" @click="send">{{ form.n }}</button>',
      data: { form: { n: 1 } }
    })
    answers['T2|submit'] = json({ ok: true })
    document.body.innerHTML = '<rvq-portal source="T2"></rvq-portal>'
    site.mountAll()
    enter()
    await settle()
    // a method of the content: its own data (form) and its client (portal)
    const btn = document.querySelector('.b') as HTMLButtonElement
    expect(btn.textContent).toBe('1')
    const app = (document.querySelector('[data-rvq-portal]') as any).__rvqPortal
    expect(app).toBeTruthy()

    const c = site.client('T2')
    const fd = new FormData()
    fd.set('name', 'Ana')
    const res = await c.action('submit', fd)
    expect(await res.json()).toEqual({ ok: true })
    const last = seen[seen.length - 1]
    expect(last.headers.get(site.headerAction)).toBe('submit')
    expect(last.body).toBe(fd)
    await c.action('draft', { a: 1 })
    expect(seen[seen.length - 1].headers.get('Content-Type')).toBe('application/json')
    expect(seen[seen.length - 1].body).toBe('{"a":1}')
  })

  it('an error says so, and tries again; an expired page asks to be reloaded', async () => {
    let fail = true
    answers['T3|'] = () => (fail ? new Response('x', { status: 500 }) : json({ template: '<i class="ok">ok</i>' })())
    answers['T4|'] = json({ error: 'expired' }, 410)
    document.body.innerHTML = '<rvq-portal source="T3" data-retry-text="De novo"></rvq-portal><rvq-portal source="T4"></rvq-portal>'
    site.mountAll()
    enter()
    await settle()
    const errs = document.querySelectorAll('.rvq-portal-error')
    expect(errs.length).toBe(2)
    expect(errs[1].textContent).toContain('out of date')
    fail = false
    ;(errs[0].querySelector('button') as HTMLButtonElement).click()
    await settle()
    expect(document.querySelector('.ok')).not.toBeNull()
  })

  it('several portals: each its own, each once', async () => {
    answers['A|'] = json({ template: '<b class="a">A</b>' })
    answers['B|'] = json({ template: '<b class="b">B</b>' })
    document.body.innerHTML = '<rvq-portal source="A"></rvq-portal><rvq-portal source="B"></rvq-portal><rvq-portal></rvq-portal>'
    expect(site.mountAll().length).toBe(2)
    expect(site.mountAll().length).toBe(0)
    enter()
    await settle()
    enter()
    await settle()
    expect(seen.length).toBe(2)
    expect(document.querySelector('.a')).not.toBeNull()
    expect(document.querySelector('.b')).not.toBeNull()
  })

  it('without IntersectionObserver: loads at once', async () => {
    delete (window as any).IntersectionObserver
    answers['T5|'] = json({ template: '<b class="x">x</b>' })
    document.body.innerHTML = '<rvq-portal source="T5"></rvq-portal>'
    site.mountAll()
    await settle()
    expect(document.querySelector('.x')).not.toBeNull()
  })

  it('a content of a kind: made by it — its template anew, its behavior a mixin', async () => {
    const mounted: string[] = []
    site.kinds['counter'] = (c, client) => ({
      // the configuration on the root of the template, read and taken off
      template: c.template.replace(' data-start="5"', ''),
      options: {
        data: { start: 5 },
        mixins: [{
          data: () => ({ n: 5 }),
          methods: { inc(this: any) { this.n++ } },
          mounted(this: any) { mounted.push(this.$el.tagName + ':' + (client ? 'client' : 'none')) }
        }]
      }
    })
    answers['K|'] = json({ kind: 'counter', template: '<div class="k" data-start="5"><button @click="inc">{{ n }}</button></div>' })
    document.body.innerHTML = '<rvq-portal source="K"></rvq-portal>'
    site.mountAll()
    enter()
    await settle()
    const k = document.querySelector('.k')!
    expect(k.getAttribute('data-start')).toBeNull()
    expect(k.textContent).toBe('5')
    ;(k.querySelector('button') as HTMLButtonElement).click()
    await settle()
    expect(k.textContent).toBe('6')
    expect(mounted).toEqual(['DIV:client'])
  })

  it('a kind not there: an error, nothing compiled', async () => {
    answers['K2|'] = json({ kind: 'nowhere', template: '<b class="no">x</b>' })
    document.body.innerHTML = '<rvq-portal source="K2"></rvq-portal>'
    site.mountAll()
    enter()
    await settle()
    expect(document.querySelector('.no')).toBeNull()
    expect(document.querySelector('.rvq-portal-error')).not.toBeNull()
  })

  it('the scripts of a content: loaded once for the page, before it shows', async () => {
    const appended: string[] = []
    const orig = document.head.appendChild.bind(document.head)
    vi.spyOn(document.head, 'appendChild').mockImplementation((el: any) => {
      if (el.tagName === 'SCRIPT') {
        appended.push(el.src)
        setTimeout(() => el.onload && el.onload(), 0)
        return el
      }
      return orig(el)
    })
    answers['S1|'] = json({ template: '<b class="s1">1</b>', scripts: ['/x/form.js'] })
    answers['S2|'] = json({ template: '<b class="s2">2</b>', scripts: ['/x/form.js'] })
    document.body.innerHTML = '<rvq-portal source="S1"></rvq-portal><rvq-portal source="S2"></rvq-portal>'
    site.mountAll()
    enter()
    await settle()
    await new Promise((r) => setTimeout(r, 5))
    await settle()
    expect(appended.filter((s) => s.endsWith('/x/form.js')).length).toBe(1)
    expect(document.querySelector('.s1')).not.toBeNull()
    expect(document.querySelector('.s2')).not.toBeNull()
    vi.restoreAllMocks()
  })

  it('a content given (a preview): shown as is, with no client', async () => {
    document.body.innerHTML = '<div id="p"></div>'
    site.mountContent(document.getElementById('p')!, { template: '<i class="pv">{{ portal === null ? "none" : "client" }}</i>' })
    await settle()
    expect(document.querySelector('.pv')!.textContent).toBe('none')
  })
})
