// The link cards of a text (rvq mdkit link-cards): the links marked made
// cards — their targets asked at once, each once —, their texts as text;
// one the server does not answer, or that links another site, stays.
import {describe, expect, it, vi} from 'vitest'
import {linkCards} from '@rvq/mdkit'

describe('link cards', () => {
  it('draws the cards the server answers, leaves the others', async () => {
    document.body.innerHTML =
      '<div class="mdkit-link-card" data-href="/en-us/about"><a href="/en-us/about">/en-us/about</a></div>' +
      '<div class="mdkit-link-card" data-href="/en-us/about"><a href="/en-us/about">/en-us/about</a></div>' +
      '<div class="mdkit-link-card" data-href="admin:posts:&lt;1&gt;:">admin:posts:&lt;1&gt;:</div>' +
      '<div class="mdkit-link-card" data-href="/gone"><a href="/gone">/gone</a></div>' +
      '<div class="mdkit-link-card" data-href="/evil"><a href="/evil">/evil</a></div>'
    const resolve = vi.fn(async (targets: string[]) => ({
      '/en-us/about': {href: '/en-us/about', title: 'About <b>us</b>', cover: '/system/c.jpg', summary: 'Who we are'},
      'admin:posts:<1>:': {href: '/admin/posts/1', title: 'A post'},
      '/evil': {href: '//evil.com/', title: 'Evil'},
      targets
    }))
    await linkCards(document.body, resolve)
    expect(resolve).toHaveBeenCalledTimes(1)
    expect(resolve.mock.calls[0][0]).toEqual(['/en-us/about', 'admin:posts:<1>:', '/gone', '/evil'])
    const cards = [...document.querySelectorAll('.mdkit-link-card')]
    expect(cards.map((c) => c.getAttribute('data-card'))).toEqual(['drawn', 'drawn', 'drawn', 'none', 'none'])
    const a = cards[0].querySelector('a.mdkit-card')!
    expect(a.getAttribute('href')).toBe('/en-us/about')
    expect(a.querySelector('.mdkit-card-title')!.textContent).toBe('About <b>us</b>')
    expect(a.querySelector('b')).toBeNull()
    expect(a.querySelector('img')!.getAttribute('src')).toBe('/system/c.jpg')
    expect(a.querySelector('.mdkit-card-summary')!.textContent).toBe('Who we are')
    expect(cards[2].querySelector('a')!.getAttribute('href')).toBe('/admin/posts/1')
    expect(cards[3].innerHTML).toBe('<a href="/gone">/gone</a>')
    expect(cards[4].querySelector('a')!.getAttribute('href')).toBe('/evil')

    // drawn once: asked again, nothing
    await linkCards(document.body, resolve)
    expect(resolve).toHaveBeenCalledTimes(1)
  })
})
