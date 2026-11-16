// The icon of a link (an icon of Material Design, mdi-…, before its text):
// the attribute icon of the mark link of any tiptap editor (LinkIcon), kept
// in the HTML as data-icon, shown by the classes of the font of the icons
// (mdi mdi-…: the icon its ::before, the text as it is). A file chosen
// gives its own (the icon of its extension, the server's).
import { Extension } from '@tiptap/core'

const iconRe = /^mdi-[a-z0-9-]+$/

// validIcon is name when it is an icon (mdi-…), else ''.
export function validIcon(name) {
  const n = String(name || '').trim()
  return iconRe.test(n) ? n : ''
}

export const LinkIcon = Extension.create({
  name: 'linkIcon',
  addGlobalAttributes() {
    return [{
      types: ['link'],
      attributes: {
        icon: {
          default: null,
          parseHTML: (el) => validIcon(el.getAttribute('data-icon')) || null,
          renderHTML: (attrs) => (validIcon(attrs.icon) ? { 'data-icon': attrs.icon, class: 'mdi ' + attrs.icon } : {}),
        },
      },
    }]
  },
})

// iconLinkHTML is the link of the Markdown editor with an icon (Markdown has
// none: HTML), else null.
export function iconLinkHTML(href, label, icon) {
  const i = validIcon(icon)
  if (!i) return null
  const esc = (s) => String(s).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
  return `<a href="${esc(href)}" data-icon="${i}">${esc(label)}</a>`
}

// iconLinkOf reads a link of the Markdown editor with an icon ({href, text,
// icon}), else null.
export function iconLinkOf(text) {
  const m = /^<a href="([^"]*)" data-icon="(mdi-[a-z0-9-]+)">([^<]*)<\/a>$/.exec(text || '')
  if (!m) return null
  const un = (s) => s.replace(/&quot;/g, '"').replace(/&gt;/g, '>').replace(/&lt;/g, '<').replace(/&amp;/g, '&')
  return { href: un(m[1]), icon: m[2], text: un(m[3]) }
}
