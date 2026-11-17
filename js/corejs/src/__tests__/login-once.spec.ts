import { describe, expect, it, vi } from 'vitest'
import { plaid } from '../builder'

declare let window: any

// The session gone under several requests at once (the widgets of a
// dashboard reloading): one login, the others waiting for it and run again
// once it is back.
describe('login once', () => {
  it('asks the login once and runs every request again after it', async () => {
    window.__goplaid = window.__goplaid || { portals: {}, index: 0 }
    let session = false
    let logins = 0
    const calls: string[] = []
    let loginScope: any = null
    global.fetch = vi.fn().mockImplementation((url: string) => {
      calls.push(url)
      if (url.includes('__execute_event__=login_dialog')) {
        logins++
        return Promise.resolve({ ok: true, status: 200, redirected: false, json: () => Promise.resolve({}) })
      }
      if (!session) {
        return Promise.resolve({
          ok: false,
          status: 401,
          json: () => Promise.resolve({ loginURI: '/!login/dialog?__execute_event__=login_dialog' })
        })
      }
      return Promise.resolve({ ok: true, status: 200, redirected: false, json: () => Promise.resolve({}) })
    })

    const origScope = (plaid() as any).constructor.prototype.scope
    ;(plaid() as any).constructor.prototype.scope = function (v: any) {
      if (v && v.onLoginSuccess) {
        loginScope = v
      }
      return origScope.call(this, v)
    }

    await Promise.all([plaid().url('/a').eventFunc('w1').go(), plaid().url('/a').eventFunc('w2').go(), plaid().url('/a').eventFunc('w3').go()])
    await new Promise((r) => setTimeout(r, 0))
    expect(logins).toBe(1)
    expect(window.__goplaid.loginPending).toBe(true)

    session = true
    loginScope.onLoginSuccess()
    await new Promise((r) => setTimeout(r, 0))
    expect(window.__goplaid.loginPending).toBe(false)
    for (const w of ['w1', 'w2', 'w3']) {
      expect(calls.filter((u) => u.includes('__execute_event__=' + w)).length).toBe(2)
    }
    ;(plaid() as any).constructor.prototype.scope = origScope
  })
})
