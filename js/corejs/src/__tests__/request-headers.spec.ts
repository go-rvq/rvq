import { describe, expect, it, vi } from 'vitest'
import { plaid, requestHeadersOf } from '../builder'

// The headers the page asks every request to carry (vars.__requestHeaders):
// a value, or a function read at each request; an empty one not sent.
describe('request headers', () => {
  it('reads values and functions, leaves out the empty ones', () => {
    let n = 0
    const vars: any = { __requestHeaders: { 'X-A': 'a', 'X-B': () => String(++n), 'X-C': () => '', 'X-D': null } }
    expect(requestHeadersOf(vars)).toEqual({ 'X-A': 'a', 'X-B': '1' })
    expect(requestHeadersOf(vars)['X-B']).toBe('2')
    expect(requestHeadersOf(undefined)).toEqual({})
  })

  it('goes in every request, with X-Plaid-Request', async () => {
    const fetch = vi.fn().mockImplementation(() => Promise.resolve(new Response('{}', { headers: { 'Content-Type': 'application/json' } })))
    global.fetch = fetch
    const vars: any = { __requestHeaders: { 'X-Rvq-Menu-Records': () => '{"k":"v"}' } }
    await plaid().vars(vars).url('/a').eventFunc('x').go()
    const headers = fetch.mock.calls[0][1].headers
    expect(headers['X-Rvq-Menu-Records']).toBe('{"k":"v"}')
    expect(headers['X-Plaid-Request']).toBe('1')
  })
})
