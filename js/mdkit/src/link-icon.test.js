import { describe, expect, it } from 'bun:test'
import { iconLinkHTML, iconLinkOf, validIcon } from './link-icon.js'
import { LinkForm } from './link-picker.js'

describe('validIcon', () => {
  it('is an icon of Material Design only', () => {
    expect(validIcon('mdi-file-pdf-box')).toBe('mdi-file-pdf-box')
    expect(validIcon(' mdi-file ')).toBe('mdi-file')
    expect(validIcon('javascript:x')).toBe('')
    expect(validIcon('mdi-"x')).toBe('')
  })
})

// a link of the Markdown editor with an icon: HTML, read back
describe('iconLinkHTML / iconLinkOf', () => {
  it('writes and reads it', () => {
    const html = iconLinkHTML('/system/files/a?x=1&y="2"', 'R&D <report>', 'mdi-file-pdf-box')
    expect(html).toBe('<a href="/system/files/a?x=1&amp;y=&quot;2&quot;" data-icon="mdi-file-pdf-box">R&amp;D &lt;report&gt;</a>')
    expect(iconLinkOf(html)).toEqual({ href: '/system/files/a?x=1&y="2"', icon: 'mdi-file-pdf-box', text: 'R&D <report>' })
    expect(iconLinkHTML('/a', 'a', '')).toBeNull()
    expect(iconLinkOf('[a](/a)')).toBeNull()
  })
})

// the icon of the form: a file chosen gives its own — in place of none or of
// the one a choice before put, not of one written
describe('the icon of the link form', () => {
  const form = (icon = '') => {
    const vm = { ...LinkForm.data.call({ href: '', target: '', text: '', icon }), iconName: icon, groups: [], q: '' }
    vm.validIcon = LinkForm.methods.validIcon
    vm.pick = LinkForm.methods.pick.bind(vm)
    return vm
  }
  it('takes the icon of the file chosen', async () => {
    const vm = form()
    await vm.pick({ target: '/system/files/a', title: 'a.pdf', icon: 'mdi-file-pdf-box' })
    expect(vm.iconName).toBe('mdi-file-pdf-box')
    await vm.pick({ target: '/system/files/b', title: 'b.xlsx', icon: 'mdi-file-excel' })
    expect(vm.iconName).toBe('mdi-file-excel')
    // a page: no icon of its own, the one there kept
    await vm.pick({ target: '/about', title: 'About' })
    expect(vm.iconName).toBe('mdi-file-excel')
  })
  it('keeps an icon written', async () => {
    const vm = form('mdi-star')
    await vm.pick({ target: '/system/files/a', title: 'a.pdf', icon: 'mdi-file-pdf-box' })
    expect(vm.iconName).toBe('mdi-star')
    expect(vm.picked).toBe('mdi-file-pdf-box')
  })
})
