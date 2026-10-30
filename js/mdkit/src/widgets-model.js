// The buttons and the tabs of a text: their value and their HTML, the same
// in the rich text (tiptap's nodes) and in Markdown (a block of HTML; a
// tab's content Markdown between blank lines), read back from either.
//   buttons: {align, items: [{href, target, title, description, color,
//            iconLeft, iconRight, newTab}]} —
//     <div data-type="buttons" data-align="…"><a data-type="button" href="…"
//     data-href="…" title="description" data-color="…" data-icon-left="mdi-…"
//     data-icon-right="mdi-…" target="_blank" rel="noopener">title</a>…</div>;
//     a title empty is the page's or the record's (target: the link cards'),
//     put by the page that shows it;
//   tabs: {align, items: [{title, color, iconLeft, iconRight, body}]} — body
//     the content (HTML, or Markdown) —
//     <div data-type="tabs" data-align="…"><div data-type="tab"
//     data-title="…" data-color="…" data-icon-left="…" data-icon-right="…">
//     …</div>…</div>.

// COLORS are the colors of the theme an item may have ('' the default).
export const COLORS = ['', 'primary', 'secondary', 'success', 'info', 'warning', 'error']

// ALIGNS are the places of the items in their line.
export const ALIGNS = ['start', 'center', 'end', 'stretch']

const iconRe = /^mdi-[a-z0-9-]+$/
const esc = (s) => String(s == null ? '' : s).replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
const color = (c) => (COLORS.includes(c) ? c : '')
const align = (a) => (ALIGNS.includes(a) ? a : 'start')
const icon = (i) => (iconRe.test(i || '') ? i : '')

// attrs is the attributes of an item's look (color, icons), those it has.
function attrs(it) {
  let s = ''
  if (color(it.color)) s += ' data-color="' + color(it.color) + '"'
  if (icon(it.iconLeft)) s += ' data-icon-left="' + icon(it.iconLeft) + '"'
  if (icon(it.iconRight)) s += ' data-icon-right="' + icon(it.iconRight) + '"'
  return s
}

// look is an item's look read from el.
function look(el) {
  return {
    color: color(el.getAttribute('data-color') || ''),
    iconLeft: icon(el.getAttribute('data-icon-left') || ''),
    iconRight: icon(el.getAttribute('data-icon-right') || ''),
  }
}

// buttonsHTML is the HTML of buttons (one line: a block of Markdown too).
export function buttonsHTML(v) {
  const items = (v.items || []).filter((it) => it.href)
  return '<div data-type="buttons" data-align="' + align(v.align) + '">' + items.map((it) => {
    let a = '<a data-type="button" href="' + esc(it.href) + '"'
    if (it.target) a += ' data-href="' + esc(it.target) + '"'
    if (it.description) a += ' title="' + esc(it.description) + '"'
    a += attrs(it)
    if (it.newTab) a += ' target="_blank" rel="noopener"'
    return a + '>' + esc(it.title || '') + '</a>'
  }).join('') + '</div>'
}

// buttonsOf is the value of the buttons of el (div[data-type=buttons]).
export function buttonsOf(el) {
  return {
    align: align(el.getAttribute('data-align')),
    items: [...el.querySelectorAll('a')].map((a) => ({
      href: a.getAttribute('href') || '',
      target: a.getAttribute('data-href') || '',
      title: (a.textContent || '').trim(),
      description: a.getAttribute('title') || '',
      newTab: a.getAttribute('target') === '_blank',
      ...look(a),
    })),
  }
}

// tabsHTML is the HTML of tabs, each tab's body its HTML.
export function tabsHTML(v) {
  return '<div data-type="tabs" data-align="' + align(v.align) + '">' + (v.items || []).map((it) =>
    '<div data-type="tab" data-title="' + esc(it.title || '') + '"' + attrs(it) + '>' + (it.body || '') + '</div>',
  ).join('') + '</div>'
}

// tabsOf is the value of the tabs of el (div[data-type=tabs]), each body its
// HTML.
export function tabsOf(el) {
  return {
    align: align(el.getAttribute('data-align')),
    items: [...el.children].filter((c) => c.getAttribute('data-type') === 'tab').map((t) => ({
      title: t.getAttribute('data-title') || '',
      body: t.innerHTML,
      ...look(t),
    })),
  }
}

// tabsMarkdown is the Markdown of tabs, each body Markdown: its HTML lines
// around each body between blank lines (Markdown there).
export function tabsMarkdown(v) {
  let s = '<div data-type="tabs" data-align="' + align(v.align) + '">\n'
  for (const it of v.items || []) {
    const body = (it.body || '').trim()
    s += '<div data-type="tab" data-title="' + esc(it.title || '') + '"' + attrs(it) + '>\n\n' + (body ? body + '\n\n' : '') + '</div>\n'
  }
  return s + '</div>'
}

const unesc = (s) => s.replace(/&quot;/g, '"').replace(/&lt;/g, '<').replace(/&gt;/g, '>').replace(/&#39;/g, "'").replace(/&amp;/g, '&')

// attrOf is the attribute name of the opening tag tag (a string).
function attrOf(tag, name) {
  const m = new RegExp('\\s' + name + '="([^"]*)"').exec(tag)
  return m ? unesc(m[1]) : ''
}

const tagLook = (tag) => ({
  color: color(attrOf(tag, 'data-color')),
  iconLeft: icon(attrOf(tag, 'data-icon-left')),
  iconRight: icon(attrOf(tag, 'data-icon-right')),
})

const opens = (line) => (line.match(/<div\b/gi) || []).length
const closes = (line) => (line.match(/<\/div>/gi) || []).length

// blockAt is the block of buttons or tabs of the Markdown text at the line
// of offset — {kind, from, to (offsets of its lines), value} — or null.
export function blockAt(text, offset) {
  const lines = text.split('\n')
  const starts = []
  let o = 0
  for (const l of lines) {
    starts.push(o)
    o += l.length + 1
  }
  let at = starts.findIndex((s, i) => offset >= s && offset <= s + lines[i].length)
  if (at < 0) at = lines.length - 1
  // the opening line of a block at or before it, whose end is at or after it
  for (let i = at; i >= 0; i--) {
    const m = /^\s*<div data-type="(buttons|tabs)"[^>]*>/.exec(lines[i])
    if (!m) continue
    let depth = 0
    let j = i
    for (; j < lines.length; j++) {
      depth += opens(lines[j]) - closes(lines[j])
      if (depth <= 0) break
    }
    if (j < at) return null
    const from = starts[i]
    const to = starts[Math.min(j, lines.length - 1)] + lines[Math.min(j, lines.length - 1)].length
    const src = text.slice(from, to)
    return { kind: m[1], from, to, value: m[1] === 'buttons' ? buttonsOfMarkdown(src) : tabsOfMarkdown(src) }
  }
  return null
}

// buttonsOfMarkdown is the value of the buttons of their HTML (a string).
export function buttonsOfMarkdown(src) {
  const open = /<div data-type="buttons"[^>]*>/.exec(src)
  const items = []
  const re = /<a\b([^>]*)>([\s\S]*?)<\/a>/g
  let m
  while ((m = re.exec(src))) {
    const tag = '<a' + m[1] + '>'
    items.push({
      href: attrOf(tag, 'href'), target: attrOf(tag, 'data-href'), title: unesc(m[2].replace(/<[^>]*>/g, '')).trim(),
      description: attrOf(tag, 'title'), newTab: attrOf(tag, 'target') === '_blank', ...tagLook(tag),
    })
  }
  return { align: align(open ? attrOf(open[0], 'data-align') : ''), items }
}

// tabsOfMarkdown is the value of the tabs of their Markdown, each body
// Markdown.
export function tabsOfMarkdown(src) {
  const lines = src.split('\n')
  const open = /<div data-type="tabs"[^>]*>/.exec(lines[0] || '')
  const items = []
  let cur = null
  let depth = 0
  for (let i = 1; i < lines.length; i++) {
    const l = lines[i]
    const tab = depth === 0 && /^\s*<div data-type="tab"[^>]*>\s*$/.exec(l)
    if (tab) {
      const tag = l.trim()
      cur = { title: attrOf(tag, 'data-title'), body: [], ...tagLook(tag) }
      items.push(cur)
      depth = 1
      continue
    }
    if (cur && depth > 0) {
      depth += opens(l) - closes(l)
      if (depth <= 0) {
        // its closing line: what is before </div> is the body's
        const before = l.slice(0, l.toLowerCase().lastIndexOf('</div>'))
        if (before.trim()) cur.body.push(before)
        depth = 0
        cur = null
        continue
      }
      cur.body.push(l)
    }
  }
  return {
    align: align(open ? attrOf(open[0], 'data-align') : ''),
    items: items.map((it) => ({ ...it, body: it.body.join('\n').trim() })),
  }
}
