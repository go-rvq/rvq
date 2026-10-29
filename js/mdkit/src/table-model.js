// The table of the table builder (mdkit-table-builder): a grid of texts —
// each Markdown (bold, italic, code, links…; its lines broken by <br> in a
// table of GFM) —,
// its merged cells, whether its first row is the header, the alignment of
// each column; made the Markdown of a table — GitHub's (GFM) when it can be:
// a header, no merged cell —, else its HTML (colspan, rowspan); read back
// from either (the table under the cursor, edited again). Text only (no DOM).
//
// The widths of its columns: '' (automatic), pixels ('120') or a part of the
// table's ('25%'); with any, HTML (GFM has none). The table fluid — as wide as
// the text, the default — or not — as wide as its columns, placed left,
// center or right (its class: mdkit-fit mdkit-<position>; HTML too) —.
//
// A cell's content placed on its own: across (halign, over its column's
// align) and up and down (valign: top, middle, bottom); with any, HTML (align,
// valign on the cell).
//
//	{ rows, cols, cells: [[text]], merges: [{ r, c, rs, cs }], header, aligns: ['' | 'left' | 'center' | 'right'], widths: ['' | '120' | '25%'], fluid, position: 'left' | 'center' | 'right', halign: [['' | 'left' | 'center' | 'right']], valign: [['' | 'top' | 'middle' | 'bottom']] }

// create is a table of rows × cols, empty, its first row the header.
export function create(rows = 3, cols = 3) {
  return {
    rows, cols,
    cells: Array.from({ length: rows }, () => Array(cols).fill('')),
    merges: [], header: true, aligns: Array(cols).fill(''), widths: Array(cols).fill(''),
    fluid: true, position: 'left', halign: grid(rows, cols), valign: grid(rows, cols),
  }
}

// grid is rows × cols of ''.
const grid = (rows, cols) => Array.from({ length: rows }, () => Array(cols).fill(''))

const positions = ['left', 'center', 'right']

// tableClass is the class of a table not fluid (mdkit-fit and its position);
// '' fluid.
export const tableClass = (t) => (t.fluid === false ? 'mdkit-fit mdkit-' + (positions.includes(t.position) ? t.position : 'left') : '')

// layoutOf is the fluid and the position of the class of a table.
function layoutOf(cls) {
  const c = String(cls || '')
  if (!/\bmdkit-fit\b/.test(c)) return { fluid: true, position: 'left' }
  const m = /\bmdkit-(left|center|right)\b/.exec(c)
  return { fluid: false, position: m ? m[1] : 'left' }
}

// normalize gives a table what it may lack (its widths: a table of before),
// as many as its columns.
export function normalize(t) {
  const fit = (a) => Array.from({ length: t.cols }, (_, i) => (a && a[i]) || '')
  t.aligns = fit(t.aligns)
  t.widths = fit(t.widths)
  const fitGrid = (g) => Array.from({ length: t.rows }, (_, r) => fit(g && g[r]))
  t.halign = fitGrid(t.halign)
  t.valign = fitGrid(t.valign)
  if (t.fluid === undefined) t.fluid = true
  if (!positions.includes(t.position)) t.position = 'left'
  return t
}

// percent says the widths are parts of the table's ('25%').
export const percent = (t) => (t.widths || []).some((w) => /%$/.test(w))

// mergeAt is the merge that covers the cell (r, c); null, none.
export function mergeAt(t, r, c) {
  return t.merges.find((m) => r >= m.r && r < m.r + m.rs && c >= m.c && c < m.c + m.cs) || null
}

// covered says the cell (r, c) is under a merge, not its first: not drawn.
export function covered(t, r, c) {
  const m = mergeAt(t, r, c)
  return !!m && (m.r !== r || m.c !== c)
}

// clean drops the merges of a single cell.
function clean(t) {
  t.merges = t.merges.filter((m) => m.rs > 0 && m.cs > 0 && m.rs * m.cs > 1)
}

// insertRow puts an empty row at the index at (0..rows); a merge it falls
// within grows.
export function insertRow(t, at) {
  t.cells.splice(at, 0, Array(t.cols).fill(''))
  t.halign.splice(at, 0, Array(t.cols).fill(''))
  t.valign.splice(at, 0, Array(t.cols).fill(''))
  for (const m of t.merges) {
    if (m.r >= at) m.r++
    else if (at < m.r + m.rs) m.rs++
  }
  t.rows++
}

// removeRow takes the row at out (the last one stays); a merge it crossed
// shrinks, its text kept in the row below.
export function removeRow(t, at) {
  if (t.rows <= 1) return
  for (const m of t.merges) {
    if (m.r === at && m.rs > 1) {
      for (const g of [t.cells, t.halign, t.valign]) g[at + 1][m.c] = g[at][m.c]
    }
  }
  t.cells.splice(at, 1)
  t.halign.splice(at, 1)
  t.valign.splice(at, 1)
  for (const m of t.merges) {
    if (m.r > at) m.r--
    else if (at < m.r + m.rs) m.rs--
  }
  t.rows--
  clean(t)
}

// insertCol puts an empty column at the index at (0..cols).
export function insertCol(t, at) {
  for (const row of [...t.cells, ...t.halign, ...t.valign]) row.splice(at, 0, '')
  t.aligns.splice(at, 0, '')
  t.widths.splice(at, 0, '')
  for (const m of t.merges) {
    if (m.c >= at) m.c++
    else if (at < m.c + m.cs) m.cs++
  }
  t.cols++
}

// removeCol takes the column at out (the last one stays).
export function removeCol(t, at) {
  if (t.cols <= 1) return
  for (const m of t.merges) {
    if (m.c === at && m.cs > 1) {
      for (const g of [t.cells, t.halign, t.valign]) g[m.r][at + 1] = g[m.r][at]
    }
  }
  for (const row of [...t.cells, ...t.halign, ...t.valign]) row.splice(at, 1)
  t.aligns.splice(at, 1)
  t.widths.splice(at, 1)
  for (const m of t.merges) {
    if (m.c > at) m.c--
    else if (at < m.c + m.cs) m.cs--
  }
  t.cols--
  clean(t)
}

// span is the range (r1, c1)–(r2, c2), in order, grown to hold whole every
// merge it touches (as a spreadsheet's selection).
export function span(t, r1, c1, r2, c2) {
  let s = { r1: Math.min(r1, r2), c1: Math.min(c1, c2), r2: Math.max(r1, r2), c2: Math.max(c1, c2) }
  for (let grown = true; grown;) {
    grown = false
    for (const m of t.merges) {
      const inside = m.r <= s.r2 && m.r + m.rs - 1 >= s.r1 && m.c <= s.c2 && m.c + m.cs - 1 >= s.c1
      if (!inside) continue
      const n = {
        r1: Math.min(s.r1, m.r), c1: Math.min(s.c1, m.c),
        r2: Math.max(s.r2, m.r + m.rs - 1), c2: Math.max(s.c2, m.c + m.cs - 1),
      }
      if (n.r1 !== s.r1 || n.c1 !== s.c1 || n.r2 !== s.r2 || n.c2 !== s.c2) {
        s = n
        grown = true
      }
    }
  }
  return s
}

// merge makes the range one cell — the merges within it gone into it —, the
// texts of its cells its text (joined by a space); false: a single cell.
export function merge(t, r1, c1, r2, c2) {
  const s = span(t, r1, c1, r2, c2)
  if (s.r1 === s.r2 && s.c1 === s.c2) return false
  const texts = []
  for (let r = s.r1; r <= s.r2; r++) {
    for (let c = s.c1; c <= s.c2; c++) {
      if (!covered(t, r, c) && t.cells[r][c].trim()) texts.push(t.cells[r][c].trim())
      t.cells[r][c] = ''
      if (r !== s.r1 || c !== s.c1) t.halign[r][c] = t.valign[r][c] = ''
    }
  }
  t.merges = t.merges.filter((m) => !(m.r >= s.r1 && m.r <= s.r2 && m.c >= s.c1 && m.c <= s.c2))
  t.merges.push({ r: s.r1, c: s.c1, rs: s.r2 - s.r1 + 1, cs: s.c2 - s.c1 + 1 })
  t.cells[s.r1][s.c1] = texts.join(' ')
  return true
}

// unmerge makes the merge over (r, c) its cells again.
export function unmerge(t, r, c) {
  const m = mergeAt(t, r, c)
  if (m) t.merges = t.merges.filter((x) => x !== m)
  return !!m
}

// placed says a cell has its content placed on its own (halign, valign).
const placed = (t) => [...(t.halign || []), ...(t.valign || [])].some((row) => row.some((a) => a))

// gfm says the table is Markdown's (GFM): a header, no merged cell, no width,
// no cell placed on its own.
export const gfm = (t) => t.header && t.merges.length === 0 && !(t.widths || []).some((w) => w) && t.fluid !== false && !placed(t)

// cellAlign is where the content of the cell (r, c) goes across: its own, else
// its column's ('' the default).
export const cellAlign = (t, r, c) => (t.halign && t.halign[r] && t.halign[r][c]) || t.aligns[c] || ''

// cellVAlign is where the content of the cell (r, c) goes up and down ('' the
// default).
export const cellVAlign = (t, r, c) => (t.valign && t.valign[r] && t.valign[r][c]) || ''

// pipeCell is a cell of a line of GFM: its pipes escaped, its lines broken
// by <br>
const pipeCell = (s) => ' ' + (s || '').trim().replace(/\|/g, '\\|').replace(/[ \t]*\n[ \t]*/g, '<br>') + ' '
const delimiter = { '': '---', left: ':--', center: ':-:', right: '--:' }
// markdownRe is a text of Markdown (not plain): its marks, a line break
const markdownRe = /[*_`~[\]!\\#<>\n]|^\s*([-+]|\d+\.)\s/

// toMarkdown is the table as its text: GFM, when it can be; else HTML.
export function toMarkdown(t) {
  if (!gfm(t)) return toHTML(t)
  const line = (cells) => '|' + cells.map(pipeCell).join('|') + '|'
  return [line(t.cells[0]), '|' + t.aligns.map((a) => ' ' + delimiter[a || ''] + ' ').join('|') + '|',
    ...t.cells.slice(1).map(line)].join('\n')
}

// toHTML is the table as HTML: its header (thead, th), its merged cells
// (colspan, rowspan), its columns' alignment (align) and widths (col width;
// in parts, the table the whole width); each cell's Markdown in
// it, between blank lines (rendered by the server, then sanitized).
export function toHTML(t) {
  const widths = t.widths || []
  const cls = tableClass(t)
  const out = ['<table' + (cls ? ` class="${cls}"` : '') + (percent(t) && !cls ? ' width="100%"' : '') + '>']
  // the widths of the columns
  if (widths.some((w) => w)) {
    out.push('<colgroup>' + Array.from({ length: t.cols }, (_, c) => (widths[c] ? `<col width="${widths[c]}">` : '<col>')).join('') + '</colgroup>')
  }
  const row = (r, tag) => {
    const cells = []
    for (let c = 0; c < t.cols; c++) {
      if (covered(t, r, c)) continue
      const m = mergeAt(t, r, c)
      let attrs = ''
      if (m && m.cs > 1) attrs += ` colspan="${m.cs}"`
      if (m && m.rs > 1) attrs += ` rowspan="${m.rs}"`
      if (cellAlign(t, r, c)) attrs += ` align="${cellAlign(t, r, c)}"`
      if (cellVAlign(t, r, c)) attrs += ` valign="${cellVAlign(t, r, c)}"`
      // its Markdown between blank lines: read as Markdown within the HTML
      // (CommonMark: a block of HTML ends at a blank line); plain text, as it is
      const text = (t.cells[r][c] || '').trim()
      cells.push(markdownRe.test(text) ? `<${tag}${attrs}>\n\n${text}\n\n</${tag}>` : `<${tag}${attrs}>${text}</${tag}>`)
    }
    return '<tr>' + cells.join('') + '</tr>'
  }
  let from = 0
  if (t.header) {
    out.push('<thead>', row(0, 'th'), '</thead>')
    from = 1
  }
  if (from < t.rows) {
    out.push('<tbody>')
    for (let r = from; r < t.rows; r++) out.push(row(r, 'td'))
    out.push('</tbody>')
  }
  out.push('</table>')
  return out.join('\n')
}

// splitRow is the cells of a line of a GFM table (its \| a pipe).
function splitRow(line) {
  let s = line.trim()
  if (s.startsWith('|')) s = s.slice(1)
  if (s.endsWith('|') && !s.endsWith('\\|')) s = s.slice(0, -1)
  const cells = []
  let cur = ''
  for (let i = 0; i < s.length; i++) {
    if (s[i] === '\\' && s[i + 1] === '|') {
      cur += '|'
      i++
    } else if (s[i] === '|') {
      cells.push(cur.trim())
      cur = ''
    } else cur += s[i]
  }
  cells.push(cur.trim())
  return cells
}

export const delimiterRe = /^\s*\|?\s*:?-{1,}:?\s*(\|\s*:?-{1,}:?\s*)*\|?\s*$/

// fromMarkdown is the table of the lines of a GFM table; null, not one.
export function fromMarkdown(lines) {
  if (lines.length < 2 || !delimiterRe.test(lines[1])) return null
  const head = splitRow(lines[0])
  const cols = head.length
  const aligns = splitRow(lines[1]).slice(0, cols).map((d) => {
    const l = d.startsWith(':')
    const r = d.endsWith(':')
    return l && r ? 'center' : r ? 'right' : l ? 'left' : ''
  })
  while (aligns.length < cols) aligns.push('')
  const fit = (cells) => Array.from({ length: cols }, (_, i) => cells[i] || '')
  const br = (cells) => cells.map((x) => x.replace(/<br\s*\/?>/gi, '\n'))
  const cells = [head, ...lines.slice(2).map(splitRow)].map(fit).map(br)
  return {
    rows: cells.length, cols, cells, merges: [], header: true, aligns, widths: Array(cols).fill(''), fluid: true, position: 'left',
    halign: grid(cells.length, cols), valign: grid(cells.length, cols),
  }
}

// inlineMarkdown is the content of a cell of HTML as Markdown: the table
// builder's (Markdown between blank lines) as it is; another's, its tags
// made Markdown (strong, em, del, code, a, br, p), the others dropped.
function inlineMarkdown(html) {
  const s = html.trim()
  if (!/<[a-z!/]/i.test(s)) return decode(s)
  return decode(s
    .replace(/<br\s*\/?>/gi, '\n')
    .replace(/<\/p>\s*<p[^>]*>/gi, '\n\n')
    .replace(/<(strong|b)\b[^>]*>([\s\S]*?)<\/\1>/gi, '**$2**')
    .replace(/<(em|i)\b[^>]*>([\s\S]*?)<\/\1>/gi, '_$2_')
    .replace(/<(del|s|strike)\b[^>]*>([\s\S]*?)<\/\1>/gi, '~~$2~~')
    .replace(/<code\b[^>]*>([\s\S]*?)<\/code>/gi, '`$1`')
    .replace(/<a\b[^>]*href\s*=\s*"([^"]*)"[^>]*>([\s\S]*?)<\/a>/gi, '[$2]($1)')
    .replace(/<[^>]*>/g, ''))
    .trim()
}

// decode is text with its entities of HTML made characters.
const decode = (s) => s.replace(/&nbsp;/g, ' ').replace(/&lt;/g, '<').replace(/&gt;/g, '>').replace(/&quot;/g, '"').replace(/&#39;/g, "'").replace(/&amp;/g, '&')

// fromHTML is the table of the HTML of a table: its rows, its cells (th, td)
// in place — their spans the merges —, its header (a first row of th), its
// columns' alignment (the align all its cells have), its cells' own (align —
// or a style text-align — other than their column's; valign, or a style
// vertical-align); null, no row.
export function fromHTML(html) {
  const rowsHTML = [...html.matchAll(/<tr\b[^>]*>([\s\S]*?)<\/tr>/gi)].map((m) => m[1])
  if (!rowsHTML.length) return null
  const grid = []
  const merges = []
  const hs = []
  const vs = []
  let header = false
  rowsHTML.forEach((rowHTML, r) => {
    grid[r] = grid[r] || []
    let c = 0
    const cells = [...rowHTML.matchAll(/<(t[hd])\b([^>]*)>([\s\S]*?)<\/\1>/gi)]
    if (r === 0) header = cells.length > 0 && cells.every((m) => m[1].toLowerCase() === 'th')
    for (const [, , attrs, inner] of cells) {
      while (grid[r][c] !== undefined) c++
      const num = (name) => Math.max(1, parseInt((new RegExp(name + '\\s*=\\s*"?(\\d+)', 'i').exec(attrs) || [])[1] || '1', 10))
      const cs = num('colspan')
      const rs = num('rowspan')
      const align = (/(?:\balign\s*=\s*"?|text-align\s*:\s*)(left|center|right)/i.exec(attrs) || [])[1] || ''
      const valign = (/(?:\bvalign\s*=\s*"?|vertical-align\s*:\s*)(top|middle|bottom)/i.exec(attrs) || [])[1] || ''
      hs.push({ r, c, a: align.toLowerCase() })
      vs.push({ r, c, a: valign.toLowerCase() })
      for (let i = 0; i < rs; i++) {
        grid[r + i] = grid[r + i] || []
        for (let j = 0; j < cs; j++) grid[r + i][c + j] = i === 0 && j === 0 ? inlineMarkdown(inner) : ''
      }
      if (cs > 1 || rs > 1) merges.push({ r, c, rs, cs })
      c += cs
    }
  })
  const rows = grid.length
  const cols = Math.max(...grid.map((row) => row.length))
  // the widths of its columns (col width), in order; else its cells'
  // (colwidth: the rich text editor's — tiptap — table, in pixels)
  const widths = []
  for (const [, attrs] of html.matchAll(/<col\b([^>]*)>/gi)) {
    const w = (/\bwidth\s*=\s*"?(\d+%?)/i.exec(attrs) || [])[1] || ''
    const span = Math.max(1, parseInt((/span\s*=\s*"?(\d+)/i.exec(attrs) || [])[1] || '1', 10))
    for (let i = 0; i < span; i++) widths.push(w)
  }
  if (!widths.some((w) => w)) {
    widths.length = 0
    for (const [, attrs] of (rowsHTML[0] || '').matchAll(/<t[hd]\b([^>]*)>/gi)) {
      const cw = (/colwidth\s*=\s*"([\d,]+)"/i.exec(attrs) || [])[1]
      const span = Math.max(1, parseInt((/colspan\s*=\s*"?(\d+)/i.exec(attrs) || [])[1] || '1', 10))
      const ws = cw ? cw.split(',') : []
      for (let i = 0; i < span; i++) widths.push(ws[i] && ws[i] !== '0' ? ws[i] : '')
    }
  }
  const cells = grid.map((row) => Array.from({ length: cols }, (_, i) => row[i] ?? ''))
  // a column's alignment: the one all its cells have; a cell's other, its own
  const aligns = Array.from({ length: cols }, (_, c) => {
    const of = hs.filter((h) => h.c === c).map((h) => h.a)
    return of.length && of.every((a) => a === of[0]) ? of[0] : ''
  })
  const halign = Array.from({ length: rows }, () => Array(cols).fill(''))
  const valign = Array.from({ length: rows }, () => Array(cols).fill(''))
  for (const h of hs) if (h.a !== aligns[h.c]) halign[h.r][h.c] = h.a || 'left'
  for (const v of vs) valign[v.r][v.c] = v.a
  return {
    ...layoutOf((/<table\b[^>]*\bclass\s*=\s*"([^"]*)"/i.exec(html) || [])[1]),
    rows, cols, cells, merges, header,
    aligns,
    widths: Array.from({ length: cols }, (_, i) => widths[i] || ''),
    halign, valign,
  }
}

// tableAt is the table of the lines around the line at (an index): its first
// and last line and the table — of GFM, or of HTML (from <table to
// </table>) —; null, the line is not of a table.
export function tableAt(lines, at) {
  // HTML: a <table> open above (or here) not closed before the line
  for (let i = at; i >= 0; i--) {
    if (/<\/table>/i.test(lines[i]) && i < at) break
    if (/<table\b/i.test(lines[i])) {
      let j = i
      while (j < lines.length && !/<\/table>/i.test(lines[j])) j++
      if (j >= lines.length || j < at) break
      const t = fromHTML(lines.slice(i, j + 1).join('\n'))
      return t ? { from: i, to: j, table: t } : null
    }
  }
  const isRow = (l) => l !== undefined && /^\s*\|/.test(l)
  if (!isRow(lines[at])) return null
  let from = at
  let to = at
  while (isRow(lines[from - 1])) from--
  while (isRow(lines[to + 1])) to++
  const t = fromMarkdown(lines.slice(from, to + 1))
  return t ? { from, to, table: t } : null
}

// toRichHTML is the table as the rich text editor's (tiptap) HTML: its cells'
// text (Markdown) made HTML by inline (a paragraph each), the header's th,
// the merged cells' spans, the widths in pixels of the columns on the cells
// (colwidth) — a part of the table, none: the table is fluid, a column as
// wide as it takes —, where a cell's content goes on the cell (align,
// valign: the editor's cell attributes, the plugin tables').
export function toRichHTML(t, inline) {
  const widths = (t.widths || []).map((w) => (/%$/.test(w) ? '' : w || ''))
  const esc = (s) => String(s).replace(/&/g, '&amp;').replace(/"/g, '&quot;')
  const rows = []
  for (let r = 0; r < t.rows; r++) {
    const cells = []
    for (let c = 0; c < t.cols; c++) {
      if (covered(t, r, c)) continue
      const m = mergeAt(t, r, c)
      const cs = m ? m.cs : 1
      const rs = m ? m.rs : 1
      const tag = t.header && r === 0 ? 'th' : 'td'
      const cw = widths.slice(c, c + cs)
      let attrs = ` colspan="${cs}" rowspan="${rs}"`
      if (cw.some((w) => w)) attrs += ` colwidth="${esc(cw.map((w) => w || '0').join(','))}"`
      if (cellAlign(t, r, c)) attrs += ` align="${cellAlign(t, r, c)}"`
      if (cellVAlign(t, r, c)) attrs += ` valign="${cellVAlign(t, r, c)}"`
      const lines = String(t.cells[r][c] || '').split('\n').map((l) => inline(l)).join('<br>')
      cells.push(`<${tag}${attrs}><p>${lines}</p></${tag}>`)
    }
    rows.push('<tr>' + cells.join('') + '</tr>')
  }
  const cls = tableClass(t)
  return '<table' + (cls ? ` class="${cls}"` : '') + '><tbody>' + rows.join('') + '</tbody></table>'
}

// separators are the separators of columns of a CSV a paste may have, the
// first the default.
export const separators = ['\t', ';', ',', '|']

// detectSeparator is the separator of the columns of text: the one of
// separators (outside quotes) found the most times in its first line,
// the same in every line; else a comma.
export function detectSeparator(text) {
  const lines = String(text).split(/\r?\n/).filter((l) => l.trim()).slice(0, 10)
  let best = ','
  let most = 0
  for (const sep of separators) {
    const counts = lines.map((l) => parseCSV(l, sep)[0]?.length - 1 || 0)
    if (counts[0] > most && counts.every((n) => n === counts[0])) {
      best = sep
      most = counts[0]
    }
  }
  if (most) return best
  // lines of different lengths: the most found in the first
  for (const sep of separators) {
    const n = (parseCSV(lines[0] || '', sep)[0]?.length || 1) - 1
    if (n > most) {
      best = sep
      most = n
    }
  }
  return best
}

// parseCSV is the rows of the cells of text, its columns separated by sep
// (none: detected): a cell between double quotes may hold sep, line breaks
// and "" (a quote); the blank lines at the end dropped.
export function parseCSV(text, sep) {
  const s = String(text).replace(/\r\n?/g, '\n')
  if (!sep) sep = detectSeparator(s)
  const rows = []
  let row = []
  let cell = ''
  let quoted = false
  for (let i = 0; i < s.length; i++) {
    const ch = s[i]
    if (quoted) {
      if (ch === '"' && s[i + 1] === '"') {
        cell += '"'
        i++
      } else if (ch === '"') quoted = false
      else cell += ch
    } else if (ch === '"' && cell === '') quoted = true
    else if (s.startsWith(sep, i)) {
      row.push(cell)
      cell = ''
      i += sep.length - 1
    } else if (ch === '\n') {
      row.push(cell)
      rows.push(row)
      row = []
      cell = ''
    } else cell += ch
  }
  row.push(cell)
  rows.push(row)
  while (rows.length && rows[rows.length - 1].every((c) => !c.trim())) rows.pop()
  return rows.map((r) => r.map((c) => c.trim()))
}

// fromGrid is the table of rows of texts: as many columns as its longest,
// its first row the header.
export function fromGrid(grid) {
  const cols = Math.max(1, ...grid.map((r) => r.length))
  const t = create(Math.max(1, grid.length), cols)
  grid.forEach((r, i) => r.forEach((c, j) => { t.cells[i][j] = c }))
  return t
}

// fromClipboard is the table of what was copied: the HTML of a table (a
// spreadsheet's, as Excel's, a page's) when there is one, else the text as a
// CSV of sep (none: detected); null, nothing.
export function fromClipboard(html, text, sep) {
  if (html && /<table\b/i.test(html)) {
    const t = fromHTML(html.slice(html.search(/<table\b/i)))
    if (t) return normalize(t)
  }
  const grid = parseCSV(text || '', sep)
  return grid.length ? fromGrid(grid) : null
}

// paste puts the table src in t from the cell (r, c) — t grown to hold it —:
// its texts, its cells' places, its merges (the merges of t it touches
// undone); the rest of t as it was.
export function paste(t, src, r = 0, c = 0) {
  normalize(src)
  while (t.rows < r + src.rows) insertRow(t, t.rows)
  while (t.cols < c + src.cols) insertCol(t, t.cols)
  const r2 = r + src.rows - 1
  const c2 = c + src.cols - 1
  t.merges = t.merges.filter((m) => m.r > r2 || m.r + m.rs - 1 < r || m.c > c2 || m.c + m.cs - 1 < c)
  for (let i = 0; i < src.rows; i++) {
    for (let j = 0; j < src.cols; j++) {
      t.cells[r + i][c + j] = src.cells[i][j] || ''
      t.halign[r + i][c + j] = src.halign[i][j] || ''
      t.valign[r + i][c + j] = src.valign[i][j] || ''
    }
  }
  for (const m of src.merges) t.merges.push({ r: m.r + r, c: m.c + c, rs: m.rs, cs: m.cs })
  return t
}
