import { describe, expect, it } from 'bun:test'
import { externalHref, metaTitle } from './link-picker.js'

// an address of another site: http or https, not of this origin
describe('externalHref', () => {
  it('tells the addresses of other sites', () => {
    const o = 'https://site.test'
    expect(externalHref('https://example.com/a', o)).toBe(true)
    expect(externalHref('http://example.com', o)).toBe(true)
    expect(externalHref('https://site.test/about', o)).toBe(false)
    expect(externalHref('/en-us/about', o)).toBe(false)
    expect(externalHref('mailto:a@b.c', o)).toBe(false)
    expect(externalHref('admin:posts:<1>:', o)).toBe(false)
  })
})

// the title asked of the page's meta: its answer's, '' of none or of an error
describe('metaTitle', () => {
  const doc = (url) => ({ querySelector: () => (url ? { getAttribute: (k) => (k === 'data-title' ? url : null) } : null) })
  it('is none without the meta', () => {
    expect(metaTitle(doc(''))).toBe(null)
  })
  it('asks the server, once', async () => {
    const calls = []
    globalThis.fetch = async (url, opts) => {
      calls.push([url, opts.body.get('url')])
      return { ok: true, json: async () => ({ data: { title: 'Contato' } }) }
    }
    const ask = metaTitle(doc('/admin?__execute_event__=presets_LinkTitle'))
    expect(await ask('https://example.com/c')).toBe('Contato')
    expect(calls).toEqual([['/admin?__execute_event__=presets_LinkTitle', 'https://example.com/c']])
  })
  it('is "" of an error or of no answer in time', async () => {
    globalThis.fetch = async () => ({ ok: false })
    expect(await metaTitle(doc('/t'))('https://x.y')).toBe('')
    globalThis.fetch = (_, opts) => new Promise((_, reject) => opts.signal.addEventListener('abort', () => reject(new Error('aborted'))))
    expect(await metaTitle(doc('/t'), 20)('https://x.y')).toBe('')
  })
})
