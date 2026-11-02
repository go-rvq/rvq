import { describe, expect, it, vi } from 'vitest'
import { plaid } from '../builder'

// The server out of reach: no alert() (its dialog stops the page), the
// page's message bar.
describe('fetch error', () => {
  it('tells it in the message bar, never by alert', async () => {
    const alert = vi.fn()
    ;(globalThis as any).alert = alert
    ;(globalThis as any).window.alert = alert
    global.fetch = vi.fn().mockImplementation(() => Promise.reject(new TypeError('network error')))
    const vars: any = {}
    await plaid().vars(vars).url('/a').eventFunc('x').go()
    expect(alert).not.toHaveBeenCalled()
    expect(vars.presetsMessage?.show).toBe(true)
    expect(vars.presetsMessage?.color).toBe('error')
  })
})
